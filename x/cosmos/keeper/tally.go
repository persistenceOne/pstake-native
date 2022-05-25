package keeper

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

func (k Keeper) Tally(ctx sdkTypes.Context, proposal cosmosTypes.Proposal) (passes bool, tallyResult map[cosmosTypes.VoteOption]sdkTypes.Dec) {
	results := make(map[cosmosTypes.VoteOption]sdkTypes.Dec)
	results[cosmosTypes.OptionYes] = sdkTypes.ZeroDec()
	results[cosmosTypes.OptionAbstain] = sdkTypes.ZeroDec()
	results[cosmosTypes.OptionNo] = sdkTypes.ZeroDec()
	results[cosmosTypes.OptionNoWithVeto] = sdkTypes.ZeroDec()

	currValidators := make(map[string]govTypes.ValidatorGovInfo)

	// fetch all the bonded validators, insert them into currValidators
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx, func(index int64, validator stakingTypes.ValidatorI) (stop bool) {
		currValidators[validator.GetOperator().String()] = govTypes.NewValidatorGovInfo(
			validator.GetOperator(),
			validator.GetBondedTokens(),
			validator.GetDelegatorShares(),
			sdkTypes.ZeroDec(),
			govTypes.WeightedVoteOptions{},
		)

		return false
	})

	k.IterateVotes(ctx, proposal.ProposalId, func(vote cosmosTypes.Vote) bool {
		voter, err := sdkTypes.AccAddressFromBech32(vote.Voter)
		if err != nil {
			panic(err)
		}

		val, _, err := k.getAllValidatorOrchestratorMappingAndFindIfExist(ctx, voter)
		if err != nil {
			panic(err)
		}

		valAddress, found := k.GetValidatorOrchestrator(ctx, val)
		if valAddress == nil {
			panic("unauthorized vote present in db")
		}

		valAddressString := string(valAddress.Bytes())
		if found {
			votingPower := currValidators[valAddressString].DelegatorShares.
				MulInt(currValidators[valAddressString].BondedTokens).
				Quo(currValidators[valAddressString].DelegatorShares)
			for _, option := range vote.Options {
				subPower := votingPower.Mul(option.Weight)
				results[option.Option] = results[option.Option].Add(subPower)
			}
		}
		k.deleteVote(ctx, proposal.ProposalId, voter)
		return false
	})

	return true, results
}
