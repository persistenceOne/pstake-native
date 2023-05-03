package keeper

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/persistenceOne/persistence-sdk/v2/utils"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func (k *Keeper) BeginBlock(ctx sdk.Context) {
	err := utils.ApplyFuncIfNoError(ctx, k.DoDelegate)
	if err != nil {
		k.Logger(ctx).Error("Unable to Delegate tokens", "err: ", err)
	}

	for _, hc := range k.GetAllHostChains(ctx) {
		consensusState, err := k.GetLatestConsensusState(ctx, hc.ConnectionId)
		if err != nil {
			k.Logger(ctx).Error("could not retrieve client state", "host_chain", hc.ChainId)
			return
		}

		// if the next validator set hash has changes, send an ICQ and update it
		if !bytes.Equal(consensusState.NextValidatorsHash, hc.NextValsetHash) ||
			bytes.Equal(hc.NextValsetHash, []byte{}) {
			k.Logger(ctx).Info(
				"New validator set detected, sending an ICQ to update it.",
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
}

func (k *Keeper) DoDelegate(ctx sdk.Context) error {
	hostChains := k.GetAllHostChains(ctx)

	// create and execute MsgDelegation txs for each host chain
	for _, hc := range hostChains {
		deposits := k.GetDelegableDepositsForChain(ctx, hc.ChainId)

		// nothing to do if there are no deposits
		if len(deposits) == 0 {
			continue
		}

		// get the total amount that can be delegated for that host chain
		totalDepositDelegation := sdk.ZeroInt()
		for _, deposit := range deposits {
			totalDepositDelegation = totalDepositDelegation.Add(deposit.Amount.Amount)
		}

		// generate the delegation messages based on the hc total amount
		messages, err := k.GenerateDelegateMessages(hc, totalDepositDelegation)
		if err != nil {
			return err
		}

		// execute the ICA transactions
		sequenceId, err := k.GenerateAndExecuteICATx(
			ctx,
			hc.ConnectionId,
			hc.ChainId+"."+types.DelegateICAType,
			messages,
		)
		if err != nil {
			return err
		}

		// if everything went well, update the deposit states and set the sequence id
		for _, deposit := range deposits {
			deposit.IbcSequenceId = sequenceId
			deposit.State = types.Deposit_DEPOSIT_DELEGATING
			k.SetDeposit(ctx, deposit)
		}
	}

	return nil
}
