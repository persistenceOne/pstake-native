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
	epochstypes "github.com/persistenceOne/persistence-sdk/x/epochs/types"
	ibchookertypes "github.com/persistenceOne/persistence-sdk/x/ibchooker/types"

	lscosmostypes "github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// BeforeEpochStart - call hook if registered
func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
}

/*
AfterEpochEnd handle the "stake", "reward" and "undelegate" epoch and their respective actions
1. "stake" generates delegate transaction for delegating the amount of stake accumulated over the "stake" epoch
2. "reward" generates delegate transaction for withdrawing and restaking the amount of stake accumulated over the "reward" epochs
and shift the amount to next epoch if the min amount is not reached
3. "undelegate" generated the undelegate transaction for undelegating the amount accumulated over the "undelegate" epoch
*/
func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	//params := k.GetParams(ctx)
	if !k.GetModuleState(ctx) {
		return
	}
	hostChainParams := k.GetHostChainParams(ctx)
	if epochIdentifier == lscosmostypes.DelegationEpochIdentifier {
		k.DelegationEpochWorkFlow(ctx, hostChainParams)
	}
	if epochIdentifier == lscosmostypes.RewardEpochIdentifier {
		k.RewardEpochEpochWorkFlow(ctx, hostChainParams)
	}
	if epochIdentifier == lscosmostypes.UndelegationEpochIdentifier {
		k.UndelegationEpochWorkFlow(ctx, hostChainParams)
	}
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
func (h EpochsHooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.BeforeEpochStart(ctx, epochIdentifier, epochNumber)
}

func (h EpochsHooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.AfterEpochEnd(ctx, epochIdentifier, epochNumber)
}

func (k Keeper) DelegationEpochWorkFlow(ctx sdk.Context, hostChainParams lscosmostypes.HostChainParams) {
	// greater than min amount, transfer from deposit to delegation, to ibctransfer.
	// Right now we only do baseDenom
	ibcDenom := ibctransfertypes.ParseDenomTrace(
		ibctransfertypes.GetPrefixedDenom(
			hostChainParams.TransferPort, hostChainParams.TransferChannel, hostChainParams.BaseDenom,
		),
	).IBCDenom()
	allBalances := k.bankKeeper.GetAllBalances(ctx, authtypes.NewModuleAddress(lscosmostypes.DepositModuleAccount))
	depositBalance := sdk.NewCoin(ibcDenom, allBalances.AmountOf(ibcDenom))
	if !depositBalance.Amount.GT(sdk.ZeroInt()) {
		return
	}
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, lscosmostypes.DepositModuleAccount, lscosmostypes.DelegationModuleAccount, sdk.NewCoins(depositBalance))
	if err != nil {
		k.Logger(ctx).Info("Could not send amount from ", lscosmostypes.DepositModuleAccount, " module account to ",
			lscosmostypes.DelegationModuleAccount)
		return
	}

	delegationState := k.GetDelegationState(ctx)
	_, clientState, err := k.channelKeeper.GetChannelClientState(ctx, hostChainParams.TransferPort, hostChainParams.TransferChannel)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("Error getting client state %s", err))
		return
	}
	timeoutHeight := clienttypes.NewHeight(clientState.GetLatestHeight().GetRevisionNumber(), clientState.GetLatestHeight().GetRevisionHeight()+lscosmostypes.IBCTimeoutHeightIncrement)

	msg := ibctransfertypes.NewMsgTransfer(hostChainParams.TransferPort, hostChainParams.TransferChannel,
		depositBalance, authtypes.NewModuleAddress(lscosmostypes.DelegationModuleAccount).String(),
		delegationState.HostChainDelegationAddress, timeoutHeight, 0)

	handler := k.msgRouter.Handler(msg)

	res, err := handler(ctx, msg)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("could not send transfer msg via MsgServiceRouter, error: %s", err))
		return
	}

	ctx.EventManager().EmitEvents(res.GetEvents())

	// move extra tokens to pstake address - anyone can send tokens to delegation address.
	// should be transferred to pstake address.
	//remainingBalance := allBalances.Sub(sdk.NewCoins(depositBalance))

}

func (k Keeper) RewardEpochEpochWorkFlow(ctx sdk.Context, hostChainParams lscosmostypes.HostChainParams) {
	// send withdraw rewards from delegators.
	delegationState := k.GetDelegationState(ctx)
	if len(delegationState.HostAccountDelegations) == 0 {
		return
	}
	withdrawRewardMsgs := make([]sdk.Msg, len(delegationState.HostAccountDelegations))
	for i, delegation := range delegationState.HostAccountDelegations {
		withdrawRewardMsgs[i] = &distributiontypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: delegationState.HostChainDelegationAddress,
			ValidatorAddress: delegation.ValidatorAddress,
		}
	}
	err := generateAndExecuteICATx(ctx, k, hostChainParams.ConnectionID, lscosmostypes.DelegationAccountPortID, withdrawRewardMsgs, false)
	if err != nil {
		return
	}

	// on Ack do icq for reward acc. balance of uatom
	// callback for sending it to delegation account
	// on Ack delegate txn
}

func (k Keeper) UndelegationEpochWorkFlow(ctx sdk.Context, hostChainParams lscosmostypes.HostChainParams) {
}

// ___________________________________________________________________________________________________

func (k Keeper) OnRecvIBCTransferPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress, transferAck ibcexported.Acknowledgement) {
}

func (k Keeper) OnAcknowledgementIBCTransferPacket(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte, relayer sdk.AccAddress, transferAckErr error) {
	if !k.GetModuleState(ctx) {
		return
	}
	if transferAckErr != nil {
		return
	}
	var ack channeltypes.Acknowledgement
	if err := ibctransfertypes.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return
	}
	if !ack.Success() {
		return
	}
	var data ibctransfertypes.FungibleTokenPacketData
	if err := ibctransfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return
	}
	// check for tokens moved from delegationModuleAccount to it's ica counterpart.
	hostChainParams := k.GetHostChainParams(ctx)
	delegationState := k.GetDelegationState(ctx)
	if packet.GetSourceChannel() != hostChainParams.TransferChannel ||
		packet.GetSourcePort() != hostChainParams.TransferPort {
		// no need to log, since most likely code is expected to enter this condition
		return
	}

	if data.GetSender() != k.GetDelegationModuleAccount(ctx).GetAddress().String() ||
		data.GetReceiver() != delegationState.HostChainDelegationAddress ||
		data.GetDenom() != ibctransfertypes.GetPrefixedDenom(hostChainParams.TransferPort, hostChainParams.TransferChannel, hostChainParams.BaseDenom) {
		return
	}
	k.Logger(ctx).Info(fmt.Sprintf("pstake tokens successfully transferred to host chain address %s, amount: %s, denom: %s", data.Receiver, data.Amount, data.Denom))

	//do ica delegate.
	amount, ok := sdk.NewIntFromString(data.GetAmount())
	if !ok {
		return
	}
	k.AddBalanceToDelegationState(ctx, sdk.NewCoin(hostChainParams.BaseDenom, amount))

	delegatableAmount := k.GetDelegationState(ctx).HostDelegationAccountBalance.AmountOf(hostChainParams.BaseDenom)
	allowlistedValidators := k.GetAllowListedValidators(ctx)
	if !delegatableAmount.GT(sdk.NewInt(int64(len(allowlistedValidators.AllowListedValidators)))) {
		k.Logger(ctx).Info(fmt.Sprintf("amount is too low to delegate, %v ", delegatableAmount))
		return
	}
	msgs := DelegateMsgs(delegationState.HostChainDelegationAddress, allowlistedValidators, delegatableAmount, hostChainParams.BaseDenom)

	err := generateAndExecuteICATx(ctx, k, hostChainParams.ConnectionID, lscosmostypes.DelegationAccountPortID, msgs, false)
	if err != nil {
		return
	}
}

func (k Keeper) OnTimeoutIBCTransferPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress, transferTimeoutErr error) {
	// Do nothing because amount will be reverted to delegationModuleAccount.
}

type IBCTransferHooks struct {
	k Keeper
}

var _ ibchookertypes.IBCHandshakeHooks = IBCTransferHooks{}

func (k Keeper) NewIBCTransferHooks() IBCTransferHooks {
	return IBCTransferHooks{k}
}

func (i IBCTransferHooks) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress, transferAck ibcexported.Acknowledgement) {
	i.k.OnRecvIBCTransferPacket(ctx, packet, relayer, transferAck)
}

func (i IBCTransferHooks) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte, relayer sdk.AccAddress, transferAckErr error) {
	i.k.OnAcknowledgementIBCTransferPacket(ctx, packet, acknowledgement, relayer, transferAckErr)
}

func (i IBCTransferHooks) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress, transferTimeoutErr error) {
	i.k.OnTimeoutIBCTransferPacket(ctx, packet, relayer, transferTimeoutErr)
}
