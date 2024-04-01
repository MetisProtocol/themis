// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package stakinginfo

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

// ISequencerInfoSequencer is an auto generated low-level Go binding around an user-defined struct.
type ISequencerInfoSequencer struct {
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
}

// StakinginfoMetaData contains all meta data concerning the Stakinginfo contract.
var StakinginfoMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_newBatchId\",\"type\":\"uint256\"}],\"name\":\"BatchSubmitReward\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"sequencerId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"totalAmount\",\"type\":\"uint256\"}],\"name\":\"ClaimRewards\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"sequencerId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"newAmount\",\"type\":\"uint256\"}],\"name\":\"LockUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"sequencerId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"activationBatch\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"total\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"signerPubkey\",\"type\":\"bytes\"}],\"name\":\"Locked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"sequencerId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"total\",\"type\":\"uint256\"}],\"name\":\"Relocked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_newMaxLock\",\"type\":\"uint256\"}],\"name\":\"SetMaxLock\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_newMinLock\",\"type\":\"uint256\"}],\"name\":\"SetMinLock\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_payer\",\"type\":\"address\"}],\"name\":\"SetRewardPayer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"sequencerId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldSigner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newSigner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"signerPubkey\",\"type\":\"bytes\"}],\"name\":\"SignerChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"sequencerId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"deactivationBatch\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"deactivationTime\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"unlockClaimTime\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"UnlockInit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"sequencerId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"total\",\"type\":\"uint256\"}],\"name\":\"Unlocked\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"bridge\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_batchId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_totalReward\",\"type\":\"uint256\"}],\"name\":\"distributeReward\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_seqId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_reward\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"_l2gas\",\"type\":\"uint32\"}],\"name\":\"finalizeUnlock\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_seqId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_locked\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_incoming\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fromReward\",\"type\":\"uint256\"}],\"name\":\"increaseLocked\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_manager\",\"type\":\"address\"}],\"name\":\"initManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_bridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_l1Token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_l2Token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_l2ChainId\",\"type\":\"uint256\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_seqId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_reward\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"_l2gas\",\"type\":\"uint32\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"activationBatch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedBatch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deactivationBatch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deactivationTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unlockClaimTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"rewardRecipient\",\"type\":\"address\"},{\"internalType\":\"enumISequencerInfo.Status\",\"name\":\"status\",\"type\":\"uint8\"}],\"internalType\":\"structISequencerInfo.Sequencer\",\"name\":\"_seq\",\"type\":\"tuple\"}],\"name\":\"initializeUnlock\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l1Token\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l2ChainId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l2Token\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_seqId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"_l2gas\",\"type\":\"uint32\"}],\"name\":\"liquidateReward\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"sequencerId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"oldSigner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"newSigner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signerPubkey\",\"type\":\"bytes\"}],\"name\":\"logSignerChange\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"manager\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxLock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minLock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_signer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_batchId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_signerPubkey\",\"type\":\"bytes\"}],\"name\":\"newSequencer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewardPayer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_maxLock\",\"type\":\"uint256\"}],\"name\":\"setMaxLock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_minLock\",\"type\":\"uint256\"}],\"name\":\"setMinLock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_payer\",\"type\":\"address\"}],\"name\":\"setRewardPayer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalLocked\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalRewardsLiquidated\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// StakinginfoABI is the input ABI used to generate the binding from.
// Deprecated: Use StakinginfoMetaData.ABI instead.
var StakinginfoABI = StakinginfoMetaData.ABI

// Stakinginfo is an auto generated Go binding around an Ethereum contract.
type Stakinginfo struct {
	StakinginfoCaller     // Read-only binding to the contract
	StakinginfoTransactor // Write-only binding to the contract
	StakinginfoFilterer   // Log filterer for contract events
}

// StakinginfoCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakinginfoCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakinginfoTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakinginfoTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakinginfoFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakinginfoFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakinginfoSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakinginfoSession struct {
	Contract     *Stakinginfo      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakinginfoCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakinginfoCallerSession struct {
	Contract *StakinginfoCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// StakinginfoTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakinginfoTransactorSession struct {
	Contract     *StakinginfoTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// StakinginfoRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakinginfoRaw struct {
	Contract *Stakinginfo // Generic contract binding to access the raw methods on
}

// StakinginfoCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakinginfoCallerRaw struct {
	Contract *StakinginfoCaller // Generic read-only contract binding to access the raw methods on
}

// StakinginfoTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakinginfoTransactorRaw struct {
	Contract *StakinginfoTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStakinginfo creates a new instance of Stakinginfo, bound to a specific deployed contract.
func NewStakinginfo(address common.Address, backend bind.ContractBackend) (*Stakinginfo, error) {
	contract, err := bindStakinginfo(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Stakinginfo{StakinginfoCaller: StakinginfoCaller{contract: contract}, StakinginfoTransactor: StakinginfoTransactor{contract: contract}, StakinginfoFilterer: StakinginfoFilterer{contract: contract}}, nil
}

// NewStakinginfoCaller creates a new read-only instance of Stakinginfo, bound to a specific deployed contract.
func NewStakinginfoCaller(address common.Address, caller bind.ContractCaller) (*StakinginfoCaller, error) {
	contract, err := bindStakinginfo(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakinginfoCaller{contract: contract}, nil
}

// NewStakinginfoTransactor creates a new write-only instance of Stakinginfo, bound to a specific deployed contract.
func NewStakinginfoTransactor(address common.Address, transactor bind.ContractTransactor) (*StakinginfoTransactor, error) {
	contract, err := bindStakinginfo(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakinginfoTransactor{contract: contract}, nil
}

// NewStakinginfoFilterer creates a new log filterer instance of Stakinginfo, bound to a specific deployed contract.
func NewStakinginfoFilterer(address common.Address, filterer bind.ContractFilterer) (*StakinginfoFilterer, error) {
	contract, err := bindStakinginfo(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakinginfoFilterer{contract: contract}, nil
}

// bindStakinginfo binds a generic wrapper to an already deployed contract.
func bindStakinginfo(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StakinginfoMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Stakinginfo *StakinginfoRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Stakinginfo.Contract.StakinginfoCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Stakinginfo *StakinginfoRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stakinginfo.Contract.StakinginfoTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Stakinginfo *StakinginfoRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Stakinginfo.Contract.StakinginfoTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Stakinginfo *StakinginfoCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Stakinginfo.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Stakinginfo *StakinginfoTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stakinginfo.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Stakinginfo *StakinginfoTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Stakinginfo.Contract.contract.Transact(opts, method, params...)
}

// Bridge is a free data retrieval call binding the contract method 0xe78cea92.
//
// Solidity: function bridge() view returns(address)
func (_Stakinginfo *StakinginfoCaller) Bridge(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Stakinginfo.contract.Call(opts, &out, "bridge")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Bridge is a free data retrieval call binding the contract method 0xe78cea92.
//
// Solidity: function bridge() view returns(address)
func (_Stakinginfo *StakinginfoSession) Bridge() (common.Address, error) {
	return _Stakinginfo.Contract.Bridge(&_Stakinginfo.CallOpts)
}

// Bridge is a free data retrieval call binding the contract method 0xe78cea92.
//
// Solidity: function bridge() view returns(address)
func (_Stakinginfo *StakinginfoCallerSession) Bridge() (common.Address, error) {
	return _Stakinginfo.Contract.Bridge(&_Stakinginfo.CallOpts)
}

// L1Token is a free data retrieval call binding the contract method 0xc01e1bd6.
//
// Solidity: function l1Token() view returns(address)
func (_Stakinginfo *StakinginfoCaller) L1Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Stakinginfo.contract.Call(opts, &out, "l1Token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L1Token is a free data retrieval call binding the contract method 0xc01e1bd6.
//
// Solidity: function l1Token() view returns(address)
func (_Stakinginfo *StakinginfoSession) L1Token() (common.Address, error) {
	return _Stakinginfo.Contract.L1Token(&_Stakinginfo.CallOpts)
}

// L1Token is a free data retrieval call binding the contract method 0xc01e1bd6.
//
// Solidity: function l1Token() view returns(address)
func (_Stakinginfo *StakinginfoCallerSession) L1Token() (common.Address, error) {
	return _Stakinginfo.Contract.L1Token(&_Stakinginfo.CallOpts)
}

// L2ChainId is a free data retrieval call binding the contract method 0xd6ae3cd5.
//
// Solidity: function l2ChainId() view returns(uint256)
func (_Stakinginfo *StakinginfoCaller) L2ChainId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Stakinginfo.contract.Call(opts, &out, "l2ChainId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// L2ChainId is a free data retrieval call binding the contract method 0xd6ae3cd5.
//
// Solidity: function l2ChainId() view returns(uint256)
func (_Stakinginfo *StakinginfoSession) L2ChainId() (*big.Int, error) {
	return _Stakinginfo.Contract.L2ChainId(&_Stakinginfo.CallOpts)
}

// L2ChainId is a free data retrieval call binding the contract method 0xd6ae3cd5.
//
// Solidity: function l2ChainId() view returns(uint256)
func (_Stakinginfo *StakinginfoCallerSession) L2ChainId() (*big.Int, error) {
	return _Stakinginfo.Contract.L2ChainId(&_Stakinginfo.CallOpts)
}

// L2Token is a free data retrieval call binding the contract method 0x56eff267.
//
// Solidity: function l2Token() view returns(address)
func (_Stakinginfo *StakinginfoCaller) L2Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Stakinginfo.contract.Call(opts, &out, "l2Token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L2Token is a free data retrieval call binding the contract method 0x56eff267.
//
// Solidity: function l2Token() view returns(address)
func (_Stakinginfo *StakinginfoSession) L2Token() (common.Address, error) {
	return _Stakinginfo.Contract.L2Token(&_Stakinginfo.CallOpts)
}

// L2Token is a free data retrieval call binding the contract method 0x56eff267.
//
// Solidity: function l2Token() view returns(address)
func (_Stakinginfo *StakinginfoCallerSession) L2Token() (common.Address, error) {
	return _Stakinginfo.Contract.L2Token(&_Stakinginfo.CallOpts)
}

// Manager is a free data retrieval call binding the contract method 0x481c6a75.
//
// Solidity: function manager() view returns(address)
func (_Stakinginfo *StakinginfoCaller) Manager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Stakinginfo.contract.Call(opts, &out, "manager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Manager is a free data retrieval call binding the contract method 0x481c6a75.
//
// Solidity: function manager() view returns(address)
func (_Stakinginfo *StakinginfoSession) Manager() (common.Address, error) {
	return _Stakinginfo.Contract.Manager(&_Stakinginfo.CallOpts)
}

// Manager is a free data retrieval call binding the contract method 0x481c6a75.
//
// Solidity: function manager() view returns(address)
func (_Stakinginfo *StakinginfoCallerSession) Manager() (common.Address, error) {
	return _Stakinginfo.Contract.Manager(&_Stakinginfo.CallOpts)
}

// MaxLock is a free data retrieval call binding the contract method 0x6c0b3e46.
//
// Solidity: function maxLock() view returns(uint256)
func (_Stakinginfo *StakinginfoCaller) MaxLock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Stakinginfo.contract.Call(opts, &out, "maxLock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxLock is a free data retrieval call binding the contract method 0x6c0b3e46.
//
// Solidity: function maxLock() view returns(uint256)
func (_Stakinginfo *StakinginfoSession) MaxLock() (*big.Int, error) {
	return _Stakinginfo.Contract.MaxLock(&_Stakinginfo.CallOpts)
}

// MaxLock is a free data retrieval call binding the contract method 0x6c0b3e46.
//
// Solidity: function maxLock() view returns(uint256)
func (_Stakinginfo *StakinginfoCallerSession) MaxLock() (*big.Int, error) {
	return _Stakinginfo.Contract.MaxLock(&_Stakinginfo.CallOpts)
}

// MinLock is a free data retrieval call binding the contract method 0xf037c630.
//
// Solidity: function minLock() view returns(uint256)
func (_Stakinginfo *StakinginfoCaller) MinLock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Stakinginfo.contract.Call(opts, &out, "minLock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinLock is a free data retrieval call binding the contract method 0xf037c630.
//
// Solidity: function minLock() view returns(uint256)
func (_Stakinginfo *StakinginfoSession) MinLock() (*big.Int, error) {
	return _Stakinginfo.Contract.MinLock(&_Stakinginfo.CallOpts)
}

// MinLock is a free data retrieval call binding the contract method 0xf037c630.
//
// Solidity: function minLock() view returns(uint256)
func (_Stakinginfo *StakinginfoCallerSession) MinLock() (*big.Int, error) {
	return _Stakinginfo.Contract.MinLock(&_Stakinginfo.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Stakinginfo *StakinginfoCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Stakinginfo.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Stakinginfo *StakinginfoSession) Owner() (common.Address, error) {
	return _Stakinginfo.Contract.Owner(&_Stakinginfo.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Stakinginfo *StakinginfoCallerSession) Owner() (common.Address, error) {
	return _Stakinginfo.Contract.Owner(&_Stakinginfo.CallOpts)
}

// RewardPayer is a free data retrieval call binding the contract method 0x6eb27154.
//
// Solidity: function rewardPayer() view returns(address)
func (_Stakinginfo *StakinginfoCaller) RewardPayer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Stakinginfo.contract.Call(opts, &out, "rewardPayer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RewardPayer is a free data retrieval call binding the contract method 0x6eb27154.
//
// Solidity: function rewardPayer() view returns(address)
func (_Stakinginfo *StakinginfoSession) RewardPayer() (common.Address, error) {
	return _Stakinginfo.Contract.RewardPayer(&_Stakinginfo.CallOpts)
}

// RewardPayer is a free data retrieval call binding the contract method 0x6eb27154.
//
// Solidity: function rewardPayer() view returns(address)
func (_Stakinginfo *StakinginfoCallerSession) RewardPayer() (common.Address, error) {
	return _Stakinginfo.Contract.RewardPayer(&_Stakinginfo.CallOpts)
}

// TotalLocked is a free data retrieval call binding the contract method 0x56891412.
//
// Solidity: function totalLocked() view returns(uint256)
func (_Stakinginfo *StakinginfoCaller) TotalLocked(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Stakinginfo.contract.Call(opts, &out, "totalLocked")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalLocked is a free data retrieval call binding the contract method 0x56891412.
//
// Solidity: function totalLocked() view returns(uint256)
func (_Stakinginfo *StakinginfoSession) TotalLocked() (*big.Int, error) {
	return _Stakinginfo.Contract.TotalLocked(&_Stakinginfo.CallOpts)
}

// TotalLocked is a free data retrieval call binding the contract method 0x56891412.
//
// Solidity: function totalLocked() view returns(uint256)
func (_Stakinginfo *StakinginfoCallerSession) TotalLocked() (*big.Int, error) {
	return _Stakinginfo.Contract.TotalLocked(&_Stakinginfo.CallOpts)
}

// TotalRewardsLiquidated is a free data retrieval call binding the contract method 0xcd6b8388.
//
// Solidity: function totalRewardsLiquidated() view returns(uint256)
func (_Stakinginfo *StakinginfoCaller) TotalRewardsLiquidated(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Stakinginfo.contract.Call(opts, &out, "totalRewardsLiquidated")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalRewardsLiquidated is a free data retrieval call binding the contract method 0xcd6b8388.
//
// Solidity: function totalRewardsLiquidated() view returns(uint256)
func (_Stakinginfo *StakinginfoSession) TotalRewardsLiquidated() (*big.Int, error) {
	return _Stakinginfo.Contract.TotalRewardsLiquidated(&_Stakinginfo.CallOpts)
}

// TotalRewardsLiquidated is a free data retrieval call binding the contract method 0xcd6b8388.
//
// Solidity: function totalRewardsLiquidated() view returns(uint256)
func (_Stakinginfo *StakinginfoCallerSession) TotalRewardsLiquidated() (*big.Int, error) {
	return _Stakinginfo.Contract.TotalRewardsLiquidated(&_Stakinginfo.CallOpts)
}

// DistributeReward is a paid mutator transaction binding the contract method 0xe3bcd27c.
//
// Solidity: function distributeReward(uint256 _batchId, uint256 _totalReward) returns()
func (_Stakinginfo *StakinginfoTransactor) DistributeReward(opts *bind.TransactOpts, _batchId *big.Int, _totalReward *big.Int) (*types.Transaction, error) {
	return _Stakinginfo.contract.Transact(opts, "distributeReward", _batchId, _totalReward)
}

// DistributeReward is a paid mutator transaction binding the contract method 0xe3bcd27c.
//
// Solidity: function distributeReward(uint256 _batchId, uint256 _totalReward) returns()
func (_Stakinginfo *StakinginfoSession) DistributeReward(_batchId *big.Int, _totalReward *big.Int) (*types.Transaction, error) {
	return _Stakinginfo.Contract.DistributeReward(&_Stakinginfo.TransactOpts, _batchId, _totalReward)
}

// DistributeReward is a paid mutator transaction binding the contract method 0xe3bcd27c.
//
// Solidity: function distributeReward(uint256 _batchId, uint256 _totalReward) returns()
func (_Stakinginfo *StakinginfoTransactorSession) DistributeReward(_batchId *big.Int, _totalReward *big.Int) (*types.Transaction, error) {
	return _Stakinginfo.Contract.DistributeReward(&_Stakinginfo.TransactOpts, _batchId, _totalReward)
}

// FinalizeUnlock is a paid mutator transaction binding the contract method 0x528ed12a.
//
// Solidity: function finalizeUnlock(address _operator, uint256 _seqId, uint256 _amount, uint256 _reward, address _recipient, uint32 _l2gas) payable returns()
func (_Stakinginfo *StakinginfoTransactor) FinalizeUnlock(opts *bind.TransactOpts, _operator common.Address, _seqId *big.Int, _amount *big.Int, _reward *big.Int, _recipient common.Address, _l2gas uint32) (*types.Transaction, error) {
	return _Stakinginfo.contract.Transact(opts, "finalizeUnlock", _operator, _seqId, _amount, _reward, _recipient, _l2gas)
}

// FinalizeUnlock is a paid mutator transaction binding the contract method 0x528ed12a.
//
// Solidity: function finalizeUnlock(address _operator, uint256 _seqId, uint256 _amount, uint256 _reward, address _recipient, uint32 _l2gas) payable returns()
func (_Stakinginfo *StakinginfoSession) FinalizeUnlock(_operator common.Address, _seqId *big.Int, _amount *big.Int, _reward *big.Int, _recipient common.Address, _l2gas uint32) (*types.Transaction, error) {
	return _Stakinginfo.Contract.FinalizeUnlock(&_Stakinginfo.TransactOpts, _operator, _seqId, _amount, _reward, _recipient, _l2gas)
}

// FinalizeUnlock is a paid mutator transaction binding the contract method 0x528ed12a.
//
// Solidity: function finalizeUnlock(address _operator, uint256 _seqId, uint256 _amount, uint256 _reward, address _recipient, uint32 _l2gas) payable returns()
func (_Stakinginfo *StakinginfoTransactorSession) FinalizeUnlock(_operator common.Address, _seqId *big.Int, _amount *big.Int, _reward *big.Int, _recipient common.Address, _l2gas uint32) (*types.Transaction, error) {
	return _Stakinginfo.Contract.FinalizeUnlock(&_Stakinginfo.TransactOpts, _operator, _seqId, _amount, _reward, _recipient, _l2gas)
}

// IncreaseLocked is a paid mutator transaction binding the contract method 0x2684b8ec.
//
// Solidity: function increaseLocked(uint256 _seqId, uint256 _nonce, address _owner, uint256 _locked, uint256 _incoming, uint256 _fromReward) returns()
func (_Stakinginfo *StakinginfoTransactor) IncreaseLocked(opts *bind.TransactOpts, _seqId *big.Int, _nonce *big.Int, _owner common.Address, _locked *big.Int, _incoming *big.Int, _fromReward *big.Int) (*types.Transaction, error) {
	return _Stakinginfo.contract.Transact(opts, "increaseLocked", _seqId, _nonce, _owner, _locked, _incoming, _fromReward)
}

// IncreaseLocked is a paid mutator transaction binding the contract method 0x2684b8ec.
//
// Solidity: function increaseLocked(uint256 _seqId, uint256 _nonce, address _owner, uint256 _locked, uint256 _incoming, uint256 _fromReward) returns()
func (_Stakinginfo *StakinginfoSession) IncreaseLocked(_seqId *big.Int, _nonce *big.Int, _owner common.Address, _locked *big.Int, _incoming *big.Int, _fromReward *big.Int) (*types.Transaction, error) {
	return _Stakinginfo.Contract.IncreaseLocked(&_Stakinginfo.TransactOpts, _seqId, _nonce, _owner, _locked, _incoming, _fromReward)
}

// IncreaseLocked is a paid mutator transaction binding the contract method 0x2684b8ec.
//
// Solidity: function increaseLocked(uint256 _seqId, uint256 _nonce, address _owner, uint256 _locked, uint256 _incoming, uint256 _fromReward) returns()
func (_Stakinginfo *StakinginfoTransactorSession) IncreaseLocked(_seqId *big.Int, _nonce *big.Int, _owner common.Address, _locked *big.Int, _incoming *big.Int, _fromReward *big.Int) (*types.Transaction, error) {
	return _Stakinginfo.Contract.IncreaseLocked(&_Stakinginfo.TransactOpts, _seqId, _nonce, _owner, _locked, _incoming, _fromReward)
}

// InitManager is a paid mutator transaction binding the contract method 0xb1fc19d3.
//
// Solidity: function initManager(address _manager) returns()
func (_Stakinginfo *StakinginfoTransactor) InitManager(opts *bind.TransactOpts, _manager common.Address) (*types.Transaction, error) {
	return _Stakinginfo.contract.Transact(opts, "initManager", _manager)
}

// InitManager is a paid mutator transaction binding the contract method 0xb1fc19d3.
//
// Solidity: function initManager(address _manager) returns()
func (_Stakinginfo *StakinginfoSession) InitManager(_manager common.Address) (*types.Transaction, error) {
	return _Stakinginfo.Contract.InitManager(&_Stakinginfo.TransactOpts, _manager)
}

// InitManager is a paid mutator transaction binding the contract method 0xb1fc19d3.
//
// Solidity: function initManager(address _manager) returns()
func (_Stakinginfo *StakinginfoTransactorSession) InitManager(_manager common.Address) (*types.Transaction, error) {
	return _Stakinginfo.Contract.InitManager(&_Stakinginfo.TransactOpts, _manager)
}

// Initialize is a paid mutator transaction binding the contract method 0xcf756fdf.
//
// Solidity: function initialize(address _bridge, address _l1Token, address _l2Token, uint256 _l2ChainId) returns()
func (_Stakinginfo *StakinginfoTransactor) Initialize(opts *bind.TransactOpts, _bridge common.Address, _l1Token common.Address, _l2Token common.Address, _l2ChainId *big.Int) (*types.Transaction, error) {
	return _Stakinginfo.contract.Transact(opts, "initialize", _bridge, _l1Token, _l2Token, _l2ChainId)
}

// Initialize is a paid mutator transaction binding the contract method 0xcf756fdf.
//
// Solidity: function initialize(address _bridge, address _l1Token, address _l2Token, uint256 _l2ChainId) returns()
func (_Stakinginfo *StakinginfoSession) Initialize(_bridge common.Address, _l1Token common.Address, _l2Token common.Address, _l2ChainId *big.Int) (*types.Transaction, error) {
	return _Stakinginfo.Contract.Initialize(&_Stakinginfo.TransactOpts, _bridge, _l1Token, _l2Token, _l2ChainId)
}

// Initialize is a paid mutator transaction binding the contract method 0xcf756fdf.
//
// Solidity: function initialize(address _bridge, address _l1Token, address _l2Token, uint256 _l2ChainId) returns()
func (_Stakinginfo *StakinginfoTransactorSession) Initialize(_bridge common.Address, _l1Token common.Address, _l2Token common.Address, _l2ChainId *big.Int) (*types.Transaction, error) {
	return _Stakinginfo.Contract.Initialize(&_Stakinginfo.TransactOpts, _bridge, _l1Token, _l2Token, _l2ChainId)
}

// InitializeUnlock is a paid mutator transaction binding the contract method 0x2243069c.
//
// Solidity: function initializeUnlock(uint256 _seqId, uint256 _reward, uint32 _l2gas, (uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,address,bytes,address,uint8) _seq) payable returns()
func (_Stakinginfo *StakinginfoTransactor) InitializeUnlock(opts *bind.TransactOpts, _seqId *big.Int, _reward *big.Int, _l2gas uint32, _seq ISequencerInfoSequencer) (*types.Transaction, error) {
	return _Stakinginfo.contract.Transact(opts, "initializeUnlock", _seqId, _reward, _l2gas, _seq)
}

// InitializeUnlock is a paid mutator transaction binding the contract method 0x2243069c.
//
// Solidity: function initializeUnlock(uint256 _seqId, uint256 _reward, uint32 _l2gas, (uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,address,bytes,address,uint8) _seq) payable returns()
func (_Stakinginfo *StakinginfoSession) InitializeUnlock(_seqId *big.Int, _reward *big.Int, _l2gas uint32, _seq ISequencerInfoSequencer) (*types.Transaction, error) {
	return _Stakinginfo.Contract.InitializeUnlock(&_Stakinginfo.TransactOpts, _seqId, _reward, _l2gas, _seq)
}

// InitializeUnlock is a paid mutator transaction binding the contract method 0x2243069c.
//
// Solidity: function initializeUnlock(uint256 _seqId, uint256 _reward, uint32 _l2gas, (uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,address,bytes,address,uint8) _seq) payable returns()
func (_Stakinginfo *StakinginfoTransactorSession) InitializeUnlock(_seqId *big.Int, _reward *big.Int, _l2gas uint32, _seq ISequencerInfoSequencer) (*types.Transaction, error) {
	return _Stakinginfo.Contract.InitializeUnlock(&_Stakinginfo.TransactOpts, _seqId, _reward, _l2gas, _seq)
}

// LiquidateReward is a paid mutator transaction binding the contract method 0x5d7878a8.
//
// Solidity: function liquidateReward(uint256 _seqId, uint256 _amount, address _recipient, uint32 _l2gas) payable returns()
func (_Stakinginfo *StakinginfoTransactor) LiquidateReward(opts *bind.TransactOpts, _seqId *big.Int, _amount *big.Int, _recipient common.Address, _l2gas uint32) (*types.Transaction, error) {
	return _Stakinginfo.contract.Transact(opts, "liquidateReward", _seqId, _amount, _recipient, _l2gas)
}

// LiquidateReward is a paid mutator transaction binding the contract method 0x5d7878a8.
//
// Solidity: function liquidateReward(uint256 _seqId, uint256 _amount, address _recipient, uint32 _l2gas) payable returns()
func (_Stakinginfo *StakinginfoSession) LiquidateReward(_seqId *big.Int, _amount *big.Int, _recipient common.Address, _l2gas uint32) (*types.Transaction, error) {
	return _Stakinginfo.Contract.LiquidateReward(&_Stakinginfo.TransactOpts, _seqId, _amount, _recipient, _l2gas)
}

// LiquidateReward is a paid mutator transaction binding the contract method 0x5d7878a8.
//
// Solidity: function liquidateReward(uint256 _seqId, uint256 _amount, address _recipient, uint32 _l2gas) payable returns()
func (_Stakinginfo *StakinginfoTransactorSession) LiquidateReward(_seqId *big.Int, _amount *big.Int, _recipient common.Address, _l2gas uint32) (*types.Transaction, error) {
	return _Stakinginfo.Contract.LiquidateReward(&_Stakinginfo.TransactOpts, _seqId, _amount, _recipient, _l2gas)
}

// LogSignerChange is a paid mutator transaction binding the contract method 0xb3285702.
//
// Solidity: function logSignerChange(uint256 sequencerId, address oldSigner, address newSigner, uint256 nonce, bytes signerPubkey) returns()
func (_Stakinginfo *StakinginfoTransactor) LogSignerChange(opts *bind.TransactOpts, sequencerId *big.Int, oldSigner common.Address, newSigner common.Address, nonce *big.Int, signerPubkey []byte) (*types.Transaction, error) {
	return _Stakinginfo.contract.Transact(opts, "logSignerChange", sequencerId, oldSigner, newSigner, nonce, signerPubkey)
}

// LogSignerChange is a paid mutator transaction binding the contract method 0xb3285702.
//
// Solidity: function logSignerChange(uint256 sequencerId, address oldSigner, address newSigner, uint256 nonce, bytes signerPubkey) returns()
func (_Stakinginfo *StakinginfoSession) LogSignerChange(sequencerId *big.Int, oldSigner common.Address, newSigner common.Address, nonce *big.Int, signerPubkey []byte) (*types.Transaction, error) {
	return _Stakinginfo.Contract.LogSignerChange(&_Stakinginfo.TransactOpts, sequencerId, oldSigner, newSigner, nonce, signerPubkey)
}

// LogSignerChange is a paid mutator transaction binding the contract method 0xb3285702.
//
// Solidity: function logSignerChange(uint256 sequencerId, address oldSigner, address newSigner, uint256 nonce, bytes signerPubkey) returns()
func (_Stakinginfo *StakinginfoTransactorSession) LogSignerChange(sequencerId *big.Int, oldSigner common.Address, newSigner common.Address, nonce *big.Int, signerPubkey []byte) (*types.Transaction, error) {
	return _Stakinginfo.Contract.LogSignerChange(&_Stakinginfo.TransactOpts, sequencerId, oldSigner, newSigner, nonce, signerPubkey)
}

// NewSequencer is a paid mutator transaction binding the contract method 0x1badded5.
//
// Solidity: function newSequencer(uint256 _id, address _owner, address _signer, uint256 _amount, uint256 _batchId, bytes _signerPubkey) returns()
func (_Stakinginfo *StakinginfoTransactor) NewSequencer(opts *bind.TransactOpts, _id *big.Int, _owner common.Address, _signer common.Address, _amount *big.Int, _batchId *big.Int, _signerPubkey []byte) (*types.Transaction, error) {
	return _Stakinginfo.contract.Transact(opts, "newSequencer", _id, _owner, _signer, _amount, _batchId, _signerPubkey)
}

// NewSequencer is a paid mutator transaction binding the contract method 0x1badded5.
//
// Solidity: function newSequencer(uint256 _id, address _owner, address _signer, uint256 _amount, uint256 _batchId, bytes _signerPubkey) returns()
func (_Stakinginfo *StakinginfoSession) NewSequencer(_id *big.Int, _owner common.Address, _signer common.Address, _amount *big.Int, _batchId *big.Int, _signerPubkey []byte) (*types.Transaction, error) {
	return _Stakinginfo.Contract.NewSequencer(&_Stakinginfo.TransactOpts, _id, _owner, _signer, _amount, _batchId, _signerPubkey)
}

// NewSequencer is a paid mutator transaction binding the contract method 0x1badded5.
//
// Solidity: function newSequencer(uint256 _id, address _owner, address _signer, uint256 _amount, uint256 _batchId, bytes _signerPubkey) returns()
func (_Stakinginfo *StakinginfoTransactorSession) NewSequencer(_id *big.Int, _owner common.Address, _signer common.Address, _amount *big.Int, _batchId *big.Int, _signerPubkey []byte) (*types.Transaction, error) {
	return _Stakinginfo.Contract.NewSequencer(&_Stakinginfo.TransactOpts, _id, _owner, _signer, _amount, _batchId, _signerPubkey)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Stakinginfo *StakinginfoTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stakinginfo.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Stakinginfo *StakinginfoSession) RenounceOwnership() (*types.Transaction, error) {
	return _Stakinginfo.Contract.RenounceOwnership(&_Stakinginfo.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Stakinginfo *StakinginfoTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Stakinginfo.Contract.RenounceOwnership(&_Stakinginfo.TransactOpts)
}

// SetMaxLock is a paid mutator transaction binding the contract method 0xcd15b2a5.
//
// Solidity: function setMaxLock(uint256 _maxLock) returns()
func (_Stakinginfo *StakinginfoTransactor) SetMaxLock(opts *bind.TransactOpts, _maxLock *big.Int) (*types.Transaction, error) {
	return _Stakinginfo.contract.Transact(opts, "setMaxLock", _maxLock)
}

// SetMaxLock is a paid mutator transaction binding the contract method 0xcd15b2a5.
//
// Solidity: function setMaxLock(uint256 _maxLock) returns()
func (_Stakinginfo *StakinginfoSession) SetMaxLock(_maxLock *big.Int) (*types.Transaction, error) {
	return _Stakinginfo.Contract.SetMaxLock(&_Stakinginfo.TransactOpts, _maxLock)
}

// SetMaxLock is a paid mutator transaction binding the contract method 0xcd15b2a5.
//
// Solidity: function setMaxLock(uint256 _maxLock) returns()
func (_Stakinginfo *StakinginfoTransactorSession) SetMaxLock(_maxLock *big.Int) (*types.Transaction, error) {
	return _Stakinginfo.Contract.SetMaxLock(&_Stakinginfo.TransactOpts, _maxLock)
}

// SetMinLock is a paid mutator transaction binding the contract method 0xaa15af6a.
//
// Solidity: function setMinLock(uint256 _minLock) returns()
func (_Stakinginfo *StakinginfoTransactor) SetMinLock(opts *bind.TransactOpts, _minLock *big.Int) (*types.Transaction, error) {
	return _Stakinginfo.contract.Transact(opts, "setMinLock", _minLock)
}

// SetMinLock is a paid mutator transaction binding the contract method 0xaa15af6a.
//
// Solidity: function setMinLock(uint256 _minLock) returns()
func (_Stakinginfo *StakinginfoSession) SetMinLock(_minLock *big.Int) (*types.Transaction, error) {
	return _Stakinginfo.Contract.SetMinLock(&_Stakinginfo.TransactOpts, _minLock)
}

// SetMinLock is a paid mutator transaction binding the contract method 0xaa15af6a.
//
// Solidity: function setMinLock(uint256 _minLock) returns()
func (_Stakinginfo *StakinginfoTransactorSession) SetMinLock(_minLock *big.Int) (*types.Transaction, error) {
	return _Stakinginfo.Contract.SetMinLock(&_Stakinginfo.TransactOpts, _minLock)
}

// SetRewardPayer is a paid mutator transaction binding the contract method 0xe8b8b413.
//
// Solidity: function setRewardPayer(address _payer) returns()
func (_Stakinginfo *StakinginfoTransactor) SetRewardPayer(opts *bind.TransactOpts, _payer common.Address) (*types.Transaction, error) {
	return _Stakinginfo.contract.Transact(opts, "setRewardPayer", _payer)
}

// SetRewardPayer is a paid mutator transaction binding the contract method 0xe8b8b413.
//
// Solidity: function setRewardPayer(address _payer) returns()
func (_Stakinginfo *StakinginfoSession) SetRewardPayer(_payer common.Address) (*types.Transaction, error) {
	return _Stakinginfo.Contract.SetRewardPayer(&_Stakinginfo.TransactOpts, _payer)
}

// SetRewardPayer is a paid mutator transaction binding the contract method 0xe8b8b413.
//
// Solidity: function setRewardPayer(address _payer) returns()
func (_Stakinginfo *StakinginfoTransactorSession) SetRewardPayer(_payer common.Address) (*types.Transaction, error) {
	return _Stakinginfo.Contract.SetRewardPayer(&_Stakinginfo.TransactOpts, _payer)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Stakinginfo *StakinginfoTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Stakinginfo.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Stakinginfo *StakinginfoSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Stakinginfo.Contract.TransferOwnership(&_Stakinginfo.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Stakinginfo *StakinginfoTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Stakinginfo.Contract.TransferOwnership(&_Stakinginfo.TransactOpts, newOwner)
}

// StakinginfoBatchSubmitRewardIterator is returned from FilterBatchSubmitReward and is used to iterate over the raw logs and unpacked data for BatchSubmitReward events raised by the Stakinginfo contract.
type StakinginfoBatchSubmitRewardIterator struct {
	Event *StakinginfoBatchSubmitReward // Event containing the contract specifics and raw log

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
func (it *StakinginfoBatchSubmitRewardIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoBatchSubmitReward)
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
		it.Event = new(StakinginfoBatchSubmitReward)
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
func (it *StakinginfoBatchSubmitRewardIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoBatchSubmitRewardIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoBatchSubmitReward represents a BatchSubmitReward event raised by the Stakinginfo contract.
type StakinginfoBatchSubmitReward struct {
	NewBatchId *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBatchSubmitReward is a free log retrieval operation binding the contract event 0x9e5aedd489785d05ba086064386f2e75b3e497d3dc00a54ed1c78bfc50a3953f.
//
// Solidity: event BatchSubmitReward(uint256 _newBatchId)
func (_Stakinginfo *StakinginfoFilterer) FilterBatchSubmitReward(opts *bind.FilterOpts) (*StakinginfoBatchSubmitRewardIterator, error) {

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "BatchSubmitReward")
	if err != nil {
		return nil, err
	}
	return &StakinginfoBatchSubmitRewardIterator{contract: _Stakinginfo.contract, event: "BatchSubmitReward", logs: logs, sub: sub}, nil
}

// WatchBatchSubmitReward is a free log subscription operation binding the contract event 0x9e5aedd489785d05ba086064386f2e75b3e497d3dc00a54ed1c78bfc50a3953f.
//
// Solidity: event BatchSubmitReward(uint256 _newBatchId)
func (_Stakinginfo *StakinginfoFilterer) WatchBatchSubmitReward(opts *bind.WatchOpts, sink chan<- *StakinginfoBatchSubmitReward) (event.Subscription, error) {

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "BatchSubmitReward")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoBatchSubmitReward)
				if err := _Stakinginfo.contract.UnpackLog(event, "BatchSubmitReward", log); err != nil {
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

// ParseBatchSubmitReward is a log parse operation binding the contract event 0x9e5aedd489785d05ba086064386f2e75b3e497d3dc00a54ed1c78bfc50a3953f.
//
// Solidity: event BatchSubmitReward(uint256 _newBatchId)
func (_Stakinginfo *StakinginfoFilterer) ParseBatchSubmitReward(log types.Log) (*StakinginfoBatchSubmitReward, error) {
	event := new(StakinginfoBatchSubmitReward)
	if err := _Stakinginfo.contract.UnpackLog(event, "BatchSubmitReward", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakinginfoClaimRewardsIterator is returned from FilterClaimRewards and is used to iterate over the raw logs and unpacked data for ClaimRewards events raised by the Stakinginfo contract.
type StakinginfoClaimRewardsIterator struct {
	Event *StakinginfoClaimRewards // Event containing the contract specifics and raw log

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
func (it *StakinginfoClaimRewardsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoClaimRewards)
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
		it.Event = new(StakinginfoClaimRewards)
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
func (it *StakinginfoClaimRewardsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoClaimRewardsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoClaimRewards represents a ClaimRewards event raised by the Stakinginfo contract.
type StakinginfoClaimRewards struct {
	SequencerId *big.Int
	Recipient   common.Address
	Amount      *big.Int
	TotalAmount *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterClaimRewards is a free log retrieval operation binding the contract event 0x18c7dc2a1800c409227dc12c0c05ada9c995ebfe0e42ae6d65f1b3ae3e6111de.
//
// Solidity: event ClaimRewards(uint256 indexed sequencerId, address recipient, uint256 indexed amount, uint256 indexed totalAmount)
func (_Stakinginfo *StakinginfoFilterer) FilterClaimRewards(opts *bind.FilterOpts, sequencerId []*big.Int, amount []*big.Int, totalAmount []*big.Int) (*StakinginfoClaimRewardsIterator, error) {

	var sequencerIdRule []interface{}
	for _, sequencerIdItem := range sequencerId {
		sequencerIdRule = append(sequencerIdRule, sequencerIdItem)
	}

	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var totalAmountRule []interface{}
	for _, totalAmountItem := range totalAmount {
		totalAmountRule = append(totalAmountRule, totalAmountItem)
	}

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "ClaimRewards", sequencerIdRule, amountRule, totalAmountRule)
	if err != nil {
		return nil, err
	}
	return &StakinginfoClaimRewardsIterator{contract: _Stakinginfo.contract, event: "ClaimRewards", logs: logs, sub: sub}, nil
}

// WatchClaimRewards is a free log subscription operation binding the contract event 0x18c7dc2a1800c409227dc12c0c05ada9c995ebfe0e42ae6d65f1b3ae3e6111de.
//
// Solidity: event ClaimRewards(uint256 indexed sequencerId, address recipient, uint256 indexed amount, uint256 indexed totalAmount)
func (_Stakinginfo *StakinginfoFilterer) WatchClaimRewards(opts *bind.WatchOpts, sink chan<- *StakinginfoClaimRewards, sequencerId []*big.Int, amount []*big.Int, totalAmount []*big.Int) (event.Subscription, error) {

	var sequencerIdRule []interface{}
	for _, sequencerIdItem := range sequencerId {
		sequencerIdRule = append(sequencerIdRule, sequencerIdItem)
	}

	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var totalAmountRule []interface{}
	for _, totalAmountItem := range totalAmount {
		totalAmountRule = append(totalAmountRule, totalAmountItem)
	}

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "ClaimRewards", sequencerIdRule, amountRule, totalAmountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoClaimRewards)
				if err := _Stakinginfo.contract.UnpackLog(event, "ClaimRewards", log); err != nil {
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

// ParseClaimRewards is a log parse operation binding the contract event 0x18c7dc2a1800c409227dc12c0c05ada9c995ebfe0e42ae6d65f1b3ae3e6111de.
//
// Solidity: event ClaimRewards(uint256 indexed sequencerId, address recipient, uint256 indexed amount, uint256 indexed totalAmount)
func (_Stakinginfo *StakinginfoFilterer) ParseClaimRewards(log types.Log) (*StakinginfoClaimRewards, error) {
	event := new(StakinginfoClaimRewards)
	if err := _Stakinginfo.contract.UnpackLog(event, "ClaimRewards", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakinginfoInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Stakinginfo contract.
type StakinginfoInitializedIterator struct {
	Event *StakinginfoInitialized // Event containing the contract specifics and raw log

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
func (it *StakinginfoInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoInitialized)
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
		it.Event = new(StakinginfoInitialized)
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
func (it *StakinginfoInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoInitialized represents a Initialized event raised by the Stakinginfo contract.
type StakinginfoInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Stakinginfo *StakinginfoFilterer) FilterInitialized(opts *bind.FilterOpts) (*StakinginfoInitializedIterator, error) {

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &StakinginfoInitializedIterator{contract: _Stakinginfo.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Stakinginfo *StakinginfoFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *StakinginfoInitialized) (event.Subscription, error) {

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoInitialized)
				if err := _Stakinginfo.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Stakinginfo *StakinginfoFilterer) ParseInitialized(log types.Log) (*StakinginfoInitialized, error) {
	event := new(StakinginfoInitialized)
	if err := _Stakinginfo.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakinginfoLockUpdateIterator is returned from FilterLockUpdate and is used to iterate over the raw logs and unpacked data for LockUpdate events raised by the Stakinginfo contract.
type StakinginfoLockUpdateIterator struct {
	Event *StakinginfoLockUpdate // Event containing the contract specifics and raw log

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
func (it *StakinginfoLockUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoLockUpdate)
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
		it.Event = new(StakinginfoLockUpdate)
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
func (it *StakinginfoLockUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoLockUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoLockUpdate represents a LockUpdate event raised by the Stakinginfo contract.
type StakinginfoLockUpdate struct {
	SequencerId *big.Int
	Nonce       *big.Int
	NewAmount   *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterLockUpdate is a free log retrieval operation binding the contract event 0xd716c027b3dd610e4534df756848128bbb159a757724c17d89fcc4d0151b1f30.
//
// Solidity: event LockUpdate(uint256 indexed sequencerId, uint256 indexed nonce, uint256 indexed newAmount)
func (_Stakinginfo *StakinginfoFilterer) FilterLockUpdate(opts *bind.FilterOpts, sequencerId []*big.Int, nonce []*big.Int, newAmount []*big.Int) (*StakinginfoLockUpdateIterator, error) {

	var sequencerIdRule []interface{}
	for _, sequencerIdItem := range sequencerId {
		sequencerIdRule = append(sequencerIdRule, sequencerIdItem)
	}
	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}
	var newAmountRule []interface{}
	for _, newAmountItem := range newAmount {
		newAmountRule = append(newAmountRule, newAmountItem)
	}

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "LockUpdate", sequencerIdRule, nonceRule, newAmountRule)
	if err != nil {
		return nil, err
	}
	return &StakinginfoLockUpdateIterator{contract: _Stakinginfo.contract, event: "LockUpdate", logs: logs, sub: sub}, nil
}

// WatchLockUpdate is a free log subscription operation binding the contract event 0xd716c027b3dd610e4534df756848128bbb159a757724c17d89fcc4d0151b1f30.
//
// Solidity: event LockUpdate(uint256 indexed sequencerId, uint256 indexed nonce, uint256 indexed newAmount)
func (_Stakinginfo *StakinginfoFilterer) WatchLockUpdate(opts *bind.WatchOpts, sink chan<- *StakinginfoLockUpdate, sequencerId []*big.Int, nonce []*big.Int, newAmount []*big.Int) (event.Subscription, error) {

	var sequencerIdRule []interface{}
	for _, sequencerIdItem := range sequencerId {
		sequencerIdRule = append(sequencerIdRule, sequencerIdItem)
	}
	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}
	var newAmountRule []interface{}
	for _, newAmountItem := range newAmount {
		newAmountRule = append(newAmountRule, newAmountItem)
	}

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "LockUpdate", sequencerIdRule, nonceRule, newAmountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoLockUpdate)
				if err := _Stakinginfo.contract.UnpackLog(event, "LockUpdate", log); err != nil {
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

// ParseLockUpdate is a log parse operation binding the contract event 0xd716c027b3dd610e4534df756848128bbb159a757724c17d89fcc4d0151b1f30.
//
// Solidity: event LockUpdate(uint256 indexed sequencerId, uint256 indexed nonce, uint256 indexed newAmount)
func (_Stakinginfo *StakinginfoFilterer) ParseLockUpdate(log types.Log) (*StakinginfoLockUpdate, error) {
	event := new(StakinginfoLockUpdate)
	if err := _Stakinginfo.contract.UnpackLog(event, "LockUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakinginfoLockedIterator is returned from FilterLocked and is used to iterate over the raw logs and unpacked data for Locked events raised by the Stakinginfo contract.
type StakinginfoLockedIterator struct {
	Event *StakinginfoLocked // Event containing the contract specifics and raw log

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
func (it *StakinginfoLockedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoLocked)
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
		it.Event = new(StakinginfoLocked)
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
func (it *StakinginfoLockedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoLockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoLocked represents a Locked event raised by the Stakinginfo contract.
type StakinginfoLocked struct {
	Signer          common.Address
	SequencerId     *big.Int
	Nonce           *big.Int
	ActivationBatch *big.Int
	Amount          *big.Int
	Total           *big.Int
	SignerPubkey    []byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterLocked is a free log retrieval operation binding the contract event 0xe6f1eb1f1d0ca344d03cf47b9e6ece8a7d3b196e38dd7dd2307cca75e26682a8.
//
// Solidity: event Locked(address indexed signer, uint256 indexed sequencerId, uint256 nonce, uint256 indexed activationBatch, uint256 amount, uint256 total, bytes signerPubkey)
func (_Stakinginfo *StakinginfoFilterer) FilterLocked(opts *bind.FilterOpts, signer []common.Address, sequencerId []*big.Int, activationBatch []*big.Int) (*StakinginfoLockedIterator, error) {

	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}
	var sequencerIdRule []interface{}
	for _, sequencerIdItem := range sequencerId {
		sequencerIdRule = append(sequencerIdRule, sequencerIdItem)
	}

	var activationBatchRule []interface{}
	for _, activationBatchItem := range activationBatch {
		activationBatchRule = append(activationBatchRule, activationBatchItem)
	}

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "Locked", signerRule, sequencerIdRule, activationBatchRule)
	if err != nil {
		return nil, err
	}
	return &StakinginfoLockedIterator{contract: _Stakinginfo.contract, event: "Locked", logs: logs, sub: sub}, nil
}

// WatchLocked is a free log subscription operation binding the contract event 0xe6f1eb1f1d0ca344d03cf47b9e6ece8a7d3b196e38dd7dd2307cca75e26682a8.
//
// Solidity: event Locked(address indexed signer, uint256 indexed sequencerId, uint256 nonce, uint256 indexed activationBatch, uint256 amount, uint256 total, bytes signerPubkey)
func (_Stakinginfo *StakinginfoFilterer) WatchLocked(opts *bind.WatchOpts, sink chan<- *StakinginfoLocked, signer []common.Address, sequencerId []*big.Int, activationBatch []*big.Int) (event.Subscription, error) {

	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}
	var sequencerIdRule []interface{}
	for _, sequencerIdItem := range sequencerId {
		sequencerIdRule = append(sequencerIdRule, sequencerIdItem)
	}

	var activationBatchRule []interface{}
	for _, activationBatchItem := range activationBatch {
		activationBatchRule = append(activationBatchRule, activationBatchItem)
	}

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "Locked", signerRule, sequencerIdRule, activationBatchRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoLocked)
				if err := _Stakinginfo.contract.UnpackLog(event, "Locked", log); err != nil {
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

// ParseLocked is a log parse operation binding the contract event 0xe6f1eb1f1d0ca344d03cf47b9e6ece8a7d3b196e38dd7dd2307cca75e26682a8.
//
// Solidity: event Locked(address indexed signer, uint256 indexed sequencerId, uint256 nonce, uint256 indexed activationBatch, uint256 amount, uint256 total, bytes signerPubkey)
func (_Stakinginfo *StakinginfoFilterer) ParseLocked(log types.Log) (*StakinginfoLocked, error) {
	event := new(StakinginfoLocked)
	if err := _Stakinginfo.contract.UnpackLog(event, "Locked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakinginfoOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Stakinginfo contract.
type StakinginfoOwnershipTransferredIterator struct {
	Event *StakinginfoOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *StakinginfoOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoOwnershipTransferred)
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
		it.Event = new(StakinginfoOwnershipTransferred)
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
func (it *StakinginfoOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoOwnershipTransferred represents a OwnershipTransferred event raised by the Stakinginfo contract.
type StakinginfoOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Stakinginfo *StakinginfoFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*StakinginfoOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &StakinginfoOwnershipTransferredIterator{contract: _Stakinginfo.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Stakinginfo *StakinginfoFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *StakinginfoOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoOwnershipTransferred)
				if err := _Stakinginfo.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Stakinginfo *StakinginfoFilterer) ParseOwnershipTransferred(log types.Log) (*StakinginfoOwnershipTransferred, error) {
	event := new(StakinginfoOwnershipTransferred)
	if err := _Stakinginfo.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakinginfoRelockedIterator is returned from FilterRelocked and is used to iterate over the raw logs and unpacked data for Relocked events raised by the Stakinginfo contract.
type StakinginfoRelockedIterator struct {
	Event *StakinginfoRelocked // Event containing the contract specifics and raw log

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
func (it *StakinginfoRelockedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoRelocked)
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
		it.Event = new(StakinginfoRelocked)
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
func (it *StakinginfoRelockedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoRelockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoRelocked represents a Relocked event raised by the Stakinginfo contract.
type StakinginfoRelocked struct {
	SequencerId *big.Int
	Amount      *big.Int
	Total       *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRelocked is a free log retrieval operation binding the contract event 0x33a87ba488658b3d1319098cd49c6d65b72a79c0f3530fec611e7afffed04395.
//
// Solidity: event Relocked(uint256 indexed sequencerId, uint256 amount, uint256 total)
func (_Stakinginfo *StakinginfoFilterer) FilterRelocked(opts *bind.FilterOpts, sequencerId []*big.Int) (*StakinginfoRelockedIterator, error) {

	var sequencerIdRule []interface{}
	for _, sequencerIdItem := range sequencerId {
		sequencerIdRule = append(sequencerIdRule, sequencerIdItem)
	}

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "Relocked", sequencerIdRule)
	if err != nil {
		return nil, err
	}
	return &StakinginfoRelockedIterator{contract: _Stakinginfo.contract, event: "Relocked", logs: logs, sub: sub}, nil
}

// WatchRelocked is a free log subscription operation binding the contract event 0x33a87ba488658b3d1319098cd49c6d65b72a79c0f3530fec611e7afffed04395.
//
// Solidity: event Relocked(uint256 indexed sequencerId, uint256 amount, uint256 total)
func (_Stakinginfo *StakinginfoFilterer) WatchRelocked(opts *bind.WatchOpts, sink chan<- *StakinginfoRelocked, sequencerId []*big.Int) (event.Subscription, error) {

	var sequencerIdRule []interface{}
	for _, sequencerIdItem := range sequencerId {
		sequencerIdRule = append(sequencerIdRule, sequencerIdItem)
	}

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "Relocked", sequencerIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoRelocked)
				if err := _Stakinginfo.contract.UnpackLog(event, "Relocked", log); err != nil {
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

// ParseRelocked is a log parse operation binding the contract event 0x33a87ba488658b3d1319098cd49c6d65b72a79c0f3530fec611e7afffed04395.
//
// Solidity: event Relocked(uint256 indexed sequencerId, uint256 amount, uint256 total)
func (_Stakinginfo *StakinginfoFilterer) ParseRelocked(log types.Log) (*StakinginfoRelocked, error) {
	event := new(StakinginfoRelocked)
	if err := _Stakinginfo.contract.UnpackLog(event, "Relocked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakinginfoSetMaxLockIterator is returned from FilterSetMaxLock and is used to iterate over the raw logs and unpacked data for SetMaxLock events raised by the Stakinginfo contract.
type StakinginfoSetMaxLockIterator struct {
	Event *StakinginfoSetMaxLock // Event containing the contract specifics and raw log

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
func (it *StakinginfoSetMaxLockIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoSetMaxLock)
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
		it.Event = new(StakinginfoSetMaxLock)
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
func (it *StakinginfoSetMaxLockIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoSetMaxLockIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoSetMaxLock represents a SetMaxLock event raised by the Stakinginfo contract.
type StakinginfoSetMaxLock struct {
	NewMaxLock *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSetMaxLock is a free log retrieval operation binding the contract event 0xbe23e9641c545443c3c625039b327c0eee88e9024040be7b03c5d73862d425e0.
//
// Solidity: event SetMaxLock(uint256 _newMaxLock)
func (_Stakinginfo *StakinginfoFilterer) FilterSetMaxLock(opts *bind.FilterOpts) (*StakinginfoSetMaxLockIterator, error) {

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "SetMaxLock")
	if err != nil {
		return nil, err
	}
	return &StakinginfoSetMaxLockIterator{contract: _Stakinginfo.contract, event: "SetMaxLock", logs: logs, sub: sub}, nil
}

// WatchSetMaxLock is a free log subscription operation binding the contract event 0xbe23e9641c545443c3c625039b327c0eee88e9024040be7b03c5d73862d425e0.
//
// Solidity: event SetMaxLock(uint256 _newMaxLock)
func (_Stakinginfo *StakinginfoFilterer) WatchSetMaxLock(opts *bind.WatchOpts, sink chan<- *StakinginfoSetMaxLock) (event.Subscription, error) {

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "SetMaxLock")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoSetMaxLock)
				if err := _Stakinginfo.contract.UnpackLog(event, "SetMaxLock", log); err != nil {
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

// ParseSetMaxLock is a log parse operation binding the contract event 0xbe23e9641c545443c3c625039b327c0eee88e9024040be7b03c5d73862d425e0.
//
// Solidity: event SetMaxLock(uint256 _newMaxLock)
func (_Stakinginfo *StakinginfoFilterer) ParseSetMaxLock(log types.Log) (*StakinginfoSetMaxLock, error) {
	event := new(StakinginfoSetMaxLock)
	if err := _Stakinginfo.contract.UnpackLog(event, "SetMaxLock", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakinginfoSetMinLockIterator is returned from FilterSetMinLock and is used to iterate over the raw logs and unpacked data for SetMinLock events raised by the Stakinginfo contract.
type StakinginfoSetMinLockIterator struct {
	Event *StakinginfoSetMinLock // Event containing the contract specifics and raw log

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
func (it *StakinginfoSetMinLockIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoSetMinLock)
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
		it.Event = new(StakinginfoSetMinLock)
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
func (it *StakinginfoSetMinLockIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoSetMinLockIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoSetMinLock represents a SetMinLock event raised by the Stakinginfo contract.
type StakinginfoSetMinLock struct {
	NewMinLock *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSetMinLock is a free log retrieval operation binding the contract event 0xabb05374bb45ebfef33afb21ec5aa52333708d8217fd8e5c0616efd2530b2145.
//
// Solidity: event SetMinLock(uint256 _newMinLock)
func (_Stakinginfo *StakinginfoFilterer) FilterSetMinLock(opts *bind.FilterOpts) (*StakinginfoSetMinLockIterator, error) {

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "SetMinLock")
	if err != nil {
		return nil, err
	}
	return &StakinginfoSetMinLockIterator{contract: _Stakinginfo.contract, event: "SetMinLock", logs: logs, sub: sub}, nil
}

// WatchSetMinLock is a free log subscription operation binding the contract event 0xabb05374bb45ebfef33afb21ec5aa52333708d8217fd8e5c0616efd2530b2145.
//
// Solidity: event SetMinLock(uint256 _newMinLock)
func (_Stakinginfo *StakinginfoFilterer) WatchSetMinLock(opts *bind.WatchOpts, sink chan<- *StakinginfoSetMinLock) (event.Subscription, error) {

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "SetMinLock")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoSetMinLock)
				if err := _Stakinginfo.contract.UnpackLog(event, "SetMinLock", log); err != nil {
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

// ParseSetMinLock is a log parse operation binding the contract event 0xabb05374bb45ebfef33afb21ec5aa52333708d8217fd8e5c0616efd2530b2145.
//
// Solidity: event SetMinLock(uint256 _newMinLock)
func (_Stakinginfo *StakinginfoFilterer) ParseSetMinLock(log types.Log) (*StakinginfoSetMinLock, error) {
	event := new(StakinginfoSetMinLock)
	if err := _Stakinginfo.contract.UnpackLog(event, "SetMinLock", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakinginfoSetRewardPayerIterator is returned from FilterSetRewardPayer and is used to iterate over the raw logs and unpacked data for SetRewardPayer events raised by the Stakinginfo contract.
type StakinginfoSetRewardPayerIterator struct {
	Event *StakinginfoSetRewardPayer // Event containing the contract specifics and raw log

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
func (it *StakinginfoSetRewardPayerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoSetRewardPayer)
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
		it.Event = new(StakinginfoSetRewardPayer)
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
func (it *StakinginfoSetRewardPayerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoSetRewardPayerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoSetRewardPayer represents a SetRewardPayer event raised by the Stakinginfo contract.
type StakinginfoSetRewardPayer struct {
	Payer common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterSetRewardPayer is a free log retrieval operation binding the contract event 0x30b92f5a89d7473895c4e9ce260fa7d0eefef2d59d5e68192e2e8cca4b9473a0.
//
// Solidity: event SetRewardPayer(address _payer)
func (_Stakinginfo *StakinginfoFilterer) FilterSetRewardPayer(opts *bind.FilterOpts) (*StakinginfoSetRewardPayerIterator, error) {

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "SetRewardPayer")
	if err != nil {
		return nil, err
	}
	return &StakinginfoSetRewardPayerIterator{contract: _Stakinginfo.contract, event: "SetRewardPayer", logs: logs, sub: sub}, nil
}

// WatchSetRewardPayer is a free log subscription operation binding the contract event 0x30b92f5a89d7473895c4e9ce260fa7d0eefef2d59d5e68192e2e8cca4b9473a0.
//
// Solidity: event SetRewardPayer(address _payer)
func (_Stakinginfo *StakinginfoFilterer) WatchSetRewardPayer(opts *bind.WatchOpts, sink chan<- *StakinginfoSetRewardPayer) (event.Subscription, error) {

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "SetRewardPayer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoSetRewardPayer)
				if err := _Stakinginfo.contract.UnpackLog(event, "SetRewardPayer", log); err != nil {
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

// ParseSetRewardPayer is a log parse operation binding the contract event 0x30b92f5a89d7473895c4e9ce260fa7d0eefef2d59d5e68192e2e8cca4b9473a0.
//
// Solidity: event SetRewardPayer(address _payer)
func (_Stakinginfo *StakinginfoFilterer) ParseSetRewardPayer(log types.Log) (*StakinginfoSetRewardPayer, error) {
	event := new(StakinginfoSetRewardPayer)
	if err := _Stakinginfo.contract.UnpackLog(event, "SetRewardPayer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakinginfoSignerChangeIterator is returned from FilterSignerChange and is used to iterate over the raw logs and unpacked data for SignerChange events raised by the Stakinginfo contract.
type StakinginfoSignerChangeIterator struct {
	Event *StakinginfoSignerChange // Event containing the contract specifics and raw log

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
func (it *StakinginfoSignerChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoSignerChange)
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
		it.Event = new(StakinginfoSignerChange)
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
func (it *StakinginfoSignerChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoSignerChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoSignerChange represents a SignerChange event raised by the Stakinginfo contract.
type StakinginfoSignerChange struct {
	SequencerId  *big.Int
	Nonce        *big.Int
	OldSigner    common.Address
	NewSigner    common.Address
	SignerPubkey []byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterSignerChange is a free log retrieval operation binding the contract event 0x086044c0612a8c965d4cccd907f0d588e40ad68438bd4c1274cac60f4c3a9d1f.
//
// Solidity: event SignerChange(uint256 indexed sequencerId, uint256 nonce, address indexed oldSigner, address indexed newSigner, bytes signerPubkey)
func (_Stakinginfo *StakinginfoFilterer) FilterSignerChange(opts *bind.FilterOpts, sequencerId []*big.Int, oldSigner []common.Address, newSigner []common.Address) (*StakinginfoSignerChangeIterator, error) {

	var sequencerIdRule []interface{}
	for _, sequencerIdItem := range sequencerId {
		sequencerIdRule = append(sequencerIdRule, sequencerIdItem)
	}

	var oldSignerRule []interface{}
	for _, oldSignerItem := range oldSigner {
		oldSignerRule = append(oldSignerRule, oldSignerItem)
	}
	var newSignerRule []interface{}
	for _, newSignerItem := range newSigner {
		newSignerRule = append(newSignerRule, newSignerItem)
	}

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "SignerChange", sequencerIdRule, oldSignerRule, newSignerRule)
	if err != nil {
		return nil, err
	}
	return &StakinginfoSignerChangeIterator{contract: _Stakinginfo.contract, event: "SignerChange", logs: logs, sub: sub}, nil
}

// WatchSignerChange is a free log subscription operation binding the contract event 0x086044c0612a8c965d4cccd907f0d588e40ad68438bd4c1274cac60f4c3a9d1f.
//
// Solidity: event SignerChange(uint256 indexed sequencerId, uint256 nonce, address indexed oldSigner, address indexed newSigner, bytes signerPubkey)
func (_Stakinginfo *StakinginfoFilterer) WatchSignerChange(opts *bind.WatchOpts, sink chan<- *StakinginfoSignerChange, sequencerId []*big.Int, oldSigner []common.Address, newSigner []common.Address) (event.Subscription, error) {

	var sequencerIdRule []interface{}
	for _, sequencerIdItem := range sequencerId {
		sequencerIdRule = append(sequencerIdRule, sequencerIdItem)
	}

	var oldSignerRule []interface{}
	for _, oldSignerItem := range oldSigner {
		oldSignerRule = append(oldSignerRule, oldSignerItem)
	}
	var newSignerRule []interface{}
	for _, newSignerItem := range newSigner {
		newSignerRule = append(newSignerRule, newSignerItem)
	}

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "SignerChange", sequencerIdRule, oldSignerRule, newSignerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoSignerChange)
				if err := _Stakinginfo.contract.UnpackLog(event, "SignerChange", log); err != nil {
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

// ParseSignerChange is a log parse operation binding the contract event 0x086044c0612a8c965d4cccd907f0d588e40ad68438bd4c1274cac60f4c3a9d1f.
//
// Solidity: event SignerChange(uint256 indexed sequencerId, uint256 nonce, address indexed oldSigner, address indexed newSigner, bytes signerPubkey)
func (_Stakinginfo *StakinginfoFilterer) ParseSignerChange(log types.Log) (*StakinginfoSignerChange, error) {
	event := new(StakinginfoSignerChange)
	if err := _Stakinginfo.contract.UnpackLog(event, "SignerChange", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakinginfoUnlockInitIterator is returned from FilterUnlockInit and is used to iterate over the raw logs and unpacked data for UnlockInit events raised by the Stakinginfo contract.
type StakinginfoUnlockInitIterator struct {
	Event *StakinginfoUnlockInit // Event containing the contract specifics and raw log

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
func (it *StakinginfoUnlockInitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoUnlockInit)
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
		it.Event = new(StakinginfoUnlockInit)
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
func (it *StakinginfoUnlockInitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoUnlockInitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoUnlockInit represents a UnlockInit event raised by the Stakinginfo contract.
type StakinginfoUnlockInit struct {
	User              common.Address
	SequencerId       *big.Int
	Nonce             *big.Int
	DeactivationBatch *big.Int
	DeactivationTime  *big.Int
	UnlockClaimTime   *big.Int
	Amount            *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterUnlockInit is a free log retrieval operation binding the contract event 0x06d9e13438f0daf13a71d63f3f8579db8bdeb299e4b651942313c73224d7af69.
//
// Solidity: event UnlockInit(address indexed user, uint256 indexed sequencerId, uint256 nonce, uint256 deactivationBatch, uint256 deactivationTime, uint256 unlockClaimTime, uint256 indexed amount)
func (_Stakinginfo *StakinginfoFilterer) FilterUnlockInit(opts *bind.FilterOpts, user []common.Address, sequencerId []*big.Int, amount []*big.Int) (*StakinginfoUnlockInitIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var sequencerIdRule []interface{}
	for _, sequencerIdItem := range sequencerId {
		sequencerIdRule = append(sequencerIdRule, sequencerIdItem)
	}

	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "UnlockInit", userRule, sequencerIdRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &StakinginfoUnlockInitIterator{contract: _Stakinginfo.contract, event: "UnlockInit", logs: logs, sub: sub}, nil
}

// WatchUnlockInit is a free log subscription operation binding the contract event 0x06d9e13438f0daf13a71d63f3f8579db8bdeb299e4b651942313c73224d7af69.
//
// Solidity: event UnlockInit(address indexed user, uint256 indexed sequencerId, uint256 nonce, uint256 deactivationBatch, uint256 deactivationTime, uint256 unlockClaimTime, uint256 indexed amount)
func (_Stakinginfo *StakinginfoFilterer) WatchUnlockInit(opts *bind.WatchOpts, sink chan<- *StakinginfoUnlockInit, user []common.Address, sequencerId []*big.Int, amount []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var sequencerIdRule []interface{}
	for _, sequencerIdItem := range sequencerId {
		sequencerIdRule = append(sequencerIdRule, sequencerIdItem)
	}

	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "UnlockInit", userRule, sequencerIdRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoUnlockInit)
				if err := _Stakinginfo.contract.UnpackLog(event, "UnlockInit", log); err != nil {
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

// ParseUnlockInit is a log parse operation binding the contract event 0x06d9e13438f0daf13a71d63f3f8579db8bdeb299e4b651942313c73224d7af69.
//
// Solidity: event UnlockInit(address indexed user, uint256 indexed sequencerId, uint256 nonce, uint256 deactivationBatch, uint256 deactivationTime, uint256 unlockClaimTime, uint256 indexed amount)
func (_Stakinginfo *StakinginfoFilterer) ParseUnlockInit(log types.Log) (*StakinginfoUnlockInit, error) {
	event := new(StakinginfoUnlockInit)
	if err := _Stakinginfo.contract.UnpackLog(event, "UnlockInit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakinginfoUnlockedIterator is returned from FilterUnlocked and is used to iterate over the raw logs and unpacked data for Unlocked events raised by the Stakinginfo contract.
type StakinginfoUnlockedIterator struct {
	Event *StakinginfoUnlocked // Event containing the contract specifics and raw log

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
func (it *StakinginfoUnlockedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoUnlocked)
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
		it.Event = new(StakinginfoUnlocked)
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
func (it *StakinginfoUnlockedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoUnlockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoUnlocked represents a Unlocked event raised by the Stakinginfo contract.
type StakinginfoUnlocked struct {
	User        common.Address
	SequencerId *big.Int
	Amount      *big.Int
	Total       *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterUnlocked is a free log retrieval operation binding the contract event 0x5245d528087a96a64f4589a764f00061e4671eab90cb1e019b1a5b24b2e4c2a8.
//
// Solidity: event Unlocked(address indexed user, uint256 indexed sequencerId, uint256 amount, uint256 total)
func (_Stakinginfo *StakinginfoFilterer) FilterUnlocked(opts *bind.FilterOpts, user []common.Address, sequencerId []*big.Int) (*StakinginfoUnlockedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var sequencerIdRule []interface{}
	for _, sequencerIdItem := range sequencerId {
		sequencerIdRule = append(sequencerIdRule, sequencerIdItem)
	}

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "Unlocked", userRule, sequencerIdRule)
	if err != nil {
		return nil, err
	}
	return &StakinginfoUnlockedIterator{contract: _Stakinginfo.contract, event: "Unlocked", logs: logs, sub: sub}, nil
}

// WatchUnlocked is a free log subscription operation binding the contract event 0x5245d528087a96a64f4589a764f00061e4671eab90cb1e019b1a5b24b2e4c2a8.
//
// Solidity: event Unlocked(address indexed user, uint256 indexed sequencerId, uint256 amount, uint256 total)
func (_Stakinginfo *StakinginfoFilterer) WatchUnlocked(opts *bind.WatchOpts, sink chan<- *StakinginfoUnlocked, user []common.Address, sequencerId []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var sequencerIdRule []interface{}
	for _, sequencerIdItem := range sequencerId {
		sequencerIdRule = append(sequencerIdRule, sequencerIdItem)
	}

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "Unlocked", userRule, sequencerIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoUnlocked)
				if err := _Stakinginfo.contract.UnpackLog(event, "Unlocked", log); err != nil {
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

// ParseUnlocked is a log parse operation binding the contract event 0x5245d528087a96a64f4589a764f00061e4671eab90cb1e019b1a5b24b2e4c2a8.
//
// Solidity: event Unlocked(address indexed user, uint256 indexed sequencerId, uint256 amount, uint256 total)
func (_Stakinginfo *StakinginfoFilterer) ParseUnlocked(log types.Log) (*StakinginfoUnlocked, error) {
	event := new(StakinginfoUnlocked)
	if err := _Stakinginfo.contract.UnpackLog(event, "Unlocked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
