package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/persistenceOne/pstake-native/x/lscosmos"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

func (suite *IntegrationTestSuite) TestCosmosIBCParamsQuery() {
	app, ctx := suite.app, suite.ctx

	suite.govHandler = lscosmos.NewLSCosmosProposalHandler(suite.app.LSCosmosKeeper)
	minDeposit := sdk.NewInt(5)
	depositFee, err := sdk.NewDecFromStr("0.01")
	suite.NoError(err)
	restakeFee, err := sdk.NewDecFromStr("0.02")
	suite.NoError(err)
	unstakeFee, err := sdk.NewDecFromStr("0.03")
	suite.NoError(err)
	params := types.NewCosmosIBCParams("connection-0", "channel-0", "transfer",
		"uatom", "ustkatom", minDeposit, depositFee, restakeFee, unstakeFee)
	suite.app.LSCosmosKeeper.SetCosmosIBCParams(ctx, params)

	c := sdk.WrapSDKContext(ctx)
	response, err := app.LSCosmosKeeper.CosmosIBCParams(c, &types.QueryCosmosIBCParamsRequest{})
	suite.NoError(err)
	minDeposit, ok := sdk.NewIntFromString("1")
	if !ok {
		err = sdkErrors.Wrap(err, "minimum deposit amount is invalid")
	}
	suite.NoError(err)
	suite.Equal(&types.QueryCosmosIBCParamsResponse{CosmosIBCParams: params}, response)

}
