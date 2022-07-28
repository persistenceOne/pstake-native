package client

import (
	"github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/persistenceOne/pstake-native/x/ls-cosmos/client/cli"
	"github.com/persistenceOne/pstake-native/x/ls-cosmos/client/rest"
)

var (
	RegisterCosmosChainProposalHandler = client.NewProposalHandler(cli.NewRegisterCosmosChainCmd, rest.RegisterChainRESTHandler)
)
