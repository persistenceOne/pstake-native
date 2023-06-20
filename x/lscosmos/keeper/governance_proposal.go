package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
)

// HandleMinDepositAndFeeChangeProposal changes host chain params for desired min-deposit and protocol fee
func HandleMinDepositAndFeeChangeProposal(ctx sdk.Context, k Keeper, content types.MinDepositAndFeeChangeProposal) error { //nolint:staticcheck
	return types.ErrDeprecated
}

// HandlePstakeFeeAddressChangeProposal changes fee collector address
func HandlePstakeFeeAddressChangeProposal(ctx sdk.Context, k Keeper, content types.PstakeFeeAddressChangeProposal) error { //nolint:staticcheck
	return types.ErrDeprecated
}

// HandleAllowListedValidatorSetChangeProposal changes the allowList validator set
func HandleAllowListedValidatorSetChangeProposal(ctx sdk.Context, k Keeper, content types.AllowListedValidatorSetChangeProposal) error { //nolint:staticcheck
	return types.ErrDeprecated
}
