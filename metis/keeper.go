package metis

import (
	"errors"
	"math/big"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/ethereum/go-ethereum/common"

	"github.com/metis-seq/themis/chainmanager"
	"github.com/metis-seq/themis/helper"
	metisTypes "github.com/metis-seq/themis/metis/types"
	"github.com/metis-seq/themis/params/subspace"
	"github.com/metis-seq/themis/staking"
	hmTypes "github.com/metis-seq/themis/types"
)

var (
	LastSpanIDKey           = []byte{0x35} // Key to store last span start block
	SpanPrefixKey           = []byte{0x36} // prefix key to store span
	ReSpanPrefixKey         = []byte{0x37} // prefix key to store re-span
	LastProcessedEthBlock   = []byte{0x38} // key to store last processed eth block for seed
	ReSpanTimePrefixKey     = []byte{0x39} // prefix key to store re-span time
	ReSpanProposerPrefixKey = []byte{0x40} // prefix key to store re-span proposer
	MetisTxPrefixKey        = []byte{0x41} // prefix key to store metis-tx
	SpanProposerPrefixKey   = []byte{0x42} // prefix key to store span proposer
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
	// chain manager keeper
	chainKeeper chainmanager.Keeper
}

// NewKeeper is the constructor of Keeper
func NewKeeper(
	cdc *codec.Codec,
	storeKey sdk.StoreKey,
	paramSpace subspace.Subspace,
	codespace sdk.CodespaceType,
	chainKeeper chainmanager.Keeper,
	stakingKeeper staking.Keeper,
	caller helper.ContractCaller,
) Keeper {
	return Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		paramSpace:     paramSpace.WithKeyTable(metisTypes.ParamKeyTable()),
		codespace:      codespace,
		chainKeeper:    chainKeeper,
		sk:             stakingKeeper,
		contractCaller: caller,
	}
}

// Codespace returns the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", metisTypes.ModuleName)
}

// GetSpanKey appends prefix to start block
func GetSpanKey(id uint64) []byte {
	return append(SpanPrefixKey, []byte(strconv.FormatUint(id, 10))...)
}

// GetReSpanKey appends prefix to start block
func GetReSpanKey(msg metisTypes.MsgReProposeSpan) []byte {
	var respanID []byte
	respanID = append(ReSpanPrefixKey, []byte(strconv.FormatUint(msg.ID, 10))...)
	respanID = append(respanID, []byte(strconv.FormatUint(msg.StartBlock, 10))...)
	respanID = append(respanID, []byte(strconv.FormatUint(msg.EndBlock, 10))...)
	respanID = append(respanID, []byte(strconv.FormatUint(msg.CurrentL2Height, 10))...)
	respanID = append(respanID, []byte(strconv.FormatUint(msg.CurrentL2Epoch, 10))...)
	return respanID
}

func GetReSpanProposerKey(msg metisTypes.MsgReProposeSpan) []byte {
	var respanProposerID []byte
	respanProposerID = append(ReSpanProposerPrefixKey, []byte(strconv.FormatUint(msg.ID, 10))...)
	respanProposerID = append(respanProposerID, []byte(strconv.FormatUint(msg.StartBlock, 10))...)
	respanProposerID = append(respanProposerID, []byte(strconv.FormatUint(msg.EndBlock, 10))...)
	respanProposerID = append(respanProposerID, []byte(strconv.FormatUint(msg.CurrentL2Height, 10))...)
	respanProposerID = append(respanProposerID, []byte(strconv.FormatUint(msg.CurrentL2Epoch, 10))...)
	respanProposerID = append(respanProposerID, []byte(msg.ChainID)...)
	return respanProposerID
}

func GetSpanProposerKey(msg metisTypes.MsgProposeSpan) []byte {
	var spanProposerID []byte
	spanProposerID = append(SpanProposerPrefixKey, []byte(strconv.FormatUint(msg.ID, 10))...)
	spanProposerID = append(spanProposerID, []byte(strconv.FormatUint(msg.StartBlock, 10))...)
	spanProposerID = append(spanProposerID, []byte(strconv.FormatUint(msg.EndBlock, 10))...)
	spanProposerID = append(spanProposerID, []byte(msg.ChainID)...)
	return spanProposerID
}

// GetReSpanTimeKey appends prefix to start block
func GetReSpanTimeKey(id uint64) []byte {
	return append(ReSpanTimePrefixKey, []byte(strconv.FormatUint(id, 10))...)
}

// AddNewSpan adds new span for metis to store
func (k *Keeper) AddNewSpan(ctx sdk.Context, span hmTypes.Span) error {
	store := ctx.KVStore(k.storeKey)

	out, err := k.cdc.MarshalBinaryBare(span)
	if err != nil {
		k.Logger(ctx).Error("Error marshalling span", "error", err)
		return err
	}

	// store set span id
	store.Set(GetSpanKey(span.ID), out)

	// update last span
	k.UpdateLastSpan(ctx, span.ID)

	return nil
}

// AddNewRawSpan adds new span for metis to store
func (k *Keeper) AddNewRawSpan(ctx sdk.Context, span hmTypes.Span) error {
	store := ctx.KVStore(k.storeKey)

	out, err := k.cdc.MarshalBinaryBare(span)
	if err != nil {
		k.Logger(ctx).Error("Error marshalling span", "error", err)
		return err
	}

	// store set span id
	store.Set(GetSpanKey(span.ID), out)

	return nil
}

// GetSpan fetches span indexed by id from store
func (k *Keeper) GetSpan(ctx sdk.Context, id uint64) (*hmTypes.Span, error) {
	store := ctx.KVStore(k.storeKey)
	spanKey := GetSpanKey(id)

	// If we are starting from 0 there will be no spanKey present
	if !store.Has(spanKey) {
		return nil, errors.New("span not found for id")
	}

	var span hmTypes.Span
	if err := k.cdc.UnmarshalBinaryBare(store.Get(spanKey), &span); err != nil {
		return nil, err
	}

	return &span, nil
}

func (k *Keeper) HasSpan(ctx sdk.Context, id uint64) bool {
	store := ctx.KVStore(k.storeKey)
	spanKey := GetSpanKey(id)

	return store.Has(spanKey)
}

func (k *Keeper) DelSpan(ctx sdk.Context, id uint64) {
	lastSpan, _ := k.GetLastSpan(ctx)
	if lastSpan != nil && lastSpan.ID == id {
		k.UpdateLastSpan(ctx, lastSpan.ID-1)
	}

	store := ctx.KVStore(k.storeKey)
	spanKey := GetSpanKey(id)
	store.Delete(spanKey)
}

func (k *Keeper) HasReSpanFinish(ctx sdk.Context, msg metisTypes.MsgReProposeSpan) bool {
	store := ctx.KVStore(k.storeKey)
	respanKey := GetReSpanKey(msg)
	return store.Has(respanKey)
}

func (k *Keeper) SetReSpanFinish(ctx sdk.Context, msg metisTypes.MsgReProposeSpan) error {
	store := ctx.KVStore(k.storeKey)
	respanKey := GetReSpanKey(msg)

	out, err := k.cdc.MarshalBinaryBare(msg)
	if err != nil {
		k.Logger(ctx).Error("Error marshalling span", "error", err)
		return err
	}

	// store set re-span msg
	store.Set(respanKey, out)
	return nil
}

func (k *Keeper) SetReSpanTime(ctx sdk.Context, id, time uint64) {
	store := ctx.KVStore(k.storeKey)
	respanTimeKey := GetReSpanTimeKey(id)
	store.Set(respanTimeKey, []byte(strconv.FormatUint(time, 10)))
}

func (k *Keeper) GetReSpanTime(ctx sdk.Context, id uint64) (uint64, error) {
	store := ctx.KVStore(k.storeKey)
	respanTimeKey := GetReSpanTimeKey(id)
	return strconv.ParseUint(string(store.Get(respanTimeKey)), 10, 64)
}

func (k *Keeper) SaveSpanProposer(ctx sdk.Context, msg metisTypes.MsgProposeSpan) error {
	store := ctx.KVStore(k.storeKey)
	spanProposerKey := GetSpanProposerKey(msg)

	// get exist proposers
	var allProposers []string
	proposerBytes := store.Get(spanProposerKey)
	if len(proposerBytes) > 0 {
		err := k.cdc.UnmarshalBinaryBare(proposerBytes, &allProposers)
		if err != nil {
			return err
		}

		// check repeat proposer
		var newProposerExist bool
		for _, proposer := range allProposers {
			if proposer == msg.Proposer.EthAddress().Hex() {
				newProposerExist = true
			}
		}

		if !newProposerExist {
			allProposers = append(allProposers, msg.Proposer.EthAddress().Hex())
		}
	} else {
		allProposers = append(allProposers, msg.Proposer.EthAddress().Hex())
	}

	allProposersBytes, err := k.cdc.MarshalBinaryBare(allProposers)
	if err != nil {
		return err
	}
	store.Set(spanProposerKey, allProposersBytes)
	return nil
}

func (k *Keeper) GetSpanAllProposers(ctx sdk.Context, msg metisTypes.MsgProposeSpan) ([]string, error) {
	store := ctx.KVStore(k.storeKey)
	spanProposerKey := GetSpanProposerKey(msg)

	// get exist proposers
	var allProposers []string
	proposerBytes := store.Get(spanProposerKey)
	if len(proposerBytes) > 0 {
		err := k.cdc.UnmarshalBinaryBare(proposerBytes, &allProposers)
		if err != nil {
			return nil, err
		}
	}

	return allProposers, nil
}

func (k *Keeper) SaveReSpanProposer(ctx sdk.Context, msg metisTypes.MsgReProposeSpan) error {
	store := ctx.KVStore(k.storeKey)
	respanProposerKey := GetReSpanProposerKey(msg)

	// get exist proposers
	var allProposers []string
	proposerBytes := store.Get(respanProposerKey)
	if len(proposerBytes) > 0 {
		err := k.cdc.UnmarshalBinaryBare(proposerBytes, &allProposers)
		if err != nil {
			return err
		}

		// check repeat proposer
		var newProposerExist bool
		for _, proposer := range allProposers {
			if proposer == msg.Proposer.EthAddress().Hex() {
				newProposerExist = true
			}
		}

		if !newProposerExist {
			allProposers = append(allProposers, msg.Proposer.EthAddress().Hex())
		}
	} else {
		allProposers = append(allProposers, msg.Proposer.EthAddress().Hex())
	}

	allProposersBytes, err := k.cdc.MarshalBinaryBare(allProposers)
	if err != nil {
		return err
	}
	store.Set(respanProposerKey, allProposersBytes)
	return nil
}

func (k *Keeper) GetReSpanAllProposers(ctx sdk.Context, msg metisTypes.MsgReProposeSpan) ([]string, error) {
	store := ctx.KVStore(k.storeKey)
	respanProposerKey := GetReSpanProposerKey(msg)

	// get exist proposers
	var allProposers []string
	proposerBytes := store.Get(respanProposerKey)
	if len(proposerBytes) > 0 {
		err := k.cdc.UnmarshalBinaryBare(proposerBytes, &allProposers)
		if err != nil {
			return nil, err
		}
	}

	return allProposers, nil
}

func (k *Keeper) DelReSpan(ctx sdk.Context, msg metisTypes.MsgReProposeSpan) {
	store := ctx.KVStore(k.storeKey)
	respanKey := GetReSpanKey(msg)
	store.Delete(respanKey)
}

// GetAllSpans fetches all indexed by id from store
func (k *Keeper) GetAllSpans(ctx sdk.Context) (spans []*hmTypes.Span) {
	// iterate through spans and create span update array
	k.IterateSpansAndApplyFn(ctx, func(span hmTypes.Span) error {
		// append to list of validatorUpdates
		spans = append(spans, &span)
		return nil
	})

	return
}

// GetSpanList returns all spans with params like page and limit
func (k *Keeper) GetSpanList(ctx sdk.Context, page uint64, limit uint64) ([]*hmTypes.Span, error) {
	store := ctx.KVStore(k.storeKey)

	// have max limit
	if limit > 20 {
		limit = 20
	}

	// get paginated iterator
	iterator := hmTypes.KVStorePrefixIteratorPaginated(store, SpanPrefixKey, uint(page), uint(limit))
	defer iterator.Close()

	// loop through validators to get valid validators
	var spans []*hmTypes.Span

	for ; iterator.Valid(); iterator.Next() {
		var span hmTypes.Span
		if err := k.cdc.UnmarshalBinaryBare(iterator.Value(), &span); err == nil {
			spans = append(spans, &span)
		}
	}

	return spans, nil
}

// GetReSpanList returns all re-spans with params like page and limit
func (k *Keeper) GetReSpanList(ctx sdk.Context, page uint64, limit uint64) ([]*metisTypes.MsgReProposeSpan, error) {
	store := ctx.KVStore(k.storeKey)

	// have max limit
	if limit > 20 {
		limit = 20
	}

	// get paginated iterator
	iterator := hmTypes.KVStorePrefixIteratorPaginated(store, ReSpanPrefixKey, uint(page), uint(limit))
	defer iterator.Close()

	// loop through validators to get valid validators
	var respans []*metisTypes.MsgReProposeSpan

	for ; iterator.Valid(); iterator.Next() {
		var respan metisTypes.MsgReProposeSpan
		if err := k.cdc.UnmarshalBinaryBare(iterator.Value(), &respan); err == nil {
			respans = append(respans, &respan)
		}
	}

	return respans, nil
}

// GetLastSpan fetches last span using lastStartBlock
func (k *Keeper) GetLastSpan(ctx sdk.Context) (*hmTypes.Span, error) {
	store := ctx.KVStore(k.storeKey)

	var lastSpanID uint64

	if store.Has(LastSpanIDKey) {
		// get last span id
		var err error
		if lastSpanID, err = strconv.ParseUint(string(store.Get(LastSpanIDKey)), 10, 64); err != nil {
			return nil, err
		}
	}

	return k.GetSpan(ctx, lastSpanID)
}

// FreezeSet freezes validator set for next span
func (k *Keeper) FreezeSet(ctx sdk.Context, id uint64, startBlock uint64, endBlock uint64, metisChainID string, seed common.Hash, unEligibleVals []hmTypes.ThemisAddress) error {
	// select next producers
	newProducers, err := k.SelectNextProducers(ctx, seed, unEligibleVals)
	if err != nil {
		return err
	}

	// increment last eth block
	k.IncrementLastEthBlock(ctx)

	// generate new span
	newSpan := hmTypes.NewSpan(
		id,
		startBlock,
		endBlock,
		k.sk.GetValidatorSet(ctx),
		newProducers,
		metisChainID,
	)

	return k.AddNewSpan(ctx, newSpan)
}

// FreezeFixedSet freezes validator set for next span
func (k *Keeper) FreezeFixedSet(ctx sdk.Context, id uint64, startBlock uint64, endBlock uint64, metisChainID string, fixedSequencer hmTypes.ThemisAddress) error {
	// increment last eth block
	k.IncrementLastEthBlock(ctx)

	var fixedProducer hmTypes.Validator
	valSet := k.sk.GetValidatorSet(ctx)
	for _, val := range valSet.Validators {
		if val.Signer == fixedSequencer {
			fixedProducer = *hmTypes.NewValidator(
				val.ID,
				val.StartBatch,
				val.EndBatch,
				val.Nonce,
				val.VotingPower,
				val.PubKey,
				val.Signer,
			)
			break
		}
	}

	if fixedProducer.Signer.EthAddress() == common.BigToAddress(big.NewInt(0)) {
		val, _ := k.sk.GetValidatorInfo(ctx, fixedSequencer[:])
		fixedProducer = *hmTypes.NewValidator(
			val.ID,
			val.StartBatch,
			val.EndBatch,
			val.Nonce,
			val.VotingPower,
			val.PubKey,
			val.Signer,
		)
	}

	span := hmTypes.NewSpan(
		id,
		startBlock,
		endBlock,
		k.sk.GetValidatorSet(ctx),
		[]hmTypes.Validator{fixedProducer},
		metisChainID,
	)

	var err error
	// check if span id exist
	if !k.HasSpan(ctx, id) {
		// generate new span
		err = k.AddNewSpan(ctx, span)
	} else {
		// update exist span
		err = k.AddNewRawSpan(ctx, span)
	}

	return err
}

// UpdateFixedSet update validator set for next span
func (k *Keeper) UpdateFixedSet(ctx sdk.Context, id uint64, startBlock uint64, endBlock uint64, metisChainID string, fixedSequencer hmTypes.ThemisAddress) error {
	// increment last eth block
	k.IncrementLastEthBlock(ctx)

	var fixedProducer hmTypes.Validator
	valSet := k.sk.GetValidatorSet(ctx)
	for _, val := range valSet.Validators {
		if val.Signer == fixedSequencer {
			fixedProducer = *hmTypes.NewValidator(
				val.ID,
				val.StartBatch,
				val.EndBatch,
				val.Nonce,
				val.VotingPower,
				val.PubKey,
				val.Signer,
			)
			break
		}
	}

	// generate new span
	newSpan := hmTypes.NewSpan(
		id,
		startBlock,
		endBlock,
		k.sk.GetValidatorSet(ctx),
		[]hmTypes.Validator{fixedProducer},
		metisChainID,
	)

	return k.AddNewSpan(ctx, newSpan)
}

// SelectNextProducers selects producers for next span
func (k *Keeper) SelectNextProducers(ctx sdk.Context, seed common.Hash, unEligibleVals []hmTypes.ThemisAddress) (vals []hmTypes.Validator, err error) {
	// spanEligibleVals are current validators who are not getting deactivated in between next span
	spanEligibleVals := k.sk.GetSpanEligibleValidators(ctx)
	producerCount := k.GetParams(ctx).ProducerCount

	// if producers to be selected is more than current validators no need to select/shuffle
	if len(spanEligibleVals) <= int(producerCount) {
		return spanEligibleVals, nil
	}

	// exclude unqualified validator
	var finalSpanEligibleVals []hmTypes.Validator
	for _, spanEligibleVal := range spanEligibleVals {
		if len(unEligibleVals) == 0 {
			finalSpanEligibleVals = append(finalSpanEligibleVals, spanEligibleVal)
		} else {
			for _, uv := range unEligibleVals {
				if spanEligibleVal.Signer != uv {
					finalSpanEligibleVals = append(finalSpanEligibleVals, spanEligibleVal)
				}
			}
		}
	}
	// select next producers using seed as blockheader hash
	fn := SelectNextProducers

	newProducersIds, err := fn(seed, finalSpanEligibleVals, producerCount)
	if err != nil {
		return vals, err
	}

	IDToPower := make(map[uint64]uint64)
	for _, ID := range newProducersIds {
		IDToPower[ID] = IDToPower[ID] + 1
	}

	for key, value := range IDToPower {
		if val, ok := k.sk.GetValidatorFromValID(ctx, hmTypes.NewValidatorID(key)); ok {
			val.VotingPower = int64(value)
			vals = append(vals, val)
		}
	}

	// sort by address
	vals = hmTypes.SortValidatorByAddress(vals)

	return vals, nil
}

// UpdateLastSpan updates the last span start block
func (k *Keeper) UpdateLastSpan(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(LastSpanIDKey, []byte(strconv.FormatUint(id, 10)))
}

// UpdateLastSpanEndBlock updates the last span end block
func (k *Keeper) UpdateLastSpanEndBlock(ctx sdk.Context, endBlock uint64) error {
	lastSpan, err := k.GetLastSpan(ctx)
	if err != nil {
		return err
	}
	lastSpan.EndBlock = endBlock

	return k.AddNewSpan(ctx, *lastSpan)
}

// UpdateFixedSpanEndBlock updates the last span end block
func (k *Keeper) UpdateFixedSpanEndBlock(ctx sdk.Context, spanId, endBlock uint64) error {
	lastSpan, err := k.GetLastSpan(ctx)
	if err != nil {
		return err
	}
	if lastSpan.ID-1 != spanId && lastSpan.ID != spanId {
		return errors.New("spanId not allowed")
	}

	span, err := k.GetSpan(ctx, spanId)
	if err != nil {
		return err
	}

	span.EndBlock = endBlock
	return k.AddNewRawSpan(ctx, *span)
}

func (k *Keeper) IsSpanValidator(ctx sdk.Context, proposer hmTypes.ThemisAddress) bool {
	currentL1Batch := k.sk.GetL1Batch(ctx)
	vals := k.sk.GetAllValidators(ctx)
	for _, val := range vals {
		if val.Signer.EthAddress().Hex() == proposer.EthAddress().Hex() && val.IsCurrentValidator(currentL1Batch) {
			return true
		}
	}
	return false
}

// IncrementLastEthBlock increment last eth block
func (k *Keeper) IncrementLastEthBlock(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)

	lastEthBlock := big.NewInt(0)
	if store.Has(LastProcessedEthBlock) {
		lastEthBlock = lastEthBlock.SetBytes(store.Get(LastProcessedEthBlock))
	}

	store.Set(LastProcessedEthBlock, lastEthBlock.Add(lastEthBlock, big.NewInt(1)).Bytes())
}

// SetLastEthBlock sets last eth block number
func (k *Keeper) SetLastEthBlock(ctx sdk.Context, blockNumber *big.Int) {
	store := ctx.KVStore(k.storeKey)
	store.Set(LastProcessedEthBlock, blockNumber.Bytes())
}

// GetLastEthBlock get last processed Eth block for seed
func (k *Keeper) GetLastEthBlock(ctx sdk.Context) *big.Int {
	store := ctx.KVStore(k.storeKey)

	lastEthBlock := big.NewInt(0)
	if store.Has(LastProcessedEthBlock) {
		lastEthBlock = lastEthBlock.SetBytes(store.Get(LastProcessedEthBlock))
	}

	return lastEthBlock
}

func (k Keeper) GetNextSpanSeed(ctx sdk.Context) (common.Hash, error) {
	// lastEthBlock := k.GetLastEthBlock(ctx)

	// increment last processed header block number
	// newEthBlock := lastEthBlock.Add(lastEthBlock, big.NewInt(1))
	// k.Logger(ctx).Debug("newEthBlock to generate seed", "newEthBlock", newEthBlock)

	// fetch block header from mainchain
	// blockHeader, err := k.contractCaller.GetMainChainBlock(newEthBlock)
	// if err != nil {
	// 	k.Logger(ctx).Error("Error fetching block header from mainchain while calculating next span seed", "error", err)
	// 	return common.Hash{}, err
	// }
	// return blockHeader.Hash(), nil

	lastSpan, err := k.GetLastSpan(ctx)
	if err != nil {
		k.Logger(ctx).Error("Error GetLastSpan while calculating next span seed", "error", err)
		return common.Hash{}, err
	}
	seedHash := sha256Hash(big.NewInt(int64(lastSpan.ID + 1)).Bytes())
	return common.BytesToHash(seedHash[:]), nil
}

// -----------------------------------------------------------------------------
// Params

// SetParams sets the metis module's parameters.
func (k *Keeper) SetParams(ctx sdk.Context, params metisTypes.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// GetParams gets the metis module's parameters.
func (k *Keeper) GetParams(ctx sdk.Context) (params metisTypes.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return
}

//
// Utils
//

// IterateSpansAndApplyFn iterates spans and apply the given function.
func (k *Keeper) IterateSpansAndApplyFn(ctx sdk.Context, f func(span hmTypes.Span) error) {
	store := ctx.KVStore(k.storeKey)

	// get span iterator
	iterator := sdk.KVStorePrefixIterator(store, SpanPrefixKey)
	defer iterator.Close()

	// loop through spans to get valid spans
	for ; iterator.Valid(); iterator.Next() {
		// unmarshall span
		var result hmTypes.Span
		if err := k.cdc.UnmarshalBinaryBare(iterator.Value(), &result); err != nil {
			k.Logger(ctx).Error("Error UnmarshalBinaryBare", "error", err)
		}
		// call function and return if required
		if err := f(result); err != nil {
			return
		}
	}
}

// GetMetisTxKey appends prefix to start block
func GetMetisTxKey(txHash string) []byte {
	return append(MetisTxPrefixKey, []byte(txHash)...)
}

func (k *Keeper) HasMetisTx(ctx sdk.Context, txHash string) bool {
	store := ctx.KVStore(k.storeKey)
	metisTxKey := GetMetisTxKey(txHash)
	return store.Has(metisTxKey)
}

func (k *Keeper) GetMetisTx(ctx sdk.Context, txHash string) (*metisTypes.MsgMetisTx, error) {
	store := ctx.KVStore(k.storeKey)
	metisTxKey := GetMetisTxKey(txHash)

	// If we are starting from 0 there will be no spanKey present
	if !store.Has(metisTxKey) {
		return nil, errors.New("span not found for id")
	}

	var metisTx metisTypes.MsgMetisTx
	if err := k.cdc.UnmarshalBinaryBare(store.Get(metisTxKey), &metisTx); err != nil {
		return nil, err
	}

	return &metisTx, nil
}

func (k *Keeper) AddMetisTx(ctx sdk.Context, metisTx metisTypes.MsgMetisTx) error {
	store := ctx.KVStore(k.storeKey)
	metisTxKey := GetMetisTxKey(metisTx.TxHash.Hex())

	out, err := k.cdc.MarshalBinaryBare(metisTx)
	if err != nil {
		k.Logger(ctx).Error("Error marshalling metisTx", "error", err)
		return err
	}

	// store set span id
	store.Set(metisTxKey, out)
	return nil
}
