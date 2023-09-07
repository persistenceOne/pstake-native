package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	icqtypes "github.com/persistenceOne/persistence-sdk/v2/x/interchainquery/types"
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
	return nil
}

// HandleDelegationCallback generates and executes delegation query
func (k Keeper) HandleDelegationCallback(ctx sdk.Context, _ []byte, _ icqtypes.Query) error {
	return nil
}
