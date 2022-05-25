package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

func (k Keeper) addToMintTokenStore(ctx sdk.Context, msg cosmosTypes.MsgMintTokensForAccount) {
	mintTokenStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyMintTokenStore)
	key := k.cdc.MustMarshal(&cosmosTypes.ChainIDHeightAndTxHashKey{ChainID: msg.ChainID, BlockHeight: msg.BlockHeight, TxHash: msg.TxHash})
	totalValidatorCount := k.GetTotalValidatorOrchestratorCount(ctx)

	// store has the key in it or not
	if !mintTokenStore.Has(key) {
		ratio := sdk.NewDec(1).Quo(sdk.NewDec(totalValidatorCount))
		mintTokenStoreValue := cosmosTypes.NewMintTokenStoreValue(msg, ratio, msg.OrchestratorAddress, ctx.BlockHeight()+cosmosTypes.StorageWindow)
		mintTokenStore.Set(key, k.cdc.MustMarshal(&mintTokenStoreValue))
		return
	}

	var mintTokenStoreValue cosmosTypes.MintTokenStoreValue
	k.cdc.MustUnmarshal(mintTokenStore.Get(key), &mintTokenStoreValue)

	// Match if the message value and stored value are same
	// if not equal then initialize by new value in store
	if !StoreValueEqualOrNotMintToken(mintTokenStoreValue, msg) {
		ratio := sdk.NewDec(1).Quo(sdk.NewDec(totalValidatorCount))
		newValue := cosmosTypes.NewMintTokenStoreValue(msg, ratio, msg.OrchestratorAddress, ctx.BlockHeight()+cosmosTypes.StorageWindow)
		mintTokenStore.Set(key, k.cdc.MustMarshal(&newValue))
		return
	}

	// if equal then check if orchestrator has already sent same details previously
	if !mintTokenStoreValue.Find(msg.OrchestratorAddress) {
		mintTokenStoreValue.UpdateValues(msg.OrchestratorAddress, totalValidatorCount)
		mintTokenStore.Set(key, k.cdc.MustMarshal(&mintTokenStoreValue))
	}
}

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

func (k Keeper) setMintedFlagInMintTokenStore(ctx sdk.Context, mv cosmosTypes.MintTokenStoreValue) {
	mintTokenStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyMintTokenStore)
	key := k.cdc.MustMarshal(&cosmosTypes.ChainIDHeightAndTxHashKey{ChainID: mv.MintTokens.ChainID, BlockHeight: mv.MintTokens.BlockHeight, TxHash: mv.MintTokens.TxHash})

	var mintTokenStoreValue cosmosTypes.MintTokenStoreValue
	k.cdc.MustUnmarshal(mintTokenStore.Get(key), &mintTokenStoreValue)
	mintTokenStoreValue.Minted = true

	mintTokenStore.Set(key, k.cdc.MustMarshal(&mintTokenStoreValue))
}

func (k Keeper) setAddedToEpochFlagInMintTokenStore(ctx sdk.Context, mv cosmosTypes.MintTokenStoreValue) {
	mintTokenStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyMintTokenStore)
	key := k.cdc.MustMarshal(&cosmosTypes.ChainIDHeightAndTxHashKey{ChainID: mv.MintTokens.ChainID, BlockHeight: mv.MintTokens.BlockHeight, TxHash: mv.MintTokens.TxHash})

	var mintTokenStoreValue cosmosTypes.MintTokenStoreValue
	k.cdc.MustUnmarshal(mintTokenStore.Get(key), &mintTokenStoreValue)
	mintTokenStoreValue.AddedToEpoch = true

	mintTokenStore.Set(key, k.cdc.MustMarshal(&mintTokenStoreValue))
}

func (k Keeper) deleteFromMintTokenStore(ctx sdk.Context, mv cosmosTypes.MintTokenStoreValue) {
	mintTokenStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyMintTokenStore)
	key := k.cdc.MustMarshal(&cosmosTypes.ChainIDHeightAndTxHashKey{ChainID: mv.MintTokens.ChainID, BlockHeight: mv.MintTokens.BlockHeight, TxHash: mv.MintTokens.TxHash})
	mintTokenStore.Delete(key)
}

func (k Keeper) ProcessAllMintingStoreValue(ctx sdk.Context) {
	listOfMintTokenStoreValue := k.getAllMintTokenStoreValue(ctx)
	for _, mv := range listOfMintTokenStoreValue {
		if mv.Ratio.GT(cosmosTypes.MinimumRatioForMajority) && !mv.AddedToEpoch && !mv.Minted {
			// step 1 : mint tokens for account
			err := k.mintTokensOnMajority(ctx, mv.MintTokens)
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
