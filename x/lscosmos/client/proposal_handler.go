package client

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"

	"github.com/persistenceOne/pstake-native/x/lscosmos/client/cli"
)

var (
	RegisterHostChainProposalHandler      = govclient.NewProposalHandler(cli.NewRegisterHostChainCmd, emptyRestHandler)
	MinDepositAndFeeChangeProposalHandler = govclient.NewProposalHandler(cli.NewMinDepositAndFeeChangeCmd, emptyRestHandler)
)

func emptyRestHandler(client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "unsupported-lscsmos-client",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Legacy REST Routes are not supported for lscosmos proposals")
		},
	}
}
