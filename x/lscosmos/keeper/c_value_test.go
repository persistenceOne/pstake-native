package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

func (suite *IntegrationTestSuite) TestCValue() {
	app, ctx := suite.app, suite.ctx

	lscosmosKeeper := app.LSCosmosKeeper

	amounts := sdk.NewCoins(sdk.NewInt64Coin(lscosmosKeeper.GetHostChainParams(ctx).MintDenom, 1000000))
	err := app.BankKeeper.MintCoins(ctx, types.ModuleName, amounts)
	suite.NoError(err)

	cValue := lscosmosKeeper.GetCValue(ctx)
	suite.Equal(sdk.OneDec(), cValue)

	tokenValue, residue := lscosmosKeeper.ConvertStkToToken(ctx, sdk.NewDecCoin(lscosmosKeeper.GetHostChainParams(ctx).MintDenom, sdk.NewInt(1000000)))
	suite.True(sdk.NewInt64Coin(lscosmosKeeper.GetIBCDenom(ctx), 1000000).IsEqual(tokenValue))
	suite.True(sdk.NewDecCoinFromDec(lscosmosKeeper.GetIBCDenom(ctx), sdk.ZeroDec()).IsEqual(residue))

	stkValue, residue := lscosmosKeeper.ConvertTokenToStk(ctx, sdk.NewDecCoin(lscosmosKeeper.GetIBCDenom(ctx), sdk.NewInt(1000000)))
	suite.True(sdk.NewInt64Coin(lscosmosKeeper.GetHostChainParams(ctx).MintDenom, 1000000).IsEqual(stkValue))
	suite.True(sdk.NewDecCoinFromDec(lscosmosKeeper.GetHostChainParams(ctx).MintDenom, sdk.ZeroDec()).IsEqual(residue))

	supply := lscosmosKeeper.GetMintedAmount(ctx)
	suite.True(amounts.AmountOf(lscosmosKeeper.GetHostChainParams(ctx).MintDenom).Equal(supply))

	delegationState := types.DelegationState{
		HostAccountDelegations: []types.HostAccountDelegation{
			{
				ValidatorAddress: "",
				Amount:           sdk.NewInt64Coin(lscosmosKeeper.GetHostChainParams(ctx).BaseDenom, 600000),
			},
			{
				ValidatorAddress: "",
				Amount:           sdk.NewInt64Coin(lscosmosKeeper.GetHostChainParams(ctx).BaseDenom, 200000),
			},
			{
				ValidatorAddress: "",
				Amount:           sdk.NewInt64Coin(lscosmosKeeper.GetHostChainParams(ctx).BaseDenom, 100000),
			},
			{
				ValidatorAddress: "",
				Amount:           sdk.NewInt64Coin(lscosmosKeeper.GetHostChainParams(ctx).BaseDenom, 100000),
			},
		},
		HostDelegationAccountBalance: sdk.NewCoins(sdk.NewInt64Coin(lscosmosKeeper.GetHostChainParams(ctx).BaseDenom, 1000)),
	}

	lscosmosKeeper.SetDelegationState(ctx, delegationState)

	stakedAmount := lscosmosKeeper.GetStakedAmount(ctx)
	suite.True(sdk.NewInt(1000000).Equal(stakedAmount))

	cValue = lscosmosKeeper.GetCValue(ctx)
	suite.Equal(sdk.NewDecWithPrec(999000999000999001, 18), cValue)

	tokenValue, residue = lscosmosKeeper.ConvertStkToToken(ctx, sdk.NewDecCoin(lscosmosKeeper.GetHostChainParams(ctx).MintDenom, sdk.NewInt(1000000)))
	suite.True(sdk.NewInt64Coin(lscosmosKeeper.GetIBCDenom(ctx), 1001000).IsEqual(tokenValue))
	suite.True(sdk.NewDecCoinFromDec(lscosmosKeeper.GetIBCDenom(ctx), sdk.ZeroDec()).IsEqual(residue))

	stkValue, residue = lscosmosKeeper.ConvertTokenToStk(ctx, sdk.NewDecCoin(lscosmosKeeper.GetIBCDenom(ctx), sdk.NewInt(1000000)))
	suite.True(sdk.NewInt64Coin(lscosmosKeeper.GetHostChainParams(ctx).MintDenom, 999000).IsEqual(stkValue))
	suite.True(sdk.NewDecCoinFromDec(lscosmosKeeper.GetHostChainParams(ctx).MintDenom, sdk.NewDecWithPrec(999000999001000000, 18)).IsEqual(residue))

	hostChainParams := lscosmosKeeper.GetHostChainParams(ctx)
	ibcDenom := ibctransfertypes.ParseDenomTrace(
		ibctransfertypes.GetPrefixedDenom(
			hostChainParams.TransferPort, hostChainParams.TransferChannel, hostChainParams.BaseDenom,
		),
	).IBCDenom()

	err = app.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(ibcDenom, sdk.NewInt(20000))))
	suite.NoError(err)

	err = app.BankKeeper.SendCoinsFromModuleToModule(ctx,
		types.ModuleName,
		types.DepositModuleAccount,
		sdk.NewCoins(sdk.NewCoin(ibcDenom, sdk.NewInt(10000))),
	)
	suite.NoError(err)

	err = app.BankKeeper.SendCoinsFromModuleToModule(ctx,
		types.ModuleName,
		types.DelegationModuleAccount,
		sdk.NewCoins(sdk.NewCoin(ibcDenom, sdk.NewInt(10000))),
	)
	suite.NoError(err)

	cValue = lscosmosKeeper.GetCValue(ctx)
	suite.Equal(sdk.NewDecWithPrec(979431929480901077, 18), cValue)

	tokenValue, residue = lscosmosKeeper.ConvertStkToToken(ctx, sdk.NewDecCoin(lscosmosKeeper.GetHostChainParams(ctx).MintDenom, sdk.NewInt(1000000)))
	suite.True(sdk.NewInt64Coin(lscosmosKeeper.GetIBCDenom(ctx), 1021000).IsEqual(tokenValue))
	suite.True(sdk.NewDecCoinFromDec(lscosmosKeeper.GetIBCDenom(ctx), sdk.ZeroDec()).IsEqual(residue))

	stkValue, residue = lscosmosKeeper.ConvertTokenToStk(ctx, sdk.NewDecCoin(lscosmosKeeper.GetIBCDenom(ctx), sdk.NewInt(1000000)))
	suite.True(sdk.NewInt64Coin(lscosmosKeeper.GetHostChainParams(ctx).MintDenom, 979431).IsEqual(stkValue))
	suite.True(sdk.NewDecCoinFromDec(lscosmosKeeper.GetHostChainParams(ctx).MintDenom, sdk.NewDecWithPrec(929480901077000000, 18)).IsEqual(residue))
}
