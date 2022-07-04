package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

/*
Adds the minting message entry to the minting store with the given validator address.
Performs the following actions :
  1. Checks if store has the key or not. If not then create new entry
  2. Checks if store has it and matches all the details present in the message. If not then create a new entry.
  3. Finally, if all the details match then append the validator address to keep track.
*/
func (k Keeper) addToMintTokenStore(ctx sdk.Context, msg cosmosTypes.MsgMintTokensForAccount, validatorAddress sdk.ValAddress) {
	mintTokenStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyMintTokenStore)
	key := k.cdc.MustMarshal(&cosmosTypes.ChainIDHeightAndTxHashKey{ChainID: msg.ChainID, BlockHeight: msg.BlockHeight, TxHash: msg.TxHash})
	totalValidatorCount := k.GetTotalValidatorOrchestratorCount(ctx)

	// store has the key in it or not
	if !mintTokenStore.Has(key) {
		ratio := sdk.NewDec(1).Quo(sdk.NewDec(totalValidatorCount))
		mintTokenStoreValue := cosmosTypes.NewMintTokenStoreValue(msg, ratio, validatorAddress, ctx.BlockHeight()+cosmosTypes.StorageWindow)
		mintTokenStore.Set(key, k.cdc.MustMarshal(&mintTokenStoreValue))
		return
	}

	var mintTokenStoreValue cosmosTypes.MintTokenStoreValue
	k.cdc.MustUnmarshal(mintTokenStore.Get(key), &mintTokenStoreValue)

	// Match if the message value and stored value are same
	// if not equal then initialize by new value in store
	if !StoreValueEqualOrNotMintToken(mintTokenStoreValue, msg) {
		ratio := sdk.NewDec(1).Quo(sdk.NewDec(totalValidatorCount))
		newValue := cosmosTypes.NewMintTokenStoreValue(msg, ratio, validatorAddress, ctx.BlockHeight()+cosmosTypes.StorageWindow)
		mintTokenStore.Set(key, k.cdc.MustMarshal(&newValue))
		return
	}

	// if equal then check if orchestrator has already sent same details previously
	if !mintTokenStoreValue.Find(validatorAddress.String()) {
		mintTokenStoreValue.UpdateValues(validatorAddress.String(), totalValidatorCount)
		mintTokenStore.Set(key, k.cdc.MustMarshal(&mintTokenStoreValue))
	}
}

//Gets all the entries in the mint token store. Used in processing all the mint requests.
func (k Keeper) getAllMintTokenStoreValue(ctx sdk.Context) (list []cosmosTypes.MintTokenStoreValue) {
	mintTokenStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyMintTokenStore)
	iterator := mintTokenStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var mintTokenStoreValue cosmosTypes.MintTokenStoreValue
		k.cdc.MustUnmarshal(iterator.Value(), &mintTokenStoreValue)
		list = append(list, mintTokenStoreValue)
	}
	return list
}

// Set the minted flag true. Used when the minting is successful for the given request
func (k Keeper) setMintedFlagInMintTokenStore(ctx sdk.Context, mv cosmosTypes.MintTokenStoreValue) {
	mintTokenStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyMintTokenStore)
	key := k.cdc.MustMarshal(&cosmosTypes.ChainIDHeightAndTxHashKey{ChainID: mv.MintTokens.ChainID, BlockHeight: mv.MintTokens.BlockHeight, TxHash: mv.MintTokens.TxHash})

	var mintTokenStoreValue cosmosTypes.MintTokenStoreValue
	k.cdc.MustUnmarshal(mintTokenStore.Get(key), &mintTokenStoreValue)
	mintTokenStoreValue.Minted = true

	mintTokenStore.Set(key, k.cdc.MustMarshal(&mintTokenStoreValue))
}

// Sets added to epoch flag true. Used when the amount has been added to epoch store for "uatom".
func (k Keeper) setAddedToEpochFlagInMintTokenStore(ctx sdk.Context, mv cosmosTypes.MintTokenStoreValue) {
	mintTokenStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyMintTokenStore)
	key := k.cdc.MustMarshal(&cosmosTypes.ChainIDHeightAndTxHashKey{ChainID: mv.MintTokens.ChainID, BlockHeight: mv.MintTokens.BlockHeight, TxHash: mv.MintTokens.TxHash})

	var mintTokenStoreValue cosmosTypes.MintTokenStoreValue
	k.cdc.MustUnmarshal(mintTokenStore.Get(key), &mintTokenStoreValue)
	mintTokenStoreValue.AddedToEpoch = true

	mintTokenStore.Set(key, k.cdc.MustMarshal(&mintTokenStoreValue))
}

// removes the details set in the mint token store. Used when the active block height is reached
func (k Keeper) deleteFromMintTokenStore(ctx sdk.Context, mv cosmosTypes.MintTokenStoreValue) {
	mintTokenStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyMintTokenStore)
	key := k.cdc.MustMarshal(&cosmosTypes.ChainIDHeightAndTxHashKey{ChainID: mv.MintTokens.ChainID, BlockHeight: mv.MintTokens.BlockHeight, TxHash: mv.MintTokens.TxHash})
	mintTokenStore.Delete(key)
}

/*
ProcessAllMintingStoreValue processes all the minting requests
This function is called every EndBlocker to perform the defined set of actions as mentioned below :
   1. Get the list of all minting requests
   2. Checks if the majority of the validator oracle have sent the minting request. Also checks the addedToEpoch and Minted flag.
   3. If majority is reached and other conditions match then tokens are minted and flags are set to true
   4. Another condition of ActiveBlockHeight is also checked whether to delete the entry or not.
*/
func (k Keeper) ProcessAllMintingStoreValue(ctx sdk.Context) {
	listOfMintTokenStoreValue := k.getAllMintTokenStoreValue(ctx)
	for _, mv := range listOfMintTokenStoreValue {
		if mv.Ratio.GT(cosmosTypes.MinimumRatioForMajority) && !mv.AddedToEpoch && !mv.Minted {
			// step 1 : mint tokens for account
			err := k.mintTokens(ctx, mv.MintTokens)
			if err != nil {
				panic(err)
			}
			// step 2 : mark minted flag true
			k.setMintedFlagInMintTokenStore(ctx, mv)
			// step 3 : add to current epoch for staking
			k.addToStakingEpoch(ctx, mv.MintTokens.Amount)
			// step 4 : make added to epoch flag true
			k.setAddedToEpochFlagInMintTokenStore(ctx, mv)
		}

		if mv.ActiveBlockHeight < ctx.BlockHeight() {
			k.deleteFromMintTokenStore(ctx, mv)
		}
	}
}

// StoreValueEqualOrNotMintToken Helper function for mint token store to check if the relevant details in the message matches or not.
func StoreValueEqualOrNotMintToken(storeValue cosmosTypes.MintTokenStoreValue, msgValue cosmosTypes.MsgMintTokensForAccount) bool {
	if storeValue.MintTokens.AddressFromMemo != msgValue.AddressFromMemo {
		return false
	}
	if !storeValue.MintTokens.Amount.IsEqual(msgValue.Amount) {
		return false
	}
	if storeValue.MintTokens.TxHash != msgValue.TxHash {
		return false
	}
	if storeValue.MintTokens.ChainID != msgValue.ChainID {
		return false
	}
	if storeValue.MintTokens.BlockHeight != msgValue.BlockHeight {
		return false
	}
	return true
}
