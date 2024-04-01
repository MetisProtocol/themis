package sidechannel_test

import (
	"math/rand"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/metis-seq/themis/app"
	"github.com/metis-seq/themis/sidechannel"
	"github.com/metis-seq/themis/sidechannel/simulation"
	"github.com/metis-seq/themis/sidechannel/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// GenesisTestSuite integrate test suite context object
type GenesisTestSuite struct {
	suite.Suite

	app *app.ThemisApp
	ctx sdk.Context
}

// SetupTest setup necessary things for genesis test
func (suite *GenesisTestSuite) SetupTest() {
	suite.app = setupWithGenesis()
	suite.ctx = suite.app.BaseApp.NewContext(true, abci.Header{})
}

// TestGenesisTestSuite
func TestGenesisTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(GenesisTestSuite))
}

// TestInitExportGenesis test import and export genesis state
func (suite *GenesisTestSuite) TestInitExportGenesis() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	genesisState := types.DefaultGenesisState()
	sidechannel.InitGenesis(ctx, app.SidechannelKeeper, genesisState)

	actualParams := sidechannel.ExportGenesis(ctx, app.SidechannelKeeper)
	require.Equal(t, genesisState, actualParams, "Default export should be default genesis state")

	// get random seed from time as source
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	genesisState = types.NewGenesisState(simulation.RandomPastCommits(r, 2, 5, 10))
	sidechannel.InitGenesis(ctx, app.SidechannelKeeper, genesisState)

	actualParams = sidechannel.ExportGenesis(ctx, app.SidechannelKeeper)

	require.Equal(t, len(genesisState.PastCommits), len(actualParams.PastCommits))
	require.Equal(t, len(genesisState.PastCommits[0].Txs), len(actualParams.PastCommits[0].Txs))
}
