package types

import (
	"github.com/metis-seq/themis/params/subspace"
)

const (
	// DefaultSendEnabled enabled
	DefaultSendEnabled = true
)

// ParamStoreKeySendEnabled is store's key for SendEnabled
var ParamStoreKeySendEnabled = []byte("sendenabled")

// ParamKeyTable type declaration for parameters
func ParamKeyTable() subspace.KeyTable {
	return subspace.NewKeyTable(
		ParamStoreKeySendEnabled, false,
	)
}
