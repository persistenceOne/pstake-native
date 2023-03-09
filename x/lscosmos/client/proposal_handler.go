package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/client/cli"
)

var (
	MinDepositAndFeeChangeProposalHandler      = govclient.NewProposalHandler(cli.NewMinDepositAndFeeChangeCmd)
	PstakeFeeAddressChangeProposalHandler      = govclient.NewProposalHandler(cli.NewPstakeFeeAddressChangeCmd)
	AllowListValidatorSetChangeProposalHandler = govclient.NewProposalHandler(cli.NewAllowListedValidatorSetChangeProposalCmd)
)
