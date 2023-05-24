package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func (suite *IntegrationTestSuite) TestGetSetHostChain() {
	tc := []struct {
		name     string
		input    types.HostChain
		expected types.HostChain
		success  bool
	}{
		{
			name:     "success test",
			input:    types.HostChain{ChainId: "hc1"},
			expected: types.HostChain{ChainId: "hc1"},
			success:  true,
		},
		{
			name:     "unsuccessful test",
			input:    types.HostChain{ChainId: "hc1"},
			expected: types.HostChain{ChainId: "hc2"},
			success:  false,
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			pstakeApp, ctx := suite.app, suite.ctx

			pstakeApp.LiquidStakeIBCKeeper.SetHostChain(ctx, &t.input)
			hc, found := pstakeApp.LiquidStakeIBCKeeper.GetHostChain(ctx, t.expected.ChainId)
			if t.success {
				suite.Require().Equal(found, true)
				suite.Require().Equal(hc.ChainId, t.expected.ChainId)
			} else {
				suite.Require().Equal(found, false)
				suite.Require().Equal(hc.ChainId, "")
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestSetHostChainValidator() {
	tc := []struct {
		name      string
		hc        types.HostChain
		validator types.Validator
	}{
		{
			name: "new validator",
			hc:   types.HostChain{ChainId: "hc1", Validators: make([]*types.Validator, 0)},
			validator: types.Validator{
				OperatorAddress: "valoper1",
				Status:          stakingtypes.BondStatusBonded,
				Weight:          sdk.OneDec(),
				DelegatedAmount: sdk.NewInt(100),
			},
		},
		{
			name: "update validator",
			hc: types.HostChain{
				ChainId: "hc1",
				Validators: []*types.Validator{
					{
						OperatorAddress: "valoper1",
						Status:          stakingtypes.BondStatusBonded,
						Weight:          sdk.OneDec(),
						DelegatedAmount: sdk.NewInt(100),
					},
				},
			},
			validator: types.Validator{
				OperatorAddress: "valoper1",
				Status:          stakingtypes.BondStatusBonded,
				Weight:          DecFromStr("0.5"),
				DelegatedAmount: sdk.NewInt(150),
			},
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			pstakeApp, ctx := suite.app, suite.ctx

			pstakeApp.LiquidStakeIBCKeeper.SetHostChainValidator(ctx, &t.hc, &t.validator)

			suite.Require().Equal(len(t.hc.Validators), 1)
			suite.Require().Equal(t.hc.Validators[0], &t.validator)
		})
	}
}

func (suite *IntegrationTestSuite) ProcessHostChainValidatorUpdates() {
	tc := []struct {
		name       string
		hc         types.HostChain
		validators []stakingtypes.Validator
		expected   []*types.Validator
	}{
		{
			name: "new validator",
			hc:   types.HostChain{ChainId: "hc1", Validators: make([]*types.Validator, 0)},
			validators: []stakingtypes.Validator{
				{
					OperatorAddress: "valoper1",
					Status:          stakingtypes.Bonded,
				},
				{
					OperatorAddress: "valoper2",
					Status:          stakingtypes.Bonded,
				},
			},
		},
		{
			name: "create and update validator",
			hc: types.HostChain{
				ChainId: "hc1",
				Validators: []*types.Validator{
					{
						OperatorAddress: "valoper1",
						Status:          stakingtypes.BondStatusBonded,
						Weight:          sdk.OneDec(),
						DelegatedAmount: sdk.NewInt(100),
					},
				},
			},
			validators: []stakingtypes.Validator{
				{
					OperatorAddress: "valoper1",
					Status:          stakingtypes.Unbonding,
				},
				{
					OperatorAddress: "valoper2",
					Status:          stakingtypes.Bonded,
				},
			},
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			pstakeApp, ctx := suite.app, suite.ctx

			err := pstakeApp.LiquidStakeIBCKeeper.ProcessHostChainValidatorUpdates(ctx, &t.hc, t.validators)

			suite.Require().Equal(err, nil)
			suite.Require().Equal(len(t.hc.Validators), len(t.validators))
			for i, validator := range t.hc.Validators {
				suite.Require().NotEqual(err, nil)
				suite.Require().Equal(validator.OperatorAddress, t.validators[i].OperatorAddress)
				suite.Require().Equal(validator.Status, t.validators[i].Status.String())
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGetAllHostChains() {
	pstakeApp, ctx := suite.app, suite.ctx
	hostChains := pstakeApp.LiquidStakeIBCKeeper.GetAllHostChains(ctx)

	suite.Require().Equal(len(hostChains), 1)
}

func (suite *IntegrationTestSuite) TestGetHostChainFromIBCDenom() {
	tc := []struct {
		name     string
		ibcDenom string
		success  bool
	}{
		{
			name:     "retrieve successfully",
			ibcDenom: "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2",
			success:  true,
		},
		{
			name:     "not any chain with ibc denom",
			ibcDenom: "ibc/1234",
			success:  false,
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			pstakeApp, ctx := suite.app, suite.ctx

			hc, found := pstakeApp.LiquidStakeIBCKeeper.GetHostChainFromIbcDenom(ctx, t.ibcDenom)
			if t.success {
				suite.Require().Equal(found, true)
				suite.Require().Equal(hc.ChainId, suite.chainB.ChainID)
			} else {
				suite.Require().Equal(found, false)
				suite.Require().Equal(hc.ChainId, "")
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGetHostChainFromDelegatorAddress() {
	tc := []struct {
		name             string
		delegatorAddress string
		success          bool
	}{
		{
			name:             "retrieve successfully",
			delegatorAddress: "cosmos1mykw6u6dq4z7qhw9aztpk5yp8j8y5n0c6usg9faqepw83y2u4nzq2qxaxc",
			success:          true,
		},
		{
			name:             "no chains with delegator",
			delegatorAddress: "valoper1",
			success:          false,
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			pstakeApp, ctx := suite.app, suite.ctx

			hc, found := pstakeApp.LiquidStakeIBCKeeper.GetHostChainFromDelegatorAddress(ctx, t.delegatorAddress)
			if t.success {
				suite.Require().Equal(found, true)
				suite.Require().Equal(hc.ChainId, suite.chainB.ChainID)
			} else {
				suite.Require().Equal(found, false)
				suite.Require().Equal(hc.ChainId, "")
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGetHostChainCValue() {
	pstakeApp, ctx := suite.app, suite.ctx

	hc, found := pstakeApp.LiquidStakeIBCKeeper.GetHostChain(ctx, suite.path.EndpointB.Chain.ChainID)
	suite.Require().Equal(found, true)

	cValue := pstakeApp.LiquidStakeIBCKeeper.GetHostChainCValue(ctx, hc)
	suite.Require().Equal(cValue, sdk.OneDec())

	testAmount := sdk.NewInt64Coin(hc.MintDenom(), 100)
	suite.Require().NoError(testutil.FundModuleAccount(pstakeApp.BankKeeper, ctx, types.ModuleName, sdk.NewCoins(testAmount)))

	hc.Validators[0].DelegatedAmount = sdk.NewInt(100)

	cValue = pstakeApp.LiquidStakeIBCKeeper.GetHostChainCValue(ctx, hc)
	suite.Require().Equal(cValue, sdk.OneDec())
}

func (suite *IntegrationTestSuite) TestUpdateHostChainValidatorWeight() {
	tc := []struct {
		name             string
		hc               types.HostChain
		validatorAddress string
		validatorWeight  string
		success          bool
	}{
		{
			name: "validator exists",
			hc: types.HostChain{
				ChainId: "hc1",
				Validators: []*types.Validator{
					{
						OperatorAddress: "valoper1",
						Status:          stakingtypes.BondStatusBonded,
						Weight:          sdk.OneDec(),
						DelegatedAmount: sdk.NewInt(100),
					},
				},
			},
			validatorAddress: "valoper1",
			validatorWeight:  "0.5",
			success:          true,
		},
		{
			name: "validator doesn't exists",
			hc: types.HostChain{
				ChainId: "hc1",
				Validators: []*types.Validator{
					{
						OperatorAddress: "valoper1",
						Status:          stakingtypes.BondStatusBonded,
						Weight:          sdk.OneDec(),
						DelegatedAmount: sdk.NewInt(100),
					},
				},
			},
			validatorAddress: "valoper2",
			validatorWeight:  "1",
			success:          false,
		},
		{
			name: "wrong weight value",
			hc: types.HostChain{
				ChainId: "hc1",
				Validators: []*types.Validator{
					{
						OperatorAddress: "valoper1",
						Status:          stakingtypes.BondStatusBonded,
						Weight:          sdk.OneDec(),
						DelegatedAmount: sdk.NewInt(100),
					},
				},
			},
			validatorAddress: "valoper1",
			validatorWeight:  "weight",
			success:          false,
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			pstakeApp, ctx := suite.app, suite.ctx

			err := pstakeApp.LiquidStakeIBCKeeper.UpdateHostChainValidatorWeight(
				ctx,
				&t.hc,
				t.validatorAddress,
				t.validatorWeight,
			)

			if t.success {
				suite.Require().NoError(err)
				suite.Require().Equal(len(t.hc.Validators), 1)
				suite.Require().Equal(t.hc.Validators[0].Weight, DecFromStr(t.validatorWeight))
			} else {
				suite.Require().Error(err)
				suite.Require().Equal(len(t.hc.Validators), 1)
			}
		})
	}
}
