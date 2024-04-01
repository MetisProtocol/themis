package simulation

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/metis-seq/themis/chainmanager/types"
	hmTypes "github.com/metis-seq/themis/types"
	"github.com/metis-seq/themis/types/module"
	"github.com/metis-seq/themis/types/simulation"
)

// Parameter keys
const (
	MainchainTxConfirmations  = "mainchain_tx_confirmations"
	MetischainTxConfirmations = "metischain_tx_confirmations"

	MetisChainID          = "metis_chain_id"
	MetisTokenAddress     = "metis_token_address"     //nolint
	StakingManagerAddress = "staking_manager_address" //nolint
	StakingInfoAddress    = "staking_info_address"    //nolint

	// Metis Chain Contracts
	ValidatorSetAddress = "validator_set_address" //nolint
)

func GenMainchainTxConfirmations(r *rand.Rand) uint64 {
	return uint64(simulation.RandIntBetween(r, 1, 100))
}

func GenMetischainTxConfirmations(r *rand.Rand) uint64 {
	return uint64(simulation.RandIntBetween(r, 1, 100))
}

func GenThemisAddress() hmTypes.ThemisAddress {
	return hmTypes.BytesToThemisAddress(simulation.RandHex(20))
}

// GenMetisChainId returns randomc chainID
func GenMetisChainId(r *rand.Rand) string {
	return strconv.Itoa(simulation.RandIntBetween(r, 0, math.MaxInt32))
}

func RandomizedGenState(simState *module.SimulationState) {
	var mainchainTxConfirmations uint64

	simState.AppParams.GetOrGenerate(simState.Cdc, MainchainTxConfirmations, &mainchainTxConfirmations, simState.Rand,
		func(r *rand.Rand) { mainchainTxConfirmations = GenMainchainTxConfirmations(r) },
	)

	var (
		metischainTxConfirmations uint64
		metisChainID              string
	)

	simState.AppParams.GetOrGenerate(simState.Cdc, MetischainTxConfirmations, &metischainTxConfirmations, simState.Rand,
		func(r *rand.Rand) { metischainTxConfirmations = GenMetischainTxConfirmations(r) },
	)

	simState.AppParams.GetOrGenerate(simState.Cdc, MetisChainID, &metisChainID, simState.Rand,
		func(r *rand.Rand) { metisChainID = GenMetisChainId(r) },
	)

	var (
		metisTokenAddress     = GenThemisAddress()
		stakingManagerAddress = GenThemisAddress()
		stakingInfoAddress    = GenThemisAddress()
		validatorSetAddress   = GenThemisAddress()
	)

	chainParams := types.ChainParams{
		MetisChainID:          metisChainID,
		MetisTokenAddress:     metisTokenAddress,
		StakingManagerAddress: stakingManagerAddress,
		StakingInfoAddress:    stakingInfoAddress,
		ValidatorSetAddress:   validatorSetAddress,
	}
	params := types.NewParams(mainchainTxConfirmations, metischainTxConfirmations, chainParams)
	chainManagerGenesis := types.NewGenesisState(params)
	fmt.Printf("Selected randomly generated chainmanager parameters:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, chainManagerGenesis))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(chainManagerGenesis)
}
