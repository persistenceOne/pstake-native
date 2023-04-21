package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	ibckeeper "github.com/cosmos/ibc-go/v6/modules/core/keeper"
	ibctmtypes "github.com/cosmos/ibc-go/v6/modules/light-clients/07-tendermint/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey

	accountKeeper       types.AccountKeeper
	bankKeeper          types.BankKeeper
	icaControllerKeeper types.ICAControllerKeeper
	scopedKeeper        types.ScopedKeeper
	ibcKeeper           *ibckeeper.Keeper

	paramSpace paramtypes.Subspace

	msgRouter *baseapp.MsgServiceRouter

	authority string
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,

	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	icaControllerKeeper types.ICAControllerKeeper,
	scopedKeeper types.ScopedKeeper,
	ibcKeeper *ibckeeper.Keeper,

	paramSpace paramtypes.Subspace,

	msgRouter *baseapp.MsgServiceRouter,

	authority string,
) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}
	return Keeper{
		cdc:                 cdc,
		accountKeeper:       accountKeeper,
		bankKeeper:          bankKeeper,
		icaControllerKeeper: icaControllerKeeper,
		scopedKeeper:        scopedKeeper,
		ibcKeeper:           ibcKeeper,
		storeKey:            storeKey,
		paramSpace:          paramSpace,
		msgRouter:           msgRouter,
		authority:           authority,
	}
}

// GetParams gets the total set of liquidstakeibc parameters.
func (k *Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of liquidstakeibc parameters.
func (k *Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// SendProtocolFee to the community pool
func (k *Keeper) SendProtocolFee(ctx sdk.Context, protocolFee sdk.Coins, moduleAccount, feeAddress string) error {
	addr, err := sdk.AccAddressFromBech32(feeAddress)
	if err != nil {
		return err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, moduleAccount, addr, protocolFee)
	if err != nil {
		return err
	}
	return nil
}

// GetClientState retrieves the client state given a connection id
func (k *Keeper) GetClientState(ctx sdk.Context, connectionID string) (*ibctmtypes.ClientState, error) {
	conn, found := k.ibcKeeper.ConnectionKeeper.GetConnection(ctx, connectionID)
	if !found {
		return nil, fmt.Errorf("invalid connection id, \"%s\" not found", connectionID)
	}

	clientState, found := k.ibcKeeper.ClientKeeper.GetClientState(ctx, conn.ClientId)
	if !found {
		return nil, fmt.Errorf("client id \"%s\" not found for connection \"%s\"", conn.ClientId, connectionID)
	}

	client, ok := clientState.(*ibctmtypes.ClientState)
	if !ok {
		return nil, fmt.Errorf("invalid client state for connection \"%s\"", connectionID)
	}

	return client, nil
}

// GetChainID gets the id of the host chain given a connection id
func (k *Keeper) GetChainID(ctx sdk.Context, connectionID string) (string, error) {
	clientState, err := k.GetClientState(ctx, connectionID)
	if err != nil {
		return "", fmt.Errorf("client state not found for connection \"%s\": \"%s\"", connectionID, err.Error())
	}

	return clientState.ChainId, nil
}
