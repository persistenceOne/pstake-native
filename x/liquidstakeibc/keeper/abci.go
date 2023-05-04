package keeper

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func (k *Keeper) BeginBlock(ctx sdk.Context) {

	// perform BeginBlocker tasks for each chain
	for _, hc := range k.GetAllHostChains(ctx) {
		// attempt to delegate
		k.DoDelegate(ctx, hc)

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
	sequenceId, err := k.GenerateAndExecuteICATx(
		ctx,
		hc.ConnectionId,
		hc.ChainId+"."+types.DelegateICAType,
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
		deposit.IbcSequenceId = sequenceId
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
