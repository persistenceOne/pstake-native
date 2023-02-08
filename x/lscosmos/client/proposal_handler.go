package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/persistenceOne/pstake-native/x/lscosmos/client/cli"
)

var (
	RegisterHostChainProposalHandler           = govclient.NewProposalHandler(cli.NewRegisterHostChainCmd)
	MinDepositAndFeeChangeProposalHandler      = govclient.NewProposalHandler(cli.NewMinDepositAndFeeChangeCmd)
	PstakeFeeAddressChangeProposalHandler      = govclient.NewProposalHandler(cli.NewPstakeFeeAddressChangeCmd)
	AllowListValidatorSetChangeProposalHandler = govclient.NewProposalHandler(cli.NewAllowListedValidatorSetChangeProposalCmd)
)

//func emptyRestHandler(client.Context) govrest.ProposalRESTHandler {
//	return govrest.ProposalRESTHandler{
//		SubRoute: "unsupported-lscsmos-client",
//		Handler: func(w http.ResponseWriter, r *http.Request) {
//			rest.WriteErrorResponse(w, http.StatusBadRequest, "Legacy REST Routes are not supported for lscosmos proposals")
//		},
//	}
//}
