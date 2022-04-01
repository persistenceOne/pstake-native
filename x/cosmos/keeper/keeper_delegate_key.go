package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
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

func (k Keeper) GetOrchestratorValidator(ctx sdkTypes.Context, orch sdkTypes.AccAddress) (validator sdkTypes.ValAddress, found bool) {
	if err := sdkTypes.VerifyAddressFormat(orch); err != nil {
		ctx.Logger().Error("invalid orch address")
		return validator, false
	}
	store := ctx.KVStore(k.storeKey)
	orchestratorValidatorStore := prefix.NewStore(store, []byte(cosmosTypes.OrchestratorValidatorStoreKey))
	valAddr := orchestratorValidatorStore.Get([]byte(cosmosTypes.GetOrchestratorAddressKey(orch)))
	validatorDetails, found := k.stakingKeeper.GetValidator(ctx, valAddr)
	validator = sdkTypes.ValAddress(validatorDetails.OperatorAddress)

	return validator, found
}

// gets the count of total validator and orchestrator mappings for ratio calculation
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
