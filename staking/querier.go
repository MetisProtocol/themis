package staking

import (
	"fmt"
	"math/big"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	jsoniter "github.com/json-iterator/go"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/metis-seq/themis/helper"
	"github.com/metis-seq/themis/staking/types"
	hmTypes "github.com/metis-seq/themis/types"
)

// NewQuerier returns querier for staking Rest endpoints
func NewQuerier(keeper Keeper, contractCaller helper.IContractCaller) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryCurrentValidatorSet:
			return handleQueryCurrentValidatorSet(ctx, req, keeper)
		case types.QueryAllValidatorSet:
			return handleQueryAllValidatorSet(ctx, req, keeper)
		case types.QuerySigner:
			return handleQuerySigner(ctx, req, keeper)
		case types.QueryValidator:
			return handleQueryValidator(ctx, req, keeper)
		case types.QueryValidatorStatus:
			return handleQueryValidatorStatus(ctx, req, keeper)
		case types.QueryProposer:
			return handleQueryProposer(ctx, req, keeper)
		case types.QueryCurrentProposer:
			return handleQueryCurrentProposer(ctx, req, keeper)
		case types.QueryStakingSequence:
			return handleQueryStakingSequence(ctx, req, keeper, contractCaller)
		case types.QueryTotalValidatorPower:
			return handleQueryTotalValidatorPower(ctx, req, keeper)
		case types.QueryCurrentL1Batch:
			return handleQueryCurrentBatch(ctx, req, keeper)

		default:
			return nil, sdk.ErrUnknownRequest("unknown staking query endpoint")
		}
	}
}

func handleQueryTotalValidatorPower(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	bz, err := jsoniter.ConfigFastest.Marshal(keeper.GetTotalPower(ctx))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func handleQueryCurrentBatch(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	bz, err := jsoniter.ConfigFastest.Marshal(keeper.GetL1Batch(ctx))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func handleQueryCurrentValidatorSet(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	// get validator set
	validatorSet := keeper.GetValidatorSet(ctx)

	// json record
	bz, err := jsoniter.ConfigFastest.Marshal(validatorSet)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func handleQueryAllValidatorSet(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	// get all validator set
	allValidatorSet := keeper.GetAllValidators(ctx)

	// json record
	bz, err := jsoniter.ConfigFastest.Marshal(allValidatorSet)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func handleQuerySigner(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QuerySignerParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	// get validator info
	validator, err := keeper.GetValidatorInfo(ctx, params.SignerAddress)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("Error while getting validator by signer", err.Error()))
	}

	// json record
	bz, err := jsoniter.ConfigFastest.Marshal(validator)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func handleQueryValidator(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryValidatorParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	// get validator info
	validator, ok := keeper.GetValidatorFromValID(ctx, params.ValidatorID)
	if !ok {
		return nil, sdk.ErrUnknownRequest("No validator found")
	}

	// json record
	bz, err := jsoniter.ConfigFastest.Marshal(validator)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func handleQueryValidatorStatus(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QuerySignerParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	// get validator status by signer address
	status := keeper.IsCurrentValidatorByAddress(ctx, params.SignerAddress)

	// json record
	bz, err := jsoniter.ConfigFastest.Marshal(status)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func handleQueryProposer(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryProposerParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	// get validator set
	validatorSet := keeper.GetValidatorSet(ctx)

	times := int(params.Times)
	if times > len(validatorSet.Validators) {
		times = len(validatorSet.Validators)
	}

	// init proposers
	proposers := make([]hmTypes.Validator, times)

	// get proposers
	for index := 0; index < times; index++ {
		proposers[index] = *(validatorSet.GetProposer())
		validatorSet.IncrementProposerPriority(1)
	}

	// json record
	bz, err := jsoniter.ConfigFastest.Marshal(proposers)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func handleQueryCurrentProposer(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	proposer := keeper.GetCurrentProposer(ctx)

	bz, err := jsoniter.ConfigFastest.Marshal(proposer)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func handleQueryStakingSequence(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, contractCallerObj helper.IContractCaller) ([]byte, sdk.Error) {
	var params types.QueryStakingSequenceParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	chainParams := keeper.chainKeeper.GetParams(ctx)

	// get main tx receipt
	receipt, err := contractCallerObj.GetConfirmedTxReceipt(hmTypes.HexToThemisHash(params.TxHash).EthHash(), chainParams.MainchainTxConfirmations)
	if err != nil || receipt == nil {
		return nil, sdk.ErrInternal("Transaction is not confirmed yet. Please wait for sometime and try again")
	}

	// sequence id

	sequence := new(big.Int).Mul(receipt.BlockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(params.LogIndex))

	// check if incoming tx already exists
	if !keeper.HasStakingSequence(ctx, sequence.String()) {
		keeper.Logger(ctx).Error("No staking sequence exist: %s %s", params.TxHash, params.LogIndex)
		return nil, nil
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, sequence)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}
