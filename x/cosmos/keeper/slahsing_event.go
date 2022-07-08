package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

/*
setSlashingEventDetails Adds the slashing message entry to the slashing store with the given validator address.
Performs the following actions :
  1. Checks if store has the key or not. If not then create new entry
  2. Checks if store has it and matches all the details present in the message. If not then create a new entry.
  3. Finally, if all the details match then append the validator address to keep track.
*/
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

// setAddedToCValueTrue Sets the addedToCValue flag true for thw given slashing store value
func (k Keeper) setAddedToCValueTrue(ctx sdk.Context, value cosmosTypes.SlashingStoreValue) {
	slashingStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeySlashingStore)
	chainIDHeightAndTxHash := cosmosTypes.NewChainIDHeightAndTxHash(value.SlashingDetails.ChainID, value.SlashingDetails.BlockHeight, value.SlashingDetails.SlashType)

	var slashingStoreValue cosmosTypes.SlashingStoreValue
	k.cdc.MustUnmarshal(slashingStore.Get(k.cdc.MustMarshal(&chainIDHeightAndTxHash)), &slashingStoreValue)

	slashingStoreValue.AddedToCValue = true
	slashingStore.Set(k.cdc.MustMarshal(&chainIDHeightAndTxHash), k.cdc.MustMarshal(&slashingStoreValue))
}

// Gets all the slashing event details present in the slashing store
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

// deleteSlashingEventDetails Removes the slashing event details corresponding to the passed values
func (k Keeper) deleteSlashingEventDetails(ctx sdk.Context, value cosmosTypes.SlashingStoreValue) {
	slashingStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeySlashingStore)
	chainIDHeightAndTxHash := cosmosTypes.NewChainIDHeightAndTxHash(value.SlashingDetails.ChainID, value.SlashingDetails.BlockHeight, value.SlashingDetails.SlashType)
	slashingStore.Delete(k.cdc.MustMarshal(&chainIDHeightAndTxHash))
}

/*
ProcessAllSlashingEvents processes all the slashing requests
This function is called every EndBlocker to perform the defined set of actions as mentioned below :
   1. Get the list of all slashing requests
   2. Checks if the majority of the validator oracle have sent the minting request. Also checks the addedToCValue flag.
   3. If majority is reached and other conditions match then slashing event is accepted and C value is updated
   4. Another condition of ActiveBlockHeight is also checked whether to delete the entry or not.
*/
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
			// set added to C value true
			k.setAddedToCValueTrue(ctx, se)
			// update current delegations of the validator with the delegation supplied with message
			// supply a zero coin with any denom in order to keep the unbonding delegations same
			k.UpdateDelegationCosmosValidator(ctx, valAddress, se.SlashingDetails.CurrentDelegation, sdk.NewCoin("test", sdk.ZeroInt()))
		}
		if se.ActiveBlockHeight < ctx.BlockHeight() && se.AddedToCValue {
			k.deleteSlashingEventDetails(ctx, se)
		}
	}
}

// StoreValueEqualOrNotSlashingEvent Helper function for slashing store to check if the relevant details in the message matches or not.
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
