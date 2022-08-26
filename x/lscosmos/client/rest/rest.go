package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	restClient "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/gorilla/mux"

	"github.com/persistenceOne/pstake-native/x/lscosmos/client/utils"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

type SendReq struct {
	BaseReq          rest.BaseReq `json:"base_req" yaml:"base_req"`
	DelegatorAddress string       `json:"delegator_address" yaml:"delegator_address"`
	Amount           sdk.Coin     `json:"amount" yaml:"amount"`
}

func RegisterHandlers(clientCtx client.Context, rtr *mux.Router) {
	r := restClient.WithHTTPDeprecationHeaders(rtr)
	r.HandleFunc("/lscosmos/liquidstake", LiquidStakeHandlerFn(clientCtx)).Methods("POST")
}

// LiquidStakeHandlerFn returned an HTTP REST handler for creating a MsgLiquidStake
func LiquidStakeHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req SendReq
		if !rest.ReadRESTReq(writer, request, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(writer) {
			return
		}

		delegatorAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddress)
		if rest.CheckBadRequestError(writer, err) {
			return
		}

		msg := types.NewMsgLiquidStake(req.Amount, delegatorAddr)
		tx.WriteGeneratedTxResponse(clientCtx, writer, req.BaseReq, msg)

	}
}

// RegisterChainRESTHandler returns a ProposalRESTHandler that exposes the param
// change REST handler with a given sub-route.
func RegisterChainRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "register_chain",
		Handler:  postRegisterChainHandlerFn(clientCtx),
	}
}

func postRegisterChainHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req utils.RegisterCosmosChainProposalReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if req.BaseReq.ValidateBasic(w) {
			return
		}

		minDeposit, ok := sdk.NewIntFromString(req.MinDeposit)
		if !ok {
			_ = rest.CheckBadRequestError(w, types.ErrInvalidIntParse)
			return
		}
		depositFee, err := sdk.NewDecFromStr(req.PStakeDepositFee)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		restakeFee, err := sdk.NewDecFromStr(req.PStakeRestakeFee)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		unstakeFee, err := sdk.NewDecFromStr(req.PStakeUnstakeFee)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		content := types.NewRegisterCosmosChainProposal(
			req.Title,
			req.Description,
			req.ModuleEnabled,
			req.IBCConnection,
			req.TokenTransferChannel,
			req.TokenTransferPort,
			req.BaseDenom,
			req.MintDenom,
			minDeposit,
			req.AllowListedValidators,
			depositFee,
			restakeFee,
			unstakeFee,
		)

		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, req.Proposer)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
