package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	hmTypes "github.com/metis-seq/themis/types"
)

//
// Propose Mpc Msg
//

var _ sdk.Msg = &MsgProposeMpcCreate{}
var _ sdk.Msg = &MsgProposeMpcSign{}
var _ sdk.Msg = &MsgMpcSign{}

// MsgProposeMpcCreate creates msg mpc create
type MsgProposeMpcCreate struct {
	ID           string                `json:"mpc_id"`
	Threshold    uint64                `json:"threshold"`
	Participants []hmTypes.PartyID     `json:"participants"`
	Proposer     hmTypes.ThemisAddress `json:"proposer"`
	MpcAddress   hmTypes.ThemisAddress `json:"mpc_address"`
	MpcPubkey    []byte                `json:"mpc_pubkey"`
	MpcType      hmTypes.MpcType       `json:"mpc_type"`
}

// NewMsgProposeMpcCreate creates new mpc create message
func NewMsgProposeMpcCreate(
	id string,
	threshold uint64,
	participants []hmTypes.PartyID,
	proposer hmTypes.ThemisAddress,
	mpcAddress hmTypes.ThemisAddress,
	mpcPubkey []byte,
	mpcType hmTypes.MpcType,
) MsgProposeMpcCreate {
	return MsgProposeMpcCreate{
		ID:           id,
		Threshold:    threshold,
		Participants: participants,
		Proposer:     proposer,
		MpcPubkey:    mpcPubkey,
		MpcType:      mpcType,
	}
}

// Type returns message type
func (msg MsgProposeMpcCreate) Type() string {
	return "propose-mpc-create"
}

// Route returns route for message
func (msg MsgProposeMpcCreate) Route() string {
	return RouterKey
}

// GetSigners returns address of the signer
func (msg MsgProposeMpcCreate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{hmTypes.ThemisAddressToAccAddress(msg.Proposer)}
}

// GetSignBytes returns sign bytes for proposeSpan message type
func (msg MsgProposeMpcCreate) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic validates the message and returns error
func (msg MsgProposeMpcCreate) ValidateBasic() sdk.Error {
	if msg.Threshold <= 0 || msg.Threshold >= uint64(len(msg.Participants)) {
		return sdk.ErrUnknownRequest("invalid threshold")
	}

	if len(msg.Participants) < 2 {
		return sdk.ErrUnknownRequest("invalid participants")
	}

	if msg.Proposer.Empty() {
		return sdk.ErrInvalidAddress(msg.Proposer.String())
	}

	return nil
}

// GetSideSignBytes returns side sign bytes
func (msg MsgProposeMpcCreate) GetSideSignBytes() []byte {
	return nil
}

// MsgProposeMpcSign creates msg mpc create
type MsgProposeMpcSign struct {
	ID       string                `json:"sign_id"`
	MpcID    string                `json:"mpc_id"`
	SignType hmTypes.SignType      `json:"sign_type"`
	SignData []byte                `json:"sign_data"`
	SignMsg  []byte                `json:"sign_msg"`
	Proposer hmTypes.ThemisAddress `json:"proposer"`
}

// NewMsgProposeMpcSign creates new mpc create message
func NewMsgProposeMpcSign(
	id string,
	mpcID string,
	signType hmTypes.SignType,
	signData []byte,
	signMsg []byte,
	proposer hmTypes.ThemisAddress,
) MsgProposeMpcSign {
	return MsgProposeMpcSign{
		ID:       id,
		MpcID:    mpcID,
		SignType: signType,
		SignData: signData,
		SignMsg:  signMsg,
		Proposer: proposer,
	}
}

// Type returns message type
func (msg MsgProposeMpcSign) Type() string {
	return "propose-mpc-sign"
}

// Route returns route for message
func (msg MsgProposeMpcSign) Route() string {
	return RouterKey
}

// GetSignBytes returns sign bytes for proposeSpan message type
func (msg MsgProposeMpcSign) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic validates the message and returns error
func (msg MsgProposeMpcSign) ValidateBasic() sdk.Error {
	if msg.ID == "" {
		return sdk.ErrUnknownRequest("invalid mpc sign id")
	}
	if msg.MpcID == "" {
		return sdk.ErrUnknownRequest("invalid mpc mpc key id")
	}

	if len(msg.SignData) <= 0 {
		return sdk.ErrUnknownRequest("invalid mpc sign data")
	}

	// if len(msg.SignMsg) <= 0 {
	// 	return sdk.ErrUnknownRequest("invalid mpc sign msg")
	// }

	return nil
}

// GetSigners returns address of the signer
func (msg MsgProposeMpcSign) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{hmTypes.ThemisAddressToAccAddress(msg.Proposer)}
}

// GetSideSignBytes returns side sign bytes
func (msg MsgProposeMpcSign) GetSideSignBytes() []byte {
	return nil
}

// MsgMpcSign creates msg mpc sign
type MsgMpcSign struct {
	ID        string                `json:"sign_id"`
	Signature string                `json:"signature"`
	Proposer  hmTypes.ThemisAddress `json:"proposer"`
}

// NewMsgMpcSign creates new mpc create message
func NewMsgMpcSign(
	id string,
	signature string,
	proposer hmTypes.ThemisAddress,
) MsgMpcSign {
	return MsgMpcSign{
		ID:        id,
		Signature: signature,
		Proposer:  proposer,
	}
}

// Type returns message type
func (msg MsgMpcSign) Type() string {
	return "mpc-sign"
}

// Route returns route for message
func (msg MsgMpcSign) Route() string {
	return RouterKey
}

// GetSignBytes returns sign bytes for proposeSpan message type
func (msg MsgMpcSign) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic validates the message and returns error
func (msg MsgMpcSign) ValidateBasic() sdk.Error {
	if msg.ID == "" {
		return sdk.ErrUnknownRequest("invalid mpc sign id")
	}

	if msg.Signature == "" {
		return sdk.ErrUnknownRequest("invalid mpc signature")
	}

	return nil
}

// GetSigners returns address of the signer
func (msg MsgMpcSign) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{hmTypes.ThemisAddressToAccAddress(msg.Proposer)}
}

// GetSideSignBytes returns side sign bytes
func (msg MsgMpcSign) GetSideSignBytes() []byte {
	return nil
}
