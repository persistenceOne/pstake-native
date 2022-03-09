package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	"net/http"
)

// SendReq defines the properties of a send request's body.
type SendReq struct {
	BaseReq            rest.BaseReq   `json:"base_req" yaml:"base_req"`
	ChainID            string         `json:"chain_id" yaml:"chain_id"`
	TxHash             string         `json:"tx_hash" yaml:"tx_hash"`
	BlockHeight        int64          `json:"block_height" yaml:"block_height"`
	DestinationAddress sdk.AccAddress `json:"destination_address" yaml:"destination_address"`
	Amount             sdk.Coins      `json:"amount" yaml:"amount"`
}

// NewSendRequestHandlerFn returns an HTTP REST handler for creating a MsgMint
// transaction.
func NewMintRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32Addr := vars["address"]

		toAddr, err := sdk.AccAddressFromBech32(bech32Addr)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		var req SendReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		orchestrator, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		msg := cosmosTypes.NewMsgMintTokensForAccount(toAddr, orchestrator, req.Amount, req.ChainID, req.TxHash, req.BlockHeight)
		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
