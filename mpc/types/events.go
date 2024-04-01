package types

// mpc module event types
const (
	EventTypeProposeMpcCreate = "propose-mpc-create"
	EventTypeProposeMpcSign   = "propose-mpc-sign"
	EventTypeMpcSign          = "mpc-sign"

	AttributeKeySuccess      = "success"
	AttributeKeyMpcID        = "mpc-id"
	AttributeKeyMpcAddress   = "mpc-address"
	AttributeKeyMpcSignID    = "mpc-sign-id"
	AttributeKeyMpcSignMsg   = "mpc-sign-msg"
	AttributeKeyMpcSignature = "mpc-signature"
	AttributeKeyMpcSignType  = "mpc-sign_type"

	AttributeValueCategory = ModuleName
)
