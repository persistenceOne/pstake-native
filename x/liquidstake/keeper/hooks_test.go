package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	chain "github.com/persistenceOne/pstake-native/v4/app"
	"github.com/persistenceOne/pstake-native/v4/x/liquidstake/keeper"
	"github.com/persistenceOne/pstake-native/v4/x/liquidstake/types"
)

type HooksTestSuite struct {
	suite.Suite

	app    *chain.PstakeApp
	ctx    sdk.Context
	keeper keeper.Keeper
	hooks  keeper.EpochHooks
}

func TestHooksTestSuite(t *testing.T) {
	suite.Run(t, new(HooksTestSuite))
}

func (s *HooksTestSuite) SetupTest() {
	s.app = chain.Setup(s.T(), false, 5)
	s.ctx = s.app.BaseApp.NewContextLegacy(false, tmproto.Header{})
	s.keeper = s.app.LiquidStakeKeeper
	s.hooks = s.keeper.EpochHooks()
	s.ctx = s.ctx.WithBlockHeight(100).WithBlockTime(chain.ParseTime("2022-03-01T00:00:00Z"))
}

func (s *HooksTestSuite) TestEpochHooks() {
	// Test EpochHooks() returns a valid EpochHooks struct
	hooks := s.keeper.EpochHooks()
	s.Require().NotNil(hooks)
}

func (s *HooksTestSuite) TestAfterEpochEnd() {
	// AfterEpochEnd should always return nil
	err := s.hooks.AfterEpochEnd(s.ctx, "any-epoch", 1)
	s.Require().NoError(err)
}

func (s *HooksTestSuite) TestBeforeEpochStart() {
	// Test BeforeEpochStart calls Keeper.BeforeEpochStart
	err := s.hooks.BeforeEpochStart(s.ctx, "test-epoch", 1)
	s.Require().NoError(err)
}

func (s *HooksTestSuite) TestKeeperBeforeEpochStart() {
	// Test when module is paused
	params := s.keeper.GetParams(s.ctx)
	params.ModulePaused = true
	s.Require().NoError(s.keeper.SetParams(s.ctx, params))

	err := s.keeper.BeforeEpochStart(s.ctx, "any-epoch", 1)
	s.Require().NoError(err)

	// Test when module is not paused but epoch is not recognized
	params.ModulePaused = false
	s.Require().NoError(s.keeper.SetParams(s.ctx, params))

	err = s.keeper.BeforeEpochStart(s.ctx, "unknown-epoch", 1)
	s.Require().NoError(err)

	// Test with autocompound epoch
	err = s.keeper.BeforeEpochStart(s.ctx, types.AutocompoundEpoch, 1)
	s.Require().NoError(err)

	// Test with rebalance epoch
	err = s.keeper.BeforeEpochStart(s.ctx, types.RebalanceEpoch, 1)
	s.Require().NoError(err)
}
