package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"

	lscosmostypes "github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// GenerateAndExecuteICATx does ica transactions with messages,
// optimistic bool does not check for channel to be open. only use to do icatxns when channel is getting created.
func (k Keeper) GenerateAndExecuteICATx(ctx sdk.Context, connectionID string, portID string, msgs []sdk.Msg) error {

	channelID, found := k.icaControllerKeeper.GetOpenActiveChannel(ctx, connectionID, portID)
	if !found {
		k.Logger(ctx).Error(fmt.Sprintf("failed to retrieve active channel for port %s", portID))
		return channeltypes.ErrInvalidChannelState
	}

	chanCap, found := k.lscosmosScopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(portID, channelID))
	if !found {
		k.Logger(ctx).Error(fmt.Sprintf("module does not own channel capability, module: %s, channelID: %s, portId: %s", lscosmostypes.ModuleName, channelID, portID))
		return channeltypes.ErrChannelCapabilityNotFound
	}

	msgData, err := icatypes.SerializeCosmosTx(k.cdc, msgs)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("could not serialize cosmostx err %v", err))
		return err
	}

	icaPacketData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: msgData,
	}
	timeoutTimestamp := ctx.BlockTime().Add(lscosmostypes.ICATimeoutTimestamp).UnixNano()
	seq, err := k.icaControllerKeeper.SendTx(ctx, chanCap, connectionID, portID, icaPacketData, uint64(timeoutTimestamp))
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("send ica txn of msgs: %s failed with err: %v", msgs, err))
		return sdkerrors.Wrapf(lscosmostypes.ErrICATxFailure, "Failed to send ica msgs with err: %v", err)
	}
	k.Logger(ctx).Info(fmt.Sprintf("sent ICA transactions with seq: %v,  channelID: %s, portId: %s, msgs: %s", seq, channelID, portID, msgs))
	return nil
}
