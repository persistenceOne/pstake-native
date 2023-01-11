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

// BeginBlock will use utils.ApplyFuncIfNoError to apply the changes made by the functions
// passed as parameters
func (k Keeper) BeginBlock(ctx sdk.Context) {
	if !k.GetModuleState(ctx) {
		return
	}

	//fork logic, halt height + 1
	if ctx.BlockHeight() == 9616501 {
		err := MintPstakeTokens(ctx, k)
		if err != nil {
			panic(err)
		}
	}

	err := utils.ApplyFuncIfNoError(ctx, k.DoDelegate)
	if err != nil {
		k.Logger(ctx).Error("Unable to Delegate tokens with ", "err: ", err)
	}
	err = utils.ApplyFuncIfNoError(ctx, k.ProcessMaturedUndelegation)
	if err != nil {
		k.Logger(ctx).Error("Unable to process matured undelegations with ", "err: ", err)
	}

}

// DoDelegate generates and executes ICA transactions based on the generated delegation state
// from DelegateMsgs
func (k Keeper) DoDelegate(ctx sdk.Context) error {
	hostChainParams := k.GetHostChainParams(ctx)
	delegationState := k.GetDelegationState(ctx)
	delegatableAmount := delegationState.HostDelegationAccountBalance.AmountOf(hostChainParams.BaseDenom)

	allowListedValidators := k.GetAllowListedValidators(ctx)
	if !delegatableAmount.IsPositive() || len(allowListedValidators.AllowListedValidators) == 0 {
		// amount to delegate is too low, return early
		return nil
	}

	// generate delegate messages based on the delegatable amount and current validators
	// delegation state
	msgs, err := k.DelegateMsgs(ctx, delegatableAmount, hostChainParams.BaseDenom, delegationState)
	if err != nil {
		return err
	}

	// get host accounts and use them to generate and execute ICA tx for delegations.
	hostAccounts := k.GetHostAccounts(ctx)
	err = k.GenerateAndExecuteICATx(ctx, hostChainParams.ConnectionID, hostAccounts.DelegatorAccountPortID(), msgs)
	if err != nil {
		return err
	}

	amountToDelegate := sdk.NewCoin(hostChainParams.BaseDenom, delegatableAmount)
	k.RemoveBalanceFromDelegationState(ctx, sdk.NewCoins(amountToDelegate))
	k.AddICADelegateToTransientStore(ctx, amountToDelegate)

	return nil
}

// ProcessMaturedUndelegation processes all the matured undelegations by fetching all the host
// account matured undelegations and processing them one by one
func (k Keeper) ProcessMaturedUndelegation(ctx sdk.Context) error {
	// check if there are any matured undelegations
	maturedUndelegations := k.GetHostAccountMaturedUndelegations(ctx)
	if len(maturedUndelegations) == 0 {
		// No matured delegations
		return nil
	}

	hostChainParams := k.GetHostChainParams(ctx)
	delegationState := k.GetDelegationState(ctx)
	hostAccounts := k.GetHostAccounts(ctx)

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

func MintPstakeTokens(ctx sdk.Context, k Keeper) error {
	if ctx.ChainID() != "core-1" {
		return nil
	}

	atomTVU := k.GetDepositAccountAmount(ctx).
		Add(k.GetIBCTransferTransientAmount(ctx)).
		Add(k.GetDelegationTransientAmount(ctx)).
		Add(k.GetStakedAmount(ctx)).
		Add(k.GetHostDelegationAccountAmount(ctx))

	mintedAmount := k.GetMintedAmount(ctx)
	mintDenom := k.GetHostChainParams(ctx).MintDenom
	if atomTVU.LTE(mintedAmount) {
		return nil
	}

	toNewMint := atomTVU.Sub(mintedAmount)

	switch ctx.ChainID() {
	case "core-1":
		pstakeFeeAddress := sdk.MustAccAddressFromBech32(k.GetHostChainParams(ctx).PstakeParams.PstakeFeeAddress)
		err := k.MintTokens(ctx, sdk.NewCoin(mintDenom, toNewMint), pstakeFeeAddress)
		if err != nil {
			k.Logger(ctx).Error("Failed to mint and send remainingAmount to pstakeFeeAddress")
			return err
		}
		return nil
	default:
		return nil
	}
}
