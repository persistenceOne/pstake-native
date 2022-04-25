package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

func (k Keeper) GetProportions(ctx sdk.Context, mintedCoin sdk.Coin, ratio sdk.Dec) sdk.Coin {
	return sdk.NewCoin(mintedCoin.Denom, mintedCoin.Amount.ToDec().Mul(ratio).TruncateInt())
}

func (k Keeper) processAllRewardsClaimed(ctx sdk.Context, rewardsAmount sdk.Coin) error {
	params := k.GetParams(ctx)
	rewardAmountInUSTK := sdk.NewCoin(params.MintDenom, rewardsAmount.Amount)
	distributionProportion := params.DistributionProportion
	totalDistributionProportion := distributionProportion.ValidatorRewards.Add(distributionProportion.DeveloperRewards)
	totalRewards := k.GetProportions(ctx, rewardAmountInUSTK, totalDistributionProportion)

	validatorRewards := k.GetProportions(ctx, totalRewards, distributionProportion.ValidatorRewards)
	developerRewards := k.GetProportions(ctx, totalRewards, distributionProportion.DeveloperRewards)

	for _, wallet := range params.ValidatorSetNativeChain {
		amount := sdk.NewCoins(k.GetProportions(ctx, validatorRewards, wallet.Weight))
		err := k.mintTokensForRewardReceivers(ctx, wallet.Address, amount)
		if err != nil {
			return err
		}
	}

	for _, wallet := range params.WeightedDeveloperRewardsReceivers {
		amount := sdk.NewCoins(k.GetProportions(ctx, developerRewards, wallet.Weight))
		err := k.mintTokensForRewardReceivers(ctx, wallet.Address, amount)
		if err != nil {
			return err
		}
	}

	//TODO : update c ratio
	return nil
}
