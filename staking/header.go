package staking

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/metis-seq/themis/staking/simulation"
	"github.com/metis-seq/themis/types"
)

// LoadValidatorSet loads validator set
func LoadValidatorSet(t *testing.T, count int, keeper Keeper, ctx sdk.Context, randomise bool, timeAlive int) types.ValidatorSet {
	t.Helper()

	var valSet types.ValidatorSet
	validators := simulation.GenRandomVal(count, 0, 10, uint64(timeAlive), randomise, 1)
	for _, validator := range validators {
		err := keeper.AddValidator(ctx, validator)
		require.NoError(t, err, "Unable to set validator, Error: %v", err)

		err = valSet.UpdateWithChangeSet([]*types.Validator{&validator})
		require.NoError(t, err)
	}

	err := keeper.UpdateValidatorSetInStore(ctx, valSet)
	require.NoError(t, err, "Unable to update validator set")

	vals := keeper.GetAllValidators(ctx)
	require.NotNil(t, vals)

	return valSet
}
