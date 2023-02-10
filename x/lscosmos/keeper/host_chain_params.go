package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// SetHostChainParams sets the host chain params in store
func (k Keeper) SetHostChainParams(ctx sdk.Context, hostChainParams types.HostChainParams) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.HostChainParamsKey, k.cdc.MustMarshal(&hostChainParams))
}

// GetHostChainParams gets the host chain params in store
func (k Keeper) GetHostChainParams(ctx sdk.Context) types.HostChainParams {
	store := ctx.KVStore(k.storeKey)

	var hostChainParams types.HostChainParams
	k.cdc.MustUnmarshal(store.Get(types.HostChainParamsKey), &hostChainParams)

	return hostChainParams
}

// GetIBCDenom returns IBC denom in form of string
func (k Keeper) GetIBCDenom(ctx sdk.Context) string {
	hostChainParams := k.GetHostChainParams(ctx)
	ibcDenom := ibctransfertypes.ParseDenomTrace(
		ibctransfertypes.GetPrefixedDenom(
			hostChainParams.TransferPort, hostChainParams.TransferChannel, hostChainParams.BaseDenom,
		),
	).IBCDenom()

	return ibcDenom
}
