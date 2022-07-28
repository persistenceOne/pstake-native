package ls_cosmos_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/suite"

	"github.com/persistenceOne/pstake-native/app"
	"github.com/persistenceOne/pstake-native/app/helpers"
	ls_cosmos "github.com/persistenceOne/pstake-native/x/ls-cosmos"
	"github.com/persistenceOne/pstake-native/x/ls-cosmos/types"
)

type HandlerTestSuite struct {
	suite.Suite

	app        *app.PstakeApp
	ctx        sdk.Context
	govHandler govtypes.Handler
}

func (suite *HandlerTestSuite) SetupTest() {
	_, app, ctx := helpers.CreateTestApp()
	suite.app = &app
	suite.ctx = ctx
	suite.govHandler = ls_cosmos.NewLSCosmosProposalHandler(suite.app.LSCosmosKeeper)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestProposalHandler() {
	testCases := []struct {
		name     string
		proposal *types.RegisterCosmosChainProposal
		expErr   bool
	}{
		{
			"all fields",
			types.NewRegisterCosmosChainProposal("title", "description", "connection", "channel-1", "transfer", "uatom", "ustkatom"),
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			err := suite.govHandler(suite.ctx, tc.proposal)
			if tc.expErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
