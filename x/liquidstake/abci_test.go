package liquidstake_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	chain "github.com/persistenceOne/pstake-native/v4/app"
	"github.com/persistenceOne/pstake-native/v4/x/liquidstake"
	"github.com/persistenceOne/pstake-native/v4/x/liquidstake/keeper"
)

type ABCITestSuite struct {
	suite.Suite

	app    *chain.PstakeApp
	ctx    sdk.Context
	keeper keeper.Keeper
}

func TestABCITestSuite(t *testing.T) {
	suite.Run(t, new(ABCITestSuite))
}

func (s *ABCITestSuite) SetupTest() {
	s.app = chain.Setup(s.T(), false, 5)
	s.ctx = s.app.BaseApp.NewContextLegacy(false, tmproto.Header{})
	s.keeper = s.app.LiquidStakeKeeper
	s.ctx = s.ctx.WithBlockHeight(100).WithBlockTime(chain.ParseTime("2022-03-01T00:00:00Z"))
}

func (s *ABCITestSuite) TestBeginBlock() {
	// Test when module is not paused
	params := s.keeper.GetParams(s.ctx)
	params.ModulePaused = false
	s.Require().NoError(s.keeper.SetParams(s.ctx, params))

	// Call BeginBlock
	liquidstake.BeginBlock(s.ctx, s.keeper)

	// Test when module is paused
	params.ModulePaused = true
	s.Require().NoError(s.keeper.SetParams(s.ctx, params))

	// Call BeginBlock
	liquidstake.BeginBlock(s.ctx, s.keeper)
}
