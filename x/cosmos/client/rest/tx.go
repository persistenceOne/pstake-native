package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

// SendReq defines the properties of a send request's body.
type SendReq struct {
	BaseReq             rest.BaseReq `json:"base_req"`
	AddressFromMemo     string       `json:"address_from_memo" yaml:"address_from_memo"`
	OrchestratorAddress string       `json:"orchestrator_address" yaml:"orchestrator_address"`
	Amount              sdk.Coins    `json:"amount" yaml:"amount"`
	TxHash              string       `json:"tx_hash" yaml:"tx_hash"`
	ChainID             string       `json:"chain_id" yaml:"chain_id"`
	BlockHeight         int64        `json:"block_height" yaml:"block_height"`
}

// NewSendRequestHandlerFn returns an HTTP REST handler for creating a MsgMint
// transaction.
func NewMintRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req SendReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}
		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		toAddr, err := sdk.AccAddressFromBech32(req.AddressFromMemo)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		orchestrator, err := sdk.AccAddressFromBech32(req.OrchestratorAddress)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		msg := cosmosTypes.NewMsgMintTokensForAccount(toAddr, orchestrator, req.Amount, req.TxHash, req.ChainID, req.BlockHeight)
		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
