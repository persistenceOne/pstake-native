package utils

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
	"github.com/stretchr/testify/require"
)

func TestNewRegisterChainJSON(t *testing.T) {
	propJSON := NewRegisterChainJSON(
		"title",
		"description",
		true,
		"cosmoshub-4",
		"connection",
		"channel-1",
		"transfer",
		"uatom",
		"ustkatom",
		"5",
		"persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9",
		"0.0",
		"0.0",
		"0.0",
		"1000stake",
		types.AllowListedValidators{
			AllowListedValidators: []types.AllowListedValidator{{
				ValidatorAddress: "Valaddr",
				TargetWeight:     sdk.OneDec(),
			}}},
	)
	require.Equal(t, "title", propJSON.Title)
	require.Equal(t, "description", propJSON.Description)
	require.Equal(t, "cosmoshub-4", propJSON.ChainID)
	require.Equal(t, "connection", propJSON.ConnectionID)
	require.Equal(t, "channel-1", propJSON.TransferChannel)
	require.Equal(t, "transfer", propJSON.TransferPort)
	require.Equal(t, "uatom", propJSON.BaseDenom)
	require.Equal(t, "ustkatom", propJSON.MintDenom)
	require.Equal(t, "5", propJSON.MinDeposit)
	require.Equal(t, "persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9", propJSON.PstakeFeeAddress)
	require.Equal(t, "0.0", propJSON.PstakeDepositFee)
	require.Equal(t, "0.0", propJSON.PstakeRestakeFee)
	require.Equal(t, "0.0", propJSON.PstakeUnstakeFee)
	require.Equal(t, "1000stake", propJSON.Deposit)
}

func TestNewMinDepositAndFeeChangeJSON(t *testing.T) {
	propJSON := NewMinDepositAndFeeChangeJSON(
		"title",
		"description",
		"5",
		"0.0",
		"0.0",
		"0.0",
		"1000stake")

	require.Equal(t, "title", propJSON.Title)
	require.Equal(t, "description", propJSON.Description)
	require.Equal(t, "5", propJSON.MinDeposit)
	require.Equal(t, "0.0", propJSON.PstakeDepositFee)
	require.Equal(t, "0.0", propJSON.PstakeRestakeFee)
	require.Equal(t, "0.0", propJSON.PstakeUnstakeFee)
	require.Equal(t, "1000stake", propJSON.Deposit)
}

func TestNewPstakeFeeAddressChangeProposalJSON(t *testing.T) {
	propJSON := NewPstakeFeeAddressChangeProposalJSON(
		"title",
		"description",
		"persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9",
		"1000stake")

	require.Equal(t, "title", propJSON.Title)
	require.Equal(t, "description", propJSON.Description)
	require.Equal(t, "persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9", propJSON.PstakeFeeAddress)
	require.Equal(t, "1000stake", propJSON.Deposit)
}

func TestNewAllowListedValidatorSetChangeProposalJSON(t *testing.T) {
	propJSON := NewAllowListedValidatorSetChangeProposalJSON(
		"title",
		"description",
		"1000stake",
		types.AllowListedValidators{
			AllowListedValidators: []types.AllowListedValidator{{
				ValidatorAddress: "Valaddr",
				TargetWeight:     sdk.OneDec(),
			}}},
	)

	require.Equal(t, "title", propJSON.Title)
	require.Equal(t, "description", propJSON.Description)
	require.Equal(t, "Valaddr", propJSON.AllowListedValidators.AllowListedValidators[0].ValidatorAddress)
	require.Equal(t, sdk.OneDec(), propJSON.AllowListedValidators.AllowListedValidators[0].TargetWeight)
	require.Equal(t, "1000stake", propJSON.Deposit)

}
