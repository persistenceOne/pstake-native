package ls_cosmos

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/persistenceOne/pstake-native/x/ls-cosmos/keeper"
	"github.com/persistenceOne/pstake-native/x/ls-cosmos/types"
)

// NewLSCosmosProposalHandler creates a new governance Handler for ls-cosmos module
func NewLSCosmosProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.RegisterCosmosChainProposal:
			return HandleRegisterCosmosChainProposal(ctx, k, *c)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized proposal content type: %T", c)
		}
	}
}

// HandleRegisterCosmosChainProposal performs the writes cosmos ICB params.
func HandleRegisterCosmosChainProposal(ctx sdk.Context, k keeper.Keeper, content types.RegisterCosmosChainProposal) error {
	k.SetCosmosIBCParams(ctx, content)
	return nil
}
