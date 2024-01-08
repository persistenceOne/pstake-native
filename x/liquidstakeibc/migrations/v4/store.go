package v4

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

// MigrateStore performs in-place store migrations from v2.x to v2.8.2.
// The migration includes:
//
// - Migrate unbondings to mark stuck pending unbondings as failed.
func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {

	RemovableUnbondings := map[string]map[int64]any{"cosmoshub-4": {312: nil}, "osmosis-1": {429: nil, 432: nil}}

	for _, removableUnbonding := range getRemovableUnbondings(ctx, storeKey, cdc, RemovableUnbondings) {
		removableUnbonding.State = types.Unbonding_UNBONDING_FAILED
		setUnbonding(ctx, storeKey, cdc, removableUnbonding)
	}

	return nil
}

func getRemovableUnbondings(
	ctx sdk.Context,
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	removableUnbondings map[string]map[int64]interface{},
) []*types.Unbonding {
	store := prefix.NewStore(ctx.KVStore(storeKey), types.UnbondingKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	unbondings := make([]*types.Unbonding, 0)
	for ; iterator.Valid(); iterator.Next() {
		unbonding := types.Unbonding{}
		cdc.MustUnmarshal(iterator.Value(), &unbonding)

		_, chain := removableUnbondings[unbonding.ChainId]
		if chain {
			_, epoch := removableUnbondings[unbonding.ChainId][unbonding.EpochNumber]
			if epoch {
				unbondings = append(unbondings, &unbonding)
			}
		}
	}

	return unbondings
}

func setUnbonding(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec, ub *types.Unbonding) {
	store := prefix.NewStore(ctx.KVStore(storeKey), types.UnbondingKey)
	bytes := cdc.MustMarshal(ub)
	store.Set(append([]byte(ub.ChainId), []byte(strconv.FormatInt(ub.EpochNumber, 10))...), bytes)
}
