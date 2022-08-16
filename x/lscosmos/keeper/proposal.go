package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibcporttypes "github.com/cosmos/ibc-go/v3/modules/core/05-port/types"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// HandleRegisterCosmosChainProposal performs the writes cosmos ICB params.
func HandleRegisterCosmosChainProposal(ctx sdk.Context, k Keeper, content types.RegisterCosmosChainProposal) error {
	minDeposit, ok := sdk.NewIntFromString(content.MinDeposit)
	if !ok {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "minimum deposit must be a positive integer")
	}

	pStakeDepositFee, err := sdk.NewDecFromStr(content.PStakeDepositFee)
	if err != nil {
		return err
	}

	if content.TokenTransferPort != ibctransfertypes.PortID {
		return sdkerrors.Wrap(ibcporttypes.ErrInvalidPort, "Only acceptable TokenTransferPort is \"transfer\"")
	}

	// This checks for channel being active
	err = k.icaControllerKeeper.RegisterInterchainAccount(ctx, content.IBCConnection, types.DelegationModuleAccount)
	if err != nil {
		return sdkerrors.Wrap(err, "Could not register ica delegation Address")
	}
	err = k.icaControllerKeeper.RegisterInterchainAccount(ctx, content.IBCConnection, types.RewardModuleAccount)
	if err != nil {
		return sdkerrors.Wrap(err, "Could not register ica reward Address")
	}

	paramsProposal := types.NewCosmosIBCParams(content.IBCConnection, content.TokenTransferChannel,
		content.TokenTransferPort, content.BaseDenom, content.MintDenom, minDeposit, pStakeDepositFee)

	k.SetCosmosIBCParams(ctx, paramsProposal)
	return nil
}
