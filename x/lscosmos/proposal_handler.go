package lscosmos

import (
	govv1beta1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/persistenceOne/pstake-native/x/lscosmos/keeper"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// NewLSCosmosProposalHandler creates a new governance Handler for lscosmos module
func NewLSCosmosProposalHandler(k keeper.Keeper) govv1beta1types.Handler {
	return func(ctx sdk.Context, content govv1beta1types.Content) error {
		switch c := content.(type) {
		case *types.RegisterHostChainProposal:
			return keeper.HandleRegisterHostChainProposal(ctx, k, *c)
		case *types.MinDepositAndFeeChangeProposal:
			return keeper.HandleMinDepositAndFeeChangeProposal(ctx, k, *c)
		case *types.PstakeFeeAddressChangeProposal:
			return keeper.HandlePstakeFeeAddressChangeProposal(ctx, k, *c)
		case *types.AllowListedValidatorSetChangeProposal:
			return keeper.HandleAllowListedValidatorSetChangeProposal(ctx, k, *c)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized proposal content type: %T", c)
		}
	}
}
