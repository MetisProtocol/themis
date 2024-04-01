package app

import (
	"github.com/cosmos/cosmos-sdk/codec"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	authTypes "github.com/metis-seq/themis/auth/types"

	hmTypes "github.com/metis-seq/themis/types"
)

// Setup initializes a new App. A Nop logger is set in App.
func Setup(isCheckTx bool) *ThemisApp {
	db := dbm.NewMemDB()
	app := NewThemisApp(log.NewNopLogger(), db)

	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		genesisState := NewDefaultGenesisState()

		stateBytes, err := codec.MarshalJSONIndent(app.Codec(), genesisState)
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators:    []abci.ValidatorUpdate{},
				AppStateBytes: stateBytes,
			},
		)
	}

	return app
}

// SetupWithGenesisAccounts initializes a new Themis with the provided genesis
// accounts and possible balances.
func SetupWithGenesisAccounts(genAccs []authTypes.GenesisAccount) *ThemisApp {
	// setup with isCheckTx
	app := Setup(true)

	// initialize the chain with the passed in genesis accounts
	genesisState := NewDefaultGenesisState()

	authGenesis := authTypes.NewGenesisState(authTypes.DefaultParams(), genAccs)
	genesisState[authTypes.ModuleName] = app.Codec().MustMarshalJSON(authGenesis)

	// bankGenesis := authTypes.NewGenesisState(authTypes.DefaultGenesisState().SendEnabled)
	// genesisState[authTypes.ModuleName] = app.Codec().MustMarshalJSON(bankGenesis)

	stateBytes, err := codec.MarshalJSONIndent(app.Codec(), genesisState)
	if err != nil {
		panic(err)
	}

	app.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)

	app.Commit()
	app.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: app.LastBlockHeight() + 1}})

	return app
}

// GenerateAccountStrategy account strategy
type GenerateAccountStrategy func(int) []hmTypes.ThemisAddress
