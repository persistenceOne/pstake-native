package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
)

// HandleMinDepositAndFeeChangeProposal changes host chain params for desired min-deposit and protocol fee
func HandleMinDepositAndFeeChangeProposal(ctx sdk.Context, k Keeper, content types.MinDepositAndFeeChangeProposal) error {
	if !k.GetModuleState(ctx) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Module not enabled")
	}

	hostChainParams := k.GetHostChainParams(ctx)
	if hostChainParams.IsEmpty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "host chain not registered")
	}

	// modify oldData with the new proposal content
	hostChainParams.MinDeposit = content.MinDeposit
	hostChainParams.PstakeParams.PstakeDepositFee = content.PstakeDepositFee
	hostChainParams.PstakeParams.PstakeRestakeFee = content.PstakeRestakeFee
	hostChainParams.PstakeParams.PstakeUnstakeFee = content.PstakeUnstakeFee
	hostChainParams.PstakeParams.PstakeRedemptionFee = content.PstakeRedemptionFee

	k.SetHostChainParams(ctx, hostChainParams)

	return nil
}

// HandlePstakeFeeAddressChangeProposal changes fee collector address
func HandlePstakeFeeAddressChangeProposal(ctx sdk.Context, k Keeper, content types.PstakeFeeAddressChangeProposal) error {
	//Do not check ModuleEnabled state or host chain params here because non-critical proposal and will help not hardcode address inside default genesis

	hostChainParams := k.GetHostChainParams(ctx)

	// modify oldData with the new proposal content
	hostChainParams.PstakeParams.PstakeFeeAddress = content.PstakeFeeAddress

	k.SetHostChainParams(ctx, hostChainParams)

	return nil
}

// HandleAllowListedValidatorSetChangeProposal changes the allowList validator set
func HandleAllowListedValidatorSetChangeProposal(ctx sdk.Context, k Keeper, content types.AllowListedValidatorSetChangeProposal) error {
	if !k.GetModuleState(ctx) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Module not enabled")
	}

	hostChainParams := k.GetHostChainParams(ctx)
	if hostChainParams.IsEmpty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "host chain not registered")
	}

	if !content.AllowListedValidators.Valid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Allow listed validators is invalid")
	}

	k.SetAllowListedValidators(ctx, content.AllowListedValidators)
	return nil
}
