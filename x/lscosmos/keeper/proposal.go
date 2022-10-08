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

// HandleRegisterHostChainProposal performs the writes host chain params.
func HandleRegisterHostChainProposal(ctx sdk.Context, k Keeper, content types.RegisterHostChainProposal) error {
	oldData := k.GetHostChainParams(ctx)
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

	hostAccounts := k.GetHostAccounts(ctx)
	if err := hostAccounts.Validate(); err != nil {
		return err
	}
	// This checks for channel being active
	err = k.icaControllerKeeper.RegisterInterchainAccount(ctx, content.ConnectionID, hostAccounts.DelegatorAccountOwnerID)
	if err != nil {
		return sdkerrors.Wrap(err, "Could not register ica delegation Address")
	}

	paramsProposal := types.NewHostChainParams(content.ChainID, content.ConnectionID, content.TransferChannel,
		content.TransferPort, content.BaseDenom, content.MintDenom, content.PstakeParams.PstakeFeeAddress,
		content.MinDeposit, content.PstakeParams.PstakeDepositFee, content.PstakeParams.PstakeRestakeFee,
		content.PstakeParams.PstakeUnstakeFee, content.PstakeParams.PstakeRedemptionFee)

	k.SetHostChainParams(ctx, paramsProposal)

	if !content.AllowListedValidators.Valid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Allow listed validators is invalid")
	}
	k.SetAllowListedValidators(ctx, content.AllowListedValidators)
	return nil
}

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
