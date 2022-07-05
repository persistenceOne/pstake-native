package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

// SetCosmosValidatorSet put a check on length of validator set and val set weights to maintain equal mapping
// sets the cosmos validator address as key and weight details as value
func (k Keeper) SetCosmosValidatorSet(ctx sdk.Context, cosmosValSetWeights []cosmosTypes.WeightedAddressAmount) {
	cosmosValSetStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyCosmosValidatorWeights)

	newSortedCosmosValSetWeights := cosmosTypes.NewWeightedAddressAmounts(cosmosValSetWeights).Sort()
	for i := range newSortedCosmosValSetWeights {
		valAddress, err := cosmosTypes.ValAddressFromBech32(newSortedCosmosValSetWeights[i].Address, cosmosTypes.Bech32PrefixValAddr)
		if err != nil {
			panic(err)
		}
		if cosmosValSetStore.Has(valAddress.Bytes()) {
			var weightedAddress cosmosTypes.WeightedAddressAmount
			k.cdc.MustUnmarshal(cosmosValSetStore.Get(valAddress.Bytes()), &weightedAddress)
			weightedAddress.Weight = newSortedCosmosValSetWeights[i].Weight
			cosmosValSetStore.Set(valAddress.Bytes(), k.cdc.MustMarshal(&weightedAddress))
		} else {
			bondDenom, err := k.GetParams(ctx).GetBondDenomOf("uatom")
			if err != nil {
				panic(err)
			}
			newSortedCosmosValSetWeights[i].Amount = sdk.ZeroInt()
			newSortedCosmosValSetWeights[i].Denom = bondDenom
			newSortedCosmosValSetWeights[i].UnbondingTokens = sdk.NewCoin(bondDenom, sdk.ZeroInt())
			cosmosValSetStore.Set(valAddress.Bytes(), k.cdc.MustMarshal(&newSortedCosmosValSetWeights[i]))
		}
	}
}

// SetCosmosValidatorWeight sets given weight details if store already has val address in it or panics
func (k Keeper) SetCosmosValidatorWeight(ctx sdk.Context, valAddress sdk.ValAddress, weight sdk.Dec) {
	cosmosValSetStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyCosmosValidatorWeights)
	if !cosmosValSetStore.Has(valAddress.Bytes()) {
		panic(fmt.Errorf("valAddress not present in kv store"))
	}
	var weightedAddress cosmosTypes.WeightedAddressAmount
	k.cdc.MustUnmarshal(cosmosValSetStore.Get(valAddress.Bytes()), &weightedAddress)
	weightedAddress.Weight = weight
	cosmosValSetStore.Set(valAddress.Bytes(), k.cdc.MustMarshal(&weightedAddress))
}

// GetAllCosmosValidatorSet gets all the cosmos validator set details
func (k Keeper) GetAllCosmosValidatorSet(ctx sdk.Context) (weightedAddresses cosmosTypes.WeightedAddressAmounts) {
	cosmosValSetStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyCosmosValidatorWeights)
	iterator := cosmosValSetStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var weightedAddress cosmosTypes.WeightedAddressAmount
		k.cdc.MustUnmarshal(iterator.Value(), &weightedAddress)
		weightedAddresses = append(weightedAddresses, weightedAddress)
	}
	return weightedAddresses.Sort()
}

// UpdateDelegationCosmosValidator updates the delegation of given cosmos validator
func (k Keeper) UpdateDelegationCosmosValidator(ctx sdk.Context, valAddress sdk.ValAddress, amount sdk.Coin, unbondingAmount sdk.Coin) {
	cosmosValSetStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyCosmosValidatorWeights)
	if !cosmosValSetStore.Has(valAddress.Bytes()) {
		panic(fmt.Errorf("valAddress not present in kv store"))
	}
	var weightedAddress cosmosTypes.WeightedAddressAmount
	k.cdc.MustUnmarshal(cosmosValSetStore.Get(valAddress.Bytes()), &weightedAddress)
	weightedAddress.Amount = amount.Amount
	weightedAddress.Denom = amount.Denom
	if !unbondingAmount.IsZero() {
		weightedAddress.UnbondingTokens = unbondingAmount
	}
	cosmosValSetStore.Set(valAddress.Bytes(), k.cdc.MustMarshal(&weightedAddress))
}

// getDelegationCosmosValidator Gets the delegation of given cosmos validator
func (k Keeper) getDelegationCosmosValidator(ctx sdk.Context, valAddress sdk.ValAddress) sdk.Coin {
	cosmosValSetStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyCosmosValidatorWeights)
	if !cosmosValSetStore.Has(valAddress.Bytes()) {
		panic(fmt.Errorf("valAddress not present in kv store"))
	}
	var weightedAddress cosmosTypes.WeightedAddressAmount
	k.cdc.MustUnmarshal(cosmosValSetStore.Get(valAddress.Bytes()), &weightedAddress)
	return weightedAddress.Coin()
}

// removeCosmosValidatorWeight removes the given validator from the cosmos validator set
func (k Keeper) removeCosmosValidatorWeight(ctx sdk.Context, valAddress sdk.ValAddress) {
	cosmosValSetStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyCosmosValidatorWeights)
	cosmosValSetStore.Delete(valAddress.Bytes())
}

//______________________________________________________________________________________________________________________

// setOracleValidatorSet put a check on length of validator set and val set weights to maintain equal mapping
// sets the native validator address as key and weight details as value
func (k Keeper) setOracleValidatorSet(ctx sdk.Context, valAddresses []sdk.ValAddress, nativeValSetWeights []cosmosTypes.WeightedAddress) {
	nativeValSetStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyNativeValidatorWeights)

	for i, va := range valAddresses {
		nativeValSetStore.Set(va.Bytes(), k.cdc.MustMarshal(&nativeValSetWeights[i]))
	}
}

// setOracleValidatorWeight sets given weight details if store already has val address in it or panics
func (k Keeper) setOracleValidatorWeight(ctx sdk.Context, valAddress sdk.ValAddress, weight cosmosTypes.WeightedAddress) {
	nativeValSetStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyNativeValidatorWeights)
	if !nativeValSetStore.Has(valAddress.Bytes()) {
		panic(fmt.Errorf("valAddress not present in kv store"))
	}
	nativeValSetStore.Set(valAddress.Bytes(), k.cdc.MustMarshal(&weight))
}

// getAllOracleValidatorSet gets the list of all oracle validator set
func (k Keeper) getAllOracleValidatorSet(ctx sdk.Context) (weightedAddresses []cosmosTypes.WeightedAddress) {
	nativeValSetStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyNativeValidatorWeights)
	iterator := nativeValSetStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var weightedAddress cosmosTypes.WeightedAddress
		k.cdc.MustUnmarshal(iterator.Value(), &weightedAddress)
		weightedAddresses = append(weightedAddresses, weightedAddress)
	}
	return weightedAddresses
}

// removeOracleValidatorWeight removes the given validator from the native validator set
func (k Keeper) removeOracleValidatorWeight(ctx sdk.Context, valAddress sdk.ValAddress) {
	nativeValSetStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyNativeValidatorWeights)
	nativeValSetStore.Delete(valAddress.Bytes())
}
