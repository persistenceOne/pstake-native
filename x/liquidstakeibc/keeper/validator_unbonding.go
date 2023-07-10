package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func (k *Keeper) SetValidatorUnbonding(ctx sdk.Context, vu *types.ValidatorUnbonding) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ValidatorUnbondingKey)
	bytes := k.cdc.MustMarshal(vu)
	store.Set(types.GetValidatorUnbondingStoreKey(vu.ChainId, vu.ValidatorAddress, vu.EpochNumber), bytes)
}

func (k *Keeper) GetValidatorUnbonding(
	ctx sdk.Context,
	chainID string,
	validatorAddress string,
	epochNumber int64,
) (*types.ValidatorUnbonding, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ValidatorUnbondingKey)
	bz := store.Get(types.GetValidatorUnbondingStoreKey(chainID, validatorAddress, epochNumber))
	if bz == nil {
		return &types.ValidatorUnbonding{}, false
	}

	var validatorUnbonding types.ValidatorUnbonding
	k.cdc.MustUnmarshal(bz, &validatorUnbonding)
	return &validatorUnbonding, true
}

func (k *Keeper) GetAllValidatorUnbondedAmount(ctx sdk.Context, hc *types.HostChain) sdk.Int { //nolint:staticcheck
	validatorUnbondings := k.FilterValidatorUnbondings(
		ctx,
		func(u types.ValidatorUnbonding) bool {
			return u.ChainId == hc.ChainId && u.MatureTime != time.Time{}
		},
	)

	amount := sdk.ZeroInt()
	for _, validatorUnbonding := range validatorUnbondings {
		amount = amount.Add(validatorUnbonding.Amount.Amount)
	}

	return amount
}

func (k *Keeper) DeleteValidatorUnbonding(ctx sdk.Context, ub *types.ValidatorUnbonding) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ValidatorUnbondingKey)
	store.Delete(types.GetValidatorUnbondingStoreKey(ub.ChainId, ub.ValidatorAddress, ub.EpochNumber))

	telemetry.IncrCounter(float32(-1), ub.ChainId, "validator_unbondings")
}

func (k *Keeper) DeleteValidatorUnbondingsForSequenceID(ctx sdk.Context, sequenceID string) {
	validatorUnbondings := k.FilterValidatorUnbondings(
		ctx,
		func(u types.ValidatorUnbonding) bool {
			return u.IbcSequenceId == sequenceID
		},
	)

	for _, validatorUnbonding := range validatorUnbondings {
		k.DeleteValidatorUnbonding(ctx, validatorUnbonding)
	}
}

func (k *Keeper) FilterValidatorUnbondings(
	ctx sdk.Context,
	filter func(u types.ValidatorUnbonding) bool,
) []*types.ValidatorUnbonding {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ValidatorUnbondingKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	validatorUnbondings := make([]*types.ValidatorUnbonding, 0)
	for ; iterator.Valid(); iterator.Next() {
		validatorUnbonding := types.ValidatorUnbonding{}
		k.cdc.MustUnmarshal(iterator.Value(), &validatorUnbonding)
		if filter(validatorUnbonding) {
			validatorUnbondings = append(validatorUnbondings, &validatorUnbonding)
		}
	}

	return validatorUnbondings
}
