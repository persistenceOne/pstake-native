/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package rest

import (
	"fmt"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/x/cosmos/types"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	restClient "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/gorilla/mux"
	"github.com/persistenceOne/pstake-native/x/cosmos/client/utils"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

func RegisterHandlers(clientCtx client.Context, rtr *mux.Router) {
	r := restClient.WithHTTPDeprecationHeaders(rtr)
	registerQueryRoutes(clientCtx, r)
	r.HandleFunc(fmt.Sprintf("/cosmos/incoming/minting"), NewMintRequestHandlerFn(clientCtx)).Methods("POST")
}

// EnableModuleProposalRESTHandler returns a EnableModuleProposalRESTHandler that exposes the community pool spend REST handler with a given sub-route.
func EnableModuleProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "module_enable",
		Handler:  postEnableModuleProposalHandlerFn(clientCtx),
	}
}

func postEnableModuleProposalHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req utils.EnableModuleProposalReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		content := types.NewEnableModuleProposal(
			req.EnableModule.Title,
			req.EnableModule.Description,
			req.EnableModule.Threshold,
			req.EnableModule.AccountNumber)

		deposit, err := sdkTypes.ParseCoinsNormalized(req.EnableModule.Deposit)
		if err != nil {
			return
		}

		depositor, err := sdkTypes.AccAddressFromBech32(req.EnableModule.Depositor)
		if err != nil {
			return
		}

		msg, err := govtypes.NewMsgSubmitProposal(content, deposit, depositor)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

// ChangeMultisigProposalRESTHandler returns a ChangeMultisigProposalRESTHandler that exposes the community pool spend REST handler with a given sub-route.
func ChangeMultisigProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "change_multisig",
		Handler:  postChangeMultisigProposalHandlerFn(clientCtx),
	}
}

func postChangeMultisigProposalHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req utils.ChangeMultisigPropsoalReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		content := types.NewChangeMultisigProposal(req.ChangeMultisig.Title,
			req.ChangeMultisig.Description,
			req.ChangeMultisig.Threshold,
			req.ChangeMultisig.OrchestratorAddresses,
			req.ChangeMultisig.AccountNumber)

		//TODO : check if correct way to do it
		deposit, err := sdkTypes.ParseCoinsNormalized(req.ChangeMultisig.Deposit)
		if err != nil {
			fmt.Println(err)
			return
		}

		depositor, err := sdkTypes.AccAddressFromBech32(req.ChangeMultisig.Depositor)
		if err != nil {
			fmt.Println(err)
			return
		}

		msg, err := govtypes.NewMsgSubmitProposal(content, deposit, depositor)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

// ChangeCosmosValidatorWeightsProposalRESTHandler returns a ChangeCosmosValidatorWeightsProposalRESTHandler that exposes the community pool spend REST handler with a given sub-route.
func ChangeCosmosValidatorWeightsProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "change_cosmos_validator_weights",
		Handler:  postChangeCosmosValidatorWeightsProposalHandlerFn(clientCtx),
	}
}

func postChangeCosmosValidatorWeightsProposalHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req utils.ChangeCosmosValidatorWeightsProposalReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		var weightedAddresses []cosmosTypes.WeightedAddressAmount

		for _, weightedAddress := range req.CosmosValidatorSet.WeightedAddresses {
			weight, err := sdkTypes.NewDecFromStr(weightedAddress.Weight)
			if err != nil {
				return
			}
			weightedAddresses = append(
				weightedAddresses,
				cosmosTypes.WeightedAddressAmount{
					Address: weightedAddress.ValAddress,
					Weight:  weight,
				})
		}

		content := types.NewChangeCosmosValidatorWeightsProposal(
			req.CosmosValidatorSet.Title,
			req.CosmosValidatorSet.Description,
			weightedAddresses)

		deposit, err := sdkTypes.ParseCoinsNormalized(req.CosmosValidatorSet.Deposit)
		if err != nil {
			return
		}

		depositor, err := sdkTypes.AccAddressFromBech32(req.CosmosValidatorSet.Depositor)
		if err != nil {
			return
		}

		msg, err := govtypes.NewMsgSubmitProposal(content, deposit, depositor)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

// ChangeOracleValidatorWeightsProposalRESTHandler returns a ChangeOracleValidatorWeightsProposalRESTHandler that exposes the community pool spend REST handler with a given sub-route.
func ChangeOracleValidatorWeightsProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "change_cosmos_validator_weights",
		Handler:  postChangeOracleValidatorWeightsProposalHandlerFn(clientCtx),
	}
}

func postChangeOracleValidatorWeightsProposalHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req utils.ChangeOracleValidatorWeightsProposalReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		var weightedAddresses []cosmosTypes.WeightedAddress

		for _, weightedAddress := range req.OracleValidatorSet.WeightedAddresses {
			weight, err := sdkTypes.NewDecFromStr(weightedAddress.Weight)
			if err != nil {
				return
			}
			weightedAddresses = append(
				weightedAddresses,
				cosmosTypes.WeightedAddress{
					Address: weightedAddress.ValAddress,
					Weight:  weight,
				})
		}

		content := types.NewChangeOracleValidatorWeightsProposal(
			req.OracleValidatorSet.Title,
			req.OracleValidatorSet.Description,
			weightedAddresses)

		deposit, err := sdkTypes.ParseCoinsNormalized(req.OracleValidatorSet.Deposit)
		if err != nil {
			return
		}

		depositor, err := sdkTypes.AccAddressFromBech32(req.OracleValidatorSet.Depositor)
		if err != nil {
			return
		}

		msg, err := govtypes.NewMsgSubmitProposal(content, deposit, depositor)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
