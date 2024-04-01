package auth_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/metis-seq/themis/app"
	authTypes "github.com/metis-seq/themis/auth/types"
	supplyTypes "github.com/metis-seq/themis/supply/types"
)

//
// Test suite
//

// ModuleTestSuite integrate test suite context object
type ModuleTestSuite struct {
	suite.Suite

	app *app.ThemisApp
	ctx sdk.Context
}

func (suite *ModuleTestSuite) SetupTest() {
	suite.app, suite.ctx = createTestApp(false)
}

func TestModuleTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ModuleTestSuite))
}

//
// Tests
//

func (suite *ModuleTestSuite) TestModuleAccount() {
	t, happ, ctx := suite.T(), suite.app, suite.ctx
	acc := happ.AccountKeeper.GetAccount(ctx, supplyTypes.NewModuleAddress(authTypes.FeeCollectorName))
	require.NotNil(t, acc)
}
