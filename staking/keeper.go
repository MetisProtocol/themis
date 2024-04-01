package staking

import (
	"encoding/hex"
	"errors"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/metis-seq/themis/chainmanager"
	"github.com/metis-seq/themis/helper"
	"github.com/metis-seq/themis/params/subspace"
	"github.com/metis-seq/themis/staking/types"
	hmTypes "github.com/metis-seq/themis/types"
)

var (
	DefaultValue = []byte{0x01} // Value to store in Cache and CacheACK & ValidatorSetChange Flag

	L1BatchKey             = []byte{0x11} // key to store L1 staking batch
	ValidatorsKey          = []byte{0x21} // prefix for each key to a validator
	ValidatorMapKey        = []byte{0x22} // prefix for each key for validator map
	CurrentValidatorSetKey = []byte{0x23} // Key to store current validator set
	StakingSequenceKey     = []byte{0x24} // prefix for each key for staking sequence map
)

// ModuleCommunicator manages different module interaction
type ModuleCommunicator interface {
	GetL1Batch(ctx sdk.Context) uint64
	SetCoins(ctx sdk.Context, addr hmTypes.ThemisAddress, amt sdk.Coins) sdk.Error
	GetCoins(ctx sdk.Context, addr hmTypes.ThemisAddress) sdk.Coins
	SendCoins(ctx sdk.Context, from hmTypes.ThemisAddress, to hmTypes.ThemisAddress, amt sdk.Coins) sdk.Error
}

// Keeper stores all related data
type Keeper struct {
	cdc *codec.Codec
	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey
	// codespacecodespace
	codespace sdk.CodespaceType
	// param space
	paramSpace subspace.Subspace
	// chain manager keeper
	chainKeeper chainmanager.Keeper
	// module communicator
	moduleCommunicator ModuleCommunicator
}

// NewKeeper create new keeper
func NewKeeper(
	cdc *codec.Codec,
	storeKey sdk.StoreKey,
	paramSpace subspace.Subspace,
	codespace sdk.CodespaceType,
	chainKeeper chainmanager.Keeper,
	moduleCommunicator ModuleCommunicator,
) Keeper {
	keeper := Keeper{
		cdc:                cdc,
		storeKey:           storeKey,
		paramSpace:         paramSpace.WithKeyTable(types.ParamKeyTable()),
		codespace:          codespace,
		chainKeeper:        chainKeeper,
		moduleCommunicator: moduleCommunicator,
	}

	return keeper
}

// Codespace returns the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ModuleName)
}

// GetValidatorKey drafts the validator key for addresses
func GetValidatorKey(address []byte) []byte {
	return append(ValidatorsKey, address...)
}

// GetValidatorMapKey returns validator map
func GetValidatorMapKey(address []byte) []byte {
	return append(ValidatorMapKey, address...)
}

// GetStakingSequenceKey returns staking sequence key
func GetStakingSequenceKey(sequence string) []byte {
	return append(StakingSequenceKey, []byte(sequence)...)
}

// AddValidator adds validator indexed with address
func (k *Keeper) AddValidator(ctx sdk.Context, validator hmTypes.Validator) error {
	store := ctx.KVStore(k.storeKey)

	bz, err := hmTypes.MarshallValidator(k.cdc, validator)
	if err != nil {
		return err
	}

	// store validator with address prefixed with validator key as index
	store.Set(GetValidatorKey(validator.Signer.Bytes()), bz)

	k.Logger(ctx).Debug("Validator stored", "key", hex.EncodeToString(GetValidatorKey(validator.Signer.Bytes())), "validator", validator.String())

	// add validator to validator ID => SignerAddress map
	k.SetValidatorIDToSignerAddr(ctx, validator.ID, validator.Signer)

	return nil
}

// IsCurrentValidatorByAddress check if validator is in current validator set by signer address
func (k *Keeper) IsCurrentValidatorByAddress(ctx sdk.Context, address []byte) bool {
	// get l1 batch
	l1Batch := k.moduleCommunicator.GetL1Batch(ctx)

	// get validator info
	validator, err := k.GetValidatorInfo(ctx, address)
	if err != nil {
		return false
	}

	// check if validator is current validator
	return validator.IsCurrentValidator(l1Batch)
}

// GetValidatorInfo returns validator
func (k *Keeper) GetValidatorInfo(ctx sdk.Context, address []byte) (validator hmTypes.Validator, err error) {
	store := ctx.KVStore(k.storeKey)

	// check if validator exists
	key := GetValidatorKey(address)
	if !store.Has(key) {
		return validator, errors.New("Validator not found")
	}

	// unmarshall validator and return
	validator, err = hmTypes.UnmarshallValidator(k.cdc, store.Get(key))
	if err != nil {
		return validator, err
	}

	// return true if validator
	return validator, nil
}

// GetActiveValidatorInfo returns active validator
func (k *Keeper) GetActiveValidatorInfo(ctx sdk.Context, address []byte) (validator hmTypes.Validator, err error) {
	validator, err = k.GetValidatorInfo(ctx, address)
	if err != nil {
		return validator, err
	}

	// get l1 batch
	l1Batch := k.moduleCommunicator.GetL1Batch(ctx)
	if !validator.IsCurrentValidator(l1Batch) {
		return validator, errors.New("Validator is not active")
	}

	// return true if validator
	return validator, nil
}

// GetCurrentValidators returns all validators who are in validator set
func (k *Keeper) GetCurrentValidators(ctx sdk.Context) (validators []hmTypes.Validator) {
	// get l1 batch
	l1Batch := k.moduleCommunicator.GetL1Batch(ctx)

	// Get validators
	// iterate through validator list
	k.IterateValidatorsAndApplyFn(ctx, func(validator hmTypes.Validator) error {
		// check if validator is valid for current batch
		if validator.IsCurrentValidator(l1Batch) {
			// append if validator is current valdiator
			validators = append(validators, validator)
		}
		return nil
	})

	return
}

func (k *Keeper) GetTotalPower(ctx sdk.Context) (totalPower int64) {
	k.IterateCurrentValidatorsAndApplyFn(ctx, func(validator *hmTypes.Validator) bool {
		totalPower += validator.VotingPower
		return true
	})

	return
}

// GetSpanEligibleValidators returns current validators who are not getting deactivated in between next span
func (k *Keeper) GetSpanEligibleValidators(ctx sdk.Context) (validators []hmTypes.Validator) {
	// get l1 batch
	l1Batch := k.moduleCommunicator.GetL1Batch(ctx)

	// Get validators and iterate through validator list
	k.IterateValidatorsAndApplyFn(ctx, func(validator hmTypes.Validator) error {
		// check if validator is valid for current batch and endBatch is not set.
		if validator.EndBatch == 0 && validator.IsCurrentValidator(l1Batch) {
			// append if validator is current valdiator
			validators = append(validators, validator)
		}
		return nil
	})

	return
}

// GetAllValidators returns all validators
func (k *Keeper) GetAllValidators(ctx sdk.Context) (validators []*hmTypes.Validator) {
	// iterate through validators and create validator update array
	k.IterateValidatorsAndApplyFn(ctx, func(validator hmTypes.Validator) error {
		// append to list of validatorUpdates
		validators = append(validators, &validator)
		return nil
	})

	return
}

// IterateValidatorsAndApplyFn interate validators and apply the given function.
func (k *Keeper) IterateValidatorsAndApplyFn(ctx sdk.Context, f func(validator hmTypes.Validator) error) {
	store := ctx.KVStore(k.storeKey)

	// get validator iterator
	iterator := sdk.KVStorePrefixIterator(store, ValidatorsKey)
	defer iterator.Close()

	// loop through validators to get valid validators
	for ; iterator.Valid(); iterator.Next() {
		// unmarshall validator
		validator, _ := hmTypes.UnmarshallValidator(k.cdc, iterator.Value())
		// call function and return if required
		if err := f(validator); err != nil {
			return
		}
	}
}

// UpdateSigner updates validator with signer and pubkey + validator => signer map
func (k *Keeper) UpdateSigner(ctx sdk.Context, newSigner hmTypes.ThemisAddress, newPubkey hmTypes.PubKey, prevSigner hmTypes.ThemisAddress) error {
	// get old validator from state and make power 0
	validator, err := k.GetValidatorInfo(ctx, prevSigner.Bytes())
	if err != nil {
		k.Logger(ctx).Error("Unable to fetch valiator from store")
		return err
	}

	// copy power to reassign below
	validatorPower := validator.VotingPower
	validator.VotingPower = 0

	// update validator
	if err := k.AddValidator(ctx, validator); err != nil {
		k.Logger(ctx).Error("UpdateSigner | AddValidator", "error", err)
	}

	//update signer in prev Validator
	validator.Signer = newSigner
	validator.PubKey = newPubkey
	validator.VotingPower = validatorPower

	// add updated validator to store with new key
	if err = k.AddValidator(ctx, validator); err != nil {
		k.Logger(ctx).Error("UpdateSigner | AddValidator", "error", err)
	}

	return nil
}

// UpdateValidatorSetInStore adds validator set to store
func (k *Keeper) UpdateValidatorSetInStore(ctx sdk.Context, newValidatorSet hmTypes.ValidatorSet) error {
	// TODO check if we may have to delay this by 1 height to sync with tendermint validator updates
	store := ctx.KVStore(k.storeKey)

	// marshall validator set
	bz, err := k.cdc.MarshalBinaryBare(newValidatorSet)
	if err != nil {
		return err
	}

	// set validator set with CurrentValidatorSetKey as key in store
	store.Set(CurrentValidatorSetKey, bz)

	return nil
}

// GetValidatorSet returns current Validator Set from store
func (k *Keeper) GetValidatorSet(ctx sdk.Context) (validatorSet hmTypes.ValidatorSet) {
	store := ctx.KVStore(k.storeKey)
	// get current validator set from store
	bz := store.Get(CurrentValidatorSetKey)
	// unmarhsall

	if err := k.cdc.UnmarshalBinaryBare(bz, &validatorSet); err != nil {
		k.Logger(ctx).Error("GetValidatorSet | UnmarshalBinaryBare", "error", err)
	}

	// return validator set
	return validatorSet
}

// IncrementAccum increments accum for validator set by n times and replace validator set in store
func (k *Keeper) IncrementAccum(ctx sdk.Context, times int) {
	// get validator set
	validatorSet := k.GetValidatorSet(ctx)

	// increment accum
	validatorSet.IncrementProposerPriority(times)

	// replace
	if err := k.UpdateValidatorSetInStore(ctx, validatorSet); err != nil {
		k.Logger(ctx).Error("IncrementAccum | UpdateValidatorSetInStore", "error", err)
	}
}

// GetNextProposer returns next proposer
func (k *Keeper) GetNextProposer(ctx sdk.Context) *hmTypes.Validator {
	// get validator set
	validatorSet := k.GetValidatorSet(ctx)

	// Increment accum in copy
	copiedValidatorSet := validatorSet.CopyIncrementProposerPriority(1)

	// get signer address for next signer
	return copiedValidatorSet.GetProposer()
}

// GetCurrentProposer returns current proposer
func (k *Keeper) GetCurrentProposer(ctx sdk.Context) *hmTypes.Validator {
	// get validator set
	validatorSet := k.GetValidatorSet(ctx)

	// return get proposer
	return validatorSet.GetProposer()
}

// SetValidatorIDToSignerAddr sets mapping for validator ID to signer address
func (k *Keeper) SetValidatorIDToSignerAddr(ctx sdk.Context, valID hmTypes.ValidatorID, signerAddr hmTypes.ThemisAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetValidatorMapKey(valID.Bytes()), signerAddr.Bytes())
}

// GetSignerFromValidatorID get signer address from validator ID
func (k *Keeper) GetSignerFromValidatorID(ctx sdk.Context, valID hmTypes.ValidatorID) (common.Address, bool) {
	store := ctx.KVStore(k.storeKey)
	key := GetValidatorMapKey(valID.Bytes())
	// check if validator address has been mapped
	if !store.Has(key) {
		return helper.ZeroAddress, false
	}
	// return address from bytes
	return common.BytesToAddress(store.Get(key)), true
}

// GetValidatorFromValID returns signer from validator ID
func (k *Keeper) GetValidatorFromValID(ctx sdk.Context, valID hmTypes.ValidatorID) (validator hmTypes.Validator, ok bool) {
	signerAddr, ok := k.GetSignerFromValidatorID(ctx, valID)
	if !ok {
		return validator, ok
	}

	// query for validator signer address
	validator, err := k.GetValidatorInfo(ctx, signerAddr.Bytes())
	if err != nil {
		return validator, false
	}

	return validator, true
}

// GetLastUpdated get last updated at for validator
func (k *Keeper) GetLastUpdated(ctx sdk.Context, valID hmTypes.ValidatorID) (updatedAt string, found bool) {
	// get validator
	validator, ok := k.GetValidatorFromValID(ctx, valID)
	if !ok {
		return "", false
	}

	return validator.LastUpdated, true
}

// IterateCurrentValidatorsAndApplyFn iterate through current validators
func (k *Keeper) IterateCurrentValidatorsAndApplyFn(ctx sdk.Context, f func(validator *hmTypes.Validator) bool) {
	currentValidatorSet := k.GetValidatorSet(ctx)
	for _, v := range currentValidatorSet.Validators {
		if stop := f(v); stop {
			return
		}
	}
}

//
// Staking sequence
//

// SetStakingSequence sets staking sequence
func (k *Keeper) SetStakingSequence(ctx sdk.Context, sequence string) {
	store := ctx.KVStore(k.storeKey)

	store.Set(GetStakingSequenceKey(sequence), DefaultValue)
}

// HasStakingSequence checks if staking sequence already exists
func (k *Keeper) HasStakingSequence(ctx sdk.Context, sequence string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(GetStakingSequenceKey(sequence))
}

// GetStakingSequences checks if Staking already exists
func (k *Keeper) GetStakingSequences(ctx sdk.Context) (sequences []string) {
	k.IterateStakingSequencesAndApplyFn(ctx, func(sequence string) error {
		sequences = append(sequences, sequence)
		return nil
	})

	return
}

// IterateStakingSequencesAndApplyFn interate validators and apply the given function.
func (k *Keeper) IterateStakingSequencesAndApplyFn(ctx sdk.Context, f func(sequence string) error) {
	store := ctx.KVStore(k.storeKey)

	// get sequence iterator
	iterator := sdk.KVStorePrefixIterator(store, StakingSequenceKey)
	defer iterator.Close()

	// loop through validators to get valid validators
	for ; iterator.Valid(); iterator.Next() {
		sequence := string(iterator.Key()[len(StakingSequenceKey):])

		// call function and return if required
		if err := f(sequence); err != nil {
			return
		}
	}
}

// Unjail a validator
func (k *Keeper) Unjail(ctx sdk.Context, valID hmTypes.ValidatorID) {
	// get validator from state and make jailed = false
	validator, found := k.GetValidatorFromValID(ctx, valID)
	if !found {
		k.Logger(ctx).Error("Unable to fetch valiator from store")
		return
	}

	if !validator.Jailed {
		k.Logger(ctx).Info("Already unjailed.")
		return
	}
	// unjail validator
	validator.Jailed = false

	// add updated validator to store with new key
	if err := k.AddValidator(ctx, validator); err != nil {
		k.Logger(ctx).Error("Failed to add validator", "Error", err)
	}
}

//
// L1 batch
//

// GetL1Batch returns current L1Batch count
func (k Keeper) GetL1Batch(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	// check if ack count is there
	if store.Has(L1BatchKey) {
		// get current L1Batch count
		l1BatchId, err := strconv.ParseUint(string(store.Get(L1BatchKey)), 10, 64)
		if err != nil {
			k.Logger(ctx).Error("Unable to convert key to int")
		} else {
			return l1BatchId
		}
	}

	return 0
}

// UpdateL1BatchWithValue updates L1Batch with value
func (k Keeper) UpdateL1BatchWithValue(ctx sdk.Context, value uint64) {
	store := ctx.KVStore(k.storeKey)

	// convert
	l1BatchId := []byte(strconv.FormatUint(value, 10))

	// update
	store.Set(L1BatchKey, l1BatchId)
}

// UpdateL1Batch updates L1Batch count by 1
func (k Keeper) UpdateL1Batch(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)

	// get current L1Batch Count
	L1Batch := k.GetL1Batch(ctx)

	// increment by 1
	L1Batchs := []byte(strconv.FormatUint(L1Batch+1, 10))

	// update
	store.Set(L1BatchKey, L1Batchs)
}
