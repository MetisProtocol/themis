package simulation

import (
	"github.com/metis-seq/themis/bank/types"
	"github.com/metis-seq/themis/types/module"
)

// RandomizedGenState returns bank genesis
func RandomizedGenState(simState *module.SimulationState) {
	bankGenesis := types.NewGenesisState(true)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(bankGenesis)
}
