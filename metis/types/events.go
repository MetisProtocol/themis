package types

// staking module event types
const (
	EventTypeProposeSpan   = "propose-span"
	EventTypeReProposeSpan = "re-propose-span"
	EventTypeDeleteReSpan  = "delete-re-span"
	EventTypeMetisTx       = "metis-tx"

	AttributeKeySuccess        = "success"
	AttributeKeySpanID         = "span-id"
	AttributeKeyOldSpanID      = "old-span-id"
	AttributeKeySpanStartBlock = "start-block"
	AttributeKeySpanEndBlock   = "end-block"
	AttributeKeySpanSigner     = "span-signer"
	AttributeKeySpanRecover    = "is-span-recover"
	AttributeKeyMetisTxHash    = "metis-tx-hash"
	AttributeKeyMetisTxData    = "metis-tx-data"

	AttributeValueCategory = ModuleName
)
