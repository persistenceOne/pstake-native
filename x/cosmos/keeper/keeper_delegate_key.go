package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"

	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

/*
SetValidatorOrchestrator sets the oracle address for the given validator address.
Address corresponding to one validator can be multiple given that it is not already mapped to another validator
*/
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

	//check if orchestrator address is already mapped to another validator
	_, exist, err := k.getAllValidatorOrchestratorMappingAndFindIfExist(ctx, orch)
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("orchestrator address already exist")
	}

	key := val.Bytes()
	//check if validator key already exists or not
	if orchestratorValidatorStore.Has(key) {
		var valStoreValue cosmosTypes.ValidatorStoreValue
		k.cdc.MustUnmarshal(orchestratorValidatorStore.Get(key), &valStoreValue)

		// if only one address is present then new one is added to the 2nd position
		valStoreValue.OrchestratorAddresses = append(valStoreValue.OrchestratorAddresses, orch.String())
		orchestratorValidatorStore.Set(key, k.cdc.MustMarshal(&valStoreValue))
		return nil
	}

	// if no address is mapped with the validator
	a := cosmosTypes.NewValidatorStoreValue(orch)
	orchestratorValidatorStore.Set(key, k.cdc.MustMarshal(&a))
	return nil
}

/*
RemoveValidatorOrchestrator removes the oracle address mapped to the given validator.
Addresses mapped to one validator can be multiple, but not all can be removed. Only the ones that are not present in
multisig account address can be removed
*/
func (k Keeper) RemoveValidatorOrchestrator(ctx sdkTypes.Context, val sdkTypes.ValAddress, orch sdkTypes.AccAddress) error {
	if err := sdkTypes.VerifyAddressFormat(val); err != nil {
		return sdkErrors.Wrap(err, "invalid val address")
	}
	if err := sdkTypes.VerifyAddressFormat(orch); err != nil {
		return sdkErrors.Wrap(err, "invalid orch address")
	}

	//checks if orch address is present in current multisig or not
	present := k.checkOrchestratorAddressPresentInMultisig(ctx, orch)
	if present {
		return fmt.Errorf("orch address present in multisig address")
	}

	orchestratorValidatorStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.ValidatorOrchestratorStoreKey)

	//checks if validator already exist in staking
	_, found := k.stakingKeeper.GetValidator(ctx, val)
	if !found {
		return fmt.Errorf("validator address : %s does not exist on chain", val.String())
	}

	//check if orchestrator address is already mapped to another validator
	_, exist, err := k.getAllValidatorOrchestratorMappingAndFindIfExist(ctx, orch)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("orchestrator address : %s does not exist", orch.String())
	}

	key := val.Bytes()
	if orchestratorValidatorStore.Has(key) {
		var valStoreValue cosmosTypes.ValidatorStoreValue
		k.cdc.MustUnmarshal(orchestratorValidatorStore.Get(key), &valStoreValue)

		if len(valStoreValue.OrchestratorAddresses) == 1 {
			return fmt.Errorf("can not remove the one and only mapping")
		}

		// find the element with address same as the orch address and remove that index from orch address array
		for index, vs := range valStoreValue.OrchestratorAddresses {
			if vs == orch.String() {
				valStoreValue.OrchestratorAddresses = RemoveIndex(valStoreValue.OrchestratorAddresses, index)
			}
		}

		orchestratorValidatorStore.Set(key, k.cdc.MustMarshal(&valStoreValue))
		return nil
	}

	return fmt.Errorf("validator address not present in kv store")
}

// CheckValidator checks if the validator exists on chain or not. Returns address and found bool.
func (k Keeper) CheckValidator(ctx sdkTypes.Context, val sdkTypes.ValAddress) (validator sdkTypes.ValAddress, found bool) {
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

// GetTotalValidatorOrchestratorCount gets the count of total validator and orchestrator mappings for ratio calculation
func (k Keeper) GetTotalValidatorOrchestratorCount(ctx sdkTypes.Context) int64 {
	orchestratorValidatorStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.ValidatorOrchestratorStoreKey)
	iterator := orchestratorValidatorStore.Iterator(nil, nil)
	defer iterator.Close()
	counter := 0
	for ; iterator.Valid(); iterator.Next() {
		counter++
	}
	return int64(counter)
}

// gets the validator mapping currently present in DB
func (k Keeper) getValidatorMapping(ctx sdkTypes.Context, valAddress sdkTypes.ValAddress) cosmosTypes.ValidatorStoreValue {
	orchestratorValidatorStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.ValidatorOrchestratorStoreKey)
	iterator := orchestratorValidatorStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		val := sdkTypes.ValAddress(iterator.Key())
		if val.Equals(valAddress) {
			var validatorStoreValue cosmosTypes.ValidatorStoreValue
			k.cdc.MustUnmarshal(iterator.Value(), &validatorStoreValue)
			return validatorStoreValue
		}
	}
	return cosmosTypes.ValidatorStoreValue{}
}

// gets all validator mapping  and checks if the given oracle address is present in it.
func (k Keeper) getAllValidatorOrchestratorMappingAndFindIfExist(ctx sdkTypes.Context,
	orch sdkTypes.AccAddress) (valAddress sdkTypes.ValAddress, found bool, err error) {
	found = false
	orchestratorValidatorStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.ValidatorOrchestratorStoreKey)
	iterator := orchestratorValidatorStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var validatorStoreValue cosmosTypes.ValidatorStoreValue
		err = k.cdc.Unmarshal(iterator.Value(), &validatorStoreValue)
		if err != nil {
			return nil, found, err
		}
		for _, address := range validatorStoreValue.OrchestratorAddresses {
			if address == orch.String() {
				val := sdkTypes.ValAddress(iterator.Key())
				valAddress = val
				found = true
			}
		}
	}
	if valAddress == nil {
		err = fmt.Errorf("validator address not found")
	}
	return valAddress, found, err
}

// checks if all the validators have oracle mapping. Used in HandleEnableModuleProposal.
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
		if count > 2 {
			return []string{}, cosmosTypes.ErrMoreThanTwoOrchestratorAddressesMapping
		}
		validatorOrchestratorMap[val.String()] = validatorStoreValue.OrchestratorAddresses
	}

	for _, val := range k.getAllOracleValidatorSet(ctx) {
		if validatorOrchestratorMap[val.Address] == nil {
			return []string{}, fmt.Errorf("validator mapping not present in KV store")
		}
		delete(validatorOrchestratorMap, val.Address)
	}

	if len(validatorOrchestratorMap) > 0 {
		return []string{}, fmt.Errorf("more than expected valdiator orchestrator mapping present")
	}

	return orchestratorList, nil
}

// RemoveIndex Remove the given index from the given slice of string
func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}
