package keeper

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	icatypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func (k *Keeper) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	channelCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) (string, error) {
	return version, nil
}

func (k *Keeper) OnChanOpenAck(
	ctx sdk.Context,
	portID string,
	channelID string,
	counterpartyChannelID string,
	counterpartyVersion string,
) error {
	// get the connection id from the port and channel identifiers
	connID, _, err := k.ibcKeeper.ChannelKeeper.GetChannelConnection(ctx, portID, channelID)
	if err != nil {
		return fmt.Errorf("unable to get connection id using port %s: %w", portID, err)
	}

	// get interchain account address
	address, found := k.icaControllerKeeper.GetInterchainAccountAddress(ctx, connID, portID)
	if !found {
		return fmt.Errorf("couldn't find address for %s/%s", connID, portID)
	}

	// get the port owner from the port id
	portOwner, found := strings.CutPrefix(portID, icatypes.ControllerPortPrefix)
	if !found {
		return fmt.Errorf("unable to parse port id %s", portID)
	}

	// create the ica account
	icaAccount := &types.ICAAccount{Address: address, Balance: sdk.Coins{}, Owner: portOwner}

	// get the chain id using the connection id
	chainID, err := k.GetChainID(ctx, connID)
	if err != nil {
		return fmt.Errorf("unable to get chain id for connection %s: %w", connID, err)
	}

	// get host chain
	hc, found := k.GetHostChain(ctx, chainID)
	if !found {
		return fmt.Errorf("host chain with id %s is not registered", chainID)
	}

	// get the ica account type from the ownership string
	_, icaAccountType, found := strings.Cut(portOwner, ".")
	if !found {
		return fmt.Errorf("unable to parse port owner %s", portOwner)
	}

	switch icaAccountType { // TODO: Query for balances upon creation ?
	case types.DelegateICAType:
		hc.DelegationAccount = icaAccount
	case types.RewardsICAType:
		hc.RewardsAccount = icaAccount
	}

	k.SetHostChain(ctx, &hc)
	return nil
}

func (k *Keeper) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	return nil
}

func (k *Keeper) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	return nil
}
