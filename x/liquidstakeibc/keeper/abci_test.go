package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
	"time"
)

func (suite *IntegrationTestSuite) TestKeeper_BeginBlockCode() {
	k := suite.app.LiquidStakeIBCKeeper
	ctx := suite.ctx
	suite.SetupHostChainAB()
	hc, _ := k.GetHostChain(ctx, suite.chainB.ChainID)
	delAddr := authtypes.NewModuleAddress("addr1").String()

	// user unbondings
	k.SetUserUnbonding(ctx, &types.UserUnbonding{
		ChainId:      hc.ChainId,
		EpochNumber:  0,
		Address:      delAddr,
		StkAmount:    sdk.NewCoin(hc.MintDenom(), sdk.NewInt(1000000)),
		UnbondAmount: sdk.NewCoin(hc.HostDenom, sdk.NewInt(1000000)),
	})
	k.SetUserUnbonding(ctx, &types.UserUnbonding{
		ChainId:      hc.ChainId,
		EpochNumber:  4,
		Address:      delAddr,
		StkAmount:    sdk.NewCoin(hc.MintDenom(), sdk.NewInt(1000000)),
		UnbondAmount: sdk.NewCoin(hc.HostDenom, sdk.NewInt(1000000)),
	})
	k.SetUserUnbonding(ctx, &types.UserUnbonding{
		ChainId:      hc.ChainId,
		EpochNumber:  8,
		Address:      delAddr,
		StkAmount:    sdk.NewCoin(hc.MintDenom(), sdk.NewInt(1000000)),
		UnbondAmount: sdk.NewCoin(hc.HostDenom, sdk.NewInt(1000000)),
	})
	k.SetUserUnbonding(ctx, &types.UserUnbonding{
		ChainId:      hc.ChainId,
		EpochNumber:  12,
		Address:      delAddr,
		StkAmount:    sdk.NewCoin(hc.MintDenom(), sdk.NewInt(1000000)),
		UnbondAmount: sdk.NewCoin(hc.HostDenom, sdk.NewInt(1000000)),
	})
	k.SetUnbonding(ctx, &types.Unbonding{
		ChainId:       hc.ChainId,
		EpochNumber:   0,
		MatureTime:    time.Unix(int64(ctx.BlockTime().Sub(time.Unix(10, 0)).Seconds()), 0),
		BurnAmount:    sdk.NewCoin(hc.MintDenom(), sdk.NewInt(1000000)),
		UnbondAmount:  sdk.NewCoin(hc.HostDenom, sdk.NewInt(1000000)),
		IbcSequenceId: "channel-0-sequence-4",
		State:         types.Unbonding_UNBONDING_CLAIMABLE,
	})
	k.SetUnbonding(ctx, &types.Unbonding{
		ChainId:       hc.ChainId,
		EpochNumber:   4,
		MatureTime:    time.Time{},
		BurnAmount:    sdk.NewCoin(hc.MintDenom(), sdk.NewInt(1000000)),
		UnbondAmount:  sdk.NewCoin(hc.HostDenom, sdk.NewInt(1000000)),
		IbcSequenceId: "channel-0-sequence-3",
		State:         types.Unbonding_UNBONDING_FAILED,
	})
	k.SetUnbonding(ctx, &types.Unbonding{
		ChainId:       hc.ChainId,
		EpochNumber:   8,
		MatureTime:    time.Unix(int64(ctx.BlockTime().Sub(time.Unix(10, 0)).Seconds()), 0),
		BurnAmount:    sdk.NewCoin(hc.MintDenom(), sdk.NewInt(1000000)),
		UnbondAmount:  sdk.NewCoin(hc.HostDenom, sdk.NewInt(1000000)),
		IbcSequenceId: "channel-0-sequence-2",
		State:         types.Unbonding_UNBONDING_MATURED,
	})
	k.SetUnbonding(ctx, &types.Unbonding{
		ChainId:       hc.ChainId,
		EpochNumber:   12,
		MatureTime:    time.Unix(int64(ctx.BlockTime().Sub(time.Unix(10, 0)).Seconds()), 0),
		BurnAmount:    sdk.NewCoin(hc.MintDenom(), sdk.NewInt(1000000)),
		UnbondAmount:  sdk.NewCoin(hc.HostDenom, sdk.NewInt(1000000)),
		IbcSequenceId: "channel-0-sequence-1",
		State:         types.Unbonding_UNBONDING_MATURING,
	})
	mintcoins := func(denom string) sdk.Coins { return sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(90000000))) }
	err := suite.app.MintKeeper.MintCoins(ctx, mintcoins(hc.MintDenom()))
	suite.Require().Nil(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, types.UndelegationModuleAccount, mintcoins(hc.MintDenom()))
	suite.Require().Nil(err)
	err = suite.app.MintKeeper.MintCoins(ctx, mintcoins(hc.IBCDenom()))
	suite.Require().Nil(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, types.UndelegationModuleAccount, mintcoins(hc.IBCDenom()))
	suite.Require().Nil(err)

	k.DoClaim(ctx, hc)
	ubd, found := k.GetUnbonding(ctx, hc.ChainId, 0)
	suite.Require().Equal(false, found)

	k.DoProcessMaturedUndelegations(ctx, hc)
	ubd, found = k.GetUnbonding(ctx, hc.ChainId, 12)
	suite.Require().Equal(true, found)
	suite.Require().Equal(types.Unbonding_UNBONDING_MATURED, ubd.State)

	// Redeemlsm
	k.SetLSMDeposit(ctx, &types.LSMDeposit{
		ChainId:          hc.ChainId,
		Amount:           sdk.OneInt(),
		Shares:           sdk.OneDec(),
		Denom:            hc.HostDenom,  // should be cosmosvaloperxx/1
		IbcDenom:         hc.IBCDenom(), // should be ibc of denom
		DelegatorAddress: delAddr,
		State:            types.LSMDeposit_DEPOSIT_RECEIVED,
		IbcSequenceId:    "",
	})

	k.DoRedeemLSMTokens(ctx, hc)
	lsmDeposit, ok := k.GetLSMDeposit(ctx, hc.ChainId, delAddr, hc.HostDenom)
	suite.Require().True(ok)
	suite.Require().Equal(types.LSMDeposit_DEPOSIT_UNTOKENIZING, lsmDeposit.State)

	k.SetRedelegationTx(ctx, &types.RedelegateTx{
		ChainId:       hc.ChainId,
		IbcSequenceId: "channel-0-sequence-1",
		State:         types.RedelegateTx_REDELEGATE_ACKED,
	})

	k.DoDeleteRedelegationTxs(ctx)
	_, found = k.GetRedelegationTx(ctx, hc.ChainId, "channel-0-sequence-1")
	suite.Require().False(found)

	k.SetRedelegations(ctx, hc.ChainId, []*stakingtypes.Redelegation{{
		DelegatorAddress:    delAddr,
		ValidatorSrcAddress: "val1",
		ValidatorDstAddress: "val2",
		Entries: []stakingtypes.RedelegationEntry{{
			CompletionTime: time.Unix(int64(ctx.BlockTime().Sub(time.Unix(10, 0)).Seconds()), 0),
		}},
	}})

	k.DoDeleteMaturedRedelegation(ctx, hc)
	redels, found := k.GetRedelegations(ctx, hc.ChainId)
	suite.Require().True(found)
	suite.Require().Equal(0, len(redels.Redelegations[0].Entries))

	k.SetValidatorUnbonding(ctx, &types.ValidatorUnbonding{
		ChainId:          hc.ChainId,
		EpochNumber:      0,
		MatureTime:       time.Unix(int64(ctx.BlockTime().Sub(time.Unix(10, 0)).Seconds()), 0),
		ValidatorAddress: "val1",
		Amount:           sdk.NewInt64Coin(hc.HostDenom, 100000),
		IbcSequenceId:    "",
	})
	k.DoProcessMaturedUndelegations(ctx, hc)
	valubd, found := k.GetValidatorUnbonding(ctx, hc.ChainId, "val1", 0)
	suite.Require().True(found)
	suite.Require().NotEqual("", valubd.IbcSequenceId)

}
