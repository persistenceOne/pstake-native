package lscosmos_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/stretchr/testify/suite"

	"github.com/persistenceOne/pstake-native/v2/app"
	"github.com/persistenceOne/pstake-native/v2/app/helpers"
	"github.com/persistenceOne/pstake-native/v2/x/lscosmos"
)

type HandlerTestSuite struct {
	suite.Suite

	app        *app.PstakeApp
	ctx        sdk.Context
	govHandler govtypes.Handler
}

func (suite *HandlerTestSuite) SetupTest() {
	_, pstakeApp, ctx := helpers.CreateTestApp(suite.T())
	suite.app = &pstakeApp
	suite.ctx = ctx
	suite.govHandler = lscosmos.NewLSCosmosProposalHandler(suite.app.LSCosmosKeeper)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
