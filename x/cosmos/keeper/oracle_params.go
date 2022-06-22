package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

// Sets oracle last update height as the last block on which oracle sent a message on native side
func (k Keeper) setOracleLastUpdateHeightNative(ctx sdk.Context, orchestratorAddress sdk.AccAddress, nativeBlockHeight int64) {
	oracleLastUpdateStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOracleLastUpdateHeightNative)
	oracleLastUpdateStore.Set(orchestratorAddress.Bytes(), cosmosTypes.Int64Bytes(nativeBlockHeight))
}

// Gets the last update height of oracle on native side
func (k Keeper) getOracleLastUpdateHeightNative(ctx sdk.Context, orchestratorAddress sdk.AccAddress) int64 {
	oracleLastUpdateStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOracleLastUpdateHeightNative)
	return cosmosTypes.Int64FromBytes(oracleLastUpdateStore.Get(orchestratorAddress.Bytes()))
}

// Removes the entry of oracle from the DB
func (k Keeper) removeOracleLastUpdateHeightNative(ctx sdk.Context, orchestratorAddress sdk.AccAddress) {
	oracleLastUpdateStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOracleLastUpdateHeightNative)
	oracleLastUpdateStore.Delete(orchestratorAddress.Bytes())
}

//______________________________________________________________________________________________________________________

// Sets oracle last update height as the last block on which oracle sent a message on cosmos side
func (k Keeper) setOracleLastUpdateHeightCosmos(ctx sdk.Context, orchestratorAddress sdk.AccAddress, cosmosBlockHeight int64) {
	oracleLastUpdateStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOracleLastUpdateHeightCosmos)
	oracleLastUpdateStore.Set(orchestratorAddress.Bytes(), cosmosTypes.Int64Bytes(cosmosBlockHeight))
}

// Gets the last update height of oracle on cosmos side
func (k Keeper) getOracleLastUpdateHeightCosmos(ctx sdk.Context, orchestratorAddress sdk.AccAddress) int64 {
	oracleLastUpdateStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOracleLastUpdateHeightCosmos)
	return cosmosTypes.Int64FromBytes(oracleLastUpdateStore.Get(orchestratorAddress.Bytes()))
}

// Removes the entry of oracle from the DB
func (k Keeper) removeOracleLastUpdateHeightCosmos(ctx sdk.Context, orchestratorAddress sdk.AccAddress) {
	oracleLastUpdateStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOracleLastUpdateHeightCosmos)
	oracleLastUpdateStore.Delete(orchestratorAddress.Bytes())
}
