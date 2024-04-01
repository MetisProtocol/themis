package types

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/metis-seq/themis/helper"
	"github.com/metis-seq/themis/params/subspace"
	hmTypes "github.com/metis-seq/themis/types"
)

// Default parameter values
const (
	DefaultMainchainTxConfirmations  uint64 = 6
	DefaultMetischainTxConfirmations uint64 = 1
)

var (
	DefaultValidatorSetAddress   hmTypes.ThemisAddress = hmTypes.HexToThemisAddress("0x0000000000000000000000000000000000001000")
	DefaultMetisTokenAddress     hmTypes.ThemisAddress = hmTypes.HexToThemisAddress("0x114f836434A9aa9ca584491E7965b16565bF5d7b")
	DefaultStakingManagerAddress hmTypes.ThemisAddress = hmTypes.HexToThemisAddress("0x73D5B3D9C5502953E51E3dDeDFf81A3e86FA874D")
	DefaultStakingInfoAddress    hmTypes.ThemisAddress = hmTypes.HexToThemisAddress("0x33CdB54Fb5B0A469adB7D294dd868f4b782E2fBA")
)

// Parameter keys
var (
	KeyMainchainTxConfirmations  = []byte("MainchainTxConfirmations")
	KeyMetischainTxConfirmations = []byte("MetischainTxConfirmations")
	KeyChainParams               = []byte("ChainParams")
)

var _ subspace.ParamSet = &Params{}

// ChainParams chain related params
type ChainParams struct {
	MainChainID           string                `json:"main_chain_id" yaml:"main_chain_id"`
	MetisChainID          string                `json:"metis_chain_id" yaml:"metis_chain_id"`
	MetisTokenAddress     hmTypes.ThemisAddress `json:"metis_token_address" yaml:"metis_token_address"`
	StakingManagerAddress hmTypes.ThemisAddress `json:"staking_manager_address" yaml:"staking_manager_address"`
	StakingInfoAddress    hmTypes.ThemisAddress `json:"staking_info_address" yaml:"staking_info_address"`

	// Metis Chain Contracts
	ValidatorSetAddress hmTypes.ThemisAddress `json:"validator_set_address" yaml:"validator_set_address"`
}

func (cp ChainParams) String() string {
	return fmt.Sprintf(`
	MainChainID: 									%s
	MetisChainID: 									%s
  MetisTokenAddress:            %s
	StakingManagerAddress:        %s
  StakingInfoAddress:           %s
	ValidatorSetAddress:					%s`,
		cp.MainChainID, cp.MetisChainID, cp.MetisTokenAddress, cp.StakingManagerAddress, cp.StakingInfoAddress, cp.ValidatorSetAddress)
}

// Params defines the parameters for the chainmanager module.
type Params struct {
	MainchainTxConfirmations  uint64      `json:"mainchain_tx_confirmations" yaml:"mainchain_tx_confirmations"`
	MetischainTxConfirmations uint64      `json:"metischain_tx_confirmations" yaml:"metischain_tx_confirmations"`
	ChainParams               ChainParams `json:"chain_params" yaml:"chain_params"`
}

// NewParams creates a new Params object
func NewParams(mainchainTxConfirmations uint64, metischainTxConfirmations uint64, chainParams ChainParams) Params {
	return Params{
		MainchainTxConfirmations:  mainchainTxConfirmations,
		MetischainTxConfirmations: metischainTxConfirmations,
		ChainParams:               chainParams,
	}
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		{KeyMainchainTxConfirmations, &p.MainchainTxConfirmations},
		{KeyMetischainTxConfirmations, &p.MetischainTxConfirmations},
		{KeyChainParams, &p.ChainParams},
	}
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)

	return bytes.Equal(bz1, bz2)
}

// String implements the stringer interface.
func (p Params) String() string {
	var sb strings.Builder

	sb.WriteString("Params: \n")
	sb.WriteString(fmt.Sprintf("MainchainTxConfirmations: %d\n", p.MainchainTxConfirmations))
	sb.WriteString(fmt.Sprintf("MetischainTxConfirmations: %d\n", p.MetischainTxConfirmations))
	sb.WriteString(fmt.Sprintf("ChainParams: %s\n", p.ChainParams.String()))

	return sb.String()
}

// Validate checks that the parameters have valid values.
func (p Params) Validate() error {
	if err := validateThemisAddress("metis_token_address", p.ChainParams.MetisTokenAddress); err != nil {
		return err
	}

	if err := validateThemisAddress("staking_manager_address", p.ChainParams.StakingManagerAddress); err != nil {
		return err
	}

	if err := validateThemisAddress("staking_info_address", p.ChainParams.StakingInfoAddress); err != nil {
		return err
	}

	if err := validateThemisAddress("validator_set_address", p.ChainParams.ValidatorSetAddress); err != nil {
		return err
	}

	return nil
}

func validateThemisAddress(key string, value hmTypes.ThemisAddress) error {
	if value.String() == "" {
		return fmt.Errorf("Invalid value %s in chain_params", key)
	}

	return nil
}

//
// Extra functions
//

// ParamKeyTable for auth module
func ParamKeyTable() subspace.KeyTable {
	return subspace.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		MainchainTxConfirmations:  DefaultMainchainTxConfirmations,
		MetischainTxConfirmations: DefaultMetischainTxConfirmations,
		ChainParams: ChainParams{
			MainChainID:           helper.DefaultMainChainID,
			MetisChainID:          helper.DefaultMetisChainID,
			ValidatorSetAddress:   DefaultValidatorSetAddress,
			MetisTokenAddress:     DefaultMetisTokenAddress,
			StakingManagerAddress: DefaultStakingManagerAddress,
			StakingInfoAddress:    DefaultStakingInfoAddress,
		},
	}
}
