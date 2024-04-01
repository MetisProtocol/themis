package bank_test

import (
	"encoding/json"
	"fmt"
	"testing"

	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/metis-seq/themis/app"
	authTypes "github.com/metis-seq/themis/auth/types"
	"github.com/metis-seq/themis/bank"
	"github.com/metis-seq/themis/bank/types"
	hmTypes "github.com/metis-seq/themis/types"
	"github.com/metis-seq/themis/types/simulation"
)

//
// Test suite
//

// QuerierTestSuite integrate test suite context object
type QuerierTestSuite struct {
	suite.Suite

	app     *app.ThemisApp
	ctx     sdk.Context
	querier sdk.Querier
}

func (suite *QuerierTestSuite) SetupTest() {
	suite.app, suite.ctx = createTestApp(false)
	suite.querier = bank.NewQuerier(suite.app.BankKeeper)
}

func TestQuerierTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(QuerierTestSuite))
}

//
// Tests
//

func (suite *QuerierTestSuite) TestInvalidQuery() {
	t, _, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	req := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	bz, err := querier(ctx, []string{"other"}, req)
	require.Error(t, err)
	require.Nil(t, bz)

	bz, err = querier(ctx, []string{types.QuerierRoute}, req)
	require.Error(t, err)
	require.Nil(t, bz)
}

func (suite *QuerierTestSuite) TestQueryBalance() {
	t, happ, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier
	cdc := happ.Codec()

	// account path
	path := []string{types.QueryBalance}

	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryBalance),
		Data: []byte{},
	}
	res, sdkErr := querier(ctx, path, req)
	require.Error(t, sdkErr)
	require.Nil(t, res)

	// balance for non-existing address (empty address)

	req.Data = cdc.MustMarshalJSON(types.NewQueryBalanceParams(hmTypes.BytesToThemisAddress([]byte(""))))
	res, sdkErr = querier(ctx, path, req)
	require.NoError(t, sdkErr)
	require.NotNil(t, res)

	// fetch balance
	var balance sdk.Coins

	require.NoError(t, json.Unmarshal(res, &balance))
	require.True(t, balance.IsZero())

	// balance for non-existing address
	_, _, addr := sdkAuth.KeyTestPubAddr()
	req.Data = cdc.MustMarshalJSON(types.NewQueryBalanceParams(hmTypes.AccAddressToThemisAddress(addr)))
	res, sdkErr = querier(ctx, path, req)
	require.NoError(t, sdkErr)
	require.NotNil(t, res)

	require.NoError(t, json.Unmarshal(res, &balance))
	require.True(t, balance.IsZero())

	// set account
	acc1 := happ.AccountKeeper.NewAccountWithAddress(ctx, hmTypes.AccAddressToThemisAddress(addr))
	amt := simulation.RandomFeeCoins()
	err := acc1.SetCoins(amt)
	require.NoError(t, err)
	happ.AccountKeeper.SetAccount(ctx, acc1)

	res, err = querier(ctx, path, req)
	require.NoError(t, err)
	require.NotNil(t, res)

	require.NoError(t, json.Unmarshal(res, &balance))
	require.True(t, balance.IsEqual(amt), "address coins stored in the store should be equal to amt")

	{
		// setting nil to account
		require.Panics(t, func() {
			happ.AccountKeeper.SetAccount(ctx, nil)
		})

		// store invalid/empty account
		store := ctx.KVStore(happ.GetKey(authTypes.StoreKey))
		store.Set(authTypes.AddressStoreKey(hmTypes.AccAddressToThemisAddress(addr)), []byte(""))
		require.Panics(t, func() {
			_, err = querier(ctx, path, req)
			require.NoError(t, err)
		})
	}
}
