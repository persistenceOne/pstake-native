package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	q "github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icqtypes "github.com/persistenceOne/persistence-sdk/v2/x/interchainquery/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

const (
	ValidatorSet = "validatorset"
	Balances     = "balances"
	Delegation   = "delegation"
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
		AddCallback(Balances, CallbackFn(BalancesCallback)).
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

	return k.ProcessHostChainValidatorUpdates(ctx, hc, response.Validators)
}

func BalancesCallback(k Keeper, ctx sdk.Context, data []byte, query icqtypes.Query) error {
	hc, found := k.GetHostChain(ctx, query.ChainId)
	if !found {
		return fmt.Errorf("host chain with id %s is not registered", query.ChainId)
	}

	response := banktypes.QueryBalanceResponse{}
	err := k.cdc.Unmarshal(data, &response)
	if err != nil {
		return fmt.Errorf("could not unmarshall ICQ balances response: %w", err)
	}

	request := banktypes.QueryBalanceRequest{}
	err = k.cdc.Unmarshal(query.Request, &request)
	if err != nil {
		return fmt.Errorf("could not unmarshall ICQ balances request: %w", err)
	}

	switch request.Address {
	case hc.DelegationAccount.Address:
		hc.DelegationAccount.Balance = *response.Balance
	case hc.RewardsAccount.Address:
		hc.RewardsAccount.Balance = *response.Balance
		if !hc.RewardsAccount.Balance.IsZero() {
			// send all the rewards account balance to the deposit account, so it can be re-staked
			_, err = k.SendICATransfer(
				ctx,
				hc,
				hc.RewardsAccount.Balance,
				hc.RewardsAccount.Address,
				authtypes.NewModuleAddress(types.DepositModuleAccount).String(),
				k.RewardsAccountPortOwner(hc.ChainId),
			)
			if err != nil {
				return fmt.Errorf("could not send ICA rewards transfer: %w", err)
			}
		}
	default:
		return fmt.Errorf("address doesn't belong to any ICA accout of the host chain with id %s", query.ChainId)
	}

	// recalculate the host chain c value after the local account data has been updated
	hc.CValue = k.GetHostChainCValue(ctx, hc)
	k.SetHostChain(ctx, hc)

	return nil
}

func DelegationCallback(k Keeper, ctx sdk.Context, data []byte, query icqtypes.Query) error {
	hc, found := k.GetHostChain(ctx, query.ChainId)
	if !found {
		return fmt.Errorf("host chain with id %s is not registered", query.ChainId)
	}

	response := stakingtypes.QueryDelegationResponse{}
	err := k.cdc.Unmarshal(data, &response)
	if err != nil {
		return fmt.Errorf("could not unmarshall ICQ delegation response: %w", err)
	}

	validator, found := hc.GetValidator(response.DelegationResponse.Delegation.ValidatorAddress)
	if !found {
		return fmt.Errorf(
			"validator %s for host chain %s not found",
			response.DelegationResponse.Delegation.ValidatorAddress,
			query.ChainId,
		)
	}

	if response.DelegationResponse.Balance.Amount.LT(validator.DelegatedAmount) {
		slashedAmount := validator.DelegatedAmount.Sub(response.DelegationResponse.Balance.Amount)

		k.Logger(ctx).Info("Validator has ben slashed !!!",
			"host-chain:", hc.ChainId,
			"validator:", validator.OperatorAddress,
			"slashed-amount:", slashedAmount,
		)

		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeSlashing,
				sdk.NewAttribute(types.AttributeValidatorAddress, validator.OperatorAddress),
				sdk.NewAttribute(types.AttributeExistingDelegation, validator.DelegatedAmount.String()),
				sdk.NewAttribute(types.AttributeUpdatedDelegation, response.DelegationResponse.Balance.Amount.String()),
				sdk.NewAttribute(types.AttributeSlashedAmount, slashedAmount.String()),
			)})
	}

	validator.DelegatedAmount = response.DelegationResponse.Balance.Amount
	k.SetHostChainValidator(ctx, hc, validator)

	return nil
}
