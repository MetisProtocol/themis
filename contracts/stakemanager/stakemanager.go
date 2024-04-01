// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package stakemanager

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// StakemanagerMetaData contains all meta data concerning the Stakemanager contract.
var StakemanagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"NoRewardRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoSuchSeq\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotSeqOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotSeqSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotWhitelisted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NullAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnedSequencer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SeqNotActive\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignerExisted\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newReward\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldReward\",\"type\":\"uint256\"}],\"name\":\"RewardUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_seqId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"SequencerOwnerChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_seqId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"}],\"name\":\"SequencerRewardRecipientChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_n\",\"type\":\"uint256\"}],\"name\":\"SetSignerUpdateThrottle\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"_yes\",\"type\":\"bool\"}],\"name\":\"SetWhitelist\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_newMpc\",\"type\":\"address\"}],\"name\":\"UpdateMpc\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_cur\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_prev\",\"type\":\"uint256\"}],\"name\":\"WithrawDelayTimeChange\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BLOCK_REWARD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"WITHDRAWAL_DELAY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_batchId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_startEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_endEpoch\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"_seqs\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_blocks\",\"type\":\"uint256[]\"}],\"name\":\"batchSubmitRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"totalReward\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"curBatchState\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"number\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endEpoch\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentBatch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"escrow\",\"outputs\":[{\"internalType\":\"contractILockingInfo\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_seqId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"_l2Gas\",\"type\":\"uint32\"}],\"name\":\"forceUnlock\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_escrow\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_signer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_signerPubkey\",\"type\":\"bytes\"}],\"name\":\"lockFor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_signer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_rewardRecipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_signerPubkey\",\"type\":\"bytes\"}],\"name\":\"lockWithRewardRecipient\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mpcAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_seqId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_lockReward\",\"type\":\"bool\"}],\"name\":\"relock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"seqOwners\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"seqId\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"seqSigners\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"seqId\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumISequencerInfo.Status\",\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"seqStatuses\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"count\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"seqId\",\"type\":\"uint256\"}],\"name\":\"sequencers\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"activationBatch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedBatch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deactivationBatch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deactivationTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unlockClaimTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"rewardRecipient\",\"type\":\"address\"},{\"internalType\":\"enumISequencerInfo.Status\",\"name\":\"status\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_yes\",\"type\":\"bool\"}],\"name\":\"setPause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_seqId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"setSequencerOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_seqId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"}],\"name\":\"setSequencerRewardRecipient\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_n\",\"type\":\"uint256\"}],\"name\":\"setSignerUpdateThrottle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_yes\",\"type\":\"bool\"}],\"name\":\"setWhitelist\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"signerUpdateThrottle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSequencers\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_seqId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"_l2Gas\",\"type\":\"uint32\"}],\"name\":\"unlock\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_seqId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"_l2Gas\",\"type\":\"uint32\"}],\"name\":\"unlockClaim\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newReward\",\"type\":\"uint256\"}],\"name\":\"updateBlockReward\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newMpc\",\"type\":\"address\"}],\"name\":\"updateMpc\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_seqId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_signerPubkey\",\"type\":\"bytes\"}],\"name\":\"updateSigner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_time\",\"type\":\"uint256\"}],\"name\":\"updateWithdrawDelayTimeValue\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"whitelist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_seqId\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"_l2Gas\",\"type\":\"uint32\"}],\"name\":\"withdrawRewards\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

// StakemanagerABI is the input ABI used to generate the binding from.
// Deprecated: Use StakemanagerMetaData.ABI instead.
var StakemanagerABI = StakemanagerMetaData.ABI

// Stakemanager is an auto generated Go binding around an Ethereum contract.
type Stakemanager struct {
	StakemanagerCaller     // Read-only binding to the contract
	StakemanagerTransactor // Write-only binding to the contract
	StakemanagerFilterer   // Log filterer for contract events
}

// StakemanagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakemanagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakemanagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakemanagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakemanagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakemanagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakemanagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakemanagerSession struct {
	Contract     *Stakemanager     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakemanagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakemanagerCallerSession struct {
	Contract *StakemanagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// StakemanagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakemanagerTransactorSession struct {
	Contract     *StakemanagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// StakemanagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakemanagerRaw struct {
	Contract *Stakemanager // Generic contract binding to access the raw methods on
}

// StakemanagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakemanagerCallerRaw struct {
	Contract *StakemanagerCaller // Generic read-only contract binding to access the raw methods on
}

// StakemanagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakemanagerTransactorRaw struct {
	Contract *StakemanagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStakemanager creates a new instance of Stakemanager, bound to a specific deployed contract.
func NewStakemanager(address common.Address, backend bind.ContractBackend) (*Stakemanager, error) {
	contract, err := bindStakemanager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Stakemanager{StakemanagerCaller: StakemanagerCaller{contract: contract}, StakemanagerTransactor: StakemanagerTransactor{contract: contract}, StakemanagerFilterer: StakemanagerFilterer{contract: contract}}, nil
}

// NewStakemanagerCaller creates a new read-only instance of Stakemanager, bound to a specific deployed contract.
func NewStakemanagerCaller(address common.Address, caller bind.ContractCaller) (*StakemanagerCaller, error) {
	contract, err := bindStakemanager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakemanagerCaller{contract: contract}, nil
}

// NewStakemanagerTransactor creates a new write-only instance of Stakemanager, bound to a specific deployed contract.
func NewStakemanagerTransactor(address common.Address, transactor bind.ContractTransactor) (*StakemanagerTransactor, error) {
	contract, err := bindStakemanager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakemanagerTransactor{contract: contract}, nil
}

// NewStakemanagerFilterer creates a new log filterer instance of Stakemanager, bound to a specific deployed contract.
func NewStakemanagerFilterer(address common.Address, filterer bind.ContractFilterer) (*StakemanagerFilterer, error) {
	contract, err := bindStakemanager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakemanagerFilterer{contract: contract}, nil
}

// bindStakemanager binds a generic wrapper to an already deployed contract.
func bindStakemanager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StakemanagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Stakemanager *StakemanagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Stakemanager.Contract.StakemanagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Stakemanager *StakemanagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stakemanager.Contract.StakemanagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Stakemanager *StakemanagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Stakemanager.Contract.StakemanagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Stakemanager *StakemanagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Stakemanager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Stakemanager *StakemanagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stakemanager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Stakemanager *StakemanagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Stakemanager.Contract.contract.Transact(opts, method, params...)
}

// BLOCKREWARD is a free data retrieval call binding the contract method 0x7f05b9ef.
//
// Solidity: function BLOCK_REWARD() view returns(uint256)
func (_Stakemanager *StakemanagerCaller) BLOCKREWARD(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Stakemanager.contract.Call(opts, &out, "BLOCK_REWARD")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BLOCKREWARD is a free data retrieval call binding the contract method 0x7f05b9ef.
//
// Solidity: function BLOCK_REWARD() view returns(uint256)
func (_Stakemanager *StakemanagerSession) BLOCKREWARD() (*big.Int, error) {
	return _Stakemanager.Contract.BLOCKREWARD(&_Stakemanager.CallOpts)
}

// BLOCKREWARD is a free data retrieval call binding the contract method 0x7f05b9ef.
//
// Solidity: function BLOCK_REWARD() view returns(uint256)
func (_Stakemanager *StakemanagerCallerSession) BLOCKREWARD() (*big.Int, error) {
	return _Stakemanager.Contract.BLOCKREWARD(&_Stakemanager.CallOpts)
}

// WITHDRAWALDELAY is a free data retrieval call binding the contract method 0x0ebb172a.
//
// Solidity: function WITHDRAWAL_DELAY() view returns(uint256)
func (_Stakemanager *StakemanagerCaller) WITHDRAWALDELAY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Stakemanager.contract.Call(opts, &out, "WITHDRAWAL_DELAY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WITHDRAWALDELAY is a free data retrieval call binding the contract method 0x0ebb172a.
//
// Solidity: function WITHDRAWAL_DELAY() view returns(uint256)
func (_Stakemanager *StakemanagerSession) WITHDRAWALDELAY() (*big.Int, error) {
	return _Stakemanager.Contract.WITHDRAWALDELAY(&_Stakemanager.CallOpts)
}

// WITHDRAWALDELAY is a free data retrieval call binding the contract method 0x0ebb172a.
//
// Solidity: function WITHDRAWAL_DELAY() view returns(uint256)
func (_Stakemanager *StakemanagerCallerSession) WITHDRAWALDELAY() (*big.Int, error) {
	return _Stakemanager.Contract.WITHDRAWALDELAY(&_Stakemanager.CallOpts)
}

// CurBatchState is a free data retrieval call binding the contract method 0xb4472970.
//
// Solidity: function curBatchState() view returns(uint256 id, uint256 number, uint256 startEpoch, uint256 endEpoch)
func (_Stakemanager *StakemanagerCaller) CurBatchState(opts *bind.CallOpts) (struct {
	Id         *big.Int
	Number     *big.Int
	StartEpoch *big.Int
	EndEpoch   *big.Int
}, error) {
	var out []interface{}
	err := _Stakemanager.contract.Call(opts, &out, "curBatchState")

	outstruct := new(struct {
		Id         *big.Int
		Number     *big.Int
		StartEpoch *big.Int
		EndEpoch   *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Id = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Number = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.StartEpoch = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.EndEpoch = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// CurBatchState is a free data retrieval call binding the contract method 0xb4472970.
//
// Solidity: function curBatchState() view returns(uint256 id, uint256 number, uint256 startEpoch, uint256 endEpoch)
func (_Stakemanager *StakemanagerSession) CurBatchState() (struct {
	Id         *big.Int
	Number     *big.Int
	StartEpoch *big.Int
	EndEpoch   *big.Int
}, error) {
	return _Stakemanager.Contract.CurBatchState(&_Stakemanager.CallOpts)
}

// CurBatchState is a free data retrieval call binding the contract method 0xb4472970.
//
// Solidity: function curBatchState() view returns(uint256 id, uint256 number, uint256 startEpoch, uint256 endEpoch)
func (_Stakemanager *StakemanagerCallerSession) CurBatchState() (struct {
	Id         *big.Int
	Number     *big.Int
	StartEpoch *big.Int
	EndEpoch   *big.Int
}, error) {
	return _Stakemanager.Contract.CurBatchState(&_Stakemanager.CallOpts)
}

// CurrentBatch is a free data retrieval call binding the contract method 0x76cd940e.
//
// Solidity: function currentBatch() view returns(uint256)
func (_Stakemanager *StakemanagerCaller) CurrentBatch(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Stakemanager.contract.Call(opts, &out, "currentBatch")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentBatch is a free data retrieval call binding the contract method 0x76cd940e.
//
// Solidity: function currentBatch() view returns(uint256)
func (_Stakemanager *StakemanagerSession) CurrentBatch() (*big.Int, error) {
	return _Stakemanager.Contract.CurrentBatch(&_Stakemanager.CallOpts)
}

// CurrentBatch is a free data retrieval call binding the contract method 0x76cd940e.
//
// Solidity: function currentBatch() view returns(uint256)
func (_Stakemanager *StakemanagerCallerSession) CurrentBatch() (*big.Int, error) {
	return _Stakemanager.Contract.CurrentBatch(&_Stakemanager.CallOpts)
}

// Escrow is a free data retrieval call binding the contract method 0xe2fdcc17.
//
// Solidity: function escrow() view returns(address)
func (_Stakemanager *StakemanagerCaller) Escrow(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Stakemanager.contract.Call(opts, &out, "escrow")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Escrow is a free data retrieval call binding the contract method 0xe2fdcc17.
//
// Solidity: function escrow() view returns(address)
func (_Stakemanager *StakemanagerSession) Escrow() (common.Address, error) {
	return _Stakemanager.Contract.Escrow(&_Stakemanager.CallOpts)
}

// Escrow is a free data retrieval call binding the contract method 0xe2fdcc17.
//
// Solidity: function escrow() view returns(address)
func (_Stakemanager *StakemanagerCallerSession) Escrow() (common.Address, error) {
	return _Stakemanager.Contract.Escrow(&_Stakemanager.CallOpts)
}

// MpcAddress is a free data retrieval call binding the contract method 0x111f4630.
//
// Solidity: function mpcAddress() view returns(address)
func (_Stakemanager *StakemanagerCaller) MpcAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Stakemanager.contract.Call(opts, &out, "mpcAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MpcAddress is a free data retrieval call binding the contract method 0x111f4630.
//
// Solidity: function mpcAddress() view returns(address)
func (_Stakemanager *StakemanagerSession) MpcAddress() (common.Address, error) {
	return _Stakemanager.Contract.MpcAddress(&_Stakemanager.CallOpts)
}

// MpcAddress is a free data retrieval call binding the contract method 0x111f4630.
//
// Solidity: function mpcAddress() view returns(address)
func (_Stakemanager *StakemanagerCallerSession) MpcAddress() (common.Address, error) {
	return _Stakemanager.Contract.MpcAddress(&_Stakemanager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Stakemanager *StakemanagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Stakemanager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Stakemanager *StakemanagerSession) Owner() (common.Address, error) {
	return _Stakemanager.Contract.Owner(&_Stakemanager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Stakemanager *StakemanagerCallerSession) Owner() (common.Address, error) {
	return _Stakemanager.Contract.Owner(&_Stakemanager.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Stakemanager *StakemanagerCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Stakemanager.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Stakemanager *StakemanagerSession) Paused() (bool, error) {
	return _Stakemanager.Contract.Paused(&_Stakemanager.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Stakemanager *StakemanagerCallerSession) Paused() (bool, error) {
	return _Stakemanager.Contract.Paused(&_Stakemanager.CallOpts)
}

// SeqOwners is a free data retrieval call binding the contract method 0x169abefc.
//
// Solidity: function seqOwners(address owner) view returns(uint256 seqId)
func (_Stakemanager *StakemanagerCaller) SeqOwners(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Stakemanager.contract.Call(opts, &out, "seqOwners", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SeqOwners is a free data retrieval call binding the contract method 0x169abefc.
//
// Solidity: function seqOwners(address owner) view returns(uint256 seqId)
func (_Stakemanager *StakemanagerSession) SeqOwners(owner common.Address) (*big.Int, error) {
	return _Stakemanager.Contract.SeqOwners(&_Stakemanager.CallOpts, owner)
}

// SeqOwners is a free data retrieval call binding the contract method 0x169abefc.
//
// Solidity: function seqOwners(address owner) view returns(uint256 seqId)
func (_Stakemanager *StakemanagerCallerSession) SeqOwners(owner common.Address) (*big.Int, error) {
	return _Stakemanager.Contract.SeqOwners(&_Stakemanager.CallOpts, owner)
}

// SeqSigners is a free data retrieval call binding the contract method 0xbeb26755.
//
// Solidity: function seqSigners(address signer) view returns(uint256 seqId)
func (_Stakemanager *StakemanagerCaller) SeqSigners(opts *bind.CallOpts, signer common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Stakemanager.contract.Call(opts, &out, "seqSigners", signer)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SeqSigners is a free data retrieval call binding the contract method 0xbeb26755.
//
// Solidity: function seqSigners(address signer) view returns(uint256 seqId)
func (_Stakemanager *StakemanagerSession) SeqSigners(signer common.Address) (*big.Int, error) {
	return _Stakemanager.Contract.SeqSigners(&_Stakemanager.CallOpts, signer)
}

// SeqSigners is a free data retrieval call binding the contract method 0xbeb26755.
//
// Solidity: function seqSigners(address signer) view returns(uint256 seqId)
func (_Stakemanager *StakemanagerCallerSession) SeqSigners(signer common.Address) (*big.Int, error) {
	return _Stakemanager.Contract.SeqSigners(&_Stakemanager.CallOpts, signer)
}

// SeqStatuses is a free data retrieval call binding the contract method 0x86d203ab.
//
// Solidity: function seqStatuses(uint8 status) view returns(uint256 count)
func (_Stakemanager *StakemanagerCaller) SeqStatuses(opts *bind.CallOpts, status uint8) (*big.Int, error) {
	var out []interface{}
	err := _Stakemanager.contract.Call(opts, &out, "seqStatuses", status)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SeqStatuses is a free data retrieval call binding the contract method 0x86d203ab.
//
// Solidity: function seqStatuses(uint8 status) view returns(uint256 count)
func (_Stakemanager *StakemanagerSession) SeqStatuses(status uint8) (*big.Int, error) {
	return _Stakemanager.Contract.SeqStatuses(&_Stakemanager.CallOpts, status)
}

// SeqStatuses is a free data retrieval call binding the contract method 0x86d203ab.
//
// Solidity: function seqStatuses(uint8 status) view returns(uint256 count)
func (_Stakemanager *StakemanagerCallerSession) SeqStatuses(status uint8) (*big.Int, error) {
	return _Stakemanager.Contract.SeqStatuses(&_Stakemanager.CallOpts, status)
}

// Sequencers is a free data retrieval call binding the contract method 0x6ba7ccff.
//
// Solidity: function sequencers(uint256 seqId) view returns(uint256 amount, uint256 reward, uint256 activationBatch, uint256 updatedBatch, uint256 deactivationBatch, uint256 deactivationTime, uint256 unlockClaimTime, uint256 nonce, address owner, address signer, bytes pubkey, address rewardRecipient, uint8 status)
func (_Stakemanager *StakemanagerCaller) Sequencers(opts *bind.CallOpts, seqId *big.Int) (struct {
	Amount            *big.Int
	Reward            *big.Int
	ActivationBatch   *big.Int
	UpdatedBatch      *big.Int
	DeactivationBatch *big.Int
	DeactivationTime  *big.Int
	UnlockClaimTime   *big.Int
	Nonce             *big.Int
	Owner             common.Address
	Signer            common.Address
	Pubkey            []byte
	RewardRecipient   common.Address
	Status            uint8
}, error) {
	var out []interface{}
	err := _Stakemanager.contract.Call(opts, &out, "sequencers", seqId)

	outstruct := new(struct {
		Amount            *big.Int
		Reward            *big.Int
		ActivationBatch   *big.Int
		UpdatedBatch      *big.Int
		DeactivationBatch *big.Int
		DeactivationTime  *big.Int
		UnlockClaimTime   *big.Int
		Nonce             *big.Int
		Owner             common.Address
		Signer            common.Address
		Pubkey            []byte
		RewardRecipient   common.Address
		Status            uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Amount = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Reward = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.ActivationBatch = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UpdatedBatch = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.DeactivationBatch = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.DeactivationTime = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.UnlockClaimTime = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.Nonce = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.Owner = *abi.ConvertType(out[8], new(common.Address)).(*common.Address)
	outstruct.Signer = *abi.ConvertType(out[9], new(common.Address)).(*common.Address)
	outstruct.Pubkey = *abi.ConvertType(out[10], new([]byte)).(*[]byte)
	outstruct.RewardRecipient = *abi.ConvertType(out[11], new(common.Address)).(*common.Address)
	outstruct.Status = *abi.ConvertType(out[12], new(uint8)).(*uint8)

	return *outstruct, err

}

// Sequencers is a free data retrieval call binding the contract method 0x6ba7ccff.
//
// Solidity: function sequencers(uint256 seqId) view returns(uint256 amount, uint256 reward, uint256 activationBatch, uint256 updatedBatch, uint256 deactivationBatch, uint256 deactivationTime, uint256 unlockClaimTime, uint256 nonce, address owner, address signer, bytes pubkey, address rewardRecipient, uint8 status)
func (_Stakemanager *StakemanagerSession) Sequencers(seqId *big.Int) (struct {
	Amount            *big.Int
	Reward            *big.Int
	ActivationBatch   *big.Int
	UpdatedBatch      *big.Int
	DeactivationBatch *big.Int
	DeactivationTime  *big.Int
	UnlockClaimTime   *big.Int
	Nonce             *big.Int
	Owner             common.Address
	Signer            common.Address
	Pubkey            []byte
	RewardRecipient   common.Address
	Status            uint8
}, error) {
	return _Stakemanager.Contract.Sequencers(&_Stakemanager.CallOpts, seqId)
}

// Sequencers is a free data retrieval call binding the contract method 0x6ba7ccff.
//
// Solidity: function sequencers(uint256 seqId) view returns(uint256 amount, uint256 reward, uint256 activationBatch, uint256 updatedBatch, uint256 deactivationBatch, uint256 deactivationTime, uint256 unlockClaimTime, uint256 nonce, address owner, address signer, bytes pubkey, address rewardRecipient, uint8 status)
func (_Stakemanager *StakemanagerCallerSession) Sequencers(seqId *big.Int) (struct {
	Amount            *big.Int
	Reward            *big.Int
	ActivationBatch   *big.Int
	UpdatedBatch      *big.Int
	DeactivationBatch *big.Int
	DeactivationTime  *big.Int
	UnlockClaimTime   *big.Int
	Nonce             *big.Int
	Owner             common.Address
	Signer            common.Address
	Pubkey            []byte
	RewardRecipient   common.Address
	Status            uint8
}, error) {
	return _Stakemanager.Contract.Sequencers(&_Stakemanager.CallOpts, seqId)
}

// SignerUpdateThrottle is a free data retrieval call binding the contract method 0xc65066d4.
//
// Solidity: function signerUpdateThrottle() view returns(uint256)
func (_Stakemanager *StakemanagerCaller) SignerUpdateThrottle(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Stakemanager.contract.Call(opts, &out, "signerUpdateThrottle")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SignerUpdateThrottle is a free data retrieval call binding the contract method 0xc65066d4.
//
// Solidity: function signerUpdateThrottle() view returns(uint256)
func (_Stakemanager *StakemanagerSession) SignerUpdateThrottle() (*big.Int, error) {
	return _Stakemanager.Contract.SignerUpdateThrottle(&_Stakemanager.CallOpts)
}

// SignerUpdateThrottle is a free data retrieval call binding the contract method 0xc65066d4.
//
// Solidity: function signerUpdateThrottle() view returns(uint256)
func (_Stakemanager *StakemanagerCallerSession) SignerUpdateThrottle() (*big.Int, error) {
	return _Stakemanager.Contract.SignerUpdateThrottle(&_Stakemanager.CallOpts)
}

// TotalSequencers is a free data retrieval call binding the contract method 0xcc3ab923.
//
// Solidity: function totalSequencers() view returns(uint256)
func (_Stakemanager *StakemanagerCaller) TotalSequencers(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Stakemanager.contract.Call(opts, &out, "totalSequencers")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSequencers is a free data retrieval call binding the contract method 0xcc3ab923.
//
// Solidity: function totalSequencers() view returns(uint256)
func (_Stakemanager *StakemanagerSession) TotalSequencers() (*big.Int, error) {
	return _Stakemanager.Contract.TotalSequencers(&_Stakemanager.CallOpts)
}

// TotalSequencers is a free data retrieval call binding the contract method 0xcc3ab923.
//
// Solidity: function totalSequencers() view returns(uint256)
func (_Stakemanager *StakemanagerCallerSession) TotalSequencers() (*big.Int, error) {
	return _Stakemanager.Contract.TotalSequencers(&_Stakemanager.CallOpts)
}

// Whitelist is a free data retrieval call binding the contract method 0x9b19251a.
//
// Solidity: function whitelist(address ) view returns(bool)
func (_Stakemanager *StakemanagerCaller) Whitelist(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _Stakemanager.contract.Call(opts, &out, "whitelist", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Whitelist is a free data retrieval call binding the contract method 0x9b19251a.
//
// Solidity: function whitelist(address ) view returns(bool)
func (_Stakemanager *StakemanagerSession) Whitelist(arg0 common.Address) (bool, error) {
	return _Stakemanager.Contract.Whitelist(&_Stakemanager.CallOpts, arg0)
}

// Whitelist is a free data retrieval call binding the contract method 0x9b19251a.
//
// Solidity: function whitelist(address ) view returns(bool)
func (_Stakemanager *StakemanagerCallerSession) Whitelist(arg0 common.Address) (bool, error) {
	return _Stakemanager.Contract.Whitelist(&_Stakemanager.CallOpts, arg0)
}

// BatchSubmitRewards is a paid mutator transaction binding the contract method 0x11c7d144.
//
// Solidity: function batchSubmitRewards(uint256 _batchId, uint256 _startEpoch, uint256 _endEpoch, address[] _seqs, uint256[] _blocks) returns(uint256 totalReward)
func (_Stakemanager *StakemanagerTransactor) BatchSubmitRewards(opts *bind.TransactOpts, _batchId *big.Int, _startEpoch *big.Int, _endEpoch *big.Int, _seqs []common.Address, _blocks []*big.Int) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "batchSubmitRewards", _batchId, _startEpoch, _endEpoch, _seqs, _blocks)
}

// BatchSubmitRewards is a paid mutator transaction binding the contract method 0x11c7d144.
//
// Solidity: function batchSubmitRewards(uint256 _batchId, uint256 _startEpoch, uint256 _endEpoch, address[] _seqs, uint256[] _blocks) returns(uint256 totalReward)
func (_Stakemanager *StakemanagerSession) BatchSubmitRewards(_batchId *big.Int, _startEpoch *big.Int, _endEpoch *big.Int, _seqs []common.Address, _blocks []*big.Int) (*types.Transaction, error) {
	return _Stakemanager.Contract.BatchSubmitRewards(&_Stakemanager.TransactOpts, _batchId, _startEpoch, _endEpoch, _seqs, _blocks)
}

// BatchSubmitRewards is a paid mutator transaction binding the contract method 0x11c7d144.
//
// Solidity: function batchSubmitRewards(uint256 _batchId, uint256 _startEpoch, uint256 _endEpoch, address[] _seqs, uint256[] _blocks) returns(uint256 totalReward)
func (_Stakemanager *StakemanagerTransactorSession) BatchSubmitRewards(_batchId *big.Int, _startEpoch *big.Int, _endEpoch *big.Int, _seqs []common.Address, _blocks []*big.Int) (*types.Transaction, error) {
	return _Stakemanager.Contract.BatchSubmitRewards(&_Stakemanager.TransactOpts, _batchId, _startEpoch, _endEpoch, _seqs, _blocks)
}

// ForceUnlock is a paid mutator transaction binding the contract method 0xca99e838.
//
// Solidity: function forceUnlock(uint256 _seqId, uint32 _l2Gas) payable returns()
func (_Stakemanager *StakemanagerTransactor) ForceUnlock(opts *bind.TransactOpts, _seqId *big.Int, _l2Gas uint32) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "forceUnlock", _seqId, _l2Gas)
}

// ForceUnlock is a paid mutator transaction binding the contract method 0xca99e838.
//
// Solidity: function forceUnlock(uint256 _seqId, uint32 _l2Gas) payable returns()
func (_Stakemanager *StakemanagerSession) ForceUnlock(_seqId *big.Int, _l2Gas uint32) (*types.Transaction, error) {
	return _Stakemanager.Contract.ForceUnlock(&_Stakemanager.TransactOpts, _seqId, _l2Gas)
}

// ForceUnlock is a paid mutator transaction binding the contract method 0xca99e838.
//
// Solidity: function forceUnlock(uint256 _seqId, uint32 _l2Gas) payable returns()
func (_Stakemanager *StakemanagerTransactorSession) ForceUnlock(_seqId *big.Int, _l2Gas uint32) (*types.Transaction, error) {
	return _Stakemanager.Contract.ForceUnlock(&_Stakemanager.TransactOpts, _seqId, _l2Gas)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _escrow) returns()
func (_Stakemanager *StakemanagerTransactor) Initialize(opts *bind.TransactOpts, _escrow common.Address) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "initialize", _escrow)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _escrow) returns()
func (_Stakemanager *StakemanagerSession) Initialize(_escrow common.Address) (*types.Transaction, error) {
	return _Stakemanager.Contract.Initialize(&_Stakemanager.TransactOpts, _escrow)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _escrow) returns()
func (_Stakemanager *StakemanagerTransactorSession) Initialize(_escrow common.Address) (*types.Transaction, error) {
	return _Stakemanager.Contract.Initialize(&_Stakemanager.TransactOpts, _escrow)
}

// LockFor is a paid mutator transaction binding the contract method 0xaf70cba5.
//
// Solidity: function lockFor(address _signer, uint256 _amount, bytes _signerPubkey) returns()
func (_Stakemanager *StakemanagerTransactor) LockFor(opts *bind.TransactOpts, _signer common.Address, _amount *big.Int, _signerPubkey []byte) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "lockFor", _signer, _amount, _signerPubkey)
}

// LockFor is a paid mutator transaction binding the contract method 0xaf70cba5.
//
// Solidity: function lockFor(address _signer, uint256 _amount, bytes _signerPubkey) returns()
func (_Stakemanager *StakemanagerSession) LockFor(_signer common.Address, _amount *big.Int, _signerPubkey []byte) (*types.Transaction, error) {
	return _Stakemanager.Contract.LockFor(&_Stakemanager.TransactOpts, _signer, _amount, _signerPubkey)
}

// LockFor is a paid mutator transaction binding the contract method 0xaf70cba5.
//
// Solidity: function lockFor(address _signer, uint256 _amount, bytes _signerPubkey) returns()
func (_Stakemanager *StakemanagerTransactorSession) LockFor(_signer common.Address, _amount *big.Int, _signerPubkey []byte) (*types.Transaction, error) {
	return _Stakemanager.Contract.LockFor(&_Stakemanager.TransactOpts, _signer, _amount, _signerPubkey)
}

// LockWithRewardRecipient is a paid mutator transaction binding the contract method 0x9ad42030.
//
// Solidity: function lockWithRewardRecipient(address _signer, address _rewardRecipient, uint256 _amount, bytes _signerPubkey) returns()
func (_Stakemanager *StakemanagerTransactor) LockWithRewardRecipient(opts *bind.TransactOpts, _signer common.Address, _rewardRecipient common.Address, _amount *big.Int, _signerPubkey []byte) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "lockWithRewardRecipient", _signer, _rewardRecipient, _amount, _signerPubkey)
}

// LockWithRewardRecipient is a paid mutator transaction binding the contract method 0x9ad42030.
//
// Solidity: function lockWithRewardRecipient(address _signer, address _rewardRecipient, uint256 _amount, bytes _signerPubkey) returns()
func (_Stakemanager *StakemanagerSession) LockWithRewardRecipient(_signer common.Address, _rewardRecipient common.Address, _amount *big.Int, _signerPubkey []byte) (*types.Transaction, error) {
	return _Stakemanager.Contract.LockWithRewardRecipient(&_Stakemanager.TransactOpts, _signer, _rewardRecipient, _amount, _signerPubkey)
}

// LockWithRewardRecipient is a paid mutator transaction binding the contract method 0x9ad42030.
//
// Solidity: function lockWithRewardRecipient(address _signer, address _rewardRecipient, uint256 _amount, bytes _signerPubkey) returns()
func (_Stakemanager *StakemanagerTransactorSession) LockWithRewardRecipient(_signer common.Address, _rewardRecipient common.Address, _amount *big.Int, _signerPubkey []byte) (*types.Transaction, error) {
	return _Stakemanager.Contract.LockWithRewardRecipient(&_Stakemanager.TransactOpts, _signer, _rewardRecipient, _amount, _signerPubkey)
}

// Relock is a paid mutator transaction binding the contract method 0x015bb180.
//
// Solidity: function relock(uint256 _seqId, uint256 _amount, bool _lockReward) returns()
func (_Stakemanager *StakemanagerTransactor) Relock(opts *bind.TransactOpts, _seqId *big.Int, _amount *big.Int, _lockReward bool) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "relock", _seqId, _amount, _lockReward)
}

// Relock is a paid mutator transaction binding the contract method 0x015bb180.
//
// Solidity: function relock(uint256 _seqId, uint256 _amount, bool _lockReward) returns()
func (_Stakemanager *StakemanagerSession) Relock(_seqId *big.Int, _amount *big.Int, _lockReward bool) (*types.Transaction, error) {
	return _Stakemanager.Contract.Relock(&_Stakemanager.TransactOpts, _seqId, _amount, _lockReward)
}

// Relock is a paid mutator transaction binding the contract method 0x015bb180.
//
// Solidity: function relock(uint256 _seqId, uint256 _amount, bool _lockReward) returns()
func (_Stakemanager *StakemanagerTransactorSession) Relock(_seqId *big.Int, _amount *big.Int, _lockReward bool) (*types.Transaction, error) {
	return _Stakemanager.Contract.Relock(&_Stakemanager.TransactOpts, _seqId, _amount, _lockReward)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Stakemanager *StakemanagerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Stakemanager *StakemanagerSession) RenounceOwnership() (*types.Transaction, error) {
	return _Stakemanager.Contract.RenounceOwnership(&_Stakemanager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Stakemanager *StakemanagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Stakemanager.Contract.RenounceOwnership(&_Stakemanager.TransactOpts)
}

// SetPause is a paid mutator transaction binding the contract method 0xbedb86fb.
//
// Solidity: function setPause(bool _yes) returns()
func (_Stakemanager *StakemanagerTransactor) SetPause(opts *bind.TransactOpts, _yes bool) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "setPause", _yes)
}

// SetPause is a paid mutator transaction binding the contract method 0xbedb86fb.
//
// Solidity: function setPause(bool _yes) returns()
func (_Stakemanager *StakemanagerSession) SetPause(_yes bool) (*types.Transaction, error) {
	return _Stakemanager.Contract.SetPause(&_Stakemanager.TransactOpts, _yes)
}

// SetPause is a paid mutator transaction binding the contract method 0xbedb86fb.
//
// Solidity: function setPause(bool _yes) returns()
func (_Stakemanager *StakemanagerTransactorSession) SetPause(_yes bool) (*types.Transaction, error) {
	return _Stakemanager.Contract.SetPause(&_Stakemanager.TransactOpts, _yes)
}

// SetSequencerOwner is a paid mutator transaction binding the contract method 0xa953791f.
//
// Solidity: function setSequencerOwner(uint256 _seqId, address _owner) returns()
func (_Stakemanager *StakemanagerTransactor) SetSequencerOwner(opts *bind.TransactOpts, _seqId *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "setSequencerOwner", _seqId, _owner)
}

// SetSequencerOwner is a paid mutator transaction binding the contract method 0xa953791f.
//
// Solidity: function setSequencerOwner(uint256 _seqId, address _owner) returns()
func (_Stakemanager *StakemanagerSession) SetSequencerOwner(_seqId *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _Stakemanager.Contract.SetSequencerOwner(&_Stakemanager.TransactOpts, _seqId, _owner)
}

// SetSequencerOwner is a paid mutator transaction binding the contract method 0xa953791f.
//
// Solidity: function setSequencerOwner(uint256 _seqId, address _owner) returns()
func (_Stakemanager *StakemanagerTransactorSession) SetSequencerOwner(_seqId *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _Stakemanager.Contract.SetSequencerOwner(&_Stakemanager.TransactOpts, _seqId, _owner)
}

// SetSequencerRewardRecipient is a paid mutator transaction binding the contract method 0xd83b0e14.
//
// Solidity: function setSequencerRewardRecipient(uint256 _seqId, address _recipient) returns()
func (_Stakemanager *StakemanagerTransactor) SetSequencerRewardRecipient(opts *bind.TransactOpts, _seqId *big.Int, _recipient common.Address) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "setSequencerRewardRecipient", _seqId, _recipient)
}

// SetSequencerRewardRecipient is a paid mutator transaction binding the contract method 0xd83b0e14.
//
// Solidity: function setSequencerRewardRecipient(uint256 _seqId, address _recipient) returns()
func (_Stakemanager *StakemanagerSession) SetSequencerRewardRecipient(_seqId *big.Int, _recipient common.Address) (*types.Transaction, error) {
	return _Stakemanager.Contract.SetSequencerRewardRecipient(&_Stakemanager.TransactOpts, _seqId, _recipient)
}

// SetSequencerRewardRecipient is a paid mutator transaction binding the contract method 0xd83b0e14.
//
// Solidity: function setSequencerRewardRecipient(uint256 _seqId, address _recipient) returns()
func (_Stakemanager *StakemanagerTransactorSession) SetSequencerRewardRecipient(_seqId *big.Int, _recipient common.Address) (*types.Transaction, error) {
	return _Stakemanager.Contract.SetSequencerRewardRecipient(&_Stakemanager.TransactOpts, _seqId, _recipient)
}

// SetSignerUpdateThrottle is a paid mutator transaction binding the contract method 0xbfd6fc3f.
//
// Solidity: function setSignerUpdateThrottle(uint256 _n) returns()
func (_Stakemanager *StakemanagerTransactor) SetSignerUpdateThrottle(opts *bind.TransactOpts, _n *big.Int) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "setSignerUpdateThrottle", _n)
}

// SetSignerUpdateThrottle is a paid mutator transaction binding the contract method 0xbfd6fc3f.
//
// Solidity: function setSignerUpdateThrottle(uint256 _n) returns()
func (_Stakemanager *StakemanagerSession) SetSignerUpdateThrottle(_n *big.Int) (*types.Transaction, error) {
	return _Stakemanager.Contract.SetSignerUpdateThrottle(&_Stakemanager.TransactOpts, _n)
}

// SetSignerUpdateThrottle is a paid mutator transaction binding the contract method 0xbfd6fc3f.
//
// Solidity: function setSignerUpdateThrottle(uint256 _n) returns()
func (_Stakemanager *StakemanagerTransactorSession) SetSignerUpdateThrottle(_n *big.Int) (*types.Transaction, error) {
	return _Stakemanager.Contract.SetSignerUpdateThrottle(&_Stakemanager.TransactOpts, _n)
}

// SetWhitelist is a paid mutator transaction binding the contract method 0x53d6fd59.
//
// Solidity: function setWhitelist(address _addr, bool _yes) returns()
func (_Stakemanager *StakemanagerTransactor) SetWhitelist(opts *bind.TransactOpts, _addr common.Address, _yes bool) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "setWhitelist", _addr, _yes)
}

// SetWhitelist is a paid mutator transaction binding the contract method 0x53d6fd59.
//
// Solidity: function setWhitelist(address _addr, bool _yes) returns()
func (_Stakemanager *StakemanagerSession) SetWhitelist(_addr common.Address, _yes bool) (*types.Transaction, error) {
	return _Stakemanager.Contract.SetWhitelist(&_Stakemanager.TransactOpts, _addr, _yes)
}

// SetWhitelist is a paid mutator transaction binding the contract method 0x53d6fd59.
//
// Solidity: function setWhitelist(address _addr, bool _yes) returns()
func (_Stakemanager *StakemanagerTransactorSession) SetWhitelist(_addr common.Address, _yes bool) (*types.Transaction, error) {
	return _Stakemanager.Contract.SetWhitelist(&_Stakemanager.TransactOpts, _addr, _yes)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Stakemanager *StakemanagerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Stakemanager *StakemanagerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Stakemanager.Contract.TransferOwnership(&_Stakemanager.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Stakemanager *StakemanagerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Stakemanager.Contract.TransferOwnership(&_Stakemanager.TransactOpts, newOwner)
}

// Unlock is a paid mutator transaction binding the contract method 0x262c0e66.
//
// Solidity: function unlock(uint256 _seqId, uint32 _l2Gas) payable returns()
func (_Stakemanager *StakemanagerTransactor) Unlock(opts *bind.TransactOpts, _seqId *big.Int, _l2Gas uint32) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "unlock", _seqId, _l2Gas)
}

// Unlock is a paid mutator transaction binding the contract method 0x262c0e66.
//
// Solidity: function unlock(uint256 _seqId, uint32 _l2Gas) payable returns()
func (_Stakemanager *StakemanagerSession) Unlock(_seqId *big.Int, _l2Gas uint32) (*types.Transaction, error) {
	return _Stakemanager.Contract.Unlock(&_Stakemanager.TransactOpts, _seqId, _l2Gas)
}

// Unlock is a paid mutator transaction binding the contract method 0x262c0e66.
//
// Solidity: function unlock(uint256 _seqId, uint32 _l2Gas) payable returns()
func (_Stakemanager *StakemanagerTransactorSession) Unlock(_seqId *big.Int, _l2Gas uint32) (*types.Transaction, error) {
	return _Stakemanager.Contract.Unlock(&_Stakemanager.TransactOpts, _seqId, _l2Gas)
}

// UnlockClaim is a paid mutator transaction binding the contract method 0x8ddc74de.
//
// Solidity: function unlockClaim(uint256 _seqId, uint32 _l2Gas) payable returns()
func (_Stakemanager *StakemanagerTransactor) UnlockClaim(opts *bind.TransactOpts, _seqId *big.Int, _l2Gas uint32) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "unlockClaim", _seqId, _l2Gas)
}

// UnlockClaim is a paid mutator transaction binding the contract method 0x8ddc74de.
//
// Solidity: function unlockClaim(uint256 _seqId, uint32 _l2Gas) payable returns()
func (_Stakemanager *StakemanagerSession) UnlockClaim(_seqId *big.Int, _l2Gas uint32) (*types.Transaction, error) {
	return _Stakemanager.Contract.UnlockClaim(&_Stakemanager.TransactOpts, _seqId, _l2Gas)
}

// UnlockClaim is a paid mutator transaction binding the contract method 0x8ddc74de.
//
// Solidity: function unlockClaim(uint256 _seqId, uint32 _l2Gas) payable returns()
func (_Stakemanager *StakemanagerTransactorSession) UnlockClaim(_seqId *big.Int, _l2Gas uint32) (*types.Transaction, error) {
	return _Stakemanager.Contract.UnlockClaim(&_Stakemanager.TransactOpts, _seqId, _l2Gas)
}

// UpdateBlockReward is a paid mutator transaction binding the contract method 0xf580ffcb.
//
// Solidity: function updateBlockReward(uint256 newReward) returns()
func (_Stakemanager *StakemanagerTransactor) UpdateBlockReward(opts *bind.TransactOpts, newReward *big.Int) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "updateBlockReward", newReward)
}

// UpdateBlockReward is a paid mutator transaction binding the contract method 0xf580ffcb.
//
// Solidity: function updateBlockReward(uint256 newReward) returns()
func (_Stakemanager *StakemanagerSession) UpdateBlockReward(newReward *big.Int) (*types.Transaction, error) {
	return _Stakemanager.Contract.UpdateBlockReward(&_Stakemanager.TransactOpts, newReward)
}

// UpdateBlockReward is a paid mutator transaction binding the contract method 0xf580ffcb.
//
// Solidity: function updateBlockReward(uint256 newReward) returns()
func (_Stakemanager *StakemanagerTransactorSession) UpdateBlockReward(newReward *big.Int) (*types.Transaction, error) {
	return _Stakemanager.Contract.UpdateBlockReward(&_Stakemanager.TransactOpts, newReward)
}

// UpdateMpc is a paid mutator transaction binding the contract method 0xd11d0681.
//
// Solidity: function updateMpc(address _newMpc) returns()
func (_Stakemanager *StakemanagerTransactor) UpdateMpc(opts *bind.TransactOpts, _newMpc common.Address) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "updateMpc", _newMpc)
}

// UpdateMpc is a paid mutator transaction binding the contract method 0xd11d0681.
//
// Solidity: function updateMpc(address _newMpc) returns()
func (_Stakemanager *StakemanagerSession) UpdateMpc(_newMpc common.Address) (*types.Transaction, error) {
	return _Stakemanager.Contract.UpdateMpc(&_Stakemanager.TransactOpts, _newMpc)
}

// UpdateMpc is a paid mutator transaction binding the contract method 0xd11d0681.
//
// Solidity: function updateMpc(address _newMpc) returns()
func (_Stakemanager *StakemanagerTransactorSession) UpdateMpc(_newMpc common.Address) (*types.Transaction, error) {
	return _Stakemanager.Contract.UpdateMpc(&_Stakemanager.TransactOpts, _newMpc)
}

// UpdateSigner is a paid mutator transaction binding the contract method 0xf41a9642.
//
// Solidity: function updateSigner(uint256 _seqId, bytes _signerPubkey) returns()
func (_Stakemanager *StakemanagerTransactor) UpdateSigner(opts *bind.TransactOpts, _seqId *big.Int, _signerPubkey []byte) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "updateSigner", _seqId, _signerPubkey)
}

// UpdateSigner is a paid mutator transaction binding the contract method 0xf41a9642.
//
// Solidity: function updateSigner(uint256 _seqId, bytes _signerPubkey) returns()
func (_Stakemanager *StakemanagerSession) UpdateSigner(_seqId *big.Int, _signerPubkey []byte) (*types.Transaction, error) {
	return _Stakemanager.Contract.UpdateSigner(&_Stakemanager.TransactOpts, _seqId, _signerPubkey)
}

// UpdateSigner is a paid mutator transaction binding the contract method 0xf41a9642.
//
// Solidity: function updateSigner(uint256 _seqId, bytes _signerPubkey) returns()
func (_Stakemanager *StakemanagerTransactorSession) UpdateSigner(_seqId *big.Int, _signerPubkey []byte) (*types.Transaction, error) {
	return _Stakemanager.Contract.UpdateSigner(&_Stakemanager.TransactOpts, _seqId, _signerPubkey)
}

// UpdateWithdrawDelayTimeValue is a paid mutator transaction binding the contract method 0x71e10cfa.
//
// Solidity: function updateWithdrawDelayTimeValue(uint256 _time) returns()
func (_Stakemanager *StakemanagerTransactor) UpdateWithdrawDelayTimeValue(opts *bind.TransactOpts, _time *big.Int) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "updateWithdrawDelayTimeValue", _time)
}

// UpdateWithdrawDelayTimeValue is a paid mutator transaction binding the contract method 0x71e10cfa.
//
// Solidity: function updateWithdrawDelayTimeValue(uint256 _time) returns()
func (_Stakemanager *StakemanagerSession) UpdateWithdrawDelayTimeValue(_time *big.Int) (*types.Transaction, error) {
	return _Stakemanager.Contract.UpdateWithdrawDelayTimeValue(&_Stakemanager.TransactOpts, _time)
}

// UpdateWithdrawDelayTimeValue is a paid mutator transaction binding the contract method 0x71e10cfa.
//
// Solidity: function updateWithdrawDelayTimeValue(uint256 _time) returns()
func (_Stakemanager *StakemanagerTransactorSession) UpdateWithdrawDelayTimeValue(_time *big.Int) (*types.Transaction, error) {
	return _Stakemanager.Contract.UpdateWithdrawDelayTimeValue(&_Stakemanager.TransactOpts, _time)
}

// WithdrawRewards is a paid mutator transaction binding the contract method 0x17396687.
//
// Solidity: function withdrawRewards(uint256 _seqId, uint32 _l2Gas) payable returns()
func (_Stakemanager *StakemanagerTransactor) WithdrawRewards(opts *bind.TransactOpts, _seqId *big.Int, _l2Gas uint32) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "withdrawRewards", _seqId, _l2Gas)
}

// WithdrawRewards is a paid mutator transaction binding the contract method 0x17396687.
//
// Solidity: function withdrawRewards(uint256 _seqId, uint32 _l2Gas) payable returns()
func (_Stakemanager *StakemanagerSession) WithdrawRewards(_seqId *big.Int, _l2Gas uint32) (*types.Transaction, error) {
	return _Stakemanager.Contract.WithdrawRewards(&_Stakemanager.TransactOpts, _seqId, _l2Gas)
}

// WithdrawRewards is a paid mutator transaction binding the contract method 0x17396687.
//
// Solidity: function withdrawRewards(uint256 _seqId, uint32 _l2Gas) payable returns()
func (_Stakemanager *StakemanagerTransactorSession) WithdrawRewards(_seqId *big.Int, _l2Gas uint32) (*types.Transaction, error) {
	return _Stakemanager.Contract.WithdrawRewards(&_Stakemanager.TransactOpts, _seqId, _l2Gas)
}

// StakemanagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Stakemanager contract.
type StakemanagerInitializedIterator struct {
	Event *StakemanagerInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakemanagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakemanagerInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakemanagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerInitialized represents a Initialized event raised by the Stakemanager contract.
type StakemanagerInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Stakemanager *StakemanagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*StakemanagerInitializedIterator, error) {

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &StakemanagerInitializedIterator{contract: _Stakemanager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Stakemanager *StakemanagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *StakemanagerInitialized) (event.Subscription, error) {

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerInitialized)
				if err := _Stakemanager.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Stakemanager *StakemanagerFilterer) ParseInitialized(log types.Log) (*StakemanagerInitialized, error) {
	event := new(StakemanagerInitialized)
	if err := _Stakemanager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakemanagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Stakemanager contract.
type StakemanagerOwnershipTransferredIterator struct {
	Event *StakemanagerOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakemanagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakemanagerOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakemanagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerOwnershipTransferred represents a OwnershipTransferred event raised by the Stakemanager contract.
type StakemanagerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Stakemanager *StakemanagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*StakemanagerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &StakemanagerOwnershipTransferredIterator{contract: _Stakemanager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Stakemanager *StakemanagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *StakemanagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerOwnershipTransferred)
				if err := _Stakemanager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Stakemanager *StakemanagerFilterer) ParseOwnershipTransferred(log types.Log) (*StakemanagerOwnershipTransferred, error) {
	event := new(StakemanagerOwnershipTransferred)
	if err := _Stakemanager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakemanagerPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Stakemanager contract.
type StakemanagerPausedIterator struct {
	Event *StakemanagerPaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakemanagerPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakemanagerPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakemanagerPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerPaused represents a Paused event raised by the Stakemanager contract.
type StakemanagerPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Stakemanager *StakemanagerFilterer) FilterPaused(opts *bind.FilterOpts) (*StakemanagerPausedIterator, error) {

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &StakemanagerPausedIterator{contract: _Stakemanager.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Stakemanager *StakemanagerFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *StakemanagerPaused) (event.Subscription, error) {

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerPaused)
				if err := _Stakemanager.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Stakemanager *StakemanagerFilterer) ParsePaused(log types.Log) (*StakemanagerPaused, error) {
	event := new(StakemanagerPaused)
	if err := _Stakemanager.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakemanagerRewardUpdateIterator is returned from FilterRewardUpdate and is used to iterate over the raw logs and unpacked data for RewardUpdate events raised by the Stakemanager contract.
type StakemanagerRewardUpdateIterator struct {
	Event *StakemanagerRewardUpdate // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakemanagerRewardUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerRewardUpdate)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakemanagerRewardUpdate)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakemanagerRewardUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerRewardUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerRewardUpdate represents a RewardUpdate event raised by the Stakemanager contract.
type StakemanagerRewardUpdate struct {
	NewReward *big.Int
	OldReward *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRewardUpdate is a free log retrieval operation binding the contract event 0xf67f33e8589d3ea0356303c0f9a8e764873692159f777ff79e4fc523d389dfcd.
//
// Solidity: event RewardUpdate(uint256 newReward, uint256 oldReward)
func (_Stakemanager *StakemanagerFilterer) FilterRewardUpdate(opts *bind.FilterOpts) (*StakemanagerRewardUpdateIterator, error) {

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "RewardUpdate")
	if err != nil {
		return nil, err
	}
	return &StakemanagerRewardUpdateIterator{contract: _Stakemanager.contract, event: "RewardUpdate", logs: logs, sub: sub}, nil
}

// WatchRewardUpdate is a free log subscription operation binding the contract event 0xf67f33e8589d3ea0356303c0f9a8e764873692159f777ff79e4fc523d389dfcd.
//
// Solidity: event RewardUpdate(uint256 newReward, uint256 oldReward)
func (_Stakemanager *StakemanagerFilterer) WatchRewardUpdate(opts *bind.WatchOpts, sink chan<- *StakemanagerRewardUpdate) (event.Subscription, error) {

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "RewardUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerRewardUpdate)
				if err := _Stakemanager.contract.UnpackLog(event, "RewardUpdate", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRewardUpdate is a log parse operation binding the contract event 0xf67f33e8589d3ea0356303c0f9a8e764873692159f777ff79e4fc523d389dfcd.
//
// Solidity: event RewardUpdate(uint256 newReward, uint256 oldReward)
func (_Stakemanager *StakemanagerFilterer) ParseRewardUpdate(log types.Log) (*StakemanagerRewardUpdate, error) {
	event := new(StakemanagerRewardUpdate)
	if err := _Stakemanager.contract.UnpackLog(event, "RewardUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakemanagerSequencerOwnerChangedIterator is returned from FilterSequencerOwnerChanged and is used to iterate over the raw logs and unpacked data for SequencerOwnerChanged events raised by the Stakemanager contract.
type StakemanagerSequencerOwnerChangedIterator struct {
	Event *StakemanagerSequencerOwnerChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakemanagerSequencerOwnerChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerSequencerOwnerChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakemanagerSequencerOwnerChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakemanagerSequencerOwnerChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerSequencerOwnerChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerSequencerOwnerChanged represents a SequencerOwnerChanged event raised by the Stakemanager contract.
type StakemanagerSequencerOwnerChanged struct {
	SeqId *big.Int
	Owner common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterSequencerOwnerChanged is a free log retrieval operation binding the contract event 0x4078101d3657bf2f1ee46f64d5c94266d244d71bb0daa460336d3d6f11c9a4ac.
//
// Solidity: event SequencerOwnerChanged(uint256 _seqId, address _owner)
func (_Stakemanager *StakemanagerFilterer) FilterSequencerOwnerChanged(opts *bind.FilterOpts) (*StakemanagerSequencerOwnerChangedIterator, error) {

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "SequencerOwnerChanged")
	if err != nil {
		return nil, err
	}
	return &StakemanagerSequencerOwnerChangedIterator{contract: _Stakemanager.contract, event: "SequencerOwnerChanged", logs: logs, sub: sub}, nil
}

// WatchSequencerOwnerChanged is a free log subscription operation binding the contract event 0x4078101d3657bf2f1ee46f64d5c94266d244d71bb0daa460336d3d6f11c9a4ac.
//
// Solidity: event SequencerOwnerChanged(uint256 _seqId, address _owner)
func (_Stakemanager *StakemanagerFilterer) WatchSequencerOwnerChanged(opts *bind.WatchOpts, sink chan<- *StakemanagerSequencerOwnerChanged) (event.Subscription, error) {

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "SequencerOwnerChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerSequencerOwnerChanged)
				if err := _Stakemanager.contract.UnpackLog(event, "SequencerOwnerChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSequencerOwnerChanged is a log parse operation binding the contract event 0x4078101d3657bf2f1ee46f64d5c94266d244d71bb0daa460336d3d6f11c9a4ac.
//
// Solidity: event SequencerOwnerChanged(uint256 _seqId, address _owner)
func (_Stakemanager *StakemanagerFilterer) ParseSequencerOwnerChanged(log types.Log) (*StakemanagerSequencerOwnerChanged, error) {
	event := new(StakemanagerSequencerOwnerChanged)
	if err := _Stakemanager.contract.UnpackLog(event, "SequencerOwnerChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakemanagerSequencerRewardRecipientChangedIterator is returned from FilterSequencerRewardRecipientChanged and is used to iterate over the raw logs and unpacked data for SequencerRewardRecipientChanged events raised by the Stakemanager contract.
type StakemanagerSequencerRewardRecipientChangedIterator struct {
	Event *StakemanagerSequencerRewardRecipientChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakemanagerSequencerRewardRecipientChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerSequencerRewardRecipientChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakemanagerSequencerRewardRecipientChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakemanagerSequencerRewardRecipientChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerSequencerRewardRecipientChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerSequencerRewardRecipientChanged represents a SequencerRewardRecipientChanged event raised by the Stakemanager contract.
type StakemanagerSequencerRewardRecipientChanged struct {
	SeqId     *big.Int
	Recipient common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSequencerRewardRecipientChanged is a free log retrieval operation binding the contract event 0x357bb123cabaf224688e3d8de5feb37d685dc3a27012a7bce1201c49bc369cb6.
//
// Solidity: event SequencerRewardRecipientChanged(uint256 _seqId, address _recipient)
func (_Stakemanager *StakemanagerFilterer) FilterSequencerRewardRecipientChanged(opts *bind.FilterOpts) (*StakemanagerSequencerRewardRecipientChangedIterator, error) {

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "SequencerRewardRecipientChanged")
	if err != nil {
		return nil, err
	}
	return &StakemanagerSequencerRewardRecipientChangedIterator{contract: _Stakemanager.contract, event: "SequencerRewardRecipientChanged", logs: logs, sub: sub}, nil
}

// WatchSequencerRewardRecipientChanged is a free log subscription operation binding the contract event 0x357bb123cabaf224688e3d8de5feb37d685dc3a27012a7bce1201c49bc369cb6.
//
// Solidity: event SequencerRewardRecipientChanged(uint256 _seqId, address _recipient)
func (_Stakemanager *StakemanagerFilterer) WatchSequencerRewardRecipientChanged(opts *bind.WatchOpts, sink chan<- *StakemanagerSequencerRewardRecipientChanged) (event.Subscription, error) {

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "SequencerRewardRecipientChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerSequencerRewardRecipientChanged)
				if err := _Stakemanager.contract.UnpackLog(event, "SequencerRewardRecipientChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSequencerRewardRecipientChanged is a log parse operation binding the contract event 0x357bb123cabaf224688e3d8de5feb37d685dc3a27012a7bce1201c49bc369cb6.
//
// Solidity: event SequencerRewardRecipientChanged(uint256 _seqId, address _recipient)
func (_Stakemanager *StakemanagerFilterer) ParseSequencerRewardRecipientChanged(log types.Log) (*StakemanagerSequencerRewardRecipientChanged, error) {
	event := new(StakemanagerSequencerRewardRecipientChanged)
	if err := _Stakemanager.contract.UnpackLog(event, "SequencerRewardRecipientChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakemanagerSetSignerUpdateThrottleIterator is returned from FilterSetSignerUpdateThrottle and is used to iterate over the raw logs and unpacked data for SetSignerUpdateThrottle events raised by the Stakemanager contract.
type StakemanagerSetSignerUpdateThrottleIterator struct {
	Event *StakemanagerSetSignerUpdateThrottle // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakemanagerSetSignerUpdateThrottleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerSetSignerUpdateThrottle)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakemanagerSetSignerUpdateThrottle)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakemanagerSetSignerUpdateThrottleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerSetSignerUpdateThrottleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerSetSignerUpdateThrottle represents a SetSignerUpdateThrottle event raised by the Stakemanager contract.
type StakemanagerSetSignerUpdateThrottle struct {
	N   *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterSetSignerUpdateThrottle is a free log retrieval operation binding the contract event 0xe58685f6b78e6d567d2ed9d7c5efb779c4cd91c63c76763a0e3204a5671f4705.
//
// Solidity: event SetSignerUpdateThrottle(uint256 _n)
func (_Stakemanager *StakemanagerFilterer) FilterSetSignerUpdateThrottle(opts *bind.FilterOpts) (*StakemanagerSetSignerUpdateThrottleIterator, error) {

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "SetSignerUpdateThrottle")
	if err != nil {
		return nil, err
	}
	return &StakemanagerSetSignerUpdateThrottleIterator{contract: _Stakemanager.contract, event: "SetSignerUpdateThrottle", logs: logs, sub: sub}, nil
}

// WatchSetSignerUpdateThrottle is a free log subscription operation binding the contract event 0xe58685f6b78e6d567d2ed9d7c5efb779c4cd91c63c76763a0e3204a5671f4705.
//
// Solidity: event SetSignerUpdateThrottle(uint256 _n)
func (_Stakemanager *StakemanagerFilterer) WatchSetSignerUpdateThrottle(opts *bind.WatchOpts, sink chan<- *StakemanagerSetSignerUpdateThrottle) (event.Subscription, error) {

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "SetSignerUpdateThrottle")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerSetSignerUpdateThrottle)
				if err := _Stakemanager.contract.UnpackLog(event, "SetSignerUpdateThrottle", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSetSignerUpdateThrottle is a log parse operation binding the contract event 0xe58685f6b78e6d567d2ed9d7c5efb779c4cd91c63c76763a0e3204a5671f4705.
//
// Solidity: event SetSignerUpdateThrottle(uint256 _n)
func (_Stakemanager *StakemanagerFilterer) ParseSetSignerUpdateThrottle(log types.Log) (*StakemanagerSetSignerUpdateThrottle, error) {
	event := new(StakemanagerSetSignerUpdateThrottle)
	if err := _Stakemanager.contract.UnpackLog(event, "SetSignerUpdateThrottle", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakemanagerSetWhitelistIterator is returned from FilterSetWhitelist and is used to iterate over the raw logs and unpacked data for SetWhitelist events raised by the Stakemanager contract.
type StakemanagerSetWhitelistIterator struct {
	Event *StakemanagerSetWhitelist // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakemanagerSetWhitelistIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerSetWhitelist)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakemanagerSetWhitelist)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakemanagerSetWhitelistIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerSetWhitelistIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerSetWhitelist represents a SetWhitelist event raised by the Stakemanager contract.
type StakemanagerSetWhitelist struct {
	User common.Address
	Yes  bool
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterSetWhitelist is a free log retrieval operation binding the contract event 0xf6019ec0a78d156d249a1ec7579e2321f6ac7521d6e1d2eacf90ba4a184dcceb.
//
// Solidity: event SetWhitelist(address _user, bool _yes)
func (_Stakemanager *StakemanagerFilterer) FilterSetWhitelist(opts *bind.FilterOpts) (*StakemanagerSetWhitelistIterator, error) {

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "SetWhitelist")
	if err != nil {
		return nil, err
	}
	return &StakemanagerSetWhitelistIterator{contract: _Stakemanager.contract, event: "SetWhitelist", logs: logs, sub: sub}, nil
}

// WatchSetWhitelist is a free log subscription operation binding the contract event 0xf6019ec0a78d156d249a1ec7579e2321f6ac7521d6e1d2eacf90ba4a184dcceb.
//
// Solidity: event SetWhitelist(address _user, bool _yes)
func (_Stakemanager *StakemanagerFilterer) WatchSetWhitelist(opts *bind.WatchOpts, sink chan<- *StakemanagerSetWhitelist) (event.Subscription, error) {

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "SetWhitelist")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerSetWhitelist)
				if err := _Stakemanager.contract.UnpackLog(event, "SetWhitelist", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSetWhitelist is a log parse operation binding the contract event 0xf6019ec0a78d156d249a1ec7579e2321f6ac7521d6e1d2eacf90ba4a184dcceb.
//
// Solidity: event SetWhitelist(address _user, bool _yes)
func (_Stakemanager *StakemanagerFilterer) ParseSetWhitelist(log types.Log) (*StakemanagerSetWhitelist, error) {
	event := new(StakemanagerSetWhitelist)
	if err := _Stakemanager.contract.UnpackLog(event, "SetWhitelist", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakemanagerUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Stakemanager contract.
type StakemanagerUnpausedIterator struct {
	Event *StakemanagerUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakemanagerUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakemanagerUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakemanagerUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerUnpaused represents a Unpaused event raised by the Stakemanager contract.
type StakemanagerUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Stakemanager *StakemanagerFilterer) FilterUnpaused(opts *bind.FilterOpts) (*StakemanagerUnpausedIterator, error) {

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &StakemanagerUnpausedIterator{contract: _Stakemanager.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Stakemanager *StakemanagerFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *StakemanagerUnpaused) (event.Subscription, error) {

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerUnpaused)
				if err := _Stakemanager.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Stakemanager *StakemanagerFilterer) ParseUnpaused(log types.Log) (*StakemanagerUnpaused, error) {
	event := new(StakemanagerUnpaused)
	if err := _Stakemanager.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakemanagerUpdateMpcIterator is returned from FilterUpdateMpc and is used to iterate over the raw logs and unpacked data for UpdateMpc events raised by the Stakemanager contract.
type StakemanagerUpdateMpcIterator struct {
	Event *StakemanagerUpdateMpc // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakemanagerUpdateMpcIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerUpdateMpc)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakemanagerUpdateMpc)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakemanagerUpdateMpcIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerUpdateMpcIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerUpdateMpc represents a UpdateMpc event raised by the Stakemanager contract.
type StakemanagerUpdateMpc struct {
	NewMpc common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterUpdateMpc is a free log retrieval operation binding the contract event 0xc6759872346bc2093e270f2fa00d99d7bcdde60a410a3e9956b34196d42fee76.
//
// Solidity: event UpdateMpc(address _newMpc)
func (_Stakemanager *StakemanagerFilterer) FilterUpdateMpc(opts *bind.FilterOpts) (*StakemanagerUpdateMpcIterator, error) {

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "UpdateMpc")
	if err != nil {
		return nil, err
	}
	return &StakemanagerUpdateMpcIterator{contract: _Stakemanager.contract, event: "UpdateMpc", logs: logs, sub: sub}, nil
}

// WatchUpdateMpc is a free log subscription operation binding the contract event 0xc6759872346bc2093e270f2fa00d99d7bcdde60a410a3e9956b34196d42fee76.
//
// Solidity: event UpdateMpc(address _newMpc)
func (_Stakemanager *StakemanagerFilterer) WatchUpdateMpc(opts *bind.WatchOpts, sink chan<- *StakemanagerUpdateMpc) (event.Subscription, error) {

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "UpdateMpc")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerUpdateMpc)
				if err := _Stakemanager.contract.UnpackLog(event, "UpdateMpc", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpdateMpc is a log parse operation binding the contract event 0xc6759872346bc2093e270f2fa00d99d7bcdde60a410a3e9956b34196d42fee76.
//
// Solidity: event UpdateMpc(address _newMpc)
func (_Stakemanager *StakemanagerFilterer) ParseUpdateMpc(log types.Log) (*StakemanagerUpdateMpc, error) {
	event := new(StakemanagerUpdateMpc)
	if err := _Stakemanager.contract.UnpackLog(event, "UpdateMpc", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakemanagerWithrawDelayTimeChangeIterator is returned from FilterWithrawDelayTimeChange and is used to iterate over the raw logs and unpacked data for WithrawDelayTimeChange events raised by the Stakemanager contract.
type StakemanagerWithrawDelayTimeChangeIterator struct {
	Event *StakemanagerWithrawDelayTimeChange // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakemanagerWithrawDelayTimeChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerWithrawDelayTimeChange)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakemanagerWithrawDelayTimeChange)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakemanagerWithrawDelayTimeChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerWithrawDelayTimeChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerWithrawDelayTimeChange represents a WithrawDelayTimeChange event raised by the Stakemanager contract.
type StakemanagerWithrawDelayTimeChange struct {
	Cur  *big.Int
	Prev *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterWithrawDelayTimeChange is a free log retrieval operation binding the contract event 0x08cb0bf599c925a6240976039d9d4d3201d88d2ba83703b890049356cdbb67e6.
//
// Solidity: event WithrawDelayTimeChange(uint256 _cur, uint256 _prev)
func (_Stakemanager *StakemanagerFilterer) FilterWithrawDelayTimeChange(opts *bind.FilterOpts) (*StakemanagerWithrawDelayTimeChangeIterator, error) {

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "WithrawDelayTimeChange")
	if err != nil {
		return nil, err
	}
	return &StakemanagerWithrawDelayTimeChangeIterator{contract: _Stakemanager.contract, event: "WithrawDelayTimeChange", logs: logs, sub: sub}, nil
}

// WatchWithrawDelayTimeChange is a free log subscription operation binding the contract event 0x08cb0bf599c925a6240976039d9d4d3201d88d2ba83703b890049356cdbb67e6.
//
// Solidity: event WithrawDelayTimeChange(uint256 _cur, uint256 _prev)
func (_Stakemanager *StakemanagerFilterer) WatchWithrawDelayTimeChange(opts *bind.WatchOpts, sink chan<- *StakemanagerWithrawDelayTimeChange) (event.Subscription, error) {

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "WithrawDelayTimeChange")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerWithrawDelayTimeChange)
				if err := _Stakemanager.contract.UnpackLog(event, "WithrawDelayTimeChange", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithrawDelayTimeChange is a log parse operation binding the contract event 0x08cb0bf599c925a6240976039d9d4d3201d88d2ba83703b890049356cdbb67e6.
//
// Solidity: event WithrawDelayTimeChange(uint256 _cur, uint256 _prev)
func (_Stakemanager *StakemanagerFilterer) ParseWithrawDelayTimeChange(log types.Log) (*StakemanagerWithrawDelayTimeChange, error) {
	event := new(StakemanagerWithrawDelayTimeChange)
	if err := _Stakemanager.contract.UnpackLog(event, "WithrawDelayTimeChange", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
