package liquidstakeibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v6/modules/core/exported"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/keeper"
)

var _ porttypes.IBCModule = &IBCModule{}

// IBCModule implements the ICS26 callbacks for the fee middleware given the
// fee keeper and the underlying application.
type IBCModule struct {
	keeper keeper.Keeper
}

func NewIBCModule(keeper keeper.Keeper) IBCModule {
	return IBCModule{
		keeper: keeper,
	}
}

func (ibcModule IBCModule) OnChanOpenInit(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID string, channelID string, channelCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, version string) (string, error) {
	return ibcModule.keeper.OnChanOpenInit()
}

func (ibcModule IBCModule) OnChanOpenTry(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID, channelID string, channelCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, counterpartyVersion string) (version string, err error) {
	return "", nil
}

func (ibcModule IBCModule) OnChanOpenAck(ctx sdk.Context, portID, channelID string, counterpartyChannelID string, counterpartyVersion string) error {
	return ibcModule.keeper.OnChanOpenAck()
}

func (ibcModule IBCModule) OnChanOpenConfirm(ctx sdk.Context, portID, channelID string) error {
	return nil
}

func (ibcModule IBCModule) OnChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	return nil
}

func (ibcModule IBCModule) OnChanCloseConfirm(ctx sdk.Context, portID, channelID string) error {
	return nil
}

func (ibcModule IBCModule) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) ibcexported.Acknowledgement {
	return nil
}

func (ibcModule IBCModule) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte, relayer sdk.AccAddress) error {
	return ibcModule.keeper.OnAcknowledgementPacket()
}

func (ibcModule IBCModule) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) error {
	return ibcModule.keeper.OnTimeoutPacket()
}
