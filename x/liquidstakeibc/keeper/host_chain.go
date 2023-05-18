package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

// SetHostChain sets a host chain in the store
func (k *Keeper) SetHostChain(ctx sdk.Context, hc *types.HostChain) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HostChainKey)
	bytes := k.cdc.MustMarshal(hc)
	store.Set([]byte(hc.ChainId), bytes)
}

// SetHostChainValidator sets a validator on the target host chain
func (k *Keeper) SetHostChainValidator(
	ctx sdk.Context,
	hc *types.HostChain,
	validator *types.Validator,
) {
	found := false
	for i, val := range hc.Validators {
		if validator.OperatorAddress == val.OperatorAddress {
			hc.Validators[i] = validator
			found = true
			break
		}
	}

	if !found {
		hc.Validators = append(hc.Validators, validator)
	}

	k.SetHostChain(ctx, hc)
}

// ProcessHostChainValidatorUpdates processes the new validator set for a host chain
func (k *Keeper) ProcessHostChainValidatorUpdates(
	ctx sdk.Context,
	hc *types.HostChain,
	validators []stakingtypes.Validator,
) error {
	for _, validator := range validators {
		val, found := hc.GetValidator(validator.OperatorAddress)

		if !found {
			hc.Validators = append(
				hc.Validators,
				&types.Validator{
					OperatorAddress: validator.OperatorAddress,
					Status:          validator.Status.String(),
					Weight:          sdk.ZeroDec(),
					DelegatedAmount: sdk.ZeroInt(),
					TotalAmount:     validator.Tokens,
				},
			)
			k.SetHostChain(ctx, hc)
		} else {
			if validator.Status.String() != val.Status {
				val.Status = validator.Status.String()
				if val.Status == stakingtypes.BondStatusUnbonding && validator.Jailed {
					// TODO: Undelegate ? What if it is Jailed
				}
				k.SetHostChainValidator(ctx, hc, val)
			}
			if !validator.Tokens.Equal(val.TotalAmount) {
				// validator has been slashed, update its delegation
				if err := k.QueryValidatorDelegation(ctx, hc, val); err != nil {
					return fmt.Errorf(
						"error while querying validator %s delegation: %s",
						val.OperatorAddress,
						err.Error(),
					)
				}
				k.SetHostChainValidator(ctx, hc, val)
			}
		}
	}

	return nil
}

// GetHostChain returns a host chain given its id
func (k *Keeper) GetHostChain(ctx sdk.Context, chainID string) (*types.HostChain, bool) {
	hc := types.HostChain{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HostChainKey)
	bytes := store.Get([]byte(chainID))
	if len(bytes) == 0 {
		return &hc, false
	}

	k.cdc.MustUnmarshal(bytes, &hc)
	return &hc, true
}

// GetAllHostChains retrieves all registered host chains
func (k *Keeper) GetAllHostChains(ctx sdk.Context) []*types.HostChain {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HostChainKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	hostChains := make([]*types.HostChain, 0)
	for ; iterator.Valid(); iterator.Next() {
		hc := types.HostChain{}
		k.cdc.MustUnmarshal(iterator.Value(), &hc)
		hostChains = append(hostChains, &hc)
	}

	return hostChains
}

// GetHostChainFromIbcDenom returns a host chain given its ibc denomination on Persistence
func (k *Keeper) GetHostChainFromIbcDenom(ctx sdk.Context, ibcDenom string) (*types.HostChain, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HostChainKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	found := false
	hc := types.HostChain{}
	for ; iterator.Valid(); iterator.Next() {
		chain := types.HostChain{}
		k.cdc.MustUnmarshal(iterator.Value(), &chain)

		if chain.IBCDenom() == ibcDenom {
			hc = chain
			found = true
			break
		}
	}

	return &hc, found
}

// GetHostChainFromHostDenom returns a host chain given its host denomination
func (k *Keeper) GetHostChainFromHostDenom(ctx sdk.Context, hostDenom string) (*types.HostChain, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HostChainKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	found := false
	hc := types.HostChain{}
	for ; iterator.Valid(); iterator.Next() {
		chain := types.HostChain{}
		k.cdc.MustUnmarshal(iterator.Value(), &chain)

		if chain.HostDenom == hostDenom {
			hc = chain
			found = true
			break
		}
	}

	return &hc, found
}

// GetHostChainFromDelegatorAddress returns a host chain given its delegator address
func (k *Keeper) GetHostChainFromDelegatorAddress(ctx sdk.Context, delegatorAddress string) (*types.HostChain, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HostChainKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	found := false
	hc := types.HostChain{}
	for ; iterator.Valid(); iterator.Next() {
		chain := types.HostChain{}
		k.cdc.MustUnmarshal(iterator.Value(), &chain)

		if chain.DelegationAccount != nil && chain.DelegationAccount.Address == delegatorAddress {
			hc = chain
			found = true
			break
		}
	}

	return &hc, found
}

// GetHostChainCValue calculates the host chain c value
func (k *Keeper) GetHostChainCValue(ctx sdk.Context, hc *types.HostChain) sdk.Dec {
	// total stk minted amount
	mintedAmount := k.bankKeeper.GetSupply(ctx, hc.MintDenom()).Amount

	// delegated amount + delegation account balance + deposit module account balance
	liquidStakedAmount := hc.GetHostChainTotalDelegations().
		Add(hc.DelegationAccount.Balance.Amount).
		Add(k.bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(types.DepositModuleAccount), hc.IBCDenom()).Amount)

	if mintedAmount.IsZero() || liquidStakedAmount.IsZero() {
		return sdk.OneDec()
	}

	return sdk.NewDecFromInt(mintedAmount).Quo(sdk.NewDecFromInt(liquidStakedAmount))
}

// UpdateHostChainValidatorWeight updates a host chain validator weight
func (k *Keeper) UpdateHostChainValidatorWeight(
	ctx sdk.Context,
	hc *types.HostChain,
	address string,
	weight string,
) error {
	newWeight, err := sdk.NewDecFromStr(weight)
	if err != nil {
		return err
	}

	found := false
	for i, validator := range hc.Validators {
		if validator.OperatorAddress == address {
			hc.Validators[i].Weight = newWeight
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("could not find validator with address %s while updating validator weight", address)
	}

	k.SetHostChain(ctx, hc)
	return nil
}
