package types

import (
	"bytes"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	hmCommon "github.com/metis-seq/themis/common"
	"github.com/metis-seq/themis/helper"
	"github.com/metis-seq/themis/types"
	hmTypes "github.com/metis-seq/themis/types"
)

var cdc = codec.New()

//
// Validator Join
//

var _ sdk.Msg = &MsgValidatorJoin{}

type MsgValidatorJoin struct {
	From            hmTypes.ThemisAddress `json:"from"`
	ID              hmTypes.ValidatorID   `json:"id"`
	ActivationBatch uint64                `json:"activationBatch"`
	Amount          sdk.Int               `json:"amount"`
	SignerPubKey    hmTypes.PubKey        `json:"pub_key"`
	TxHash          hmTypes.ThemisHash    `json:"tx_hash"`
	LogIndex        uint64                `json:"log_index"`
	BlockNumber     uint64                `json:"block_number"`
	Nonce           uint64                `json:"nonce"`
}

// NewMsgValidatorJoin creates new validator-join
func NewMsgValidatorJoin(
	from hmTypes.ThemisAddress,
	id uint64,
	activationBatch uint64,
	amount sdk.Int,
	pubkey hmTypes.PubKey,
	txhash hmTypes.ThemisHash,
	logIndex uint64,
	blockNumber uint64,
	nonce uint64,
) MsgValidatorJoin {
	return MsgValidatorJoin{
		From:            from,
		ID:              hmTypes.NewValidatorID(id),
		ActivationBatch: activationBatch,
		Amount:          amount,
		SignerPubKey:    pubkey,
		TxHash:          txhash,
		LogIndex:        logIndex,
		BlockNumber:     blockNumber,
		Nonce:           nonce,
	}
}

func (msg MsgValidatorJoin) Type() string {
	return "validator-join"
}

func (msg MsgValidatorJoin) Route() string {
	return RouterKey
}

func (msg MsgValidatorJoin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{hmTypes.ThemisAddressToAccAddress(msg.From)}
}

func (msg MsgValidatorJoin) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

func (msg MsgValidatorJoin) ValidateBasic() sdk.Error {
	if msg.ID == 0 {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid validator ID %v", msg.ID)
	}

	if bytes.Equal(msg.SignerPubKey.Bytes(), helper.ZeroPubKey.Bytes()) {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid pub key %v", msg.SignerPubKey.String())
	}

	if msg.From.Empty() {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid proposer %v", msg.From.String())
	}

	return nil
}

// GetTxHash Returns tx hash
func (msg MsgValidatorJoin) GetTxHash() types.ThemisHash {
	return msg.TxHash
}

// GetLogIndex Returns log index
func (msg MsgValidatorJoin) GetLogIndex() uint64 {
	return msg.LogIndex
}

// GetSideSignBytes returns side sign bytes
func (msg MsgValidatorJoin) GetSideSignBytes() []byte {
	return nil
}

// GetNonce Returns nonce
func (msg MsgValidatorJoin) GetNonce() uint64 {
	return msg.Nonce
}

//
// Stake update
//

//
// validator exit
//

var _ sdk.Msg = &MsgStakeUpdate{}

// MsgStakeUpdate represents stake update
type MsgStakeUpdate struct {
	From        hmTypes.ThemisAddress `json:"from"`
	ID          hmTypes.ValidatorID   `json:"id"`
	NewAmount   sdk.Int               `json:"amount"`
	TxHash      hmTypes.ThemisHash    `json:"tx_hash"`
	LogIndex    uint64                `json:"log_index"`
	BlockNumber uint64                `json:"block_number"`
	Nonce       uint64                `json:"nonce"`
}

// NewMsgStakeUpdate represents stake update
func NewMsgStakeUpdate(from hmTypes.ThemisAddress, id uint64, newAmount sdk.Int, txhash hmTypes.ThemisHash, logIndex uint64, blockNumber uint64, nonce uint64) MsgStakeUpdate {
	return MsgStakeUpdate{
		From:        from,
		ID:          hmTypes.NewValidatorID(id),
		NewAmount:   newAmount,
		TxHash:      txhash,
		LogIndex:    logIndex,
		BlockNumber: blockNumber,
		Nonce:       nonce,
	}
}

func (msg MsgStakeUpdate) Type() string {
	return "validator-stake-update"
}

func (msg MsgStakeUpdate) Route() string {
	return RouterKey
}

func (msg MsgStakeUpdate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{hmTypes.ThemisAddressToAccAddress(msg.From)}
}

func (msg MsgStakeUpdate) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

func (msg MsgStakeUpdate) ValidateBasic() sdk.Error {
	if msg.ID == 0 {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid validator ID %v", msg.ID)
	}

	if msg.From.Empty() {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid proposer %v", msg.From.String())
	}

	return nil
}

// GetTxHash Returns tx hash
func (msg MsgStakeUpdate) GetTxHash() types.ThemisHash {
	return msg.TxHash
}

// GetLogIndex Returns log index
func (msg MsgStakeUpdate) GetLogIndex() uint64 {
	return msg.LogIndex
}

// GetSideSignBytes returns side sign bytes
func (msg MsgStakeUpdate) GetSideSignBytes() []byte {
	return nil
}

// GetNonce Returns nonce
func (msg MsgStakeUpdate) GetNonce() uint64 {
	return msg.Nonce
}

// validator update
var _ sdk.Msg = &MsgSignerUpdate{}

// MsgSignerUpdate signer update struct
// TODO add old signer sig check
type MsgSignerUpdate struct {
	From            hmTypes.ThemisAddress `json:"from"`
	ID              hmTypes.ValidatorID   `json:"id"`
	NewSignerPubKey hmTypes.PubKey        `json:"pubKey"`
	TxHash          hmTypes.ThemisHash    `json:"tx_hash"`
	LogIndex        uint64                `json:"log_index"`
	BlockNumber     uint64                `json:"block_number"`
	Nonce           uint64                `json:"nonce"`
}

func NewMsgSignerUpdate(
	from hmTypes.ThemisAddress,
	id uint64,
	pubKey hmTypes.PubKey,
	txhash hmTypes.ThemisHash,
	logIndex uint64,
	blockNumber uint64,
	nonce uint64,
) MsgSignerUpdate {
	return MsgSignerUpdate{
		From:            from,
		ID:              hmTypes.NewValidatorID(id),
		NewSignerPubKey: pubKey,
		TxHash:          txhash,
		LogIndex:        logIndex,
		BlockNumber:     blockNumber,
		Nonce:           nonce,
	}
}

func (msg MsgSignerUpdate) Type() string {
	return "signer-update"
}

func (msg MsgSignerUpdate) Route() string {
	return RouterKey
}

func (msg MsgSignerUpdate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{hmTypes.ThemisAddressToAccAddress(msg.From)}
}

func (msg MsgSignerUpdate) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

func (msg MsgSignerUpdate) ValidateBasic() sdk.Error {
	if msg.ID == 0 {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid validator ID %v", msg.ID)
	}

	if msg.From.Empty() {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid proposer %v", msg.From.String())
	}

	if bytes.Equal(msg.NewSignerPubKey.Bytes(), helper.ZeroPubKey.Bytes()) {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid pub key %v", msg.NewSignerPubKey.String())
	}

	return nil
}

// GetTxHash Returns tx hash
func (msg MsgSignerUpdate) GetTxHash() types.ThemisHash {
	return msg.TxHash
}

// GetLogIndex Returns log index
func (msg MsgSignerUpdate) GetLogIndex() uint64 {
	return msg.LogIndex
}

// GetSideSignBytes returns side sign bytes
func (msg MsgSignerUpdate) GetSideSignBytes() []byte {
	return nil
}

// GetNonce Returns nonce
func (msg MsgSignerUpdate) GetNonce() uint64 {
	return msg.Nonce
}

//
// validator exit
//

var _ sdk.Msg = &MsgValidatorExit{}

type MsgValidatorExit struct {
	From              hmTypes.ThemisAddress `json:"from"`
	ID                hmTypes.ValidatorID   `json:"id"`
	DeactivationBatch uint64                `json:"deactivationBatch"`
	TxHash            hmTypes.ThemisHash    `json:"tx_hash"`
	LogIndex          uint64                `json:"log_index"`
	BlockNumber       uint64                `json:"block_number"`
	Nonce             uint64                `json:"nonce"`
}

func NewMsgValidatorExit(from hmTypes.ThemisAddress, id uint64, deactivationBatch uint64, txhash hmTypes.ThemisHash, logIndex uint64, blockNumber uint64, nonce uint64) MsgValidatorExit {
	return MsgValidatorExit{
		From:              from,
		ID:                hmTypes.NewValidatorID(id),
		DeactivationBatch: deactivationBatch,
		TxHash:            txhash,
		LogIndex:          logIndex,
		BlockNumber:       blockNumber,
		Nonce:             nonce,
	}
}

func (msg MsgValidatorExit) Type() string {
	return "validator-exit"
}

func (msg MsgValidatorExit) Route() string {
	return RouterKey
}

func (msg MsgValidatorExit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{hmTypes.ThemisAddressToAccAddress(msg.From)}
}

func (msg MsgValidatorExit) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

func (msg MsgValidatorExit) ValidateBasic() sdk.Error {
	if msg.ID == 0 {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid validator ID %v", msg.ID)
	}

	if msg.From.Empty() {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid proposer %v", msg.From.String())
	}

	return nil
}

// GetTxHash Returns tx hash
func (msg MsgValidatorExit) GetTxHash() types.ThemisHash {
	return msg.TxHash
}

// GetLogIndex Returns log index
func (msg MsgValidatorExit) GetLogIndex() uint64 {
	return msg.LogIndex
}

// GetSideSignBytes returns side sign bytes
func (msg MsgValidatorExit) GetSideSignBytes() []byte {
	return nil
}

// GetNonce Returns nonce
func (msg MsgValidatorExit) GetNonce() uint64 {
	return msg.Nonce
}

//
// batch submit reward
//

var _ sdk.Msg = &MsgBatchSubmitReward{}

type MsgBatchSubmitReward struct {
	From        hmTypes.ThemisAddress `json:"from"`
	ID          hmTypes.ValidatorID   `json:"id"`
	BatchID     uint64                `json:"batchId"`
	TxHash      hmTypes.ThemisHash    `json:"tx_hash"`
	LogIndex    uint64                `json:"log_index"`
	BlockNumber uint64                `json:"block_number"`
	Nonce       uint64                `json:"nonce"`
}

func NewMsgBatchSubmitReward(from hmTypes.ThemisAddress, id, batchId uint64, txhash hmTypes.ThemisHash, logIndex, blockNumber, nonce uint64) MsgBatchSubmitReward {
	return MsgBatchSubmitReward{
		From:        from,
		ID:          hmTypes.NewValidatorID(id),
		BatchID:     batchId,
		TxHash:      txhash,
		LogIndex:    logIndex,
		BlockNumber: blockNumber,
		Nonce:       nonce,
	}
}

func (msg MsgBatchSubmitReward) Type() string {
	return "batch-submit-reward"
}

func (msg MsgBatchSubmitReward) Route() string {
	return RouterKey
}

func (msg MsgBatchSubmitReward) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{hmTypes.ThemisAddressToAccAddress(msg.From)}
}

func (msg MsgBatchSubmitReward) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

func (msg MsgBatchSubmitReward) ValidateBasic() sdk.Error {
	if msg.ID == 0 {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid validator ID %v", msg.ID)
	}

	if msg.From.Empty() {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid proposer %v", msg.From.String())
	}

	return nil
}

// GetTxHash Returns tx hash
func (msg MsgBatchSubmitReward) GetTxHash() types.ThemisHash {
	return msg.TxHash
}

// GetLogIndex Returns log index
func (msg MsgBatchSubmitReward) GetLogIndex() uint64 {
	return msg.LogIndex
}

// GetSideSignBytes returns side sign bytes
func (msg MsgBatchSubmitReward) GetSideSignBytes() []byte {
	return nil
}

// GetNonce Returns nonce
func (msg MsgBatchSubmitReward) GetNonce() uint64 {
	return msg.Nonce
}
