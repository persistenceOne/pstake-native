package client

import (
	"github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/persistenceOne/pstake-native/x/cosmos/client/cli"
	"github.com/persistenceOne/pstake-native/x/cosmos/client/rest"
)

var (
	EnableModuleProposalHandler                 = client.NewProposalHandler(cli.NewEnableModuleCmd, rest.EnableModuleProposalRESTHandler)
	ChangeMultisigProposalHandler               = client.NewProposalHandler(cli.NewChangeMultisigCmd, rest.ChangeMultisigProposalRESTHandler)
	ChangeCosmosValidatorWeightsProposalHandler = client.NewProposalHandler(cli.NewChangeCosmosValidatorWeightsCmd, rest.ChangeCosmosValidatorWeightsProposalRESTHandler)
	ChangeOracleValidatorWeightsProposalHandler = client.NewProposalHandler(cli.NewChangeOracleValidatorWeightsCmd, rest.ChangeOracleValidatorWeightsProposalRESTHandler)
)
