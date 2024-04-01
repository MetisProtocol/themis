package sqlite

type MetisTxCache struct {
	ID     uint64 `json:"id" sql:"id"`
	TxData string `json:"tx_data" sql:"tx_data"`
}
