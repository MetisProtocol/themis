package common

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type CodeType = sdk.CodeType

const (
	DefaultCodespace sdk.CodespaceType = "1"

	CodeInvalidMsg CodeType = 1400
	CodeOldTx      CodeType = 1401

	CodeInvalidProposerInput CodeType = 1500
	CodeInvalidBlockInput    CodeType = 1501
	CodeInvalidACK           CodeType = 1502
	CodeNoACK                CodeType = 1503
	CodeBadTimeStamp         CodeType = 1504
	CodeInvalidNoACK         CodeType = 1505
	CodeTooManyNoAck         CodeType = 1506
	CodeLowBal               CodeType = 1507

	CodeOldValidator        CodeType = 2500
	CodeNoValidator         CodeType = 2501
	CodeValSignerMismatch   CodeType = 2502
	CodeValidatorExitDeny   CodeType = 2503
	CodeValAlreadyUnbonded  CodeType = 2504
	CodeSignerSynced        CodeType = 2505
	CodeValSave             CodeType = 2506
	CodeValAlreadyJoined    CodeType = 2507
	CodeSignerUpdateError   CodeType = 2508
	CodeNoConn              CodeType = 2509
	CodeWaitFrConfirmation  CodeType = 2510
	CodeValPubkeyMismatch   CodeType = 2511
	CodeErrDecodeEvent      CodeType = 2512
	CodeNoSignerChangeError CodeType = 2513
	CodeNonce               CodeType = 2514
	CodeInvalidBatchID      CodeType = 2515

	CodeSpanNotCountinuous     CodeType = 3501
	CodeUnableToFreezeSet      CodeType = 3502
	CodeSpanNotFound           CodeType = 3503
	CodeValSetMisMatch         CodeType = 3504
	CodeProducerMisMatch       CodeType = 3505
	CodeInvalidMetisChainID    CodeType = 3506
	CodeInvalidSpanDuration    CodeType = 3507
	CodeUpdateSpanFailed       CodeType = 3508
	CodeInvalidSpanID          CodeType = 3509
	CodeErrReSpanOldID         CodeType = 3510
	CodeErrReSpanVoteNotEnough CodeType = 3511
	CodeMetisTxExist           CodeType = 3512
	CodeInvalidSeed            CodeType = 3513
	CodeSpanHadExit            CodeType = 3514
	CodeInvalidNewSequencer    CodeType = 3516

	CodeErrComputeGenesisAccountRoot CodeType = 4503
	CodeAccountRootMismatch          CodeType = 4504
	CodeErrAccountRootHash           CodeType = 4505

	CodeInvalidReceipt         CodeType = 5501
	CodeSideTxValidationFailed CodeType = 5502

	CodeValSigningInfoSave     CodeType = 6501
	CodeErrValUnjail           CodeType = 6502
	CodeSlashInfoDetails       CodeType = 6503
	CodeTickNotInContinuity    CodeType = 6504
	CodeTickAckNotInContinuity CodeType = 6505

	CodeMpcNotFound              CodeType = 7501
	CodeMpcExist                 CodeType = 7502
	CodeUpdateMpcFailed          CodeType = 7503
	CodeMpcSignNotFound          CodeType = 7504
	CodeMpcSignExist             CodeType = 7505
	CodeUpdateMpcSignFailed      CodeType = 7506
	CodeMpcAlreadyExist          CodeType = 7507
	CodeMpcSignAlreadyExist      CodeType = 7508
	CodeMpcSignFailed            CodeType = 7509
	CodeMpcSignatureAlreadyExist CodeType = 7510
	CodeMpcInvalidSignature      CodeType = 7511
	CodeMpcInvalidType           CodeType = 7512
)

// -------- Invalid msg

func ErrInvalidMsg(codespace sdk.CodespaceType, format string, args ...interface{}) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidMsg, format, args...)
}

// ----------- Staking Errors

func ErrOldValidator(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeOldValidator, "Start Batch behind Current Batch")
}

func ErrNoValidator(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeNoValidator, "Validator information not found")
}

func ErrNonce(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeNonce, "Incorrect validator nonce")
}

func ErrValSignerPubKeyMismatch(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeValPubkeyMismatch, "Signer Pubkey mismatch between event and msg")
}

func ErrValSignerMismatch(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeValSignerMismatch, "Signer Address doesnt match pubkey address")
}

func ErrValIsNotCurrentVal(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeValidatorExitDeny, "Validator is not in validator set, exit not possible")
}

func ErrValUnbonded(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeValAlreadyUnbonded, "Validator already unbonded , cannot exit")
}

func ErrSignerUpdateError(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeSignerUpdateError, "Signer update error")
}

func ErrNoSignerChange(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeNoSignerChangeError, "New signer same as old signer")
}

func ErrOldTx(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeOldTx, "Old txhash not allowed")
}

func ErrValidatorAlreadySynced(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeSignerSynced, "No signer update found, invalid message")
}

func ErrValidatorSave(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeValSave, "Cannot save validator")
}

func ErrValidatorNotDeactivated(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeValSave, "Validator Not Deactivated")
}

func ErrValidatorAlreadyJoined(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeValAlreadyJoined, "Validator already joined")
}

func ErrInvalidBatchID(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeInvalidBatchID, "Invalid batch id")
}

// Metis Errors --------------------------------

func ErrInvalidMetisChainID(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeInvalidMetisChainID, "Invalid Metis chain id")
}

func ErrSpanNotInContinuity(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeSpanNotCountinuous, "Span not countinuous")
}

func ErrInvalidBlock(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeInvalidBlockInput, "Invalid start/end block")
}

func ErrInvalidSpanDuration(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeInvalidSpanDuration, "wrong span duration")
}

func ErrInvalidSpanID(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeInvalidSpanID, "InvaliSpanId")
}

func ErrInvalidSeed(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeInvalidSeed, "Invalid seed")
}

func ErrSpanNotFound(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeSpanNotFound, "Span not found")
}

func ErrSpanHadExist(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeSpanHadExit, "Span had exist")
}

func ErrReSpanOldID(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeErrReSpanOldID, "oldSpanID must less than or equal L2 epochID")
}

func ErrUpdateSpan(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeUpdateSpanFailed, "Update last span failed")
}

func ErrUnableToFreezeValSet(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeUnableToFreezeSet, "Unable to freeze validator set for next span")
}

func ErrValSetMisMatch(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeValSetMisMatch, "Validator set mismatch")
}

func ErrProducerMisMatch(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeProducerMisMatch, "Producer set mismatch")
}

func ErrReSpanVoteNotEnough(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeErrReSpanVoteNotEnough, "ReSpan vote not enough")
}

func ErrMetisTxExist(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeMetisTxExist, "Metis tx exist")
}

func ErrInvalidNewSequencer(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeInvalidNewSequencer, "Invalid sequencer")
}

// Mpc Errors --------------------------------

func ErrMpcNotFound(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeMpcNotFound, "Mpc not found")
}

func ErrMpcAlreadyExist(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeMpcAlreadyExist, "Mpc already exist")
}

func ErrMpcInvalidType(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeMpcInvalidType, "Invalid mpc type")
}

func ErrUpdateMpc(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeUpdateMpcFailed, "Update last span failed")
}

func ErrMpcSignNotFound(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeMpcSignNotFound, "MpcSign not found")
}

func ErrUpdateMpcSign(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeUpdateMpcSignFailed, "Update last span failed")
}

func ErrMpcSignAlreadyExist(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeMpcSignAlreadyExist, "MpcSign already exist")
}

func ErrMpcSignFailed(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeMpcSignAlreadyExist, "MpcSign failed")
}

func ErrMpcAlreadyHadSignature(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeMpcSignatureAlreadyExist, "Mpc signature already exist")
}

//
// Side-tx errors
//

// ErrorSideTx represents side-tx error
func ErrorSideTx(codespace sdk.CodespaceType, code CodeType) (res abci.ResponseDeliverSideTx) {
	res.Code = uint32(code)
	res.Codespace = string(codespace)
	res.Result = abci.SideTxResultType_Skip // skip side-tx vote in-case of error

	return
}

func ErrSideTxValidation(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeSideTxValidationFailed, "External call majority validation failed. ")
}

//
// Private methods
//

func CodeToDefaultMsg(code CodeType) string {
	switch code {
	// case CodeInvalidBlockInput:
	// 	return "Invalid Block Input"
	case CodeInvalidMsg:
		return "Invalid Message"
	case CodeInvalidProposerInput:
		return "Proposer is not valid"
	case CodeInvalidBlockInput:
		return "Wrong roothash for given start and end block numbers"
	case CodeInvalidACK:
		return "Ack Not Valid"
	case CodeBadTimeStamp:
		return "Invalid time stamp. It must be in near past."
	case CodeTooManyNoAck:
		return "Too many no-acks"
	case CodeLowBal:
		return "Insufficient balance"

	case CodeOldValidator:
		return "Start Batch behind Current Batch"
	case CodeNoValidator:
		return "Validator information not found"
	case CodeValSignerMismatch:
		return "Signer Address doesnt match pubkey address"
	case CodeValidatorExitDeny:
		return "Validator is not in validator set, exit not possible"
	case CodeValAlreadyUnbonded:
		return "Validator already unbonded , cannot exit"
	case CodeSignerSynced:
		return "No signer update found, invalid message"
	case CodeValSave:
		return "Cannot save validator"
	case CodeValAlreadyJoined:
		return "Validator already joined"
	case CodeSignerUpdateError:
		return "Signer update error"
	case CodeNoConn:
		return "Unable to connect to chain"
	case CodeWaitFrConfirmation:
		return "wait for confirmation time before sending transaction"
	case CodeValPubkeyMismatch:
		return "Signer Pubkey mismatch between event and msg"
	case CodeSpanNotCountinuous:
		return "Span not countinuous"
	case CodeUnableToFreezeSet:
		return "Unable to freeze validator set for next span"
	case CodeSpanNotFound:
		return "Span not found"
	case CodeValSetMisMatch:
		return "Validator set mismatch"
	case CodeProducerMisMatch:
		return "Producer set mismatch"
	case CodeInvalidMetisChainID:
		return "Invalid Metis chain id"
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

func msgOrDefaultMsg(msg string, code CodeType) string {
	if msg != "" {
		return msg
	}

	return CodeToDefaultMsg(code)
}

func newError(codespace sdk.CodespaceType, code CodeType, msg string) sdk.Error {
	msg = msgOrDefaultMsg(msg, code)
	return sdk.NewError(codespace, code, msg)
}

// Slashing errors
func ErrValidatorSigningInfoSave(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeValSigningInfoSave, "Cannot save validator signing info")
}

func ErrUnjailValidator(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeErrValUnjail, "Error while unJail validator")
}

func ErrSlashInfoDetails(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeSlashInfoDetails, "Wrong slash info details")
}

func ErrTickNotInContinuity(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeTickNotInContinuity, "Tick not in continuity")
}

func ErrTickAckNotInContinuity(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeTickAckNotInContinuity, "Tick-ack not in continuity")
}
