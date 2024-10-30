package mpc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/metis-seq/themis/common"
	"github.com/metis-seq/themis/mpc/types"
	hmTypes "github.com/metis-seq/themis/types"
)

// NewHandler returns a handler for "mpc" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgProposeMpcCreate:
			return HandleMsgProposeMpcCreate(ctx, msg, k)
		case types.MsgProposeMpcSign:
			return HandleMsgProposeMpcSign(ctx, msg, k)
		case types.MsgMpcSign:
			return HandleMsgMpcSign(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("Invalid message in mpc module").Result()
		}
	}
}

// HandleMsgProposeMpcCreate handles proposeSpan msg
func HandleMsgProposeMpcCreate(ctx sdk.Context, msg types.MsgProposeMpcCreate, k Keeper) sdk.Result {
	k.Logger(ctx).Info("Validating proposed mpc create msg",
		"id", msg.ID,
	)

	// check if mpc id exist.
	if k.HasMpc(ctx, msg.ID) {
		k.Logger(ctx).Info("mpc id already exist, ignore it")
		return common.ErrMpcAlreadyExist(k.Codespace()).Result()
	}

	// check mpc type
	switch msg.MpcType {
	case hmTypes.CommonMpcType, hmTypes.StateSubmitMpcType, hmTypes.RewardSubmitMpcType, hmTypes.BlobSubmitMpcType:
	default:
		k.Logger(ctx).Info("invalid mpc type, ignore it")
		return common.ErrMpcInvalidType(k.Codespace()).Result()
	}

	// add events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeProposeMpcCreate,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyMpcID, msg.ID),
			sdk.NewAttribute(types.AttributeKeyMpcAddress, msg.MpcAddress.EthAddress().Hex()),
		),
	})

	// draft result with events
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// HandleMsgProposeMpcSign handles mpc msg
func HandleMsgProposeMpcSign(ctx sdk.Context, msg types.MsgProposeMpcSign, k Keeper) sdk.Result {
	k.Logger(ctx).Info("Validating proposed mpc sign msg",
		"id", msg.ID,
	)

	// check if mpc sign id exist.
	if k.HasMpcSign(ctx, msg.ID) {
		k.Logger(ctx).Info("mpc sign id already exist, ignore it")
		return common.ErrMpcSignAlreadyExist(k.Codespace()).Result()
	}

	// add events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeProposeMpcSign,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyMpcSignID, msg.ID),
		),
	})

	// draft result with events
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// HandleMsgMpcSign handles mpc msg
func HandleMsgMpcSign(ctx sdk.Context, msg types.MsgMpcSign, k Keeper) sdk.Result {
	k.Logger(ctx).Info("Validating proposed mpc sign msg",
		"id", msg.ID,
	)

	// check if mpc sign id exist.
	if !k.HasMpcSign(ctx, msg.ID) {
		k.Logger(ctx).Info("mpc sign not exist, ignore it")
		return common.ErrMpcSignNotFound(k.Codespace()).Result()
	}

	// add events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMpcSign,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyMpcSignID, msg.ID),
		),
	})

	// draft result with events
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
