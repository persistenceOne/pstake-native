package keeper_test

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/pkg/errors"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func decFromStr(str string) sdk.Dec {
	dec, _ := sdk.NewDecFromStr(str)
	return dec
}

func (suite *IntegrationTestSuite) TestGenerateDelegateMessages() {
	hc, found := suite.app.LiquidStakeIBCKeeper.GetHostChain(suite.ctx, suite.chainB.ChainID)
	suite.Require().Equal(found, true)

	tc := []struct {
		name                  string
		validators            []*types.Validator
		expected              map[string]int64
		totalDelegationAmount int64
		err                   error
	}{
		{
			name: "Case 1",
			validators: []*types.Validator{
				{
					OperatorAddress: hc.Validators[0].OperatorAddress,
					Weight:          decFromStr("0.3"),
					DelegatedAmount: sdk.NewInt(50),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: hc.Validators[1].OperatorAddress,
					Weight:          decFromStr("0.2"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: hc.Validators[2].OperatorAddress,
					Weight:          decFromStr("0.1"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: hc.Validators[3].OperatorAddress,
					Weight:          decFromStr("0.4"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
			},
			expected: map[string]int64{
				hc.Validators[1].OperatorAddress: int64(30),
				hc.Validators[2].OperatorAddress: int64(15),
				hc.Validators[3].OperatorAddress: int64(55),
			},
			totalDelegationAmount: int64(100),
		},
		{
			name: "Case 2",
			validators: []*types.Validator{
				{
					OperatorAddress: hc.Validators[0].OperatorAddress,
					Weight:          decFromStr("0.6"),
					DelegatedAmount: sdk.NewInt(50),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: hc.Validators[1].OperatorAddress,
					Weight:          decFromStr("0.25"),
					DelegatedAmount: sdk.NewInt(60),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: hc.Validators[2].OperatorAddress,
					Weight:          decFromStr("0.15"),
					DelegatedAmount: sdk.NewInt(10),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: hc.Validators[3].OperatorAddress,
					Weight:          decFromStr("0"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: hc.Validators[3].OperatorAddress,
					Weight:          decFromStr("0"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusUnbonded,
				},
			},
			expected: map[string]int64{
				hc.Validators[0].OperatorAddress: int64(58),
				hc.Validators[2].OperatorAddress: int64(2),
			},
			totalDelegationAmount: int64(60),
		},
		{
			name: "Case 3",
			validators: []*types.Validator{
				{
					OperatorAddress: hc.Validators[0].OperatorAddress,
					Weight:          decFromStr("0.25"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: hc.Validators[1].OperatorAddress,
					Weight:          decFromStr("0.25"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: hc.Validators[2].OperatorAddress,
					Weight:          decFromStr("0.25"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: hc.Validators[3].OperatorAddress,
					Weight:          decFromStr("0.25"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
			},
			expected: map[string]int64{
				hc.Validators[0].OperatorAddress: int64(25),
				hc.Validators[1].OperatorAddress: int64(25),
				hc.Validators[2].OperatorAddress: int64(25),
				hc.Validators[3].OperatorAddress: int64(25),
			},
			totalDelegationAmount: int64(100),
		},
		{
			name: "Case 4",
			validators: []*types.Validator{
				{
					OperatorAddress: hc.Validators[0].OperatorAddress,
					Weight:          decFromStr("0.2"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: hc.Validators[1].OperatorAddress,
					Weight:          decFromStr("0.2"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: hc.Validators[2].OperatorAddress,
					Weight:          decFromStr("0.33"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: hc.Validators[3].OperatorAddress,
					Weight:          decFromStr("0.27"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
			},
			expected:              map[string]int64{},
			totalDelegationAmount: int64(0),
			err:                   errorsmod.Wrap(types.ErrInvalidMessages, "no messages to delegate"),
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			hc.Validators = t.validators

			messages, err := suite.app.LiquidStakeIBCKeeper.GenerateDelegateMessages(
				hc,
				sdk.NewInt(t.totalDelegationAmount),
			)

			suite.Require().Equal(errors.Cause(t.err), errors.Cause(err))
			suite.Require().Equal(len(t.expected), len(messages))

			if err == nil {
				var totalAmount int64
				for _, message := range messages {
					msgDelegate := message.(*stakingtypes.MsgDelegate)

					suite.Require().Equal(t.expected[msgDelegate.ValidatorAddress], msgDelegate.Amount.Amount.Int64())

					totalAmount += msgDelegate.Amount.Amount.Int64()
				}
				suite.Require().Equal(t.totalDelegationAmount, totalAmount)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGenerateUndelegateMessages() {
	hc, found := suite.app.LiquidStakeIBCKeeper.GetHostChain(suite.ctx, suite.chainB.ChainID)
	suite.Require().Equal(found, true)

	tc := []struct {
		name               string
		validators         []*types.Validator
		expected           map[string]int64
		undelegationAmount int64
		err                error
	}{
		{
			name: "Case 1",
			validators: []*types.Validator{
				{
					OperatorAddress: hc.Validators[0].OperatorAddress,
					Weight:          decFromStr("0.3"),
					DelegatedAmount: sdk.NewInt(45000),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: hc.Validators[1].OperatorAddress,
					Weight:          decFromStr("0.2"),
					DelegatedAmount: sdk.NewInt(25000),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: hc.Validators[2].OperatorAddress,
					Weight:          decFromStr("0.1"),
					DelegatedAmount: sdk.NewInt(10000),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: hc.Validators[3].OperatorAddress,
					Weight:          decFromStr("0.4"),
					DelegatedAmount: sdk.NewInt(56000),
					Status:          stakingtypes.BondStatusBonded,
				},
			},
			expected: map[string]int64{
				hc.Validators[0].OperatorAddress: int64(8700),
				hc.Validators[1].OperatorAddress: int64(800),
				hc.Validators[3].OperatorAddress: int64(5500),
			},
			undelegationAmount: int64(15000),
		},
		{
			name: "Case 2",
			validators: []*types.Validator{
				{
					OperatorAddress: hc.Validators[0].OperatorAddress,
					Weight:          decFromStr("0.6"),
					DelegatedAmount: sdk.NewInt(88000),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: hc.Validators[1].OperatorAddress,
					Weight:          decFromStr("0.25"),
					DelegatedAmount: sdk.NewInt(42000),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: hc.Validators[2].OperatorAddress,
					Weight:          decFromStr("0.15"),
					DelegatedAmount: sdk.NewInt(23000),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: hc.Validators[3].OperatorAddress,
					Weight:          decFromStr("0"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
			},
			expected: map[string]int64{
				hc.Validators[0].OperatorAddress: int64(17800),
				hc.Validators[1].OperatorAddress: int64(12750),
				hc.Validators[2].OperatorAddress: int64(5450),
			},
			undelegationAmount: int64(36000),
		},
		{
			name: "Case 3",
			validators: []*types.Validator{
				{
					OperatorAddress: hc.Validators[0].OperatorAddress,
					Weight:          decFromStr("0"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: hc.Validators[1].OperatorAddress,
					Weight:          decFromStr("0"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: hc.Validators[2].OperatorAddress,
					Weight:          decFromStr("0"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
				{
					OperatorAddress: hc.Validators[3].OperatorAddress,
					Weight:          decFromStr("0"),
					DelegatedAmount: sdk.NewInt(0),
					Status:          stakingtypes.BondStatusBonded,
				},
			},
			expected:           map[string]int64{},
			undelegationAmount: int64(10000),
			err:                errorsmod.Wrap(types.ErrInvalidMessages, "no messages to undelegate"),
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			hc.Validators = t.validators

			messages, err := suite.app.LiquidStakeIBCKeeper.GenerateUndelegateMessages(
				hc,
				sdk.NewInt(t.undelegationAmount),
			)

			suite.Require().Equal(errors.Cause(t.err), errors.Cause(err))
			suite.Require().Equal(len(t.expected), len(messages))

			if err == nil {
				var totalAmount int64
				for _, message := range messages {
					msgUndelegate := message.(*stakingtypes.MsgUndelegate)

					suite.Require().Equal(
						t.expected[msgUndelegate.ValidatorAddress],
						msgUndelegate.Amount.Amount.Int64(),
					)

					totalAmount += msgUndelegate.Amount.Amount.Int64()
				}
				suite.Require().Equal(t.undelegationAmount, totalAmount)
			}
		})
	}
}
