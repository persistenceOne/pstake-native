package lscosmos

import (
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/persistenceOne/pstake-native/x/lscosmos/keeper"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// NewLSCosmosProposalHandler creates a new governance Handler for lscosmos module
func NewLSCosmosProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.RegisterHostChainProposal:
			return keeper.HandleRegisterHostChainProposal(ctx, k, *c)
		case *types.MinDepositAndFeeChangeProposal:
			return keeper.HandleMinDepositAndFeeChangeProposal(ctx, k, *c)
		case *types.PstakeFeeAddressChangeProposal:
			return keeper.HandlePstakeFeeAddressChangeProposal(ctx, k, *c)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized proposal content type: %T", c)
		}
	}
}
