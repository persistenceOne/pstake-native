package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"

	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
)

// Keeper of this module maintains the state of whole module
type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   store.StoreKey
	memKey     store.StoreKey
	paramstore paramtypes.Subspace

	bankKeeper           types.BankKeeper
	accountKeeper        types.AccountKeeper
	epochKeeper          types.EpochKeeper
	ics4WrapperKeeper    types.ICS4WrapperKeeper
	channelKeeper        types.ChannelKeeper
	portKeeper           types.PortKeeper
	ibcTransferKeeper    types.IBCTransferKeeper
	icaControllerKeeper  types.ICAControllerKeeper
	icqKeeper            types.ICQKeeper
	liquidStakeIBCKeeper types.LiquidStakeIBCKeeper
	lscosmosScopedKeeper types.ScopedKeeper

	msgRouter *baseapp.MsgServiceRouter
}

// NewKeeper returns a new instance of ls cosmos module keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey store.StoreKey,
	ps paramtypes.Subspace,
	bankKeeper types.BankKeeper,
	accKeeper types.AccountKeeper,
	epochKeeper types.EpochKeeper,
	ics4WrapperKeeper types.ICS4WrapperKeeper,
	channelKeeper types.ChannelKeeper,
	portKeeper types.PortKeeper,
	ibcTransferKeeper types.IBCTransferKeeper,
	icaControllerKeeper types.ICAControllerKeeper,
	icqKeeper types.ICQKeeper,
	liquidStakeIBCKeeper types.LiquidStakeIBCKeeper,
	lscosmosScopedKeeper types.ScopedKeeper,
	msgRouter *baseapp.MsgServiceRouter,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		bankKeeper:           bankKeeper,
		accountKeeper:        accKeeper,
		epochKeeper:          epochKeeper,
		ics4WrapperKeeper:    ics4WrapperKeeper,
		channelKeeper:        channelKeeper,
		portKeeper:           portKeeper,
		ibcTransferKeeper:    ibcTransferKeeper,
		icaControllerKeeper:  icaControllerKeeper,
		icqKeeper:            icqKeeper,
		liquidStakeIBCKeeper: liquidStakeIBCKeeper,
		lscosmosScopedKeeper: lscosmosScopedKeeper,
		cdc:                  cdc,
		storeKey:             storeKey,
		memKey:               memKey,
		paramstore:           ps,
		msgRouter:            msgRouter,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// ChanCloseInit defines a wrapper function for the channel Keeper's function
func (k Keeper) ChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	capName := host.ChannelCapabilityPath(portID, channelID)
	chanCap, ok := k.lscosmosScopedKeeper.GetCapability(ctx, capName)
	if !ok {
		return errorsmod.Wrapf(channeltypes.ErrChannelCapabilityNotFound, "could not retrieve channel capability at: %s", capName)
	}
	return k.channelKeeper.ChanCloseInit(ctx, portID, channelID, chanCap)
}

// IsBound checks if the module is already bound to the desired port
func (k Keeper) IsBound(ctx sdk.Context, portID string) bool {
	_, ok := k.lscosmosScopedKeeper.GetCapability(ctx, host.PortPath(portID))
	return ok
}

// BindPort defines a wrapper function for the ort Keeper's function in
// order to expose it to module's InitGenesis function
func (k Keeper) BindPort(ctx sdk.Context, portID string) error {
	capability := k.portKeeper.BindPort(ctx, portID)
	return k.ClaimCapability(ctx, capability, host.PortPath(portID))
}

// AuthenticateCapability wraps the lscosmosScopedKeeper's AuthenticateCapability function
func (k Keeper) AuthenticateCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) bool {
	return k.lscosmosScopedKeeper.AuthenticateCapability(ctx, cap, name)
}

// ClaimCapability allows the module that can claim a capability that IBC module passes to it
func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.lscosmosScopedKeeper.ClaimCapability(ctx, cap, name)
}

// NewCapability allows the module that can initiate and claim a capability that IBC module passes to it
func (k Keeper) NewCapability(ctx sdk.Context, name string) error {
	_, err := k.lscosmosScopedKeeper.NewCapability(ctx, name)
	return err
}

// GetDepositModuleAccount returns deposit module account interface
func (k Keeper) GetDepositModuleAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.DepositModuleAccount)
}

// GetDelegationModuleAccount returns the delegation module account interface
func (k Keeper) GetDelegationModuleAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.DelegationModuleAccount)
}

// GetRewardModuleAccount returns the reward module account interface
func (k Keeper) GetRewardModuleAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.RewardModuleAccount)
}

// GetUndelegationModuleAccount returns the undelegation module account interface
func (k Keeper) GetUndelegationModuleAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.UndelegationModuleAccount)
}

// GetRewardBoosterModuleAccount returns the rewards booster module account interface
func (k Keeper) GetRewardBoosterModuleAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.RewardBoosterModuleAccount)
}

// MintTokens in the given account
func (k Keeper) MintTokens(ctx sdk.Context, mintCoin sdk.Coin, delegatorAddress sdk.AccAddress) error {

	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(mintCoin))
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delegatorAddress, sdk.NewCoins(mintCoin))
	if err != nil {
		return err
	}

	return nil
}

// SendTokensToDepositModule sends the tokens to DepositModuleAccount
func (k Keeper) SendTokensToDepositModule(ctx sdk.Context, depositCoin sdk.Coins, senderAddress sdk.AccAddress) error {
	return k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddress, types.DepositModuleAccount, depositCoin)
}

// SendProtocolFee to the community pool
func (k Keeper) SendProtocolFee(ctx sdk.Context, protocolFee sdk.Coins, moduleAccount, pstakeFeeAddressString string) error {
	addr, err := sdk.AccAddressFromBech32(pstakeFeeAddressString)
	if err != nil {
		return err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, moduleAccount, addr, protocolFee)
	if err != nil {
		return err
	}
	return nil
}
