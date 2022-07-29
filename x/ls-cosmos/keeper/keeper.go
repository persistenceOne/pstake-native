package keeper

import (
	"fmt"
	accountKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	ibcTransferKeeper "github.com/cosmos/ibc-go/v3/modules/apps/transfer/keeper"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/persistenceOne/pstake-native/x/ls-cosmos/types"
)

type (
	Keeper struct {
		cdc                codec.BinaryCodec
		storeKey           sdk.StoreKey
		memKey             sdk.StoreKey
		paramstore         paramtypes.Subspace
		bankKeeper         bankKeeper.BaseKeeper
		distributionKeeper distrkeeper.Keeper
		accountKeeper      accountKeeper.AccountKeeper
		ibcTransKeeper     ibcTransferKeeper.Keeper
		ibcKeeepr          ibckeeper.Keeper
		scopedKeeper       capabilitykeeper.ScopedKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,
	bankKeeper bankKeeper.BaseKeeper,
	disributionKeeper distrkeeper.Keeper,
	accKeeper accountKeeper.AccountKeeper,
	ibckeeper ibckeeper.Keeper,
	ibcTransferKeeper ibcTransferKeeper.Keeper,
	scopedKeeper capabilitykeeper.ScopedKeeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		bankKeeper:         bankKeeper,
		distributionKeeper: disributionKeeper,
		accountKeeper:      accKeeper,
		ibcKeeepr:          ibckeeper,
		ibcTransKeeper:     ibcTransferKeeper,
		scopedKeeper:       scopedKeeper,
		cdc:                cdc,
		storeKey:           storeKey,
		memKey:             memKey,
		paramstore:         ps,
	}
}

func (k Keeper) IbcTransferKeeper() ibcTransferKeeper.Keeper { return k.ibcTransKeeper }

func (k Keeper) BankKeeper() bankKeeper.BaseKeeper { return k.bankKeeper }

func (k Keeper) DistributionKeeper() distrkeeper.Keeper { return k.distributionKeeper }

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// ChanCloseInit defines a wrapper function for the channel Keeper's function
func (k Keeper) ChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	capName := host.ChannelCapabilityPath(portID, channelID)
	chanCap, ok := k.scopedKeeper.GetCapability(ctx, capName)
	if !ok {
		return sdkerrors.Wrapf(channeltypes.ErrChannelCapabilityNotFound, "could not retrieve channel capability at: %s", capName)
	}
	return k.ibcKeeepr.ChannelKeeper.ChanCloseInit(ctx, portID, channelID, chanCap)
}

// IsBound checks if the module is already bound to the desired port
func (k Keeper) IsBound(ctx sdk.Context, portID string) bool {
	_, ok := k.scopedKeeper.GetCapability(ctx, host.PortPath(portID))
	return ok
}

// BindPort defines a wrapper function for the ort Keeper's function in
// order to expose it to module's InitGenesis function
func (k Keeper) BindPort(ctx sdk.Context, portID string) error {
	capability := k.ibcKeeepr.PortKeeper.BindPort(ctx, portID)
	return k.ClaimCapability(ctx, capability, host.PortPath(portID))
}

// GetPort returns the portID for the module. Used in ExportGenesis
func (k Keeper) GetPort(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(types.PortKey))
}

// SetPort sets the portID for the module. Used in InitGenesis
func (k Keeper) SetPort(ctx sdk.Context, portID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PortKey, []byte(portID))
}

// AuthenticateCapability wraps the scopedKeeper's AuthenticateCapability function
func (k Keeper) AuthenticateCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) bool {
	return k.scopedKeeper.AuthenticateCapability(ctx, cap, name)
}

// ClaimCapability allows the module that can claim a capability that IBC module passes to it
func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}

//MintTokens in the given account
func (k Keeper) MintTokens(ctx sdk.Context, mintCoin sdk.Coin, mintAddress sdk.AccAddress) error {
	if mintCoin.Amount.GT(sdk.NewInt(5)) {
		err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(mintCoin))
		if err != nil {
			return err
		}

		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, mintAddress, sdk.NewCoins(mintCoin))
		if err != nil {
			return err
		}
	}
	return nil
}

//SendTokensToDepositAddress
func (k Keeper) SendTokensToDepositAddress(ctx sdk.Context, depositCoin sdk.Coins, depositAddress sdk.AccAddress, senderAddress sdk.AccAddress) error {
	err := k.bankKeeper.SendCoins(ctx, senderAddress, depositAddress, depositCoin)
	if err != nil {
		return err
	}
	return nil
}

//SendResidueToCommunityPool sends the residue stk token to community pool
func (k Keeper) SendResidueToCommunityPool(ctx sdk.Context, residue []sdk.DecCoin) {
	feePool := k.distributionKeeper.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(residue...)
	k.distributionKeeper.SetFeePool(ctx, feePool)
}
