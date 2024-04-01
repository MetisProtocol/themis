// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package sequencerset

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

// MetisSequencerSetEpoch is an auto generated low-level Go binding around an user-defined struct.
type MetisSequencerSetEpoch struct {
	Number     *big.Int
	Signer     common.Address
	StartBlock *big.Int
	EndBlock   *big.Int
}

// SequencersetMetaData contains all meta data concerning the Sequencerset contract.
var SequencersetMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_newLength\",\"type\":\"uint256\"}],\"name\":\"EpochUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_newMpcAddress\",\"type\":\"address\"}],\"name\":\"MpcAddressUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"epochId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"endBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"NewEpoch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"oldEpochId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"newEpochId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"curEpochId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"endBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newSigner\",\"type\":\"address\"}],\"name\":\"ReCommitEpoch\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newLength\",\"type\":\"uint256\"}],\"name\":\"UpdateEpochLength\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newMpc\",\"type\":\"address\"}],\"name\":\"UpdateMpcAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endBlock\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"commitEpoch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentEpoch\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"number\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endBlock\",\"type\":\"uint256\"}],\"internalType\":\"structMetisSequencerSet.Epoch\",\"name\":\"epoch\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentEpochNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"epochs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"number\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endBlock\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"number\",\"type\":\"uint256\"}],\"name\":\"getEpochByBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"number\",\"type\":\"uint256\"}],\"name\":\"getMetisSequencer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_initialSequencer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_mpcAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_firstStartBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_firstEndBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_epochLength\",\"type\":\"uint256\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mpcAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"oldEpochId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"newEpochId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endBlock\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newSigner\",\"type\":\"address\"}],\"name\":\"recommitEpoch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// SequencersetABI is the input ABI used to generate the binding from.
// Deprecated: Use SequencersetMetaData.ABI instead.
var SequencersetABI = SequencersetMetaData.ABI

// Sequencerset is an auto generated Go binding around an Ethereum contract.
type Sequencerset struct {
	SequencersetCaller     // Read-only binding to the contract
	SequencersetTransactor // Write-only binding to the contract
	SequencersetFilterer   // Log filterer for contract events
}

// SequencersetCaller is an auto generated read-only Go binding around an Ethereum contract.
type SequencersetCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SequencersetTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SequencersetTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SequencersetFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SequencersetFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SequencersetSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SequencersetSession struct {
	Contract     *Sequencerset     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SequencersetCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SequencersetCallerSession struct {
	Contract *SequencersetCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// SequencersetTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SequencersetTransactorSession struct {
	Contract     *SequencersetTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// SequencersetRaw is an auto generated low-level Go binding around an Ethereum contract.
type SequencersetRaw struct {
	Contract *Sequencerset // Generic contract binding to access the raw methods on
}

// SequencersetCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SequencersetCallerRaw struct {
	Contract *SequencersetCaller // Generic read-only contract binding to access the raw methods on
}

// SequencersetTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SequencersetTransactorRaw struct {
	Contract *SequencersetTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSequencerset creates a new instance of Sequencerset, bound to a specific deployed contract.
func NewSequencerset(address common.Address, backend bind.ContractBackend) (*Sequencerset, error) {
	contract, err := bindSequencerset(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Sequencerset{SequencersetCaller: SequencersetCaller{contract: contract}, SequencersetTransactor: SequencersetTransactor{contract: contract}, SequencersetFilterer: SequencersetFilterer{contract: contract}}, nil
}

// NewSequencersetCaller creates a new read-only instance of Sequencerset, bound to a specific deployed contract.
func NewSequencersetCaller(address common.Address, caller bind.ContractCaller) (*SequencersetCaller, error) {
	contract, err := bindSequencerset(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SequencersetCaller{contract: contract}, nil
}

// NewSequencersetTransactor creates a new write-only instance of Sequencerset, bound to a specific deployed contract.
func NewSequencersetTransactor(address common.Address, transactor bind.ContractTransactor) (*SequencersetTransactor, error) {
	contract, err := bindSequencerset(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SequencersetTransactor{contract: contract}, nil
}

// NewSequencersetFilterer creates a new log filterer instance of Sequencerset, bound to a specific deployed contract.
func NewSequencersetFilterer(address common.Address, filterer bind.ContractFilterer) (*SequencersetFilterer, error) {
	contract, err := bindSequencerset(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SequencersetFilterer{contract: contract}, nil
}

// bindSequencerset binds a generic wrapper to an already deployed contract.
func bindSequencerset(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SequencersetMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Sequencerset *SequencersetRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Sequencerset.Contract.SequencersetCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Sequencerset *SequencersetRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Sequencerset.Contract.SequencersetTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Sequencerset *SequencersetRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Sequencerset.Contract.SequencersetTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Sequencerset *SequencersetCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Sequencerset.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Sequencerset *SequencersetTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Sequencerset.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Sequencerset *SequencersetTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Sequencerset.Contract.contract.Transact(opts, method, params...)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns((uint256,address,uint256,uint256) epoch)
func (_Sequencerset *SequencersetCaller) CurrentEpoch(opts *bind.CallOpts) (MetisSequencerSetEpoch, error) {
	var out []interface{}
	err := _Sequencerset.contract.Call(opts, &out, "currentEpoch")

	if err != nil {
		return *new(MetisSequencerSetEpoch), err
	}

	out0 := *abi.ConvertType(out[0], new(MetisSequencerSetEpoch)).(*MetisSequencerSetEpoch)

	return out0, err

}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns((uint256,address,uint256,uint256) epoch)
func (_Sequencerset *SequencersetSession) CurrentEpoch() (MetisSequencerSetEpoch, error) {
	return _Sequencerset.Contract.CurrentEpoch(&_Sequencerset.CallOpts)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns((uint256,address,uint256,uint256) epoch)
func (_Sequencerset *SequencersetCallerSession) CurrentEpoch() (MetisSequencerSetEpoch, error) {
	return _Sequencerset.Contract.CurrentEpoch(&_Sequencerset.CallOpts)
}

// CurrentEpochNumber is a free data retrieval call binding the contract method 0x6903beb4.
//
// Solidity: function currentEpochNumber() view returns(uint256)
func (_Sequencerset *SequencersetCaller) CurrentEpochNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Sequencerset.contract.Call(opts, &out, "currentEpochNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentEpochNumber is a free data retrieval call binding the contract method 0x6903beb4.
//
// Solidity: function currentEpochNumber() view returns(uint256)
func (_Sequencerset *SequencersetSession) CurrentEpochNumber() (*big.Int, error) {
	return _Sequencerset.Contract.CurrentEpochNumber(&_Sequencerset.CallOpts)
}

// CurrentEpochNumber is a free data retrieval call binding the contract method 0x6903beb4.
//
// Solidity: function currentEpochNumber() view returns(uint256)
func (_Sequencerset *SequencersetCallerSession) CurrentEpochNumber() (*big.Int, error) {
	return _Sequencerset.Contract.CurrentEpochNumber(&_Sequencerset.CallOpts)
}

// EpochLength is a free data retrieval call binding the contract method 0x57d775f8.
//
// Solidity: function epochLength() view returns(uint256)
func (_Sequencerset *SequencersetCaller) EpochLength(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Sequencerset.contract.Call(opts, &out, "epochLength")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EpochLength is a free data retrieval call binding the contract method 0x57d775f8.
//
// Solidity: function epochLength() view returns(uint256)
func (_Sequencerset *SequencersetSession) EpochLength() (*big.Int, error) {
	return _Sequencerset.Contract.EpochLength(&_Sequencerset.CallOpts)
}

// EpochLength is a free data retrieval call binding the contract method 0x57d775f8.
//
// Solidity: function epochLength() view returns(uint256)
func (_Sequencerset *SequencersetCallerSession) EpochLength() (*big.Int, error) {
	return _Sequencerset.Contract.EpochLength(&_Sequencerset.CallOpts)
}

// Epochs is a free data retrieval call binding the contract method 0xc6b61e4c.
//
// Solidity: function epochs(uint256 ) view returns(uint256 number, address signer, uint256 startBlock, uint256 endBlock)
func (_Sequencerset *SequencersetCaller) Epochs(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Number     *big.Int
	Signer     common.Address
	StartBlock *big.Int
	EndBlock   *big.Int
}, error) {
	var out []interface{}
	err := _Sequencerset.contract.Call(opts, &out, "epochs", arg0)

	outstruct := new(struct {
		Number     *big.Int
		Signer     common.Address
		StartBlock *big.Int
		EndBlock   *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Number = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Signer = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.StartBlock = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.EndBlock = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Epochs is a free data retrieval call binding the contract method 0xc6b61e4c.
//
// Solidity: function epochs(uint256 ) view returns(uint256 number, address signer, uint256 startBlock, uint256 endBlock)
func (_Sequencerset *SequencersetSession) Epochs(arg0 *big.Int) (struct {
	Number     *big.Int
	Signer     common.Address
	StartBlock *big.Int
	EndBlock   *big.Int
}, error) {
	return _Sequencerset.Contract.Epochs(&_Sequencerset.CallOpts, arg0)
}

// Epochs is a free data retrieval call binding the contract method 0xc6b61e4c.
//
// Solidity: function epochs(uint256 ) view returns(uint256 number, address signer, uint256 startBlock, uint256 endBlock)
func (_Sequencerset *SequencersetCallerSession) Epochs(arg0 *big.Int) (struct {
	Number     *big.Int
	Signer     common.Address
	StartBlock *big.Int
	EndBlock   *big.Int
}, error) {
	return _Sequencerset.Contract.Epochs(&_Sequencerset.CallOpts, arg0)
}

// GetEpochByBlock is a free data retrieval call binding the contract method 0x46df33d2.
//
// Solidity: function getEpochByBlock(uint256 number) view returns(uint256)
func (_Sequencerset *SequencersetCaller) GetEpochByBlock(opts *bind.CallOpts, number *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Sequencerset.contract.Call(opts, &out, "getEpochByBlock", number)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetEpochByBlock is a free data retrieval call binding the contract method 0x46df33d2.
//
// Solidity: function getEpochByBlock(uint256 number) view returns(uint256)
func (_Sequencerset *SequencersetSession) GetEpochByBlock(number *big.Int) (*big.Int, error) {
	return _Sequencerset.Contract.GetEpochByBlock(&_Sequencerset.CallOpts, number)
}

// GetEpochByBlock is a free data retrieval call binding the contract method 0x46df33d2.
//
// Solidity: function getEpochByBlock(uint256 number) view returns(uint256)
func (_Sequencerset *SequencersetCallerSession) GetEpochByBlock(number *big.Int) (*big.Int, error) {
	return _Sequencerset.Contract.GetEpochByBlock(&_Sequencerset.CallOpts, number)
}

// GetMetisSequencer is a free data retrieval call binding the contract method 0x3edae769.
//
// Solidity: function getMetisSequencer(uint256 number) view returns(address)
func (_Sequencerset *SequencersetCaller) GetMetisSequencer(opts *bind.CallOpts, number *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Sequencerset.contract.Call(opts, &out, "getMetisSequencer", number)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetMetisSequencer is a free data retrieval call binding the contract method 0x3edae769.
//
// Solidity: function getMetisSequencer(uint256 number) view returns(address)
func (_Sequencerset *SequencersetSession) GetMetisSequencer(number *big.Int) (common.Address, error) {
	return _Sequencerset.Contract.GetMetisSequencer(&_Sequencerset.CallOpts, number)
}

// GetMetisSequencer is a free data retrieval call binding the contract method 0x3edae769.
//
// Solidity: function getMetisSequencer(uint256 number) view returns(address)
func (_Sequencerset *SequencersetCallerSession) GetMetisSequencer(number *big.Int) (common.Address, error) {
	return _Sequencerset.Contract.GetMetisSequencer(&_Sequencerset.CallOpts, number)
}

// MpcAddress is a free data retrieval call binding the contract method 0x111f4630.
//
// Solidity: function mpcAddress() view returns(address)
func (_Sequencerset *SequencersetCaller) MpcAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Sequencerset.contract.Call(opts, &out, "mpcAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MpcAddress is a free data retrieval call binding the contract method 0x111f4630.
//
// Solidity: function mpcAddress() view returns(address)
func (_Sequencerset *SequencersetSession) MpcAddress() (common.Address, error) {
	return _Sequencerset.Contract.MpcAddress(&_Sequencerset.CallOpts)
}

// MpcAddress is a free data retrieval call binding the contract method 0x111f4630.
//
// Solidity: function mpcAddress() view returns(address)
func (_Sequencerset *SequencersetCallerSession) MpcAddress() (common.Address, error) {
	return _Sequencerset.Contract.MpcAddress(&_Sequencerset.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Sequencerset *SequencersetCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Sequencerset.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Sequencerset *SequencersetSession) Owner() (common.Address, error) {
	return _Sequencerset.Contract.Owner(&_Sequencerset.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Sequencerset *SequencersetCallerSession) Owner() (common.Address, error) {
	return _Sequencerset.Contract.Owner(&_Sequencerset.CallOpts)
}

// UpdateEpochLength is a paid mutator transaction binding the contract method 0x24316ccb.
//
// Solidity: function UpdateEpochLength(uint256 _newLength) returns()
func (_Sequencerset *SequencersetTransactor) UpdateEpochLength(opts *bind.TransactOpts, _newLength *big.Int) (*types.Transaction, error) {
	return _Sequencerset.contract.Transact(opts, "UpdateEpochLength", _newLength)
}

// UpdateEpochLength is a paid mutator transaction binding the contract method 0x24316ccb.
//
// Solidity: function UpdateEpochLength(uint256 _newLength) returns()
func (_Sequencerset *SequencersetSession) UpdateEpochLength(_newLength *big.Int) (*types.Transaction, error) {
	return _Sequencerset.Contract.UpdateEpochLength(&_Sequencerset.TransactOpts, _newLength)
}

// UpdateEpochLength is a paid mutator transaction binding the contract method 0x24316ccb.
//
// Solidity: function UpdateEpochLength(uint256 _newLength) returns()
func (_Sequencerset *SequencersetTransactorSession) UpdateEpochLength(_newLength *big.Int) (*types.Transaction, error) {
	return _Sequencerset.Contract.UpdateEpochLength(&_Sequencerset.TransactOpts, _newLength)
}

// UpdateMpcAddress is a paid mutator transaction binding the contract method 0x643dbfce.
//
// Solidity: function UpdateMpcAddress(address _newMpc) returns()
func (_Sequencerset *SequencersetTransactor) UpdateMpcAddress(opts *bind.TransactOpts, _newMpc common.Address) (*types.Transaction, error) {
	return _Sequencerset.contract.Transact(opts, "UpdateMpcAddress", _newMpc)
}

// UpdateMpcAddress is a paid mutator transaction binding the contract method 0x643dbfce.
//
// Solidity: function UpdateMpcAddress(address _newMpc) returns()
func (_Sequencerset *SequencersetSession) UpdateMpcAddress(_newMpc common.Address) (*types.Transaction, error) {
	return _Sequencerset.Contract.UpdateMpcAddress(&_Sequencerset.TransactOpts, _newMpc)
}

// UpdateMpcAddress is a paid mutator transaction binding the contract method 0x643dbfce.
//
// Solidity: function UpdateMpcAddress(address _newMpc) returns()
func (_Sequencerset *SequencersetTransactorSession) UpdateMpcAddress(_newMpc common.Address) (*types.Transaction, error) {
	return _Sequencerset.Contract.UpdateMpcAddress(&_Sequencerset.TransactOpts, _newMpc)
}

// CommitEpoch is a paid mutator transaction binding the contract method 0x4fb71bdd.
//
// Solidity: function commitEpoch(uint256 newEpoch, uint256 startBlock, uint256 endBlock, address signer) returns()
func (_Sequencerset *SequencersetTransactor) CommitEpoch(opts *bind.TransactOpts, newEpoch *big.Int, startBlock *big.Int, endBlock *big.Int, signer common.Address) (*types.Transaction, error) {
	return _Sequencerset.contract.Transact(opts, "commitEpoch", newEpoch, startBlock, endBlock, signer)
}

// CommitEpoch is a paid mutator transaction binding the contract method 0x4fb71bdd.
//
// Solidity: function commitEpoch(uint256 newEpoch, uint256 startBlock, uint256 endBlock, address signer) returns()
func (_Sequencerset *SequencersetSession) CommitEpoch(newEpoch *big.Int, startBlock *big.Int, endBlock *big.Int, signer common.Address) (*types.Transaction, error) {
	return _Sequencerset.Contract.CommitEpoch(&_Sequencerset.TransactOpts, newEpoch, startBlock, endBlock, signer)
}

// CommitEpoch is a paid mutator transaction binding the contract method 0x4fb71bdd.
//
// Solidity: function commitEpoch(uint256 newEpoch, uint256 startBlock, uint256 endBlock, address signer) returns()
func (_Sequencerset *SequencersetTransactorSession) CommitEpoch(newEpoch *big.Int, startBlock *big.Int, endBlock *big.Int, signer common.Address) (*types.Transaction, error) {
	return _Sequencerset.Contract.CommitEpoch(&_Sequencerset.TransactOpts, newEpoch, startBlock, endBlock, signer)
}

// Initialize is a paid mutator transaction binding the contract method 0xd13f90b4.
//
// Solidity: function initialize(address _initialSequencer, address _mpcAddress, uint256 _firstStartBlock, uint256 _firstEndBlock, uint256 _epochLength) returns()
func (_Sequencerset *SequencersetTransactor) Initialize(opts *bind.TransactOpts, _initialSequencer common.Address, _mpcAddress common.Address, _firstStartBlock *big.Int, _firstEndBlock *big.Int, _epochLength *big.Int) (*types.Transaction, error) {
	return _Sequencerset.contract.Transact(opts, "initialize", _initialSequencer, _mpcAddress, _firstStartBlock, _firstEndBlock, _epochLength)
}

// Initialize is a paid mutator transaction binding the contract method 0xd13f90b4.
//
// Solidity: function initialize(address _initialSequencer, address _mpcAddress, uint256 _firstStartBlock, uint256 _firstEndBlock, uint256 _epochLength) returns()
func (_Sequencerset *SequencersetSession) Initialize(_initialSequencer common.Address, _mpcAddress common.Address, _firstStartBlock *big.Int, _firstEndBlock *big.Int, _epochLength *big.Int) (*types.Transaction, error) {
	return _Sequencerset.Contract.Initialize(&_Sequencerset.TransactOpts, _initialSequencer, _mpcAddress, _firstStartBlock, _firstEndBlock, _epochLength)
}

// Initialize is a paid mutator transaction binding the contract method 0xd13f90b4.
//
// Solidity: function initialize(address _initialSequencer, address _mpcAddress, uint256 _firstStartBlock, uint256 _firstEndBlock, uint256 _epochLength) returns()
func (_Sequencerset *SequencersetTransactorSession) Initialize(_initialSequencer common.Address, _mpcAddress common.Address, _firstStartBlock *big.Int, _firstEndBlock *big.Int, _epochLength *big.Int) (*types.Transaction, error) {
	return _Sequencerset.Contract.Initialize(&_Sequencerset.TransactOpts, _initialSequencer, _mpcAddress, _firstStartBlock, _firstEndBlock, _epochLength)
}

// RecommitEpoch is a paid mutator transaction binding the contract method 0x2c91c679.
//
// Solidity: function recommitEpoch(uint256 oldEpochId, uint256 newEpochId, uint256 startBlock, uint256 endBlock, address newSigner) returns()
func (_Sequencerset *SequencersetTransactor) RecommitEpoch(opts *bind.TransactOpts, oldEpochId *big.Int, newEpochId *big.Int, startBlock *big.Int, endBlock *big.Int, newSigner common.Address) (*types.Transaction, error) {
	return _Sequencerset.contract.Transact(opts, "recommitEpoch", oldEpochId, newEpochId, startBlock, endBlock, newSigner)
}

// RecommitEpoch is a paid mutator transaction binding the contract method 0x2c91c679.
//
// Solidity: function recommitEpoch(uint256 oldEpochId, uint256 newEpochId, uint256 startBlock, uint256 endBlock, address newSigner) returns()
func (_Sequencerset *SequencersetSession) RecommitEpoch(oldEpochId *big.Int, newEpochId *big.Int, startBlock *big.Int, endBlock *big.Int, newSigner common.Address) (*types.Transaction, error) {
	return _Sequencerset.Contract.RecommitEpoch(&_Sequencerset.TransactOpts, oldEpochId, newEpochId, startBlock, endBlock, newSigner)
}

// RecommitEpoch is a paid mutator transaction binding the contract method 0x2c91c679.
//
// Solidity: function recommitEpoch(uint256 oldEpochId, uint256 newEpochId, uint256 startBlock, uint256 endBlock, address newSigner) returns()
func (_Sequencerset *SequencersetTransactorSession) RecommitEpoch(oldEpochId *big.Int, newEpochId *big.Int, startBlock *big.Int, endBlock *big.Int, newSigner common.Address) (*types.Transaction, error) {
	return _Sequencerset.Contract.RecommitEpoch(&_Sequencerset.TransactOpts, oldEpochId, newEpochId, startBlock, endBlock, newSigner)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Sequencerset *SequencersetTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Sequencerset.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Sequencerset *SequencersetSession) RenounceOwnership() (*types.Transaction, error) {
	return _Sequencerset.Contract.RenounceOwnership(&_Sequencerset.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Sequencerset *SequencersetTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Sequencerset.Contract.RenounceOwnership(&_Sequencerset.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Sequencerset *SequencersetTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Sequencerset.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Sequencerset *SequencersetSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Sequencerset.Contract.TransferOwnership(&_Sequencerset.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Sequencerset *SequencersetTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Sequencerset.Contract.TransferOwnership(&_Sequencerset.TransactOpts, newOwner)
}

// SequencersetEpochUpdatedIterator is returned from FilterEpochUpdated and is used to iterate over the raw logs and unpacked data for EpochUpdated events raised by the Sequencerset contract.
type SequencersetEpochUpdatedIterator struct {
	Event *SequencersetEpochUpdated // Event containing the contract specifics and raw log

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
func (it *SequencersetEpochUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SequencersetEpochUpdated)
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
		it.Event = new(SequencersetEpochUpdated)
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
func (it *SequencersetEpochUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SequencersetEpochUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SequencersetEpochUpdated represents a EpochUpdated event raised by the Sequencerset contract.
type SequencersetEpochUpdated struct {
	NewLength *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterEpochUpdated is a free log retrieval operation binding the contract event 0xb33a1f54dde4e0082c45281b338d78b2c4b5be163b6ffffa5d0d6d1050ba5a58.
//
// Solidity: event EpochUpdated(uint256 _newLength)
func (_Sequencerset *SequencersetFilterer) FilterEpochUpdated(opts *bind.FilterOpts) (*SequencersetEpochUpdatedIterator, error) {

	logs, sub, err := _Sequencerset.contract.FilterLogs(opts, "EpochUpdated")
	if err != nil {
		return nil, err
	}
	return &SequencersetEpochUpdatedIterator{contract: _Sequencerset.contract, event: "EpochUpdated", logs: logs, sub: sub}, nil
}

// WatchEpochUpdated is a free log subscription operation binding the contract event 0xb33a1f54dde4e0082c45281b338d78b2c4b5be163b6ffffa5d0d6d1050ba5a58.
//
// Solidity: event EpochUpdated(uint256 _newLength)
func (_Sequencerset *SequencersetFilterer) WatchEpochUpdated(opts *bind.WatchOpts, sink chan<- *SequencersetEpochUpdated) (event.Subscription, error) {

	logs, sub, err := _Sequencerset.contract.WatchLogs(opts, "EpochUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SequencersetEpochUpdated)
				if err := _Sequencerset.contract.UnpackLog(event, "EpochUpdated", log); err != nil {
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

// ParseEpochUpdated is a log parse operation binding the contract event 0xb33a1f54dde4e0082c45281b338d78b2c4b5be163b6ffffa5d0d6d1050ba5a58.
//
// Solidity: event EpochUpdated(uint256 _newLength)
func (_Sequencerset *SequencersetFilterer) ParseEpochUpdated(log types.Log) (*SequencersetEpochUpdated, error) {
	event := new(SequencersetEpochUpdated)
	if err := _Sequencerset.contract.UnpackLog(event, "EpochUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SequencersetInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Sequencerset contract.
type SequencersetInitializedIterator struct {
	Event *SequencersetInitialized // Event containing the contract specifics and raw log

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
func (it *SequencersetInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SequencersetInitialized)
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
		it.Event = new(SequencersetInitialized)
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
func (it *SequencersetInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SequencersetInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SequencersetInitialized represents a Initialized event raised by the Sequencerset contract.
type SequencersetInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Sequencerset *SequencersetFilterer) FilterInitialized(opts *bind.FilterOpts) (*SequencersetInitializedIterator, error) {

	logs, sub, err := _Sequencerset.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &SequencersetInitializedIterator{contract: _Sequencerset.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Sequencerset *SequencersetFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *SequencersetInitialized) (event.Subscription, error) {

	logs, sub, err := _Sequencerset.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SequencersetInitialized)
				if err := _Sequencerset.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Sequencerset *SequencersetFilterer) ParseInitialized(log types.Log) (*SequencersetInitialized, error) {
	event := new(SequencersetInitialized)
	if err := _Sequencerset.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SequencersetMpcAddressUpdatedIterator is returned from FilterMpcAddressUpdated and is used to iterate over the raw logs and unpacked data for MpcAddressUpdated events raised by the Sequencerset contract.
type SequencersetMpcAddressUpdatedIterator struct {
	Event *SequencersetMpcAddressUpdated // Event containing the contract specifics and raw log

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
func (it *SequencersetMpcAddressUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SequencersetMpcAddressUpdated)
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
		it.Event = new(SequencersetMpcAddressUpdated)
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
func (it *SequencersetMpcAddressUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SequencersetMpcAddressUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SequencersetMpcAddressUpdated represents a MpcAddressUpdated event raised by the Sequencerset contract.
type SequencersetMpcAddressUpdated struct {
	NewMpcAddress common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterMpcAddressUpdated is a free log retrieval operation binding the contract event 0x9416d5b743f4c1409b0213ee5d1e57e8515f4be68a32fbbda85f838891b421da.
//
// Solidity: event MpcAddressUpdated(address _newMpcAddress)
func (_Sequencerset *SequencersetFilterer) FilterMpcAddressUpdated(opts *bind.FilterOpts) (*SequencersetMpcAddressUpdatedIterator, error) {

	logs, sub, err := _Sequencerset.contract.FilterLogs(opts, "MpcAddressUpdated")
	if err != nil {
		return nil, err
	}
	return &SequencersetMpcAddressUpdatedIterator{contract: _Sequencerset.contract, event: "MpcAddressUpdated", logs: logs, sub: sub}, nil
}

// WatchMpcAddressUpdated is a free log subscription operation binding the contract event 0x9416d5b743f4c1409b0213ee5d1e57e8515f4be68a32fbbda85f838891b421da.
//
// Solidity: event MpcAddressUpdated(address _newMpcAddress)
func (_Sequencerset *SequencersetFilterer) WatchMpcAddressUpdated(opts *bind.WatchOpts, sink chan<- *SequencersetMpcAddressUpdated) (event.Subscription, error) {

	logs, sub, err := _Sequencerset.contract.WatchLogs(opts, "MpcAddressUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SequencersetMpcAddressUpdated)
				if err := _Sequencerset.contract.UnpackLog(event, "MpcAddressUpdated", log); err != nil {
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

// ParseMpcAddressUpdated is a log parse operation binding the contract event 0x9416d5b743f4c1409b0213ee5d1e57e8515f4be68a32fbbda85f838891b421da.
//
// Solidity: event MpcAddressUpdated(address _newMpcAddress)
func (_Sequencerset *SequencersetFilterer) ParseMpcAddressUpdated(log types.Log) (*SequencersetMpcAddressUpdated, error) {
	event := new(SequencersetMpcAddressUpdated)
	if err := _Sequencerset.contract.UnpackLog(event, "MpcAddressUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SequencersetNewEpochIterator is returned from FilterNewEpoch and is used to iterate over the raw logs and unpacked data for NewEpoch events raised by the Sequencerset contract.
type SequencersetNewEpochIterator struct {
	Event *SequencersetNewEpoch // Event containing the contract specifics and raw log

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
func (it *SequencersetNewEpochIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SequencersetNewEpoch)
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
		it.Event = new(SequencersetNewEpoch)
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
func (it *SequencersetNewEpochIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SequencersetNewEpochIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SequencersetNewEpoch represents a NewEpoch event raised by the Sequencerset contract.
type SequencersetNewEpoch struct {
	EpochId    *big.Int
	StartBlock *big.Int
	EndBlock   *big.Int
	Signer     common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterNewEpoch is a free log retrieval operation binding the contract event 0x9030849b7c46dbbea7911d67c5814a2ab19c0704c448defc9b87589447844cc6.
//
// Solidity: event NewEpoch(uint256 indexed epochId, uint256 startBlock, uint256 endBlock, address signer)
func (_Sequencerset *SequencersetFilterer) FilterNewEpoch(opts *bind.FilterOpts, epochId []*big.Int) (*SequencersetNewEpochIterator, error) {

	var epochIdRule []interface{}
	for _, epochIdItem := range epochId {
		epochIdRule = append(epochIdRule, epochIdItem)
	}

	logs, sub, err := _Sequencerset.contract.FilterLogs(opts, "NewEpoch", epochIdRule)
	if err != nil {
		return nil, err
	}
	return &SequencersetNewEpochIterator{contract: _Sequencerset.contract, event: "NewEpoch", logs: logs, sub: sub}, nil
}

// WatchNewEpoch is a free log subscription operation binding the contract event 0x9030849b7c46dbbea7911d67c5814a2ab19c0704c448defc9b87589447844cc6.
//
// Solidity: event NewEpoch(uint256 indexed epochId, uint256 startBlock, uint256 endBlock, address signer)
func (_Sequencerset *SequencersetFilterer) WatchNewEpoch(opts *bind.WatchOpts, sink chan<- *SequencersetNewEpoch, epochId []*big.Int) (event.Subscription, error) {

	var epochIdRule []interface{}
	for _, epochIdItem := range epochId {
		epochIdRule = append(epochIdRule, epochIdItem)
	}

	logs, sub, err := _Sequencerset.contract.WatchLogs(opts, "NewEpoch", epochIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SequencersetNewEpoch)
				if err := _Sequencerset.contract.UnpackLog(event, "NewEpoch", log); err != nil {
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

// ParseNewEpoch is a log parse operation binding the contract event 0x9030849b7c46dbbea7911d67c5814a2ab19c0704c448defc9b87589447844cc6.
//
// Solidity: event NewEpoch(uint256 indexed epochId, uint256 startBlock, uint256 endBlock, address signer)
func (_Sequencerset *SequencersetFilterer) ParseNewEpoch(log types.Log) (*SequencersetNewEpoch, error) {
	event := new(SequencersetNewEpoch)
	if err := _Sequencerset.contract.UnpackLog(event, "NewEpoch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SequencersetOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Sequencerset contract.
type SequencersetOwnershipTransferredIterator struct {
	Event *SequencersetOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *SequencersetOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SequencersetOwnershipTransferred)
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
		it.Event = new(SequencersetOwnershipTransferred)
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
func (it *SequencersetOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SequencersetOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SequencersetOwnershipTransferred represents a OwnershipTransferred event raised by the Sequencerset contract.
type SequencersetOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Sequencerset *SequencersetFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SequencersetOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Sequencerset.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SequencersetOwnershipTransferredIterator{contract: _Sequencerset.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Sequencerset *SequencersetFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SequencersetOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Sequencerset.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SequencersetOwnershipTransferred)
				if err := _Sequencerset.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Sequencerset *SequencersetFilterer) ParseOwnershipTransferred(log types.Log) (*SequencersetOwnershipTransferred, error) {
	event := new(SequencersetOwnershipTransferred)
	if err := _Sequencerset.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SequencersetReCommitEpochIterator is returned from FilterReCommitEpoch and is used to iterate over the raw logs and unpacked data for ReCommitEpoch events raised by the Sequencerset contract.
type SequencersetReCommitEpochIterator struct {
	Event *SequencersetReCommitEpoch // Event containing the contract specifics and raw log

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
func (it *SequencersetReCommitEpochIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SequencersetReCommitEpoch)
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
		it.Event = new(SequencersetReCommitEpoch)
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
func (it *SequencersetReCommitEpochIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SequencersetReCommitEpochIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SequencersetReCommitEpoch represents a ReCommitEpoch event raised by the Sequencerset contract.
type SequencersetReCommitEpoch struct {
	OldEpochId *big.Int
	NewEpochId *big.Int
	CurEpochId *big.Int
	StartBlock *big.Int
	EndBlock   *big.Int
	NewSigner  common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterReCommitEpoch is a free log retrieval operation binding the contract event 0x2cfee2bd8227abd70bbef63ce6f35d2e365c731d553d9d2231073f80b01831da.
//
// Solidity: event ReCommitEpoch(uint256 indexed oldEpochId, uint256 indexed newEpochId, uint256 curEpochId, uint256 startBlock, uint256 endBlock, address newSigner)
func (_Sequencerset *SequencersetFilterer) FilterReCommitEpoch(opts *bind.FilterOpts, oldEpochId []*big.Int, newEpochId []*big.Int) (*SequencersetReCommitEpochIterator, error) {

	var oldEpochIdRule []interface{}
	for _, oldEpochIdItem := range oldEpochId {
		oldEpochIdRule = append(oldEpochIdRule, oldEpochIdItem)
	}
	var newEpochIdRule []interface{}
	for _, newEpochIdItem := range newEpochId {
		newEpochIdRule = append(newEpochIdRule, newEpochIdItem)
	}

	logs, sub, err := _Sequencerset.contract.FilterLogs(opts, "ReCommitEpoch", oldEpochIdRule, newEpochIdRule)
	if err != nil {
		return nil, err
	}
	return &SequencersetReCommitEpochIterator{contract: _Sequencerset.contract, event: "ReCommitEpoch", logs: logs, sub: sub}, nil
}

// WatchReCommitEpoch is a free log subscription operation binding the contract event 0x2cfee2bd8227abd70bbef63ce6f35d2e365c731d553d9d2231073f80b01831da.
//
// Solidity: event ReCommitEpoch(uint256 indexed oldEpochId, uint256 indexed newEpochId, uint256 curEpochId, uint256 startBlock, uint256 endBlock, address newSigner)
func (_Sequencerset *SequencersetFilterer) WatchReCommitEpoch(opts *bind.WatchOpts, sink chan<- *SequencersetReCommitEpoch, oldEpochId []*big.Int, newEpochId []*big.Int) (event.Subscription, error) {

	var oldEpochIdRule []interface{}
	for _, oldEpochIdItem := range oldEpochId {
		oldEpochIdRule = append(oldEpochIdRule, oldEpochIdItem)
	}
	var newEpochIdRule []interface{}
	for _, newEpochIdItem := range newEpochId {
		newEpochIdRule = append(newEpochIdRule, newEpochIdItem)
	}

	logs, sub, err := _Sequencerset.contract.WatchLogs(opts, "ReCommitEpoch", oldEpochIdRule, newEpochIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SequencersetReCommitEpoch)
				if err := _Sequencerset.contract.UnpackLog(event, "ReCommitEpoch", log); err != nil {
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

// ParseReCommitEpoch is a log parse operation binding the contract event 0x2cfee2bd8227abd70bbef63ce6f35d2e365c731d553d9d2231073f80b01831da.
//
// Solidity: event ReCommitEpoch(uint256 indexed oldEpochId, uint256 indexed newEpochId, uint256 curEpochId, uint256 startBlock, uint256 endBlock, address newSigner)
func (_Sequencerset *SequencersetFilterer) ParseReCommitEpoch(log types.Log) (*SequencersetReCommitEpoch, error) {
	event := new(SequencersetReCommitEpoch)
	if err := _Sequencerset.contract.UnpackLog(event, "ReCommitEpoch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
