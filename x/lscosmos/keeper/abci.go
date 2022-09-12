package keeper

import (
	"github.com/persistenceOne/persistence-sdk/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"

	lscosmostypes "github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

func (k Keeper) BeginBlock(ctx sdk.Context) {
	if !k.GetModuleState(ctx) {
		return
	}

	hostChainParams := k.GetHostChainParams(ctx)
	delegationState := k.GetDelegationState(ctx)

	delegateWrapperFn := func(ctx sdk.Context) error {
		return k.DoDelegate(ctx, delegationState.HostChainDelegationAddress, hostChainParams.ConnectionID, hostChainParams.BaseDenom)
	}
	err := utils.ApplyFuncIfNoError(ctx, delegateWrapperFn)
	if err != nil {
		k.Logger(ctx).Error("Unable to Delegate tokens with ", "err: ", err)
	}

}

func (k Keeper) DoDelegate(ctx sdk.Context, hostChainDelegationAddress, connectionID, baseDenom string) error {
	delegatableAmount := k.GetDelegationState(ctx).HostDelegationAccountBalance.AmountOf(baseDenom)
	allowlistedValidators := k.GetAllowListedValidators(ctx)
	if !delegatableAmount.GT(sdk.NewInt(int64(len(allowlistedValidators.AllowListedValidators)))) {
		// amount to delegate is too low, return early
		return nil
	}
	msgs := DelegateMsgs(hostChainDelegationAddress, allowlistedValidators, delegatableAmount, baseDenom)
	err := k.GenerateAndExecuteICATx(ctx, connectionID, lscosmostypes.DelegationAccountPortID, msgs)
	if err != nil {
		return err
	}

	amountToDelegate := sdk.NewCoin(baseDenom, delegatableAmount)
	k.RemoveBalanceFromDelegationState(ctx, sdk.NewCoins(amountToDelegate))
	k.AddICADelegateToTransitionStore(ctx, amountToDelegate)

	return nil
}
