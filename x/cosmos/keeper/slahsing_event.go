package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

func (k Keeper) setSlashingEventDetails(ctx sdk.Context, msg cosmosTypes.MsgSlashingEventOnCosmosChain, validatorAddress sdk.ValAddress) {
	slashingStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeySlashingStore)
	chainIDHeightAndTxHash := cosmosTypes.NewChainIDHeightAndTxHash(msg.ChainID, msg.BlockHeight, msg.SlashType)
	key := k.cdc.MustMarshal(&chainIDHeightAndTxHash)
	totalValidatorCount := k.GetTotalValidatorOrchestratorCount(ctx)

	// store has the key in it or not
	if !slashingStore.Has(key) {
		ratio := sdk.NewDec(1).Quo(sdk.NewDec(totalValidatorCount))
		newValue := cosmosTypes.NewSlashingStoreValue(msg, ratio, validatorAddress, ctx.BlockHeight()+cosmosTypes.StorageWindow)
		slashingStore.Set(key, k.cdc.MustMarshal(&newValue))
		return
	}

	var slashingStoreValue cosmosTypes.SlashingStoreValue
	k.cdc.MustUnmarshal(slashingStore.Get(key), &slashingStoreValue)

	// Match if the message value and stored value are same
	// if not equal then initialize by new value in store
	if !StoreValueEqualOrNotSlashingEvent(slashingStoreValue.SlashingDetails, msg) {
		ratio := sdk.NewDec(1).Quo(sdk.NewDec(totalValidatorCount))
		newValue := cosmosTypes.NewSlashingStoreValue(msg, ratio, validatorAddress, ctx.BlockHeight()+cosmosTypes.StorageWindow)
		slashingStore.Set(key, k.cdc.MustMarshal(&newValue))
		return
	}

	// if equal then check if orchestrator has already sent same details previously
	if !slashingStoreValue.Find(validatorAddress.String()) {
		slashingStoreValue.UpdateValues(validatorAddress.String(), totalValidatorCount)
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
	chainIDHeightAndTxHash := cosmosTypes.NewChainIDHeightAndTxHash(value.SlashingDetails.ChainID, value.SlashingDetails.BlockHeight, value.SlashingDetails.SlashType)
	slashingStore.Delete(k.cdc.MustMarshal(&chainIDHeightAndTxHash))
}

func (k Keeper) ProcessAllSlashingEvents(ctx sdk.Context) {
	slashingEventList := k.getAllSlashingEventDetails(ctx)
	for _, se := range slashingEventList {
		if se.Ratio.GT(cosmosTypes.MinimumRatioForMajority) && !se.AddedToCValue {
			valAddress, err := cosmosTypes.ValAddressFromBech32(se.SlashingDetails.ValidatorAddress, cosmosTypes.Bech32PrefixValAddr)
			if err != nil {
				panic(err)
			}

			// get current delegation of validator
			delegation := k.getDelegationCosmosValidator(ctx, valAddress)
			// calculate slashed amount
			slashedAmount := delegation.Sub(se.SlashingDetails.CurrentDelegation)
			// update C value based on slashed amount
			k.SlashingEvent(ctx, slashedAmount)
			// update current delegations of the validator with the delegation supplied with message
			k.UpdateDelegationCosmosValidator(ctx, valAddress, se.SlashingDetails.CurrentDelegation)
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
	if !storeValue.CurrentDelegation.IsEqual(msgValue.CurrentDelegation) {
		return false
	}
	if storeValue.SlashType != msgValue.SlashType {
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
