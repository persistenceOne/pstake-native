package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hashicorp/golang-lru/simplelru"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	"sync"
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
		panic(err)
	}
}

func (k Keeper) AddToMintedAmount(ctx sdk.Context, newlyMinted sdk.Coin) {
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

func (k Keeper) SubFromMintedAmount(ctx sdk.Context, burntAmount sdk.Coin) {
	cValueMu.Lock()
	defer cValueMu.Unlock()
	cValueCache.Remove(cValueCacheKey)

	store := ctx.KVStore(k.storeKey)

	var newMintedAmount sdk.Coin
	k.cdc.MustUnmarshal(store.Get(cosmosTypes.KeyMintedAmount), &newMintedAmount)
	newMintedAmount = newMintedAmount.Sub(burntAmount)
	store.Set(cosmosTypes.KeyMintedAmount, k.cdc.MustMarshal(&newMintedAmount))
}

func (k Keeper) AddToVirtuallyStakedAmount(ctx sdk.Context, notStakedAmount sdk.Coin) {
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

func (k Keeper) SubFromVirtuallyStakedAmount(ctx sdk.Context, notStakedAmount sdk.Coin) {
	cValueMu.Lock()
	defer cValueMu.Unlock()
	cValueCache.Remove(cValueCacheKey)

	store := ctx.KVStore(k.storeKey)

	var newVirtuallyStakedAmount sdk.Coin
	k.cdc.MustUnmarshal(store.Get(cosmosTypes.KeyVirtuallyStakedAmount), &newVirtuallyStakedAmount)
	newVirtuallyStakedAmount = newVirtuallyStakedAmount.Sub(notStakedAmount)
	store.Set(cosmosTypes.KeyVirtuallyStakedAmount, k.cdc.MustMarshal(&newVirtuallyStakedAmount))
}

func (k Keeper) AddToStakedAmount(ctx sdk.Context, stakedAmount sdk.Coin) {
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

func (k Keeper) SubFromStakedAmount(ctx sdk.Context, stakedAmount sdk.Coin) {
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
	newStakedAmount = newStakedAmount.Sub(stakedAmount)
	store.Set(cosmosTypes.KeyStakedAmount, k.cdc.MustMarshal(&newStakedAmount))
}

func (k Keeper) GetMintedAmount(ctx sdk.Context) sdk.Coin {
	store := ctx.KVStore(k.storeKey)
	if store.Has(cosmosTypes.KeyMintedAmount) {
		var mintedAmount sdk.Coin
		k.cdc.MustUnmarshal(store.Get(cosmosTypes.KeyMintedAmount), &mintedAmount)
		return mintedAmount
	}
	return sdk.NewInt64Coin(k.GetParams(ctx).MintDenom, 0)
}

func (k Keeper) GetVirtuallyStakedAmount(ctx sdk.Context) sdk.Coin {
	store := ctx.KVStore(k.storeKey)
	if store.Has(cosmosTypes.KeyVirtuallyStakedAmount) {
		var mintedAmount sdk.Coin
		k.cdc.MustUnmarshal(store.Get(cosmosTypes.KeyVirtuallyStakedAmount), &mintedAmount)
		return mintedAmount
	}
	bondDenom, err := k.GetParams(ctx).GetBondDenomOf("uatom")
	if err != nil {
		panic(err)
	}
	return sdk.NewInt64Coin(bondDenom, 0)
}

func (k Keeper) GetStakedAmount(ctx sdk.Context) sdk.Coin {
	store := ctx.KVStore(k.storeKey)
	if store.Has(cosmosTypes.KeyStakedAmount) {
		var mintedAmount sdk.Coin
		k.cdc.MustUnmarshal(store.Get(cosmosTypes.KeyStakedAmount), &mintedAmount)
		return mintedAmount
	}
	bondDenom, err := k.GetParams(ctx).GetBondDenomOf("uatom")
	if err != nil {
		panic(err)
	}
	return sdk.NewInt64Coin(bondDenom, 0)
}

func (k Keeper) GetCValue(ctx sdk.Context) sdk.Dec {
	cValueMu.Lock()
	defer cValueMu.Unlock()
	cValue, ok1 := cValueCache.Get(cValueCacheKey)
	if !ok1 {
		// calculate C value and set it and return
		totalStaked := k.GetVirtuallyStakedAmount(ctx).Amount.Add(k.GetStakedAmount(ctx).Amount)
		if totalStaked.IsZero() {
			cValueCache.Add(cValueCacheKey, CValue{cValue: sdk.NewDec(1), blockHeight: ctx.BlockHeight()})
			return sdk.NewDec(1)
		}
		calculatedCValue := sdk.NewDecFromInt(k.GetMintedAmount(ctx).Amount).Quo(sdk.NewDecFromInt(totalStaked))
		cValueCache.Add(cValueCacheKey, CValue{cValue: calculatedCValue, blockHeight: ctx.BlockHeight()})
		return calculatedCValue
	}

	cValueStruct := cValue.(CValue)
	// if the block has not changed then return the cached value
	if cValueStruct.blockHeight == ctx.BlockHeight() {
		return cValueStruct.cValue
	}

	// if the block has changed or struct is not properly converted then calculate new value
	totalStaked := k.GetVirtuallyStakedAmount(ctx).Amount.Add(k.GetStakedAmount(ctx).Amount)
	if totalStaked.IsZero() {
		cValueCache.Add(cValueCacheKey, CValue{cValue: sdk.NewDec(1), blockHeight: ctx.BlockHeight()})
		return sdk.NewDec(1)
	}
	calculatedCValue := sdk.NewDecFromInt(k.GetMintedAmount(ctx).Amount).Quo(sdk.NewDecFromInt(totalStaked))
	cValueCache.Add(cValueCacheKey, CValue{cValue: calculatedCValue, blockHeight: ctx.BlockHeight()})
	return calculatedCValue
}

func (k Keeper) SlashingEvent(ctx sdk.Context, slashedAmount sdk.Coin) {
	k.SubFromStakedAmount(ctx, slashedAmount)
	cValueMu.Lock()
	defer cValueMu.Unlock()
	cValueCache.Remove(cValueCacheKey)
}

//TODO : can add difference between actual value and our value and halt module on the basis of cut-off
