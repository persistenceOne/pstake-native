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
	proposal := types.NewRegisterCosmosChainProposal("title", "description", "connection-0", "channel-0", "transfer", "uatom", "ustkatom", minDeposit.String(), depositFee.String())
	params := types.NewCosmosIBCParams(proposal.IBCConnection, proposal.TokenTransferChannel, proposal.TokenTransferPort, proposal.BaseDenom, proposal.MintDenom, minDeposit, depositFee)
	suite.app.LSCosmosKeeper.SetCosmosIBCParams(ctx, params)

	c := sdk.WrapSDKContext(ctx)
	response, err := app.LSCosmosKeeper.CosmosIBCParams(c, &types.QueryCosmosIBCParamsRequest{})
	suite.NoError(err)
	minDeposit, ok := sdk.NewIntFromString(proposal.MinDeposit)
	if !ok {
		err = sdkErrors.Wrap(err, "minimum deposit amount is invalid")
	}
	suite.NoError(err)
	pStakeDepositFee, err := sdk.NewDecFromStr(proposal.PStakeDepositFee)
	suite.NoError(err)
	cosmoIBCparams := types.NewCosmosIBCParams(proposal.IBCConnection, proposal.TokenTransferChannel, proposal.TokenTransferPort, proposal.BaseDenom, proposal.MintDenom, minDeposit, pStakeDepositFee)
	suite.Equal(&types.QueryCosmosIBCParamsResponse{CosmosIBCParams: cosmoIBCparams}, response)

}
