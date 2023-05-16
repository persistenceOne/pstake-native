package keeper

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	"github.com/gogo/protobuf/proto"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func (k *Keeper) BeginBlock(ctx sdk.Context) {

	// perform BeginBlocker tasks for each chain
	for _, hc := range k.GetAllHostChains(ctx) {
		// attempt to recreate closed ICA channels
		k.DoRecreateICA(ctx, hc)

		// attempt to delegate
		k.DoDelegate(ctx, hc)

		// attempt to automatically claim matured undelegations
		k.DoClaim(ctx, hc)

		// attempt to process any matured unbondings
		k.DoProcessMaturedUndelegations(ctx, hc)

		// attempt to update the validator set if there are any changes
		k.DoUpdateValidatorSet(ctx, hc)
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
		k.DelegateAccountPortOwner(hc.ChainId),
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

func (k *Keeper) DoUpdateValidatorSet(ctx sdk.Context, hc *types.HostChain) {
	consensusState, err := k.GetLatestConsensusState(ctx, hc.ConnectionId)
	if err != nil {
		k.Logger(ctx).Error("could not retrieve client state", "host_chain", hc.ChainId)
		return
	}

	// if the next validator set hash has changes, send an ICQ and update it
	if !bytes.Equal(consensusState.NextValidatorsHash, hc.NextValsetHash) ||
		bytes.Equal(hc.NextValsetHash, []byte{}) {
		k.Logger(ctx).Info(
			"new validator set detected, sending an ICQ to update it.",
			"chain_id",
			hc.ChainId,
		)
		if err = k.QueryHostChainValidators(ctx, hc, stakingtypes.QueryValidatorsRequest{}); err != nil {
			k.Logger(ctx).Error(
				"error sending ICQ for host chain validators",
				"host_chain",
				hc.ChainId,
			)
		}

		// update the validator set next hash
		hc.NextValsetHash = consensusState.NextValidatorsHash
		k.SetHostChain(ctx, hc)
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
	if !k.IsICAChannelActive(ctx, hc, k.GetPortID(k.DelegateAccountPortOwner(hc.ChainId))) &&
		hc.DelegationAccount.ChannelState != types.ICAAccount_ICA_CHANNEL_CREATING {
		if err := k.RegisterICAAccount(ctx, hc.ConnectionId, k.DelegateAccountPortOwner(hc.ChainId)); err != nil {
			k.Logger(ctx).Error("error recreating %s delegate ica: %w", hc.ChainId, err)
		}

		k.Logger(ctx).Info("Recreating delegate ICA.", "chain", hc.ChainId)

		hc.DelegationAccount.ChannelState = types.ICAAccount_ICA_CHANNEL_CREATING
		k.SetHostChain(ctx, hc)
	}

	// if the channel is closed, and it is not being recreated, recreate it
	if !k.IsICAChannelActive(ctx, hc, k.GetPortID(k.RewardsAccountPortOwner(hc.ChainId))) &&
		hc.RewardsAccount.ChannelState != types.ICAAccount_ICA_CHANNEL_CREATING {
		if err := k.RegisterICAAccount(ctx, hc.ConnectionId, k.RewardsAccountPortOwner(hc.ChainId)); err != nil {
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
			return ctx.BlockTime().After(u.MatureTime) && u.State == types.Unbonding_UNBONDING_MATURING
		},
	)

	for _, unbonding := range unbondings {
		channel, found := k.ibcKeeper.ChannelKeeper.GetChannel(ctx, hc.PortId, hc.ChannelId)
		if !found {
			k.Logger(ctx).Error(
				"could not retrieve channel while processing mature undelegations",
				"host_chain",
				hc.ChainId,
			)
			continue
		}

		timeoutHeight := clienttypes.NewHeight(
			clienttypes.GetSelfHeight(ctx).GetRevisionNumber(),
			clienttypes.GetSelfHeight(ctx).GetRevisionHeight()+types.IBCTimeoutHeightIncrement,
		)

		// prepare the msg transfer to bring the undelegation back
		msgTransfer := ibctransfertypes.NewMsgTransfer(
			channel.Counterparty.PortId,
			channel.Counterparty.ChannelId,
			unbonding.UnbondAmount,
			hc.DelegationAccount.Address,
			authtypes.NewModuleAddress(types.UndelegationModuleAccount).String(),
			timeoutHeight,
			0,
			"",
		)

		// execute the transfers
		sequenceID, err := k.GenerateAndExecuteICATx(
			ctx,
			hc.ConnectionId,
			k.DelegateAccountPortOwner(hc.ChainId),
			[]proto.Message{msgTransfer},
		)
		if err != nil {
			k.Logger(ctx).Error(
				"could not send ICA transfer txs",
				"host_chain",
				hc.ChainId,
			)
			continue
		}

		// update the unbonding sequence id and state
		unbonding.IbcSequenceId = sequenceID
		unbonding.State = types.Unbonding_UNBONDING_MATURED
		k.SetUnbonding(ctx, unbonding)
	}
}
