package types

import (
	"bytes"
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/yaml.v3"

	"github.com/ethereum/go-ethereum/common"
)

const (
	// AddrLen defines a valid address length
	AddrLen = 20
)

// Ensure that different address types implement the interface
var _ sdk.Address = ThemisAddress{}
var _ yaml.Marshaler = ThemisAddress{}

// ThemisAddress represents themis address
type ThemisAddress common.Address

// ZeroThemisAddress represents zero address
var ZeroThemisAddress = ThemisAddress{}

// EthAddress get eth address
func (aa ThemisAddress) EthAddress() common.Address {
	return common.Address(aa)
}

// Equals returns boolean for whether two AccAddresses are Equal
func (aa ThemisAddress) Equals(aa2 sdk.Address) bool {
	if aa.Empty() && aa2.Empty() {
		return true
	}

	return bytes.Equal(aa.Bytes(), aa2.Bytes())
}

// Empty returns boolean for whether an AccAddress is empty
func (aa ThemisAddress) Empty() bool {
	return bytes.Equal(aa.Bytes(), ZeroThemisAddress.Bytes())
}

// Marshal returns the raw address bytes. It is needed for protobuf
// compatibility.
func (aa ThemisAddress) Marshal() ([]byte, error) {
	return aa.Bytes(), nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (aa *ThemisAddress) Unmarshal(data []byte) error {
	*aa = ThemisAddress(common.BytesToAddress(data))
	return nil
}

// MarshalJSON marshals to JSON using Bech32.
func (aa ThemisAddress) MarshalJSON() ([]byte, error) {
	return jsoniter.ConfigFastest.Marshal(aa.String())
}

// MarshalYAML marshals to YAML using Bech32.
func (aa ThemisAddress) MarshalYAML() (interface{}, error) {
	return aa.String(), nil
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (aa *ThemisAddress) UnmarshalJSON(data []byte) error {
	var s string
	if err := jsoniter.ConfigFastest.Unmarshal(data, &s); err != nil {
		return err
	}

	*aa = HexToThemisAddress(s)

	return nil
}

// UnmarshalYAML unmarshals from JSON assuming Bech32 encoding.
func (aa *ThemisAddress) UnmarshalYAML(data []byte) error {
	var s string
	if err := yaml.Unmarshal(data, &s); err != nil {
		return err
	}

	*aa = HexToThemisAddress(s)

	return nil
}

// Bytes returns the raw address bytes.
func (aa ThemisAddress) Bytes() []byte {
	return aa[:]
}

// String implements the Stringer interface.
func (aa ThemisAddress) String() string {
	return "0x" + hex.EncodeToString(aa.Bytes())
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (aa ThemisAddress) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(aa.String()))
	case 'p':
		s.Write([]byte(fmt.Sprintf("%p", aa)))
	default:
		s.Write([]byte(fmt.Sprintf("%X", aa.Bytes())))
	}
}

//
// Address utils
//

// BytesToThemisAddress returns Address with value b.
func BytesToThemisAddress(b []byte) ThemisAddress {
	return ThemisAddress(common.BytesToAddress(b))
}

// HexToThemisAddress returns Address with value b.
func HexToThemisAddress(b string) ThemisAddress {
	return ThemisAddress(common.HexToAddress(b))
}

// AccAddressToThemisAddress returns Address with value b.
func AccAddressToThemisAddress(b sdk.AccAddress) ThemisAddress {
	return BytesToThemisAddress(b[:])
}

// ThemisAddressToAccAddress returns Address with value b.
func ThemisAddressToAccAddress(b ThemisAddress) sdk.AccAddress {
	return sdk.AccAddress(b.Bytes())
}

// SampleThemisAddress returns sample address
func SampleThemisAddress(s string) ThemisAddress {
	return BytesToThemisAddress([]byte(s))
}
