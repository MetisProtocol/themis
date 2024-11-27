package helper

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	lru "github.com/hashicorp/golang-lru"

	"github.com/metis-seq/themis/contracts/erc20"

	"github.com/metis-seq/themis/contracts/sequencerset"
	"github.com/metis-seq/themis/contracts/stakemanager"
	"github.com/metis-seq/themis/contracts/stakinginfo"
	"github.com/metis-seq/themis/types"
)

// smart contracts' events names
const (
	newHeaderBlockEvent    = "NewHeaderBlock"
	stakedEvent            = "Locked"
	stakeUpdateEvent       = "LockUpdate"
	UnstakeInitEvent       = "UnlockInit"
	signerChangeEvent      = "SignerChange"
	BatchSubmitRewardEvent = "BatchSubmitReward"
)

// ContractsABIsMap is a cached map holding the ABIs of the contracts
var ContractsABIsMap = make(map[string]*abi.ABI)

// IContractCaller represents contract caller
type IContractCaller interface {
	GetValidatorInfo(valID types.ValidatorID, stakingManagerInstance *stakemanager.Stakemanager) (validator types.Validator, err error)
	GetBalance(address common.Address) (*big.Int, error)
	GetMainChainBlock(*big.Int) (*ethTypes.Header, error)
	GetMetisChainBlock(*big.Int) (*ethTypes.Header, error)
	IsTxConfirmed(common.Hash, uint64) bool
	GetConfirmedTxReceipt(common.Hash, uint64) (*ethTypes.Receipt, error)
	GetBlockNumberFromTxHash(common.Hash) (*big.Int, error)

	// decode validator events
	DecodeValidatorJoinEvent(common.Address, *ethTypes.Receipt, uint64) (*stakinginfo.StakinginfoLocked, error)
	DecodeValidatorStakeUpdateEvent(common.Address, *ethTypes.Receipt, uint64) (*stakinginfo.StakinginfoLockUpdate, error)
	DecodeValidatorExitEvent(common.Address, *ethTypes.Receipt, uint64) (*stakinginfo.StakinginfoUnlockInit, error)
	DecodeSignerUpdateEvent(common.Address, *ethTypes.Receipt, uint64) (*stakinginfo.StakinginfoSignerChange, error)
	DecodeBatchSubmitRewardEvent(common.Address, *ethTypes.Receipt, uint64) (*stakinginfo.StakinginfoBatchSubmitReward, error)

	GetMainTxReceipt(common.Hash) (*ethTypes.Receipt, error)
	GetMetisTxReceipt(common.Hash) (*ethTypes.Receipt, error)
	ApproveTokens(*big.Int, common.Address, common.Address, *erc20.Erc20) error
	StakeFor(common.Address, *big.Int, *big.Int, bool, common.Address, *stakemanager.Stakemanager) error

	// metis related contracts
	CurrentSpanNumber(validatorSet *sequencerset.Sequencerset) (Number *big.Int)
	// GetSpanDetails(id *big.Int, validatorSet *sequencerset.Sequencerset) (*big.Int, *big.Int, *big.Int, error)
	CheckIfBlocksExist(end uint64) bool

	GetStakingInfoInstance(stakingInfoAddress common.Address) (*stakinginfo.Stakinginfo, error)
	GetSequencerSetInstance(sequencerSetAddress common.Address) (*sequencerset.Sequencerset, error)
	GetStakeManagerInstance(stakingManagerAddress common.Address) (*stakemanager.Stakemanager, error)
	GetMetisTokenInstance(metisTokenAddress common.Address) (*erc20.Erc20, error)
}

// ContractCaller contract caller
type ContractCaller struct {
	MainChainClient   *ethclient.Client
	MainChainRPC      *rpc.Client
	MainChainTimeout  time.Duration
	MetisChainClient  *ethclient.Client
	MetisChainRPC     *rpc.Client
	MetisChainTimeout time.Duration

	StakingInfoABI  abi.ABI
	ValidatorSetABI abi.ABI
	StakeManagerABI abi.ABI
	MetisTokenABI   abi.ABI

	ReceiptCache *lru.Cache

	ContractInstanceCache map[common.Address]interface{}
}

// use global variable due to the ContractCaller is not a point
var contractInstanceMutx sync.Mutex

type txExtraInfo struct {
	BlockNumber *string         `json:"blockNumber,omitempty"`
	BlockHash   *common.Hash    `json:"blockHash,omitempty"`
	From        *common.Address `json:"from,omitempty"`
}

type rpcTransaction struct {
	txExtraInfo
}

// NewContractCaller contract caller
func NewContractCaller() (contractCallerObj ContractCaller, err error) {
	config := GetConfig()
	contractCallerObj.MainChainClient = GetMainClient()
	contractCallerObj.MainChainTimeout = config.EthRPCTimeout
	contractCallerObj.MetisChainClient = GetMetisClient()
	contractCallerObj.MetisChainTimeout = config.MetisRPCTimeout
	contractCallerObj.MainChainRPC = GetMainChainRPCClient()
	contractCallerObj.MetisChainRPC = GetMetisRPCClient()
	contractCallerObj.ReceiptCache, err = lru.New(1000)

	if err != nil {
		return contractCallerObj, err
	}

	// listeners and processors instance cache (address->ABI)
	contractCallerObj.ContractInstanceCache = make(map[common.Address]interface{})
	// package global cache (string->ABI)
	if err = populateABIs(&contractCallerObj); err != nil {
		return contractCallerObj, err
	}

	return contractCallerObj, nil
}

// GetStakingInfoInstance returns stakingInfo contract instance for selected base chain
func (c *ContractCaller) GetStakingInfoInstance(stakingInfoAddress common.Address) (*stakinginfo.Stakinginfo, error) {
	contractInstance, ok := c.ContractInstanceCache[stakingInfoAddress]
	if !ok {
		contractInstanceMutx.Lock()
		defer contractInstanceMutx.Unlock()

		ci, err := stakinginfo.NewStakinginfo(stakingInfoAddress, c.MainChainClient)
		if err != nil {
			return nil, err
		}

		c.ContractInstanceCache[stakingInfoAddress] = ci
		return ci, nil
	}

	return contractInstance.(*stakinginfo.Stakinginfo), nil
}

// GetSequencerSetInstance returns sequencer set contract instance for metis chain
func (c *ContractCaller) GetSequencerSetInstance(sequencerSetAddress common.Address) (*sequencerset.Sequencerset, error) {
	contractInstance, ok := c.ContractInstanceCache[sequencerSetAddress]
	if !ok {
		contractInstanceMutx.Lock()
		defer contractInstanceMutx.Unlock()

		ci, err := sequencerset.NewSequencerset(sequencerSetAddress, c.MetisChainClient)
		if err != nil {
			return nil, err
		}
		c.ContractInstanceCache[sequencerSetAddress] = ci

		return ci, nil
	}

	return contractInstance.(*sequencerset.Sequencerset), nil
}

// GetStakeManagerInstance returns stakingInfo contract instance for selected base chain
func (c *ContractCaller) GetStakeManagerInstance(stakingManagerAddress common.Address) (*stakemanager.Stakemanager, error) {
	contractInstance, ok := c.ContractInstanceCache[stakingManagerAddress]
	if !ok {
		contractInstanceMutx.Lock()
		defer contractInstanceMutx.Unlock()

		ci, err := stakemanager.NewStakemanager(stakingManagerAddress, c.MainChainClient)
		if err != nil {
			return nil, err
		}
		c.ContractInstanceCache[stakingManagerAddress] = ci
		return ci, nil
	}

	return contractInstance.(*stakemanager.Stakemanager), nil
}

// GetMetisTokenInstance returns stakingInfo contract instance for selected base chain
func (c *ContractCaller) GetMetisTokenInstance(metisTokenAddress common.Address) (*erc20.Erc20, error) {
	contractInstance, ok := c.ContractInstanceCache[metisTokenAddress]
	if !ok {
		contractInstanceMutx.Lock()
		defer contractInstanceMutx.Unlock()

		ci, err := erc20.NewErc20(metisTokenAddress, c.MainChainClient)
		if err != nil {
			return nil, err
		}

		c.ContractInstanceCache[metisTokenAddress] = ci

		return ci, nil
	}

	return contractInstance.(*erc20.Erc20), nil
}

// GetBalance get balance of account (returns big.Int balance wont fit in uint64)
func (c *ContractCaller) GetBalance(address common.Address) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MainChainTimeout)
	defer cancel()

	balance, err := c.MainChainClient.BalanceAt(ctx, address, nil)
	if err != nil {
		Logger.Error("Unable to fetch balance of account from root chain", "Address", address.String(), "error", err)
		return big.NewInt(0), err
	}

	return balance, nil
}

// GetValidatorInfo get validator info
func (c *ContractCaller) GetValidatorInfo(valID types.ValidatorID, stakingManagerInstance *stakemanager.Stakemanager) (validator types.Validator, err error) {
	stakerDetails, err := stakingManagerInstance.Sequencers(nil, big.NewInt(int64(valID)))
	if err != nil {
		Logger.Error("Error fetching validator information from stake manager", "validatorId", valID, "status", stakerDetails.Status, "error", err)
		return
	}

	newAmount, err := GetPowerFromAmount(stakerDetails.Amount)
	if err != nil {
		return
	}

	// newAmount
	validator = types.Validator{
		ID:          valID,
		VotingPower: newAmount.Int64(),
		StartBatch:  stakerDetails.ActivationBatch.Uint64(),
		EndBatch:    stakerDetails.DeactivationBatch.Uint64(),
		Signer:      types.BytesToThemisAddress(stakerDetails.Signer.Bytes()),
	}

	return validator, nil
}

// GetMainChainBlock returns main chain block header
func (c *ContractCaller) GetMainChainBlock(blockNum *big.Int) (header *ethTypes.Header, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MainChainTimeout)
	defer cancel()

	latestBlock, err := c.MainChainClient.HeaderByNumber(ctx, blockNum)
	if err != nil {
		Logger.Error("Unable to connect to main chain", "error", err)
		return
	}

	return latestBlock, nil
}

// GetMainChainFinalizedBlock returns finalized main chain block header (post-merge)
func (c *ContractCaller) GetMainChainFinalizedBlock() (header *ethTypes.Header, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MainChainTimeout)
	defer cancel()

	latestFinalizedBlock, err := c.MainChainClient.HeaderByNumber(ctx, big.NewInt(int64(rpc.FinalizedBlockNumber)))
	if err != nil {
		Logger.Error("Unable to connect to main chain", "error", err)
		return
	}

	return latestFinalizedBlock, nil
}

// GetMainChainBlockTime returns main chain block time
func (c *ContractCaller) GetMainChainBlockTime(ctx context.Context, blockNum uint64) (time.Time, error) {
	ctx, cancel := context.WithTimeout(ctx, c.MainChainTimeout)
	defer cancel()

	latestBlock, err := c.MainChainClient.BlockByNumber(ctx, big.NewInt(0).SetUint64(blockNum))
	if err != nil {
		Logger.Error("Unable to connect to main chain", "error", err)
		return time.Time{}, err
	}

	return time.Unix(int64(latestBlock.Time()), 0), nil
}

// GetMainChainID returns main chain block time
func (c *ContractCaller) GetMainChainID() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MainChainTimeout)
	defer cancel()

	chainID, err := c.MainChainClient.ChainID(ctx)
	if err != nil {
		Logger.Error("Unable to connect to main chain", "error", err)
		return 0, err
	}

	return chainID.Uint64(), nil
}

// GetMetisChainBlock returns child chain block header
func (c *ContractCaller) GetMetisChainBlock(blockNum *big.Int) (header *ethTypes.Header, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MetisChainTimeout)
	defer cancel()

	latestBlock, err := c.MetisChainClient.HeaderByNumber(ctx, blockNum)
	if err != nil {
		Logger.Error("Unable to connect to metis chain", "error", err)
		return
	}

	return latestBlock, nil
}

// GetMetisChainID returns child chain block header
func (c *ContractCaller) GetMetisChainID() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MetisChainTimeout)
	defer cancel()

	chainID, err := c.MetisChainClient.ChainID(ctx)
	if err != nil {
		Logger.Error("Unable to connect to main chain", "error", err)
		return 0, err
	}

	return chainID.Uint64(), nil
}

// GetMetisNonce returns child chain block header
func (c *ContractCaller) GetMetisNonce(account common.Address) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MetisChainTimeout)
	defer cancel()

	nonce, err := c.MetisChainClient.NonceAt(ctx, account, nil)
	if err != nil {
		Logger.Error("Unable to connect to main chain", "error", err)
		return 0, err
	}

	return nonce, nil
}

func (c *ContractCaller) GetMetisGasprice() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MetisChainTimeout)
	defer cancel()

	gasPrice, err := c.MetisChainClient.SuggestGasPrice(ctx)
	if err != nil {
		Logger.Error("Unable to connect to main chain", "error", err)
		return nil, err
	}

	return gasPrice, nil
}

// GetBlockNumberFromTxHash gets block number of transaction
func (c *ContractCaller) GetBlockNumberFromTxHash(tx common.Hash) (*big.Int, error) {
	var rpcTx rpcTransaction
	if err := c.MainChainRPC.CallContext(context.Background(), &rpcTx, "eth_getTransactionByHash", tx); err != nil {
		return nil, err
	}

	if rpcTx.BlockNumber == nil {
		return nil, errors.New("no tx found")
	}

	blkNum := big.NewInt(0)

	blkNum, ok := blkNum.SetString(*rpcTx.BlockNumber, 0)
	if !ok {
		return nil, errors.New("unable to set string")
	}

	return blkNum, nil
}

// IsTxConfirmed is tx confirmed
func (c *ContractCaller) IsTxConfirmed(tx common.Hash, requiredConfirmations uint64) bool {
	// get main tx receipt
	receipt, err := c.GetConfirmedTxReceipt(tx, requiredConfirmations)
	if receipt == nil || err != nil {
		return false
	}

	return true
}

// GetConfirmedTxReceipt returns confirmed tx receipt
func (c *ContractCaller) GetConfirmedTxReceipt(tx common.Hash, requiredConfirmations uint64) (*ethTypes.Receipt, error) {
	var receipt *ethTypes.Receipt

	receiptCache, ok := c.ReceiptCache.Get(tx.String())
	if !ok {
		var err error

		// get main tx receipt
		receipt, err = c.GetMainTxReceipt(tx)
		if err != nil {
			Logger.Error("Error while fetching mainChain receipt", "txHash", tx.Hex(), "error", err)
			return nil, err
		}

		c.ReceiptCache.Add(tx.String(), receipt)
	} else {
		receipt, _ = receiptCache.(*ethTypes.Receipt)
	}

	receiptBlockNumber := receipt.BlockNumber.Uint64()

	Logger.Debug("Tx included in block", "block", receiptBlockNumber, "tx", tx)

	latestBlk, err := c.GetMainChainBlock(nil)
	if err != nil {
		Logger.Error("error getting latest block from main chain", "error", err)
		return nil, err
	}

	diff := latestBlk.Number.Uint64() - receiptBlockNumber
	Logger.Debug("Latest block on main chain obtained", "Block", latestBlk.Number.Uint64(), "receipt block", receiptBlockNumber, "diff", diff, "requiredConfirmations", requiredConfirmations)

	if diff < requiredConfirmations {
		Logger.Info("not enough confirmations")
		return nil, errors.New("not enough confirmations")
	}

	return receipt, nil
}

//
// Validator decode events
//

// DecodeValidatorJoinEvent represents validator staked event
func (c *ContractCaller) DecodeValidatorJoinEvent(contractAddress common.Address, receipt *ethTypes.Receipt, logIndex uint64) (*stakinginfo.StakinginfoLocked, error) {
	event := new(stakinginfo.StakinginfoLocked)

	found := false

	for _, vLog := range receipt.Logs {
		if uint64(vLog.Index) == logIndex && bytes.Equal(vLog.Address.Bytes(), contractAddress.Bytes()) {
			found = true

			if err := UnpackLog(&c.StakingInfoABI, event, stakedEvent, vLog); err != nil {
				return nil, err
			}

			break
		}
	}

	if !found {
		return nil, errors.New("event not found")
	}

	return event, nil
}

// DecodeValidatorStakeUpdateEvent represents validator stake update event
func (c *ContractCaller) DecodeValidatorStakeUpdateEvent(contractAddress common.Address, receipt *ethTypes.Receipt, logIndex uint64) (*stakinginfo.StakinginfoLockUpdate, error) {
	var (
		event = new(stakinginfo.StakinginfoLockUpdate)
		found = false
	)

	for _, vLog := range receipt.Logs {
		if uint64(vLog.Index) == logIndex && bytes.Equal(vLog.Address.Bytes(), contractAddress.Bytes()) {
			found = true

			if err := UnpackLog(&c.StakingInfoABI, event, stakeUpdateEvent, vLog); err != nil {
				return nil, err
			}

			break
		}
	}

	if !found {
		return nil, errors.New("event not found")
	}

	return event, nil
}

// DecodeValidatorExitEvent represents validator stake unStake event
func (c *ContractCaller) DecodeValidatorExitEvent(contractAddress common.Address, receipt *ethTypes.Receipt, logIndex uint64) (*stakinginfo.StakinginfoUnlockInit, error) {
	var (
		event = new(stakinginfo.StakinginfoUnlockInit)
		found = false
	)

	for _, vLog := range receipt.Logs {
		if uint64(vLog.Index) == logIndex && bytes.Equal(vLog.Address.Bytes(), contractAddress.Bytes()) {
			found = true

			if err := UnpackLog(&c.StakingInfoABI, event, UnstakeInitEvent, vLog); err != nil {
				return nil, err
			}

			break
		}
	}

	if !found {
		return nil, errors.New("event not found")
	}

	return event, nil
}

// DecodeSignerUpdateEvent represents sig update event
func (c *ContractCaller) DecodeSignerUpdateEvent(contractAddress common.Address, receipt *ethTypes.Receipt, logIndex uint64) (*stakinginfo.StakinginfoSignerChange, error) {
	var (
		event = new(stakinginfo.StakinginfoSignerChange)
		found = false
	)

	for _, vLog := range receipt.Logs {
		if uint64(vLog.Index) == logIndex && bytes.Equal(vLog.Address.Bytes(), contractAddress.Bytes()) {
			found = true

			if err := UnpackLog(&c.StakingInfoABI, event, signerChangeEvent, vLog); err != nil {
				return nil, err
			}

			break
		}
	}

	if !found {
		return nil, errors.New("event not found")
	}

	return event, nil
}

// DecodeBatchSubmitRewardEvent represents sig update event
func (c *ContractCaller) DecodeBatchSubmitRewardEvent(contractAddress common.Address, receipt *ethTypes.Receipt, logIndex uint64) (*stakinginfo.StakinginfoBatchSubmitReward, error) {
	var (
		event = new(stakinginfo.StakinginfoBatchSubmitReward)
		found = false
	)

	for _, vLog := range receipt.Logs {
		if uint64(vLog.Index) == logIndex && bytes.Equal(vLog.Address.Bytes(), contractAddress.Bytes()) {
			found = true

			if err := UnpackLog(&c.StakingInfoABI, event, BatchSubmitRewardEvent, vLog); err != nil {
				return nil, err
			}

			break
		}
	}

	if !found {
		return nil, errors.New("event not found")
	}

	return event, nil
}

//
// Span related functions
//

// CurrentSpanNumber get current span
func (c *ContractCaller) CurrentSpanNumber(validatorSetInstance *sequencerset.Sequencerset) (Number *big.Int) {
	result, err := validatorSetInstance.CurrentEpochNumber(nil)
	if err != nil {
		Logger.Error("Unable to get current span number", "error", err)
		return nil
	}

	return result
}

// GetSpanDetails get span details
func (c *ContractCaller) GetSpanDetails(id *big.Int, validatorSetInstance *sequencerset.Sequencerset) (
	*big.Int,
	*big.Int,
	*big.Int,
	error,
) {
	// d, err := validatorSetInstance.GetEpochByBlock(nil, id)
	// return d.Number, d.StartBlock, d.EndBlock, err
	return new(big.Int), new(big.Int), new(big.Int), nil
}

// CheckIfBlocksExist - check if the given block exists on local chain
func (c *ContractCaller) CheckIfBlocksExist(end uint64) bool {
	// Get block by number.
	var block *ethTypes.Header

	err := c.MetisChainRPC.Call(&block, "eth_getBlockByNumber", fmt.Sprintf("0x%x", end), false)
	if err != nil {
		return false
	}

	return end == block.Number.Uint64()
}

//
// Receipt functions
//

// GetMainTxReceipt returns main tx receipt
func (c *ContractCaller) GetMainTxReceipt(txHash common.Hash) (*ethTypes.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MainChainTimeout)
	defer cancel()

	return c.getTxReceipt(ctx, c.MainChainClient, txHash)
}

// GetMetisTxReceipt returns metis tx receipt
func (c *ContractCaller) GetMetisTxReceipt(txHash common.Hash) (*ethTypes.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MetisChainTimeout)
	defer cancel()

	return c.getTxReceipt(ctx, c.MetisChainClient, txHash)
}

func (c *ContractCaller) getTxReceipt(ctx context.Context, client *ethclient.Client, txHash common.Hash) (*ethTypes.Receipt, error) {
	return client.TransactionReceipt(ctx, txHash)
}

// utility and helper methods

// populateABIs fills the package level cache for contracts' ABIs
// When called the first time, ContractsABIsMap will be filled and getABI method won't be invoked the next times
// This reduces the number of calls to json decode methods made by the contract caller
// It uses ABIs' definitions instead of contracts addresses, as the latter might not be available at init time
func populateABIs(contractCallerObj *ContractCaller) error {
	var ccAbi *abi.ABI

	var err error

	contractsABIs := [4]string{
		stakinginfo.StakinginfoABI,
		sequencerset.SequencersetABI,
		stakemanager.StakemanagerABI,
		erc20.Erc20ABI}

	// iterate over supported ABIs
	for _, contractABI := range contractsABIs {
		ccAbi, err = chooseContractCallerABI(contractCallerObj, contractABI)
		if err != nil {
			Logger.Error("Error while fetching contract caller ABI", "error", err)
			return err
		}

		if ContractsABIsMap[contractABI] == nil {
			// fills cached abi map
			if *ccAbi, err = getABI(contractABI); err != nil {
				Logger.Error("Error while getting ABI for contract caller", "name", contractABI, "error", err)
				return err
			} else {
				// init ABI
				ContractsABIsMap[contractABI] = ccAbi
			}
		} else {
			// use cached abi
			*ccAbi = *ContractsABIsMap[contractABI]
		}
	}

	return nil
}

// chooseContractCallerABI extracts and returns the abo.ABI object from the contractCallerObj based on its abi string
func chooseContractCallerABI(contractCallerObj *ContractCaller, abi string) (*abi.ABI, error) {
	switch abi {
	case stakinginfo.StakinginfoABI:
		return &contractCallerObj.StakingInfoABI, nil
	case sequencerset.SequencersetABI:
		return &contractCallerObj.ValidatorSetABI, nil
	case stakemanager.StakemanagerABI:
		return &contractCallerObj.StakeManagerABI, nil
	case erc20.Erc20ABI:
		return &contractCallerObj.MetisTokenABI, nil
	}

	return nil, errors.New("no ABI associated with such data")
}

// getABI returns the contract's ABI struct from on its JSON representation
func getABI(data string) (abi.ABI, error) {
	return abi.JSON(strings.NewReader(data))
}
