package keeper

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/proto"

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

		// attempt to redeem LSM tokens
		if hc.Flags.Lsm {
			k.DoRedeemLSMTokens(ctx, hc)
		}
	}
	// delete acked redelegations
	k.DoDeleteRedelegationTxs(ctx)
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

		// emit the delegation event for every deposit
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeDoDelegationDeposit,
				sdk.NewAttribute(types.AttributeChainID, hc.ChainId),
				sdk.NewAttribute(types.AttributeEpoch, strconv.FormatInt(deposit.Epoch, 10)),
				sdk.NewAttribute(types.AttributeDelegatedAmount, sdk.NewCoin(hc.HostDenom, deposit.Amount.Amount).String()),
				sdk.NewAttribute(types.AttributeIBCSequenceID, sequenceID),
			),
		)
	}

	// emit the delegation event
	encMsgs, _ := json.Marshal(&messages)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDoDelegation,
			sdk.NewAttribute(types.AttributeChainID, hc.ChainId),
			sdk.NewAttribute(types.AttributeTotalDelegatedAmount, sdk.NewCoin(hc.HostDenom, totalDepositDelegation).String()),
			sdk.NewAttribute(types.AttributeICAMessages, base64.StdEncoding.EncodeToString(encMsgs)),
		),
	)
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
				k.Logger(ctx).Error(
					"could not send unbonded tokens from module account to delegator",
					"host_chain",
					hc.ChainId,
					"epoch",
					userUnbonding.EpochNumber,
				)

				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						types.EventFailedClaimUnbondings,
						sdk.NewAttribute(types.AttributeChainID, hc.ChainId),
						sdk.NewAttribute(types.AttributeEpoch, strconv.FormatInt(epochNumber, 10)),
						sdk.NewAttribute(types.AttributeClaimAddress, userUnbonding.Address),
						sdk.NewAttribute(types.AttributeClaimStatus, unbonding.State.String()),
					),
				)

				continue
			}

			var claimableCoins sdk.Coins
			var eventAmount sdk.Coin // used for claim events
			switch unbonding.State {
			case types.Unbonding_UNBONDING_CLAIMABLE:
				claimableCoins = sdk.NewCoins(sdk.NewCoin(hc.IBCDenom(), userUnbonding.UnbondAmount.Amount))
				eventAmount = sdk.NewCoin(hc.HostDenom, userUnbonding.UnbondAmount.Amount)
				unbonding.UnbondAmount = unbonding.UnbondAmount.Sub(userUnbonding.UnbondAmount)
			case types.Unbonding_UNBONDING_FAILED:
				claimableCoins = sdk.NewCoins(sdk.NewCoin(hc.MintDenom(), userUnbonding.StkAmount.Amount))
				eventAmount = sdk.NewCoin(hc.MintDenom(), userUnbonding.StkAmount.Amount)
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

				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						types.EventFailedClaimUnbondings,
						sdk.NewAttribute(types.AttributeChainID, hc.ChainId),
						sdk.NewAttribute(types.AttributeEpoch, strconv.FormatInt(epochNumber, 10)),
						sdk.NewAttribute(types.AttributeClaimAmount, eventAmount.String()),
						sdk.NewAttribute(types.AttributeClaimAddress, userUnbonding.Address),
						sdk.NewAttribute(types.AttributeClaimStatus, unbonding.State.String()),
					),
				)

				continue
			}

			// update the unbonding remaining amount and delete it if it reaches zero
			if unbonding.UnbondAmount.IsZero() || unbonding.BurnAmount.IsZero() {
				k.DeleteUnbonding(ctx, unbonding)
			} else {
				k.SetUnbonding(ctx, unbonding)
			}

			k.DeleteUserUnbonding(ctx, userUnbonding)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeClaimedUnbondings,
					sdk.NewAttribute(types.AttributeChainID, hc.ChainId),
					sdk.NewAttribute(types.AttributeEpoch, strconv.FormatInt(epochNumber, 10)),
					sdk.NewAttribute(types.AttributeClaimAmount, eventAmount.String()),
					sdk.NewAttribute(types.AttributeClaimAddress, userUnbonding.Address),
				),
			)
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
		} else {
			k.Logger(ctx).Info("Recreating delegate ICA.", "chain", hc.ChainId)

			hc.DelegationAccount.ChannelState = types.ICAAccount_ICA_CHANNEL_CREATING
			k.SetHostChain(ctx, hc)
		}
	}

	// if the channel is closed, and it is not being recreated, recreate it
	if !k.IsICAChannelActive(ctx, hc, k.GetPortID(hc.RewardsAccount.Owner)) &&
		hc.RewardsAccount.ChannelState != types.ICAAccount_ICA_CHANNEL_CREATING {
		if err := k.RegisterICAAccount(ctx, hc.ConnectionId, hc.RewardsAccount.Owner); err != nil {
			k.Logger(ctx).Error("error recreating %s rewards ica: %w", hc.ChainId, err)
		} else {
			k.Logger(ctx).Info("Recreating rewards ICA.", "chain", hc.ChainId)

			hc.RewardsAccount.ChannelState = types.ICAAccount_ICA_CHANNEL_CREATING
			k.SetHostChain(ctx, hc)
		}
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

func (k *Keeper) DoRedeemLSMTokens(ctx sdk.Context, hc *types.HostChain) {
	deposits := k.GetRedeemableLSMDeposits(ctx, hc.ChainId)

	// generate the ICA messages
	messagesChunks := make([][]proto.Message, 0)
	depositsChunks := make([][]*types.LSMDeposit, 0)
	for i := 0; i < len(deposits); i += types.ICAMessagesChunkSize {
		end := i + types.ICAMessagesChunkSize

		// avoid slicing past the deposits length
		if end > len(deposits) {
			end = len(deposits)
		}

		// create a redeem message for each deposit in the current chunk
		depositsChunk := deposits[i:end]
		messagesChunk := make([]proto.Message, 0)
		for _, deposit := range depositsChunk {
			messagesChunk = append(
				messagesChunk,
				&stakingtypes.MsgRedeemTokensForShares{
					DelegatorAddress: hc.DelegationAccount.Address,
					Amount:           sdk.NewCoin(deposit.Denom, deposit.Shares.TruncateInt()),
				},
			)
		}

		// save both chunks in the respective arrays
		messagesChunks = append(messagesChunks, messagesChunk)
		depositsChunks = append(depositsChunks, depositsChunk)
	}

	for i, messagesChunk := range messagesChunks {
		if len(messagesChunk) == 0 {
			continue
		}

		// execute the ICA transaction
		sequenceID, err := k.GenerateAndExecuteICATx(
			ctx,
			hc.ConnectionId,
			hc.DelegationAccount.Owner,
			messagesChunk,
		)
		if err != nil {
			k.Logger(ctx).Error("could not send ICA untokenize tx", "host_chain", hc.ChainId)
			return
		}

		// update the deposits state and add the IBC sequence
		k.UpdateLSMDepositsStateAndSequence(
			ctx,
			depositsChunks[i],
			types.LSMDeposit_DEPOSIT_UNTOKENIZING,
			sequenceID,
		)

		k.Logger(ctx).Info(
			fmt.Sprintf("Redeeming %v deposits.", len(depositsChunks[i])),
			"host chain",
			hc.ChainId,
			"sequence-id",
			sequenceID,
		)

		// emit the untokenize event
		encMsgs, _ := json.Marshal(&messagesChunk)
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeRedeemTokensForShares,
				sdk.NewAttribute(types.AttributeChainID, hc.ChainId),
				sdk.NewAttribute(types.AttributeICAMessages, base64.StdEncoding.EncodeToString(encMsgs)),
				sdk.NewAttribute(types.AttributeIBCSequenceID, sequenceID),
			),
		)
	}
}

func (k *Keeper) DoDeleteRedelegationTxs(ctx sdk.Context) {
	redelegationTxs := k.GetAllRedelegationTx(ctx)
	for _, redelegationTx := range redelegationTxs {
		if redelegationTx.State == types.RedelegateTx_REDELEGATE_ACKED {
			k.DeleteRedelegationTx(ctx, redelegationTx.ChainId, redelegationTx.IbcSequenceId)
		}
	}
}

func (k *Keeper) DoDeleteMaturedRedelegation(ctx sdk.Context, hc *types.HostChain) {
	redelegations, ok := k.GetRedelegations(ctx, hc.ChainId)
	if !ok {
		// no redelegations, return early
		return
	}
	for i := range redelegations.Redelegations {
		var entries []stakingtypes.RedelegationEntry
		for j := range redelegations.Redelegations[i].Entries {
			if !redelegations.Redelegations[i].Entries[j].IsMature(ctx.BlockTime()) {
				entries = append(entries, redelegations.Redelegations[i].Entries[j])
			}
		}
		redelegations.Redelegations[i].Entries = entries
	}
	k.SetRedelegations(ctx, hc.ChainId, redelegations.Redelegations)
}
