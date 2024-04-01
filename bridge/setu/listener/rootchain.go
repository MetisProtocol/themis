package listener

import (
	"context"
	"math/big"
	"os"
	"strconv"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethCommon "github.com/ethereum/go-ethereum/common"

	"github.com/metis-seq/themis/bridge/setu/util"
	chainmanagerTypes "github.com/metis-seq/themis/chainmanager/types"
	"github.com/metis-seq/themis/helper"
)

// RootChainListenerContext root chain listener context
type RootChainListenerContext struct {
	ChainmanagerParams *chainmanagerTypes.Params
}

// RootChainListener - Listens to and process events from rootchain
type RootChainListener struct {
	BaseListener
	// ABIs
	abis []*abi.ABI

	stakingInfoAbi *abi.ABI
}

const (
	lastRootBlockKey = "listener-rootchain-last-block" // storage key
)

// NewRootChainListener - constructor func
func NewRootChainListener() *RootChainListener {
	contractCaller, err := helper.NewContractCaller()
	if err != nil {
		panic(err)
	}

	abis := []*abi.ABI{
		&contractCaller.StakingInfoABI,
	}

	return &RootChainListener{
		abis:           abis,
		stakingInfoAbi: &contractCaller.StakingInfoABI,
	}
}

// Start starts new block subscription
func (rl *RootChainListener) Start() error {
	rl.Logger.Info("Starting")

	// create cancellable context
	ctx, cancelSubscription := context.WithCancel(context.Background())
	rl.cancelSubscription = cancelSubscription

	// create cancellable context
	headerCtx, cancelHeaderProcess := context.WithCancel(context.Background())
	rl.cancelHeaderProcess = cancelHeaderProcess

	// start header process
	go rl.StartHeaderProcess(headerCtx)

	// start go routine to poll for new header using client object
	rl.Logger.Info("Start polling for rootchain header blocks", "pollInterval", helper.GetConfig().RootPollInterval)

	// wait mpc generate first
	for {
		time.Sleep(5 * time.Second)
		if !util.MpcCommonGenerated || !util.MpcStateCommitGenerated {
			break
		}
		rl.Logger.Info("Start polling for rootchain header blocks, wait mpc generate first")
	}

	// start polling for the finalized block in main chain (available post-merge)
	go rl.StartPolling(ctx, helper.GetConfig().RootPollInterval, nil)

	return nil
}

// ProcessHeader - process headerblock from rootchain
func (rl *RootChainListener) ProcessHeader(newHeader *blockHeader) {
	rl.Logger.Debug("New block detected", "blockNumber", newHeader.header.Number)

	// fetch context
	rootchainContext, err := rl.getRootChainContext()
	if err != nil {
		return
	}

	requiredConfirmations := rootchainContext.ChainmanagerParams.MainchainTxConfirmations
	headerNumber := newHeader.header.Number
	from := headerNumber

	// If incoming header is a `finalized` header, it can directly be considered as
	// the upper cap (i.e. the `to` value)
	//
	// If incoming header is a `latest` header, rely on `requiredConfirmations` to get
	// finalized block range.
	if !newHeader.isFinalized {
		// This check is only useful when the L1 blocks received are < requiredConfirmations
		// just for the below headerNumber -= requiredConfirmations math operation
		confirmationBlocks := big.NewInt(0).SetUint64(requiredConfirmations)
		if headerNumber.Cmp(confirmationBlocks) <= 0 {
			rl.Logger.Error("Block number less than Confirmations required", "blockNumber", headerNumber.Uint64, "confirmationsRequired", confirmationBlocks.Uint64)
			return
		}

		// subtract the `confirmationBlocks` to only consider blocks before that
		headerNumber = headerNumber.Sub(headerNumber, confirmationBlocks)

		// update the `from` value
		from = headerNumber
	}

	// get last block from storage
	hasLastBlock, _ := rl.storageClient.Has([]byte(lastRootBlockKey), nil)
	if hasLastBlock {
		lastBlockBytes, err := rl.storageClient.Get([]byte(lastRootBlockKey), nil)
		if err != nil {
			rl.Logger.Info("Error while fetching last block bytes from storage", "error", err)
			return
		}

		rl.Logger.Debug("Got rootchain last block from bridge storage", "lastBlock", string(lastBlockBytes))

		if result, err := strconv.ParseUint(string(lastBlockBytes), 10, 64); err == nil {
			if result >= headerNumber.Uint64() {
				return
			}

			from = big.NewInt(0).SetUint64(result + 1)
		}
	}

	// check env start height
	lockingStartHeightStr := os.Getenv("LOCKING_START_HEIGHT")
	if lockingStartHeightStr != "" && !hasLastBlock {
		lockingStartHeight, _ := strconv.ParseInt(lockingStartHeightStr, 10, 64)
		if lockingStartHeight > 0 {
			from = big.NewInt(lockingStartHeight)
		}
	}

	to := headerNumber

	// Prepare block range
	if to.Cmp(from) == -1 {
		from = to
	}

	// Handle events
	rl.queryAndBroadcastEvents(rootchainContext, from, to)

	// Set last block to storage
	if err = rl.storageClient.Put([]byte(lastRootBlockKey), []byte(to.String()), nil); err != nil {
		rl.Logger.Error("rl.storageClient.Put", "Error", err)
	}
}

// queryAndBroadcastEvents fetches supported events from the rootchain and handles all of them
func (rl *RootChainListener) queryAndBroadcastEvents(rootchainContext *RootChainListenerContext, fromBlock *big.Int, toBlock *big.Int) {
	// get chain params
	chainParams := rootchainContext.ChainmanagerParams.ChainParams

	rl.Logger.Info("Query rootchain event logs", "fromBlock", fromBlock, "toBlock", toBlock, "watchContract", chainParams.StakingInfoAddress.EthAddress())

	ctx, cancel := context.WithTimeout(context.Background(), rl.contractConnector.MainChainTimeout)
	defer cancel()

	// Fetch events from the rootchain
	logs, err := rl.contractConnector.MainChainClient.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: []ethCommon.Address{
			chainParams.StakingInfoAddress.EthAddress(),
		},
	})
	if err != nil {
		rl.Logger.Error("Error while filtering logs", "error", err)
		return
	} else if len(logs) > 0 {
		rl.Logger.Info("New logs found", "numberOfLogs", len(logs))
	}

	// Process filtered log
	for _, vLog := range logs {
		topic := vLog.Topics[0].Bytes()
		for _, abiObject := range rl.abis {
			selectedEvent := helper.EventByID(abiObject, topic)
			if selectedEvent == nil {
				continue
			}

			rl.handleLog(vLog, selectedEvent)
		}
	}
}

// getRootChainContext returns the root chain context
func (rl *RootChainListener) getRootChainContext() (*RootChainListenerContext, error) {
	chainmanagerParams, err := util.GetChainmanagerParams(rl.cliCtx)
	if err != nil {
		rl.Logger.Error("Error while fetching chain manager params", "error", err)
		return nil, err
	}

	return &RootChainListenerContext{
		ChainmanagerParams: chainmanagerParams,
	}, nil
}
