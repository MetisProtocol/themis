package bank_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/metis-seq/themis/app"
	authTypes "github.com/metis-seq/themis/auth/types"
	"github.com/metis-seq/themis/bank"
	"github.com/metis-seq/themis/bank/types"
	"github.com/metis-seq/themis/helper/mocks"
	hmTypes "github.com/metis-seq/themis/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite

	app            *app.ThemisApp
	ctx            sdk.Context
	handler        sdk.Handler
	contractCaller mocks.IContractCaller
}

// SetupTest setup all necessary things for querier tesing
func (suite *HandlerTestSuite) SetupTest() {
	suite.app, suite.ctx = createTestApp(false)

	suite.contractCaller = mocks.IContractCaller{}
	suite.handler = bank.NewHandler(suite.app.BankKeeper, &suite.contractCaller)
}

// TestHandlerTestSuite
func TestHandlerTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestHandleMsgUnknown() {
	t, _, ctx := suite.T(), suite.app, suite.ctx

	result := suite.handler(ctx, nil)
	require.False(t, result.IsOK())
}

func (suite *HandlerTestSuite) TestHandlerMsgSend() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	amount := int64(10000000)
	from := hmTypes.HexToThemisAddress("123")
	to := hmTypes.HexToThemisAddress("456")
	_, err := app.BankKeeper.AddCoins(ctx, from, sdk.NewCoins(sdk.NewCoin(authTypes.FeeToken, sdk.NewInt(amount*10))))
	require.NoError(t, err)

	msgSend := types.NewMsgSend(
		from,
		to,
		sdk.NewCoins(sdk.NewCoin(authTypes.FeeToken, sdk.NewInt(amount))),
	)
	result := suite.handler(ctx, msgSend)
	require.True(t, result.IsOK(), "Expected New msg to be sent")

	fromAcc := app.BankKeeper.GetCoins(ctx, to)
	require.Less(t, fromAcc.AmountOf(authTypes.FeeToken).Int64(), sdk.NewInt(amount*10).Int64())

	toAcc := app.BankKeeper.GetCoins(ctx, to)
	require.Equal(t, sdk.NewInt(amount), toAcc.AmountOf(authTypes.FeeToken))
}
