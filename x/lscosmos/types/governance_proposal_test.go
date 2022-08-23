package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

func TestParameterChangeProposal(t *testing.T) {
	pcp := types.NewRegisterCosmosChainProposal(
		"title",
		"description",
		true,
		"connection-0",
		"channel-1",
		"transfer",
		"uatom",
		"ustkatom",
		sdk.OneInt().MulRaw(5),
		types.AllowListedValidators{AllowListedValidators: []types.AllowListedValidator{{ValidatorAddress: "addr", TargetWeight: sdk.OneDec()}}},
		sdk.ZeroDec(),
		sdk.ZeroDec(),
		sdk.ZeroDec(),
	)

	require.Equal(t, "title", pcp.GetTitle())
	require.Equal(t, "connection-0", pcp.IBCConnection)
	require.Equal(t, true, pcp.ModuleEnabled)
	require.Equal(t, "channel-1", pcp.TokenTransferChannel)
	require.Equal(t, "transfer", pcp.TokenTransferPort)
	require.Equal(t, "uatom", pcp.BaseDenom)
	require.Equal(t, "ustkatom", pcp.MintDenom)
	require.Equal(t, sdk.NewInt(5), pcp.MinDeposit)
	require.Equal(t, "addr", pcp.AllowListedValidators.AllowListedValidators[0].ValidatorAddress)
	require.Equal(t, sdk.OneDec(), pcp.AllowListedValidators.AllowListedValidators[0].TargetWeight)
	require.Equal(t, sdk.ZeroDec(), pcp.PStakeDepositFee)
	require.Equal(t, sdk.ZeroDec(), pcp.PStakeRestakeFee)
	require.Equal(t, sdk.ZeroDec(), pcp.PStakeUnstakeFee)
	require.Equal(t, types.RouterKey, pcp.ProposalRoute())
	require.Equal(t, types.ProposalTypeRegisterCosmosChain, pcp.ProposalType())
	require.Nil(t, pcp.ValidateBasic())
}
