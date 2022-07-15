package keeper_test

import (
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

var cosmosTx = cosmosTypes.CosmosTx{
	Tx: sdkTx.Tx{
		Body: &sdkTx.TxBody{
			Messages:      []*codecTypes.Any{},
			Memo:          "",
			TimeoutHeight: 0,
		},
		AuthInfo: &sdkTx.AuthInfo{
			SignerInfos: nil,
			Fee: &sdkTx.Fee{
				Amount:   nil,
				GasLimit: cosmosTypes.MinGasFee,
				Payer:    "",
			},
		},
		Signatures: nil,
	},
	EventEmitted:      false,
	Status:            "",
	TxHash:            "",
	ActiveBlockHeight: 10000000,
	SignerAddress:     "address",
}

var txStatus = cosmosTypes.MsgTxStatus{
	OrchestratorAddress: "address",
	TxHash:              "ABCD",
	Status:              "success",
	AccountNumber:       1,
	SequenceNumber:      1,
	Balance:             sdk.NewCoins(sdk.NewCoin("test", sdk.NewInt(10))),
	ValidatorDetails:    []cosmosTypes.ValidatorDetails{},
	BlockHeight:         10,
}

func (suite *IntegrationTestSuite) TestKeeper_SetNewTxnInOutgoingPool() {
	app, ctx := suite.app, suite.ctx
	appCosmosKeeper := app.CosmosKeeper

	// set new txID as 1
	txID := uint64(1)

	// set transaction in outgoing pool
	appCosmosKeeper.SetNewTxnInOutgoingPool(ctx, txID, cosmosTx)

	// get transaction from outgoing pool
	response, err := appCosmosKeeper.GetTxnFromOutgoingPoolByID(ctx, txID)
	suite.NoError(err)
	// check certain fields only as the txBody and authInfo are pointers
	suite.Equal(cosmosTx.EventEmitted, response.CosmosTxDetails.EventEmitted)
	suite.Equal(cosmosTx.ActiveBlockHeight, response.CosmosTxDetails.ActiveBlockHeight)
	suite.Equal(cosmosTx.SignerAddress, response.CosmosTxDetails.SignerAddress)
}

func (suite *IntegrationTestSuite) TestTransactionQueue() {
	app, ctx := suite.app, suite.ctx
	appCosmosKeeper := app.CosmosKeeper

	moduleStatus := appCosmosKeeper.GetParams(ctx).ModuleEnabled
	suite.Equal(true, moduleStatus)
	txID := uint64(1)

	appCosmosKeeper.SetNewInTransactionQueue(ctx, txID)
	appCosmosKeeper.SetNewInTransactionQueue(ctx, txID+1)
	appCosmosKeeper.SetNewInTransactionQueue(ctx, txID+2)

	activeID := appCosmosKeeper.GetNextFromTransactionQueue(ctx)
	suite.Equal(txID, activeID)

	suite.Equal(activeID, appCosmosKeeper.GetActiveFromTransactionQueue(ctx))
	for i := 0; i <= 10; i++ {
		appCosmosKeeper.IncrementRetryCounterInTransactionQueue(ctx, activeID)
	}
	suite.Equal(false, appCosmosKeeper.GetParams(ctx).ModuleEnabled)
}
