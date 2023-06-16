package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	q "github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icqtypes "github.com/persistenceOne/persistence-sdk/v2/x/interchainquery/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

const (
	ValidatorSet              = "validatorset"
	Validator                 = "validator"
	RewardAccountBalances     = "reward-balances"
	DelegationAccountBalances = "delegation-balances"
	Delegation                = "validator-delegation"
)

type CallbackFn func(Keeper, sdk.Context, []byte, icqtypes.Query) error

type Callbacks struct {
	k         Keeper
	callbacks map[string]CallbackFn
}

var _ icqtypes.QueryCallbacks = Callbacks{}

func (k *Keeper) CallbackHandler() Callbacks {
	return Callbacks{*k, make(map[string]CallbackFn)}
}
func (c Callbacks) AddCallback(id string, fn interface{}) icqtypes.QueryCallbacks {
	c.callbacks[id] = fn.(CallbackFn)
	return c
}

func (c Callbacks) Call(ctx sdk.Context, id string, args []byte, query icqtypes.Query) error {
	return c.callbacks[id](c.k, ctx, args, query)
}

func (c Callbacks) Has(id string) bool {
	_, found := c.callbacks[id]
	return found
}

func (c Callbacks) RegisterCallbacks() icqtypes.QueryCallbacks {
	a := c.
		AddCallback(ValidatorSet, CallbackFn(ValidatorSetCallback)).
		AddCallback(Validator, CallbackFn(ValidatorCallback)).
		AddCallback(RewardAccountBalances, CallbackFn(RewardsAccountBalanceCallback)).
		AddCallback(DelegationAccountBalances, CallbackFn(DelegationAccountBalanceCallback)).
		AddCallback(Delegation, CallbackFn(DelegationCallback))

	return a.(Callbacks)
}

// Callbacks

func ValidatorSetCallback(k Keeper, ctx sdk.Context, data []byte, query icqtypes.Query) error {
	hc, found := k.GetHostChain(ctx, query.ChainId)
	if !found {
		return fmt.Errorf("host chain with id %s is not registered", query.ChainId)
	}

	response := stakingtypes.QueryValidatorsResponse{}
	err := k.cdc.Unmarshal(data, &response)
	if err != nil {
		return fmt.Errorf("could not unmarshall ICQ validatorset response: %w", err)
	}

	// if the result is not complete, submit an ICQ query to gather the next chunk
	if response.Pagination != nil && len(response.Pagination.NextKey) > 0 {
		request := stakingtypes.QueryValidatorsRequest{}
		err = k.cdc.Unmarshal(query.Request, &request)
		if err != nil {
			return fmt.Errorf("could not unmarshall ICQ validatorset request: %w", err)
		}

		request.Pagination = new(q.PageRequest)
		request.Pagination.Key = response.Pagination.NextKey
		if err = k.QueryHostChainValidators(ctx, hc, request); err != nil {
			return errorsmod.Wrapf(types.ErrFailedICQRequest, "error submitting validators icq: %s", err.Error())
		}
	}

	for _, validator := range response.Validators {
		val, found := hc.GetValidator(validator.OperatorAddress)

		// if it is a new validator or any of the attributes we track has changed, query for it
		if !found || (val != nil && (validator.Status.String() != val.Status ||
			!validator.DelegatorShares.Equal(val.DelegatorShares) || !validator.Tokens.Equal(val.TotalAmount))) {
			if err := k.QueryHostChainValidator(ctx, hc, validator.OperatorAddress); err != nil {
				return errorsmod.Wrapf(types.ErrFailedICQRequest, "error querying for validator: %s", err.Error())
			}
		}
	}

	return nil
}

func ValidatorCallback(k Keeper, ctx sdk.Context, data []byte, query icqtypes.Query) error {
	hc, found := k.GetHostChain(ctx, query.ChainId)
	if !found {
		return fmt.Errorf("host chain with id %s is not registered", query.ChainId)
	}

	validator, err := stakingtypes.UnmarshalValidator(k.cdc, data)
	if err != nil {
		return fmt.Errorf("could not unmarshall ICQ validator response: %w", err)
	}

	return k.ProcessHostChainValidatorUpdates(ctx, hc, validator)
}

func DelegationCallback(k Keeper, ctx sdk.Context, data []byte, query icqtypes.Query) error {
	hc, found := k.GetHostChain(ctx, query.ChainId)
	if !found {
		return fmt.Errorf("host chain with id %s is not registered", query.ChainId)
	}

	delegation, err := stakingtypes.UnmarshalDelegation(k.cdc, data)
	if err != nil {
		return fmt.Errorf("could not unmarshall ICQ delegation response: %w", err)
	}

	validator, found := hc.GetValidator(delegation.ValidatorAddress)
	if !found {
		return fmt.Errorf(
			"validator %s for host chain %s not found",
			delegation.ValidatorAddress,
			query.ChainId,
		)
	}

	delegatedAmount := validator.SharesToTokens(delegation.Shares)
	slashedAmount := validator.DelegatedAmount.Sub(delegatedAmount)

	if slashedAmount.IsPositive() {
		k.Logger(ctx).Info("Validator has been slashed !!!",
			"host-chain:", hc.ChainId,
			"validator:", validator.OperatorAddress,
			"slashed-amount:", slashedAmount,
		)

		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeSlashing,
				sdk.NewAttribute(types.AttributeValidatorAddress, validator.OperatorAddress),
				sdk.NewAttribute(types.AttributeExistingDelegation, validator.DelegatedAmount.String()),
				sdk.NewAttribute(types.AttributeUpdatedDelegation, delegatedAmount.String()),
				sdk.NewAttribute(types.AttributeSlashedAmount, slashedAmount.String()),
			)})
	}

	validator.DelegatedAmount = delegatedAmount
	k.SetHostChainValidator(ctx, hc, validator)

	return nil
}

func DelegationAccountBalanceCallback(k Keeper, ctx sdk.Context, data []byte, query icqtypes.Query) error {
	hc, found := k.GetHostChain(ctx, query.ChainId)
	if !found {
		return fmt.Errorf("host chain with id %s is not registered", query.ChainId)
	}

	balance, err := bankkeeper.UnmarshalBalanceCompat(k.cdc, data, hc.HostDenom)
	if err != nil {
		return fmt.Errorf("could unmarshal balance from ICQ balances request: %w", err)
	}

	hc.DelegationAccount.Balance = balance

	k.SetHostChain(ctx, hc)

	return nil
}

func RewardsAccountBalanceCallback(k Keeper, ctx sdk.Context, data []byte, query icqtypes.Query) error {
	hc, found := k.GetHostChain(ctx, query.ChainId)
	if !found {
		return fmt.Errorf("host chain with id %s is not registered", query.ChainId)
	}

	balance, err := bankkeeper.UnmarshalBalanceCompat(k.cdc, data, hc.HostDenom)
	if err != nil {
		return fmt.Errorf("could unmarshal balance from ICQ balances request: %w", err)
	}

	hc.RewardsAccount.Balance = balance
	if !hc.RewardsAccount.Balance.IsZero() {
		// send all the rewards account balance to the deposit account, so it can be re-staked
		_, err = k.SendICATransfer(
			ctx,
			hc,
			hc.RewardsAccount.Balance,
			hc.RewardsAccount.Address,
			authtypes.NewModuleAddress(types.DepositModuleAccount).String(),
			hc.RewardsAccount.Owner,
		)
		if err != nil {
			return fmt.Errorf("could not send ICA rewards transfer: %w", err)
		}
	}

	k.SetHostChain(ctx, hc)

	return nil
}
