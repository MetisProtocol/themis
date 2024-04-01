// nolint
package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"

	"github.com/ethereum/go-ethereum/common"

	restClient "github.com/metis-seq/themis/client/rest"
	"github.com/metis-seq/themis/metis/types"
	hmTypes "github.com/metis-seq/themis/types"
	"github.com/metis-seq/themis/types/rest"
)

// It represents Propose Span msg.
//
//swagger:response metisProposeSpanResponse
type metisProposeSpanResponse struct {
	//in:body
	Output output `json:"output"`
}

type output struct {
	Type  string `json:"type"`
	Value value  `json:"value"`
}

type value struct {
	Msg       msg    `json:"msg"`
	Signature string `json:"signature"`
	Memo      string `json:"memo"`
}

type msg struct {
	Type  string `json:"type"`
	Value val    `json:"value"`
}

type val struct {
	SpanID       string `json:"span_id"`
	Proposer     string `json:"proposer"`
	StartBlock   string `json:"start_block"`
	EndBlock     string `json:"end_block"`
	MetisChainId string `json:"metis_chain_id"`
	Seed         string `json:"seed"`
}

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/metis/propose-span",
		postProposeSpanHandlerFn(cliCtx),
	).Methods("POST")
}

// ProposeSpanReq struct for proposing new span
type ProposeSpanReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	ID           uint64 `json:"span_id"`
	StartBlock   uint64 `json:"start_block"`
	MetisChainID string `json:"metis_chain_id"`
	L2EpochID    uint64 `json:"l2_epoch_id"`
}

//swagger:parameters metisProposeSpan
type metisProposeSpan struct {

	//Body
	//required:true
	//in:body
	Input SendReqInput `json:"input"`
}

type SendReqInput struct {

	//required:true
	//in:body
	BaseReq BaseReq `json:"base_req"`

	//required:true
	//in:body
	ID uint64 `json:"span_id"`

	//required:true
	//in:body
	StartBlock uint64 `json:"start_block"`

	//required:true
	//in:body
	MetisChainID string `json:"metis_chain_id"`
}

type BaseReq struct {

	//Address of the sender
	//required:true
	//in:body
	From string `json:"address"`

	//Chain ID of Themis
	//required:true
	//in:body
	ChainID string `json:"chain_id"`
}

// swagger:route POST /metis/propose-span metis metisProposeSpan
// It returns the prepared msg for proposing the span
// responses:
//   200: metisProposeSpanResponse

func postProposeSpanHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read req from request
		var req ProposeSpanReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		//
		// Get span duration
		//

		// fetch duration
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryParams, types.ParamSpan), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, errors.New("Span duration not found ").Error())
			return
		}

		var spanDuration uint64
		if err = jsoniter.ConfigFastest.Unmarshal(res, &spanDuration); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// fetch seed
		res, _, err = cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryNextSpanSeed), nil)
		if err != nil {
			RestLogger.Error("Error while fetching next span seed  ", "Error", err.Error())
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		var seed common.Hash
		if err = jsoniter.ConfigFastest.Unmarshal(res, &seed); err != nil {
			return
		}

		// draft a propose span message
		msg := types.NewMsgProposeSpan(
			req.ID,
			hmTypes.HexToThemisAddress(req.BaseReq.From),
			req.L2EpochID,
			req.StartBlock,
			req.StartBlock+spanDuration-1,
			req.MetisChainID,
			seed,
			false,
			hmTypes.HexToThemisAddress(req.BaseReq.From),
			nil,
		)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
