package types

import (
	"encoding/json"

	chainmanagerTypes "github.com/metis-seq/themis/chainmanager/types"
	"github.com/metis-seq/themis/gov/types"
	hmTypes "github.com/metis-seq/themis/types"
)

// GenesisState is the metis state that must be provided at genesis.
type GenesisState struct {
	Params Params          `json:"params" yaml:"params"`
	Spans  []*hmTypes.Span `json:"spans" yaml:"spans"` // list of spans
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(params Params, spans []*hmTypes.Span) GenesisState {
	return GenesisState{
		Params: params,
		Spans:  spans,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultParams(), nil)
}

// ValidateGenesis performs basic validation of metis genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	return nil
}

// genFirstSpan generates default first valdiator producer set
func genFirstSpan(valset hmTypes.ValidatorSet, chainId string, starBlock, endBlock uint64) []*hmTypes.Span {
	var (
		firstSpan         []*hmTypes.Span
		selectedProducers []hmTypes.Validator
	)

	if len(valset.Validators) > int(DefaultProducerCount) {
		// pop top validators and select
		for i := 0; uint64(i) < DefaultProducerCount; i++ {
			selectedProducers = append(selectedProducers, *valset.Validators[i])
		}
	} else {
		for _, val := range valset.Validators {
			selectedProducers = append(selectedProducers, *val)
		}
	}

	newSpan := hmTypes.NewSpan(0, starBlock, endBlock, valset, selectedProducers, chainId)
	firstSpan = append(firstSpan, &newSpan)

	return firstSpan
}

// GetGenesisStateFromAppState returns staking GenesisState given raw application genesis state
func GetGenesisStateFromAppState(appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		types.ModuleCdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return genesisState
}

// SetGenesisStateToAppState sets state into app state
func SetGenesisStateToAppState(appState map[string]json.RawMessage, currentValSet hmTypes.ValidatorSet,
	startBlock, endBlock, epochLength uint64) (map[string]json.RawMessage, error) {
	// set state to metis state
	metisState := GetGenesisStateFromAppState(appState)
	chainState := chainmanagerTypes.GetGenesisStateFromAppState(appState)
	metisState.Params.SpanDuration = epochLength
	metisState.Spans = genFirstSpan(currentValSet,
		chainState.Params.ChainParams.MetisChainID, startBlock, endBlock)

	appState[ModuleName] = types.ModuleCdc.MustMarshalJSON(metisState)

	return appState, nil
}
