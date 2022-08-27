package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	ibcporttypes "github.com/cosmos/ibc-go/v3/modules/core/05-port/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// HandleRegisterCosmosChainProposal performs the writes cosmos IBC params.
func HandleRegisterCosmosChainProposal(ctx sdk.Context, k Keeper, content types.RegisterCosmosChainProposal) error {
	oldData := k.GetCosmosParams(ctx)
	if !oldData.IsEmpty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Module was already registered")
	}
	if !content.ModuleEnabled {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Module should also be enabled while passing register proposal")
	}
	if content.TransferPort != ibctransfertypes.PortID {
		return sdkerrors.Wrap(ibcporttypes.ErrInvalidPort, "Only acceptable TransferPort is \"transfer\"")
	}

	// checks for valid and active channel
	channel, found := k.channelKeeper.GetChannel(ctx, content.TransferPort, content.TransferChannel)
	if !found {
		return sdkerrors.Wrap(ibcchanneltypes.ErrChannelNotFound, fmt.Sprintf("channel for ibc transfer: %s not found", content.TransferChannel))
	}
	if channel.State != ibcchanneltypes.OPEN {
		return sdkerrors.Wrapf(
			ibcchanneltypes.ErrInvalidChannelState,
			"channel state is not OPEN (got %s)", channel.State.String(),
		)
	}
	// TODO Understand capabilities and see if it has to be/ should be claimed in lsscopedkeeper. If it even matters.
	_, err := k.lscosmosScopedKeeper.NewCapability(ctx, host.ChannelCapabilityPath(content.TransferPort, content.TransferChannel))
	if err != nil {
		return sdkerrors.Wrapf(err, "Failed to create and claim capability for ibc transfer port and channel")
	}

	// This checks for channel being active
	err = k.icaControllerKeeper.RegisterInterchainAccount(ctx, content.ConnectionID, types.DelegationModuleAccount)
	if err != nil {
		return sdkerrors.Wrap(err, "Could not register ica delegation Address")
	}
	err = k.icaControllerKeeper.RegisterInterchainAccount(ctx, content.ConnectionID, types.RewardModuleAccount)
	if err != nil {
		return sdkerrors.Wrap(err, "Could not register ica reward Address")
	}

	paramsProposal := types.NewCosmosParams(content.ConnectionID, content.TransferChannel,
		content.TransferPort, content.BaseDenom, content.MintDenom, content.MinDeposit,
		content.PstakeDepositFee, content.PstakeRestakeFee, content.PstakeUnstakeFee)

	k.SetCosmosParams(ctx, paramsProposal)

	if !content.AllowListedValidators.Valid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Allow listed validators is invalid")
	}
	k.SetAllowListedValidators(ctx, content.AllowListedValidators)
	return nil
}
