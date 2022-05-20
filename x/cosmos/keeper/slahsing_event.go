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
		newValue := cosmosTypes.NewSlashingStoreValue(msg, ratio, msg.OrchestratorAddress)
		slashingStore.Set(key, k.cdc.MustMarshal(&newValue))
		return
	}

	var slashingStoreValue cosmosTypes.SlashingStoreValue
	k.cdc.MustUnmarshal(slashingStore.Get(key), &slashingStoreValue)

	// Match if the message value and stored value are same
	// if not equal then initialize by new value in store
	if !StoreValueEqualOrNotSlashingEvent(slashingStoreValue.SlashingDetails, msg) {
		ratio := sdk.NewDec(1).Quo(sdk.NewDec(totalValidatorCount))
		newValue := cosmosTypes.NewSlashingStoreValue(msg, ratio, msg.OrchestratorAddress)
		slashingStore.Set(key, k.cdc.MustMarshal(&newValue))
		return
	}

	// if equal then check if orchestrator has already sent same details previously
	if !slashingStoreValue.Find(msg.OrchestratorAddress) {
		slashingStoreValue.UpdateValues(msg.OrchestratorAddress, totalValidatorCount)
		slashingStore.Set(key, k.cdc.MustMarshal(&slashingStoreValue))
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
