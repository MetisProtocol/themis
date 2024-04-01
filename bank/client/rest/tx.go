// nolint
package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"

	bankTypes "github.com/metis-seq/themis/bank/types"
	restClient "github.com/metis-seq/themis/client/rest"
	"github.com/metis-seq/themis/types"
	"github.com/metis-seq/themis/types/rest"
)

// It represents transfer msg.
//
//swagger:response bankBalanceTransferResponse
type bankBalanceTransferResponse struct {
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
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`
	Amount      []coin `json:"amount"`
}

type coin struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/bank/accounts/{address}/transfers", SendRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/bank/balances/{address}", QueryBalancesRequestHandlerFn(cliCtx)).Methods("GET")
}

// SendReq defines the properties of a send request's body.
type SendReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

	Amount sdk.Coins `json:"amount" yaml:"amount"`
}

//swagger:parameters bankBalanceTransfer
type bankBalanceTransfer struct {

	//Address of the account
	//required:true
	//in:path
	Address string `json:"address"`

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
	Amount []coin `json:"amount"`
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

// swagger:route POST /bank/accounts/{address}/transfers bank bankBalanceTransfer
// It returns the prepared msg for the transfer of balance from one account to another.
// responses:
//
//	200: bankBalanceTransferResponse
//
// SendRequestHandlerFn - http request handler to send coins to a address.
func SendRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		// get to address
		toAddr := types.HexToThemisAddress(vars["address"])

		var req SendReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// get from address
		fromAddr := types.HexToThemisAddress(req.BaseReq.From)

		msg := bankTypes.NewMsgSend(fromAddr, toAddr, req.Amount)
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
