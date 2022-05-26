package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

// put a check on length of validator set and val set weights to maintain equal mapping
// sets the cosmos validator address as key and weight details as value
func (k Keeper) setCosmosValidatorSet(ctx sdk.Context, cosmosValSetWeights []cosmosTypes.WeightedAddressAmount) {
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
			cosmosValSetStore.Set(valAddress.Bytes(), k.cdc.MustMarshal(&newSortedCosmosValSetWeights[i]))
		}
	}
}

// sets given weight details if store already has val address in it or panics
func (k Keeper) setCosmosValidatorWeight(ctx sdk.Context, valAddress sdk.ValAddress, weight sdk.Dec) {
	cosmosValSetStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyCosmosValidatorWeights)
	if !cosmosValSetStore.Has(valAddress.Bytes()) {
		panic(fmt.Errorf("valAddress not present in kv store"))
	}
	var weightedAddress cosmosTypes.WeightedAddressAmount
	k.cdc.MustUnmarshal(cosmosValSetStore.Get(valAddress.Bytes()), &weightedAddress)
	weightedAddress.Weight = weight
	cosmosValSetStore.Set(valAddress.Bytes(), k.cdc.MustMarshal(&weightedAddress))
}

func (k Keeper) getAllCosmosValidatorSet(ctx sdk.Context) (weightedAddresses cosmosTypes.WeightedAddressAmounts) {
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

func (k Keeper) updateCurrentDelegatedAmountOfCosmosValidator(ctx sdk.Context, valAddress sdk.ValAddress, amount sdk.Coin) {
	cosmosValSetStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyCosmosValidatorWeights)
	if !cosmosValSetStore.Has(valAddress.Bytes()) {
		panic(fmt.Errorf("valAddress not present in kv store"))
	}
	var weightedAddress cosmosTypes.WeightedAddressAmount
	k.cdc.MustUnmarshal(cosmosValSetStore.Get(valAddress.Bytes()), &weightedAddress)
	weightedAddress.Amount = amount.Amount
	weightedAddress.Denom = amount.Denom
	cosmosValSetStore.Set(valAddress.Bytes(), k.cdc.MustMarshal(&weightedAddress))
}

func (k Keeper) updateCosmosValidatorStakingParams(ctx sdk.Context, msgs []sdk.Msg) error {
	uatomDenom, err := k.GetParams(ctx).GetBondDenomOf("uatom")
	if err != nil {
		return err
	}
	totalAmountInDelegateMsgs := sdk.NewInt64Coin(uatomDenom, 0)

	//TODO : MsgUndelegate
	msgsMap := make(map[string]stakingTypes.MsgDelegate, len(msgs))
	for _, msg := range msgs {
		delegateMsg := msg.(*stakingTypes.MsgDelegate)
		totalAmountInDelegateMsgs = totalAmountInDelegateMsgs.Add(delegateMsg.Amount)
		msgsMap[delegateMsg.ValidatorAddress] = *delegateMsg
	}

	k.setTotalDelegatedAmountTillDate(ctx, totalAmountInDelegateMsgs)

	internalWeightedAddressCosmos := k.getAllCosmosValidatorSet(ctx)
	for _, element := range internalWeightedAddressCosmos {
		if val, ok := msgsMap[element.Address]; ok {
			if element.Denom != val.Amount.Denom {
				continue
			}
			element.Amount.Add(val.Amount.Amount)
			valAddress, err := cosmosTypes.ValAddressFromBech32(element.Address, cosmosTypes.Bech32PrefixValAddr)
			if err != nil {
				panic(err)
			}
			k.updateCurrentDelegatedAmountOfCosmosValidator(ctx, valAddress, element.Coin())
		}
	}
	return nil
	//TODO : Update c token ratio
}

func (k Keeper) getCurrentDelegatedAmountOfCosmosValidator(ctx sdk.Context, valAddress sdk.ValAddress) sdk.Coin {
	cosmosValSetStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyCosmosValidatorWeights)
	if !cosmosValSetStore.Has(valAddress.Bytes()) {
		panic(fmt.Errorf("valAddress not present in kv store"))
	}
	var weightedAddress cosmosTypes.WeightedAddressAmount
	k.cdc.MustUnmarshal(cosmosValSetStore.Get(valAddress.Bytes()), &weightedAddress)
	return weightedAddress.Coin()
}

// removes the given validator from the cosmos validator set
func (k Keeper) removeCosmosValidatorWeight(ctx sdk.Context, valAddress sdk.ValAddress) {
	cosmosValSetStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyCosmosValidatorWeights)
	cosmosValSetStore.Delete(valAddress.Bytes())
}

//______________________________________________________________________________________________________________________

// put a check on length of validator set and val set weights to maintain equal mapping
// sets the native validator address as key and weight details as value
func (k Keeper) setOracleValidatorSet(ctx sdk.Context, valAddresses []sdk.ValAddress, nativeValSetWeights []cosmosTypes.WeightedAddress) {
	nativeValSetStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyNativeValidatorWeights)

	for i, va := range valAddresses {
		nativeValSetStore.Set(va.Bytes(), k.cdc.MustMarshal(&nativeValSetWeights[i]))
	}
}

// sets given weight details if store already has val address in it or panics
func (k Keeper) setOracleValidatorWeight(ctx sdk.Context, valAddress sdk.ValAddress, weight cosmosTypes.WeightedAddress) {
	nativeValSetStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyNativeValidatorWeights)
	if !nativeValSetStore.Has(valAddress.Bytes()) {
		panic(fmt.Errorf("valAddress not present in kv store"))
	}
	nativeValSetStore.Set(valAddress.Bytes(), k.cdc.MustMarshal(&weight))
}

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

// removes the given validator from the native validator set
func (k Keeper) removeOracleValidatorWeight(ctx sdk.Context, valAddress sdk.ValAddress) {
	nativeValSetStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyNativeValidatorWeights)
	nativeValSetStore.Delete(valAddress.Bytes())
}
