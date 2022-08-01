package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

// setOrchestratorLastUpdateHeightNative Sets orchestrator last update height as the last block on which orchestrator sent a message on native side
func (k Keeper) setOrchestratorLastUpdateHeightNative(ctx sdk.Context, orchestratorAddress sdk.AccAddress, nativeBlockHeight int64) {
	orchestratorLastUpdateStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOrchestratorLastUpdateHeightNative)
	orchestratorLastUpdateStore.Set(orchestratorAddress.Bytes(), cosmosTypes.Int64Bytes(nativeBlockHeight))
}

// getOrchestratorLastUpdateHeightNative Gets the last update height of orchestrator on native side
func (k Keeper) getOrchestratorLastUpdateHeightNative(ctx sdk.Context, orchestratorAddress sdk.AccAddress) int64 {
	orchestratorLastUpdateStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOrchestratorLastUpdateHeightNative)
	return cosmosTypes.Int64FromBytes(orchestratorLastUpdateStore.Get(orchestratorAddress.Bytes()))
}

// removeOrchestratorLastUpdateHeightNative Removes the entry of orchestrator from the DB
func (k Keeper) removeOrchestratorLastUpdateHeightNative(ctx sdk.Context, orchestratorAddress sdk.AccAddress) {
	orchestratorLastUpdateStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOrchestratorLastUpdateHeightNative)
	orchestratorLastUpdateStore.Delete(orchestratorAddress.Bytes())
}

//______________________________________________________________________________________________________________________

// setOrchestratorLastUpdateHeightCosmos Sets orchestrator last update height as the last block on which orchestrator sent a message on cosmos side
func (k Keeper) setOrchestratorLastUpdateHeightCosmos(ctx sdk.Context, orchestratorAddress sdk.AccAddress, cosmosBlockHeight int64) {
	orchestratorLastUpdateStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOrchestratorLastUpdateHeightCosmos)
	orchestratorLastUpdateStore.Set(orchestratorAddress.Bytes(), cosmosTypes.Int64Bytes(cosmosBlockHeight))
}

// getOrchestratorLastUpdateHeightCosmos Gets the last update height of orchestrator on cosmos side
func (k Keeper) getOrchestratorLastUpdateHeightCosmos(ctx sdk.Context, orchestratorAddress sdk.AccAddress) int64 {
	orchestratorLastUpdateStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOrchestratorLastUpdateHeightCosmos)
	return cosmosTypes.Int64FromBytes(orchestratorLastUpdateStore.Get(orchestratorAddress.Bytes()))
}

// removeOrchestratorLastUpdateHeightCosmos Removes the entry of orchestrator from the DB
func (k Keeper) removeOrchestratorLastUpdateHeightCosmos(ctx sdk.Context, orchestratorAddress sdk.AccAddress) {
	orchestratorLastUpdateStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOrchestratorLastUpdateHeightCosmos)
	orchestratorLastUpdateStore.Delete(orchestratorAddress.Bytes())
}
