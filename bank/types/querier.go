package types

import (
	hmTyps "github.com/metis-seq/themis/types"
)

const (
	QueryBalance = "balances"
)

// QueryBalanceParams defines the params for querying an account balance.
type QueryBalanceParams struct {
	Address hmTyps.ThemisAddress
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryBalanceParams(addr hmTyps.ThemisAddress) QueryBalanceParams {
	return QueryBalanceParams{Address: addr}
}
