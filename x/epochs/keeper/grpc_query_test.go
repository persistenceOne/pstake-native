package keeper_test

import (
	gocontext "context"
	"time"

	"github.com/persistenceOne/pstake-native/x/epochs/types"
)

func (suite *KeeperTestSuite) TestQueryEpochInfos() {
	suite.SetupTest()
	queryClient := suite.queryClient

	chainStartTime := suite.ctx.BlockTime()

	// Invalid param
	epochInfosResponse, err := queryClient.EpochInfos(gocontext.Background(), &types.QueryEpochsInfoRequest{})
	suite.Require().NoError(err)
	suite.Require().Len(epochInfosResponse.Epochs, 3)

	// check if EpochInfos are correct
	suite.Require().Equal(epochInfosResponse.Epochs[0].Identifier, "reward")
	suite.Require().Equal(epochInfosResponse.Epochs[0].StartTime, chainStartTime)
	suite.Require().Equal(epochInfosResponse.Epochs[0].Duration, time.Minute)
	suite.Require().Equal(epochInfosResponse.Epochs[0].CurrentEpoch, int64(0))
	suite.Require().Equal(epochInfosResponse.Epochs[0].CurrentEpochStartTime, chainStartTime)
	suite.Require().Equal(epochInfosResponse.Epochs[0].EpochCountingStarted, false)
	suite.Require().Equal(epochInfosResponse.Epochs[1].Identifier, "stake")
	suite.Require().Equal(epochInfosResponse.Epochs[1].StartTime, chainStartTime)
	suite.Require().Equal(epochInfosResponse.Epochs[1].Duration, time.Minute)
	suite.Require().Equal(epochInfosResponse.Epochs[1].CurrentEpoch, int64(0))
	suite.Require().Equal(epochInfosResponse.Epochs[1].CurrentEpochStartTime, chainStartTime)
	suite.Require().Equal(epochInfosResponse.Epochs[1].EpochCountingStarted, false)
	suite.Require().Equal(epochInfosResponse.Epochs[2].Identifier, "undelegate")
	suite.Require().Equal(epochInfosResponse.Epochs[2].StartTime, chainStartTime)
	suite.Require().Equal(epochInfosResponse.Epochs[2].Duration, time.Minute*2)
	suite.Require().Equal(epochInfosResponse.Epochs[2].CurrentEpoch, int64(0))
	suite.Require().Equal(epochInfosResponse.Epochs[2].CurrentEpochStartTime, chainStartTime)
	suite.Require().Equal(epochInfosResponse.Epochs[2].EpochCountingStarted, false)
}
