package types_test

import (
	"fmt"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/pstake-native/v2/app"
	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
)

func init() {
	app.SetAddressPrefixes()
}

func TestNewMinDepositAndFeeChangeProposal(t *testing.T) {
	testCases := []struct {
		testName, expectedString string
		proposal                 types.MinDepositAndFeeChangeProposal
		expectedError            error
	}{
		{
			testName: "correct proposal content",
			proposal: *types.NewMinDepositAndFeeChangeProposal(
				"title",
				"description",
				sdk.OneInt().MulRaw(5),
				sdk.ZeroDec(),
				sdk.ZeroDec(),
				sdk.ZeroDec(),
				sdk.ZeroDec(),
			),
			expectedError:  nil,
			expectedString: "MinDepositAndFeeChange:\nTitle:                 title\nDescription:           description\nMinDeposit:             5\nPstakeDepositFee:\t   0.000000000000000000\nPstakeRestakeFee: \t   0.000000000000000000\nPstakeUnstakeFee: \t   0.000000000000000000\nPstakeRedemptionFee:   0.000000000000000000\n\n",
		},
		{
			testName: "invalid title length",
			proposal: *types.NewMinDepositAndFeeChangeProposal(
				"",
				"description",
				sdk.OneInt().MulRaw(5),
				sdk.ZeroDec(),
				sdk.ZeroDec(),
				sdk.ZeroDec(),
				sdk.ZeroDec(),
			),
			expectedError: sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal title cannot be blank"),
		},
		{
			testName: "invalid title length",
			proposal: *types.NewMinDepositAndFeeChangeProposal(
				strings.Repeat("-", govtypes.MaxTitleLength+1),
				"description",
				sdk.OneInt().MulRaw(5),
				sdk.ZeroDec(),
				sdk.ZeroDec(),
				sdk.ZeroDec(),
				sdk.ZeroDec(),
			),
			expectedError: sdkerrors.Wrapf(govtypes.ErrInvalidProposalContent, "proposal title is longer than max length of %d", govtypes.MaxTitleLength),
		},
		{
			testName: "invalid description length",
			proposal: *types.NewMinDepositAndFeeChangeProposal(
				"title",
				"",
				sdk.OneInt().MulRaw(5),
				sdk.ZeroDec(),
				sdk.ZeroDec(),
				sdk.ZeroDec(),
				sdk.ZeroDec(),
			),
			expectedError: sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal description cannot be blank"),
		},
		{
			testName: "invalid description length",
			proposal: *types.NewMinDepositAndFeeChangeProposal(
				"title",
				strings.Repeat("-", govtypes.MaxDescriptionLength+1),
				sdk.OneInt().MulRaw(5),
				sdk.ZeroDec(),
				sdk.ZeroDec(),
				sdk.ZeroDec(),
				sdk.ZeroDec(),
			),
			expectedError: sdkerrors.Wrapf(govtypes.ErrInvalidProposalContent, "proposal description is longer than max length of %d", govtypes.MaxDescriptionLength),
		},
		{
			testName: "incorrect pstake deposit fee",
			proposal: *types.NewMinDepositAndFeeChangeProposal(
				"title",
				"description",
				sdk.OneInt().MulRaw(5),
				sdk.NewDec(10),
				sdk.ZeroDec(),
				sdk.ZeroDec(),
				sdk.ZeroDec(),
			),
			expectedError: sdkerrors.Wrapf(types.ErrInvalidFee, "pstake deposit fee must be between %s and %s", sdk.ZeroDec(), types.MaxPstakeDepositFee),
		},
		{
			testName: "incorrect pstake restake fee",
			proposal: *types.NewMinDepositAndFeeChangeProposal(
				"title",
				"description",
				sdk.OneInt().MulRaw(5),
				sdk.ZeroDec(),
				sdk.NewDec(10),
				sdk.ZeroDec(),
				sdk.ZeroDec(),
			),
			expectedError: sdkerrors.Wrapf(types.ErrInvalidFee, "pstake restake fee must be between %s and %s", sdk.ZeroDec(), types.MaxPstakeRestakeFee),
		},
		{
			testName: "incorrect pstake unstake fee",
			proposal: *types.NewMinDepositAndFeeChangeProposal(
				"title",
				"description",
				sdk.OneInt().MulRaw(5),
				sdk.ZeroDec(),
				sdk.ZeroDec(),
				sdk.NewDec(10),
				sdk.ZeroDec(),
			),
			expectedError: sdkerrors.Wrapf(types.ErrInvalidFee, "pstake unstake fee must be between %s and %s", sdk.ZeroDec(), types.MaxPstakeUnstakeFee),
		},
		{
			testName: "incorrect pstake unstake fee",
			proposal: *types.NewMinDepositAndFeeChangeProposal(
				"title",
				"description",
				sdk.OneInt().MulRaw(5),
				sdk.ZeroDec(),
				sdk.ZeroDec(),
				sdk.ZeroDec(),
				sdk.NewDec(10),
			),
			expectedError: sdkerrors.Wrapf(types.ErrInvalidFee, "pstake redemption fee must be between %s and %s", sdk.ZeroDec(), types.MaxPstakeRedemptionFee),
		},
		{
			testName: "incorrect deposit",
			proposal: *types.NewMinDepositAndFeeChangeProposal(
				"title",
				"description",
				sdk.OneInt().MulRaw(-1),
				sdk.ZeroDec(),
				sdk.ZeroDec(),
				sdk.ZeroDec(),
				sdk.ZeroDec(),
			),
			expectedError: sdkerrors.Wrapf(types.ErrInvalidDeposit, "min deposit must be positive"),
		},
	}

	for _, tc := range testCases {
		require.Equal(t, types.RouterKey, tc.proposal.ProposalRoute())
		require.Equal(t, types.ProposalTypeMinDepositAndFeeChange, tc.proposal.ProposalType())

		if tc.expectedError != nil {
			require.Equal(t, tc.expectedError.Error(), tc.proposal.ValidateBasic().Error())
		}
		if tc.expectedError == nil {
			require.Equal(t, "title", tc.proposal.GetTitle())
			require.Equal(t, "description", tc.proposal.GetDescription())
			require.Equal(t, tc.expectedString, tc.proposal.String())
		}
	}
}

func TestNewPstakeFeeAddressChangeProposal(t *testing.T) {
	testCases := []struct {
		testName, expectedString string
		proposal                 types.PstakeFeeAddressChangeProposal
		expectedError            error
	}{
		{
			testName: "correct proposal content",
			proposal: *types.NewPstakeFeeAddressChangeProposal(
				"title",
				"description",
				"persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9",
			),
			expectedError:  nil,
			expectedString: "PstakeFeeAddressChange:\nTitle:                 title\nDescription:           description\nPstakeFeeAddress: \t   persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9\n\n",
		},
		{
			testName: "invalid title length",
			proposal: *types.NewPstakeFeeAddressChangeProposal(
				"",
				"description",
				"persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9",
			),
			expectedError: sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal title cannot be blank"),
		},
		{
			testName: "invalid title length",
			proposal: *types.NewPstakeFeeAddressChangeProposal(
				strings.Repeat("-", govtypes.MaxTitleLength+1),
				"description",
				"persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9",
			),
			expectedError: sdkerrors.Wrapf(govtypes.ErrInvalidProposalContent, "proposal title is longer than max length of %d", govtypes.MaxTitleLength),
		},
		{
			testName: "invalid description length",
			proposal: *types.NewPstakeFeeAddressChangeProposal(
				"title",
				"",
				"persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9",
			),
			expectedError: sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal description cannot be blank"),
		},
		{
			testName: "invalid description length",
			proposal: *types.NewPstakeFeeAddressChangeProposal(
				"title",
				strings.Repeat("-", govtypes.MaxDescriptionLength+1),
				"persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9",
			),
			expectedError: sdkerrors.Wrapf(govtypes.ErrInvalidProposalContent, "proposal description is longer than max length of %d", govtypes.MaxDescriptionLength),
		},
		{
			testName: "invalid pstake fee address length",
			proposal: *types.NewPstakeFeeAddressChangeProposal(
				"title",
				"description",
				"cosmos1hcqg5wj9t42zawqkqucs7la85ffyv08lum327c",
			),
			expectedError: fmt.Errorf("invalid Bech32 prefix; expected persistence, got cosmos"),
		},
	}
	for _, tc := range testCases {
		require.Equal(t, types.RouterKey, tc.proposal.ProposalRoute())
		require.Equal(t, types.ProposalPstakeFeeAddressChange, tc.proposal.ProposalType())

		if tc.expectedError != nil {
			require.Equal(t, tc.expectedError.Error(), tc.proposal.ValidateBasic().Error())
		}
		if tc.expectedError == nil {
			require.Equal(t, "title", tc.proposal.GetTitle())
			require.Equal(t, "description", tc.proposal.GetDescription())
			require.Equal(t, tc.expectedString, tc.proposal.String())
		}
	}
}

func TestNewAllowListedValidatorSetChangeProposal(t *testing.T) {
	testCases := []struct {
		testName, expectedString string
		proposal                 types.AllowListedValidatorSetChangeProposal
		expectedError            error
	}{
		{
			testName: "correct proposal content",
			proposal: *types.NewAllowListedValidatorSetChangeProposal(
				"title",
				"description",
				types.AllowListedValidators{AllowListedValidators: []types.AllowListedValidator{{ValidatorAddress: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", TargetWeight: sdk.OneDec()}}},
			),
			expectedError:  nil,
			expectedString: "AllowListedValidatorSetChange:\nTitle:                 title\nDescription:           description\nAllowListedValidators: \t   {[{cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt 1.000000000000000000}]}\n\n",
		},
		{
			testName: "invalid title length",
			proposal: *types.NewAllowListedValidatorSetChangeProposal(
				"",
				"description",
				types.AllowListedValidators{AllowListedValidators: []types.AllowListedValidator{{ValidatorAddress: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", TargetWeight: sdk.OneDec()}}},
			),
			expectedError: sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal title cannot be blank"),
		},
		{
			testName: "invalid title length",
			proposal: *types.NewAllowListedValidatorSetChangeProposal(
				strings.Repeat("-", govtypes.MaxTitleLength+1),
				"description",
				types.AllowListedValidators{AllowListedValidators: []types.AllowListedValidator{{ValidatorAddress: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", TargetWeight: sdk.OneDec()}}},
			),
			expectedError: sdkerrors.Wrapf(govtypes.ErrInvalidProposalContent, "proposal title is longer than max length of %d", govtypes.MaxTitleLength),
		},
		{
			testName: "invalid description length",
			proposal: *types.NewAllowListedValidatorSetChangeProposal(
				"title",
				"",
				types.AllowListedValidators{AllowListedValidators: []types.AllowListedValidator{{ValidatorAddress: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", TargetWeight: sdk.OneDec()}}},
			),
			expectedError: sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal description cannot be blank"),
		},
		{
			testName: "invalid description length",
			proposal: *types.NewAllowListedValidatorSetChangeProposal(
				"title",
				strings.Repeat("-", govtypes.MaxDescriptionLength+1),
				types.AllowListedValidators{AllowListedValidators: []types.AllowListedValidator{{ValidatorAddress: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", TargetWeight: sdk.OneDec()}}},
			),
			expectedError: sdkerrors.Wrapf(govtypes.ErrInvalidProposalContent, "proposal description is longer than max length of %d", govtypes.MaxDescriptionLength),
		},
		{
			testName: "incorrect allow listed validators",
			proposal: *types.NewAllowListedValidatorSetChangeProposal(
				"title",
				"description",
				types.AllowListedValidators{AllowListedValidators: []types.AllowListedValidator{{ValidatorAddress: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", TargetWeight: sdk.ZeroDec()}}},
			),
			expectedError: sdkerrors.Wrapf(types.ErrInValidAllowListedValidators, "allow listed validators is not valid"),
		},
	}

	for _, tc := range testCases {
		require.Equal(t, types.RouterKey, tc.proposal.ProposalRoute())
		require.Equal(t, types.ProposalAllowListedValidatorSetChange, tc.proposal.ProposalType())

		if tc.expectedError != nil {
			require.Equal(t, tc.expectedError.Error(), tc.proposal.ValidateBasic().Error())
		}
		if tc.expectedError == nil {
			require.Equal(t, "title", tc.proposal.GetTitle())
			require.Equal(t, "description", tc.proposal.GetDescription())
			require.Equal(t, tc.expectedString, tc.proposal.String())
		}
	}
}
