package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/app"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

func init() {
	app.SetAddressPrefixes()
}

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
		types.AllowListedValidators{AllowListedValidators: []types.AllowListedValidator{{ValidatorAddress: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", TargetWeight: sdk.OneDec()}}},
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
	require.Equal(t, "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", pcp.AllowListedValidators.AllowListedValidators[0].ValidatorAddress)
	require.Equal(t, sdk.OneDec(), pcp.AllowListedValidators.AllowListedValidators[0].TargetWeight)
	require.Equal(t, sdk.ZeroDec(), pcp.PstakeDepositFee)
	require.Equal(t, sdk.ZeroDec(), pcp.PstakeRestakeFee)
	require.Equal(t, sdk.ZeroDec(), pcp.PstakeUnstakeFee)
	require.Equal(t, types.RouterKey, pcp.ProposalRoute())
	require.Equal(t, types.ProposalTypeRegisterHostChain, pcp.ProposalType())
	require.Nil(t, pcp.ValidateBasic())
}

func TestNewMinDepositAndFeeChangeProposal(t *testing.T) {
	pcp := types.NewMinDepositAndFeeChangeProposal(
		"title",
		"description",
		sdk.OneInt().MulRaw(5),
		sdk.ZeroDec(),
		sdk.ZeroDec(),
		sdk.ZeroDec(),
	)

	require.Equal(t, "title", pcp.GetTitle())
	require.Equal(t, sdk.NewInt(5), pcp.MinDeposit)
	require.Equal(t, sdk.ZeroDec(), pcp.PstakeDepositFee)
	require.Equal(t, sdk.ZeroDec(), pcp.PstakeRestakeFee)
	require.Equal(t, sdk.ZeroDec(), pcp.PstakeUnstakeFee)
	require.Equal(t, types.RouterKey, pcp.ProposalRoute())
	require.Equal(t, types.ProposalTypeMinDepositAndFeeChange, pcp.ProposalType())
	require.Nil(t, pcp.ValidateBasic())

}

func TestNewPstakeFeeAddressChangeProposal(t *testing.T) {
	pcp := types.NewPstakeFeeAddressChangeProposal(
		"title",
		"description",
		"persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9",
	)

	require.Equal(t, "title", pcp.GetTitle())
	require.Equal(t, "persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9", pcp.PstakeFeeAddress)
	require.Equal(t, types.RouterKey, pcp.ProposalRoute())
	require.Equal(t, types.ProposalPstakeFeeAddressChange, pcp.ProposalType())
	require.Nil(t, pcp.ValidateBasic())
}

func TestNewAllowListedValidatorSetChangeProposal(t *testing.T) {
	pcp := types.NewAllowListedValidatorSetChangeProposal(
		"title",
		"description",
		types.AllowListedValidators{AllowListedValidators: []types.AllowListedValidator{{ValidatorAddress: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", TargetWeight: sdk.OneDec()}}},
	)

	require.Equal(t, "title", pcp.GetTitle())
	require.Equal(t, "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", pcp.AllowListedValidators.AllowListedValidators[0].ValidatorAddress)
	require.Equal(t, sdk.OneDec(), pcp.AllowListedValidators.AllowListedValidators[0].TargetWeight)
	require.Equal(t, types.RouterKey, pcp.ProposalRoute())
	require.Equal(t, types.ProposalAllowListedValidatorSetChange, pcp.ProposalType())
	require.Nil(t, pcp.ValidateBasic())
}
