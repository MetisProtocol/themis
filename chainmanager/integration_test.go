package chainmanager_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/metis-seq/themis/app"
	chainManagerTypes "github.com/metis-seq/themis/chainmanager/types"
)

//
// Create test app
//

// returns context and app with params set on chainmanager keeper
func createTestApp(isCheckTx bool) (*app.ThemisApp, sdk.Context) {
	app := app.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})
	app.ChainKeeper.SetParams(ctx, chainManagerTypes.DefaultParams())

	return app, ctx
}
