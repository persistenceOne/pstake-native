package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

func (k Keeper) setOracleLastUpdateHeightNative(ctx sdk.Context, orchestratorAddress sdk.AccAddress, nativeBlockHeight int64) {
	oracleLastUpdateStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOracleLastUpdateHeightNative)
	oracleLastUpdateStore.Set(orchestratorAddress.Bytes(), cosmosTypes.Int64Bytes(nativeBlockHeight))
}

func (k Keeper) getOracleLastUpdateHeightNative(ctx sdk.Context, orchestratorAddress sdk.AccAddress) int64 {
	oracleLastUpdateStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOracleLastUpdateHeightNative)
	return cosmosTypes.Int64FromBytes(oracleLastUpdateStore.Get(orchestratorAddress.Bytes()))
}

func (k Keeper) removeOracleLastUpdateHeightNative(ctx sdk.Context, orchestratorAddress sdk.AccAddress) {
	oracleLastUpdateStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOracleLastUpdateHeightNative)
	oracleLastUpdateStore.Delete(orchestratorAddress.Bytes())
}

//______________________________________________________________________________________________________________________

func (k Keeper) setOracleLastUpdateHeightCosmos(ctx sdk.Context, orchestratorAddress sdk.AccAddress, cosmosBlockHeight int64) {
	oracleLastUpdateStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOracleLastUpdateHeightCosmos)
	oracleLastUpdateStore.Set(orchestratorAddress.Bytes(), cosmosTypes.Int64Bytes(cosmosBlockHeight))
}

func (k Keeper) getOracleLastUpdateHeightCosmos(ctx sdk.Context, orchestratorAddress sdk.AccAddress) int64 {
	oracleLastUpdateStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOracleLastUpdateHeightCosmos)
	return cosmosTypes.Int64FromBytes(oracleLastUpdateStore.Get(orchestratorAddress.Bytes()))
}

func (k Keeper) removeOracleLastUpdateHeightCosmos(ctx sdk.Context, orchestratorAddress sdk.AccAddress) {
	oracleLastUpdateStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOracleLastUpdateHeightCosmos)
	oracleLastUpdateStore.Delete(orchestratorAddress.Bytes())
}
