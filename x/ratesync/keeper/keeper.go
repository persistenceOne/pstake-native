package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	ibctmtypes "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	ibclocalhosttypes "github.com/cosmos/ibc-go/v7/modules/light-clients/09-localhost"

	"github.com/persistenceOne/pstake-native/v2/x/ratesync/types"
)

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey storetypes.StoreKey

		epochsKeeper        types.EpochsKeeper
		icaControllerKeeper types.ICAControllerKeeper
		ibcKeeper           *ibckeeper.Keeper
		liquidStakeKeeper   types.LiquidStakeKeeper

		msgRouter *baseapp.MsgServiceRouter

		authority string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	epochsKeeper types.EpochsKeeper,
	icaControllerKeeper types.ICAControllerKeeper,
	ibcKeeper *ibckeeper.Keeper,
	msgRouter *baseapp.MsgServiceRouter,
	authority string,
) *Keeper {
	return &Keeper{
		cdc:                 cdc,
		storeKey:            storeKey,
		epochsKeeper:        epochsKeeper,
		icaControllerKeeper: icaControllerKeeper,
		ibcKeeper:           ibcKeeper,
		msgRouter:           msgRouter,
		authority:           authority,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetClientState retrieves the client state given a connection id
func (k *Keeper) GetClientState(ctx sdk.Context, connectionID string) (exported.ClientState, error) {
	conn, found := k.ibcKeeper.ConnectionKeeper.GetConnection(ctx, connectionID)
	if !found {
		return nil, fmt.Errorf("invalid connection id, \"%s\" not found", connectionID)
	}

	clientState, found := k.ibcKeeper.ClientKeeper.GetClientState(ctx, conn.ClientId)
	if !found {
		return nil, fmt.Errorf("client id \"%s\" not found for connection \"%s\"", conn.ClientId, connectionID)
	}

	return clientState, nil
}

// GetChainID gets the id of the host chain given a connection id
func (k *Keeper) GetChainID(ctx sdk.Context, connectionID string) (string, error) {
	clientState, err := k.GetClientState(ctx, connectionID)
	if err != nil {
		return "", fmt.Errorf("client state not found for connection \"%s\": \"%s\"", connectionID, err.Error())
	}

	switch clientType := clientState.(type) {
	case *ibctmtypes.ClientState:
		return clientType.ChainId, nil
	case *ibclocalhosttypes.ClientState:
		return ctx.ChainID(), nil
	default:
		return "", fmt.Errorf("unexpected type of client, cannot determine chain-id: clientType: %s, connectionid: %s", clientState.ClientType(), connectionID)
	}
}
