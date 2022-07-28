package keeper

import (
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hashicorp/golang-lru/simplelru"

	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

const (
	cValueCacheKey = "cValue"
)

type CValue struct {
	cValue      sdk.Dec
	blockHeight int64
}

var (
	cValueCache *simplelru.LRU
	cValueMu    sync.Mutex
)

func init() {
	var err error

	if cValueCache, err = simplelru.NewLRU(2, nil); err != nil {
		panic(any(err))
	}
}

// AddToMinted adds to the total minted amount
// used in case when tokens are minted
func (k Keeper) AddToMinted(ctx sdk.Context, newlyMinted sdk.Coin) {
	cValueMu.Lock()
	defer cValueMu.Unlock()
	cValueCache.Remove(cValueCacheKey)

	store := ctx.KVStore(k.storeKey)
	if !store.Has(cosmosTypes.KeyMintedAmount) {
		store.Set(cosmosTypes.KeyMintedAmount, k.cdc.MustMarshal(&newlyMinted))
		return
	}

	var newMintedAmount sdk.Coin
	k.cdc.MustUnmarshal(store.Get(cosmosTypes.KeyMintedAmount), &newMintedAmount)
	newMintedAmount = newMintedAmount.Add(newlyMinted)
	store.Set(cosmosTypes.KeyMintedAmount, k.cdc.MustMarshal(&newMintedAmount))
}

// SubFromMinted subtracts amount from total minted
// used in case when tokens are burnt
func (k Keeper) SubFromMinted(ctx sdk.Context, burntAmount sdk.Coin) {
	cValueMu.Lock()
	defer cValueMu.Unlock()
	cValueCache.Remove(cValueCacheKey)

	store := ctx.KVStore(k.storeKey)

	var newMintedAmount sdk.Coin
	k.cdc.MustUnmarshal(store.Get(cosmosTypes.KeyMintedAmount), &newMintedAmount)
	newMintedAmount = newMintedAmount.Sub(burntAmount)
	if newMintedAmount.IsNegative() {
		k.disableModule(ctx)
		panic(any("minted amount is negative"))
	}
	store.Set(cosmosTypes.KeyMintedAmount, k.cdc.MustMarshal(&newMintedAmount))
}

// GetMintedAmount gets minted amount
func (k Keeper) GetMintedAmount(ctx sdk.Context) sdk.Coin {
	store := ctx.KVStore(k.storeKey)
	if store.Has(cosmosTypes.KeyMintedAmount) {
		var mintedAmount sdk.Coin
		k.cdc.MustUnmarshal(store.Get(cosmosTypes.KeyMintedAmount), &mintedAmount)
		return mintedAmount
	}
	return sdk.NewInt64Coin(k.GetParams(ctx).MintDenom, 0)
}

//______________________________________________________________________________________________________________________

// AddToVirtuallyStaked adds to the total virtually staked amount
// used in case when the tokens have been minted but not yet staked
func (k Keeper) AddToVirtuallyStaked(ctx sdk.Context, notStakedAmount sdk.Coin) {
	cValueMu.Lock()
	defer cValueMu.Unlock()
	cValueCache.Remove(cValueCacheKey)

	store := ctx.KVStore(k.storeKey)
	if !store.Has(cosmosTypes.KeyVirtuallyStakedAmount) {
		store.Set(cosmosTypes.KeyVirtuallyStakedAmount, k.cdc.MustMarshal(&notStakedAmount))
		return
	}

	var newVirtuallyStakedAmount sdk.Coin
	k.cdc.MustUnmarshal(store.Get(cosmosTypes.KeyVirtuallyStakedAmount), &newVirtuallyStakedAmount)
	newVirtuallyStakedAmount = newVirtuallyStakedAmount.Add(notStakedAmount)
	store.Set(cosmosTypes.KeyVirtuallyStakedAmount, k.cdc.MustMarshal(&newVirtuallyStakedAmount))
}

// SubFromVirtuallyStaked subtracts from virtually staked amount
// used in case when then amount has been staked and being shifted to total staked amount
func (k Keeper) SubFromVirtuallyStaked(ctx sdk.Context, notStakedAmount sdk.Coin) {
	cValueMu.Lock()
	defer cValueMu.Unlock()
	cValueCache.Remove(cValueCacheKey)

	store := ctx.KVStore(k.storeKey)

	var newVirtuallyStakedAmount sdk.Coin
	k.cdc.MustUnmarshal(store.Get(cosmosTypes.KeyVirtuallyStakedAmount), &newVirtuallyStakedAmount)
	newVirtuallyStakedAmount = newVirtuallyStakedAmount.Sub(notStakedAmount)
	if newVirtuallyStakedAmount.IsNegative() {
		k.disableModule(ctx)
		panic(any("virtually staked amount is negative"))
	}
	store.Set(cosmosTypes.KeyVirtuallyStakedAmount, k.cdc.MustMarshal(&newVirtuallyStakedAmount))
}

// GetVirtuallyStakedAmount gets virtually staked amount
func (k Keeper) GetVirtuallyStakedAmount(ctx sdk.Context) sdk.Coin {
	store := ctx.KVStore(k.storeKey)
	if store.Has(cosmosTypes.KeyVirtuallyStakedAmount) {
		var mintedAmount sdk.Coin
		k.cdc.MustUnmarshal(store.Get(cosmosTypes.KeyVirtuallyStakedAmount), &mintedAmount)
		return mintedAmount
	}
	bondDenom, err := k.GetParams(ctx).GetBondDenomOf(cosmosTypes.DefaultStakingDenom)
	if err != nil {
		panic(any(err))
	}
	return sdk.NewInt64Coin(bondDenom, 0)
}

//______________________________________________________________________________________________________________________

// AddToStaked adds to total staked amount
// used in case when the tokens have been successfully staked
func (k Keeper) AddToStaked(ctx sdk.Context, stakedAmount sdk.Coin) {
	cValueMu.Lock()
	defer cValueMu.Unlock()
	cValueCache.Remove(cValueCacheKey)

	store := ctx.KVStore(k.storeKey)
	if !store.Has(cosmosTypes.KeyStakedAmount) {
		store.Set(cosmosTypes.KeyStakedAmount, k.cdc.MustMarshal(&stakedAmount))
		return
	}

	var newStakedAmount sdk.Coin
	k.cdc.MustUnmarshal(store.Get(cosmosTypes.KeyStakedAmount), &newStakedAmount)
	newStakedAmount = newStakedAmount.Add(stakedAmount)
	store.Set(cosmosTypes.KeyStakedAmount, k.cdc.MustMarshal(&newStakedAmount))
}

// SubFromStaked subtracts from the total staked amount
// used in case when the undelegate is successfully executed
func (k Keeper) SubFromStaked(ctx sdk.Context, stakedAmount sdk.Coin) {
	cValueMu.Lock()
	defer cValueMu.Unlock()
	cValueCache.Remove(cValueCacheKey)

	store := ctx.KVStore(k.storeKey)

	var newStakedAmount sdk.Coin
	k.cdc.MustUnmarshal(store.Get(cosmosTypes.KeyStakedAmount), &newStakedAmount)
	newStakedAmount = newStakedAmount.Sub(stakedAmount)
	if newStakedAmount.IsNegative() {
		k.disableModule(ctx)
		panic(any("staked amount is negative"))
	}
	store.Set(cosmosTypes.KeyStakedAmount, k.cdc.MustMarshal(&newStakedAmount))
}

// GetStakedAmount gets staked amount
func (k Keeper) GetStakedAmount(ctx sdk.Context) sdk.Coin {
	store := ctx.KVStore(k.storeKey)
	if store.Has(cosmosTypes.KeyStakedAmount) {
		var mintedAmount sdk.Coin
		k.cdc.MustUnmarshal(store.Get(cosmosTypes.KeyStakedAmount), &mintedAmount)
		return mintedAmount
	}
	bondDenom, err := k.GetParams(ctx).GetBondDenomOf(cosmosTypes.DefaultStakingDenom)
	if err != nil {
		panic(any(err))
	}
	return sdk.NewInt64Coin(bondDenom, 0)
}

//______________________________________________________________________________________________________________________

// AddToVirtuallyUnbonded adds to virtually unbonded amount
// used in case when the token have been withdrawn but not yet unbonded
func (k Keeper) AddToVirtuallyUnbonded(ctx sdk.Context, virtuallyUnbonded sdk.Coin) {
	cValueMu.Lock()
	defer cValueMu.Unlock()
	cValueCache.Remove(cValueCacheKey)

	store := ctx.KVStore(k.storeKey)
	if !store.Has(cosmosTypes.KeyVirtuallyUnbonded) {
		store.Set(cosmosTypes.KeyVirtuallyUnbonded, k.cdc.MustMarshal(&virtuallyUnbonded))
	}

	var newVirtuallyUnbonded sdk.Coin
	k.cdc.MustUnmarshal(store.Get(cosmosTypes.KeyVirtuallyUnbonded), &newVirtuallyUnbonded)
	newVirtuallyUnbonded = newVirtuallyUnbonded.Add(virtuallyUnbonded)
	store.Set(cosmosTypes.KeyVirtuallyUnbonded, k.cdc.MustMarshal(&newVirtuallyUnbonded))

}

// SubFromVirtuallyUnbonded subtracts from the total staked amount
// used in case when the unbond has been successfully executed and is being subtracted from total staked
func (k Keeper) SubFromVirtuallyUnbonded(ctx sdk.Context, virtuallyUnbonded sdk.Coin) {
	cValueMu.Lock()
	defer cValueMu.Unlock()
	cValueCache.Remove(cValueCacheKey)

	store := ctx.KVStore(k.storeKey)

	var newVirtuallyUnbonded sdk.Coin
	k.cdc.MustUnmarshal(store.Get(cosmosTypes.KeyVirtuallyUnbonded), &newVirtuallyUnbonded)
	newVirtuallyUnbonded = newVirtuallyUnbonded.Sub(virtuallyUnbonded)
	if newVirtuallyUnbonded.IsNegative() {
		k.disableModule(ctx)
		panic(any("virtual unbonded amount is negative"))
	}
	store.Set(cosmosTypes.KeyVirtuallyUnbonded, k.cdc.MustMarshal(&newVirtuallyUnbonded))
}

// GetVirtuallyUnbonded gets virtually unbonded amount
func (k Keeper) GetVirtuallyUnbonded(ctx sdk.Context) sdk.Coin {
	store := ctx.KVStore(k.storeKey)
	if store.Has(cosmosTypes.KeyVirtuallyUnbonded) {
		var virtuallyUnbonded sdk.Coin
		k.cdc.MustUnmarshal(store.Get(cosmosTypes.KeyVirtuallyUnbonded), &virtuallyUnbonded)
		return virtuallyUnbonded
	}
	bondDenom, err := k.GetParams(ctx).GetBondDenomOf(cosmosTypes.DefaultStakingDenom)
	if err != nil {
		panic(any(err))
	}
	return sdk.NewInt64Coin(bondDenom, 0)
}

//______________________________________________________________________________________________________________________

// GetCValue gets the C cached C value if cache is valid or re-calculates if expired
// returns 1 in case where total staked amount is 0
func (k Keeper) GetCValue(ctx sdk.Context) sdk.Dec {
	cValueMu.Lock()
	defer cValueMu.Unlock()
	cValue, ok1 := cValueCache.Get(cValueCacheKey)
	if !ok1 {
		// calculate C value and set it and return
		totalStaked := k.GetVirtuallyStakedAmount(ctx).Amount.Add(k.GetStakedAmount(ctx).Amount).Sub(k.GetVirtuallyUnbonded(ctx).Amount)
		if totalStaked.IsZero() {
			cValueCache.Add(cValueCacheKey, CValue{cValue: sdk.NewDec(1), blockHeight: ctx.BlockHeight()})
			return sdk.NewDec(1)
		}
		calculatedCValue := sdk.NewDecFromInt(k.GetMintedAmount(ctx).Amount).Quo(sdk.NewDecFromInt(totalStaked))
		if calculatedCValue.IsNegative() {
			k.disableModule(ctx)
		}
		cValueCache.Add(cValueCacheKey, CValue{cValue: calculatedCValue, blockHeight: ctx.BlockHeight()})
		return calculatedCValue
	}

	cValueStruct := cValue.(CValue)
	// if the block has not changed then return the cached value
	if cValueStruct.blockHeight == ctx.BlockHeight() {
		return cValueStruct.cValue
	}

	// if the block has changed or struct is not properly converted then calculate new value
	totalStaked := k.GetVirtuallyStakedAmount(ctx).Amount.Add(k.GetStakedAmount(ctx).Amount).Sub(k.GetVirtuallyUnbonded(ctx).Amount)
	if totalStaked.IsZero() {
		cValueCache.Add(cValueCacheKey, CValue{cValue: sdk.NewDec(1), blockHeight: ctx.BlockHeight()})
		return sdk.NewDec(1)
	}
	calculatedCValue := sdk.NewDecFromInt(k.GetMintedAmount(ctx).Amount).Quo(sdk.NewDecFromInt(totalStaked))
	if calculatedCValue.IsNegative() {
		k.disableModule(ctx)
	}
	cValueCache.Add(cValueCacheKey, CValue{cValue: calculatedCValue, blockHeight: ctx.BlockHeight()})
	return calculatedCValue
}

//______________________________________________________________________________________________________________________

// SlashingEvent resets the C value and subtracts the staked amount with the given slashed amount
func (k Keeper) SlashingEvent(ctx sdk.Context, slashedAmount sdk.Coin) {
	k.SubFromStaked(ctx, slashedAmount)
	cValueMu.Lock()
	defer cValueMu.Unlock()
	cValueCache.Remove(cValueCacheKey)
}

//TODO : can add difference between actual value and our value and halt module on the basis of cut-off
