package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *IntegrationTestSuite) TestHostChainParams() {
	app, ctx := suite.app, suite.ctx

	depositFee, err := sdk.NewDecFromStr("0.01")
	suite.NoError(err)

	restakeFee, err := sdk.NewDecFromStr("0.02")
	suite.NoError(err)

	unstakeFee, err := sdk.NewDecFromStr("0.03")
	suite.NoError(err)

	params := app.LSCosmosKeeper.GetHostChainParams(ctx)
	suite.Equal(ChainID, params.ChainID)
	suite.Equal(ConnectionID, params.ConnectionID)
	suite.Equal(TransferChannel, params.TransferChannel)
	suite.Equal(TransferPort, params.TransferPort)
	suite.Equal(BaseDenom, params.BaseDenom)
	suite.Equal(MintDenom, params.MintDenom)
	suite.Equal(PstakeFeeAddress, params.PstakeParams.PstakeFeeAddress)
	suite.Equal(sdk.NewInt(5), params.MinDeposit)
	suite.Equal(depositFee, params.PstakeParams.PstakeDepositFee)
	suite.Equal(restakeFee, params.PstakeParams.PstakeRestakeFee)
	suite.Equal(unstakeFee, params.PstakeParams.PstakeUnstakeFee)
}
