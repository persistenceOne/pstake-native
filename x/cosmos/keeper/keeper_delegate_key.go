package keeper

import (
	"time"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

// SetOrchestratorValidator sets the Orchestrator key for a given validator
func (k Keeper) SetOrchestratorValidator(ctx sdkTypes.Context, val sdkTypes.ValAddress, orch sdkTypes.AccAddress) {
	if err := sdkTypes.VerifyAddressFormat(val); err != nil {
		panic(sdkErrors.Wrap(err, "invalid val address"))
	}
	if err := sdkTypes.VerifyAddressFormat(orch); err != nil {
		panic(sdkErrors.Wrap(err, "invalid orch address"))
	}
	store := ctx.KVStore(k.storeKey)
	orchestratorValidatorStore := prefix.NewStore(store, []byte(cosmosTypes.OrchestratorValidatorStoreKey))
	orchestratorValidatorStore.Set([]byte(cosmosTypes.GetOrchestratorAddressKey(orch)), val.Bytes())
}

func (k Keeper) GetOrchestratorValidator(ctx sdkTypes.Context, orch sdkTypes.AccAddress) (validator stakingTypes.Validator, found bool) {
	if err := sdkTypes.VerifyAddressFormat(orch); err != nil {
		ctx.Logger().Error("invalid orch address")
		return validator, false
	}
	store := ctx.KVStore(k.storeKey)
	orchestratorValidatorStore := prefix.NewStore(store, []byte(cosmosTypes.OrchestratorValidatorStoreKey))
	valAddr := orchestratorValidatorStore.Get([]byte(cosmosTypes.GetOrchestratorAddressKey(orch)))
	if valAddr == nil {
		return stakingTypes.Validator{
			OperatorAddress: "",
			ConsensusPubkey: &codecTypes.Any{
				TypeUrl:              "",
				Value:                []byte{},
				XXX_NoUnkeyedLiteral: struct{}{},
				XXX_unrecognized:     []byte{},
				XXX_sizecache:        0,
			},
			Jailed:          false,
			Status:          0,
			Tokens:          sdkTypes.Int{},
			DelegatorShares: sdkTypes.Dec{},
			Description: stakingTypes.Description{
				Moniker:         "",
				Identity:        "",
				Website:         "",
				SecurityContact: "",
				Details:         "",
			},
			UnbondingHeight: 0,
			UnbondingTime:   time.Time{},
			Commission: stakingTypes.Commission{
				CommissionRates: stakingTypes.CommissionRates{
					Rate:          sdkTypes.Dec{},
					MaxRate:       sdkTypes.Dec{},
					MaxChangeRate: sdkTypes.Dec{},
				},
				UpdateTime: time.Time{},
			},
			MinSelfDelegation: sdkTypes.Int{},
		}, false
	}
	validator, found = k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return stakingTypes.Validator{
			OperatorAddress: "",
			ConsensusPubkey: &codecTypes.Any{
				TypeUrl:              "",
				Value:                []byte{},
				XXX_NoUnkeyedLiteral: struct{}{},
				XXX_unrecognized:     []byte{},
				XXX_sizecache:        0,
			},
			Jailed:          false,
			Status:          0,
			Tokens:          sdkTypes.Int{},
			DelegatorShares: sdkTypes.Dec{},
			Description: stakingTypes.Description{
				Moniker:         "",
				Identity:        "",
				Website:         "",
				SecurityContact: "",
				Details:         "",
			},
			UnbondingHeight: 0,
			UnbondingTime:   time.Time{},
			Commission: stakingTypes.Commission{
				CommissionRates: stakingTypes.CommissionRates{
					Rate:          sdkTypes.Dec{},
					MaxRate:       sdkTypes.Dec{},
					MaxChangeRate: sdkTypes.Dec{},
				},
				UpdateTime: time.Time{},
			},
			MinSelfDelegation: sdkTypes.Int{},
		}, false
	}

	return validator, true
}

func (k Keeper) getTotalValidatorOrchestratorCount(ctx sdkTypes.Context) int64 {
	store := ctx.KVStore(k.storeKey)
	orchestratorValidatorStore := prefix.NewStore(store, []byte(cosmosTypes.OrchestratorValidatorStoreKey))
	iterator := orchestratorValidatorStore.Iterator(nil, nil)
	defer iterator.Close()
	counter := 0
	for ; iterator.Valid(); iterator.Next() {
		counter++
	}
	return int64(counter)
}
