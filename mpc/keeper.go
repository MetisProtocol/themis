package mpc

import (
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/metis-seq/themis/helper"
	"github.com/metis-seq/themis/mpc/types"
	"github.com/metis-seq/themis/params/subspace"
	"github.com/metis-seq/themis/staking"
	hmTypes "github.com/metis-seq/themis/types"
)

var (
	LastMpcIDKey     = []byte{0x45} // Key to store last mpc id
	MpcPrefixKey     = []byte{0x46} // prefix key to store mpc
	MpcSetPrefixKey  = []byte{0x47} // prefix key to store mpc set
	MpcSignPrefixKey = []byte{0x48} // prefix key to store mpc sign
)

// Keeper stores all related data
type Keeper struct {
	cdc *codec.Codec
	sk  staking.Keeper
	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey
	// codespace
	codespace sdk.CodespaceType
	// param space
	paramSpace subspace.Subspace
	// contract caller
	contractCaller helper.ContractCaller

	cm ModuleCommunicator
}

// NewKeeper is the constructor of Keeper
func NewKeeper(
	cdc *codec.Codec,
	storeKey sdk.StoreKey,
	paramSpace subspace.Subspace,
	codespace sdk.CodespaceType,
	caller helper.ContractCaller,
	stakingKeeper staking.Keeper,
	cm ModuleCommunicator,
) Keeper {
	return Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		paramSpace:     paramSpace,
		codespace:      codespace,
		contractCaller: caller,
		sk:             stakingKeeper,
		cm:             cm,
	}
}

// Codespace returns the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ModuleName)
}

// GetMpcKey appends prefix to start block
func GetMpcKey(id string) []byte {
	return append(MpcPrefixKey, []byte(id)...)
}

// GetMpcSetKey appends prefix
func GetMpcSetKey(mpcPartyID string) []byte {
	return append(MpcSetPrefixKey, []byte(mpcPartyID)...)
}

// GetMpcSignKey appends prefix
func GetMpcSignKey(id string) []byte {
	return append(MpcSignPrefixKey, []byte(id)...)
}

// AddNewMpcSet adds new mpc to store
func (k *Keeper) AddNewMpcSet(ctx sdk.Context, party hmTypes.PartyID) error {
	store := ctx.KVStore(k.storeKey)

	out, err := k.cdc.MarshalBinaryBare(party)
	if err != nil {
		k.Logger(ctx).Error("Error marshalling party", "error", err)
		return err
	}

	// store set party id
	store.Set(GetMpcSetKey(party.ID), out)

	return nil
}

// AddNewRawMpc adds new mpc set for metis to store
func (k *Keeper) AddNewRawMpcSet(ctx sdk.Context, party hmTypes.PartyID) error {
	store := ctx.KVStore(k.storeKey)

	out, err := k.cdc.MarshalBinaryBare(party)
	if err != nil {
		k.Logger(ctx).Error("Error marshalling party", "error", err)
		return err
	}

	// store set party id
	store.Set(GetMpcSetKey(party.ID), out)

	return nil
}

func (k *Keeper) GetAllMpcSets(ctx sdk.Context) (mpcSets []*hmTypes.PartyID) {
	// iterate through mpcs and create mpc update array
	k.IterateMpcSetsAndApplyFn(ctx, func(val hmTypes.PartyID) error {
		// append to list of validatorUpdates
		mpcSets = append(mpcSets, &val)
		return nil
	})

	return
}

// AddNewMpc adds new mpc to store
func (k *Keeper) AddNewMpc(ctx sdk.Context, mpc hmTypes.Mpc) error {
	store := ctx.KVStore(k.storeKey)

	out, err := k.cdc.MarshalBinaryBare(mpc)
	if err != nil {
		k.Logger(ctx).Error("Error marshalling mpc", "error", err)
		return err
	}

	// store set mpc id
	store.Set(GetMpcKey(mpc.ID), out)

	return nil
}

// AddNewMpcSign adds new mpc to store
func (k *Keeper) AddNewMpcSign(ctx sdk.Context, sign hmTypes.MpcSign) error {
	store := ctx.KVStore(k.storeKey)

	out, err := k.cdc.MarshalBinaryBare(sign)
	if err != nil {
		k.Logger(ctx).Error("Error marshalling mpc", "error", err)
		return err
	}

	// store set mpc id
	store.Set(GetMpcSignKey(sign.ID), out)

	return nil
}

// UpdateMpcSign adds new mpc to store
func (k *Keeper) UpdateMpcSign(ctx sdk.Context, sign hmTypes.MpcSign) error {
	store := ctx.KVStore(k.storeKey)

	if !k.HasMpcSign(ctx, sign.ID) {
		k.Logger(ctx).Error("Error UpdateMpcSign sign", "mpc sign not exist", sign.ID)
		return errors.New("mpc sign not exist:" + sign.ID)
	}

	out, err := k.cdc.MarshalBinaryBare(sign)
	if err != nil {
		k.Logger(ctx).Error("Error marshalling sign", "error", err)
		return err
	}

	// store set mpc id
	store.Set(GetMpcSignKey(sign.ID), out)

	return nil
}

// GetMpc fetches mpc indexed by id from store
func (k *Keeper) GetMpc(ctx sdk.Context, id string) (*hmTypes.Mpc, error) {
	store := ctx.KVStore(k.storeKey)
	mpcKey := GetMpcKey(id)

	// If we are starting from 0 there will be no mpcKey present
	if !store.Has(mpcKey) {
		return nil, errors.New("mpc not found for id")
	}

	var mpc hmTypes.Mpc
	if err := k.cdc.UnmarshalBinaryBare(store.Get(mpcKey), &mpc); err != nil {
		return nil, err
	}

	return &mpc, nil
}

// GetMpcSign fetches mpc indexed by id from store
func (k *Keeper) GetMpcSign(ctx sdk.Context, id string) (*hmTypes.MpcSign, error) {
	store := ctx.KVStore(k.storeKey)
	signKey := GetMpcSignKey(id)

	// If we are starting from 0 there will be no signKey present
	if !store.Has(signKey) {
		return nil, errors.New("mpc sign not found for id:" + id)
	}

	var sign hmTypes.MpcSign
	if err := k.cdc.UnmarshalBinaryBare(store.Get(signKey), &sign); err != nil {
		return nil, err
	}

	return &sign, nil
}

func (k *Keeper) HasMpc(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	mpcKey := GetMpcKey(id)
	return store.Has(mpcKey)
}

func (k *Keeper) HasMpcSign(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	signKey := GetMpcSignKey(id)
	return store.Has(signKey)
}

// GetAllMpcs fetches all indexed by id from store
func (k *Keeper) GetAllMpcs(ctx sdk.Context) (mpcs []*hmTypes.Mpc) {
	// iterate through mpcs and create mpc update array
	k.IterateMpcsAndApplyFn(ctx, func(mpc hmTypes.Mpc) error {
		// append to list of validatorUpdates
		mpcs = append(mpcs, &mpc)
		return nil
	})

	return
}

// GetMpcList returns all mpcs with params like page and limit
func (k *Keeper) GetMpcList(ctx sdk.Context, page uint64, limit uint64) ([]hmTypes.Mpc, error) {
	store := ctx.KVStore(k.storeKey)

	// have max limit
	if limit > 20 {
		limit = 20
	}

	// get paginated iterator
	iterator := hmTypes.KVStorePrefixIteratorPaginated(store, MpcPrefixKey, uint(page), uint(limit))

	// loop through validators to get valid validators
	var mpcs []hmTypes.Mpc

	for ; iterator.Valid(); iterator.Next() {
		var mpc hmTypes.Mpc
		if err := k.cdc.UnmarshalBinaryBare(iterator.Value(), &mpc); err == nil {
			mpcs = append(mpcs, mpc)
		}
	}

	return mpcs, nil
}

// GetLastMpc fetches last mpc using lastStartBlock
func (k *Keeper) GetLastMpc(ctx sdk.Context, mpcType hmTypes.MpcType) (*hmTypes.Mpc, error) {
	store := ctx.KVStore(k.storeKey)

	var lastMpcID string
	lastMpcIDKey := append(LastMpcIDKey, []byte(fmt.Sprintf("%v", mpcType))...)

	if store.Has(lastMpcIDKey) {
		lastMpcID = string(store.Get(lastMpcIDKey))
	}

	return k.GetMpc(ctx, lastMpcID)
}

// UpdateLastMpc updates the last mpc start block
func (k *Keeper) UpdateLastMpc(ctx sdk.Context, id string, mpcType hmTypes.MpcType) {
	store := ctx.KVStore(k.storeKey)
	store.Set(append(LastMpcIDKey, []byte(fmt.Sprintf("%v", mpcType))...), []byte(id))
}

//
// Utils
//

// IterateMpcsAndApplyFn iterates mpcs and apply the given function.
func (k *Keeper) IterateMpcsAndApplyFn(ctx sdk.Context, f func(mpc hmTypes.Mpc) error) {
	store := ctx.KVStore(k.storeKey)

	// get mpc iterator
	iterator := sdk.KVStorePrefixIterator(store, MpcPrefixKey)
	defer iterator.Close()

	// loop through mpcs to get valid mpcs
	for ; iterator.Valid(); iterator.Next() {
		// unmarshall mpc
		var result hmTypes.Mpc
		if err := k.cdc.UnmarshalBinaryBare(iterator.Value(), &result); err != nil {
			k.Logger(ctx).Error("Error UnmarshalBinaryBare", "error", err)
		}
		// call function and return if required
		if err := f(result); err != nil {
			return
		}
	}
}

func (k *Keeper) IterateMpcSetsAndApplyFn(ctx sdk.Context, f func(mpc hmTypes.PartyID) error) {
	store := ctx.KVStore(k.storeKey)

	// get mpc set iterator
	iterator := sdk.KVStorePrefixIterator(store, MpcSetPrefixKey)
	defer iterator.Close()

	// loop through mpcs to get valid mpcs
	for ; iterator.Valid(); iterator.Next() {
		// unmarshall mpc set
		var result hmTypes.PartyID
		if err := k.cdc.UnmarshalBinaryBare(iterator.Value(), &result); err != nil {
			k.Logger(ctx).Error("Error UnmarshalBinaryBare", "error", err)
		}
		// call function and return if required
		if err := f(result); err != nil {
			return
		}
	}
}

func (k *Keeper) IsValidator(ctx sdk.Context, proposer hmTypes.ThemisAddress) bool {
	vals := k.sk.GetAllValidators(ctx)
	for _, val := range vals {
		if val.Signer.EthAddress().Hex() == proposer.EthAddress().Hex() && val.IsCurrentValidator(k.sk.GetL1Batch(ctx)) {
			return true
		}
	}
	return false
}
