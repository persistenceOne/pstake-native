package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

func (k Keeper) GetProportions(ctx sdk.Context, mintedCoin sdk.Coin, ratio sdk.Dec) sdk.Coin {
	return sdk.NewCoin(mintedCoin.Denom, mintedCoin.Amount.ToDec().Mul(ratio).TruncateInt())
}

func (k Keeper) processAllRewardsClaimed(ctx sdk.Context, rewardsAmount sdk.Coin) error {
	// get amount in Stk assets form
	params := k.GetParams(ctx)
	rewardAmountInUSTK := sdk.NewCoin(params.MintDenom, rewardsAmount.Amount)

	// get distribution proportions for minting stk assets
	distributionProportion := params.DistributionProportion
	totalDistributionProportion := distributionProportion.ValidatorRewards.Add(distributionProportion.DeveloperRewards)
	totalRewards := k.GetProportions(ctx, rewardAmountInUSTK, totalDistributionProportion)

	// calculate rewards for developers and validators
	validatorRewards := k.GetProportions(ctx, totalRewards, distributionProportion.ValidatorRewards)
	developerRewards := k.GetProportions(ctx, totalRewards, distributionProportion.DeveloperRewards)

	for _, wallet := range k.getAllOracleValidatorSet(ctx) {
		amount := k.GetProportions(ctx, validatorRewards, wallet.Weight)
		accAddress, err := cosmosTypes.AccAddressFromBech32(wallet.Address, "persistencevaloper")
		if err != nil {
			return err
		}
		err = k.mintTokensForRewardReceivers(ctx, accAddress, amount)
		if err != nil {
			return err
		}
	}

	for _, wallet := range params.WeightedDeveloperRewardsReceivers {
		amount := k.GetProportions(ctx, developerRewards, wallet.Weight)
		accAddress, err := sdk.AccAddressFromBech32(wallet.Address)
		if err != nil {
			return err
		}
		err = k.mintTokensForRewardReceivers(ctx, accAddress, amount)
		if err != nil {
			return err
		}
	}

	// add to virtually staked amount
	k.AddToVirtuallyStakedAmount(ctx, rewardsAmount)

	return nil
}
