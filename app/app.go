package app

import (
	"fmt"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	jsoniter "github.com/json-iterator/go"
	"github.com/metis-seq/themis/metis"
	"github.com/metis-seq/themis/mpc"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/metis-seq/themis/auth"
	authTypes "github.com/metis-seq/themis/auth/types"
	"github.com/metis-seq/themis/bank"
	bankTypes "github.com/metis-seq/themis/bank/types"
	"github.com/metis-seq/themis/chainmanager"
	chainmanagerTypes "github.com/metis-seq/themis/chainmanager/types"
	"github.com/metis-seq/themis/common"
	gov "github.com/metis-seq/themis/gov"
	govTypes "github.com/metis-seq/themis/gov/types"
	"github.com/metis-seq/themis/helper"
	metisTypes "github.com/metis-seq/themis/metis/types"
	mpcTypes "github.com/metis-seq/themis/mpc/types"
	"github.com/metis-seq/themis/params"
	paramsClient "github.com/metis-seq/themis/params/client"
	"github.com/metis-seq/themis/params/subspace"
	paramsTypes "github.com/metis-seq/themis/params/types"
	"github.com/metis-seq/themis/sidechannel"
	sidechannelTypes "github.com/metis-seq/themis/sidechannel/types"

	// slashingTypes "github.com/metis-seq/themis/slashing/types"
	"github.com/metis-seq/themis/staking"
	stakingTypes "github.com/metis-seq/themis/staking/types"
	"github.com/metis-seq/themis/supply"
	supplyTypes "github.com/metis-seq/themis/supply/types"

	"github.com/metis-seq/themis/types"
	hmModule "github.com/metis-seq/themis/types/module"
	"github.com/metis-seq/themis/version"
)

const (
	// AppName denotes app name
	AppName = "Themis"
)

// Assertion for Themis app
var _ App = &ThemisApp{}

var (
	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		params.AppModuleBasic{},
		sidechannel.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		supply.AppModuleBasic{},
		chainmanager.AppModuleBasic{},
		staking.AppModuleBasic{},
		metis.AppModuleBasic{},
		mpc.AppModuleBasic{},
		gov.NewAppModuleBasic(paramsClient.ProposalHandler),
	)

	// module account permissions
	maccPerms = map[string][]string{
		authTypes.FeeCollectorName: nil,
		govTypes.ModuleName:        {},
	}
)

// ThemisApp main themis app
type ThemisApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tkeys map[string]*sdk.TransientStoreKey

	// subspaces
	subspaces map[string]subspace.Subspace

	// side router
	sideRouter types.SideRouter

	// keepers
	SidechannelKeeper sidechannel.Keeper
	AccountKeeper     auth.AccountKeeper
	BankKeeper        bank.Keeper
	SupplyKeeper      supply.Keeper
	GovKeeper         gov.Keeper
	ChainKeeper       chainmanager.Keeper
	StakingKeeper     staking.Keeper
	MetisKeeper       metis.Keeper
	MpcKeeper         mpc.Keeper

	// param keeper
	ParamsKeeper params.Keeper

	// contract keeper
	caller helper.ContractCaller

	//  total coins supply
	TotalCoinsSupply sdk.Coins

	// the module manager
	mm *module.Manager

	// simulation module manager
	sm *hmModule.SimulationManager
}

var logger = helper.Logger.With("module", "app")

//
// Module communicator
//

// ModuleCommunicator retriever
type ModuleCommunicator struct {
	App *ThemisApp
}

// GetACKCount returns ack count
func (d ModuleCommunicator) GetL1Batch(ctx sdk.Context) uint64 {
	return d.App.StakingKeeper.GetL1Batch(ctx)
}

// IsCurrentValidatorByAddress check if validator is current validator
func (d ModuleCommunicator) IsCurrentValidatorByAddress(ctx sdk.Context, address []byte) bool {
	return d.App.StakingKeeper.IsCurrentValidatorByAddress(ctx, address)
}

// GetAllDividendAccounts fetches all dividend accounts from topup module
func (d ModuleCommunicator) GetAllDividendAccounts(ctx sdk.Context) []types.DividendAccount {
	return nil
}

// GetValidatorFromValID get validator from validator id
func (d ModuleCommunicator) GetValidatorFromValID(ctx sdk.Context, valID types.ValidatorID) (validator types.Validator, ok bool) {
	return d.App.StakingKeeper.GetValidatorFromValID(ctx, valID)
}

// SetCoins sets coins
func (d ModuleCommunicator) SetCoins(ctx sdk.Context, addr types.ThemisAddress, amt sdk.Coins) sdk.Error {
	return d.App.BankKeeper.SetCoins(ctx, addr, amt)
}

// GetCoins gets coins
func (d ModuleCommunicator) GetCoins(ctx sdk.Context, addr types.ThemisAddress) sdk.Coins {
	return d.App.BankKeeper.GetCoins(ctx, addr)
}

// SendCoins transfers coins
func (d ModuleCommunicator) SendCoins(ctx sdk.Context, fromAddr types.ThemisAddress, toAddr types.ThemisAddress, amt sdk.Coins) sdk.Error {
	return d.App.BankKeeper.SendCoins(ctx, fromAddr, toAddr, amt)
}

// GetSpanById get span info by id
func (d ModuleCommunicator) GetSpanById(ctx sdk.Context, id uint64) (*types.Span, error) {
	return d.App.MetisKeeper.GetSpan(ctx, id)
}

func (d ModuleCommunicator) GetChainParams(ctx sdk.Context) chainmanagerTypes.Params {
	return d.App.ChainKeeper.GetParams(ctx)
}

//
// Themis app
//

// NewThemisApp creates themis app
func NewThemisApp(logger log.Logger, db dbm.DB, baseAppOptions ...func(*bam.BaseApp)) *ThemisApp {
	// create and register app-level codec for TXs and accounts
	cdc := MakeCodec()

	// set prefix
	config := sdk.GetConfig()
	config.Seal()

	// base app
	bApp := bam.NewBaseApp(AppName, logger, db, authTypes.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(nil)
	bApp.SetAppVersion(version.Version)

	// keys
	keys := sdk.NewKVStoreKeys(
		bam.MainStoreKey,
		sidechannelTypes.StoreKey,
		authTypes.StoreKey,
		bankTypes.StoreKey,
		supplyTypes.StoreKey,
		govTypes.StoreKey,
		chainmanagerTypes.StoreKey,
		stakingTypes.StoreKey,
		metisTypes.StoreKey,
		mpcTypes.StoreKey,
		paramsTypes.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(paramsTypes.TStoreKey)

	// create themis app
	var app = &ThemisApp{
		cdc:       cdc,
		BaseApp:   bApp,
		keys:      keys,
		tkeys:     tkeys,
		subspaces: make(map[string]subspace.Subspace),
	}

	// init params keeper and subspaces
	app.ParamsKeeper = params.NewKeeper(app.cdc, keys[paramsTypes.StoreKey], tkeys[paramsTypes.TStoreKey], paramsTypes.DefaultCodespace)
	app.subspaces[sidechannelTypes.ModuleName] = app.ParamsKeeper.Subspace(sidechannelTypes.DefaultParamspace)
	app.subspaces[authTypes.ModuleName] = app.ParamsKeeper.Subspace(authTypes.DefaultParamspace)
	app.subspaces[bankTypes.ModuleName] = app.ParamsKeeper.Subspace(bankTypes.DefaultParamspace)
	app.subspaces[supplyTypes.ModuleName] = app.ParamsKeeper.Subspace(supplyTypes.DefaultParamspace)
	app.subspaces[govTypes.ModuleName] = app.ParamsKeeper.Subspace(govTypes.DefaultParamspace).WithKeyTable(govTypes.ParamKeyTable())
	app.subspaces[chainmanagerTypes.ModuleName] = app.ParamsKeeper.Subspace(chainmanagerTypes.DefaultParamspace)
	app.subspaces[stakingTypes.ModuleName] = app.ParamsKeeper.Subspace(stakingTypes.DefaultParamspace)
	app.subspaces[metisTypes.ModuleName] = app.ParamsKeeper.Subspace(metisTypes.DefaultParamspace)
	app.subspaces[mpcTypes.ModuleName] = app.ParamsKeeper.Subspace(mpcTypes.DefaultParamspace)

	//
	// Contract caller
	//

	contractCallerObj, err := helper.NewContractCaller()
	if err != nil {
		cmn.Exit(err.Error())
	}

	app.caller = contractCallerObj

	//
	// module communicator
	//

	moduleCommunicator := ModuleCommunicator{App: app}

	//
	// keepers
	//

	// create side channel keeper
	app.SidechannelKeeper = sidechannel.NewKeeper(
		app.cdc,
		keys[sidechannelTypes.StoreKey], // target store
		app.subspaces[sidechannelTypes.ModuleName],
		common.DefaultCodespace,
	)

	// create chain keeper
	app.ChainKeeper = chainmanager.NewKeeper(
		app.cdc,
		keys[chainmanagerTypes.StoreKey], // target store
		app.subspaces[chainmanagerTypes.ModuleName],
		common.DefaultCodespace,
		app.caller,
	)

	// account keeper
	app.AccountKeeper = auth.NewAccountKeeper(
		app.cdc,
		keys[authTypes.StoreKey], // target store
		app.subspaces[authTypes.ModuleName],
		authTypes.ProtoBaseAccount, // prototype
	)

	app.StakingKeeper = staking.NewKeeper(
		app.cdc,
		keys[stakingTypes.StoreKey], // target store
		app.subspaces[stakingTypes.ModuleName],
		common.DefaultCodespace,
		app.ChainKeeper,
		moduleCommunicator,
	)

	// bank keeper
	app.BankKeeper = bank.NewKeeper(
		app.cdc,
		keys[bankTypes.StoreKey], // target store
		app.subspaces[bankTypes.ModuleName],
		bankTypes.DefaultCodespace,
		app.AccountKeeper,
		moduleCommunicator,
	)

	// supply keeper
	app.SupplyKeeper = supply.NewKeeper(
		app.cdc,
		keys[supplyTypes.StoreKey], // target store
		app.subspaces[supplyTypes.ModuleName],
		maccPerms,
		app.AccountKeeper,
		app.BankKeeper,
	)

	// register the proposal types
	govRouter := gov.NewRouter()
	govRouter.
		AddRoute(govTypes.RouterKey, govTypes.ProposalHandler).
		AddRoute(paramsTypes.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper))

	app.GovKeeper = gov.NewKeeper(
		app.cdc,
		keys[govTypes.StoreKey],
		app.subspaces[govTypes.ModuleName],
		app.SupplyKeeper,
		app.StakingKeeper,
		govTypes.DefaultCodespace,
		govRouter,
	)

	app.MetisKeeper = metis.NewKeeper(
		app.cdc,
		keys[metisTypes.StoreKey], // target store
		app.subspaces[metisTypes.ModuleName],
		common.DefaultCodespace,
		app.ChainKeeper,
		app.StakingKeeper,
		app.caller,
	)

	app.MpcKeeper = mpc.NewKeeper(
		app.cdc,
		keys[mpcTypes.StoreKey], // target store
		app.subspaces[mpcTypes.ModuleName],
		common.DefaultCodespace,
		app.caller,
		app.StakingKeeper,
		moduleCommunicator,
	)

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		sidechannel.NewAppModule(app.SidechannelKeeper),
		auth.NewAppModule(app.AccountKeeper, &app.caller, []authTypes.AccountProcessor{
			supplyTypes.AccountProcessor,
		}),
		bank.NewAppModule(app.BankKeeper, &app.caller),
		supply.NewAppModule(app.SupplyKeeper, &app.caller),
		gov.NewAppModule(app.GovKeeper, app.SupplyKeeper),
		chainmanager.NewAppModule(app.ChainKeeper, &app.caller),
		staking.NewAppModule(app.StakingKeeper, &app.caller),
		metis.NewAppModule(app.MetisKeeper, app.StakingKeeper, &app.caller),
		mpc.NewAppModule(app.MpcKeeper),
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		sidechannelTypes.ModuleName,
		authTypes.ModuleName,
		bankTypes.ModuleName,
		govTypes.ModuleName,
		chainmanagerTypes.ModuleName,
		supplyTypes.ModuleName,
		stakingTypes.ModuleName,
		metisTypes.ModuleName,
		mpcTypes.ModuleName,
	)

	// register message routes and query routes
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// side router
	app.sideRouter = types.NewSideRouter()
	for _, m := range app.mm.Modules {
		if m.Route() != "" {
			if sm, ok := m.(hmModule.SideModule); ok {
				app.sideRouter.AddRoute(m.Route(), &types.SideHandlers{
					SideTxHandler: sm.NewSideTxHandler(),
					PostTxHandler: sm.NewPostTxHandler(),
				})
			}
		}
	}

	app.sideRouter.Seal()

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	app.sm = hmModule.NewSimulationManager(
		auth.NewAppModule(app.AccountKeeper, &app.caller, []authTypes.AccountProcessor{
			supplyTypes.AccountProcessor,
		}),

		chainmanager.NewAppModule(app.ChainKeeper, &app.caller),
		staking.NewAppModule(app.StakingKeeper, &app.caller),
		bank.NewAppModule(app.BankKeeper, &app.caller),
	)
	app.sm.RegisterStoreDecoders()

	// mount the multistore and load the latest state
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)

	// perform initialization logic
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	app.SetAnteHandler(
		auth.NewAnteHandler(
			app.AccountKeeper,
			app.StakingKeeper,
			app.ChainKeeper,
			app.SupplyKeeper,
			&app.caller,
			auth.DefaultSigVerificationGasConsumer,
		),
	)
	// side-tx processor
	app.SetPostDeliverTxHandler(app.PostDeliverTxHandler)
	app.SetBeginSideBlocker(app.BeginSideBlocker)
	app.SetDeliverSideTxHandler(app.DeliverSideTxHandler)

	// load latest version
	err = app.LoadLatestVersion(app.keys[bam.MainStoreKey])
	if err != nil {
		cmn.Exit(err.Error())
	}

	app.Seal()

	return app
}

// MakeCodec create codec
func MakeCodec() *codec.Codec {
	cdc := codec.New()

	codec.RegisterCrypto(cdc)
	sdk.RegisterCodec(cdc)
	ModuleBasics.RegisterCodec(cdc)

	cdc.Seal()

	return cdc
}

// Name returns the name of the App
func (app *ThemisApp) Name() string { return app.BaseApp.Name() }

// InitChainer initializes chain
func (app *ThemisApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState

	if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	// get validator updates
	if err := ModuleBasics.ValidateGenesis(genesisState); err != nil {
		panic(err)
	}

	// check fee collector module account
	if moduleAcc := app.SupplyKeeper.GetModuleAccount(ctx, authTypes.FeeCollectorName); moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", authTypes.FeeCollectorName))
	}

	// init genesis
	app.mm.InitGenesis(ctx, genesisState)

	stakingState := stakingTypes.GetGenesisStateFromAppState(genesisState)

	// check if validator is current validator
	// add to val updates else skip
	var valUpdates []abci.ValidatorUpdate

	for _, validator := range stakingState.Validators {
		if validator.IsCurrentValidator(stakingState.StakingBatch) {
			// convert to Validator Update
			updateVal := abci.ValidatorUpdate{
				Power:  validator.VotingPower,
				PubKey: validator.PubKey.ABCIPubKey(),
			}
			// Add validator to validator updated to be processed below
			valUpdates = append(valUpdates, updateVal)
		}
	}

	// TODO make sure old validtors dont go in validator updates ie deactivated validators have to be removed
	// udpate validators
	return abci.ResponseInitChain{
		// validator updates
		Validators: valUpdates,
	}
}

// BeginBlocker application updates every begin block
func (app *ThemisApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	app.AccountKeeper.SetBlockProposer(
		ctx,
		types.BytesToThemisAddress(req.Header.GetProposerAddress()),
	)

	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker executes on each end block
func (app *ThemisApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	// transfer fees to current proposer
	if proposer, ok := app.AccountKeeper.GetBlockProposer(ctx); ok {
		moduleAccount := app.SupplyKeeper.GetModuleAccount(ctx, authTypes.FeeCollectorName)

		amount := moduleAccount.GetCoins().AmountOf(authTypes.FeeToken)
		if !amount.IsZero() {
			coins := sdk.Coins{sdk.Coin{Denom: authTypes.FeeToken, Amount: amount}}
			if err := app.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, authTypes.FeeCollectorName, proposer, coins); err != nil {
				logger.Error("EndBlocker | SendCoinsFromModuleToAccount", "Error", err)
			}
		}

		// remove block proposer
		app.AccountKeeper.RemoveBlockProposer(ctx)
	}

	var tmValUpdates []abci.ValidatorUpdate

	// --- Start update to new validators
	currentValidatorSet := app.StakingKeeper.GetValidatorSet(ctx)
	allValidators := app.StakingKeeper.GetAllValidators(ctx)
	ackCount := app.StakingKeeper.GetL1Batch(ctx)

	// get validator updates
	setUpdates := helper.GetUpdatedValidators(
		&currentValidatorSet, // pointer to current validator set -- UpdateValidators will modify it
		allValidators,        // All validators
		ackCount,             // ack count, L1 batch
	)

	if len(setUpdates) > 0 {
		// create new validator set
		if err := currentValidatorSet.UpdateWithChangeSet(setUpdates); err != nil {
			// return with nothing
			logger.Error("Unable to update current validator set", "Error", err)
			return abci.ResponseEndBlock{}
		}

		// increment proposer priority
		currentValidatorSet.IncrementProposerPriority(1)

		// validator set change
		logger.Debug("[ENDBLOCK] Updated current validator set", "proposer", currentValidatorSet.GetProposer())

		// save set in store
		if err := app.StakingKeeper.UpdateValidatorSetInStore(ctx, currentValidatorSet); err != nil {
			// return with nothing
			logger.Error("Unable to update current validator set in state", "Error", err)
			return abci.ResponseEndBlock{}
		}

		// convert updates from map to array
		for _, v := range setUpdates {
			tmValUpdates = append(tmValUpdates, abci.ValidatorUpdate{
				Power:  v.VotingPower,
				PubKey: v.PubKey.ABCIPubKey(),
			})
		}
	}

	// end block
	app.mm.EndBlock(ctx, req)

	// send validator updates to peppermint
	return abci.ResponseEndBlock{
		ValidatorUpdates: tmValUpdates,
	}
}

// LoadHeight loads a particular height
func (app *ThemisApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *ThemisApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supplyTypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// Codec returns ThemisApp's codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *ThemisApp) Codec() *codec.Codec {
	return app.cdc
}

// SetCodec set codec to app
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *ThemisApp) SetCodec(cdc *codec.Codec) {
	app.cdc = cdc
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *ThemisApp) GetKey(storeKey string) *sdk.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *ThemisApp) GetTKey(storeKey string) *sdk.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *ThemisApp) GetSubspace(moduleName string) subspace.Subspace {
	return app.subspaces[moduleName]
}

// GetSideRouter returns side-tx router
func (app *ThemisApp) GetSideRouter() types.SideRouter {
	return app.sideRouter
}

// SetSideRouter sets side-tx router
// Testing ONLY
func (app *ThemisApp) SetSideRouter(r types.SideRouter) {
	app.sideRouter = r
}

// GetModuleManager returns module manager
//
// NOTE: This is solely to be used for testing purposes.
func (app *ThemisApp) GetModuleManager() *module.Manager {
	return app.mm
}

// SimulationManager implements the SimulationApp interface
func (app *ThemisApp) SimulationManager() *hmModule.SimulationManager {
	return app.sm
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}

	return dupMaccPerms
}
