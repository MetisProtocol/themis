package mpc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	chainmanagerTypes "github.com/metis-seq/themis/chainmanager/types"
	"github.com/metis-seq/themis/types"
)

type ModuleCommunicator interface {
	GetSpanById(ctx sdk.Context, id uint64) (*types.Span, error)
	GetChainParams(ctx sdk.Context) chainmanagerTypes.Params
}
