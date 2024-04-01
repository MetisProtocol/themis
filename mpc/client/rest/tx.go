// nolint
package rest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethTypes "github.com/ethereum/go-ethereum/core/types"

	restClient "github.com/metis-seq/themis/client/rest"
	"github.com/metis-seq/themis/helper"
	"github.com/metis-seq/themis/mpc/types"
	hmTypes "github.com/metis-seq/themis/types"
	"github.com/metis-seq/themis/types/rest"
)

// It represents Propose MpcCreate msg.
//
//swagger:response ProposeMpcCreateResponse
type ProposeMpcCreateResponse struct {
	//in:body
	// Output output `json:"output"`
}

// It represents Propose MpcSign msg.
//
//swagger:response ProposeMpcSignResponse
type ProposeMpcSignResponse struct {
	//in:body
	// Output output `json:"output"`
}

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/mpc/propose-mpc-create",
		postProposeMpcCreateHandlerFn(cliCtx),
	).Methods("POST")

	r.HandleFunc(
		"/mpc/propose-mpc-sign",
		postProposeMpcSignHandlerFn(cliCtx),
	).Methods("POST")

	r.HandleFunc(
		"/mpc/mpc-sign",
		postMpcSignHandlerFn(cliCtx),
	).Methods("POST")
}

// ProposeMpcCreateReq struct for proposing new mpc
type ProposeMpcCreateReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	ID         string `json:"mpc_id"`
	Threshold  uint64 `json:"threshold"`
	Parties    string `json:"parties"`
	Proposer   string `json:"proposer"`
	MpcAddress string `json:"mpc_address"`
	MpcPubkey  string `json:"mpc_pubkey"`
	MpcType    uint64 `json:"mpc_type"`
}

// ProposeMpcSignReq struct for proposing mpc sign
type ProposeMpcSignReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	ID       string `json:"sign_id"`
	MpcID    string `json:"mpc_id"`
	SignType uint64 `json:"sign_type"`
	SignData string `json:"sign_data"`
	SignMsg  string `json:"sign_msg"`
	Proposer string `json:"proposer"`
}

type MpcSignReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	ID        string `json:"sign_id"`
	Signature string `json:"signature"`
	Proposer  string `json:"proposer"`
}

//swagger:parameters ProposeMpcCreate
type ProposeMpcCreate struct {

	//Body
	//required:true
	//in:body
	Input ReqInput `json:"input"`
}

type ReqInput struct {

	//required:true
	//in:body
	BaseReq BaseReq `json:"base_req"`

	//required:true
	//in:body
	ID uint64 `json:"mpc_id"`

	//required:true
	//in:body
	Threshold uint64 `json:"threshold"`

	//required:true
	//in:body
	Parties string `json:"parties"`
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

// swagger:route POST /mpc/propose-mpc-create
// It returns the prepared msg for proposing the mpc create
// responses:
//   200: ProposeMpcCreateResponse

func postProposeMpcCreateHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read req from request
		var req ProposeMpcCreateReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		// req.BaseReq = req.BaseReq.Sanitize()
		// if !req.BaseReq.ValidateBasic(w) {
		// 	return
		// }

		// mpc
		if req.ID == "" {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "mpc id cannot be empty")
			return
		}

		if req.Threshold <= 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid mpc threshold")
			return
		}

		if req.Parties == "" {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid mpc parties")
			return
		}

		var mpcParty []hmTypes.PartyID
		err := json.Unmarshal([]byte(req.Parties), &mpcParty)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid mpc parties")
			return
		}

		if req.Proposer == "" {
			req.Proposer = req.BaseReq.From
		}
		proposer := hmTypes.HexToThemisAddress(req.Proposer)
		if proposer.Empty() {
			proposer = helper.GetFromAddress(cliCtx)
		}

		if req.MpcPubkey == "" {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid mpc pubkey")
			return
		}
		mpcPubkey, err := hex.DecodeString(req.MpcPubkey)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid mpc pubkey")
			return
		}

		// draft a propose mpc message
		msg := types.NewMsgProposeMpcCreate(
			req.ID,
			req.Threshold,
			mpcParty,
			proposer,
			hmTypes.HexToThemisAddress(req.MpcAddress),
			mpcPubkey,
			hmTypes.MpcType(req.MpcType),
		)

		// send response
		req.BaseReq.ChainID = helper.GetGenesisDoc().ChainID
		restClient.WriteBroadcastStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// swagger:route POST /mpc/propose-mpc-sign
// It returns the prepared msg for proposing the mpc sign
// responses:
//   200: ProposeMpcSignResponse

func postProposeMpcSignHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read req from request
		var req ProposeMpcSignReq
		if !rest.ReadRESTJSONReq(w, r, &req) {
			return
		}
		RestLogger.Info("postProposeMpcSignHandlerFn req", "req", req)

		// req.BaseReq = req.BaseReq.Sanitize()
		// if !req.BaseReq.ValidateBasic(w) {
		// 	return
		// }

		proposer := hmTypes.BytesToThemisAddress(helper.GetAddress())

		if req.MpcID == "" {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "mpc id cannot be empty")
			return
		}

		// sign id
		if req.ID == "" {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid sign id")
			return
		}

		// sign type
		switch req.SignType {
		case uint64(hmTypes.BatchSubmit), uint64(hmTypes.BatchReward), uint64(hmTypes.CommitEpochToMetis), uint64(hmTypes.L1UpdateMpcAddress), uint64(hmTypes.L2UpdateMpcAddress):
		default:
			rest.WriteErrorResponse(w, http.StatusBadRequest, "unsupport sign type")
			return
		}

		// sign data
		if req.SignData == "" {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid mpc sign data")
			return
		}
		signData, err := hex.DecodeString(strings.TrimPrefix(req.SignData, "0x"))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid mpc sign data")
			return
		}
		// tx, err := parseJsonTx(req.SignData)
		// if err != nil {
		// 	rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		// 	return
		// }
		// signData, _ := tx.MarshalBinary()

		var signMsg []byte
		// sign msg
		if req.SignMsg != "" {
			signMsg, err = hex.DecodeString(strings.TrimPrefix(req.SignMsg, "0x"))
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid mpc sign msg")
				return
			}
		}

		// draft a propose sign message
		msg := types.NewMsgProposeMpcSign(
			req.ID,
			req.MpcID,
			hmTypes.SignType(req.SignType),
			signData,
			signMsg,
			proposer,
		)

		req.BaseReq.ChainID = helper.GetGenesisDoc().ChainID
		// send response
		restClient.WriteBroadcastStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

type JsonTransaction struct {
	Nonce    string `json:"nonce"`
	GasPrice string `json:"gasPrice"`
	GasLimit string `json:"gasLimit"`
	To       string `json:"to"`
	Value    string `json:"value"`
	Data     string `json:"data"`
}

func parseJsonTx(jsonTx string) (*ethTypes.Transaction, error) {
	var tx JsonTransaction
	err := json.Unmarshal([]byte(jsonTx), &tx)
	if err != nil {
		fmt.Println("Error decode json:", err)
		return nil, err
	}

	nonce, err := hexutil.DecodeBig(tx.Nonce)
	if err != nil {
		fmt.Println("Error nonce:", err)
		return nil, err
	}
	value, err := hexutil.DecodeBig(tx.Value)
	if err != nil {
		fmt.Println("Error value:", err)
		return nil, err
	}
	gasLimit, err := hexutil.DecodeBig(tx.GasLimit)
	if err != nil {
		fmt.Println("Error gasLimit:", err)
		return nil, err
	}
	gasPrice, err := hexutil.DecodeBig(tx.GasPrice)
	if err != nil {
		fmt.Println("Error gasPrice:", err)
		return nil, err
	}

	txData, err := hex.DecodeString(strings.TrimPrefix(tx.Data, "0x"))
	if err != nil {
		fmt.Println("Error DecodeString:", err)
		return nil, err
	}
	return ethTypes.NewTransaction(nonce.Uint64(), ethCommon.HexToAddress(tx.To), value, gasLimit.Uint64(), gasPrice, txData), nil
}

// swagger:route POST /mpc/mpc-create
// It returns the result of mpc sign
// responses:
//   200: MpcSignResponse

func postMpcSignHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read req from request
		var req MpcSignReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		// req.BaseReq = req.BaseReq.Sanitize()
		// if !req.BaseReq.ValidateBasic(w) {
		// 	return
		// }

		if req.Proposer == "" {
			req.Proposer = req.BaseReq.From
		}
		proposer := hmTypes.HexToThemisAddress(req.Proposer)
		if proposer.Empty() {
			proposer = helper.GetFromAddress(cliCtx)
		}

		// sign id
		if req.ID == "" {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid sign id")
			return
		}

		// signature
		if req.Signature == "" {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid sign signature")
			return
		}

		// draft a sign message
		msg := types.NewMsgMpcSign(
			req.ID,
			req.Signature,
			proposer,
		)

		// send response
		req.BaseReq.ChainID = helper.GetGenesisDoc().ChainID
		restClient.WriteBroadcastStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
