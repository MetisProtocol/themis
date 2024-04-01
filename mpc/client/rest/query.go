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
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/metis-seq/themis/mpc/types"
	hmRest "github.com/metis-seq/themis/types/rest"
)

// It represents the list of mpcs
//
//swagger:response mpcListResponse
type mpcListResponse struct {
	//in:body
	Output mpcList `json:"output"`
}

type mpcList struct {
	Height string `json:"height"`
	Result []mpc  `json:"result"`
}

// It represents the mpc
//
//swagger:response mpcResponse
type mpcResponse struct {
	//in:body
	Output mpc `json:"output"`
}

type mpc struct {
	Height string `json:"height"`
	Result Mpc    `json:"result"`
}

type Mpc struct {
}

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/mpc/list", mpcListHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/mpc/latest/{type}", latestMpcHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/mpc/set", mpcSetHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/mpc/{id}", mpcHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/mpc/sign/{id}", mpcSignHandlerFn(cliCtx)).Methods("GET")
}

// swagger:route GET mpc/list mpcList
// It returns the list of Metis Mpc
// responses:
//
//	200: mpcListResponse
func mpcListHandlerFn(
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// query mpcs
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryMpcList), nil)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "No mpc found"); !ok {
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//swagger:parameters mpcById
type mpcById struct {

	//Id number of the mpc
	//required:true
	//type:integer
	//in:path
	Id int `json:"id"`
}

// swagger:route GET /mpc/{id} mpcById
// It returns the mpc based on ID
// responses:
//
//	200: mpcResponse
func mpcHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		mpcID := vars["id"]

		var (
			res    []byte
			height int64
		)

		// get query params
		queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryMpcParams(mpcID))
		if err != nil {
			return
		}

		// fetch mpc
		res, height, err = cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryMpc), queryParams)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "No mpc found"); !ok {
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		hmRest.PostProcessResponse(w, cliCtx, res)
	}
}

//swagger:parameters latestMpcByType
type latestMpcByType struct {

	//type of the mpc
	//required:true
	//type:string
	//in:path
	Type int `json:"type"`
}

// swagger:route GET /mpc/latest-mpc/{type} mpcLatest
// It returns the latest-mpc
// responses:
//
//	200: mpcResponse
func latestMpcHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		mpcType := vars["type"]
		mpcTypeInt, _ := strconv.ParseInt(mpcType, 10, 64)
		// get query params
		queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryLatestMpcParams(int(mpcTypeInt)))
		if err != nil {
			return
		}

		// fetch latest mpc
		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLatestMpc), queryParams)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "No latest mpc found"); !ok {
			return
		}

		// return result
		cliCtx = cliCtx.WithHeight(height)
		hmRest.PostProcessResponse(w, cliCtx, res)
	}
}

// swagger:route GET /mpc/mpc-set mpcSet
// It returns the mpc parties set
// responses:
//
//	200: mpcSetResponse
func mpcSetHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// fetch  mpc set
		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryMpcSet), nil)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "No mpc set found"); !ok {
			return
		}

		// return result
		cliCtx = cliCtx.WithHeight(height)
		hmRest.PostProcessResponse(w, cliCtx, res)
	}
}

// swagger:route GET /mpc/sign/{id} mpSignById
// It returns the mpc based on ID
// responses:
//
//	200: mpcSignResponse
func mpcSignHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		mpcSignID := vars["id"]

		var (
			res    []byte
			height int64
		)

		// get query params
		queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryMpcSignParams(mpcSignID))
		if err != nil {
			return
		}

		// fetch mpc sign
		res, height, err = cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryMpcSign), queryParams)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "No mpc sign found"); !ok {
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		hmRest.PostProcessResponse(w, cliCtx, res)
	}
}
