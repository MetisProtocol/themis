package bank_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/metis-seq/themis/app"
	authTypes "github.com/metis-seq/themis/auth/types"
	bankTypes "github.com/metis-seq/themis/bank/types"
)

//
// Create test app
//

// returns context and app with params set on account keeper
func createTestApp(isCheckTx bool) (*app.ThemisApp, sdk.Context) {
	app := app.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})
	app.AccountKeeper.SetParams(ctx, authTypes.DefaultParams())
	app.BankKeeper.SetSendEnabled(ctx, bankTypes.DefaultSendEnabled)

	return app, ctx
}
