package keeper_test

import (
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (suite *IntegrationTestSuite) TestKeeper_AddGrantTransactions() {
	app, ctx := suite.app, suite.ctx
	appCosmosKeeper := app.CosmosKeeper

	orcastratorAddress := "persistence1hcqg5wj9t42zawqkqucs7la85ffyv08ljhhesu"
	prvKey, err := GetSDKPivKeyAndAddressR("persistence", 118, "together chief must vocal account off apart dinosaur move canvas spring whisper improve cruise idea earn reflect flash goat illegal mistake blood earn ridge")
	suite.NoError(nil, err)
	acc := &authTypes.BaseAccount{
		Address:       orcastratorAddress,
		PubKey:        nil,
		AccountNumber: 1,
		Sequence:      0,
	}
	acc.SetPubKey(prvKey.PubKey())
	appCosmosKeeper.SetAccountState(ctx, acc)
	appCosmosKeeper.SetCurrentAddress(ctx, acc.GetAddress())

	txID := appCosmosKeeper.AddGrantTransactions(ctx, appCosmosKeeper.GetAccountState(ctx, acc.GetAddress()))
	suite.Equal(uint64(1), txID)
}

func (suite *IntegrationTestSuite) TestKeeper_AddFeegrantTransaction() {
	app, ctx := suite.app, suite.ctx
	appCosmosKeeper := app.CosmosKeeper

	orcastratorAddress := "persistence1hcqg5wj9t42zawqkqucs7la85ffyv08ljhhesu"
	prvKey, err := GetSDKPivKeyAndAddressR("persistence", 118, "together chief must vocal account off apart dinosaur move canvas spring whisper improve cruise idea earn reflect flash goat illegal mistake blood earn ridge")
	suite.NoError(nil, err)
	acc := &authTypes.BaseAccount{
		Address:       orcastratorAddress,
		PubKey:        nil,
		AccountNumber: 1,
		Sequence:      0,
	}
	acc.SetPubKey(prvKey.PubKey())
	appCosmosKeeper.SetAccountState(ctx, acc)
	appCosmosKeeper.SetCurrentAddress(ctx, acc.GetAddress())

	txID := appCosmosKeeper.AddFeegrantTransaction(ctx, appCosmosKeeper.GetAccountState(ctx, acc.GetAddress()))
	suite.Equal(uint64(1), txID)
}

func (suite *IntegrationTestSuite) TestKeeper_AddRevokeTransactions() {
	app, ctx := suite.app, suite.ctx
	appCosmosKeeper := app.CosmosKeeper

	orcastratorAddress := "persistence1hcqg5wj9t42zawqkqucs7la85ffyv08ljhhesu"
	prvKey, err := GetSDKPivKeyAndAddressR("persistence", 118, "together chief must vocal account off apart dinosaur move canvas spring whisper improve cruise idea earn reflect flash goat illegal mistake blood earn ridge")
	suite.NoError(nil, err)
	acc := &authTypes.BaseAccount{
		Address:       orcastratorAddress,
		PubKey:        nil,
		AccountNumber: 1,
		Sequence:      0,
	}
	acc.SetPubKey(prvKey.PubKey())
	appCosmosKeeper.SetAccountState(ctx, acc)
	appCosmosKeeper.SetCurrentAddress(ctx, acc.GetAddress())

	txID := appCosmosKeeper.AddRevokeTransactions(ctx, appCosmosKeeper.GetAccountState(ctx, acc.GetAddress()))
	suite.Equal(uint64(1), txID)
}
