package simulation

import (
	"fmt"
	"math/rand"

	"github.com/metis-seq/themis/chainmanager/types"
	"github.com/metis-seq/themis/simulation"
	simtypes "github.com/metis-seq/themis/types/simulation"
)

const (
	KeyMainchainTxConfirmations  = "MainchainTxConfirmations"
	KeyMetischainTxConfirmations = "MetischainTxConfirmations"
	KeyChainParams               = "ChainParams"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, KeyMainchainTxConfirmations,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenMainchainTxConfirmations(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, KeyMetischainTxConfirmations,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenMetischainTxConfirmations(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, KeyChainParams,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenMetisChainId(r))
			},
		),
	}
}
