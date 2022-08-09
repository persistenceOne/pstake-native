package lscosmos

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/persistenceOne/pstake-native/x/lscosmos/keeper"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// NewLSCosmosProposalHandler creates a new governance Handler for lscosmos module
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
	minDeposit, ok := sdk.NewIntFromString(content.MinDeposit)
	if !ok {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "minimum deposit must be a positive integer")
	}
	pStakeDepositFee, err := sdk.NewDecFromStr(content.PStakeDepositFee)
	if err != nil {
		return err
	}
	paramsProposal := types.NewCosmosIBCParams(content.IBCConnection, content.TokenTransferChannel,
		content.TokenTransferPort, content.BaseDenom, content.MintDenom, minDeposit, pStakeDepositFee)

	k.SetCosmosIBCParams(ctx, paramsProposal)
	err = k.RegisterICAAccounts(ctx, paramsProposal)
	if err != nil {
		return err
	}
	return nil
}
