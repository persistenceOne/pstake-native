package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

func (suite *IntegrationTestSuite) TestTransitionStore() {
	app, ctx := suite.app, suite.ctx

	lscosmosKeeper := app.LSCosmosKeeper

	ibcAmountTransitionStore := types.IbcAmountTransitionStore{
		IbcTransfer: sdk.NewCoins(sdk.NewCoin(lscosmosKeeper.GetHostChainParams(ctx).BaseDenom, sdk.NewInt(10000))),
		IcaDelegate: sdk.NewCoin(lscosmosKeeper.GetHostChainParams(ctx).BaseDenom, sdk.NewInt(10000)),
	}
	lscosmosKeeper.SetIBCTransitionStore(ctx, ibcAmountTransitionStore)

	newIbcAmountTransitionStore := lscosmosKeeper.GetIBCTransitionStore(ctx)
	suite.Equal(ibcAmountTransitionStore, newIbcAmountTransitionStore)
}
