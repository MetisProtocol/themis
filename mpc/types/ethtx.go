package types

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/holiman/uint256"
)

type UnsignedTxData interface {
	TxData() types.TxData
}

// UnsignedLegacyTx is the unsigned transaction data of the original Ethereum transactions.
type UnsignedLegacyTx struct {
	Nonce    uint64          // nonce of sender account
	GasPrice *big.Int        // wei per gas
	Gas      uint64          // gas limit
	To       *common.Address `rlp:"nil"` // nil means contract creation
	Value    *big.Int        // wei amount
	Data     []byte          // contract invocation input data
}

func (t *UnsignedLegacyTx) TxData() types.TxData {
	return &types.LegacyTx{
		Nonce:    t.Nonce,
		GasPrice: t.GasPrice,
		Gas:      t.Gas,
		To:       t.To,
		Value:    t.Value,
		Data:     t.Data,
	}
}

// UnsignedAccessListTx is the unsigned data of EIP-2930 access list transactions.
type UnsignedAccessListTx struct {
	ChainID    *big.Int         // destination chain ID
	Nonce      uint64           // nonce of sender account
	GasPrice   *big.Int         // wei per gas
	Gas        uint64           // gas limit
	To         *common.Address  `rlp:"nil"` // nil means contract creation
	Value      *big.Int         // wei amount
	Data       []byte           // contract invocation input data
	AccessList types.AccessList // EIP-2930 access list
}

func (t *UnsignedAccessListTx) TxData() types.TxData {
	return &types.AccessListTx{
		ChainID:    t.ChainID,
		Nonce:      t.Nonce,
		GasPrice:   t.GasPrice,
		Gas:        t.Gas,
		To:         t.To,
		Value:      t.Value,
		Data:       t.Data,
		AccessList: t.AccessList,
	}
}

// UnsignedDynamicFeeTx represents an unsigned EIP-1559 transaction.
type UnsignedDynamicFeeTx struct {
	ChainID    *big.Int
	Nonce      uint64
	GasTipCap  *big.Int // a.k.a. maxPriorityFeePerGas
	GasFeeCap  *big.Int // a.k.a. maxFeePerGas
	Gas        uint64
	To         *common.Address `rlp:"nil"` // nil means contract creation
	Value      *big.Int
	Data       []byte
	AccessList types.AccessList
}

func (t *UnsignedDynamicFeeTx) TxData() types.TxData {
	return &types.DynamicFeeTx{
		ChainID:    t.ChainID,
		Nonce:      t.Nonce,
		GasTipCap:  t.GasTipCap,
		GasFeeCap:  t.GasFeeCap,
		Gas:        t.Gas,
		To:         t.To,
		Value:      t.Value,
		Data:       t.Data,
		AccessList: t.AccessList,
	}
}

// UnsignedBlobTx represents an unsigned EIP-4844 transaction.
type UnsignedBlobTx struct {
	ChainID    *uint256.Int
	Nonce      uint64
	GasTipCap  *uint256.Int // a.k.a. maxPriorityFeePerGas
	GasFeeCap  *uint256.Int // a.k.a. maxFeePerGas
	Gas        uint64
	To         common.Address
	Value      *uint256.Int
	Data       []byte
	AccessList types.AccessList
	BlobFeeCap *uint256.Int // a.k.a. maxFeePerBlobGas
	BlobHashes []common.Hash
}

func (t *UnsignedBlobTx) TxData() types.TxData {
	return &types.BlobTx{
		ChainID:    t.ChainID,
		Nonce:      t.Nonce,
		GasTipCap:  t.GasTipCap,
		GasFeeCap:  t.GasFeeCap,
		Gas:        t.Gas,
		To:         t.To,
		Value:      t.Value,
		Data:       t.Data,
		AccessList: t.AccessList,
		BlobFeeCap: t.BlobFeeCap,
		BlobHashes: t.BlobHashes,
	}
}

// DecodeUnsignedTx decodes unsigned tx from the sign data, since the tx we received here is missing the r,s,v,
// so when using the tx.UnmarshalBinary, it is most likely to trigger the "too few elements" error in RLP
func DecodeUnsignedTx(b []byte) (*types.Transaction, error) {
	var tx *types.Transaction
	if len(b) > 0 && b[0] > 0x7f {
		// It's a legacy transaction.
		var data UnsignedLegacyTx
		err := rlp.DecodeBytes(b, &data)
		if err != nil {
			return nil, err
		}
		tx = types.NewTx(data.TxData())
		return tx, nil
	}
	// It's an EIP-2718 typed transaction envelope.
	if len(b) <= 1 {
		return nil, errors.New("typed transaction too short")
	}
	var inner UnsignedTxData
	switch b[0] {
	case types.AccessListTxType:
		inner = new(UnsignedAccessListTx)
	case types.DynamicFeeTxType:
		inner = new(UnsignedDynamicFeeTx)
	case types.BlobTxType:
		inner = new(UnsignedBlobTx)
	default:
		return nil, types.ErrTxTypeNotSupported
	}
	err := rlp.DecodeBytes(b[1:], inner)
	if err != nil {
		return nil, err
	}
	return types.NewTx(inner.TxData()), nil
}
