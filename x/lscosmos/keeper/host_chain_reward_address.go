package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"

	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
)

// SetHostChainRewardAddress sets host chain reward address
func (k Keeper) SetHostChainRewardAddress(ctx sdk.Context, hostChainRewardAddress types.HostChainRewardAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.HostChainRewardAddressKey, k.cdc.MustMarshal(&hostChainRewardAddress))
}

// GetHostChainRewardAddress gets host chain reward address
func (k Keeper) GetHostChainRewardAddress(ctx sdk.Context) types.HostChainRewardAddress {
	store := ctx.KVStore(k.storeKey)
	var hostChainRewardAddress types.HostChainRewardAddress
	k.cdc.MustUnmarshal(store.Get(types.HostChainRewardAddressKey), &hostChainRewardAddress)
	return hostChainRewardAddress
}

// SetHostChainRewardAddressIfEmpty  sets host chain reward address
func (k Keeper) SetHostChainRewardAddressIfEmpty(ctx sdk.Context, hostChainRewardAddress types.HostChainRewardAddress) error {
	addr := k.GetHostChainRewardAddress(ctx)
	if addr.Address == "" {
		k.SetHostChainRewardAddress(ctx, hostChainRewardAddress)
		return nil
	}
	return icatypes.ErrInterchainAccountAlreadySet
}
