package keeper_test

import (
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	testhelpers "github.com/persistenceOne/pstake-native/v2/app/helpers"
	"github.com/persistenceOne/pstake-native/v2/x/liquidstake/types"
)

func (s *KeeperTestSuite) TestRebalancingCase1() {
	_, valOpers, pks := s.CreateValidators([]int64{1000000, 1000000, 1000000, 1000000, 1000000})
	s.ctx = s.ctx.WithBlockHeight(100).WithBlockTime(testhelpers.ParseTime("2022-03-01T00:00:00Z"))
	params := s.keeper.GetParams(s.ctx)
	params.UnstakeFeeRate = sdk.ZeroDec()
	params.MinLiquidStakeAmount = math.NewInt(10000)
	s.keeper.SetParams(s.ctx, params)
	s.keeper.UpdateLiquidValidatorSet(s.ctx)

	stakingAmt := math.NewInt(49998)
	// add active validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(10)},
	}
	s.keeper.SetParams(s.ctx, params)
	reds := s.keeper.UpdateLiquidValidatorSet(s.ctx)
	s.Require().Len(reds, 0)

	newShares, stkXPRTMintAmt, err := s.keeper.LiquidStake(s.ctx, types.LiquidStakeProxyAcc, s.delAddrs[0], sdk.NewCoin(sdk.DefaultBondDenom, stakingAmt))
	s.Require().NoError(err)
	s.Require().Equal(newShares, stakingAmt.ToLegacyDec())
	s.Require().Equal(stkXPRTMintAmt, stakingAmt)
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx)
	s.Require().Len(reds, 0)

	proxyAccDel1, found := s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[0])
	s.Require().True(found)
	proxyAccDel2, found := s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[1])
	s.Require().True(found)
	proxyAccDel3, found := s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[2])
	s.Require().True(found)

	s.Require().EqualValues(proxyAccDel1.Shares.TruncateInt(), math.NewInt(16666))
	s.Require().EqualValues(proxyAccDel2.Shares.TruncateInt(), math.NewInt(16666))
	s.Require().EqualValues(proxyAccDel3.Shares.TruncateInt(), math.NewInt(16666))
	totalLiquidTokens, _ := s.keeper.GetAllLiquidValidators(s.ctx).TotalLiquidTokens(s.ctx, s.app.StakingKeeper, false)
	s.Require().EqualValues(stakingAmt, totalLiquidTokens)
	s.printRedelegationsLiquidTokens()

	// update whitelist validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[3].String(), TargetWeight: math.NewInt(10)},
	}
	s.keeper.SetParams(s.ctx, params)
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx)
	s.Require().Len(reds, 3)

	proxyAccDel1, found = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[0])
	s.Require().True(found)
	proxyAccDel2, found = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[1])
	s.Require().True(found)
	proxyAccDel3, found = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[2])
	s.Require().True(found)
	proxyAccDel4, found := s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[3])
	s.Require().True(found)

	s.Require().EqualValues(proxyAccDel1.Shares.TruncateInt(), math.NewInt(12501))
	s.Require().EqualValues(proxyAccDel2.Shares.TruncateInt(), math.NewInt(12499))
	s.Require().EqualValues(proxyAccDel3.Shares.TruncateInt(), math.NewInt(12499))
	s.Require().EqualValues(proxyAccDel4.Shares.TruncateInt(), math.NewInt(12499))
	totalLiquidTokens, _ = s.keeper.GetAllLiquidValidators(s.ctx).TotalLiquidTokens(s.ctx, s.app.StakingKeeper, false)
	s.Require().EqualValues(stakingAmt, totalLiquidTokens)
	s.printRedelegationsLiquidTokens()

	//reds := s.app.StakingKeeper.GetRedelegations(s.ctx, types.LiquidStakeProxyAcc, 20)
	s.Require().Len(reds, 3)

	testhelpers.PP("before complete")
	testhelpers.PP(s.keeper.GetAllLiquidValidatorStates(s.ctx))
	testhelpers.PP(s.keeper.GetNetAmountState(s.ctx))

	// advance block time and height for complete redelegations
	s.completeRedelegationUnbonding()

	testhelpers.PP("after complete")
	testhelpers.PP(s.keeper.GetAllLiquidValidatorStates(s.ctx))
	testhelpers.PP(s.keeper.GetNetAmountState(s.ctx))

	// update whitelist validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[3].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[4].String(), TargetWeight: math.NewInt(10)},
	}
	s.keeper.SetParams(s.ctx, params)
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx)
	s.Require().Len(reds, 4)

	proxyAccDel1, found = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[0])
	s.Require().True(found)
	proxyAccDel2, found = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[1])
	s.Require().True(found)
	proxyAccDel3, found = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[2])
	s.Require().True(found)
	proxyAccDel4, found = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[3])
	s.Require().True(found)
	proxyAccDel5, found := s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[4])
	s.Require().True(found)

	s.printRedelegationsLiquidTokens()
	s.Require().EqualValues(proxyAccDel1.Shares.TruncateInt(), math.NewInt(10002))
	s.Require().EqualValues(proxyAccDel2.Shares.TruncateInt(), math.NewInt(9999))
	s.Require().EqualValues(proxyAccDel3.Shares.TruncateInt(), math.NewInt(9999))
	s.Require().EqualValues(proxyAccDel4.Shares.TruncateInt(), math.NewInt(9999))
	s.Require().EqualValues(proxyAccDel5.Shares.TruncateInt(), math.NewInt(9999))
	totalLiquidTokens, _ = s.keeper.GetAllLiquidValidators(s.ctx).TotalLiquidTokens(s.ctx, s.app.StakingKeeper, false)
	s.Require().EqualValues(stakingAmt, totalLiquidTokens)

	// advance block time and height for complete redelegations
	s.completeRedelegationUnbonding()

	// remove whitelist validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[3].String(), TargetWeight: math.NewInt(10)},
	}

	testhelpers.PP(s.keeper.GetAllLiquidValidatorStates(s.ctx))
	s.keeper.SetParams(s.ctx, params)
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx)
	s.Require().Len(reds, 4)
	testhelpers.PP(s.keeper.GetAllLiquidValidatorStates(s.ctx))

	proxyAccDel1, found = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[0])
	s.Require().True(found)
	proxyAccDel2, found = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[1])
	s.Require().True(found)
	proxyAccDel3, found = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[2])
	s.Require().True(found)
	proxyAccDel4, found = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[3])
	s.Require().True(found)
	proxyAccDel5, found = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[4])
	s.Require().False(found)

	s.printRedelegationsLiquidTokens()
	s.Require().EqualValues(proxyAccDel1.Shares.TruncateInt(), math.NewInt(12501))
	s.Require().EqualValues(proxyAccDel2.Shares.TruncateInt(), math.NewInt(12499))
	s.Require().EqualValues(proxyAccDel3.Shares.TruncateInt(), math.NewInt(12499))
	s.Require().EqualValues(proxyAccDel4.Shares.TruncateInt(), math.NewInt(12499))
	totalLiquidTokens, _ = s.keeper.GetAllLiquidValidators(s.ctx).TotalLiquidTokens(s.ctx, s.app.StakingKeeper, false)
	s.Require().EqualValues(stakingAmt, totalLiquidTokens)

	// advance block time and height for complete redelegations
	s.completeRedelegationUnbonding()

	// remove whitelist validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(10)},
	}

	s.keeper.SetParams(s.ctx, params)
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx)
	s.Require().Len(reds, 3)

	proxyAccDel1, found = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[0])
	s.Require().True(found)
	proxyAccDel2, found = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[1])
	s.Require().True(found)
	proxyAccDel3, found = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[2])
	s.Require().False(found)
	proxyAccDel4, found = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[3])
	s.Require().False(found)
	proxyAccDel5, found = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[4])
	s.Require().False(found)

	s.printRedelegationsLiquidTokens()
	s.Require().EqualValues(proxyAccDel1.Shares.TruncateInt(), math.NewInt(24999))
	s.Require().EqualValues(proxyAccDel2.Shares.TruncateInt(), math.NewInt(24999))
	totalLiquidTokens, _ = s.keeper.GetAllLiquidValidators(s.ctx).TotalLiquidTokens(s.ctx, s.app.StakingKeeper, false)
	s.Require().EqualValues(stakingAmt, totalLiquidTokens)

	// advance block time and height for complete redelegations
	s.completeRedelegationUnbonding()

	// double sign, tombstone, slash, jail
	s.doubleSign(valOpers[1], sdk.ConsAddress(pks[1].Address()))

	// check inactive with zero weight after tombstoned
	lvState, found := s.keeper.GetLiquidValidatorState(s.ctx, proxyAccDel2.GetValidatorAddr())
	s.Require().True(found)
	s.Require().Equal(lvState.Status, types.ValidatorStatusInactive)
	s.Require().Equal(lvState.Weight, sdk.ZeroInt())
	s.Require().NotEqualValues(lvState.DelShares, sdk.ZeroDec())
	s.Require().NotEqualValues(lvState.LiquidTokens, sdk.ZeroInt())

	// rebalancing, remove tombstoned liquid validator
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx)
	s.Require().Len(reds, 1)

	// all redelegated, no delShares
	proxyAccDel2, found = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[1])
	s.Require().False(found)

	// liquid validator removed, invalid after tombstoned
	lvState, found = s.keeper.GetLiquidValidatorState(s.ctx, valOpers[1])
	s.Require().False(found)
	s.Require().Equal(lvState.OperatorAddress, valOpers[1].String())
	s.Require().Equal(lvState.Status, types.ValidatorStatusUnspecified)
	s.Require().EqualValues(lvState.DelShares, sdk.ZeroDec())
	s.Require().EqualValues(lvState.LiquidTokens, sdk.ZeroInt())

	// jail last liquid validator, undelegate all liquid tokens to proxy acc
	s.doubleSign(valOpers[0], sdk.ConsAddress(pks[0].Address()))
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx)
	s.Require().Len(reds, 0)

	// no delegation of proxy acc
	proxyAccDel1, found = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[0])
	s.Require().False(found)
	val1, found := s.app.StakingKeeper.GetValidator(s.ctx, valOpers[0])
	s.Require().True(found)
	s.Require().Equal(val1.Status, stakingtypes.Unbonding)

	// check unbonding delegation to proxy acc
	ubd, found := s.app.StakingKeeper.GetUnbondingDelegation(s.ctx, types.LiquidStakeProxyAcc, val1.GetOperator())
	s.Require().True(found)

	// complete unbonding
	s.completeRedelegationUnbonding()

	// check validator Unbonded
	val1, found = s.app.StakingKeeper.GetValidator(s.ctx, valOpers[0])
	s.Require().True(found)
	s.Require().Equal(val1.Status, stakingtypes.Unbonded)

	// no rewards, delShares, liquid tokens
	nas := s.keeper.GetNetAmountState(s.ctx)
	s.Require().EqualValues(nas.TotalRemainingRewards, sdk.ZeroDec())
	s.Require().EqualValues(nas.TotalDelShares, sdk.ZeroDec())
	s.Require().EqualValues(nas.TotalLiquidTokens, sdk.ZeroInt())

	// unbonded to balance, equal with netAmount
	s.Require().EqualValues(ubd.Entries[0].Balance, nas.ProxyAccBalance)
	s.Require().EqualValues(nas.NetAmount.TruncateInt(), nas.ProxyAccBalance)

	// mintRate over 1 due to slashing
	s.Require().True(nas.MintRate.GT(sdk.OneDec()))
	stkXPRTBalanceBefore := s.app.BankKeeper.GetBalance(s.ctx, s.delAddrs[0], params.LiquidBondDenom).Amount
	nativeTokenBalanceBefore := s.app.BankKeeper.GetBalance(s.ctx, s.delAddrs[0], sdk.DefaultBondDenom).Amount
	s.Require().EqualValues(nas.StkxprtTotalSupply, stkXPRTBalanceBefore)

	// withdraw directly unstaking when no totalLiquidTokens
	s.Require().NoError(s.liquidUnstaking(s.delAddrs[0], stkXPRTBalanceBefore, false))
	stkXPRTBalanceAfter := s.app.BankKeeper.GetBalance(s.ctx, s.delAddrs[0], params.LiquidBondDenom).Amount
	nativeTokenBalanceAfter := s.app.BankKeeper.GetBalance(s.ctx, s.delAddrs[0], sdk.DefaultBondDenom).Amount
	s.Require().EqualValues(stkXPRTBalanceAfter, sdk.ZeroInt())
	s.Require().EqualValues(nativeTokenBalanceAfter.Sub(nativeTokenBalanceBefore), nas.NetAmount.TruncateInt())

	// zero net amount states
	s.RequireNetAmountStateZero()
}

func (s *KeeperTestSuite) TestRebalancingConsecutiveCase() {
	_, valOpers, _ := s.CreateValidators([]int64{
		1000000, 1000000, 1000000, 1000000, 1000000,
		1000000, 1000000, 1000000, 1000000, 1000000,
		1000000, 1000000, 1000000, 1000000, 1000000})
	s.ctx = s.ctx.WithBlockHeight(100).WithBlockTime(testhelpers.ParseTime("2022-03-01T00:00:00Z"))
	params := s.keeper.GetParams(s.ctx)
	params.UnstakeFeeRate = sdk.ZeroDec()
	params.MinLiquidStakeAmount = math.NewInt(10000)
	s.keeper.SetParams(s.ctx, params)
	s.keeper.UpdateLiquidValidatorSet(s.ctx)

	stakingAmt := math.NewInt(10000000000000)
	s.fundAddr(s.delAddrs[0], sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, stakingAmt)))
	// add active validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[3].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[4].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[5].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[6].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[7].String(), TargetWeight: math.NewInt(1)},
	}
	s.keeper.SetParams(s.ctx, params)
	reds := s.keeper.UpdateLiquidValidatorSet(s.ctx)
	s.Require().Len(reds, 0)

	newShares, stkXPRTMintAmt, err := s.keeper.LiquidStake(s.ctx, types.LiquidStakeProxyAcc, s.delAddrs[0], sdk.NewCoin(sdk.DefaultBondDenom, stakingAmt))
	s.Require().NoError(err)
	s.Require().Equal(newShares, stakingAmt.ToLegacyDec())
	s.Require().Equal(stkXPRTMintAmt, stakingAmt)
	// assert rebalanced
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx)
	s.Require().Len(reds, 0)
	s.Require().Equal(s.redelegationsErrorCount(reds), 0)
	s.printRedelegationsLiquidTokens()

	// add active validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[3].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[4].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[5].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[6].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[7].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[8].String(), TargetWeight: math.NewInt(5)},
	}
	s.keeper.SetParams(s.ctx, params)
	s.ctx = s.ctx.WithBlockHeight(s.ctx.BlockHeight() + 100).WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 24))
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx)
	s.Require().Len(reds, 8)
	s.Require().Equal(s.redelegationsErrorCount(reds), 0)
	s.printRedelegationsLiquidTokens()
	// assert rebalanced
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx)
	s.Require().Len(reds, 0)

	// add active validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[3].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[4].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[5].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[6].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[7].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[8].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[9].String(), TargetWeight: math.NewInt(1)},
	}
	s.keeper.SetParams(s.ctx, params)
	s.ctx = s.ctx.WithBlockHeight(s.ctx.BlockHeight() + 100).WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 24))
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx)
	s.Require().Len(reds, 9)
	s.Require().Equal(s.redelegationsErrorCount(reds), 0)
	s.printRedelegationsLiquidTokens()
	// assert rebalanced
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx)
	s.Require().Len(reds, 0)

	// complete redelegations
	s.ctx = s.ctx.WithBlockHeight(s.ctx.BlockHeight() + 100).WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 24 * 20).Add(time.Hour))
	staking.EndBlocker(s.ctx, s.app.StakingKeeper)
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx)
	s.Require().Len(reds, 0)
	// assert rebalanced
	s.Require().Equal(s.redelegationsErrorCount(reds), 0)
	s.printRedelegationsLiquidTokens()

	// remove active validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[3].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[4].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[5].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[6].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[7].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[8].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[9].String(), TargetWeight: math.NewInt(1)},
	}
	s.keeper.SetParams(s.ctx, params)
	s.ctx = s.ctx.WithBlockHeight(s.ctx.BlockHeight() + 100).WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 24))
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx)
	s.Require().Len(reds, 9)
	s.Require().Equal(s.redelegationsErrorCount(reds), 0)
	s.printRedelegationsLiquidTokens()
	// assert rebalanced
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx)
	s.Require().Len(reds, 0)

	// add active validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[3].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[4].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[5].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[6].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[7].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[8].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[9].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[10].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[11].String(), TargetWeight: math.NewInt(1)},
		{ValidatorAddress: valOpers[12].String(), TargetWeight: math.NewInt(1)},
	}
	s.keeper.SetParams(s.ctx, params)
	s.ctx = s.ctx.WithBlockHeight(s.ctx.BlockHeight() + 100).WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 24))
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx)
	s.Require().Len(reds, 11)
	// fail rebalancing due to redelegation hopping
	s.Require().Equal(s.redelegationsErrorCount(reds), 11)
	s.printRedelegationsLiquidTokens()

	// complete redelegation and retry
	s.completeRedelegationUnbonding()
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx)
	s.printRedelegationsLiquidTokens()
	s.Require().Len(reds, 11)
	s.Require().Equal(s.redelegationsErrorCount(reds), 0)

	// assert rebalanced
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx)
	s.Require().Len(reds, 0)

	// modify weight
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[3].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[4].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[5].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[6].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[7].String(), TargetWeight: math.NewInt(5)},
		{ValidatorAddress: valOpers[8].String(), TargetWeight: math.NewInt(5)},
		{ValidatorAddress: valOpers[9].String(), TargetWeight: math.NewInt(5)},
		{ValidatorAddress: valOpers[10].String(), TargetWeight: math.NewInt(5)},
		{ValidatorAddress: valOpers[11].String(), TargetWeight: math.NewInt(5)},
		{ValidatorAddress: valOpers[12].String(), TargetWeight: math.NewInt(5)},
	}
	s.keeper.SetParams(s.ctx, params)
	s.ctx = s.ctx.WithBlockHeight(s.ctx.BlockHeight() + 100).WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 24))
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx)
	s.Require().Len(reds, 6)
	// fail rebalancing partially due to redelegation hopping
	s.Require().Equal(s.redelegationsErrorCount(reds), 3)
	s.printRedelegationsLiquidTokens()

	// additional liquid stake when not rebalanced
	_, _, err = s.keeper.LiquidStake(s.ctx, types.LiquidStakeProxyAcc, s.delAddrs[0], sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1000000000)))
	s.Require().NoError(err)
	s.printRedelegationsLiquidTokens()

	// complete some redelegations
	s.ctx = s.ctx.WithBlockHeight(s.ctx.BlockHeight() + 100).WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 24 * 20).Add(time.Hour))
	staking.EndBlocker(s.ctx, s.app.StakingKeeper)
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx)
	s.Require().Len(reds, 9)

	// failed redelegations with small amount (less than rebalancing trigger)
	s.Require().Equal(s.redelegationsErrorCount(reds), 6)
	s.printRedelegationsLiquidTokens()

	// assert rebalanced
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx)
	s.Require().Len(reds, 0)
	s.Require().Equal(s.redelegationsErrorCount(reds), 0)
	s.printRedelegationsLiquidTokens()
}

func (s *KeeperTestSuite) TestAutocompoundStakingRewards() {
	_, valOpers, _ := s.CreateValidators([]int64{2000000, 2000000, 2000000})
	params := s.keeper.GetParams(s.ctx)

	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(10)},
	}
	s.keeper.SetParams(s.ctx, params)
	s.keeper.UpdateLiquidValidatorSet(s.ctx)

	stakingAmt := math.NewInt(100000000)
	s.Require().NoError(s.liquidStaking(s.delAddrs[0], stakingAmt))

	// no rewards
	totalRewards, totalDelShares, totalLiquidTokens := s.keeper.CheckDelegationStates(s.ctx, types.LiquidStakeProxyAcc)
	s.EqualValues(totalRewards, sdk.ZeroDec())
	s.EqualValues(totalDelShares, stakingAmt.ToLegacyDec(), totalLiquidTokens)

	// allocate rewards
	s.advanceHeight(100, false)
	totalRewards, totalDelShares, totalLiquidTokens = s.keeper.CheckDelegationStates(s.ctx, types.LiquidStakeProxyAcc)
	s.NotEqualValues(totalRewards, sdk.ZeroDec())
	s.NotEqualValues(totalLiquidTokens, sdk.ZeroDec())

	// withdraw rewards and re-staking
	whitelistedValsMap := types.GetWhitelistedValsMap(params.WhitelistedValidators)
	s.keeper.AutocompoundStakingRewards(s.ctx, whitelistedValsMap)
	totalRewardsAfter, totalDelSharesAfter, totalLiquidTokensAfter := s.keeper.CheckDelegationStates(s.ctx, types.LiquidStakeProxyAcc)
	s.EqualValues(totalRewardsAfter, sdk.ZeroDec())

	autocompoundFee := params.AutocompoundFeeRate.Mul(totalRewards)
	s.EqualValues(totalDelSharesAfter, totalRewards.Sub(autocompoundFee).Add(totalDelShares).TruncateDec(), totalLiquidTokensAfter)

	stakingParams := s.app.StakingKeeper.GetParams(s.ctx)
	feeAccountBalance := s.app.BankKeeper.GetBalance(
		s.ctx,
		sdk.MustAccAddressFromBech32(params.FeeAccountAddress),
		stakingParams.BondDenom,
	)
	s.EqualValues(autocompoundFee.TruncateInt(), feeAccountBalance.Amount)
}

func (s *KeeperTestSuite) TestRemoveAllLiquidValidator() {
	_, valOpers, _ := s.CreateValidators([]int64{2000000, 2000000, 2000000})
	params := s.keeper.GetParams(s.ctx)

	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(10)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(10)},
	}
	s.Require().NoError(s.keeper.SetParams(s.ctx, params))
	s.keeper.UpdateLiquidValidatorSet(s.ctx)

	stakingAmt := math.NewInt(100000000)
	s.Require().NoError(s.liquidStaking(s.delAddrs[0], stakingAmt))

	// allocate rewards
	s.advanceHeight(1, false)
	nasBefore := s.keeper.GetNetAmountState(s.ctx)
	s.Require().NotEqualValues(sdk.ZeroDec(), nasBefore.TotalRemainingRewards)
	s.Require().NotEqualValues(sdk.ZeroDec(), nasBefore.TotalDelShares)
	s.Require().NotEqualValues(sdk.ZeroDec(), nasBefore.NetAmount)
	s.Require().NotEqualValues(sdk.ZeroInt(), nasBefore.TotalLiquidTokens)
	s.Require().EqualValues(sdk.ZeroInt(), nasBefore.ProxyAccBalance)

	// remove all whitelist
	params.WhitelistedValidators = []types.WhitelistedValidator{}
	s.Require().NoError(s.keeper.SetParams(s.ctx, params))
	s.keeper.UpdateLiquidValidatorSet(s.ctx)

	// no liquid validator
	lvs := s.keeper.GetAllLiquidValidators(s.ctx)
	s.Require().Len(lvs, 0)

	autocompoundFee := params.AutocompoundFeeRate.Mul(nasBefore.TotalRemainingRewards)

	nasAfter := s.keeper.GetNetAmountState(s.ctx)
	s.Require().EqualValues(sdk.ZeroDec(), nasAfter.TotalRemainingRewards)
	s.Require().EqualValues(nasBefore.TotalRemainingRewards.Sub(autocompoundFee).TruncateInt(), nasAfter.ProxyAccBalance)
	s.Require().EqualValues(sdk.ZeroDec(), nasAfter.TotalDelShares)
	s.Require().EqualValues(sdk.ZeroInt(), nasAfter.TotalLiquidTokens)
	s.Require().EqualValues(nasBefore.NetAmount.TruncateInt(), nasAfter.NetAmount.Add(autocompoundFee).TruncateInt())

	s.completeRedelegationUnbonding()
	nasAfter2 := s.keeper.GetNetAmountState(s.ctx)
	s.Require().EqualValues(nasAfter.ProxyAccBalance.Add(nasBefore.TotalLiquidTokens), nasAfter2.ProxyAccBalance)
	s.Require().EqualValues(nasBefore.NetAmount.TruncateInt(), nasAfter2.NetAmount.Add(autocompoundFee).TruncateInt())

	stakingParams := s.app.StakingKeeper.GetParams(s.ctx)
	feeAccountBalance := s.app.BankKeeper.GetBalance(
		s.ctx,
		sdk.MustAccAddressFromBech32(params.FeeAccountAddress),
		stakingParams.BondDenom,
	)

	s.EqualValues(autocompoundFee.TruncateInt(), feeAccountBalance.Amount)
}
