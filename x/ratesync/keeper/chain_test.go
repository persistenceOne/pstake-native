package keeper_test

import (
	"strconv"

	"github.com/persistenceOne/pstake-native/v2/x/ratesync/keeper"
	"github.com/persistenceOne/pstake-native/v2/x/ratesync/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNChain(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.HostChain {
	items := make([]types.HostChain, n)
	for i := range items {
		items[i].ID = uint64(i)

		keeper.SetHostChain(ctx, items[i])
	}
	return items
}

func (suite *IntegrationTestSuite) TestHostChainGet() {
	keeper, ctx := suite.app.RatesyncKeeper, suite.ctx
	items := createNChain(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetHostChain(ctx,
			item.ID,
		)
		suite.Require().True(found)
		suite.Require().Equal(&item, &rst)
	}
}
func (suite *IntegrationTestSuite) TestHostChainRemove() {
	keeper, ctx := suite.app.RatesyncKeeper, suite.ctx
	items := createNChain(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveHostChain(ctx,
			item.ID,
		)
		_, found := keeper.GetHostChain(ctx,
			item.ID,
		)
		suite.Require().False(found)
	}
}

func (suite *IntegrationTestSuite) TestHostChainGetAll() {
	keeper, ctx := suite.app.RatesyncKeeper, suite.ctx
	items := createNChain(keeper, ctx, 10)
	suite.Require().ElementsMatch(items, keeper.GetAllHostChain(ctx))
}
