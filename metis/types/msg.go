package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	hmTypes "github.com/metis-seq/themis/types"
)

//
// Propose Span Msg
//

var _ sdk.Msg = &MsgProposeSpan{}
var _ sdk.Msg = &MsgReProposeSpan{}
var _ sdk.Msg = &MsgMetisTx{}

// MsgProposeSpan creates msg propose span
type MsgProposeSpan struct {
	ID              uint64                  `json:"span_id"`
	Proposer        hmTypes.ThemisAddress   `json:"proposer"`
	CurrentL2Height uint64                  `json:"current_l2_height"`
	StartBlock      uint64                  `json:"start_block"`
	EndBlock        uint64                  `json:"end_block"`
	ChainID         string                  `json:"metis_chain_id"`
	Seed            common.Hash             `json:"seed"`
	IsRecover       bool                    `json:"is_recover"`
	RecoverSigner   hmTypes.ThemisAddress   `json:"recover_signer"`
	BlackList       []hmTypes.ThemisAddress `json:"black_list"`
}

// NewMsgProposeSpan creates new propose span message
func NewMsgProposeSpan(
	id uint64,
	proposer hmTypes.ThemisAddress,
	currentL2Height uint64,
	startBlock uint64,
	endBlock uint64,
	chainID string,
	seed common.Hash,
	isRecover bool,
	recoverSigner hmTypes.ThemisAddress,
	blackList []hmTypes.ThemisAddress,
) MsgProposeSpan {
	return MsgProposeSpan{
		ID:              id,
		Proposer:        proposer,
		CurrentL2Height: currentL2Height,
		StartBlock:      startBlock,
		EndBlock:        endBlock,
		ChainID:         chainID,
		Seed:            seed,
		IsRecover:       isRecover,
		RecoverSigner:   recoverSigner,
		BlackList:       blackList,
	}
}

// Type returns message type
func (msg MsgProposeSpan) Type() string {
	return "propose-span"
}

// Route returns route for message
func (msg MsgProposeSpan) Route() string {
	return RouterKey
}

// GetSigners returns address of the signer
func (msg MsgProposeSpan) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{hmTypes.ThemisAddressToAccAddress(msg.Proposer)}
}

// GetSignBytes returns sign bytes for proposeSpan message type
func (msg MsgProposeSpan) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic validates the message and returns error
func (msg MsgProposeSpan) ValidateBasic() sdk.Error {
	if msg.Proposer.Empty() {
		return sdk.ErrInvalidAddress(msg.Proposer.String())
	}

	return nil
}

// GetSideSignBytes returns side sign bytes
func (msg MsgProposeSpan) GetSideSignBytes() []byte {
	return nil
}

// MsgReProposeSpan creates msg propose span
type MsgReProposeSpan struct {
	ID              uint64                `json:"span_id"`
	Proposer        hmTypes.ThemisAddress `json:"proposer"`
	CurrentProducer hmTypes.ThemisAddress `json:"current_producer"`
	NextProducer    hmTypes.ThemisAddress `json:"next_producer"`
	CurrentL2Height uint64                `json:"current_l2_height"`
	CurrentL2Epoch  uint64                `json:"current_l2_epoch"`
	StartBlock      uint64                `json:"start_block"`
	EndBlock        uint64                `json:"end_block"`
	ChainID         string                `json:"metis_chain_id"`
	Seed            common.Hash           `json:"seed"`
}

// NewMsgReProposeSpan creates new propose span message
func NewMsgReProposeSpan(
	id uint64,
	proposer hmTypes.ThemisAddress,
	currentProducer hmTypes.ThemisAddress,
	nextProducer hmTypes.ThemisAddress,
	currentL2Height uint64,
	currentL2Epoch uint64,
	startBlock uint64,
	endBlock uint64,
	chainID string,
	seed common.Hash,
) MsgReProposeSpan {
	return MsgReProposeSpan{
		ID:              id,
		Proposer:        proposer,
		CurrentProducer: currentProducer,
		NextProducer:    nextProducer,
		CurrentL2Height: currentL2Height,
		CurrentL2Epoch:  currentL2Epoch,
		StartBlock:      startBlock,
		EndBlock:        endBlock,
		ChainID:         chainID,
		Seed:            seed,
	}
}

// Type returns message type
func (msg MsgReProposeSpan) Type() string {
	return "re-propose-span"
}

// Route returns route for message
func (msg MsgReProposeSpan) Route() string {
	return RouterKey
}

// GetSigners returns address of the signer
func (msg MsgReProposeSpan) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{hmTypes.ThemisAddressToAccAddress(msg.Proposer)}
}

// GetSignBytes returns sign bytes for proposeSpan message type
func (msg MsgReProposeSpan) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic validates the message and returns error
func (msg MsgReProposeSpan) ValidateBasic() sdk.Error {
	if msg.Proposer.Empty() {
		return sdk.ErrInvalidAddress(msg.Proposer.String())
	}

	return nil
}

// GetSideSignBytes returns side sign bytes
func (msg MsgReProposeSpan) GetSideSignBytes() []byte {
	return nil
}

// MsgMetisTx creates msg metis tx
type MsgMetisTx struct {
	Proposer hmTypes.ThemisAddress `json:"proposer"`
	TxHash   hmTypes.ThemisHash    `json:"tx_hash"`
	LogIndex uint64                `json:"log_index"`
	TxData   string                `json:"tx_data"`
}

// NewMsgMetisTx creates new metis tx message
func NewMsgMetisTx(
	proposer hmTypes.ThemisAddress,
	txHash hmTypes.ThemisHash,
	LogIndex uint64,
	txData string,
) MsgMetisTx {
	return MsgMetisTx{
		Proposer: proposer,
		TxHash:   txHash,
		LogIndex: LogIndex,
		TxData:   txData,
	}
}

// Type returns message type
func (msg MsgMetisTx) Type() string {
	return "metis-tx"
}

// Route returns route for message
func (msg MsgMetisTx) Route() string {
	return RouterKey
}

// GetSigners returns address of the signer
func (msg MsgMetisTx) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{hmTypes.ThemisAddressToAccAddress(msg.Proposer)}
}

// GetSignBytes returns sign bytes for proposeSpan message type
func (msg MsgMetisTx) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic validates the message and returns error
func (msg MsgMetisTx) ValidateBasic() sdk.Error {
	if msg.Proposer.Empty() {
		return sdk.ErrInvalidAddress(msg.Proposer.String())
	}

	return nil
}

// GetSideSignBytes returns side sign bytes
func (msg MsgMetisTx) GetSideSignBytes() []byte {
	return nil
}
