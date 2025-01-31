package types

import (
	"encoding/json"

	"github.com/metis-seq/themis/gov/types"
)

//
// Gensis state
//

// GenesisState - all chainmanager state that must be provided at genesis
type GenesisState struct {
	Params Params `json:"params" yaml:"params"`
}

// NewGenesisState - Create a new genesis state
func NewGenesisState(params Params) GenesisState {
	return GenesisState{
		Params: params,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultParams())
}

// ValidateGenesis performs basic validation of auth genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	return nil
}

// GetGenesisStateFromAppState returns staking GenesisState given raw application genesis state
func GetGenesisStateFromAppState(appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		ModuleCdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return genesisState
}

// SetGenesisStateToAppState sets state into app state
func SetGenesisStateToAppState(appState map[string]json.RawMessage, params Params) (map[string]json.RawMessage, error) {
	genesisState := GetGenesisStateFromAppState(appState)
	genesisState.Params = params
	appState[ModuleName] = types.ModuleCdc.MustMarshalJSON(genesisState)

	return appState, nil
}
