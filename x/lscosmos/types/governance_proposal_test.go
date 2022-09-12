package types_test

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/persistenceOne/pstake-native/app"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

func init() {
	app.SetAddressPrefixes()
}

func TestParameterChangeProposal(t *testing.T) {
	testCases := []struct {
		testName      string
		proposal      types.RegisterHostChainProposal
		expectedError error
	}{
		{
			testName: "correct proposal content",
			proposal: *types.NewRegisterHostChainProposal("title", "description", true,
				"cosmoshub-4", "connection-0", "channel-1", "transfer",
				"uatom", "ustkatom", "persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9",
				sdk.OneInt().MulRaw(5),
				types.AllowListedValidators{AllowListedValidators: []types.AllowListedValidator{{ValidatorAddress: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", TargetWeight: sdk.OneDec()}}},
				sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(),
			),
			expectedError: nil,
		},
		{
			testName: "invalid title length",
			proposal: *types.NewRegisterHostChainProposal("", "description", true,
				"cosmoshub-4", "connection-0", "channel-1", "transfer",
				"uatom", "ustkatom", "persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9",
				sdk.OneInt().MulRaw(5),
				types.AllowListedValidators{AllowListedValidators: []types.AllowListedValidator{{ValidatorAddress: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", TargetWeight: sdk.OneDec()}}},
				sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(),
			),
			expectedError: sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal title cannot be blank"),
		},
		{
			testName: "invalid title length",
			proposal: *types.NewRegisterHostChainProposal(strings.Repeat("-", 141), "description", true,
				"cosmoshub-4", "connection-0", "channel-1", "transfer",
				"uatom", "ustkatom", "persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9",
				sdk.OneInt().MulRaw(5),
				types.AllowListedValidators{AllowListedValidators: []types.AllowListedValidator{{ValidatorAddress: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", TargetWeight: sdk.OneDec()}}},
				sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(),
			),
			expectedError: sdkerrors.Wrapf(govtypes.ErrInvalidProposalContent, "proposal title is longer than max length of %d", govtypes.MaxTitleLength),
		},
		{
			testName: "invalid description length",
			proposal: *types.NewRegisterHostChainProposal("title", "", true,
				"cosmoshub-4", "connection-0", "channel-1", "transfer",
				"uatom", "ustkatom", "persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9",
				sdk.OneInt().MulRaw(5),
				types.AllowListedValidators{AllowListedValidators: []types.AllowListedValidator{{ValidatorAddress: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", TargetWeight: sdk.OneDec()}}},
				sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(),
			),
			expectedError: sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal description cannot be blank"),
		},
		{
			testName: "invalid description length",
			proposal: *types.NewRegisterHostChainProposal("title", strings.Repeat("-", govtypes.MaxDescriptionLength+1), true,
				"cosmoshub-4", "connection-0", "channel-1", "transfer",
				"uatom", "ustkatom", "persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9",
				sdk.OneInt().MulRaw(5),
				types.AllowListedValidators{AllowListedValidators: []types.AllowListedValidator{{ValidatorAddress: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", TargetWeight: sdk.OneDec()}}},
				sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(),
			),
			expectedError: sdkerrors.Wrapf(govtypes.ErrInvalidProposalContent, "proposal description is longer than max length of %d", govtypes.MaxDescriptionLength),
		},
		{
			testName: "incorrect allow listed validators",
			proposal: *types.NewRegisterHostChainProposal("title", "description", true,
				"cosmoshub-4", "connection-0", "channel-1", "transfer",
				"uatom", "ustkatom", "persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9",
				sdk.OneInt().MulRaw(5),
				types.AllowListedValidators{AllowListedValidators: []types.AllowListedValidator{{ValidatorAddress: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", TargetWeight: sdk.ZeroDec()}}},
				sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(),
			),
			expectedError: sdkerrors.Wrapf(types.ErrInValidAllowListedValidators, "allow listed validators is not valid"),
		},
		{
			testName: "incorrect pstake deposit fee",
			proposal: *types.NewRegisterHostChainProposal("title", "description", true,
				"cosmoshub-4", "connection-0", "channel-1", "transfer",
				"uatom", "ustkatom", "persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9",
				sdk.OneInt().MulRaw(5),
				types.AllowListedValidators{AllowListedValidators: []types.AllowListedValidator{{ValidatorAddress: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", TargetWeight: sdk.OneDec()}}},
				sdk.NewDec(10), sdk.ZeroDec(), sdk.ZeroDec(),
			),
			expectedError: sdkerrors.Wrapf(types.ErrInvalidFee, "pstake deposit fee must be between 0 and 1"),
		},
		{
			testName: "incorrect pstake restake fee",
			proposal: *types.NewRegisterHostChainProposal("title", "description", true,
				"cosmoshub-4", "connection-0", "channel-1", "transfer",
				"uatom", "ustkatom", "persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9",
				sdk.OneInt().MulRaw(5),
				types.AllowListedValidators{AllowListedValidators: []types.AllowListedValidator{{ValidatorAddress: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", TargetWeight: sdk.OneDec()}}},
				sdk.ZeroDec(), sdk.NewDec(10), sdk.ZeroDec(),
			),
			expectedError: sdkerrors.Wrapf(types.ErrInvalidFee, "pstake restake fee must be between 0 and 1"),
		},
		{
			testName: "incorrect pstake unstake fee",
			proposal: *types.NewRegisterHostChainProposal("title", "description", true,
				"cosmoshub-4", "connection-0", "channel-1", "transfer",
				"uatom", "ustkatom", "persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9",
				sdk.OneInt().MulRaw(5),
				types.AllowListedValidators{AllowListedValidators: []types.AllowListedValidator{{ValidatorAddress: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", TargetWeight: sdk.OneDec()}}},
				sdk.ZeroDec(), sdk.ZeroDec(), sdk.NewDec(10),
			),
			expectedError: sdkerrors.Wrapf(types.ErrInvalidFee, "pstake unstake fee must be between 0 and 1"),
		},
		{
			testName: "incorrect deposit",
			proposal: *types.NewRegisterHostChainProposal("title", "description", true,
				"cosmoshub-4", "connection-0", "channel-1", "transfer",
				"uatom", "ustkatom", "persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9",
				sdk.OneInt().MulRaw(-1),
				types.AllowListedValidators{AllowListedValidators: []types.AllowListedValidator{{ValidatorAddress: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", TargetWeight: sdk.OneDec()}}},
				sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(),
			),
			expectedError: sdkerrors.Wrapf(types.ErrInvalidDeposit, "min deposit must be positive"),
		},
	}

	for _, tc := range testCases {
		require.Equal(t, types.RouterKey, tc.proposal.ProposalRoute())
		require.Equal(t, types.ProposalTypeRegisterHostChain, tc.proposal.ProposalType())

		if tc.expectedError != nil {
			require.Equal(t, tc.expectedError.Error(), tc.proposal.ValidateBasic().Error())
		}
	}

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
