package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/proto"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func (k *Keeper) BeginBlock(ctx sdk.Context) {

	// perform BeginBlocker tasks for each chain
	for _, hc := range k.GetAllHostChains(ctx) {
		if !hc.Active {
			// don't do anything on inactive chains
			continue
		}

		// attempt to recreate closed ICA channels
		k.DoRecreateICA(ctx, hc)

		// attempt to delegate
		k.DoDelegate(ctx, hc)

		// attempt to automatically claim matured undelegations
		k.DoClaim(ctx, hc)

		// attempt to process any matured unbondings
		k.DoProcessMaturedUndelegations(ctx, hc)

		// attempt to transfer all LSM token deposits
		k.DoTransferLSMTokens(ctx, hc)

		// attempt to redeem all LSM token deposits
		k.DoRedeemLSMTokens(ctx, hc)
	}
}

func (k *Keeper) DoDelegate(ctx sdk.Context, hc *types.HostChain) {
	deposits := k.GetDelegableDepositsForChain(ctx, hc.ChainId)

	// nothing to do if there are no deposits
	if len(deposits) == 0 {
		return
	}

	// get the total amount that can be delegated for that host chain
	totalDepositDelegation := sdk.ZeroInt()
	for _, deposit := range deposits {
		totalDepositDelegation = totalDepositDelegation.Add(deposit.Amount.Amount)
	}

	// generate the delegation messages based on the hc total amount
	messages, err := k.GenerateDelegateMessages(hc, totalDepositDelegation)
	if err != nil {
		k.Logger(ctx).Error(
			"could not generate delegate messages",
			"host_chain",
			hc.ChainId,
		)
		return
	}

	// execute the ICA transactions
	sequenceID, err := k.GenerateAndExecuteICATx(
		ctx,
		hc.ConnectionId,
		hc.DelegationAccount.Owner,
		messages,
	)
	if err != nil {
		k.Logger(ctx).Error(
			"could not send ICA delegate txs",
			"host_chain",
			hc.ChainId,
		)
		return
	}

	// if everything went well, update the deposit states and set the sequence id
	for _, deposit := range deposits {
		deposit.IbcSequenceId = sequenceID
		deposit.State = types.Deposit_DEPOSIT_DELEGATING
		k.SetDeposit(ctx, deposit)
	}
}

func (k *Keeper) DoClaim(ctx sdk.Context, hc *types.HostChain) {
	claimableUnbondings := k.FilterUnbondings(
		ctx,
		func(u types.Unbonding) bool {
			return u.ChainId == hc.ChainId &&
				(u.State == types.Unbonding_UNBONDING_CLAIMABLE || u.State == types.Unbonding_UNBONDING_FAILED)
		},
	)

	for _, unbonding := range claimableUnbondings {
		epochNumber := unbonding.EpochNumber
		userUnbondings := k.FilterUserUnbondings(
			ctx,
			func(u types.UserUnbonding) bool {
				return u.ChainId == hc.ChainId && u.EpochNumber == epochNumber
			},
		)

		for _, userUnbonding := range userUnbondings {
			address, err := sdk.AccAddressFromBech32(userUnbonding.Address)
			if err != nil {
				return
			}

			var claimableCoins sdk.Coins
			switch unbonding.State {
			case types.Unbonding_UNBONDING_CLAIMABLE:
				claimableCoins = sdk.NewCoins(sdk.NewCoin(hc.IBCDenom(), userUnbonding.UnbondAmount.Amount))
				unbonding.UnbondAmount = unbonding.UnbondAmount.Sub(userUnbonding.UnbondAmount)
			case types.Unbonding_UNBONDING_FAILED:
				claimableCoins = sdk.NewCoins(sdk.NewCoin(hc.MintDenom(), userUnbonding.StkAmount.Amount))
				unbonding.BurnAmount = unbonding.BurnAmount.Sub(userUnbonding.StkAmount)
			}

			// send coin to the delegator address from the undelegation module account
			if err = k.bankKeeper.SendCoinsFromModuleToAccount(
				ctx,
				types.UndelegationModuleAccount,
				address,
				claimableCoins,
			); err != nil {
				k.Logger(ctx).Error(
					"could not send unbonded tokens from module account to delegator",
					"host_chain",
					hc.ChainId,
					"epoch",
					userUnbonding.EpochNumber,
				)
				return
			}

			// update the unbonding remaining amount and delete it if it reaches zero
			if unbonding.UnbondAmount.IsZero() || unbonding.BurnAmount.IsZero() {
				k.DeleteUnbonding(ctx, unbonding)
			} else {
				k.SetUnbonding(ctx, unbonding)
			}

			k.DeleteUserUnbonding(ctx, userUnbonding)
		}
	}
}

func (k *Keeper) DoRecreateICA(ctx sdk.Context, hc *types.HostChain) {
	// return early if any of the accounts is currently being recreated
	if (hc.DelegationAccount == nil || hc.RewardsAccount == nil) ||
		(hc.DelegationAccount.ChannelState == types.ICAAccount_ICA_CHANNEL_CREATING ||
			hc.RewardsAccount.ChannelState == types.ICAAccount_ICA_CHANNEL_CREATING) {
		return
	}

	// if the channel is closed, and it is not being recreated, recreate it
	if !k.IsICAChannelActive(ctx, hc, k.GetPortID(hc.DelegationAccount.Owner)) &&
		hc.DelegationAccount.ChannelState != types.ICAAccount_ICA_CHANNEL_CREATING {
		if err := k.RegisterICAAccount(ctx, hc.ConnectionId, hc.DelegationAccount.Owner); err != nil {
			k.Logger(ctx).Error("error recreating %s delegate ica: %w", hc.ChainId, err)
		}

		k.Logger(ctx).Info("Recreating delegate ICA.", "chain", hc.ChainId)

		hc.DelegationAccount.ChannelState = types.ICAAccount_ICA_CHANNEL_CREATING
		k.SetHostChain(ctx, hc)
	}

	// if the channel is closed, and it is not being recreated, recreate it
	if !k.IsICAChannelActive(ctx, hc, k.GetPortID(hc.RewardsAccount.Owner)) &&
		hc.RewardsAccount.ChannelState != types.ICAAccount_ICA_CHANNEL_CREATING {
		if err := k.RegisterICAAccount(ctx, hc.ConnectionId, hc.RewardsAccount.Owner); err != nil {
			k.Logger(ctx).Error("error recreating %s rewards ica: %w", hc.ChainId, err)
		}

		k.Logger(ctx).Info("Recreating rewards ICA.", "chain", hc.ChainId)

		hc.RewardsAccount.ChannelState = types.ICAAccount_ICA_CHANNEL_CREATING
		k.SetHostChain(ctx, hc)
	}
}

func (k *Keeper) DoProcessMaturedUndelegations(ctx sdk.Context, hc *types.HostChain) {
	// get all the unbondings that are matured
	unbondings := k.FilterUnbondings(
		ctx,
		func(u types.Unbonding) bool {
			return u.ChainId == hc.ChainId &&
				ctx.BlockTime().After(u.MatureTime) &&
				u.State == types.Unbonding_UNBONDING_MATURING
		},
	)

	for _, unbonding := range unbondings {
		sequenceID, err := k.SendICATransfer(
			ctx,
			hc,
			unbonding.UnbondAmount,
			hc.DelegationAccount.Address,
			authtypes.NewModuleAddress(types.UndelegationModuleAccount).String(),
			hc.DelegationAccount.Owner,
		)
		if err != nil {
			k.Logger(ctx).Error(
				"Could not process mature undelegations.",
				"host_chain",
				hc.ChainId,
				"error",
				err.Error(),
			)
			continue
		}

		// update the unbonding sequence id and state
		unbonding.IbcSequenceId = sequenceID
		unbonding.State = types.Unbonding_UNBONDING_MATURED
		k.SetUnbonding(ctx, unbonding)
	}

	// get all the validator unbondings that are matured
	validatorUnbondings := k.FilterValidatorUnbondings(
		ctx,
		func(u types.ValidatorUnbonding) bool {
			return u.ChainId == hc.ChainId && u.MatureTime != time.Time{} &&
				ctx.BlockTime().After(u.MatureTime) && u.IbcSequenceId == ""
		},
	)

	for _, validatorUnbonding := range validatorUnbondings {
		sequenceID, err := k.SendICATransfer(
			ctx,
			hc,
			validatorUnbonding.Amount,
			hc.DelegationAccount.Address,
			k.GetDepositModuleAccount(ctx).GetAddress().String(),
			hc.DelegationAccount.Owner,
		)
		if err != nil {
			k.Logger(ctx).Error(
				"Could not process mature validator undelegations.",
				"host_chain",
				hc.ChainId,
				"validator",
				validatorUnbonding.ValidatorAddress,
				"error",
				err.Error(),
			)
			continue
		}

		// update the validator unbonding sequence id and state
		validatorUnbonding.IbcSequenceId = sequenceID
		k.SetValidatorUnbonding(ctx, validatorUnbonding)
	}
}

func (k *Keeper) DoTransferLSMTokens(ctx sdk.Context, hc *types.HostChain) {
	for _, deposit := range k.GetTransferableLSMDeposits(ctx, hc.ChainId) {
		clientState, err := k.GetClientState(ctx, hc.ConnectionId)
		if err != nil {
			// we can't error out here as all the deposits need to be executed
			continue
		}

		timeoutHeight := clienttypes.NewHeight(
			clientState.GetLatestHeight().GetRevisionNumber(),
			clientState.GetLatestHeight().GetRevisionHeight()+types.IBCTimeoutHeightIncrement,
		)

		// craft the IBC message
		msg := ibctransfertypes.NewMsgTransfer(
			ibctransfertypes.PortID,
			hc.ChannelId,
			sdk.NewCoin(deposit.IbcDenom, deposit.Shares.TruncateInt()),
			authtypes.NewModuleAddress(types.DepositModuleAccount).String(),
			hc.DelegationAccount.Address,
			timeoutHeight,
			0,
			"",
		)

		// send the message
		handler := k.msgRouter.Handler(msg)
		res, err := handler(ctx, msg)
		if err != nil {
			k.Logger(ctx).Error(fmt.Sprintf("could not send transfer msg via MsgServiceRouter, error: %s", err))
			// we can't error out here as all the deposits need to be executed
			continue
		}
		ctx.EventManager().EmitEvents(res.GetEvents())

		var msgTransferResponse ibctransfertypes.MsgTransferResponse
		if err = k.cdc.Unmarshal(res.MsgResponses[0].Value, &msgTransferResponse); err != nil {
			// we can't error out here as all the deposits need to be executed
			continue
		}

		// update the deposit state and add the IBC sequence id
		k.UpdateLSMDepositsStateAndSequence(
			ctx,
			[]*types.LSMDeposit{deposit},
			types.LSMDeposit_DEPOSIT_SENT,
			k.GetTransactionSequenceID(hc.ChannelId, msgTransferResponse.Sequence),
		)
	}
}

func (k *Keeper) DoRedeemLSMTokens(ctx sdk.Context, hc *types.HostChain) {
	// generate the ICA messages
	messages := make([]proto.Message, 0)
	deposits := k.GetRedeemableLSMDeposits(ctx, hc.ChainId)
	for _, deposit := range deposits {
		messages = append(
			messages,
			&stakingtypes.MsgRedeemTokensForShares{
				DelegatorAddress: hc.DelegationAccount.Address,
				Amount:           sdk.NewCoin(deposit.Denom, deposit.Shares.TruncateInt()),
			},
		)
	}

	// execute the ICA transaction
	sequenceID, err := k.GenerateAndExecuteICATx(
		ctx,
		hc.ConnectionId,
		hc.DelegationAccount.Owner,
		messages,
	)
	if err != nil {
		k.Logger(ctx).Error("could not send ICA untokenize tx", "host_chain", hc.ChainId)
		return
	}

	// update the deposits state and add the IBC sequence
	k.UpdateLSMDepositsStateAndSequence(
		ctx,
		deposits,
		types.LSMDeposit_DEPOSIT_UNTOKENIZING,
		sequenceID,
	)

	k.Logger(ctx).Info(
		fmt.Sprintf("Redeeming %v deposits.", len(deposits)),
		"host chain",
		hc.ChainId,
		"sequence-id",
		sequenceID,
	)
}
