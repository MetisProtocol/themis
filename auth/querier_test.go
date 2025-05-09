package auth_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth/types"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/metis-seq/themis/app"
	"github.com/metis-seq/themis/auth"
	"github.com/metis-seq/themis/auth/exported"
	"github.com/metis-seq/themis/auth/types"
	authTypes "github.com/metis-seq/themis/auth/types"
	hmTypes "github.com/metis-seq/themis/types"
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
	suite.querier = auth.NewQuerier(suite.app.AccountKeeper)
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

func (suite *QuerierTestSuite) TestQueryAccount() {
	t, happ, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier
	cdc := happ.Codec()

	// account path
	path := []string{types.QueryAccount}

	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAccount),
		Data: []byte{},
	}
	res, err := querier(ctx, path, req)
	require.Error(t, err)
	require.Nil(t, res)

	req.Data = cdc.MustMarshalJSON(types.NewQueryAccountParams(hmTypes.BytesToThemisAddress([]byte(""))))
	res, err = querier(ctx, path, req)
	require.Error(t, err)
	require.Nil(t, res)

	_, _, addr := sdkAuth.KeyTestPubAddr()
	req.Data = cdc.MustMarshalJSON(types.NewQueryAccountParams(hmTypes.AccAddressToThemisAddress(addr)))
	res, err = querier(ctx, path, req)
	require.Error(t, err)
	require.Nil(t, res)

	happ.AccountKeeper.SetAccount(ctx, happ.AccountKeeper.NewAccountWithAddress(ctx, hmTypes.AccAddressToThemisAddress(addr)))
	res, err = querier(ctx, path, req)
	require.NoError(t, err)
	require.NotNil(t, res)

	res, err = querier(ctx, path, req)
	require.NoError(t, err)
	require.NotNil(t, res)

	var account exported.Account
	err2 := cdc.UnmarshalJSON(res, &account)
	require.Nil(t, err2)
	require.Equal(t, account.GetAddress().Bytes(), addr.Bytes())

	{
		// setting tnil to account
		require.Panics(t, func() {
			happ.AccountKeeper.SetAccount(ctx, nil)
		})

		// store invalid/empty account
		store := ctx.KVStore(happ.GetKey(authTypes.StoreKey))
		store.Set(types.AddressStoreKey(hmTypes.AccAddressToThemisAddress(addr)), []byte(""))
		require.Panics(t, func() {
			_, err = querier(ctx, path, req)
			require.NoError(t, err)
		})
	}
}

func (suite *QuerierTestSuite) TestQueryParams() {
	t, happ, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	path := []string{types.QueryParams}
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParams),
		Data: []byte{},
	}
	res, err := querier(ctx, path, req)
	require.NoError(t, err)
	require.NotNil(t, res)

	// default params
	defaultParams := authTypes.DefaultParams()

	var params types.Params

	err2 := jsoniter.ConfigFastest.Unmarshal(res, &params)
	require.Nil(t, err2)
	require.Equal(t, defaultParams.MaxMemoCharacters, params.MaxMemoCharacters)
	require.Equal(t, defaultParams.TxSigLimit, params.TxSigLimit)
	require.Equal(t, defaultParams.TxSizeCostPerByte, params.TxSizeCostPerByte)
	require.Equal(t, defaultParams.SigVerifyCostED25519, params.SigVerifyCostED25519)
	require.Equal(t, defaultParams.SigVerifyCostSecp256k1, params.SigVerifyCostSecp256k1)
	require.Equal(t, defaultParams.MaxTxGas, params.MaxTxGas)
	require.Equal(t, defaultParams.TxFees, params.TxFees)

	// set max characters
	params.MaxMemoCharacters = 10
	params.TxSizeCostPerByte = 8
	happ.AccountKeeper.SetParams(ctx, params)
	res, err = querier(ctx, path, req)
	require.NoError(t, err)
	require.NotEmpty(t, string(res))

	var params3 types.Params
	err3 := jsoniter.ConfigFastest.Unmarshal(res, &params3)
	require.NoError(t, err3)
	require.Equal(t, uint64(10), params.MaxMemoCharacters)
	require.Equal(t, uint64(8), params.TxSizeCostPerByte)

	{
		happ := app.Setup(true)
		ctx := happ.BaseApp.NewContext(true, abci.Header{})
		querier := auth.NewQuerier(happ.AccountKeeper)
		require.Panics(t, func() {
			_, err = querier(ctx, path, req)
			require.NoError(t, err)
		})
	}
}
