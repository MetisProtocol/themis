package types

import (
	"bytes"
	"encoding/hex"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"gopkg.in/yaml.v3"

	"github.com/ethereum/go-ethereum/common"
)

// Ensure that different address types implement the interface
var _ yaml.Marshaler = ThemisHash{}

// ThemisHash represents themis address
type ThemisHash common.Hash

// ZeroThemisHash represents zero address
var ZeroThemisHash = ThemisHash{}

// EthHash get eth hash
func (aa ThemisHash) EthHash() common.Hash {
	return common.Hash(aa)
}

// Equals returns boolean for whether two ThemisHash are Equal
func (aa ThemisHash) Equals(aa2 ThemisHash) bool {
	if aa.Empty() && aa2.Empty() {
		return true
	}

	return bytes.Equal(aa.Bytes(), aa2.Bytes())
}

// Empty returns boolean for whether an AccAddress is empty
func (aa ThemisHash) Empty() bool {
	return bytes.Equal(aa.Bytes(), ZeroThemisHash.Bytes())
}

// Marshal returns the raw address bytes. It is needed for protobuf
// compatibility.
func (aa ThemisHash) Marshal() ([]byte, error) {
	return aa.Bytes(), nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (aa *ThemisHash) Unmarshal(data []byte) error {
	*aa = ThemisHash(common.BytesToHash(data))
	return nil
}

// MarshalJSON marshals to JSON using Bech32.
func (aa ThemisHash) MarshalJSON() ([]byte, error) {
	return jsoniter.ConfigFastest.Marshal(aa.String())
}

// MarshalYAML marshals to YAML using Bech32.
func (aa ThemisHash) MarshalYAML() (interface{}, error) {
	return aa.String(), nil
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (aa *ThemisHash) UnmarshalJSON(data []byte) error {
	var s string
	if err := jsoniter.ConfigFastest.Unmarshal(data, &s); err != nil {
		return err
	}

	*aa = HexToThemisHash(s)

	return nil
}

// UnmarshalYAML unmarshals from JSON assuming Bech32 encoding.
func (aa *ThemisHash) UnmarshalYAML(data []byte) error {
	var s string
	if err := yaml.Unmarshal(data, &s); err != nil {
		return err
	}

	*aa = HexToThemisHash(s)

	return nil
}

// Bytes returns the raw address bytes.
func (aa ThemisHash) Bytes() []byte {
	return aa[:]
}

// String implements the Stringer interface.
func (aa ThemisHash) String() string {
	if aa.Empty() {
		return ""
	}

	return "0x" + hex.EncodeToString(aa.Bytes())
}

// Hex returns hex string
func (aa ThemisHash) Hex() string {
	return aa.String()
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (aa ThemisHash) Format(s fmt.State, verb rune) {
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
// hash utils
//

// BytesToThemisHash returns Address with value b.
func BytesToThemisHash(b []byte) ThemisHash {
	return ThemisHash(common.BytesToHash(b))
}

// HexToThemisHash returns Address with value b.
func HexToThemisHash(b string) ThemisHash {
	return ThemisHash(common.HexToHash(b))
}
