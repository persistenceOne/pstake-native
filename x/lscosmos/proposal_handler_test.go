package lscosmos_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/suite"

	"github.com/persistenceOne/pstake-native/app"
	"github.com/persistenceOne/pstake-native/app/helpers"
	"github.com/persistenceOne/pstake-native/x/lscosmos"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

type HandlerTestSuite struct {
	suite.Suite

	app        *app.PstakeApp
	ctx        sdk.Context
	govHandler govtypes.Handler
}

func (suite *HandlerTestSuite) SetupTest() {
	_, pstakeApp, ctx := helpers.CreateTestApp()
	suite.app = &pstakeApp
	suite.ctx = ctx
	suite.govHandler = lscosmos.NewLSCosmosProposalHandler(suite.app.LSCosmosKeeper)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestProposalHandler() {
	testCases := []struct {
		name     string
		proposal *types.RegisterHostChainProposal
		expErr   bool
	}{
		{
			"all fields",
			types.NewRegisterHostChainProposal("title",
				"description",
				true,
				"cosmoshub-4",
				"connection-0",
				"channel-1",
				"transfer",
				"uatom",
				"stk/atom",
				"persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9",
				sdk.OneInt().MulRaw(5),
				types.AllowListedValidators{AllowListedValidators: types.AllowListedVals{{ValidatorAddress: "addr", TargetWeight: sdk.OneDec()}}},
				sdk.ZeroDec(),
				sdk.ZeroDec(),
				sdk.ZeroDec(),
				sdk.ZeroDec()),
			true,
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
