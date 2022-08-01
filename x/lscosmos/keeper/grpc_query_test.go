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
	propsal := testProposal("title", "description", "connection", "channel-1", "transfer", "uatom", "ustkatom", "5", "0.1")
	err := suite.govHandler(ctx, propsal)
	suite.NoError(err)

	c := sdk.WrapSDKContext(ctx)
	response, err := app.LSCosmosKeeper.CosmosIBCParams(c, &types.QueryCosmosIBCParamsRequest{})
	suite.NoError(err)
	minDeposit, ok := sdk.NewIntFromString(propsal.MinDeposit)
	if !ok {
		err = sdkErrors.Wrap(err, "minimum deposit amount is invalid")
	}
	suite.NoError(err)
	pStakeFee, err := sdk.NewDecFromStr(propsal.PStakeFee)
	suite.NoError(err)
	cosmoIBCparams := types.NewCosmosIBCParams(propsal.IBCConnection, propsal.TokenTransferChannel, propsal.TokenTransferPort, propsal.BaseDenom, propsal.MintDenom, minDeposit, pStakeFee)
	suite.Equal(&types.QueryCosmosIBCParamsResponse{CosmosIBCParams: cosmoIBCparams}, response)

}
