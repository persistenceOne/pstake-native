package ratesync

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"

	"github.com/persistenceOne/pstake-native/v2/x/ratesync/keeper"
)

var _ porttypes.IBCModule = &IBCModule{}

// IBCModule implements the ICS26 callbacks for the fee middleware given the
// fee keeper and the underlying application.
type IBCModule struct {
	appStack porttypes.IBCModule
	keeper   keeper.Keeper
}

func NewIBCModule(appStack porttypes.IBCModule, keeper keeper.Keeper) IBCModule {
	return IBCModule{
		appStack: appStack,
		keeper:   keeper,
	}
}

func (m IBCModule) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	channelCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) (string, error) {
	return m.appStack.OnChanOpenInit(
		ctx,
		order,
		connectionHops,
		portID,
		channelID,
		channelCap,
		counterparty,
		version,
	)
}

func (m IBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID string,
	channelID string,
	counterpartyChannelID string,
	counterpartyVersion string,
) error {
	err := m.keeper.OnChanOpenAck(
		ctx,
		portID,
		channelID,
		counterpartyChannelID,
		counterpartyVersion,
	)
	if err != nil {
		return err
	}
	return m.appStack.OnChanOpenAck(ctx, portID, channelID, counterpartyChannelID, counterpartyVersion)

}

func (m IBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	err := m.keeper.OnAcknowledgementPacket(
		ctx,
		packet,
		acknowledgement,
		relayer,
	)
	if err != nil {
		return err
	}
	return m.appStack.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
}

func (m IBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	err := m.keeper.OnTimeoutPacket(
		ctx,
		packet,
		relayer,
	)
	if err != nil {
		return err
	}
	return m.appStack.OnTimeoutPacket(ctx, packet, relayer)
}

func (m IBCModule) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	channelCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	counterpartyVersion string,
) (version string, err error) {
	return m.appStack.OnChanOpenTry(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, counterpartyVersion)
}

func (m IBCModule) OnChanOpenConfirm(ctx sdk.Context, portID, channelID string) error {
	return m.appStack.OnChanOpenConfirm(ctx, portID, channelID)
}

func (m IBCModule) OnChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	return m.appStack.OnChanCloseInit(ctx, portID, channelID)
}

func (m IBCModule) OnChanCloseConfirm(ctx sdk.Context, portID, channelID string) error {
	return m.appStack.OnChanCloseConfirm(ctx, portID, channelID)
}

func (m IBCModule) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) ibcexported.Acknowledgement {
	return m.appStack.OnRecvPacket(ctx, packet, relayer)
}
