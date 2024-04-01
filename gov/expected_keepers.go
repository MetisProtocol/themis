package gov

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	supplyTypes "github.com/metis-seq/themis/supply/types"
	hmTypes "github.com/metis-seq/themis/types"
)

// SupplyKeeper defines the supply Keeper for module accounts
type SupplyKeeper interface {
	GetModuleAddress(name string) hmTypes.ThemisAddress
	GetModuleAccount(ctx sdk.Context, name string) supplyTypes.ModuleAccountInterface

	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	SetModuleAccount(sdk.Context, supplyTypes.ModuleAccountInterface)

	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr hmTypes.ThemisAddress, amt sdk.Coins) sdk.Error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr hmTypes.ThemisAddress, recipientModule string, amt sdk.Coins) sdk.Error
}
