package keeper

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icatypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func (k *Keeper) BeginBlock(ctx sdk.Context) {

	// perform BeginBlocker tasks for each chain
	for _, hc := range k.GetAllHostChains(ctx) {
		// attempt to recreate closed ICA channels
		k.DoRecreateICA(ctx, hc)

		// attempt to delegate
		k.DoDelegate(ctx, hc)

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

	_, isDelegateActive := k.icaControllerKeeper.GetOpenActiveChannel(
		ctx,
		hc.ConnectionId,
		icatypes.ControllerPortPrefix+k.DelegateAccountPortOwner(hc.ChainId),
	)
	// if the channel is closed, and it is not being recreated, recreate it
	if !isDelegateActive && hc.DelegationAccount.ChannelState != types.ICAAccount_ICA_CHANNEL_CREATING {
		if err := k.RegisterICAAccount(ctx, hc.ConnectionId, k.DelegateAccountPortOwner(hc.ChainId)); err != nil {
			k.Logger(ctx).Error("error recreating %s delegate ica: %w", hc.ChainId, err)
		}
		k.Logger(ctx).Info("Recreating delegate ICA.", "chain", hc.ChainId)
		hc.DelegationAccount.ChannelState = types.ICAAccount_ICA_CHANNEL_CREATING
		k.SetHostChain(ctx, hc)
	}

	_, isRewardsActive := k.icaControllerKeeper.GetOpenActiveChannel(
		ctx,
		hc.ConnectionId,
		icatypes.ControllerPortPrefix+k.RewardsAccountPortOwner(hc.ChainId),
	)
	// if the channel is closed, and it is not being recreated, recreate it
	if !isRewardsActive && hc.RewardsAccount.ChannelState != types.ICAAccount_ICA_CHANNEL_CREATING {
		if err := k.RegisterICAAccount(ctx, hc.ConnectionId, k.RewardsAccountPortOwner(hc.ChainId)); err != nil {
			k.Logger(ctx).Error("error recreating %s rewards ica: %w", hc.ChainId, err)
		}
		k.Logger(ctx).Info("Recreating rewards ICA.", "chain", hc.ChainId)
		hc.RewardsAccount.ChannelState = types.ICAAccount_ICA_CHANNEL_CREATING
		k.SetHostChain(ctx, hc)
	}
}

func (k *Keeper) DoProcessMaturedUndelegations(ctx sdk.Context, hc *types.HostChain) {}
