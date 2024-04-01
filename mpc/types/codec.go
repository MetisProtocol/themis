package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgProposeMpcCreate{}, "mpc/MsgProposeMpcCreate", nil)
	cdc.RegisterConcrete(MsgProposeMpcSign{}, "mpc/MsgProposeMpcSign", nil)
	cdc.RegisterConcrete(MsgMpcSign{}, "mpc/MsgMpcSign", nil)
}

// ModuleCdc generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	RegisterCodec(cdc)
	ModuleCdc = cdc.Seal()
}
