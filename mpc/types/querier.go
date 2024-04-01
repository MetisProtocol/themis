package types

// query endpoints supported by the auth Querier
const (
	QueryMpc       = "mpc"
	QueryMpcList   = "mpc-list"
	QueryLatestMpc = "latest-mpc"
	QueryMpcSet    = "mpc-set"
	QueryMpcSign   = "mpc-sign"
)

// QueryMpcParams defines the params for querying accounts.
type QueryMpcParams struct {
	MpcID string
}

// NewQueryMpcParams creates a new instance of QueryMpcParams.
func NewQueryMpcParams(mpcID string) QueryMpcParams {
	return QueryMpcParams{MpcID: mpcID}
}

// QueryLatestMpcParams defines the params for querying accounts.
type QueryLatestMpcParams struct {
	Type int
}

// NewQueryLatestMpcParams creates a new instance of QueryLatestMpcParams.
func NewQueryLatestMpcParams(mpcType int) QueryLatestMpcParams {
	return QueryLatestMpcParams{Type: mpcType}
}

// QueryMpcSignParams defines the params for querying accounts.
type QueryMpcSignParams struct {
	SignID string
}

// NewQueryMpcSignParams creates a new instance of QueryMpcSignParams.
func NewQueryMpcSignParams(signID string) QueryMpcSignParams {
	return QueryMpcSignParams{SignID: signID}
}
