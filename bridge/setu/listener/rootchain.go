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

var maxBlockRange = big.NewInt(500) // max event polling range

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
}

func (rl *RootChainListener) StartPolling(ctx context.Context, pollInterval time.Duration, _ *big.Int) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	// fetch context
	rootchainCtx, err := rl.getRootChainContext()
	if err != nil {
		rl.Logger.Error("Failed to get root chain context", "error", err)
		return
	}
	confirmationCount := rootchainCtx.ChainmanagerParams.MainchainTxConfirmations
	for {
		select {
		case <-ctx.Done():
			rl.Logger.Info("RootChainListener polling stopped")
			return
		case <-ticker.C:
			currentHeight, err := rl.chainClient.BlockNumber(ctx)
			if err != nil {
				rl.Logger.Error("Failed to get current block height", "error", err)
				continue
			}

			maxProcessHeight := currentHeight - confirmationCount

			lastProcessedHeight, err := rl.getLastProcessedHeight()
			if err != nil {
				rl.Logger.Error("Failed to get last processed height", "error", err)
				continue
			}

			lockingStartHeight := rl.getLockingStartHeight()
			if lastProcessedHeight.Cmp(lockingStartHeight) < 0 {
				lastProcessedHeight = lockingStartHeight
			}

			if lastProcessedHeight.Cmp(big.NewInt(int64(maxProcessHeight))) >= 0 {
				continue
			}

			startHeight := new(big.Int).Add(lastProcessedHeight, big.NewInt(1))
			endHeight := new(big.Int).Add(startHeight, maxBlockRange)
			endHeight.Sub(endHeight, big.NewInt(1))

			if endHeight.Cmp(big.NewInt(int64(maxProcessHeight))) > 0 {
				endHeight = big.NewInt(int64(maxProcessHeight))
			}

			err = rl.queryAndBroadcastEvents(rootchainCtx, startHeight, endHeight)
			if err != nil {
				rl.Logger.Error("Failed to process blocks", err)
			} else {
				err = rl.storageClient.Put([]byte(lastRootBlockKey), []byte(endHeight.String()), nil)
				if err != nil {
					rl.Logger.Error("Error while saving lastRootBlockKey", "error", err)
				} else {
					rl.Logger.Info("Root chain listener processed successfully", "from", startHeight, "to", endHeight)
				}
			}
		}
	}
}

func (rl *RootChainListener) getLockingStartHeight() *big.Int {
	lockingStartHeightStr := os.Getenv("LOCKING_START_HEIGHT")
	if lockingStartHeightStr != "" {
		lockingStartHeight, err := strconv.ParseInt(lockingStartHeightStr, 10, 64)
		if err == nil && lockingStartHeight > 0 {
			return big.NewInt(lockingStartHeight)
		}
	}
	return big.NewInt(0)
}

func (rl *RootChainListener) getLastProcessedHeight() (*big.Int, error) {
	hasLastBlock, err := rl.storageClient.Has([]byte(lastRootBlockKey), nil)
	if err != nil {
		rl.Logger.Error("Error while checking existence of last block in storage", "error", err)
		return nil, err
	}

	if !hasLastBlock {
		rl.Logger.Debug("No last block found in storage")
		return big.NewInt(0), nil
	}

	lastBlockBytes, err := rl.storageClient.Get([]byte(lastRootBlockKey), nil)
	if err != nil {
		rl.Logger.Info("Error while fetching last block bytes from storage", "error", err)
		return nil, err
	}

	rl.Logger.Debug("Got rootchain last block from storage", "lastBlock", string(lastBlockBytes))

	lastBlockHeight, err := strconv.ParseUint(string(lastBlockBytes), 10, 64)
	if err != nil {
		rl.Logger.Info("Error while parsing last block height from storage", "error", err)
		return nil, err
	}

	return big.NewInt(0).SetUint64(lastBlockHeight), nil
}

// queryAndBroadcastEvents fetches supported events from the rootchain and handles all of them
func (rl *RootChainListener) queryAndBroadcastEvents(rootchainContext *RootChainListenerContext, fromBlock *big.Int, toBlock *big.Int) error {
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
		return err
	} else if len(logs) > 0 {
		rl.Logger.Info("New logs found", "numberOfLogs", len(logs))

		if len(rl.abis) == 0 {
			rl.Logger.Error("No ABI objects available")
		}
	}

	// Process filtered log
	for ix, vLog := range logs {
		rl.Logger.Info("Starting to handle log", "logIndex", ix, "logBlock", vLog.BlockNumber, "logTx", vLog.TxHash, "topic", vLog.Topics[0])
		topic := vLog.Topics[0].Bytes()
		for _, abiObject := range rl.abis {
			selectedEvent := helper.EventByID(abiObject, topic)
			if selectedEvent == nil {
				rl.Logger.Info("No matching event found", "topic", vLog.Topics[0])
				continue
			}

			err = rl.handleLog(vLog, selectedEvent)
			if err != nil {
				// UNCHECK: Only log error, not sure what the impact of repeated transmission events would be
				rl.Logger.Error("Error while handle log", "logBlock", vLog.BlockNumber, "logTx", vLog.TxHash, "error", err)

				// try once more
				err = rl.handleLog(vLog, selectedEvent)
				if err != nil {
					rl.Logger.Error("Error while handle log again", "logBlock", vLog.BlockNumber, "logTx", vLog.TxHash, "error", err)
				} else {
					rl.Logger.Info("Event saved while second handle log", "logBlock", vLog.BlockNumber, "logTx", vLog.TxHash, "eventName", selectedEvent.Name)
				}
			} else {
				rl.Logger.Info("Event saved while first handle log", "logBlock", vLog.BlockNumber, "logTx", vLog.TxHash, "eventName", selectedEvent.Name)
			}
		}
	}
	return nil
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
