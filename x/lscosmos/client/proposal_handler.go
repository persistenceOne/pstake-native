package client

import (
	"github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/persistenceOne/pstake-native/x/lscosmos/client/cli"
	"github.com/persistenceOne/pstake-native/x/lscosmos/client/rest"
)

var (
	RegisterCosmosChainProposalHandler = client.NewProposalHandler(cli.NewRegisterCosmosChainCmd, rest.RegisterChainRESTHandler)
)
