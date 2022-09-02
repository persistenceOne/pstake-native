package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

func TestParameterChangeProposal(t *testing.T) {
	pcp := types.NewRegisterHostChainProposal(
		"title",
		"description",
		true,
		"cosmoshub-4",
		"connection-0",
		"channel-1",
		"transfer",
		"uatom",
		"ustkatom",
		"persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9",
		sdk.OneInt().MulRaw(5),
		types.AllowListedValidators{AllowListedValidators: []types.AllowListedValidator{{ValidatorAddress: "addr", TargetWeight: sdk.OneDec()}}},
		sdk.ZeroDec(),
		sdk.ZeroDec(),
		sdk.ZeroDec(),
	)

	require.Equal(t, "title", pcp.GetTitle())
	require.Equal(t, "cosmoshub-4", pcp.ChainID)
	require.Equal(t, "connection-0", pcp.ConnectionID)
	require.Equal(t, true, pcp.ModuleEnabled)
	require.Equal(t, "channel-1", pcp.TransferChannel)
	require.Equal(t, "transfer", pcp.TransferPort)
	require.Equal(t, "uatom", pcp.BaseDenom)
	require.Equal(t, "ustkatom", pcp.MintDenom)
	require.Equal(t, "persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9", pcp.PstakeFeeAddress)
	require.Equal(t, sdk.NewInt(5), pcp.MinDeposit)
	require.Equal(t, "addr", pcp.AllowListedValidators.AllowListedValidators[0].ValidatorAddress)
	require.Equal(t, sdk.OneDec(), pcp.AllowListedValidators.AllowListedValidators[0].TargetWeight)
	require.Equal(t, sdk.ZeroDec(), pcp.PstakeDepositFee)
	require.Equal(t, sdk.ZeroDec(), pcp.PstakeRestakeFee)
	require.Equal(t, sdk.ZeroDec(), pcp.PstakeUnstakeFee)
	require.Equal(t, types.RouterKey, pcp.ProposalRoute())
	require.Equal(t, types.ProposalTypeRegisterHostChain, pcp.ProposalType())
	require.Nil(t, pcp.ValidateBasic())
}
