package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

const (
	TestDelegatorAddress = "persistence1234"
	TestLSMDenom         = "cosmosvaloper1234/1"
)

func (suite *IntegrationTestSuite) TestGetSetLSMDeposit() {
	tc := []struct {
		name     string
		input    types.LSMDeposit
		expected types.LSMDeposit
		found    bool
	}{
		{
			name: "Success",
			input: types.LSMDeposit{
				ChainId:          suite.chainB.ChainID,
				DelegatorAddress: TestDelegatorAddress,
				Denom:            "cosmosvaloper1234/1",
			},
			expected: types.LSMDeposit{
				ChainId:          suite.chainB.ChainID,
				DelegatorAddress: TestDelegatorAddress,
				Denom:            "cosmosvaloper1234/1",
			},
			found: true,
		},
		{
			name: "NotFound",
			input: types.LSMDeposit{
				ChainId:          suite.chainB.ChainID,
				DelegatorAddress: TestDelegatorAddress,
				Denom:            "cosmosvaloper1234/1",
			},
			expected: types.LSMDeposit{ChainId: ""},
			found:    false,
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			suite.app.LiquidStakeIBCKeeper.SetLSMDeposit(suite.ctx, &t.input)

			hc, found := suite.app.LiquidStakeIBCKeeper.GetLSMDeposit(
				suite.ctx,
				t.expected.ChainId,
				t.expected.DelegatorAddress,
				t.expected.Denom,
			)
			suite.Require().Equal(t.found, found)
			suite.Require().Equal(hc.ChainId, t.expected.ChainId)
			suite.Require().Equal(hc.DelegatorAddress, t.expected.DelegatorAddress)
			suite.Require().Equal(hc.Denom, t.expected.Denom)
		})
	}
}

func (suite *IntegrationTestSuite) TestDeleteLSMDeposit() {
	deposit := &types.LSMDeposit{
		ChainId:          suite.chainB.ChainID,
		DelegatorAddress: TestDelegatorAddress,
		Denom:            "cosmosvaloper1234/1",
	}

	suite.app.LiquidStakeIBCKeeper.SetLSMDeposit(suite.ctx, deposit)
	suite.app.LiquidStakeIBCKeeper.DeleteLSMDeposit(suite.ctx, deposit)
	deposits := suite.app.LiquidStakeIBCKeeper.FilterLSMDeposits(
		suite.ctx,
		func(d types.LSMDeposit) bool {
			return true
		},
	)

	// preexisting deposit
	suite.Require().Equal(0, len(deposits))
}

func (suite *IntegrationTestSuite) TestGetLSMDepositsFromIbcDenom() {
	suite.app.LiquidStakeIBCKeeper.SetLSMDeposit(
		suite.ctx,
		&types.LSMDeposit{
			ChainId:          suite.chainB.ChainID,
			DelegatorAddress: TestDelegatorAddress,
			Denom:            TestLSMDenom,
			IbcDenom:         "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2",
		},
	)

	tc := []struct {
		name        string
		ibcDenom    string
		expectedLen int
		found       bool
	}{
		{
			name:        "Success",
			ibcDenom:    "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2",
			expectedLen: 1,
			found:       true,
		},
		{
			name:        "NotFound",
			ibcDenom:    "ibc/1234",
			expectedLen: 0,
			found:       false,
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			deposits := suite.app.LiquidStakeIBCKeeper.GetLSMDepositsFromIbcDenom(suite.ctx, t.ibcDenom)

			suite.Require().Equal(t.expectedLen, len(deposits))
			if t.found {
				suite.Require().Equal(
					suite.chainB.ChainID, deposits[0].ChainId,
					TestDelegatorAddress, deposits[0].DelegatorAddress,
					TestLSMDenom, deposits[0].Denom,
				)
			}

		})
	}
}

func (suite *IntegrationTestSuite) TestGetLSMDepositsFromIbcSequenceID() {
	suite.app.LiquidStakeIBCKeeper.SetLSMDeposit(
		suite.ctx,
		&types.LSMDeposit{
			ChainId:          suite.chainB.ChainID,
			DelegatorAddress: TestDelegatorAddress,
			Denom:            TestLSMDenom,
			IbcSequenceId:    "1",
		},
	)

	tc := []struct {
		name        string
		ibcSequence string
		expectedLen int
		found       bool
	}{
		{
			name:        "Success",
			ibcSequence: "1",
			expectedLen: 1,
			found:       true,
		},
		{
			name:        "NotFound",
			ibcSequence: "2",
			expectedLen: 0,
			found:       false,
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			deposits := suite.app.LiquidStakeIBCKeeper.GetLSMDepositsFromIbcSequenceID(suite.ctx, t.ibcSequence)

			suite.Require().Equal(t.expectedLen, len(deposits))
			if t.found {
				suite.Require().Equal(
					suite.chainB.ChainID, deposits[0].ChainId,
					TestDelegatorAddress, deposits[0].DelegatorAddress,
					TestLSMDenom, deposits[0].Denom,
				)
			}

		})
	}
}

func (suite *IntegrationTestSuite) TestGetTransferableLSMDeposits() {
	suite.app.LiquidStakeIBCKeeper.SetLSMDeposit(
		suite.ctx,
		&types.LSMDeposit{
			ChainId:          suite.chainB.ChainID,
			DelegatorAddress: TestDelegatorAddress,
			Denom:            TestLSMDenom,
			State:            types.LSMDeposit_DEPOSIT_PENDING,
		},
	)

	tc := []struct {
		name        string
		chainID     string
		expectedLen int
		found       bool
	}{
		{
			name:        "Success",
			chainID:     suite.chainB.ChainID,
			expectedLen: 1,
			found:       true,
		},
		{
			name:        "NotFound",
			chainID:     suite.chainA.ChainID,
			expectedLen: 0,
			found:       false,
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			deposits := suite.app.LiquidStakeIBCKeeper.GetTransferableLSMDeposits(suite.ctx, t.chainID)

			suite.Require().Equal(t.expectedLen, len(deposits))
			if t.found {
				suite.Require().Equal(
					suite.chainB.ChainID, deposits[0].ChainId,
					TestDelegatorAddress, deposits[0].DelegatorAddress,
					TestLSMDenom, deposits[0].Denom,
				)
			}

		})
	}
}

func (suite *IntegrationTestSuite) TestGetRedeemableLSMDeposits() {
	suite.app.LiquidStakeIBCKeeper.SetLSMDeposit(
		suite.ctx,
		&types.LSMDeposit{
			ChainId:          suite.chainB.ChainID,
			DelegatorAddress: TestDelegatorAddress,
			Denom:            TestLSMDenom,
			State:            types.LSMDeposit_DEPOSIT_RECEIVED,
		},
	)

	tc := []struct {
		name        string
		chainID     string
		expectedLen int
		found       bool
	}{
		{
			name:        "Success",
			chainID:     suite.chainB.ChainID,
			expectedLen: 1,
			found:       true,
		},
		{
			name:        "NotFound",
			chainID:     suite.chainA.ChainID,
			expectedLen: 0,
			found:       false,
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			deposits := suite.app.LiquidStakeIBCKeeper.GetRedeemableLSMDeposits(suite.ctx, t.chainID)

			suite.Require().Equal(t.expectedLen, len(deposits))
			if t.found {
				suite.Require().Equal(
					suite.chainB.ChainID, deposits[0].ChainId,
					TestDelegatorAddress, deposits[0].DelegatorAddress,
					TestLSMDenom, deposits[0].Denom,
				)
			}

		})
	}
}

func (suite *IntegrationTestSuite) TestRevertLSMDepositState() {
	// ibc sequence id is used as index
	deposits := []*types.LSMDeposit{
		{
			ChainId:          suite.chainB.ChainID,
			DelegatorAddress: TestDelegatorAddress,
			Denom:            TestLSMDenom,
			State:            types.LSMDeposit_DEPOSIT_PENDING,
			IbcSequenceId:    "1",
		},
		{
			ChainId:          suite.chainB.ChainID,
			DelegatorAddress: "cosmos2",
			Denom:            "cosmosvaloper1/2",
			State:            types.LSMDeposit_DEPOSIT_SENT,
			IbcSequenceId:    "2",
		},
		{
			ChainId:          suite.chainB.ChainID,
			DelegatorAddress: "cosmos3",
			Denom:            "cosmosvaloper1/3",
			State:            types.LSMDeposit_DEPOSIT_RECEIVED,
			IbcSequenceId:    "3",
		},
		{
			ChainId:          suite.chainB.ChainID,
			DelegatorAddress: "cosmos4",
			Denom:            "cosmosvaloper1/4",
			State:            types.LSMDeposit_DEPOSIT_UNTOKENIZING,
			IbcSequenceId:    "4",
		},
	}

	suite.app.LiquidStakeIBCKeeper.RevertLSMDepositsState(suite.ctx, deposits)

	updatedDeposits := suite.app.LiquidStakeIBCKeeper.FilterLSMDeposits(
		suite.ctx,
		func(d types.LSMDeposit) bool {
			return true
		},
	)

	for _, deposit := range updatedDeposits {
		suite.Assert().Equal("", deposit.IbcDenom)
		switch deposit.IbcSequenceId {
		case "1":
			suite.Assert().Equal(deposit.State, types.Deposit_DEPOSIT_PENDING)
		case "2":
			suite.Assert().Equal(deposit.State, types.Deposit_DEPOSIT_PENDING)
		case "3":
			suite.Assert().Equal(deposit.State, types.Deposit_DEPOSIT_SENT)
		case "4":
			suite.Assert().Equal(deposit.State, types.Deposit_DEPOSIT_RECEIVED)
		}
	}
}

func (suite *IntegrationTestSuite) TestUpdateLSMDepositsStateAndSequence() {
	deposits := []*types.LSMDeposit{
		{
			ChainId:          suite.chainB.ChainID,
			DelegatorAddress: TestDelegatorAddress,
			Denom:            TestLSMDenom,
			State:            types.LSMDeposit_DEPOSIT_PENDING,
			IbcSequenceId:    "1",
		},
		{
			ChainId:          suite.chainB.ChainID,
			DelegatorAddress: "cosmos2",
			Denom:            "cosmosvaloper1/2",
			State:            types.LSMDeposit_DEPOSIT_SENT,
			IbcSequenceId:    "2",
		},
		{
			ChainId:          suite.chainB.ChainID,
			DelegatorAddress: "cosmos3",
			Denom:            "cosmosvaloper1/3",
			State:            types.LSMDeposit_DEPOSIT_RECEIVED,
			IbcSequenceId:    "3",
		},
		{
			ChainId:          suite.chainB.ChainID,
			DelegatorAddress: "cosmos4",
			Denom:            "cosmosvaloper1/4",
			State:            types.LSMDeposit_DEPOSIT_UNTOKENIZING,
			IbcSequenceId:    "4",
		},
	}

	suite.app.LiquidStakeIBCKeeper.UpdateLSMDepositsStateAndSequence(suite.ctx, deposits, types.LSMDeposit_DEPOSIT_RECEIVED, "32")

	updatedDeposits := suite.app.LiquidStakeIBCKeeper.FilterLSMDeposits(
		suite.ctx,
		func(d types.LSMDeposit) bool {
			return true
		},
	)

	for _, deposit := range updatedDeposits {
		suite.Assert().Equal(types.LSMDeposit_DEPOSIT_RECEIVED, deposit.State)
		suite.Assert().Equal("32", deposit.IbcSequenceId)
	}
}

func (suite *IntegrationTestSuite) TestFilterLSMDeposits() {
	deposit := &types.LSMDeposit{
		ChainId:          suite.chainB.ChainID,
		DelegatorAddress: "cosmos1234",
		Denom:            "cosmosvaloper1234/1",
	}

	suite.app.LiquidStakeIBCKeeper.SetLSMDeposit(suite.ctx, deposit)

	deposits := suite.app.LiquidStakeIBCKeeper.FilterLSMDeposits(
		suite.ctx,
		func(d types.LSMDeposit) bool {
			return d.ChainId == suite.chainB.ChainID
		},
	)

	suite.Require().Equal(1, len(deposits))
	suite.Require().Equal(suite.chainB.ChainID, deposits[0].ChainId)
	suite.Require().Equal("cosmos1234", deposits[0].DelegatorAddress)
	suite.Require().Equal("cosmosvaloper1234/1", deposits[0].Denom)
}

func (suite *IntegrationTestSuite) TestGetLSMDepositAmountUntokenized() {
	deposits := []*types.LSMDeposit{
		{
			ChainId:          suite.chainB.ChainID,
			DelegatorAddress: TestDelegatorAddress,
			Denom:            TestLSMDenom,
			Amount:           sdk.NewInt(1000),
		},
		{
			ChainId:          suite.chainB.ChainID,
			DelegatorAddress: TestDelegatorAddress,
			Denom:            "cosmosvaloper1/2",
			Amount:           sdk.NewInt(1000),
		},
		{
			ChainId:          suite.chainB.ChainID,
			DelegatorAddress: TestDelegatorAddress,
			Denom:            "cosmosvaloper1/3",
			Amount:           sdk.NewInt(1000),
		},
		{
			ChainId:          suite.chainB.ChainID,
			DelegatorAddress: TestDelegatorAddress,
			Denom:            "cosmosvaloper1/4",
			Amount:           sdk.NewInt(1000),
		},
	}

	for _, deposit := range deposits {
		suite.app.LiquidStakeIBCKeeper.SetLSMDeposit(suite.ctx, deposit)
	}

	untokenizedAmount := suite.app.LiquidStakeIBCKeeper.GetLSMDepositAmountUntokenized(suite.ctx, suite.chainB.ChainID)

	suite.Assert().Equal(int64(1000*len(deposits)), untokenizedAmount.Int64())
}
