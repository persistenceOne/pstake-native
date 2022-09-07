package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v3/modules/core/exported"
	"github.com/persistenceOne/persistence-sdk/utils"
	epochstypes "github.com/persistenceOne/persistence-sdk/x/epochs/types"
	ibchookertypes "github.com/persistenceOne/persistence-sdk/x/ibchooker/types"

	lscosmostypes "github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// BeforeEpochStart - call hook if registered
func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return nil
}

/*
AfterEpochEnd handle the "stake", "reward" and "undelegate" epoch and their respective actions
1. "stake" generates delegate transaction for delegating the amount of stake accumulated over the "stake" epoch
2. "reward" generates delegate transaction for withdrawing and restaking the amount of stake accumulated over the "reward" epochs
and shift the amount to next epoch if the min amount is not reached
3. "undelegate" generated the undelegate transaction for undelegating the amount accumulated over the "undelegate" epoch
*/
func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	//params := k.GetParams(ctx)
	if !k.GetModuleState(ctx) {
		return lscosmostypes.ErrModuleDisabled
	}
	hostChainParams := k.GetHostChainParams(ctx)
	if epochIdentifier == lscosmostypes.DelegationEpochIdentifier {
		wrapperFn := func(ctx sdk.Context) error {
			return k.DelegationEpochWorkFlow(ctx, hostChainParams)
		}
		err := utils.ApplyFuncIfNoError(ctx, wrapperFn)
		k.Logger(ctx).Error("Failed DelegationEpochIdentifier Function with:", "err: ", err)
	}
	if epochIdentifier == lscosmostypes.RewardEpochIdentifier {
		wrapperFn := func(ctx sdk.Context) error {
			return k.RewardEpochEpochWorkFlow(ctx, hostChainParams)
		}
		err := utils.ApplyFuncIfNoError(ctx, wrapperFn)
		k.Logger(ctx).Error("Failed RewardEpochIdentifier Function with:", "err: ", err)
	}
	if epochIdentifier == lscosmostypes.UndelegationEpochIdentifier {
		wrapperFn := func(ctx sdk.Context) error {
			return k.UndelegationEpochWorkFlow(ctx, hostChainParams)
		}
		err := utils.ApplyFuncIfNoError(ctx, wrapperFn)
		k.Logger(ctx).Error("Failed UndelegationEpochIdentifier Function with:", "err: ", err)
	}
	return nil
}

// ___________________________________________________________________________________________________

// EpochsHooks wrapper struct for incentives keeper
type EpochsHooks struct {
	k Keeper
}

var _ epochstypes.EpochHooks = EpochsHooks{}

// Return the wrapper struct
func (k Keeper) NewEpochHooks() EpochsHooks {
	return EpochsHooks{k}
}

// epochs hooks
func (h EpochsHooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return h.k.BeforeEpochStart(ctx, epochIdentifier, epochNumber)
}

func (h EpochsHooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return h.k.AfterEpochEnd(ctx, epochIdentifier, epochNumber)
}

func (k Keeper) DelegationEpochWorkFlow(ctx sdk.Context, hostChainParams lscosmostypes.HostChainParams) error {
	// greater than min amount, transfer from deposit to delegation, to ibctransfer.
	// Right now we only do baseDenom
	ibcDenom := k.GetIBCDenom(ctx)
	allDepositBalances := k.bankKeeper.GetAllBalances(ctx, authtypes.NewModuleAddress(lscosmostypes.DepositModuleAccount))
	depositBalance := sdk.NewCoin(ibcDenom, allDepositBalances.AmountOf(ibcDenom))
	if depositBalance.Amount.GT(sdk.ZeroInt()) {
		err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, lscosmostypes.DepositModuleAccount, lscosmostypes.DelegationModuleAccount, sdk.NewCoins(depositBalance))
		if err != nil {
			k.Logger(ctx).Error("Could not send amount from ", lscosmostypes.DepositModuleAccount, " module account to ",
				lscosmostypes.DelegationModuleAccount)
			return err
		}
	}
	allDelegationBalances := k.bankKeeper.GetAllBalances(ctx, authtypes.NewModuleAddress(lscosmostypes.DelegationModuleAccount))
	delegationBalance := sdk.NewCoin(ibcDenom, allDelegationBalances.AmountOf(ibcDenom))
	delegationState := k.GetDelegationState(ctx)
	_, clientState, err := k.channelKeeper.GetChannelClientState(ctx, hostChainParams.TransferPort, hostChainParams.TransferChannel)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("Error getting client state %s", err))
		return err
	}
	timeoutHeight := clienttypes.NewHeight(clientState.GetLatestHeight().GetRevisionNumber(), clientState.GetLatestHeight().GetRevisionHeight()+lscosmostypes.IBCTimeoutHeightIncrement)

	msg := ibctransfertypes.NewMsgTransfer(hostChainParams.TransferPort, hostChainParams.TransferChannel,
		delegationBalance, authtypes.NewModuleAddress(lscosmostypes.DelegationModuleAccount).String(),
		delegationState.HostChainDelegationAddress, timeoutHeight, 0)

	handler := k.msgRouter.Handler(msg)

	res, err := handler(ctx, msg)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("could not send transfer msg via MsgServiceRouter, error: %s", err))
		return err
	}

	ctx.EventManager().EmitEvents(res.GetEvents())

	// move extra tokens to pstake address - anyone can send tokens to delegation address.
	// deposit address is deny-listed address - can only accept tokens via transactions, so should not have any extra tokens
	// should be transferred to pstake address.
	remainingDelegationBalance := allDelegationBalances.Sub(sdk.NewCoins(delegationBalance))

	if !remainingDelegationBalance.Empty() {
		feeAddr := sdk.MustAccAddressFromBech32(hostChainParams.PstakeFeeAddress)
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, lscosmostypes.DelegationModuleAccount, feeAddr, remainingDelegationBalance)
		if err != nil {
			k.Logger(ctx).Error(fmt.Sprintf("could not send remaining balance: %s in delegationModuleAccount: %s with error: %s", remainingDelegationBalance, lscosmostypes.DelegationModuleAccount, err))
			return err
		}
	}

	return nil
}

func (k Keeper) RewardEpochEpochWorkFlow(ctx sdk.Context, hostChainParams lscosmostypes.HostChainParams) error {
	// send withdraw rewards from delegators.
	delegationState := k.GetDelegationState(ctx)
	if len(delegationState.HostAccountDelegations) == 0 {
		return lscosmostypes.ErrNoHostChainDelegations
	}
	withdrawRewardMsgs := make([]sdk.Msg, len(delegationState.HostAccountDelegations))
	for i, delegation := range delegationState.HostAccountDelegations {
		withdrawRewardMsgs[i] = &distributiontypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: delegationState.HostChainDelegationAddress,
			ValidatorAddress: delegation.ValidatorAddress,
		}
	}
	err := k.GenerateAndExecuteICATx(ctx, hostChainParams.ConnectionID, lscosmostypes.DelegationAccountPortID, withdrawRewardMsgs)
	return err
	// on Ack do icq for reward acc. balance of uatom
	// callback for sending it to delegation account
	// on Ack delegate txn
}

func (k Keeper) UndelegationEpochWorkFlow(ctx sdk.Context, hostChainParams lscosmostypes.HostChainParams) error {
	return nil
}

// ___________________________________________________________________________________________________

func (k Keeper) OnRecvIBCTransferPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress, transferAck ibcexported.Acknowledgement) error {
	return nil
}

func (k Keeper) OnAcknowledgementIBCTransferPacket(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte, relayer sdk.AccAddress, transferAckErr error) error {
	if !k.GetModuleState(ctx) {
		return lscosmostypes.ErrModuleDisabled
	}
	if transferAckErr != nil {
		return transferAckErr
	}
	var ack channeltypes.Acknowledgement
	if err := ibctransfertypes.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return err
	}
	if !ack.Success() {
		return channeltypes.ErrInvalidAcknowledgement
	}
	var data ibctransfertypes.FungibleTokenPacketData
	if err := ibctransfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return err
	}
	// check for tokens moved from delegationModuleAccount to it's ica counterpart.
	hostChainParams := k.GetHostChainParams(ctx)
	delegationState := k.GetDelegationState(ctx)
	if packet.GetSourceChannel() != hostChainParams.TransferChannel ||
		packet.GetSourcePort() != hostChainParams.TransferPort {
		// no need to return err, since most likely code is expected to enter this condition
		return nil
	}

	if data.GetSender() != k.GetDelegationModuleAccount(ctx).GetAddress().String() ||
		data.GetReceiver() != delegationState.HostChainDelegationAddress ||
		data.GetDenom() != ibctransfertypes.GetPrefixedDenom(hostChainParams.TransferPort, hostChainParams.TransferChannel, hostChainParams.BaseDenom) {
		// no need to return err, since most likely code is expected to enter this condition
		return nil
	}
	k.Logger(ctx).Info(fmt.Sprintf("pstake tokens successfully transferred to host chain address %s, amount: %s, denom: %s", data.Receiver, data.Amount, data.Denom))

	//do ica delegate.
	amount, ok := sdk.NewIntFromString(data.GetAmount())
	if !ok {
		return ibctransfertypes.ErrInvalidAmount
	}
	k.AddBalanceToDelegationState(ctx, sdk.NewCoin(hostChainParams.BaseDenom, amount))

	return nil
}

func (k Keeper) OnTimeoutIBCTransferPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress, transferTimeoutErr error) error {
	// Do nothing because amount will be reverted to delegationModuleAccount.
	return nil
}

type IBCTransferHooks struct {
	k Keeper
}

var _ ibchookertypes.IBCHandshakeHooks = IBCTransferHooks{}

func (k Keeper) NewIBCTransferHooks() IBCTransferHooks {
	return IBCTransferHooks{k}
}

func (i IBCTransferHooks) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress, transferAck ibcexported.Acknowledgement) error {
	return i.k.OnRecvIBCTransferPacket(ctx, packet, relayer, transferAck)
}

func (i IBCTransferHooks) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte, relayer sdk.AccAddress, transferAckErr error) error {
	return i.k.OnAcknowledgementIBCTransferPacket(ctx, packet, acknowledgement, relayer, transferAckErr)

}

func (i IBCTransferHooks) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress, transferTimeoutErr error) error {
	return i.k.OnTimeoutIBCTransferPacket(ctx, packet, relayer, transferTimeoutErr)
}
