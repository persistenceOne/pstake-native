package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	icqtypes "github.com/persistenceOne/persistence-sdk/x/interchainquery/types"

	lscosmostypes "github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

const (
	RewardsAccountBalance = "reward_account_balance"
)

// Callbacks wrapper struct for interchainstaking keeper
type CallbackFn func(Keeper, sdk.Context, []byte, icqtypes.Query) error

type Callbacks struct {
	k         Keeper
	callbacks map[string]CallbackFn
}

var _ icqtypes.QueryCallbacks = Callbacks{}

func (k Keeper) CallbackHandler() Callbacks {
	return Callbacks{k, make(map[string]CallbackFn)}
}

func (c Callbacks) AddCallback(id string, fn interface{}) icqtypes.QueryCallbacks {
	c.callbacks[id] = fn.(CallbackFn)
	return c
}

func (c Callbacks) RegisterCallbacks() icqtypes.QueryCallbacks {
	a := c.
		AddCallback(RewardsAccountBalance, CallbackFn(RewardsAccountBalanceCallback))

	return a.(Callbacks)
}

func (c Callbacks) Call(ctx sdk.Context, id string, args []byte, query icqtypes.Query) error {
	return c.callbacks[id](c.k, ctx, args, query)
}

func (c Callbacks) Has(id string) bool {
	_, found := c.callbacks[id]
	return found
}

func RewardsAccountBalanceCallback(k Keeper, ctx sdk.Context, response []byte, query icqtypes.Query) error {
	return k.HandleRewardsAccountBalanceCallback(ctx, response, query)
}

func (k Keeper) HandleRewardsAccountBalanceCallback(ctx sdk.Context, response []byte, query icqtypes.Query) error {
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
	//send coins to delegation account.
	msg := &banktypes.MsgSend{
		FromAddress: rewardsAddress.Address,
		ToAddress:   delegationState.HostChainDelegationAddress,
		Amount:      sdk.NewCoins(*resp.GetBalance()),
	}
	return k.GenerateAndExecuteICATx(ctx, hostChainParams.ConnectionID, lscosmostypes.RewardAccountPortID, []sdk.Msg{msg})
}
