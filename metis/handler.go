package metis

import (
	"bytes"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/metis-seq/themis/common"
	"github.com/metis-seq/themis/helper"
	"github.com/metis-seq/themis/metis/types"
	hmTypes "github.com/metis-seq/themis/types"
)

// NewHandler returns a handler for "metis" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgProposeSpan:
			return HandleMsgProposeSpan(ctx, msg, k)
		case types.MsgReProposeSpan:
			return HandleMsgReProposeSpan(ctx, msg, k)
		case types.MsgMetisTx:
			return HandleMsgMetisTx(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("Invalid message in metis module").Result()
		}
	}
}

// HandleMsgProposeSpan handles proposeSpan msg
func HandleMsgProposeSpan(ctx sdk.Context, msg types.MsgProposeSpan, k Keeper) sdk.Result {
	k.Logger(ctx).Info("✅ Validating proposed span msg",
		"spanId", msg.ID,
		"currentL2Height", msg.CurrentL2Height,
		"startBlock", msg.StartBlock,
		"endBlock", msg.EndBlock,
		"seed", msg.Seed.String(),
		"isRecover", msg.IsRecover,
	)

	// Check if the proposer is allowed
	if !k.IsSpanValidator(ctx, msg.Proposer) {
		k.Logger(ctx).Error("Invalid msg proposer", "proposer", msg.Proposer)
		return common.ErrNoValidator(k.Codespace()).Result()
	}

	// check for replay
	if k.HasSpan(ctx, msg.ID) {
		k.Logger(ctx).Info("Skipping new span as it's already processed")
		return common.ErrSpanHadExist(k.Codespace()).Result()
	}

	// chainManager params
	params := k.chainKeeper.GetParams(ctx)
	chainParams := params.ChainParams

	// check chain id
	if chainParams.MetisChainID != msg.ChainID {
		k.Logger(ctx).Error("Invalid Metis chain id", "msgChainID", msg.ChainID)
		return common.ErrInvalidMetisChainID(k.Codespace()).Result()
	}

	// check if last span is up or if greater diff than threshold is found between validator set
	lastSpan, err := k.GetLastSpan(ctx)
	if err != nil {
		k.Logger(ctx).Error("Unable to fetch last span", "Error", err)
		return common.ErrSpanNotFound(k.Codespace()).Result()
	}

	// Validate span continuity
	if (lastSpan.ID+1 != msg.ID || msg.StartBlock != lastSpan.EndBlock+1 || msg.EndBlock < msg.StartBlock) && !msg.IsRecover {
		k.Logger(ctx).Error("Blocks not in continuity",
			"lastSpanId", lastSpan.ID,
			"spanId", msg.ID,
			"lastSpanStartBlock", lastSpan.StartBlock,
			"lastSpanEndBlock", lastSpan.EndBlock,
			"spanStartBlock", msg.StartBlock,
			"spanEndBlock", msg.EndBlock,
		)

		return common.ErrSpanNotInContinuity(k.Codespace()).Result()
	}

	// calculate next span seed locally,check if span seed matches or not.
	nextSpanSeed, _ := k.GetNextSpanSeed(ctx)
	if !bytes.Equal(msg.Seed.Bytes(), nextSpanSeed.Bytes()) {
		k.Logger(ctx).Error(
			"Span Seed does not match",
			"msgSeed", msg.Seed.String(),
			"nextSpanSeed", nextSpanSeed.String(),
		)

		return common.ErrInvalidSeed(k.Codespace()).Result()
	}

	k.Logger(ctx).Debug("Persisting span state")

	// freeze for new span
	if msg.IsRecover {
		err := k.FreezeFixedSet(ctx, msg.ID, msg.StartBlock, msg.EndBlock, msg.ChainID, msg.RecoverSigner)
		if err != nil {
			k.Logger(ctx).Error("Unable to freeze validator set for span", "Error", err)
			return common.ErrUnableToFreezeValSet(k.Codespace()).Result()
		}
	} else {
		// save span new proposer
		err := k.SaveSpanProposer(ctx, msg)
		if err != nil {
			k.Logger(ctx).Error("HandleMsgProposeSpan SaveSpanProposer failed", "err", err)
			return common.ErrSideTxValidation(k.Codespace()).Result()
		}
		k.Logger(ctx).Debug("HandleMsgProposeSpan SaveSpanProposer success")

		// get all exist proposers
		allProposers, err := k.GetSpanAllProposers(ctx, msg)
		if err != nil {
			k.Logger(ctx).Error("HandleMsgProposeSpan GetReSpanAllProposers failed", "err", err)
			return common.ErrSideTxValidation(k.Codespace()).Result()
		}
		k.Logger(ctx).Debug("HandleMsgProposeSpan GetReSpanAllProposers", "allProposers", allProposers)

		needVotes := float64(len(lastSpan.ValidatorSet.Validators)*2) / 3
		allProposersCount := float64(len(allProposers))

		k.Logger(ctx).Info("HandleMsgProposeSpan Check span votes", "needVotes", needVotes, "currentVotes", allProposersCount)
		// check re-span votes
		if needVotes > allProposersCount {
			k.Logger(ctx).Error("HandleMsgProposeSpan Check span votes not enough")
			return sdk.Result{
				Events: ctx.EventManager().Events(),
			}
		}
		k.Logger(ctx).Error("HandleMsgProposeSpan Check span votes enough")

		// recheck for replay
		if k.HasSpan(ctx, msg.ID) {
			k.Logger(ctx).Info("Skipping new span as it's already processed")
			return common.ErrSpanHadExist(k.Codespace()).Result()
		}

		err = k.FreezeSet(ctx, msg.ID, msg.StartBlock, msg.EndBlock, msg.ChainID, msg.Seed, msg.BlackList)
		if err != nil {
			k.Logger(ctx).Error("Unable to freeze validator set for span", "Error", err)
			return common.ErrUnableToFreezeValSet(k.Codespace()).Result()
		}
	}

	isRecoverStr := "false"
	if msg.IsRecover {
		isRecoverStr = "true"
	}

	// TX bytes
	txBytes := ctx.TxBytes()
	hash := tmTypes.Tx(txBytes).Hash()

	// add events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeProposeSpan,
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),                                // action
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),              // module name
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmTypes.BytesToThemisHash(hash).Hex()), // tx hash
			sdk.NewAttribute(types.AttributeKeySpanID, strconv.FormatUint(msg.ID, 10)),
			sdk.NewAttribute(types.AttributeKeySpanSigner, msg.Proposer.EthAddress().Hex()),
			sdk.NewAttribute(types.AttributeKeySpanRecover, isRecoverStr),
		),
	})

	// draft result with events
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// HandleMsgReProposeSpan handles proposeSpan msg
func HandleMsgReProposeSpan(ctx sdk.Context, msg types.MsgReProposeSpan, k Keeper) sdk.Result {
	k.Logger(ctx).Info("✅ Validating proposed re-span msg",
		"spanId", msg.ID,
		"startBlock", msg.StartBlock,
		"endBlock", msg.EndBlock,
	)

	// Check if the proposer is allowed
	if !k.IsSpanValidator(ctx, msg.Proposer) {
		k.Logger(ctx).Error("Invalid msg proposer", "proposer", msg.Proposer)
		return common.ErrNoValidator(k.Codespace()).Result()
	}

	// check oldSpan and newSpan ID
	if msg.ID != msg.CurrentL2Epoch+1 {
		k.Logger(ctx).Error("Check span not found")
		return common.ErrSpanNotFound(k.Codespace()).Result()
	}

	oldSpanID, newSpanID := msg.CurrentL2Epoch, msg.CurrentL2Epoch+1
	// check old span
	oldSpan, err := k.GetSpan(ctx, oldSpanID)
	if err != nil {
		k.Logger(ctx).Error("HandleMsgReProposeSpan Check old span not found")
		return common.ErrSpanNotFound(k.Codespace()).Result()
	}
	if msg.StartBlock-1 < oldSpan.StartBlock {
		k.Logger(ctx).Error("HandleMsgReProposeSpan invalid msg start block")
		return common.ErrSpanNotFound(k.Codespace()).Result()
	}

	// check re-span msg status
	if k.HasReSpanFinish(ctx, msg) {
		k.Logger(ctx).Info("Skipping new respan as it's already processed")
		return common.ErrOldTx(k.Codespace()).Result()
	}

	// chainManager params
	params := k.chainKeeper.GetParams(ctx)
	chainParams := params.ChainParams

	// check chain id
	if chainParams.MetisChainID != msg.ChainID {
		k.Logger(ctx).Error("Invalid Metis chain id", "msgChainID", msg.ChainID)
		return common.ErrInvalidMetisChainID(k.Codespace()).Result()
	}

	// check newSequencer
	currentBatch := k.sk.GetL1Batch(ctx)
	valSet := k.sk.GetValidatorSet(ctx)

	newSequencer := helper.CalcMetisSequencerWithSeed(msg.CurrentProducer.EthAddress().Hex(), msg.CurrentL2Height, currentBatch, &valSet, msg.Seed)
	if !strings.EqualFold(newSequencer, msg.NextProducer.EthAddress().Hex()) {
		k.Logger(ctx).Error("Producer mismatch", "newSequencer", newSequencer, "NextProducer", msg.NextProducer)
		return common.ErrInvalidNewSequencer(k.Codespace()).Result()
	}

	if msg.CurrentL2Height != msg.StartBlock-1 {
		k.Logger(ctx).Error("block mismatch", "CurrentL2Height", msg.CurrentL2Height, "msgBlock", msg.StartBlock-1)
		return common.ErrInvalidBlock(k.Codespace()).Result()
	}

	// endBlock must large than startBlock
	if msg.EndBlock <= msg.StartBlock {
		k.Logger(ctx).Error("invalid end block", "startBlock", msg.StartBlock, "endBlock", msg.EndBlock)
		return common.ErrInvalidBlock(k.Codespace()).Result()
	}

	// Not in the scope of watching
	if msg.CurrentL2Height < msg.StartBlock-1 || msg.CurrentL2Height > msg.EndBlock {
		k.Logger(ctx).Error("Error child check", "spanId", msg.ID, "CurrentL2Height", msg.CurrentL2Height, "newSpanStartBlock", msg.StartBlock, "msgEndBlock", msg.EndBlock)
		return common.ErrInvalidBlock(k.Codespace()).Result()
	}

	// save re-span new proposer
	err = k.SaveReSpanProposer(ctx, msg)
	if err != nil {
		k.Logger(ctx).Error("HandleMsgReProposeSpan SaveReSpanProposer failed", "err", err)
		return common.ErrSideTxValidation(k.Codespace()).Result()
	}
	k.Logger(ctx).Debug("HandleMsgReProposeSpan SaveReSpanProposer success")

	// get all exist proposers
	allProposers, err := k.GetReSpanAllProposers(ctx, msg)
	if err != nil {
		k.Logger(ctx).Error("HandleMsgReProposeSpan GetReSpanAllProposers failed", "err", err)
		return common.ErrSideTxValidation(k.Codespace()).Result()
	}
	k.Logger(ctx).Debug("HandleMsgReProposeSpan GetReSpanAllProposers", "allProposers", allProposers)

	allValidators := k.sk.GetCurrentValidators(ctx)
	needVotes := float64(len(allValidators)*2) / 3
	allProposersCount := float64(len(allProposers))
	k.Logger(ctx).Info("HandleMsgReProposeSpan Check re-span votes", "needVotes", needVotes, "currentVotes", allProposersCount)
	// check re-span votes
	if needVotes > allProposersCount {
		k.Logger(ctx).Error("HandleMsgReProposeSpan Check re-span votes not enough")
		return sdk.Result{
			Events: ctx.EventManager().Events(),
		}
	}
	k.Logger(ctx).Info("HandleMsgReProposeSpan Check re-span votes enough")

	// recheck re-span msg status
	if k.HasReSpanFinish(ctx, msg) {
		k.Logger(ctx).Info("Skipping new respan as it's already processed")
		return common.ErrOldTx(k.Codespace()).Result()
	}

	// check latest span
	latestSpan, err := k.GetLastSpan(ctx)
	if err != nil {
		k.Logger(ctx).Error("HandleMsgReProposeSpan GetLastSpan", "err", err)
		return common.ErrSpanNotFound(k.Codespace()).Result()
	}
	k.Logger(ctx).Debug("HandleMsgReProposeSpan check latest span id", "latestSpan", latestSpan.ID, "reSpanNewSpan", newSpanID)
	if latestSpan.ID > newSpanID {
		for {
			latestSpan, _ := k.GetLastSpan(ctx)
			if latestSpan.ID <= newSpanID {
				break
			}
			k.Logger(ctx).Info("HandleMsgReProposeSpan DelSpan", "id", latestSpan.ID)
			k.DelSpan(ctx, latestSpan.ID)
		}
	}

	// update span end block
	err = k.UpdateFixedSpanEndBlock(ctx, oldSpanID, msg.StartBlock-1)
	if err != nil {
		k.Logger(ctx).Error("HandleMsgReProposeSpan Unable to update last span for last span", "Error", err)
		return common.ErrUpdateSpan(k.Codespace()).Result()
	}

	// generate new span or update last span
	err = k.FreezeFixedSet(ctx, newSpanID, msg.StartBlock, msg.EndBlock, msg.ChainID, msg.NextProducer)
	if err != nil {
		k.Logger(ctx).Error("HandleMsgReProposeSpan Unable to freeze validator set for span", "Error", err)
		return common.ErrUnableToFreezeValSet(k.Codespace()).Result()
	}

	// re-check replay
	if k.HasReSpanFinish(ctx, msg) {
		k.Logger(ctx).Error("HandleMsgReProposeSpan Skipping new respan as it's already processed")
		return common.ErrOldTx(k.Codespace()).Result()
	}
	// set respan record
	k.SetReSpanFinish(ctx, msg)
	k.SetReSpanTime(ctx, msg.ID, uint64(ctx.BlockTime().Unix()))
	k.Logger(ctx).Info("HandleMsgReProposeSpan success", "msg", msg)

	// TX bytes
	txBytes := ctx.TxBytes()
	hash := tmTypes.Tx(txBytes).Hash()

	// add events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeReProposeSpan,
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),                                // action
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),              // module name
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmTypes.BytesToThemisHash(hash).Hex()), // tx hash
			sdk.NewAttribute(types.AttributeKeyOldSpanID, strconv.FormatUint(oldSpanID, 10)),
			sdk.NewAttribute(types.AttributeKeySpanID, strconv.FormatUint(newSpanID, 10)),
			sdk.NewAttribute(types.AttributeKeySpanStartBlock, strconv.FormatUint(msg.StartBlock, 10)),
			sdk.NewAttribute(types.AttributeKeySpanEndBlock, strconv.FormatUint(msg.EndBlock, 10)),
			sdk.NewAttribute(types.AttributeKeySpanSigner, msg.Proposer.EthAddress().Hex()),
		),
	})

	// draft result with events
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// HandleMsgMetisTx handles delReSpan msg
func HandleMsgMetisTx(ctx sdk.Context, msg types.MsgMetisTx, k Keeper) sdk.Result {
	k.Logger(ctx).Info("✅ Validating metis tx msg",
		"txHash", msg.TxHash,
		"logIndex", msg.LogIndex,
	)

	// Check if the proposer is allowed
	if !k.IsSpanValidator(ctx, msg.Proposer) {
		k.Logger(ctx).Error("Invalid msg proposer", "proposer", msg.Proposer)
		return common.ErrNoValidator(k.Codespace()).Result()
	}

	// TX bytes
	txBytes := ctx.TxBytes()
	hash := tmTypes.Tx(txBytes).Hash()

	// add events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMetisTx,
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),                                // action
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),              // module name
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmTypes.BytesToThemisHash(hash).Hex()), // tx hash
			sdk.NewAttribute(types.AttributeKeyMetisTxHash, msg.TxHash.Hex()),
			sdk.NewAttribute(types.AttributeKeyMetisTxData, msg.TxData),
		),
	})

	// draft result with events
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
