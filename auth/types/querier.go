package types

import (
	"github.com/metis-seq/themis/types"
)

// query endpoints supported by the auth Querier
const (
	QueryParams  = "params"
	QueryAccount = "account"
)

// QueryAccountParams defines the params for querying accounts.
type QueryAccountParams struct {
	Address types.ThemisAddress
}

// NewQueryAccountParams creates a new instance of QueryAccountParams.
func NewQueryAccountParams(addr types.ThemisAddress) QueryAccountParams {
	return QueryAccountParams{Address: addr}
}
