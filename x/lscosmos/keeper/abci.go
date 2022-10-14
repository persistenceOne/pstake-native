package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	"github.com/persistenceOne/persistence-sdk/utils"

	lscosmostypes "github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

func (k Keeper) BeginBlock(ctx sdk.Context) {
	if !k.GetModuleState(ctx) {
		return
	}

	err := utils.ApplyFuncIfNoError(ctx, k.DoDelegate)
	if err != nil {
		k.Logger(ctx).Error("Unable to Delegate tokens with ", "err: ", err)
	}
	err = utils.ApplyFuncIfNoError(ctx, k.ProcessMaturedUndelegation)
	if err != nil {
		k.Logger(ctx).Error("Unable to Delegate tokens with ", "err: ", err)
	}

}

func (k Keeper) DoDelegate(ctx sdk.Context) error {
	hostChainParams := k.GetHostChainParams(ctx)
	delegationState := k.GetDelegationState(ctx)
	hostAccounts := k.GetHostAccounts(ctx)

	delegatableAmount := delegationState.HostDelegationAccountBalance.AmountOf(hostChainParams.BaseDenom)
	if delegatableAmount.LT(hostChainParams.MinDeposit) {
		// amount to delegate is too low, return early
		return nil
	}
	msgs, err := k.DelegateMsgs(ctx, delegationState.HostChainDelegationAddress, delegatableAmount, hostChainParams.BaseDenom)
	if err != nil {
		return err
	}
	err = k.GenerateAndExecuteICATx(ctx, hostChainParams.ConnectionID, hostAccounts.DelegatorAccountPortID(), msgs)
	if err != nil {
		return err
	}

	amountToDelegate := sdk.NewCoin(hostChainParams.BaseDenom, delegatableAmount)
	k.RemoveBalanceFromDelegationState(ctx, sdk.NewCoins(amountToDelegate))
	k.AddICADelegateToTransientStore(ctx, amountToDelegate)

	return nil
}

func (k Keeper) ProcessMaturedUndelegation(ctx sdk.Context) error {
	hostChainParams := k.GetHostChainParams(ctx)
	delegationState := k.GetDelegationState(ctx)
	hostAccounts := k.GetHostAccounts(ctx)

	maturedUndelegations := k.GetHostAccountMaturedUndelegations(ctx)
	if len(maturedUndelegations) == 0 {
		// No matured delegations
		return nil
	}
	for _, maturedUndelegation := range maturedUndelegations {
		//do ica ibc transfer + delete the entries
		atomsUnbonded := k.GetUnbondingEpochCValue(ctx, maturedUndelegation.EpochNumber).AmountUnbonded

		channel, found := k.channelKeeper.GetChannel(ctx, hostChainParams.TransferPort, hostChainParams.TransferChannel)
		if !found {
			return channeltypes.ErrChannelNotFound
		}

		selfHeight := clienttypes.GetSelfHeight(ctx)
		timeoutHeight := clienttypes.NewHeight(selfHeight.GetRevisionNumber(), selfHeight.GetRevisionHeight()+lscosmostypes.IBCTimeoutHeightIncrement)

		msg := ibctransfertypes.NewMsgTransfer(channel.Counterparty.PortId, channel.Counterparty.ChannelId,
			atomsUnbonded, delegationState.HostChainDelegationAddress, authtypes.NewModuleAddress(lscosmostypes.UndelegationModuleAccount).String(), timeoutHeight, 0)
		err := k.GenerateAndExecuteICATx(ctx, hostChainParams.ConnectionID, hostAccounts.DelegatorAccountPortID(), []sdk.Msg{msg})
		if err != nil {
			return err
		}
		err = k.RemoveHostAccountUndelegation(ctx, maturedUndelegation.EpochNumber)
		if err != nil {
			return err
		}
		k.AddUndelegationTransferToTransientStore(ctx, lscosmostypes.TransientUndelegationTransfer{
			EpochNumber:    maturedUndelegation.EpochNumber,
			AmountUnbonded: atomsUnbonded,
		})
	}
	return nil
}
