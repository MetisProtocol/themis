package app

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	authTypes "github.com/metis-seq/themis/auth/types"
	govTypes "github.com/metis-seq/themis/gov/types"
	paramTypes "github.com/metis-seq/themis/params/types"
	"github.com/metis-seq/themis/simulation"
	stakingTypes "github.com/metis-seq/themis/staking/types"
	supplyTypes "github.com/metis-seq/themis/supply/types"
)

// Get flags every time the simulator is run
func init() {
	GetSimulatorFlags()
}

type StoreKeysPrefixes struct {
	A        sdk.StoreKey
	B        sdk.StoreKey
	Prefixes [][]byte
}

func TestFullAppSimulation(t *testing.T) {
	t.Parallel()

	config, db, dir, logger, skip, err := SetupSimulation("leveldb-app-sim", "Simulation")
	if skip {
		t.Skip("skipping application simulation")
	}

	require.NoError(t, err, "simulation setup failed")
	require.NotNil(t, db, "DB should not be nil")

	defer func() {
		db.Close()
		require.NoError(t, os.RemoveAll(dir))
	}()

	app := NewThemisApp(logger, db)
	require.Equal(t, AppName, app.Name())

	// run randomized simulation
	_, simParams, simErr := simulation.SimulateFromSeed(
		t, os.Stdout, app.BaseApp, AppStateFn(app.Codec(), app.SimulationManager()),
		SimulationOperations(app, app.Codec(), config),
		app.ModuleAccountAddrs(), config,
	)

	// export state and simParams before the simulation error is checked
	err = CheckExportSimulation(app, config, simParams)
	require.NoError(t, err)
	require.NoError(t, simErr)

	if config.Commit {
		PrintStats(db)
	}
}

func TestAppImportExport(t *testing.T) {
	t.Parallel()

	config, db, dir, logger, skip, err := SetupSimulation("leveldb-app-sim", "Simulation")
	if skip {
		t.Skip("skipping application import/export simulation")
	}

	require.NoError(t, err, "simulation setup failed")

	defer func() {
		db.Close()
		require.NoError(t, os.RemoveAll(dir))
	}()

	app := NewThemisApp(logger, db)
	require.Equal(t, AppName, app.Name())

	// Run randomized simulation
	_, simParams, simErr := simulation.SimulateFromSeed(
		t, os.Stdout, app.BaseApp, AppStateFn(app.Codec(), app.SimulationManager()),
		SimulationOperations(app, app.Codec(), config),
		app.ModuleAccountAddrs(), config,
	)

	// export state and simParams before the simulation error is checked
	err = CheckExportSimulation(app, config, simParams)
	require.NoError(t, err)
	require.NoError(t, simErr)

	if config.Commit {
		PrintStats(db)
	}

	fmt.Printf("exporting genesis...\n")

	appState, _, err := app.ExportAppStateAndValidators()
	require.NoError(t, err)

	fmt.Printf("importing genesis...\n")

	_, newDB, newDir, _, _, err := SetupSimulation("leveldb-app-sim-2", "Simulation-2")
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		newDB.Close()
		require.NoError(t, os.RemoveAll(newDir))
	}()

	newApp := NewThemisApp(logger, newDB)
	require.Equal(t, AppName, newApp.Name())

	var genesisState GenesisState
	err = app.Codec().UnmarshalJSON(appState, &genesisState)
	require.NoError(t, err)

	ctxA := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})
	ctxB := newApp.NewContext(true, abci.Header{Height: app.LastBlockHeight()})
	newApp.mm.InitGenesis(ctxB, genesisState)

	fmt.Printf("comparing stores...\n")

	storeKeysPrefixes := []StoreKeysPrefixes{
		{app.keys[baseapp.MainStoreKey], newApp.keys[baseapp.MainStoreKey], [][]byte{}},
		{app.keys[authTypes.StoreKey], newApp.keys[authTypes.StoreKey], [][]byte{}},
		{app.keys[stakingTypes.StoreKey], newApp.keys[stakingTypes.StoreKey], [][]byte{}},
		{app.keys[supplyTypes.StoreKey], newApp.keys[supplyTypes.StoreKey], [][]byte{}},
		{app.keys[paramTypes.StoreKey], newApp.keys[paramTypes.StoreKey], [][]byte{}},
		{app.keys[govTypes.StoreKey], newApp.keys[govTypes.StoreKey], [][]byte{}},
	}

	for _, skp := range storeKeysPrefixes {
		storeA := ctxA.KVStore(skp.A)
		storeB := ctxB.KVStore(skp.B)

		_, _, _, equal := sdk.DiffKVStores(storeA, storeB, skp.Prefixes)
		require.True(t, equal, "unequal sets of key-values to compare")
	}
}

func TestAppSimulationAfterImport(t *testing.T) {
	t.Parallel()

	config, db, dir, logger, skip, err := SetupSimulation("leveldb-app-sim", "Simulation")
	require.NoError(t, err, "simulation setup failed")

	if skip {
		t.Skip("skipping application simulation after import")
	}

	defer func() {
		db.Close()
		require.NoError(t, os.RemoveAll(dir))
	}()

	app := NewThemisApp(logger, db)
	require.Equal(t, AppName, app.Name())

	// Run randomized simulation
	stopEarly, simParams, simErr := simulation.SimulateFromSeed(
		t, os.Stdout, app.BaseApp, AppStateFn(app.Codec(), app.SimulationManager()),
		SimulationOperations(app, app.Codec(), config),
		app.ModuleAccountAddrs(), config,
	)

	// export state and simParams before the simulation error is checked
	err = CheckExportSimulation(app, config, simParams)
	require.NoError(t, err)
	require.NoError(t, simErr)

	if config.Commit {
		PrintStats(db)
	}

	if stopEarly {
		fmt.Println("can't export or import a zero-validator genesis, exiting test...")
		return
	}

	fmt.Printf("exporting genesis...\n")

	appState, _, err := app.ExportAppStateAndValidators()
	require.NoError(t, err)

	fmt.Printf("importing genesis...\n")

	_, newDB, newDir, _, _, err := SetupSimulation("leveldb-app-sim-2", "Simulation-2")
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		newDB.Close()
		require.NoError(t, os.RemoveAll(newDir))
	}()

	newApp := NewThemisApp(logger, newDB)

	require.Equal(t, AppName, app.Name())

	newApp.InitChain(abci.RequestInitChain{
		AppStateBytes: appState,
	})

	stopEarly, _, err = simulation.SimulateFromSeed(
		t, os.Stdout, newApp.BaseApp, AppStateFn(app.Codec(), app.SimulationManager()),
		SimulationOperations(newApp, newApp.Codec(), config),
		newApp.ModuleAccountAddrs(), config,
	)
	require.False(t, stopEarly)
	require.NoError(t, err)
}

// TODO: Make another test for the fuzzer itself, which just has noOp txs
// and doesn't depend on the application.
func TestAppStateDeterminism(t *testing.T) {
	t.Parallel()

	if !FlagEnabledValue {
		t.Skip("skipping application simulation")
	}

	config := NewConfigFromFlags()
	config.InitialBlockHeight = 1
	config.ExportParamsPath = ""
	config.OnOperation = false
	config.AllInvariants = false
	config.ChainID = SimAppChainID

	numSeeds := 3
	numTimesToRunPerSeed := 5
	appHashList := make([]json.RawMessage, numTimesToRunPerSeed)

	for i := 0; i < numSeeds; i++ {
		config.Seed = rand.Int63()

		for j := 0; j < numTimesToRunPerSeed; j++ {
			var logger log.Logger
			if FlagVerboseValue {
				logger = log.TestingLogger()
			} else {
				logger = log.NewNopLogger()
			}

			db := dbm.NewMemDB()

			// app := NewSimApp(logger, db, nil, true, map[int64]bool{}, DefaultNodeHome, FlagPeriodValue, interBlockCacheOpt())
			app := NewThemisApp(logger, db, nil)

			err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
			require.NoError(t, err)
			require.Equal(t, AppName, app.Name())

			fmt.Printf(
				"running non-determinism simulation; seed %d: %d/%d, attempt: %d/%d\n",
				config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
			)

			_, _, err = simulation.SimulateFromSeed(
				t, os.Stdout, app.BaseApp, AppStateFn(app.Codec(), app.SimulationManager()),
				SimulationOperations(app, app.Codec(), config),
				app.ModuleAccountAddrs(), config,
			)
			require.NoError(t, err)

			if config.Commit {
				PrintStats(db)
			}

			appHash := app.LastCommitID().Hash
			appHashList[j] = appHash

			if j != 0 {
				require.Equal(
					t, appHashList[0], appHashList[j],
					"non-determinism in seed %d: %d/%d, attempt: %d/%d\n", config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
				)
			}
		}
	}
}
