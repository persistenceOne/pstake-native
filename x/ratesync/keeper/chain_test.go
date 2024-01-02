package keeper_test

import (
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
	"strconv"

	"github.com/persistenceOne/pstake-native/v2/x/ratesync/keeper"
	"github.com/persistenceOne/pstake-native/v2/x/ratesync/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

var ValidHostChainInMsg = func(id uint64) types.HostChain {
	return types.HostChain{
		ID:           id,
		ChainID:      "test-1",
		ConnectionID: ibcexported.LocalhostConnectionID,
		ICAAccount: liquidstakeibctypes.ICAAccount{
			Address:      "",
			Balance:      sdk.Coin{Denom: "", Amount: sdk.ZeroInt()},
			Owner:        types.DefaultPortOwner(id),
			ChannelState: 0,
		},
		Features: types.Feature{
			LiquidStakeIBC: types.LiquidStake{
				FeatureType:     0,
				CodeID:          0,
				Instantiation:   0,
				ContractAddress: "",
				Denoms:          nil,
				Enabled:         false,
			},
			LiquidStake: types.LiquidStake{
				FeatureType:     1,
				CodeID:          0,
				Instantiation:   0,
				ContractAddress: "",
				Denoms:          nil,
				Enabled:         false,
			}},
	}
}

func createNChain(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.HostChain {
	items := make([]types.HostChain, n)
	for i := range items {
		items[i] = ValidHostChainInMsg(uint64(i))
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
