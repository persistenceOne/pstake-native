package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/pstake-native/x/ls-cosmos/types"
)

func TestParameterChangeProposal(t *testing.T) {
	pcp := types.NewRegisterCosmosChainProposal(
		"title",
		"description",
		"connection",
		"channel-1",
		"transfer",
		"uatom",
		"ustkatom",
	)

	require.Equal(t, "title", pcp.GetTitle())
	require.Equal(t, "description", pcp.GetDescription())
	require.Equal(t, types.RouterKey, pcp.ProposalRoute())
	require.Equal(t, types.ProposalTypeRegisterCosmosChain, pcp.ProposalType())
	require.Nil(t, pcp.ValidateBasic())
}
