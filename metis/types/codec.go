package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgProposeSpan{}, "metis/MsgProposeSpan", nil)
	cdc.RegisterConcrete(MsgReProposeSpan{}, "metis/MsgReProposeSpan", nil)
	cdc.RegisterConcrete(MsgMetisTx{}, "metis/MsgMetisTx", nil)
}

// ModuleCdc generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	RegisterCodec(cdc)
	ModuleCdc = cdc.Seal()
}
