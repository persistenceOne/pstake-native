package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"

	lscosmostypes "github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

func generateAndExecuteICATx(ctx sdk.Context, k Keeper, connectionID string, portID string, msgs []sdk.Msg) error {
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
		return err
	}
	k.Logger(ctx).Info(fmt.Sprintf("sent ICA transactions with seq: %v,  channelID: %s, portId: %s\"", seq, channelID, portID))
	return nil
}
