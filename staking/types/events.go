package types

var (
	EventTypeNewProposer       = "new-proposer"
	EventTypeValidatorJoin     = "validator-join"
	EventTypeSignerUpdate      = "signer-update"
	EventTypeStakeUpdate       = "stake-update"
	EventTypeValidatorExit     = "validator-exit"
	EventTypeBatchSubmitReward = "batch-submit-reward"

	AttributeKeySigner            = "signer"
	AttributeKeyDeactivationBatch = "deactivation-batch"
	AttributeKeyActivationBatch   = "activation-batch"
	AttributeKeyValidatorID       = "validator-id"
	AttributeKeyBatchID           = "batch-id"
	AttributeKeyValidatorNonce    = "validator-nonce"
	AttributeKeyUpdatedAt         = "updated-at"

	AttributeValueCategory = ModuleName
)
