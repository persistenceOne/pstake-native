/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/persistenceOne/pstake-native/x/cosmos/types"
)

func registerQueryRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc(
		"/cosmos/parameters",
		queryParamsHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/cosmos/txByID/{txID}",
		queryTxByIDHandlerFn(cliCtx),
	).Methods("GET")

}

func queryParamsHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParameters)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryTxByIDHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		endpoint := fmt.Sprintf("custom/%s", types.QuerierRoute)
		vars := mux.Vars(r)
		txID, err := strconv.ParseUint(vars["txID"], 10, 64)
		if err != nil {
			return
		}
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}
		query := &types.QueryOutgoingTxByIDRequest{
			TxID: txID,
		}
		bz, err := clientCtx.LegacyAmino.MarshalJSON(query)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		res, height, err := clientCtx.QueryWithData(endpoint, bz)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}
