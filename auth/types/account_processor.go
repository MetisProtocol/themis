package types

import (
	"github.com/metis-seq/themis/auth/exported"
)

// AccountProcessor is an interface to process account as per module
type AccountProcessor func(*GenesisAccount, *BaseAccount) exported.Account
