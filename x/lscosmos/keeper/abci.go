package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	lscosmostypes "github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

func (k Keeper) BeginBlock(ctx sdk.Context) {
	if !k.GetModuleState(ctx) {
		return
	}

	hostChainParams := k.GetHostChainParams(ctx)
	delegationState := k.GetDelegationState(ctx)
	//TODO handle err like osmosis
	_ = k.DoDelegate(ctx, delegationState.HostChainDelegationAddress, hostChainParams.ConnectionID, hostChainParams.BaseDenom)

}

func (k Keeper) DoDelegate(ctx sdk.Context, hostChainDelegationAddress, connectionID, baseDenom string) error {
	delegatableAmount := k.GetDelegationState(ctx).HostDelegationAccountBalance.AmountOf(baseDenom)
	allowlistedValidators := k.GetAllowListedValidators(ctx)
	if !delegatableAmount.GT(sdk.NewInt(int64(len(allowlistedValidators.AllowListedValidators)))) {
		k.Logger(ctx).Info(fmt.Sprintf("amount is too low to delegate, %v ", delegatableAmount))
		return nil
	}
	msgs := DelegateMsgs(hostChainDelegationAddress, allowlistedValidators, delegatableAmount, baseDenom)
	err := k.GenerateAndExecuteICATx(ctx, connectionID, lscosmostypes.DelegationAccountPortID, msgs)
	if err != nil {
		return err
	}
	k.RemoveBalanceFromDelegationState(ctx, sdk.NewCoins(sdk.NewCoin(baseDenom, delegatableAmount)))
	return nil
}
