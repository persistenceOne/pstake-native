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

// BeforeEpochStart - call hook if registered
func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
}

/*
AfterEpochEnd handle the "stake", "reward" and "undelegate" epoch and their respective actions
1. "stake" generates delegate transaction for delegating the amount of stake accumulated over the "stake" epoch
2. "reward" generates delegate transaction for delegating the amount of stake accumulated over the "reward" epochs
and shift the amount to next epoch if the min amount is not reached
3. "undelegate" generated the undelegate transaction for undelegating the amount accumulated over the "undelegate" epoch
*/
func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	//params := k.GetParams(ctx)
	//
	//if epochIdentifier == params.StakingEpochIdentifier {
	//	amount := k.getAmountFromStakingEpoch(ctx, epochNumber)
	//	if !amount.IsZero() {
	//		listOfValidatorsToStake, err := k.FetchValidatorsToDelegate(ctx, amount)
	//		if err != nil {
	//			panic(err)
	//		}
	//		err = k.generateDelegateOutgoingEvent(ctx, listOfValidatorsToStake)
	//		if err != nil {
	//			panic(err)
	//		}
	//	}
	//	k.deleteFromStakingEpoch(ctx, epochNumber)
	//}
	//
	//if epochIdentifier == params.RewardEpochIdentifier {
	//	rewardsToDelegate := k.getFromRewardsInCurrentEpochAmount(ctx, epochNumber)
	//	if !rewardsToDelegate.IsZero() {
	//		listOfValidatorsToStake, err := k.FetchValidatorsToDelegate(ctx, rewardsToDelegate)
	//		if err != nil {
	//			panic(err)
	//		}
	//
	//		// if length of validators to delegate is 0 and amount already is non-zero
	//		// then the rewards can be shifted to next epoch number
	//		if len(listOfValidatorsToStake) == 0 {
	//			k.shiftRewardsToNextEpoch(ctx, epochNumber)
	//		}
	//
	//		// if the list of validator is not empty then generate a list to delegate and mint the rewards
	//		if len(listOfValidatorsToStake) != 0 {
	//			err = k.generateDelegateOutgoingEvent(ctx, listOfValidatorsToStake)
	//			if err != nil {
	//				panic(err)
	//			}
	//
	//			err = k.MintRewardsClaimed(ctx, rewardsToDelegate)
	//			if err != nil {
	//				panic(err)
	//			}
	//		}
	//
	//		// delete rewards from the current epoch as soon as all the surrounding process is complete
	//		k.deleteFromRewardsInCurrentEpoch(ctx, epochNumber)
	//	}
	//
	//}
	//
	//if epochIdentifier == params.UndelegateEpochIdentifier {
	//	withdrawTxns := k.fetchWithdrawTxnsWithCurrentEpochInfo(ctx, epochNumber)
	//	unbondDenom, err := params.GetBondDenomOf(cosmosTypes.DefaultStakingDenom)
	//	if err != nil {
	//		panic(err)
	//	}
	//	totalWithdrawal := k.totalAmountToBeUnbonded(withdrawTxns, unbondDenom)
	//
	//	// calculate uatoms to be unbonded after incorporating C value
	//	cValue := k.GetCValue(ctx)
	//	toBeUnbondedAmount, _ := sdk.NewDecCoinFromDec(totalWithdrawal.Denom, totalWithdrawal.Amount.ToDec().Mul(sdk.NewDec(1).Quo(cValue))).TruncateDecimal()
	//	//check if amount is zero then do not emit event
	//	if !toBeUnbondedAmount.IsZero() {
	//		listOfValidatorsAndUnbondingAmount, err := k.FetchValidatorsToUndelegate(ctx, toBeUnbondedAmount)
	//		if err != nil {
	//			panic(err)
	//		}
	//		k.generateUnbondingOutgoingTxn(ctx, listOfValidatorsAndUnbondingAmount, epochNumber, cValue)
	//
	//		// convert the total withdrawal back to mint denom to burn the same amount as the unbonding transaction has been added to the queue
	//		burnCoin := sdk.NewCoin(k.GetParams(ctx).MintDenom, totalWithdrawal.Amount)
	//
	//		// burn coins
	//		err = k.bankKeeper.BurnCoins(ctx, cosmosTypes.ModuleName, sdk.NewCoins(burnCoin))
	//		if err != nil {
	//			panic(err)
	//		}
	//
	//		// subtract from minted amount as the tokens have been burnt from module
	//		k.SubFromMinted(ctx, burnCoin)
	//
	//		// subtract from staked amount as the transaction has been added to the queue
	//		k.AddToVirtuallyUnbonded(ctx, toBeUnbondedAmount)
	//	}
	//}
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
