package module

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/metis-seq/themis/types"
)

// ThemisModuleBasic is the standard form for basic non-dependant elements of an application module.
type ThemisModuleBasic interface {
	module.AppModuleBasic

	// verify genesis
	VerifyGenesis(map[string]json.RawMessage) error
}

// SideModule is the standard form for side tx elements of an application module
type SideModule interface {
	NewSideTxHandler() types.SideTxHandler
	NewPostTxHandler() types.PostTxHandler
}
