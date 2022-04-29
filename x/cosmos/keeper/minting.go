package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

/*
//TODO : add mint pool structure as comment
*/
// add a transaction to minting pool for tallying how many orchs have sent a request for minting
func (k Keeper) addToMintingPoolTx(ctx sdkTypes.Context, txHash string, destinationAddress sdkTypes.AccAddress, orchestratorAddress sdkTypes.AccAddress, amount sdkTypes.Coins) error {
	store := ctx.KVStore(k.storeKey)
	mintingPoolStore := prefix.NewStore(store, []byte(cosmosTypes.MintingPoolStoreKey))
	key := []byte(cosmosTypes.GetDestinationAddressAmountAndTxHashKey(destinationAddress, amount, txHash))
	if mintingPoolStore.Has(key) {
		var txnDetails cosmosTypes.IncomingMintTx
		bz := mintingPoolStore.Get(key)
		err := k.cdc.Unmarshal(bz, &txnDetails)
		if err != nil {
			return err
		}

		found := txnDetails.Find(orchestratorAddress.String())
		if !found {
			txnDetails.AddAndIncrement(orchestratorAddress.String())
		}

		bz, err = k.cdc.Marshal(&txnDetails)
		if err != nil {
			return err
		}
		mintingPoolStore.Set(key, bz)
	} else {
		txnDetails := cosmosTypes.NewIncomingMintTx(orchestratorAddress, 1)
		bz, _ := k.cdc.Marshal(&txnDetails)
		mintingPoolStore.Set(key, bz)
	}
	return nil
}

// Fetches the list of items in minting pool
func (k Keeper) fetchFromMintPoolTx(ctx sdkTypes.Context, keyAndValueForMinting []cosmosTypes.KeyAndValueForMinting) []cosmosTypes.KeyAndValueForMinting {
	store := ctx.KVStore(k.storeKey)
	mintingPoolStore := prefix.NewStore(store, []byte(cosmosTypes.MintingPoolStoreKey))
	totalCount := float32(k.getTotalValidatorOrchestratorCount(ctx))
	for i := range keyAndValueForMinting {
		destinationAddress, err := sdkTypes.AccAddressFromBech32(keyAndValueForMinting[i].Value.DestinationAddress)
		if err != nil {
			panic("Error in parsing destination address")
		}

		key := []byte(cosmosTypes.GetDestinationAddressAmountAndTxHashKey(destinationAddress, sdkTypes.NewCoins(keyAndValueForMinting[i].Value.Amount), keyAndValueForMinting[i].Key.TxHash))
		bz := mintingPoolStore.Get(key)

		var txnDetails cosmosTypes.IncomingMintTx
		err = k.cdc.Unmarshal(bz, &txnDetails)
		if err != nil {
			panic("Error in unmarshalling txn Details")
		}

		keyAndValueForMinting[i].Ratio = float32(len(txnDetails.OrchAddresses)) / totalCount

	}
	return keyAndValueForMinting
}

// deletes an item from mint pool
func (k Keeper) deleteFromMintPoolTx(ctx sdkTypes.Context, destinationAddress sdkTypes.AccAddress, amount sdkTypes.Coin, txHash string) {
	store := ctx.KVStore(k.storeKey)
	mintingPoolStore := prefix.NewStore(store, []byte(cosmosTypes.MintingPoolStoreKey))
	mintingPoolStore.Delete([]byte(cosmosTypes.GetDestinationAddressAmountAndTxHashKey(destinationAddress, sdkTypes.NewCoins(amount), txHash)))
}

//______________________________________________________________________________________________
/*
TODO : Add structure
*/
func (k Keeper) setMintAddressAndAmount(ctx sdkTypes.Context, chainID string, blockHeight int64, txHash string, destinationAddress sdkTypes.AccAddress, amount sdkTypes.Coin) {
	store := ctx.KVStore(k.storeKey)
	mintAddressAndAmountStore := prefix.NewStore(store, []byte(cosmosTypes.AddressAndAmountStoreKey))

	chainIDHeightAndTxHash := cosmosTypes.NewChainIDHeightAndTxHash(chainID, blockHeight, txHash)
	key, err := k.cdc.Marshal(&chainIDHeightAndTxHash)
	if err != nil {
		panic("error in marshaling chainID, height and txHash")
	}

	if !mintAddressAndAmountStore.Has(key) {
		addressAndAmount := cosmosTypes.NewAddressAndAmount(destinationAddress, amount, ctx.BlockHeight())
		bz, err := k.cdc.Marshal(&addressAndAmount)
		if err != nil {
			panic("error in marshaling address and amount")
		}
		mintAddressAndAmountStore.Set(key, bz)
	}

}

func (k Keeper) getAllMintAddressAndAmount(ctx sdkTypes.Context) (list []cosmosTypes.KeyAndValueForMinting, err error) {
	store := ctx.KVStore(k.storeKey)
	mintAddressAndAmountStore := prefix.NewStore(store, []byte(cosmosTypes.AddressAndAmountStoreKey))

	iterator := mintAddressAndAmountStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var chainIDHeightAndTxHash cosmosTypes.ChainIDHeightAndTxHashKey

		err = k.cdc.Unmarshal(iterator.Key(), &chainIDHeightAndTxHash)
		if err != nil {
			return nil, err
		}

		var addressAndAmount cosmosTypes.AddressAndAmountKey

		err = k.cdc.Unmarshal(iterator.Value(), &addressAndAmount)
		if err != nil {
			return nil, err
		}

		a := cosmosTypes.KeyAndValueForMinting{
			Key:   chainIDHeightAndTxHash,
			Value: addressAndAmount,
		}

		list = append(list, a)
	}
	return list, nil
}

func (k Keeper) deleteMintedAddressAndAmountKeys(ctx sdkTypes.Context, keyHash cosmosTypes.ChainIDHeightAndTxHashKey) {
	store := ctx.KVStore(k.storeKey)
	mintAddressAndAmountStore := prefix.NewStore(store, []byte(cosmosTypes.AddressAndAmountStoreKey))

	chainIDHeightAndTxHash := cosmosTypes.NewChainIDHeightAndTxHash(keyHash.ChainID, keyHash.BlockHeight, keyHash.TxHash)
	key, err := k.cdc.Marshal(&chainIDHeightAndTxHash)
	if err != nil {
		panic("error in marshaling chainID, height and txHash")
	}

	mintAddressAndAmountStore.Delete(key)
}

func (k Keeper) setMintedFlagTrue(ctx sdkTypes.Context, keyHash cosmosTypes.ChainIDHeightAndTxHashKey) {
	store := ctx.KVStore(k.storeKey)
	mintAddressAndAmountStore := prefix.NewStore(store, []byte(cosmosTypes.AddressAndAmountStoreKey))

	chainIDHeightAndTxHash := cosmosTypes.NewChainIDHeightAndTxHash(keyHash.ChainID, keyHash.BlockHeight, keyHash.TxHash)
	key, err := k.cdc.Marshal(&chainIDHeightAndTxHash)
	if err != nil {
		panic("error in marshaling chainID, height and txHash")
	}

	bz := mintAddressAndAmountStore.Get(key)
	var a cosmosTypes.AddressAndAmountKey
	err = k.cdc.Unmarshal(bz, &a)
	if err != nil {
		panic("error in unmarshalling address and amount")
	}
	a.Minted = true

	bz, err = k.cdc.Marshal(&a)
	if err != nil {
		panic("error in marshaling address and amount")
	}

	mintAddressAndAmountStore.Set(key, bz)
}

func (k Keeper) setAcknowledgmentFlagTrue(ctx sdkTypes.Context, keyHash cosmosTypes.ChainIDHeightAndTxHashKey) {
	store := ctx.KVStore(k.storeKey)
	mintAddressAndAmountStore := prefix.NewStore(store, []byte(cosmosTypes.AddressAndAmountStoreKey))

	chainIDHeightAndTxHash := cosmosTypes.NewChainIDHeightAndTxHash(keyHash.ChainID, keyHash.BlockHeight, keyHash.TxHash)
	key, err := k.cdc.Marshal(&chainIDHeightAndTxHash)
	if err != nil {
		panic("error in marshaling chainID, height and txHash")
	}

	bz := mintAddressAndAmountStore.Get(key)
	var a cosmosTypes.AddressAndAmountKey
	err = k.cdc.Unmarshal(bz, &a)
	if err != nil {
		panic("error in unmarshalling address and amount")
	}
	a.Acknowledgment = true

	bz, err = k.cdc.Marshal(&a)
	if err != nil {
		panic("error in marshaling address and amount")
	}

	mintAddressAndAmountStore.Set(key, bz)
}

//______________________________________________________________________________________________

// ProcessAllMintingTransactions Process all minting transactions
func (k Keeper) ProcessAllMintingTransactions(ctx sdkTypes.Context) error {
	listNew, err := k.getAllMintAddressAndAmount(ctx)
	if err != nil {
		return err
	}
	listWithRatio := k.fetchFromMintPoolTx(ctx, listNew)

	for _, addressToMintTokens := range listWithRatio {
		if addressToMintTokens.Ratio > cosmosTypes.MinimumRatioForMajority && !addressToMintTokens.Value.Acknowledgment {
			addressToMintTokens.Value.Acknowledgment = true
			k.addToStakingEpoch(ctx, addressToMintTokens)
			k.setAcknowledgmentFlagTrue(ctx, addressToMintTokens.Key)
		}

		if addressToMintTokens.Value.NativeBlockHeight+cosmosTypes.StorageWindow < ctx.BlockHeight() {
			k.deleteMintedAddressAndAmountKeys(ctx, addressToMintTokens.Key)
			destinationAddress, err := sdkTypes.AccAddressFromBech32(addressToMintTokens.Value.DestinationAddress)
			if err != nil {
				return err
			}
			k.deleteFromMintPoolTx(ctx, destinationAddress, addressToMintTokens.Value.Amount, addressToMintTokens.Key.TxHash)
		}
	}

	epochNumberAndTxIDStatusList, err := k.fetchAllInEpochPoolForMinting(ctx)
	for _, en := range epochNumberAndTxIDStatusList {
		count := 0
		for _, txIDStatus := range en.mintingEpochValue.TxIDAndStatus {
			if txIDStatus.Status == true {
				count++
			}
		}
		if count == len(en.mintingEpochValue.TxIDAndStatus) {
			// distribute to be minted tokens from the all deposits in the given epoch
			stakingEpochValue, err := k.getFromStakingEpoch(ctx, en.epochNumber)
			if err != nil {
				return err
			}
			for _, e := range stakingEpochValue.EpochMintingTxns {
				err = k.mintTokensOnMajority(ctx, e.Key, e.Value)
				if err != nil {
					return err
				}
			}

			// process all rewards distribution
			rewardAmnt, err := k.getFromRewardsInCurrentEpochAmount(ctx, en.epochNumber)
			if err != nil {
				return err
			}
			err = k.processAllRewardsClaimed(ctx, rewardAmnt)
			if err != nil {
				return err
			}
		}
		k.deleteFromRewardsInCurrentEpoch(ctx, en.epochNumber)
		k.deleteFromStakingEpoch(ctx, en.epochNumber)
		k.deleteInEpochPoolForMinting(ctx, en.epochNumber)
	}
	return nil
}
