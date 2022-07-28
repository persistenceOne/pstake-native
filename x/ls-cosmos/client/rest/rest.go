package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/persistenceOne/pstake-native/x/ls-cosmos/client/utils"
	"github.com/persistenceOne/pstake-native/x/ls-cosmos/types"
)

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
