package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
	tmLog "github.com/tendermint/tendermint/libs/log"

	"github.com/metis-seq/themis/helper"
)

// RestLogger for mpc module logger
var RestLogger tmLog.Logger

func init() {
	RestLogger = helper.Logger.With("module", "mpc/rest")
}

// RegisterRoutes registers  mpc-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}
