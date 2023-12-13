package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/proto"
	icqtypes "github.com/persistenceOne/persistence-sdk/v2/x/interchainquery/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

const (
	Validator                            = "validator"
	Delegation                           = "validator-delegation"
	RewardAccountBalances                = "reward-balances"
	NonCompoundableRewardAccountBalances = "non-compoundable-reward-balances"
	DelegationAccountBalances            = "delegation-balances"
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
		AddCallback(Validator, CallbackFn(ValidatorCallback)).
		AddCallback(RewardAccountBalances, CallbackFn(RewardsAccountBalanceCallback)).
		AddCallback(NonCompoundableRewardAccountBalances, CallbackFn(NonCompoundableRewardsAccountBalanceCallback)).
		AddCallback(DelegationAccountBalances, CallbackFn(DelegationAccountBalanceCallback)).
		AddCallback(Delegation, CallbackFn(DelegationCallback))

	return a.(Callbacks)
}

// Callbacks

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

	delegatedAmount := validator.ExchangeRate.Mul(delegation.Shares)
	slashedAmount := sdk.NewDecFromInt(validator.DelegatedAmount).Sub(delegatedAmount)

	if slashedAmount.IsPositive() {
		k.Logger(ctx).Info("Validator has been slashed !!!",
			"host-chain:", hc.ChainId,
			"validator:", validator.OperatorAddress,
			"slashed-amount:", slashedAmount,
		)

		// update the delegated amount to the slashed amount
		validator.DelegatedAmount = delegatedAmount.TruncateInt()
		k.SetHostChainValidator(ctx, hc, validator)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeSlashing,
				sdk.NewAttribute(types.AttributeChainID, hc.ChainId),
				sdk.NewAttribute(types.AttributeValidatorAddress, validator.OperatorAddress),
				sdk.NewAttribute(types.AttributeExistingDelegation, validator.DelegatedAmount.String()),
				sdk.NewAttribute(types.AttributeUpdatedDelegation, delegatedAmount.String()),
				sdk.NewAttribute(types.AttributeSlashedAmount, slashedAmount.String()),
			),
		)
	}

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

		// limit the auto-compounded rewards to the host chain autocompound factor
		var autocompoundRewards sdk.Coin
		maxAmountToTransfer := sdk.NewDecFromInt(hc.GetHostChainTotalDelegations()).Mul(hc.AutoCompoundFactor).TruncateInt()
		if maxAmountToTransfer.GT(hc.RewardsAccount.Balance.Amount) {
			autocompoundRewards = hc.RewardsAccount.Balance
		} else {
			autocompoundRewards = sdk.NewCoin(hc.RewardsAccount.Balance.Denom, maxAmountToTransfer)
		}

		// send all the rewards account balance to the deposit account, so it can be re-staked
		_, err = k.SendICATransfer(
			ctx,
			hc,
			autocompoundRewards,
			hc.RewardsAccount.Address,
			authtypes.NewModuleAddress(types.DepositModuleAccount).String(),
			hc.RewardsAccount.Owner,
		)
		if err != nil {
			return fmt.Errorf("could not send ICA rewards transfer: %w", err)
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeRewardsTransfer,
				sdk.NewAttribute(types.AttributeChainID, hc.ChainId),
				sdk.NewAttribute(types.AttributeRewardsTransferAmount, sdk.NewCoin(hc.HostDenom, autocompoundRewards.Amount).String()),
				sdk.NewAttribute(types.AttributeRewardsBalanceAmount, sdk.NewCoin(hc.HostDenom, hc.RewardsAccount.Balance.Amount.Sub(autocompoundRewards.Amount)).String()),
			),
		)
	}

	k.SetHostChain(ctx, hc)

	return nil
}

func NonCompoundableRewardsAccountBalanceCallback(k Keeper, ctx sdk.Context, data []byte, query icqtypes.Query) error {
	hc, found := k.GetHostChain(ctx, query.ChainId)
	if !found {
		return fmt.Errorf("host chain with id %s is not registered", query.ChainId)
	}

	balance, err := bankkeeper.UnmarshalBalanceCompat(k.cdc, data, hc.RewardParams.Denom)
	if err != nil {
		return fmt.Errorf("could unmarshal balance from ICQ balances request: %w", err)
	}

	if !balance.IsZero() {
		// build the transfer message to send the rewards to the swapping address
		msgTransfer := &banktypes.MsgSend{
			FromAddress: hc.RewardsAccount.Address,
			ToAddress:   hc.RewardParams.Destination,
			Amount:      sdk.NewCoins(balance),
		}

		// execute the ICA msgSend transaction
		_, err = k.GenerateAndExecuteICATx(
			ctx,
			hc.ConnectionId,
			hc.RewardsAccount.Owner,
			[]proto.Message{msgTransfer},
		)
		if err != nil {
			k.Logger(ctx).Error(
				"could not send ICA non-compoundable rewards transfer tx",
				"host_chain",
				hc.ChainId,
			)
			return fmt.Errorf("could not send ICA non-compoundable rewards transfer tx: %w", err)
		}
	}

	k.SetHostChain(ctx, hc)

	return nil
}
