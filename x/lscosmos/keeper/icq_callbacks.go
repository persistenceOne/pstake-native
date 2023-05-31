package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/proto"
	icqtypes "github.com/persistenceOne/persistence-sdk/v2/x/interchainquery/types"

	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
)

const (
	RewardsAccountBalance = "reward_account_balance"
	Delegation            = "delegation"
)

// CallbackFn wrapper struct for interchainstaking keeper
type CallbackFn func(Keeper, sdk.Context, []byte, icqtypes.Query) error

type Callbacks struct {
	k         Keeper
	callbacks map[string]CallbackFn
}

var _ icqtypes.QueryCallbacks = Callbacks{}

// CallbackHandler returns Callbacks with empty entries
func (k Keeper) CallbackHandler() Callbacks {
	return Callbacks{k, make(map[string]CallbackFn)}
}

// AddCallback adds callback using the input id and interface
func (c Callbacks) AddCallback(id string, fn interface{}) icqtypes.QueryCallbacks {
	c.callbacks[id] = fn.(CallbackFn)
	return c
}

// RegisterCallbacks adds callbacks
func (c Callbacks) RegisterCallbacks() icqtypes.QueryCallbacks {
	a := c.
		AddCallback(RewardsAccountBalance, CallbackFn(RewardsAccountBalanceCallback)).
		AddCallback(Delegation, CallbackFn(DelegationCallback))

	return a.(Callbacks)
}

// Call returns callback based on the input id, args and query
func (c Callbacks) Call(ctx sdk.Context, id string, args []byte, query icqtypes.Query) error {
	return c.callbacks[id](c.k, ctx, args, query)
}

// Has checks and returns if input id is present in callbacks
func (c Callbacks) Has(id string) bool {
	_, found := c.callbacks[id]
	return found
}

// RewardsAccountBalanceCallback returns response of HandleRewardsAccountBalanceCallback
func RewardsAccountBalanceCallback(k Keeper, ctx sdk.Context, response []byte, query icqtypes.Query) error {
	return k.HandleRewardsAccountBalanceCallback(ctx, response, query)
}

// DelegationCallback returns response of HandleDelegationCallback
func DelegationCallback(k Keeper, ctx sdk.Context, response []byte, query icqtypes.Query) error {
	return k.HandleDelegationCallback(ctx, response, query)
}

// HandleRewardsAccountBalanceCallback generates and executes rewards account balance query
func (k Keeper) HandleRewardsAccountBalanceCallback(ctx sdk.Context, response []byte, _ icqtypes.Query) error {
	resp := banktypes.QueryBalanceResponse{}
	err := k.cdc.Unmarshal(response, &resp)
	if err != nil {
		return err
	}
	k.Logger(ctx).Info("Callback for Rewards account balance", "Balances", resp.GetBalance())

	if resp.Balance.Amount.Equal(sdk.ZeroInt()) {
		k.Logger(ctx).Info("No amount in rewards account to restake - noop.")
		return nil
	}

	hostChainParams := k.GetHostChainParams(ctx)
	delegationState := k.GetDelegationState(ctx)
	rewardsAddress := k.GetHostChainRewardAddress(ctx)
	hostAccounts := k.GetHostAccounts(ctx)

	// Cap the re-staking amount so exchange rate doesn't change drastically.
	cValue := k.GetCValue(ctx)
	stkAssetSupply := k.bankKeeper.GetSupply(ctx, hostChainParams.MintDenom)

	atomTVU := sdk.NewDecFromInt(stkAssetSupply.Amount).Quo(cValue)
	atomTVUCap := atomTVU.Mul(types.RestakeCapPerDay).TruncateInt()
	sendCoinAmt := resp.Balance.Amount
	if resp.Balance.Amount.GT(atomTVUCap) {
		sendCoinAmt = atomTVUCap
	}

	//send coins to delegation account.
	msg := &banktypes.MsgSend{
		FromAddress: rewardsAddress.Address,
		ToAddress:   delegationState.HostChainDelegationAddress,
		Amount:      sdk.NewCoins(sdk.NewCoin(resp.Balance.Denom, sendCoinAmt)),
	}
	return k.GenerateAndExecuteICATx(ctx, hostChainParams.ConnectionID, hostAccounts.RewardsAccountOwnerID, []proto.Message{msg})
}

// HandleDelegationCallback generates and executes delegation query
func (k Keeper) HandleDelegationCallback(ctx sdk.Context, response []byte, _ icqtypes.Query) error {
	resp := stakingtypes.QueryDelegationResponse{}
	err := k.cdc.Unmarshal(response, &resp)
	if err != nil {
		return err
	}
	k.Logger(ctx).Info("Callback for Validator Delegation", "Response: ", resp.GetDelegationResponse())

	//check ack sequences for ica accs, return error for the callback, so it can be retried.
	pending, err := k.CheckPendingICATxs(ctx)
	if pending {
		return err
	}

	existingDelegation := k.GetHostAccountDelegation(ctx, resp.GetDelegationResponse().Delegation.ValidatorAddress)
	if resp.GetDelegationResponse().GetBalance().IsLT(existingDelegation.Amount) {
		//log slashing
		k.Logger(ctx).Info("Received delegation less than delegation-state ",
			"validator:", resp.GetDelegationResponse().Delegation.ValidatorAddress,
			"delegationState:", existingDelegation.Amount,
			"hostDelegation:", resp.GetDelegationResponse().Balance)
		// emit event slashing fixed
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypePerformSlashing,
				sdk.NewAttribute(types.AttributeValidatorAddress, resp.GetDelegationResponse().Delegation.ValidatorAddress),
				sdk.NewAttribute(types.AttributeExistingDelegation, existingDelegation.Amount.String()),
				sdk.NewAttribute(types.AttributeUpdatedDelegation, resp.GetDelegationResponse().Balance.String()),
				sdk.NewAttribute(types.AttributeSlashedAmount, existingDelegation.Amount.Sub(resp.GetDelegationResponse().Balance).String()),
			)})
		k.ForceUpdateHostAccountDelegation(ctx, types.NewHostAccountDelegation(resp.GetDelegationResponse().Delegation.ValidatorAddress, resp.GetDelegationResponse().GetBalance()))
	}
	return nil
}
