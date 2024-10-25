package mpc

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"

	abci "github.com/tendermint/tendermint/abci/types"
	tmTypes "github.com/tendermint/tendermint/types"

	hmCommon "github.com/metis-seq/themis/common"
	"github.com/metis-seq/themis/helper"
	"github.com/metis-seq/themis/mpc/types"
	hmTypes "github.com/metis-seq/themis/types"
)

// NewSideTxHandler returns a side handler for "mpc" type messages.
func NewSideTxHandler(k Keeper) hmTypes.SideTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgProposeMpcCreate:
			return SideHandleMsgMpcCreate(ctx, k, msg)
		case types.MsgProposeMpcSign:
			return SideHandleMsgProposeMpcSign(ctx, k, msg)
		case types.MsgMpcSign:
			return SideHandleMsgMpcSign(ctx, k, msg)
		default:
			return abci.ResponseDeliverSideTx{
				Code: uint32(sdk.CodeUnknownRequest),
			}
		}
	}
}

// NewPostTxHandler returns a post handler for "mpc" type messages.
func NewPostTxHandler(k Keeper) hmTypes.PostTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg, sideTxResult abci.SideTxResultType) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgProposeMpcCreate:
			return PostHandleMsgEventMpcCreate(ctx, k, msg, sideTxResult)
		case types.MsgProposeMpcSign:
			return PostHandleMsgProposeMpcSign(ctx, k, msg, sideTxResult)
		case types.MsgMpcSign:
			return PostHandleMsgMpcSign(ctx, k, msg, sideTxResult)
		default:
			errMsg := "Unrecognized Mpc Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// SideHandleMsgMpcCreate validates external calls required for processing proposed mpc
func SideHandleMsgMpcCreate(ctx sdk.Context, k Keeper, msg types.MsgProposeMpcCreate) (result abci.ResponseDeliverSideTx) {
	k.Logger(ctx).Info("Validating External call for mpc create msg",
		"mpcId", msg.ID,
	)

	// check if mpc id exist.
	if k.HasMpc(ctx, msg.ID) {
		k.Logger(ctx).Info("mpc id already exist, ignore it")
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeMpcExist)
	}

	// check threshold
	if msg.Threshold <= 0 || msg.Threshold > uint64(len(msg.Participants)-1) {
		k.Logger(ctx).Error("invalid mpc threshold", "threshold", msg.Threshold, "partiesLen", len(msg.Participants))
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
	}

	// check mpc type
	switch msg.MpcType {
	case hmTypes.CommonMpcType, hmTypes.StateSubmitMpcType, hmTypes.RewardSubmitMpcType:
	default:
		k.Logger(ctx).Info("invalid mpc type, ignore it")
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeMpcInvalidType)
	}

	// get mpc set info
	mpcSetInfo := make(map[hmTypes.ThemisAddress]struct{})
	mpcSet := k.GetAllMpcSets(ctx)
	for _, set := range mpcSet {
		signer := hmTypes.HexToThemisAddress(set.Moniker)
		mpcSetInfo[signer] = struct{}{}
	}

	// check parties
	for _, party := range msg.Participants {
		signer := hmTypes.HexToThemisAddress(party.Moniker)
		if _, exist := mpcSetInfo[signer]; !exist {
			k.Logger(ctx).Error("invalid mpc party moniker", "moniker", party.Moniker)
			return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
		}
	}

	// check pubkey and address
	pubKey, err := btcec.ParsePubKey(msg.MpcPubkey, btcec.S256())
	if err != nil {
		k.Logger(ctx).Error("ParsePubKey err", "err", err)
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
	}
	if crypto.PubkeyToAddress(*pubKey.ToECDSA()) != msg.MpcAddress.EthAddress() {
		k.Logger(ctx).Error("pubkey and mpc address mismatch")
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
	}

	// check Proposer
	if _, exist := mpcSetInfo[msg.Proposer]; !exist {
		k.Logger(ctx).Error("invalid mpc proposer", "proposer", msg.Proposer)
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
	}

	// check mpc info from mpc server
	mpcPubKey, _, err := helper.GetMpcKey(msg.ID)
	if err != nil {
		k.Logger(ctx).Error("GetMpcKey err", "err", err)
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
	}
	if !bytes.Equal(mpcPubKey, pubKey.SerializeCompressed()) {
		k.Logger(ctx).Error("MpcCreate check err: pubkey mismatch", "mpcPubKey", hex.EncodeToString(mpcPubKey), "msgPubKey", pubKey, hex.EncodeToString(pubKey.SerializeCompressed()))
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
	}

	k.Logger(ctx).Info("Successfully validated External call for mpc create msg")

	result.Result = abci.SideTxResultType_Yes
	return
}

// SideHandleMsgProposeMpcSign validates external calls required for processing proposed mpc
func SideHandleMsgProposeMpcSign(ctx sdk.Context, k Keeper, msg types.MsgProposeMpcSign) (result abci.ResponseDeliverSideTx) {
	k.Logger(ctx).Info("Validating External call for mpc sign msg",
		"id", msg.ID,
	)

	// check if sign id exist.
	if k.HasMpcSign(ctx, msg.ID) {
		k.Logger(ctx).Info("mpc sign id already exist, ignore it", "id", msg.ID)
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeMpcSignExist)
	}

	// check if sign mpc id exist.
	if !k.HasMpc(ctx, msg.MpcID) {
		k.Logger(ctx).Info("mpc id not exist", "id", msg.ID)
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeMpcNotFound)
	}

	// check Proposer
	if !k.IsValidator(ctx, msg.Proposer) {
		k.Logger(ctx).Error("invalid mpc proposer", "proposer", msg.Proposer, "id", msg.ID)
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
	}

	// check sign msg
	_, txSignMsg, txData, err := GetCheckTxSignMsgData(ctx, k, msg.SignType, msg.SignData)
	if err != nil {
		k.Logger(ctx).Error("invalid mpc sign msg", "id", msg.ID, "sign msg", hex.EncodeToString(msg.SignMsg), "error", err)
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
	}
	if msg.SignMsg != nil && txSignMsg != nil && !bytes.Equal(txSignMsg, msg.SignMsg) {
		k.Logger(ctx).Error("invalid mpc sign msg", "id", msg.ID, "sign msg", hex.EncodeToString(msg.SignMsg))
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
	}

	// check sign data logic
	switch msg.SignType {
	case hmTypes.BatchSubmit:
		err = CheckBatchSubmitTxLogic(ctx, k, txData)
		if err != nil {
			k.Logger(ctx).Error("CheckBatchSubmitTxLogic failed", "id", msg.ID, "error", err)
			return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
		}
	case hmTypes.BatchReward:
		err = CheckBatchRewardTxLogic(ctx, k, txData, msg.SignMsg)
		if err != nil {
			k.Logger(ctx).Error("CheckBatchRewardTxLogic failed", "id", msg.ID, "error", err)
			return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
		}
	case hmTypes.CommitEpochToMetis:
		err = CheckCommitEpochToMetisTxLogic(ctx, k, txData)
		if err != nil {
			k.Logger(ctx).Error("CheckCommitEpochToMetisTxLogic failed", "id", msg.ID, "error", err)
			return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
		}
	case hmTypes.ReCommitEpochToMetis:
		err = CheckReCommitEpochToMetisTxLogic(ctx, k, txData)
		if err != nil {
			k.Logger(ctx).Error("ReCommitEpochToMetis failed", "id", msg.ID, "error", err)
			return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
		}
	case hmTypes.L1UpdateMpcAddress:
		err = CheckL1UpdateMpcAddressTxLogic(ctx, k, txData)
		if err != nil {
			k.Logger(ctx).Error("L1UpdateMpcAddress failed", "id", msg.ID, "error", err)
			return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
		}
	case hmTypes.L2UpdateMpcAddress:
		err = CheckL2UpdateMpcAddressTxLogic(ctx, k, txData)
		if err != nil {
			k.Logger(ctx).Error("L2UpdateMpcAddress failed", "id", msg.ID, "error", err)
			return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
		}
	default:
	}

	k.Logger(ctx).Info("Successfully validated External call for mpc sign msg", "id", msg.ID)

	result.Result = abci.SideTxResultType_Yes
	return
}

// SideHandleMsgMpcSign validates external calls required for processing mpc sign
func SideHandleMsgMpcSign(ctx sdk.Context, k Keeper, msg types.MsgMpcSign) (result abci.ResponseDeliverSideTx) {
	k.Logger(ctx).Info("Validating External call for mpc sign msg",
		"id", msg.ID,
	)

	// check if sign id exist.
	if !k.HasMpcSign(ctx, msg.ID) {
		k.Logger(ctx).Info("mpc sign id not exist, ignore it", "id", msg.ID)
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeMpcSignNotFound)
	}

	// check sign signature
	sign, err := k.GetMpcSign(ctx, msg.ID)
	if err != nil {
		k.Logger(ctx).Error("GetMpcSign err", "id", msg.ID, "err", err)
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeMpcSignNotFound)
	}
	if len(sign.Signature) > 0 {
		k.Logger(ctx).Error("GetMpcSign err", "id", msg.ID, "err", err)
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeMpcSignatureAlreadyExist)
	}

	// get mpc info
	mpcInfo, err := k.GetMpc(ctx, sign.MpcID)
	if err != nil {
		k.Logger(ctx).Error("GetMpcSign err", "id", msg.ID, "err", err)
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeMpcNotFound)
	}

	// verify signature
	pubKey, err := btcec.ParsePubKey(mpcInfo.MpcPubkey, btcec.S256())
	if err != nil {
		k.Logger(ctx).Error("ParsePubKey err", "id", msg.ID, "err", err)
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
	}
	publicKeyBytes := crypto.FromECDSAPub(pubKey.ToECDSA())

	sig, _ := hex.DecodeString(msg.Signature)
	if sig[64] > 1 {
		sig[64] -= 27
	}

	sigPublicKey, err := crypto.Ecrecover(sign.SignMsg, sig)
	if err != nil {
		k.Logger(ctx).Error("Ecrecover err", "id", msg.ID, "err", err)
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeMpcInvalidSignature)
	}

	if !bytes.Equal(sigPublicKey, publicKeyBytes) {
		k.Logger(ctx).Error("signature verify failed, pubkey mismatch", "id", msg.ID)
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeMpcInvalidSignature)
	}

	// get mpc set info
	mpcSetInfo := make(map[hmTypes.ThemisAddress]struct{})
	mpcSet := k.GetAllMpcSets(ctx)
	for _, set := range mpcSet {
		signer := hmTypes.HexToThemisAddress(set.Moniker)
		mpcSetInfo[signer] = struct{}{}
	}

	// check Proposer
	if _, exist := mpcSetInfo[msg.Proposer]; !exist {
		k.Logger(ctx).Error("invalid mpc proposer", "id", msg.ID, "proposer", msg.Proposer)
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
	}

	k.Logger(ctx).Info("Successfully validated External call for mpc sign msg", "id", msg.ID)

	result.Result = abci.SideTxResultType_Yes
	return
}

// PostHandleMsgEventMpcCreate handles state persisting mpc msg
func PostHandleMsgEventMpcCreate(ctx sdk.Context, k Keeper, msg types.MsgProposeMpcCreate, sideTxResult abci.SideTxResultType) sdk.Result {
	// Skip handler if mpc is not approved
	if sideTxResult != abci.SideTxResultType_Yes {
		k.Logger(ctx).Debug("Skipping new mpc since side-tx didn't get yes votes", "id", msg.ID)
		return hmCommon.ErrSideTxValidation(k.Codespace()).Result()
	}

	// check for replay
	if k.HasMpc(ctx, msg.ID) {
		k.Logger(ctx).Debug("Skipping new mpc as it's already processed", "id", msg.ID)
		return hmCommon.ErrMpcAlreadyExist(k.Codespace()).Result()
	}
	k.Logger(ctx).Debug("Persisting mpc state", "id", msg.ID, "sideTxResult", sideTxResult)

	newMpc := hmTypes.NewMpc(msg.ID, msg.Threshold, msg.Participants, msg.MpcAddress, msg.MpcPubkey, msg.MpcType)
	err := k.AddNewMpc(ctx, newMpc)
	if err != nil {
		k.Logger(ctx).Error("AddNewMpc err", "id", msg.ID, "err", err)
		return hmCommon.ErrInvalidMsg(k.Codespace(), err.Error()).Result()
	}

	// set latest mpc id
	k.UpdateLastMpc(ctx, msg.ID, msg.MpcType)

	// TX bytes
	txBytes := ctx.TxBytes()
	hash := tmTypes.Tx(txBytes).Hash()

	// add events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeProposeMpcCreate,
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),                                // action
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),              // module name
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmTypes.BytesToThemisHash(hash).Hex()), // tx hash
			sdk.NewAttribute(hmTypes.AttributeKeySideTxResult, sideTxResult.String()),           // result
			sdk.NewAttribute(types.AttributeKeyMpcID, msg.ID),
			sdk.NewAttribute(types.AttributeKeyMpcAddress, msg.MpcAddress.EthAddress().Hex()),
		),
	})

	// draft result with events
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// PostHandleMsgEventMpcSign handles state persisting mpc msg
func PostHandleMsgProposeMpcSign(ctx sdk.Context, k Keeper, msg types.MsgProposeMpcSign, sideTxResult abci.SideTxResultType) sdk.Result {
	// Skip handler if mpc is not approved
	if sideTxResult != abci.SideTxResultType_Yes {
		k.Logger(ctx).Debug("Skipping new re-mpc since side-tx didn't get yes votes", "id", msg.ID)
		return hmCommon.ErrSideTxValidation(k.Codespace()).Result()
	}

	// check for replay
	if k.HasMpcSign(ctx, msg.ID) {
		k.Logger(ctx).Info("mpc sign had exist", "id", msg.ID)
		return hmCommon.ErrMpcSignAlreadyExist(k.Codespace()).Result()
	}

	k.Logger(ctx).Debug("Persisting mpc state", "id", msg.ID, "sideTxResult", sideTxResult)

	// double check sign msg
	_, txSignMsg, _, _ := GetCheckTxSignMsgData(ctx, k, msg.SignType, msg.SignData)
	if msg.SignMsg != nil {
		txSignMsg = msg.SignMsg
	}

	// store sign info
	newMpcSign := hmTypes.NewMpcSign(
		msg.ID,
		msg.MpcID,
		msg.SignType,
		msg.SignData,
		txSignMsg,
		msg.Proposer,
	)

	err := k.AddNewMpcSign(ctx, newMpcSign)
	if err != nil {
		k.Logger(ctx).Error("AddNewMpcSign err", "id", msg.ID, "err", err)
		return hmCommon.ErrInvalidMsg(k.Codespace(), err.Error()).Result()
	}

	// TX bytes
	txBytes := ctx.TxBytes()
	hash := tmTypes.Tx(txBytes).Hash()

	// add events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeProposeMpcSign,
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),                                // action
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),              // module name
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmTypes.BytesToThemisHash(hash).Hex()), // tx hash
			sdk.NewAttribute(hmTypes.AttributeKeySideTxResult, sideTxResult.String()),           // result
			sdk.NewAttribute(types.AttributeKeyMpcSignID, msg.ID),
			sdk.NewAttribute(types.AttributeKeyMpcSignType, fmt.Sprintf("%v", msg.SignType)),
			sdk.NewAttribute(types.AttributeKeyMpcSignMsg, hex.EncodeToString(txSignMsg)),
			sdk.NewAttribute(types.AttributeKeyMpcID, msg.MpcID),
		),
	})

	// draft result with events
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// PostHandleMsgMpcSign handles state persisting mpc msg
func PostHandleMsgMpcSign(ctx sdk.Context, k Keeper, msg types.MsgMpcSign, sideTxResult abci.SideTxResultType) sdk.Result {
	// Skip handler if mpc is not approved
	if sideTxResult != abci.SideTxResultType_Yes {
		k.Logger(ctx).Debug("Skipping since side-tx didn't get yes votes", "id", msg.ID)
		return hmCommon.ErrSideTxValidation(k.Codespace()).Result()
	}

	// check for replay
	if !k.HasMpcSign(ctx, msg.ID) {
		k.Logger(ctx).Info("mpc sign not exist")
		return hmCommon.ErrMpcSignNotFound(k.Codespace()).Result()
	}

	// check sign signature
	sign, err := k.GetMpcSign(ctx, msg.ID)
	if err != nil {
		k.Logger(ctx).Error("GetMpcSign err", "id", msg.ID, "err", err)
		return hmCommon.ErrInvalidMsg(k.Codespace(), err.Error()).Result()
	}
	if len(sign.Signature) > 0 {
		k.Logger(ctx).Error("GetMpcSign err", "id", msg.ID, "err", err)
		return hmCommon.ErrMpcAlreadyHadSignature(k.Codespace()).Result()
	}

	// store sign info
	existMpcSign := hmTypes.NewMpcSign(
		sign.ID,
		sign.MpcID,
		sign.SignType,
		sign.SignData,
		sign.SignMsg,
		sign.Proposer,
	)

	sig, err := hex.DecodeString(strings.TrimPrefix(msg.Signature, "0x"))
	if err != nil {
		k.Logger(ctx).Error("DecodeString err", "id", msg.ID, "err", err)
		return hmCommon.ErrInvalidMsg(k.Codespace(), err.Error()).Result()
	}
	if sig[64] > 1 {
		sig[64] -= 27
	}

	existMpcSign.Signature = sig

	// BatchReward sign type does not need to generate signed tx

	txSigner, _, _, err := GetCheckTxSignMsgData(ctx, k, sign.SignType, sign.SignData)
	if err != nil {
		k.Logger(ctx).Error("GetCheckTxSignMsgData err", "id", msg.ID, "err", err)
		return hmCommon.ErrInvalidMsg(k.Codespace(), err.Error()).Result()
	}
	k.Logger(ctx).Debug("before build tx signer with chainID", "id", msg.ID, "signerChainID", txSigner.ChainID())

	tx, err := types.DecodeUnsignedTx(sign.SignData)
	if err != nil {
		k.Logger(ctx).Error("PostHandleMsgMpcSign tx unmarlshal failed", "id", msg.ID, "err", err)
		return hmCommon.ErrInvalidMsg(k.Codespace(), err.Error()).Result()
	}
	signedTx, err := tx.WithSignature(txSigner, sig)
	if err != nil {
		k.Logger(ctx).Error("PostHandleMsgMpcSign WithSignature", "id", msg.ID, "error", err)
		return hmCommon.ErrInvalidMsg(k.Codespace(), err.Error()).Result()
	}
	k.Logger(ctx).Debug("after build tx signer with chainID", "id", msg.ID, "signerChainID", txSigner.ChainID())
	k.Logger(ctx).Debug("after build tx signer with chainID", "id", msg.ID, "signedTxChainID", signedTx.ChainId())
	existMpcSign.SignedTx, _ = signedTx.MarshalBinary()
	k.Logger(ctx).Debug("signedTx build", "id", msg.ID, "signedTx", hex.EncodeToString(existMpcSign.SignedTx))

	// update signature
	err = k.UpdateMpcSign(ctx, existMpcSign)
	if err != nil {
		k.Logger(ctx).Error("UpdateMpcSign err", "id", msg.ID, "err", err)
		return hmCommon.ErrInvalidMsg(k.Codespace(), err.Error()).Result()
	}

	// TX bytes
	txBytes := ctx.TxBytes()
	hash := tmTypes.Tx(txBytes).Hash()

	// add events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMpcSign,
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),                                // action
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),              // module name
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmTypes.BytesToThemisHash(hash).Hex()), // tx hash
			sdk.NewAttribute(hmTypes.AttributeKeySideTxResult, sideTxResult.String()),           // result
			sdk.NewAttribute(types.AttributeKeyMpcSignID, msg.ID),
			sdk.NewAttribute(types.AttributeKeyMpcSignature, msg.Signature),
			sdk.NewAttribute(types.AttributeKeyMpcSignType, fmt.Sprintf("%v", sign.SignType)),
		),
	})

	// draft result with events
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
