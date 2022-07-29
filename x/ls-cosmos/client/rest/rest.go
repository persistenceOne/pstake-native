package rest

import (
	restClient "github.com/cosmos/cosmos-sdk/client/rest"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/persistenceOne/pstake-native/x/ls-cosmos/client/utils"
	"github.com/persistenceOne/pstake-native/x/ls-cosmos/types"
)

type SendReq struct {
	BaseReq        rest.BaseReq `json:"base_req" yaml:"base_req"`
	MintAddress    string       `json:"mint_address" yaml:"mint_address"`
	Amount         sdk.Coin     `json:"amount" yaml:"amount"`
	DepositAddress string       `json:"deposit_address" yaml:"deposit_address"`
}

func RegisterHandlers(clientCtx client.Context, rtr *mux.Router) {
	r := restClient.WithHTTPDeprecationHeaders(rtr)
	r.HandleFunc("/ls-cosmos/incoming/minting", LiquidStakeHandlerFn(clientCtx)).Methods("POST")
}

//	LiquidStakeHandlreFn returnd an HTTP REST handler for creating a MsgLiquidStake
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

		mintAddr, err := sdk.AccAddressFromBech32(req.MintAddress)
		if rest.CheckBadRequestError(writer, err) {
			return
		}

		depositAddr, err := sdk.AccAddressFromBech32(req.DepositAddress)
		if rest.CheckBadRequestError(writer, err) {
			return
		}

		msg := types.NewMsgLiquidStake(req.Amount, mintAddr, depositAddr)
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

		content := types.NewRegisterCosmosChainProposal(
			req.Title,
			req.Description,
			req.IBCConnection,
			req.TokenTransferChannel,
			req.TokenTransferPort,
			req.BaseDenom,
			req.MintDenom,
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
