package types

import (
	"encoding/json"

	"github.com/metis-seq/themis/gov/types"
	hmTypes "github.com/metis-seq/themis/types"
)

// GenesisState is the mpc state that must be provided at genesis.
type GenesisState struct {
	MpcSet []hmTypes.PartyID `json:"mpc_set" yaml:"mpc_set"` // list of mpc
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(mpcs []hmTypes.PartyID) GenesisState {
	return GenesisState{
		MpcSet: mpcs,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(nil)
}

// ValidateGenesis performs basic validation of metis genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	return nil
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
func SetGenesisStateToAppState(appState map[string]json.RawMessage, currentValSet []hmTypes.PartyID) (map[string]json.RawMessage, error) {
	// set state to metis state
	mpcState := GetGenesisStateFromAppState(appState)
	mpcState.MpcSet = currentValSet
	appState[ModuleName] = types.ModuleCdc.MustMarshalJSON(mpcState)

	return appState, nil
}
