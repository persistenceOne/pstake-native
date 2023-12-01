package keeper

import (
	"strconv"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstake/types"
)

func (k Keeper) GetProxyAccBalance(ctx sdk.Context, proxyAcc sdk.AccAddress) (balance sdk.Coin) {
	bondDenom := k.stakingKeeper.BondDenom(ctx)
	return sdk.NewCoin(bondDenom, k.bankKeeper.SpendableCoins(ctx, proxyAcc).AmountOf(bondDenom))
}

// TryRedelegation attempts redelegation, which is applied only when successful through cached context because there is a constraint that fails if already receiving redelegation.
func (k Keeper) TryRedelegation(ctx sdk.Context, re types.Redelegation) (completionTime time.Time, err error) {
	dstVal := re.DstValidator.GetOperator()
	srcVal := re.SrcValidator.GetOperator()

	// check the source validator already has receiving transitive redelegation
	hasReceiving := k.stakingKeeper.HasReceivingRedelegation(ctx, re.Delegator, srcVal)
	if hasReceiving {
		return time.Time{}, stakingtypes.ErrTransitiveRedelegation
	}

	// calculate delShares from tokens with validation
	shares, err := k.stakingKeeper.ValidateUnbondAmount(
		ctx, re.Delegator, srcVal, re.Amount,
	)
	if err != nil {
		return time.Time{}, err
	}

	// when last, full redelegation of shares from delegation
	if re.Last {
		shares = re.SrcValidator.GetDelShares(ctx, k.stakingKeeper)
	}
	cachedCtx, writeCache := ctx.CacheContext()
	completionTime, err = k.stakingKeeper.BeginRedelegation(cachedCtx, re.Delegator, srcVal, dstVal, shares)
	if err != nil {
		return time.Time{}, err
	}
	writeCache()
	return completionTime, nil
}

// Rebalance argument liquidVals containing ValidatorStatusActive which is containing just added on whitelist(liquidToken 0) and ValidatorStatusInactive to delist
func (k Keeper) Rebalance(ctx sdk.Context, proxyAcc sdk.AccAddress, liquidVals types.LiquidValidators, whitelistedValsMap types.WhitelistedValsMap, rebalancingTrigger sdk.Dec) (redelegations []types.Redelegation) {
	logger := k.Logger(ctx)
	totalLiquidTokens, liquidTokenMap := liquidVals.TotalLiquidTokens(ctx, k.stakingKeeper, false)
	if !totalLiquidTokens.IsPositive() {
		return []types.Redelegation{}
	}

	weightMap, totalWeight := k.GetWeightMap(ctx, liquidVals, whitelistedValsMap)

	// no active liquid validators
	if !totalWeight.IsPositive() {
		return []types.Redelegation{}
	}

	// calculate rebalancing target map
	targetMap := map[string]math.Int{}
	totalTargetMap := sdk.ZeroInt()
	for _, val := range liquidVals {
		targetMap[val.OperatorAddress] = totalLiquidTokens.Mul(weightMap[val.OperatorAddress]).Quo(totalWeight)
		totalTargetMap = totalTargetMap.Add(targetMap[val.OperatorAddress])
	}
	crumb := totalLiquidTokens.Sub(totalTargetMap)
	if !totalTargetMap.IsPositive() {
		return []types.Redelegation{}
	}
	// crumb to first non zero liquid validator
	for _, val := range liquidVals {
		if targetMap[val.OperatorAddress].IsPositive() {
			targetMap[val.OperatorAddress] = targetMap[val.OperatorAddress].Add(crumb)
			break
		}
	}

	failCount := 0
	rebalancingThresholdAmt := rebalancingTrigger.Mul(sdk.NewDecFromInt(totalLiquidTokens)).TruncateInt()

	for i := 0; i < liquidVals.Len(); i++ {
		// get min, max of liquid token gap
		minVal, maxVal, amountNeeded, last := liquidVals.MinMaxGap(targetMap, liquidTokenMap)
		if amountNeeded.IsZero() || (i == 0 && !amountNeeded.GT(rebalancingThresholdAmt)) {
			break
		}

		// sync liquidTokenMap applied rebalancing
		liquidTokenMap[maxVal.OperatorAddress] = liquidTokenMap[maxVal.OperatorAddress].Sub(amountNeeded)
		liquidTokenMap[minVal.OperatorAddress] = liquidTokenMap[minVal.OperatorAddress].Add(amountNeeded)

		// try redelegation from max validator to min validator
		redelegation := types.Redelegation{
			Delegator:    proxyAcc,
			SrcValidator: maxVal,
			DstValidator: minVal,
			Amount:       amountNeeded,
			Last:         last,
		}
		_, err := k.TryRedelegation(ctx, redelegation)
		if err != nil {
			redelegation.Error = err
			failCount++
		}
		redelegations = append(redelegations, redelegation)
	}
	if failCount > 0 {
		logger.Error("rebalancing failed due to redelegation hopping", "redelegations", redelegations)
	}
	if len(redelegations) != 0 {
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeBeginRebalancing,
				sdk.NewAttribute(types.AttributeKeyDelegator, types.LiquidStakeProxyAcc.String()),
				sdk.NewAttribute(types.AttributeKeyRedelegationCount, strconv.Itoa(len(redelegations))),
				sdk.NewAttribute(types.AttributeKeyRedelegationFailCount, strconv.Itoa(failCount)),
			),
		})
		logger.Info(types.EventTypeBeginRebalancing,
			types.AttributeKeyDelegator, types.LiquidStakeProxyAcc.String(),
			types.AttributeKeyRedelegationCount, strconv.Itoa(len(redelegations)),
			types.AttributeKeyRedelegationFailCount, strconv.Itoa(failCount))
	}
	return redelegations
}

func (k Keeper) UpdateLiquidValidatorSet(ctx sdk.Context) []types.Redelegation {
	logger := k.Logger(ctx)
	params := k.GetParams(ctx)
	liquidValidators := k.GetAllLiquidValidators(ctx)
	liquidValsMap := liquidValidators.Map()
	whitelistedValsMap := types.GetWhitelistedValsMap(params.WhitelistedValidators)

	// Set Liquid validators for added whitelist validators
	for _, wv := range params.WhitelistedValidators {
		if _, ok := liquidValsMap[wv.ValidatorAddress]; !ok {
			lv := types.LiquidValidator{
				OperatorAddress: wv.ValidatorAddress,
			}
			if k.IsActiveLiquidValidator(ctx, lv, whitelistedValsMap) {
				k.SetLiquidValidator(ctx, lv)
				liquidValidators = append(liquidValidators, lv)
				ctx.EventManager().EmitEvents(sdk.Events{
					sdk.NewEvent(
						types.EventTypeAddLiquidValidator,
						sdk.NewAttribute(types.AttributeKeyLiquidValidator, lv.OperatorAddress),
					),
				})
				logger.Info(types.EventTypeAddLiquidValidator, types.AttributeKeyLiquidValidator, lv.OperatorAddress)
			}
		}
	}

	// rebalancing based updated liquid validators status with threshold, try by cachedCtx
	// tombstone status also handled on Rebalance
	reds := k.Rebalance(ctx, types.LiquidStakeProxyAcc, liquidValidators, whitelistedValsMap, types.RebalancingTrigger)

	// unbond all delShares to proxyAcc if delShares exist on inactive liquid validators
	for _, lv := range liquidValidators {
		if !k.IsActiveLiquidValidator(ctx, lv, whitelistedValsMap) {
			delShares := lv.GetDelShares(ctx, k.stakingKeeper)
			if delShares.IsPositive() {
				cachedCtx, writeCache := ctx.CacheContext()
				completionTime, returnAmount, _, err := k.LiquidUnbond(cachedCtx, types.LiquidStakeProxyAcc, types.LiquidStakeProxyAcc, lv.GetOperator(), delShares, false)
				if err != nil {
					logger.Error("liquid unbonding of inactive liquid validator failed", "error", err)
					continue
				}
				writeCache()
				unbondingAmount := sdk.Coin{Denom: k.stakingKeeper.BondDenom(ctx), Amount: returnAmount}.String()
				ctx.EventManager().EmitEvents(sdk.Events{
					sdk.NewEvent(
						types.EventTypeUnbondInactiveLiquidTokens,
						sdk.NewAttribute(types.AttributeKeyLiquidValidator, lv.OperatorAddress),
						sdk.NewAttribute(types.AttributeKeyUnbondingAmount, unbondingAmount),
						sdk.NewAttribute(types.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
					),
				})
				logger.Info(types.EventTypeUnbondInactiveLiquidTokens,
					types.AttributeKeyLiquidValidator, lv.OperatorAddress,
					types.AttributeKeyUnbondingAmount, unbondingAmount,
					types.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339))
			}
			_, found := k.stakingKeeper.GetDelegation(ctx, types.LiquidStakeProxyAcc, lv.GetOperator())
			if !found {
				k.RemoveLiquidValidator(ctx, lv)
				ctx.EventManager().EmitEvents(sdk.Events{
					sdk.NewEvent(
						types.EventTypeRemoveLiquidValidator,
						sdk.NewAttribute(types.AttributeKeyLiquidValidator, lv.OperatorAddress),
					),
				})
				logger.Info(types.EventTypeRemoveLiquidValidator, types.AttributeKeyLiquidValidator, lv.OperatorAddress)
			}
		}
	}

	return reds
}

// AutocompoundStakingRewards withdraws staking rewards and re-stakes when over threshold.
func (k Keeper) AutocompoundStakingRewards(ctx sdk.Context, whitelistedValsMap types.WhitelistedValsMap) {
	totalRemainingRewards, _, totalLiquidTokens := k.CheckDelegationStates(ctx, types.LiquidStakeProxyAcc)

	// checking over types.AutocompoundTrigger and execute GetRewards
	proxyAccBalance := k.GetProxyAccBalance(ctx, types.LiquidStakeProxyAcc)
	rewardsThreshold := types.AutocompoundTrigger.Mul(sdk.NewDecFromInt(totalLiquidTokens))

	// skip If it doesn't exceed the rewards threshold
	if !sdk.NewDecFromInt(proxyAccBalance.Amount).Add(totalRemainingRewards).GT(rewardsThreshold) {
		return
	}

	// Withdraw rewards of LiquidStakeProxyAcc and re-staking
	k.WithdrawLiquidRewards(ctx, types.LiquidStakeProxyAcc)

	// re-staking with proxyAccBalance, due to auto-withdraw on add staking by f1
	proxyAccBalance = k.GetProxyAccBalance(ctx, types.LiquidStakeProxyAcc)

	// skip when no active liquid validator
	activeVals := k.GetActiveLiquidValidators(ctx, whitelistedValsMap)
	if len(activeVals) == 0 {
		return
	}

	// re-staking
	cachedCtx, writeCache := ctx.CacheContext()
	_, err := k.LiquidDelegate(cachedCtx, types.LiquidStakeProxyAcc, activeVals, proxyAccBalance.Amount, whitelistedValsMap)
	if err != nil {
		logger := k.Logger(ctx)
		logger.Error("re-staking failed", "error", err)
		return
	}
	writeCache()
	logger := k.Logger(ctx)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAutocompound,
			sdk.NewAttribute(types.AttributeKeyDelegator, types.LiquidStakeProxyAcc.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, proxyAccBalance.String()),
		),
	})
	logger.Info(types.EventTypeAutocompound,
		types.AttributeKeyDelegator, types.LiquidStakeProxyAcc.String(),
		sdk.AttributeKeyAmount, proxyAccBalance.String())
}
