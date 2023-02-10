package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	icqtypes "github.com/persistenceOne/persistence-sdk/v2/x/interchainquery/types"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

const (
	RewardsAccountBalance = "reward_account_balance"
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
		AddCallback(RewardsAccountBalance, CallbackFn(RewardsAccountBalanceCallback))

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

// HandleRewardsAccountBalanceCallback generates and executes rewards account balance query
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
	hostAccounts := k.GetHostAccounts(ctx)

	// Cap the re-staking amount so exchange rate doesn't change drastically.
	cValue := k.GetCValue(ctx)
	stkAssetSupply := k.bankKeeper.GetSupply(ctx, hostChainParams.MintDenom)
	atomTVU := stkAssetSupply.Amount.ToDec().Quo(cValue)
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
	return k.GenerateAndExecuteICATx(ctx, hostChainParams.ConnectionID, hostAccounts.RewardsAccountPortID(), []sdk.Msg{msg})
}
