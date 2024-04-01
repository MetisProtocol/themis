package mpc

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	jsoniter "github.com/json-iterator/go"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/metis-seq/themis/mpc/types"
	hmTypes "github.com/metis-seq/themis/types"
)

// NewQuerier creates a querier for auth REST endpoints
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryMpc:
			return handleQueryMpc(ctx, req, keeper)
		case types.QueryMpcList:
			return handleQueryMpcList(ctx, req, keeper)
		case types.QueryLatestMpc:
			return handleQueryLatestMpc(ctx, req, keeper)
		case types.QueryMpcSet:
			return handleQueryMpcSet(ctx, req, keeper)
		case types.QueryMpcSign:
			return handleQueryMpcSign(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown auth query endpoint")
		}
	}
}

func handleQueryMpc(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryMpcParams

	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	mpc, err := keeper.GetMpc(ctx, params.MpcID)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not get mpc", err.Error()))
	}

	// return error if mpc doesn't exist
	if mpc == nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("mpc %v does not exist", params.MpcID))
	}

	// json record
	bz, err := jsoniter.ConfigFastest.Marshal(mpc)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func handleQueryMpcList(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params hmTypes.QueryPaginationParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	res, err := keeper.GetMpcList(ctx, params.Page, params.Limit)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr(fmt.Sprintf("could not fetch mpc list with page %v and limit %v", params.Page, params.Limit), err.Error()))
	}

	bz, err := jsoniter.ConfigFastest.Marshal(res)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func handleQueryLatestMpc(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var defaultMpc hmTypes.Mpc

	mpcs := keeper.GetAllMpcs(ctx)
	if len(mpcs) == 0 {
		// json record
		bz, err := jsoniter.ConfigFastest.Marshal(defaultMpc)
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}

		return bz, nil
	}

	var params types.QueryLatestMpcParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	keeper.Logger(ctx).Info("handleQueryLatestMpc req", "reqData", string(req.Data), "prams", params)

	// explcitly fetch the last mpc
	mpc, err := keeper.GetLastMpc(ctx, hmTypes.MpcType(params.Type))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not get mpc", err.Error()))
	}

	// return error if mpc doesn't exist
	if mpc == nil {
		return nil, sdk.ErrInternal("latest mpc does not exist")
	}

	// json record
	bz, err := jsoniter.ConfigFastest.Marshal(mpc)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func handleQueryMpcSet(ctx sdk.Context, _ abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var defaultMpcSet []*hmTypes.Validator

	mpcSets := keeper.GetAllMpcSets(ctx)
	if len(mpcSets) == 0 {
		// json record
		bz, err := jsoniter.ConfigFastest.Marshal(defaultMpcSet)
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}

		return bz, nil
	}

	// json record
	bz, err := jsoniter.ConfigFastest.Marshal(mpcSets)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func handleQueryMpcSign(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryMpcSignParams

	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	sign, err := keeper.GetMpcSign(ctx, params.SignID)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not get sign", err.Error()))
	}

	// return error if sign doesn't exist
	if sign == nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("sign %v does not exist", params.SignID))
	}

	// json record
	bz, err := jsoniter.ConfigFastest.Marshal(sign)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}
