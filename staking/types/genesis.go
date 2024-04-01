package types

import (
	"encoding/json"
	"errors"

	"github.com/metis-seq/themis/metis/types"
	hmTypes "github.com/metis-seq/themis/types"
)

// GenesisValidator genesis validator
type GenesisValidator struct {
	ID         hmTypes.ValidatorID   `json:"id"`
	StartBatch uint64                `json:"start_batch"`
	EndBatch   uint64                `json:"end_batch"`
	Nonce      uint64                `json:"nonce"`
	Power      uint64                `json:"power"` // aka Amount
	PubKey     hmTypes.PubKey        `json:"pub_key"`
	Signer     hmTypes.ThemisAddress `json:"signer"`
}

// ThemisValidator converts genesis validator validator to Themis validator
func (v *GenesisValidator) ThemisValidator() hmTypes.Validator {
	return hmTypes.Validator{
		ID:          v.ID,
		PubKey:      v.PubKey,
		VotingPower: int64(v.Power),
		StartBatch:  v.StartBatch,
		EndBatch:    v.EndBatch,
		Nonce:       v.Nonce,
		Signer:      v.Signer,
	}
}

// GenesisState is the validator state that must be provided at genesis.
type GenesisState struct {
	Validators       []*hmTypes.Validator `json:"validators" yaml:"validators"`
	CurrentValSet    hmTypes.ValidatorSet `json:"current_val_set" yaml:"current_val_set"`
	StakingSequences []string             `json:"staking_sequences" yaml:"staking_sequences"`
	StakingBatch     uint64               `json:"staking_batch" yaml:"staking_batch"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(
	validators []*hmTypes.Validator,
	currentValSet hmTypes.ValidatorSet,
	stakingSequences []string,
) GenesisState {
	return GenesisState{
		Validators:       validators,
		CurrentValSet:    currentValSet,
		StakingSequences: stakingSequences,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(nil, hmTypes.ValidatorSet{}, nil)
}

// ValidateGenesis performs basic validation of metis genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	for _, validator := range data.Validators {
		if !validator.ValidateBasic() {
			return errors.New("Invalid validator")
		}
	}

	for _, sq := range data.StakingSequences {
		if sq == "" {
			return errors.New("Invalid Sequence")
		}
	}

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
func SetGenesisStateToAppState(appState map[string]json.RawMessage, validators []*hmTypes.Validator, currentValSet hmTypes.ValidatorSet) (map[string]json.RawMessage, error) {
	// set state to staking state
	stakingState := GetGenesisStateFromAppState(appState)
	stakingState.Validators = validators
	stakingState.CurrentValSet = currentValSet

	appState[ModuleName] = types.ModuleCdc.MustMarshalJSON(stakingState)

	return appState, nil
}
