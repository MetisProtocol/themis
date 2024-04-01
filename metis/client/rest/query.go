// Package classification ThemisRest API
//
//	    Schemes: http
//	    BasePath: /
//	    Version: 0.0.1
//	    title: Themis APIs
//	    Consumes:
//	    - application/json
//		   Host:localhost:1317
//	    - application/json
//
// nolint
//
//swagger:meta
package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"

	"github.com/metis-seq/themis/helper"
	"github.com/metis-seq/themis/metis/types"
	stakingTypes "github.com/metis-seq/themis/staking/types"
	hmTypes "github.com/metis-seq/themis/types"
	hmRest "github.com/metis-seq/themis/types/rest"
)

type ThemisSpanResultWithHeight struct {
	Height int64
	Result []byte
}

type validator struct {
	ID           int    `json:"ID"`
	StartBatch   int    `json:"startBatch"`
	EndBatch     int    `json:"endBatch"`
	Nonce        int    `json:"nonce"`
	Power        int    `json:"power"`
	PubKey       string `json:"pubKey"`
	Signer       string `json:"signer"`
	Last_Updated string `json:"last_updated"`
	Jailed       bool   `json:"jailed"`
	Accum        int    `json:"accum"`
}

type span struct {
	SpanID     int `json:"span_id"`
	StartBlock int `json:"start_block"`
	EndBlock   int `json:"end_block"`
	//in:body
	ValidatorSet      validatorSet `json:"validator_set"`
	SelectedProducers []validator  `json:"selected_producer"`
	MetisChainId      string       `json:"metis_chain_id"`
}

type validatorSet struct {
	Validators []validator `json:"validators"`
	Proposer   validator   `json:"Proposer"`
}

// It represents the list of spans
//
//swagger:response metisSpanListResponse
type metisSpanListResponse struct {
	//in:body
	Output metisSpanList `json:"output"`
}

type metisSpanList struct {
	Height string `json:"height"`
	Result []span `json:"result"`
}

// It represents the span
//
//swagger:response metisSpanResponse
type metisSpanResponse struct {
	//in:body
	Output metisSpan `json:"output"`
}

type metisSpan struct {
	Height string `json:"height"`
	Result span   `json:"result"`
}

// It represents the metis span parameters
//
//swagger:response metisSpanParamsResponse
type metisSpanParamsResponse struct {
	//in:body
	Output metisSpanParams `json:"output"`
}

type metisSpanParams struct {
	Height string     `json:"height"`
	Result spanParams `json:"result"`
}

type spanParams struct {

	//type:integer
	SprintDuration int64 `json:"sprint_duration"`
	//type:integer
	SpanDuration int64 `json:"span_duration"`
	//type:integer
	ProducerCount int64 `json:"producer_count"`
}

// It represents the next span seed
//
//swagger:response metisNextSpanSeedRespose
type metisNextSpanSeedRespose struct {
	//in:body
	Output spanSeed `json:"output"`
}

type spanSeed struct {
	Height string `json:"height"`
	Result string `json:"result"`
}

var spanOverrides map[uint64]*ThemisSpanResultWithHeight = nil

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/metis/span/list", spanListHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/metis/respan/list", respanListHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/metis/span/{id}", spanHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/metis/latest-span", latestSpanHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/metis/prepare-next-span", prepareNextSpanHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/metis/next-span-seed", fetchNextSpanSeedHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/metis/params", paramsHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/metis/tx/{id}", metisTxHandlerFn(cliCtx)).Methods("GET")
}

// swagger:route GET /metis/next-span-seed metis metisNextSpanSeed
// It returns the seed for the next span
// responses:
//   200: metisNextSpanSeedRespose

func fetchNextSpanSeedHandlerFn(
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryNextSpanSeed), nil)
		if err != nil {
			RestLogger.Error("Error while fetching next span seed  ", "Error", err.Error())
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		RestLogger.Debug("nextSpanSeed querier response", "res", res)

		// error if span seed found
		if !hmRest.ReturnNotFoundIfNoContent(w, res, "NextSpanSeed not found") {
			RestLogger.Error("NextSpanSeed not found ", "Error", err.Error())
			return
		}

		// return result
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//swagger:parameters metisSpanList
type metisSpanListParam struct {

	//Page Number
	//required:true
	//type:integer
	//in:query
	Page int `json:"page"`

	//Limit
	//required:true
	//type:integer
	//in:query
	Limit int `json:"limit"`
}

// swagger:route GET /metis/span/list metis metisSpanList
// It returns the list of Metis Span
// responses:
//
//	200: metisSpanListResponse
func spanListHandlerFn(
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := r.URL.Query()

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// get page
		page, ok := rest.ParseUint64OrReturnBadRequest(w, vars.Get("page"))
		if !ok {
			return
		}

		// get limit
		limit, ok := rest.ParseUint64OrReturnBadRequest(w, vars.Get("limit"))
		if !ok {
			return
		}

		// get query params
		queryParams, err := cliCtx.Codec.MarshalJSON(hmTypes.NewQueryPaginationParams(page, limit))
		if err != nil {
			return
		}

		// query spans
		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySpanList), queryParams)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		RestLogger.Debug("span-list querier response", "res", res)

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "No spans found"); !ok {
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// swagger:route GET /metis/respan/list metis metisReSpanList
// It returns the list of Metis ReSpan
// responses:
//
//	200: metisSpanListResponse
func respanListHandlerFn(
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := r.URL.Query()

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// get page
		page, ok := rest.ParseUint64OrReturnBadRequest(w, vars.Get("page"))
		if !ok {
			return
		}

		// get limit
		limit, ok := rest.ParseUint64OrReturnBadRequest(w, vars.Get("limit"))
		if !ok {
			return
		}

		// get query params
		queryParams, err := cliCtx.Codec.MarshalJSON(hmTypes.NewQueryPaginationParams(page, limit))
		if err != nil {
			return
		}

		// query spans
		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryReSpanList), queryParams)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		RestLogger.Debug("respan-list querier response", "res", res)

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "No spans found"); !ok {
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//swagger:parameters metisSpanById
type metisSpanById struct {

	//Id number of the span
	//required:true
	//type:integer
	//in:path
	Id int `json:"id"`
}

// swagger:route GET /metis/span/{id} metis metisSpanById
// It returns the span based on ID
// responses:
//
//	200: metisSpanResponse
func spanHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)

		// get to address
		spanID, ok := rest.ParseUint64OrReturnBadRequest(w, vars["id"])
		if !ok {
			return
		}

		var (
			res            []byte
			height         int64
			spanOverridden bool
		)

		if spanOverrides == nil {
			loadSpanOverrides()
		}

		if span, ok := spanOverrides[spanID]; ok {
			res = span.Result
			height = span.Height
			spanOverridden = true
		}

		if !spanOverridden {
			// get query params
			queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQuerySpanParams(spanID))
			if err != nil {
				return
			}

			// fetch span
			res, height, err = cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySpan), queryParams)
			if err != nil {
				hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "No span found"); !ok {
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		hmRest.PostProcessResponse(w, cliCtx, res)
	}
}

//swagger:parameters metisTxById
type metisTxById struct {

	//Id number of the span
	//required:true
	//type:integer
	//in:path
	Id string `json:"id"`
}

// swagger:route GET /metis/tx/{id} metis metisTxById
// It returns the tx based on ID
// responses:
//
//	200: metisTxResponse
func metisTxHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		txId := vars["id"]

		// get query params
		queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryMetisTxParams(txId))
		if err != nil {
			return
		}

		// fetch span
		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryMetisTx), queryParams)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		hmRest.PostProcessResponse(w, cliCtx, res)
	}
}

// swagger:route GET /metis/latest-span metis metisSpanLatest
// It returns the latest-span
// responses:
//
//	200: metisSpanResponse
func latestSpanHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// fetch latest span
		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLatestSpan), nil)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "No latest span found"); !ok {
			return
		}

		// return result
		cliCtx = cliCtx.WithHeight(height)
		hmRest.PostProcessResponse(w, cliCtx, res)
	}
}

//swagger:parameters metisPrepareNextSpan
type metisPrepareNextSpanParam struct {

	//Start Block
	//required:true
	//type:integer
	//in:query
	StartBlock int `json:"start_block"`

	//Span ID of the span
	//required:true
	//type:integer
	//in:query
	SpanId int `json:"span_id"`

	//Chain ID of the network
	//required:true
	//type:integer
	//in:query
	ChainId int `json:"chain_id"`
}

// swagger:route GET /metis/prepare-next-span metis metisPrepareNextSpan
// It returns the prepared next span
// responses:
//
//	200: metisSpanResponse
func prepareNextSpanHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := r.URL.Query()

		spanID, ok := rest.ParseUint64OrReturnBadRequest(w, params.Get("span_id"))
		if !ok {
			return
		}

		startBlock, ok := rest.ParseUint64OrReturnBadRequest(w, params.Get("start_block"))
		if !ok {
			return
		}

		chainID := params.Get("chain_id")

		//
		// Get span duration
		//

		// fetch duration
		spanDurationBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryParams, types.ParamSpan), nil)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, spanDurationBytes, "No span duration"); !ok {
			return
		}

		var spanDuration uint64
		if err := jsoniter.ConfigFastest.Unmarshal(spanDurationBytes, &spanDuration); err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		validatorSetBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", stakingTypes.QuerierRoute, stakingTypes.QueryCurrentValidatorSet), nil)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// check content
		if !hmRest.ReturnNotFoundIfNoContent(w, validatorSetBytes, "No current validator set found") {
			return
		}

		var _validatorSet hmTypes.ValidatorSet
		if err = jsoniter.ConfigFastest.Unmarshal(validatorSetBytes, &_validatorSet); err != nil {
			hmRest.WriteErrorResponse(w, http.StatusNoContent, errors.New("unable to unmarshall JSON").Error())
			return
		}

		//
		// Fetching SelectedProducers
		//

		nextProducerBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryNextProducers), nil)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, nextProducerBytes, "Next Producers not found"); !ok {
			return
		}

		var selectedProducers []hmTypes.Validator
		if err := jsoniter.ConfigFastest.Unmarshal(nextProducerBytes, &selectedProducers); err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		selectedProducers = hmTypes.SortValidatorByAddress(selectedProducers)

		// draft a propose span message
		msg := hmTypes.NewSpan(
			spanID,
			startBlock,
			startBlock+spanDuration-1,
			_validatorSet,
			selectedProducers,
			chainID,
		)

		result, err := jsoniter.ConfigFastest.Marshal(&msg)
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		hmRest.PostProcessResponse(w, cliCtx, result)
	}
}

// swagger:route GET /metis/params metis metisSpanParams
// It returns the span parameters
// responses:
//
//	200: metisSpanParamsResponse
func paramsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParams)

		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func loadSpanOverrides() {
	spanOverrides = map[uint64]*ThemisSpanResultWithHeight{}

	j, ok := SPAN_OVERRIDES[helper.GenesisDoc.ChainID]
	if !ok {
		return
	}

	var spans []*types.ResponseWithHeight
	if err := jsoniter.ConfigFastest.Unmarshal(j, &spans); err != nil {
		return
	}

	for _, span := range spans {
		var themisSpan types.ThemisSpan
		if err := jsoniter.ConfigFastest.Unmarshal(span.Result, &themisSpan); err != nil {
			continue
		}

		height, err := strconv.ParseInt(span.Height, 10, 64)
		if err != nil {
			continue
		}

		spanOverrides[themisSpan.ID] = &ThemisSpanResultWithHeight{
			Height: height,
			Result: span.Result,
		}
	}
}

//swagger:parameters metisSpanList metisSpanById metisPrepareNextSpan metisSpanLatest metisSpanParams metisNextSpanSeed
type Height struct {

	//Block Height
	//in:query
	Height string `json:"height"`
}
