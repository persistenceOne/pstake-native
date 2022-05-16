package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

// SetValidatorOrchestrator sets the Orchestrator key for a given validator
func (k Keeper) SetValidatorOrchestrator(ctx sdkTypes.Context, val sdkTypes.ValAddress, orch sdkTypes.AccAddress) error {
	if err := sdkTypes.VerifyAddressFormat(val); err != nil {
		panic(sdkErrors.Wrap(err, "invalid val address"))
	}
	if err := sdkTypes.VerifyAddressFormat(orch); err != nil {
		panic(sdkErrors.Wrap(err, "invalid orch address"))
	}
	orchestratorValidatorStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.ValidatorOrchestratorStoreKey)

	//checks if validator already exist in staking
	_, found := k.stakingKeeper.GetValidator(ctx, val)
	if !found {
		return fmt.Errorf("validator address does not exist")
	}

	//make key out of it
	key := val.Bytes()

	//check if validator key already exists or not
	if orchestratorValidatorStore.Has(key) {
		return fmt.Errorf("validator orchestrator mapping already presnet")
	}

	//check if orchestrator address is already mapped to another validator
	_, _, exist, err := k.getAllValidatorOrchestratorMappingAndFindIfExist(ctx, orch)
	if err != nil {
		return err
	}
	if exist == true {
		return fmt.Errorf("orchestrator address already exist")
	}

	//set in store
	a := cosmosTypes.NewValidatorStoreValue(orch)
	bz, err := k.cdc.Marshal(&a)
	if err != nil {
		return err
	}
	orchestratorValidatorStore.Set(key, bz)
	return nil
}

func (k Keeper) GetValidatorOrchestrator(ctx sdkTypes.Context, val sdkTypes.ValAddress) (validator sdkTypes.ValAddress, found bool) {
	if err := sdkTypes.VerifyAddressFormat(val); err != nil {
		ctx.Logger().Error("invalid orch address")
		return validator, false
	}

	_, found = k.stakingKeeper.GetValidator(ctx, val)
	if found {
		return val, found
	}

	return nil, found
}

// gets the count of total validator and orchestrator mappings for ratio calculation
func (k Keeper) getTotalValidatorOrchestratorCount(ctx sdkTypes.Context) int64 {
	orchestratorValidatorStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.ValidatorOrchestratorStoreKey)
	iterator := orchestratorValidatorStore.Iterator(nil, nil)
	defer iterator.Close()
	counter := 0
	for ; iterator.Valid(); iterator.Next() {
		counter++
	}
	return int64(counter)
}

func (k Keeper) getAllValidatorOrchestratorMappingAndFindIfExist(ctx sdkTypes.Context, orch sdkTypes.AccAddress) (orchAddresses []string, valAddress sdkTypes.ValAddress, found bool, err error) {
	found = false
	orchestratorValidatorStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.ValidatorOrchestratorStoreKey)
	iterator := orchestratorValidatorStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var validatorStoreValue cosmosTypes.ValidatorStoreValue
		err = k.cdc.Unmarshal(iterator.Value(), &validatorStoreValue)
		if err != nil {
			return orchAddresses, nil, found, err
		}
		for _, address := range validatorStoreValue.OrchestratorAddresses {
			if address == orch.String() {
				val := sdkTypes.ValAddress(iterator.Key())
				valAddress = val
				found = true
			}
			orchAddresses = append(orchAddresses, address)
		}
	}
	return orchAddresses, valAddress, found, err
}

func (k Keeper) checkAllValidatorsHaveOrchestrators(ctx sdkTypes.Context) ([]string, error) {
	orchestratorValidatorStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.ValidatorOrchestratorStoreKey)
	validatorOrchestratorMap := make(map[string][]string)
	var orchestratorList []string
	iterator := orchestratorValidatorStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		val := sdkTypes.ValAddress(iterator.Key())
		var validatorStoreValue cosmosTypes.ValidatorStoreValue
		err := k.cdc.Unmarshal(iterator.Value(), &validatorStoreValue)
		if err != nil {
			return []string{}, err
		}
		count := 0
		for _, address := range validatorStoreValue.OrchestratorAddresses {
			if address != "" {
				orchestratorList = append(orchestratorList, address)
				count++
			}
		}
		if count == 0 {
			return []string{}, cosmosTypes.ErrValidatorOrchestratorMappingNotFound
		}
		if count > 1 {
			return []string{}, cosmosTypes.ErrMoreThanOneMapping
		}
		validatorOrchestratorMap[val.String()] = validatorStoreValue.OrchestratorAddresses
	}

	for _, val := range k.GetParams(ctx).ValidatorSetNativeChain {
		if validatorOrchestratorMap[val.Address] == nil {
			return []string{}, fmt.Errorf("validator mapping not present in KV store")
		}
		delete(validatorOrchestratorMap, val.Address)
	}

	for key := range validatorOrchestratorMap {
		if len(key) > 0 {
			return []string{}, fmt.Errorf("more than expected valdiator orchestrator mapping present")
		}
	}

	return orchestratorList, nil
}
