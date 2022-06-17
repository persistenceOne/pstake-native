package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	epochsTypes "github.com/persistenceOne/pstake-native/x/epochs/types"
)

// Implements GovHooks interface
var _ cosmosTypes.GovHooks = Keeper{}

// AfterProposalSubmission - call hook if registered
func (k Keeper) AfterProposalSubmission(ctx sdk.Context, proposalID uint64) {
	if k.hooks != nil {
		k.hooks.AfterProposalSubmission(ctx, proposalID)
	}
}

// AfterProposalVote - call hook if registered
func (k Keeper) AfterProposalVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) {
	if k.hooks != nil {
		k.hooks.AfterProposalVote(ctx, proposalID, voterAddr)
	}
}

// AfterProposalVotingPeriodEnded - call hook if registered
func (k Keeper) AfterProposalVotingPeriodEnded(ctx sdk.Context, proposalID uint64) {
	if k.hooks != nil {
		k.hooks.AfterProposalVotingPeriodEnded(ctx, proposalID)
	}
}

func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
}

func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	params := k.GetParams(ctx)

	if epochIdentifier == params.StakingEpochIdentifier {
		amount := k.getAmountFromStakingEpoch(ctx, epochNumber)
		if !amount.IsZero() {
			listOfValidatorsToStake, err := k.FetchValidatorsToDelegate(ctx, amount, epochIdentifier)
			if err != nil {
				panic(err)
			}
			err = k.generateDelegateOutgoingEvent(ctx, listOfValidatorsToStake, epochNumber)
			if err != nil {
				panic(err)
			}
		}
		k.deleteFromStakingEpoch(ctx, epochNumber)
	}

	if epochIdentifier == params.RewardEpochIdentifier {
		rewardsToDelegate := k.getFromRewardsInCurrentEpochAmount(ctx, epochNumber)
		if !rewardsToDelegate.IsZero() {
			listOfValidatorsToStake, err := k.FetchValidatorsToDelegate(ctx, rewardsToDelegate, epochIdentifier)
			if err != nil {
				panic(err)
			}
			err = k.generateDelegateOutgoingEvent(ctx, listOfValidatorsToStake, epochNumber)
			if err != nil {
				panic(err)
			}

			err = k.processAllRewardsClaimed(ctx, rewardsToDelegate)
			if err != nil {
				panic(err)
			}
			k.deleteFromRewardsInCurrentEpoch(ctx, epochNumber)
		}

	}

	if epochIdentifier == params.UndelegateEpochIdentifier {
		withdrawTxns, err := k.fetchWithdrawTxnsWithCurrentEpochInfo(ctx, epochNumber)
		if err != nil {
			panic(err)
		}
		unbondDenom, err := params.GetBondDenomOf("uatom")
		if err != nil {
			panic(err)
		}
		totalWithdrawal := k.totalAmountToBeUnbonded(withdrawTxns, unbondDenom)

		// calculate uatoms to be unbonded after incorporating C value
		toBeUnbondedAmount, _ := sdk.NewDecCoinFromDec(totalWithdrawal.Denom, totalWithdrawal.Amount.ToDec().Mul(sdk.NewDec(1).Quo(k.GetCValue(ctx)))).TruncateDecimal()
		//check if amount is zero then do not emit event
		if !toBeUnbondedAmount.IsZero() {
			listOfValidatorsAndUnbondingAmount, err := k.FetchValidatorsToUndelegate(ctx, toBeUnbondedAmount)
			if err != nil {
				panic(err)
			}
			k.generateUnbondingOutgoingEvent(ctx, listOfValidatorsAndUnbondingAmount, epochNumber)

			// convert the total withdrawal back to mint denom to burn the same amount as the unbonding transaction has been added to the queue
			burnCoin := sdk.NewCoin(k.GetParams(ctx).MintDenom, totalWithdrawal.Amount)

			// burn coins
			err = k.bankKeeper.BurnCoins(ctx, cosmosTypes.ModuleName, sdk.NewCoins(burnCoin))
			if err != nil {
				panic(err)
			}

			// subtract from minted amount as the tokens have been burnt from module
			k.SubFromMintedAmount(ctx, burnCoin)

			// todo : check if it is right place to sub from staked amount or after MsgUndelegate success in batch.go by introducing a virtually unstaked amount
			// subtract from staked amount as the transaction has been added to the queue
			k.SubFromStakedAmount(ctx, toBeUnbondedAmount)

			// delete the entry for withdrawal from current epoch info once processed
			k.deleteWithdrawTxnWithCurrentEpochInfo(ctx, epochNumber)
		}
	}
}

// ___________________________________________________________________________________________________

// RewardsHooks wrapper struct for incentives keeper
type RewardsHooks struct {
	k Keeper
}

var _ epochsTypes.EpochHooks = RewardsHooks{}

// Return the wrapper struct
func (k Keeper) Hooks() RewardsHooks {
	return RewardsHooks{k}
}

// epochs hooks
func (h RewardsHooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.BeforeEpochStart(ctx, epochIdentifier, epochNumber)
}

func (h RewardsHooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.AfterEpochEnd(ctx, epochIdentifier, epochNumber)
}
