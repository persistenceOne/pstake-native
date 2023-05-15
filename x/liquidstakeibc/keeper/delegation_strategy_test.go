package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func DecFromStr(str string) sdk.Dec {
	dec, _ := sdk.NewDecFromStr(str)
	return dec
}

func (suite *IntegrationTestSuite) TestGenerateDelegateMessages() {
	pstakeApp, ctx := suite.app, suite.ctx
	hc, found := pstakeApp.LiquidStakeIBCKeeper.GetHostChain(ctx, suite.path.EndpointB.Chain.ChainID)
	suite.Require().Equal(found, true)

	tc := []struct {
		name                  string
		hc                    *types.HostChain
		validators            []*types.Validator
		expected              map[string]sdk.Int
		totalDelegationAmount sdk.Int
		success               bool
	}{
		{
			name: "one validator has delegated tokens",
			hc:   hc,
			validators: []*types.Validator{
				{
					OperatorAddress: "valoper1",
					Weight:          DecFromStr("0.3"),
					DelegatedAmount: sdk.NewInt(50),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: "valoper2",
					Weight:          DecFromStr("0.2"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: "valoper3",
					Weight:          DecFromStr("0.1"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: "valoper4",
					Weight:          DecFromStr("0.4"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
			},
			expected: map[string]sdk.Int{
				"valoper2": sdk.NewInt(30),
				"valoper3": sdk.NewInt(15),
				"valoper4": sdk.NewInt(55),
			},
			totalDelegationAmount: sdk.NewInt(100),
			success:               true,
		},
		{
			name: "validators with 0 weight and not bonded",
			hc:   hc,
			validators: []*types.Validator{
				{
					OperatorAddress: "valoper1",
					Weight:          DecFromStr("0.6"),
					DelegatedAmount: sdk.NewInt(50),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: "valoper2",
					Weight:          DecFromStr("0.25"),
					DelegatedAmount: sdk.NewInt(60),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: "valoper3",
					Weight:          DecFromStr("0.15"),
					DelegatedAmount: sdk.NewInt(10),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: "valoper4",
					Weight:          DecFromStr("0"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: "valoper5",
					Weight:          DecFromStr("0"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusUnbonded,
				},
			},
			expected: map[string]sdk.Int{
				"valoper1": sdk.NewInt(58),
				"valoper3": sdk.NewInt(2),
			},
			totalDelegationAmount: sdk.NewInt(60),
			success:               true,
		},
		{
			name: "validators with 0 weight and not bonded",
			hc:   hc,
			validators: []*types.Validator{
				{
					OperatorAddress: "valoper1",
					Weight:          DecFromStr("0.25"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: "valoper2",
					Weight:          DecFromStr("0.25"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: "valoper3",
					Weight:          DecFromStr("0.25"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: "valoper4",
					Weight:          DecFromStr("0.25"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
			},
			expected: map[string]sdk.Int{
				"valoper1": sdk.NewInt(25),
				"valoper2": sdk.NewInt(25),
				"valoper3": sdk.NewInt(25),
				"valoper4": sdk.NewInt(25),
			},
			totalDelegationAmount: sdk.NewInt(100),
			success:               true,
		},
		{
			name: "all validators have 0 weight",
			hc:   hc,
			validators: []*types.Validator{
				{
					OperatorAddress: "valoper1",
					Weight:          DecFromStr("0"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: "valoper2",
					Weight:          DecFromStr("0"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: "valoper3",
					Weight:          DecFromStr("0"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: "valoper4",
					Weight:          DecFromStr("0"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
			},
			expected:              map[string]sdk.Int{},
			totalDelegationAmount: sdk.NewInt(100),
			success:               false,
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			hc.Validators = t.validators

			if t.success {
				messages, err := pstakeApp.LiquidStakeIBCKeeper.GenerateDelegateMessages(
					hc,
					t.totalDelegationAmount,
				)
				suite.Require().Equal(err, nil)
				suite.Require().Equal(len(messages), len(t.expected))

				totalAmount := int64(0)
				for _, message := range messages {
					msgDelegate := message.(*stakingtypes.MsgDelegate)

					suite.Require().Equal(t.expected[msgDelegate.ValidatorAddress], msgDelegate.Amount.Amount)

					totalAmount += msgDelegate.Amount.Amount.Int64()
				}

				suite.Require().Equal(t.totalDelegationAmount.Int64(), totalAmount)
			} else {
				_, err := pstakeApp.LiquidStakeIBCKeeper.GenerateDelegateMessages(
					hc,
					t.totalDelegationAmount,
				)
				suite.Error(err)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGenerateUndelegateMessages() {
	pstakeApp, ctx := suite.app, suite.ctx
	hc, found := pstakeApp.LiquidStakeIBCKeeper.GetHostChain(ctx, suite.path.EndpointB.Chain.ChainID)
	suite.Require().Equal(found, true)

	tc := []struct {
		name               string
		hc                 *types.HostChain
		validators         []*types.Validator
		expected           map[string]sdk.Int
		undelegationAmount sdk.Int
		success            bool
	}{
		{
			name: "most frequent case",
			hc:   hc,
			validators: []*types.Validator{
				{
					OperatorAddress: "valoper1",
					Weight:          DecFromStr("0.3"),
					DelegatedAmount: sdk.NewInt(45000),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: "valoper2",
					Weight:          DecFromStr("0.2"),
					DelegatedAmount: sdk.NewInt(25000),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: "valoper3",
					Weight:          DecFromStr("0.1"),
					DelegatedAmount: sdk.NewInt(10000),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: "valoper4",
					Weight:          DecFromStr("0.4"),
					DelegatedAmount: sdk.NewInt(56000),
					Status:          stakingtypes.BondStatusBonded,
				},
			},
			expected: map[string]sdk.Int{
				"valoper1": sdk.NewInt(8700),
				"valoper2": sdk.NewInt(800),
				"valoper4": sdk.NewInt(5500),
			},
			undelegationAmount: sdk.NewInt(15000),
			success:            true,
		},
		{
			name: "validators with 0 weight and not bonded",
			hc:   hc,
			validators: []*types.Validator{
				{
					OperatorAddress: "valoper1",
					Weight:          DecFromStr("0.6"),
					DelegatedAmount: sdk.NewInt(88000),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: "valoper2",
					Weight:          DecFromStr("0.25"),
					DelegatedAmount: sdk.NewInt(42000),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: "valoper3",
					Weight:          DecFromStr("0.15"),
					DelegatedAmount: sdk.NewInt(23000),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: "valoper4",
					Weight:          DecFromStr("0"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: "valoper5",
					Weight:          DecFromStr("0"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusUnbonded,
				},
			},
			expected: map[string]sdk.Int{
				"valoper1": sdk.NewInt(17800),
				"valoper2": sdk.NewInt(12750),
				"valoper3": sdk.NewInt(5450),
			},
			undelegationAmount: sdk.NewInt(36000),
			success:            true,
		},
		{
			name: "all validators have 0 weight",
			hc:   hc,
			validators: []*types.Validator{
				{
					OperatorAddress: "valoper1",
					Weight:          DecFromStr("0"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: "valoper2",
					Weight:          DecFromStr("0"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: "valoper3",
					Weight:          DecFromStr("0"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: "valoper4",
					Weight:          DecFromStr("0"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
			},
			expected:           map[string]sdk.Int{},
			undelegationAmount: sdk.NewInt(10000),
			success:            false,
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			hc.Validators = t.validators

			if t.success {
				messages, err := pstakeApp.LiquidStakeIBCKeeper.GenerateUndelegateMessages(
					hc,
					t.undelegationAmount,
				)
				suite.Require().Equal(err, nil)
				suite.Require().Equal(len(messages), len(t.expected))

				totalAmount := int64(0)
				for _, message := range messages {
					msgUndelegate := message.(*stakingtypes.MsgUndelegate)

					suite.Require().Equal(t.expected[msgUndelegate.ValidatorAddress], msgUndelegate.Amount.Amount)

					totalAmount += msgUndelegate.Amount.Amount.Int64()
				}

				suite.Require().Equal(t.undelegationAmount.Int64(), totalAmount)
			} else {
				_, err := pstakeApp.LiquidStakeIBCKeeper.GenerateUndelegateMessages(
					hc,
					t.undelegationAmount,
				)
				suite.Error(err)
			}
		})
	}
}
