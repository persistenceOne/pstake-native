package keeper

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	validator stakingtypes.Validator,
) error {
	val, found := hc.GetValidator(validator.OperatorAddress)
	if !found {
		return fmt.Errorf("validator with address %s not registered", validator.OperatorAddress)
	}

	// process status update
	if validator.Status.String() != val.Status {
		// validator transitioned into unbonding
		if validator.Status.String() != stakingtypes.BondStatusBonded {
			epochNumber := k.epochsKeeper.GetEpochInfo(ctx, types.UndelegationEpoch).CurrentEpoch
			val.UnbondingEpoch = types.CurrentUnbondingEpoch(hc.UnbondingFactor, epochNumber)
		}
		// validator transitioned into bonded
		if validator.Status.String() == stakingtypes.BondStatusBonded {
			val.UnbondingEpoch = 0
		}

		// emit the status update event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeValidatorStatusUpdate,
				sdk.NewAttribute(types.AttributeChainID, hc.ChainId),
				sdk.NewAttribute(types.AttributeValidatorAddress, val.OperatorAddress),
				sdk.NewAttribute(types.AttributeKeyValidatorOldStatus, val.Status),
				sdk.NewAttribute(types.AttributeKeyValidatorNewStatus, validator.Status.String()),
			),
		)

		val.Status = validator.Status.String()
		k.SetHostChainValidator(ctx, hc, val)
	}

	// process exchange rate update
	var exchangeRate sdk.Dec
	if validator.DelegatorShares.IsZero() {
		exchangeRate = sdk.OneDec()
	} else {
		exchangeRate = sdk.NewDecFromInt(validator.Tokens).Quo(validator.DelegatorShares)
	}
	if !exchangeRate.Equal(val.ExchangeRate) {
		if val.DelegatedAmount.GT(sdk.ZeroInt()) {
			if err := k.QueryValidatorDelegation(ctx, hc, val); err != nil {
				return fmt.Errorf(
					"error while querying validator %s delegation: %s",
					val.OperatorAddress,
					err.Error(),
				)
			}
		}

		// emit the exchange rate event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeValidatorExchangeRateUpdate,
				sdk.NewAttribute(types.AttributeChainID, hc.ChainId),
				sdk.NewAttribute(types.AttributeValidatorAddress, val.OperatorAddress),
				sdk.NewAttribute(types.AttributeKeyValidatorOldExchangeRate, val.ExchangeRate.String()),
				sdk.NewAttribute(types.AttributeKeyValidatorNewExchangeRate, exchangeRate.String()),
			),
		)

		val.ExchangeRate = exchangeRate
		k.SetHostChainValidator(ctx, hc, val)
	}

	// process LSM cap updates
	if hc.Flags.Lsm {
		// check if the validator has reached the LSM validator bond
		var validatorHasRoomForDelegations bool
		if validator.DelegatorShares.IsZero() {
			validatorHasRoomForDelegations = true // if no shares are issued yet, the validator can accept more delegations
		} else {
			validatorHasRoomForDelegations = validator.LiquidShares.Quo(validator.DelegatorShares).LT(hc.Params.LsmValidatorCap)
		}

		// check if the validator has reached the bonded shares cap
		var validatorHasEnoughBond bool
		// this is the default value for the bond factor, which disables the functionality
		// https://github.com/cosmos/cosmos-sdk/blob/0af2f4da004cbea6414a8bad56e8bdcd45badf1e/x/staking/types/params.go#L36-L73
		if hc.Params.LsmBondFactor.Equal(sdk.NewDecFromInt(sdk.NewInt(-1))) {
			validatorHasEnoughBond = true
		} else {
			validatorHasEnoughBond = validator.LiquidShares.LT(validator.ValidatorBondShares.Mul(hc.Params.LsmBondFactor))
		}

		// save the old delegable flag for event purposes
		oldDelegableFlag := val.Delegable

		// update the validator if its delegable status has changed
		val.Delegable = validatorHasRoomForDelegations && validatorHasEnoughBond
		k.SetHostChainValidator(ctx, hc, val)

		// this part of the code checks whether there is actually room to delegate on the validator.
		// it can happen that a validator will have not reached any of the caps but the amount that
		// pStake wants to stake is higher than the room available until reaching the cap.
		// if that happens, delegations will keep failing and start the loop that we wanted to avoid.
		// in order to avoid this, we simulate the generation of delegation messages to calculate the
		// exact amount that would be delegated to each validator, and then compare it to the room left on that
		// validator both on the validator cap and on the bond factor side.

		// if the validator is not delegable, no messages will be generated for it, so there is nothing else to do
		if !val.Delegable {
			// emit the delegable status event
			if oldDelegableFlag != val.Delegable {
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						types.EventTypeValidatorDelegableStateUpdate,
						sdk.NewAttribute(types.AttributeChainID, hc.ChainId),
						sdk.NewAttribute(types.AttributeValidatorAddress, val.OperatorAddress),
						sdk.NewAttribute(types.AttributeKeyValidatorDelegable, strconv.FormatBool(val.Delegable)),
					),
				)
			}

			return nil
		}

		// we use both deposits waiting to be delegated and deposits that are being delegated as the delegating workflow
		// is run every block, and there could be situations where not all deposits are being taken into account.
		delegableDeposits := append(k.GetDelegableDepositsForChain(ctx, hc.ChainId), k.GetDelegatingDepositsForChain(ctx, hc.ChainId)...)
		totalDepositDelegation := sdk.ZeroInt()
		for _, deposit := range delegableDeposits {
			totalDepositDelegation = totalDepositDelegation.Add(deposit.Amount.Amount)
		}

		// get a copy of the host chain, this is needed as GetHostChain returns a pointer.
		// the GenerateDelegateMessages messes with weights and later in the code we store the validators,
		// so in order to not modify the original chain object (and its validators) we need a copy of it.
		hcCopy, _ := k.GetHostChain(ctx, hc.ChainId)
		messages, err := k.GenerateDelegateMessages(hcCopy, totalDepositDelegation)
		if err != nil {
			k.Logger(ctx).Error(
				"could not simulate generating delegate messages after ICQ validator update",
				"host_chain",
				hc.ChainId,
			)
		}

		// there is a case where a validator can have much more delegation than what its weight dictates.
		// in this case, there won't be a message for that validator within the generated messages.
		// in order to not mess with the delegable state of the validator, and since we won't need to calculate the
		// room for it because it is not receiving any delegation, the default value of the room flag is true.
		validatorHasEnoughRoom := true
		for _, message := range messages {
			msgDelegate := message.(*stakingtypes.MsgDelegate)
			if validator.OperatorAddress == msgDelegate.ValidatorAddress {

				// calculate the amount of shares left to reach the validator lsm cap
				// shares * validator_lsm_cap - liquid_shares
				capRoom := validator.DelegatorShares.Mul(hc.Params.LsmValidatorCap).Sub(validator.LiquidShares)

				// if the bond factor functionality is disabled, calculate available room based only on the cap
				if hc.Params.LsmBondFactor.Equal(sdk.NewDecFromInt(sdk.NewInt(-1))) {
					validatorHasEnoughRoom = sdk.NewDecFromInt(msgDelegate.Amount.Amount).Quo(val.ExchangeRate).LT(capRoom)
					continue
				}

				// calculate the amount of shares left to reach the validator bond cap
				// bond_shares * bond_factor - liquid_shares
				bondRoom := validator.ValidatorBondShares.Mul(hc.Params.LsmBondFactor).Sub(validator.LiquidShares)

				// if there is room on both caps, the validator can accept the delegation we will send on the next
				// delegation workflow (next block).
				validatorHasEnoughRoom = sdk.NewDecFromInt(msgDelegate.Amount.Amount).Quo(val.ExchangeRate).LT(capRoom) &&
					sdk.NewDecFromInt(msgDelegate.Amount.Amount).Quo(val.ExchangeRate).LT(bondRoom)
			}
		}

		// recalculate the delegable state of the validator with the new flag
		val.Delegable = validatorHasRoomForDelegations && validatorHasEnoughBond && validatorHasEnoughRoom
		k.SetHostChainValidator(ctx, hc, val)

		// emit the delegable status event
		if oldDelegableFlag != val.Delegable {
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeValidatorDelegableStateUpdate,
					sdk.NewAttribute(types.AttributeChainID, hc.ChainId),
					sdk.NewAttribute(types.AttributeValidatorAddress, val.OperatorAddress),
					sdk.NewAttribute(types.AttributeKeyValidatorDelegable, strconv.FormatBool(val.Delegable)),
				),
			)
		}
	}

	return nil
}

func (k *Keeper) RedistributeValidatorWeight(ctx sdk.Context, hc *types.HostChain, validator *types.Validator) {
	validatorsWithWeight := make([]*types.Validator, 0)
	for _, val := range hc.Validators {
		if val.Weight.GT(sdk.ZeroDec()) && val.OperatorAddress != validator.OperatorAddress {
			validatorsWithWeight = append(validatorsWithWeight, val)
		}
	}

	weightDiff := validator.Weight.Quo(sdk.NewDec(int64(len(validatorsWithWeight))))
	for _, val := range validatorsWithWeight {
		val.Weight = val.Weight.Add(weightDiff)
		k.SetHostChainValidator(ctx, hc, val)
	}

	validator.Weight = sdk.ZeroDec()
	k.SetHostChainValidator(ctx, hc, validator)
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

// GetHostChainFromChannelID returns a host chain given its channel id
func (k *Keeper) GetHostChainFromChannelID(ctx sdk.Context, channelID string) (*types.HostChain, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HostChainKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	found := false
	hc := types.HostChain{}
	for ; iterator.Valid(); iterator.Next() {
		chain := types.HostChain{}
		k.cdc.MustUnmarshal(iterator.Value(), &chain)

		if chain.ChannelId == channelID {
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
