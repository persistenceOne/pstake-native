package keeper

import (
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	epochsTypes "github.com/persistenceOne/pstake-native/x/epochs/types"
)

// Implements GovHooks interface
var _ cosmosTypes.GovHooks = Keeper{}

// AfterProposalSubmission - call hook if registered
func (keeper Keeper) AfterProposalSubmission(ctx sdk.Context, proposalID uint64) {
	if keeper.hooks != nil {
		keeper.hooks.AfterProposalSubmission(ctx, proposalID)
	}
}

// AfterProposalVote - call hook if registered
func (keeper Keeper) AfterProposalVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) {
	if keeper.hooks != nil {
		keeper.hooks.AfterProposalVote(ctx, proposalID, voterAddr)
	}
}

// AfterProposalVotingPeriodEnded - call hook if registered
func (keeper Keeper) AfterProposalVotingPeriodEnded(ctx sdk.Context, proposalID uint64) {
	if keeper.hooks != nil {
		keeper.hooks.AfterProposalVotingPeriodEnded(ctx, proposalID)
	}
}

func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
}

func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	params := k.GetParams(ctx)

	if epochIdentifier == params.StakingEpochIdentifier {
		listOfMintTxns, err := k.getFromStakingEpoch(ctx, epochNumber)
		stakingDenom, err := params.GetBondDenomOf("uatom")
		if err != nil {
			panic(err)
		}

		rewardsToBeClaimed, err := k.getFromRewardsInCurrentEpochAmount(ctx, epochNumber)
		if err != nil {
			panic(err)
		}

		amt := getTotalStakingAmount(listOfMintTxns, stakingDenom)
		amt.Add(rewardsToBeClaimed)

		if !amt.IsZero() {
			listOfValidatorsToStake := k.fetchValidatorsToDelegate(ctx, amt)
			err = k.generateDelegateOutgoingEvent(ctx, listOfValidatorsToStake, epochNumber)
			if err != nil {
				panic(err)
			}
		}

		//cosmosValidators := params.ValidatorSetCosmosChain
		////TODO : Check if some amount has been delegated on cosmos chain or not. If there is then claim event is generated.
		//var withdrawMessages []*codecTypes.Any
		//for _, validator := range cosmosValidators {
		//	msg := distrTypes.MsgWithdrawDelegatorReward{
		//		DelegatorAddress: params.CustodialAddress,
		//		ValidatorAddress: validator.Address,
		//	}
		//	anyMsg, err := codecTypes.NewAnyWithValue(&msg)
		//	if err != nil {
		//		panic(err)
		//	}
		//	withdrawMessages = append(withdrawMessages, anyMsg)
		//}
		//chuckMsgs := ChunkSlice(withdrawMessages, params.ChunkSize)
		//for _, chunk := range chuckMsgs {
		//	k.generateWithdrawRewardsEvent(ctx, chunk)
		//}
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
		amount := k.totalAmountToBeUnbonded(withdrawTxns, unbondDenom)
		//check if amount is zero then do not emit event
		if !amount.IsZero() {
			listOfValidatorsAndUnbondingAmount := k.fetchValidatorsToUndelegate(ctx, amount)
			k.generateUnbondingOutgoingEvent(ctx, listOfValidatorsAndUnbondingAmount, epochNumber)
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

func ChunkSlice(slice []*codecTypes.Any, chunkSize int64) (chunks [][]*codecTypes.Any) {
	for {
		if len(slice) == 0 {
			break
		}

		// necessary check to avoid slicing beyond
		// slice capacity
		if int64(len(slice)) < chunkSize {
			chunkSize = int64(len(slice))
		}

		chunks = append(chunks, slice[0:chunkSize])
		slice = slice[chunkSize:]
	}

	return chunks
}
