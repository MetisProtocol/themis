package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "bank"

	// StoreKey is the store key string for metis
	StoreKey = ModuleName

	// RouterKey is the message route for metis
	RouterKey = ModuleName

	// QuerierRoute is the querier route for metis
	QuerierRoute = ModuleName

	// DefaultParamspace default name for parameter store
	DefaultParamspace = ModuleName

	// DefaultCodespace default code space
	DefaultCodespace sdk.CodespaceType = ModuleName
)
