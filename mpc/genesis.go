package mpc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/metis-seq/themis/mpc/types"
	hmTypes "github.com/metis-seq/themis/types"
)

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	if len(data.MpcSet) > 0 {
		for _, mpcVal := range data.MpcSet {
			if err := keeper.AddNewMpcSet(ctx, mpcVal); err != nil {
				keeper.Logger(ctx).Error("Error AddNewMpcSet", "error", err)
			}
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	allMpcSets := keeper.GetAllMpcSets(ctx)

	if len(allMpcSets) == 0 {
		return types.NewGenesisState(nil)
	}

	var sets []hmTypes.PartyID
	for _, mpcSet := range allMpcSets {
		sets = append(sets, *mpcSet)
	}

	return types.NewGenesisState(
		// TODO think better way to export all mpcSet
		sets,
	)
}
