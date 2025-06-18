package v3

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v3/x/liquidstakeibc/types"
)

// MigrateStore performs in-place store migrations from v2.3.0 to v2.
// The migration includes:
//
// - Migrate host chains to include the MaxEntries and AcceptableDelta attribute.
func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	for _, hc := range getAllHostChains(ctx, storeKey, cdc) {
		hc.Params.MaxEntries = 7
		hc.Params.RedelegationAcceptableDelta = sdk.NewInt(1000000000)

		setHostChain(ctx, storeKey, cdc, hc)
	}

	return nil
}

func getAllHostChains(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) []*types.HostChain {
	store := prefix.NewStore(ctx.KVStore(storeKey), types.HostChainKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	hostChains := make([]*types.HostChain, 0)
	for ; iterator.Valid(); iterator.Next() {
		hc := types.HostChain{}
		cdc.MustUnmarshal(iterator.Value(), &hc)
		hostChains = append(hostChains, &hc)
	}

	return hostChains
}

func setHostChain(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec, hc *types.HostChain) {
	store := prefix.NewStore(ctx.KVStore(storeKey), types.HostChainKey)
	bytes := cdc.MustMarshal(hc)
	store.Set([]byte(hc.ChainId), bytes)
}
