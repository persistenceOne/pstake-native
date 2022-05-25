package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

func (k Keeper) setSlashingEventDetails(ctx sdk.Context, msg cosmosTypes.MsgSlashingEventOnCosmosChain) {
	slashingStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeySlashingStore)
	chainIDHeightAndTxHash := cosmosTypes.NewChainIDHeightAndTxHash(msg.ChainID, msg.BlockHeight, msg.TxHash)
	key := k.cdc.MustMarshal(&chainIDHeightAndTxHash)
	totalValidatorCount := k.GetTotalValidatorOrchestratorCount(ctx)

	// store has the key in it or not
	if !slashingStore.Has(key) {
		ratio := sdk.NewDec(1).Quo(sdk.NewDec(totalValidatorCount))
		newValue := cosmosTypes.NewSlashingStoreValue(msg, ratio, msg.OrchestratorAddress, ctx.BlockHeight()+cosmosTypes.StorageWindow)
		slashingStore.Set(key, k.cdc.MustMarshal(&newValue))
		return
	}

	var slashingStoreValue cosmosTypes.SlashingStoreValue
	k.cdc.MustUnmarshal(slashingStore.Get(key), &slashingStoreValue)

	// Match if the message value and stored value are same
	// if not equal then initialize by new value in store
	if !StoreValueEqualOrNotSlashingEvent(slashingStoreValue.SlashingDetails, msg) {
		ratio := sdk.NewDec(1).Quo(sdk.NewDec(totalValidatorCount))
		newValue := cosmosTypes.NewSlashingStoreValue(msg, ratio, msg.OrchestratorAddress, ctx.BlockHeight()+cosmosTypes.StorageWindow)
		slashingStore.Set(key, k.cdc.MustMarshal(&newValue))
		return
	}

	// if equal then check if orchestrator has already sent same details previously
	if !slashingStoreValue.Find(msg.OrchestratorAddress) {
		slashingStoreValue.UpdateValues(msg.OrchestratorAddress, totalValidatorCount)
		slashingStore.Set(key, k.cdc.MustMarshal(&slashingStoreValue))
	}
}

func (k Keeper) getAllSlashingEventDetails(ctx sdk.Context) (list []cosmosTypes.SlashingStoreValue) {
	slashingStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeySlashingStore)
	iterator := slashingStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var slashingStoreValue cosmosTypes.SlashingStoreValue
		k.cdc.MustUnmarshal(iterator.Value(), &slashingStoreValue)
		list = append(list, slashingStoreValue)
	}
	return list
}

func (k Keeper) deleteSlashingEventDetails(ctx sdk.Context, value cosmosTypes.SlashingStoreValue) {
	slashingStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeySlashingStore)
	chainIDHeightAndTxHash := cosmosTypes.NewChainIDHeightAndTxHash(value.SlashingDetails.ChainID, value.SlashingDetails.BlockHeight, value.SlashingDetails.TxHash)
	slashingStore.Delete(k.cdc.MustMarshal(&chainIDHeightAndTxHash))
}

func (k Keeper) ProcessAllSlashingEvents(ctx sdk.Context) {
	slashingEventList := k.getAllSlashingEventDetails(ctx)
	for _, se := range slashingEventList {
		if se.Ratio.GT(cosmosTypes.MinimumRatioForMajority) && !se.AddedToCValue {
			// todo : edit C value
		}
		if se.ActiveBlockHeight < ctx.BlockHeight() {
			k.deleteSlashingEventDetails(ctx, se)
		}
	}
}

func StoreValueEqualOrNotSlashingEvent(storeValue cosmosTypes.MsgSlashingEventOnCosmosChain, msgValue cosmosTypes.MsgSlashingEventOnCosmosChain) bool {
	if storeValue.ValidatorAddress != msgValue.ValidatorAddress {
		return false
	}
	if !storeValue.Amount.IsEqual(msgValue.Amount) {
		return false
	}
	if storeValue.TxHash != msgValue.TxHash {
		return false
	}
	if storeValue.ChainID != msgValue.ChainID {
		return false
	}
	if storeValue.BlockHeight != msgValue.BlockHeight {
		return false
	}
	return true
}
