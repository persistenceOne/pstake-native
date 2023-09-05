package keeper_test

import (
	"fmt"

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
		found    bool
	}{
		{
			name:     "Success",
			input:    types.HostChain{ChainId: suite.chainB.ChainID},
			expected: types.HostChain{ChainId: suite.chainB.ChainID},
			found:    true,
		},
		{
			name:     "NotFound",
			input:    types.HostChain{ChainId: suite.chainB.ChainID},
			expected: types.HostChain{ChainId: ""},
			found:    false,
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			suite.app.LiquidStakeIBCKeeper.SetHostChain(suite.ctx, &t.input)

			hc, found := suite.app.LiquidStakeIBCKeeper.GetHostChain(suite.ctx, t.expected.ChainId)
			suite.Require().Equal(t.found, found)
			suite.Require().Equal(hc.ChainId, t.expected.ChainId)
		})
	}
}

func (suite *IntegrationTestSuite) TestSetHostChainValidator() {
	hcs := suite.app.LiquidStakeIBCKeeper.GetAllHostChains(suite.ctx)

	tc := []struct {
		name      string
		hc        types.HostChain
		validator types.Validator
		amount    int
	}{
		{
			name:      "Create",
			hc:        *hcs[0],
			validator: types.Validator{OperatorAddress: TestAddress},
			amount:    5,
		},
		{
			name:      "Update",
			hc:        *hcs[0],
			validator: *hcs[0].Validators[0],
			amount:    4,
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			suite.app.LiquidStakeIBCKeeper.SetHostChainValidator(suite.ctx, &t.hc, &t.validator)

			suite.Require().Equal(t.amount, len(t.hc.Validators))
		})
	}
}

func (suite *IntegrationTestSuite) TestProcessHostChainValidatorUpdates() {
	epoch := suite.app.EpochsKeeper.GetEpochInfo(suite.ctx, types.DelegationEpoch).CurrentEpoch
	hcs := suite.app.LiquidStakeIBCKeeper.GetAllHostChains(suite.ctx)

	tc := []struct {
		name                   string
		hc                     types.HostChain
		hcValidators           []*types.Validator
		validators             []stakingtypes.Validator
		expected               []*types.Validator
		expectedDelegableState map[string]bool
		err                    error
	}{
		{
			name: "UpdateState",
			hc:   *hcs[0],
			hcValidators: []*types.Validator{
				{
					OperatorAddress: "valoper1",
					Status:          stakingtypes.BondStatusBonded,
					UnbondingEpoch:  0,
					ExchangeRate:    sdk.NewDec(1),
					Delegable:       true,
				},
				{
					OperatorAddress: "valoper2",
					Status:          stakingtypes.BondStatusUnbonding,
					UnbondingEpoch:  types.CurrentUnbondingEpoch(hcs[0].UnbondingFactor, epoch),
					ExchangeRate:    sdk.NewDec(1),
					Delegable:       true,
				},
			},
			validators: []stakingtypes.Validator{
				{
					OperatorAddress:     "valoper1",
					Status:              stakingtypes.Unbonding,
					LiquidShares:        sdk.NewDec(0),
					ValidatorBondShares: sdk.NewDec(0),
					Tokens:              sdk.NewInt(100),
					DelegatorShares:     sdk.NewDec(100),
				},
				{
					OperatorAddress:     "valoper2",
					Status:              stakingtypes.Bonded,
					LiquidShares:        sdk.NewDec(0),
					ValidatorBondShares: sdk.NewDec(10),
					Tokens:              sdk.NewInt(100),
					DelegatorShares:     sdk.NewDec(100),
				},
			},
			expectedDelegableState: map[string]bool{
				"valoper1": false,
				"valoper2": true,
			},
		}, {
			name: "Validator not found.",
			hc:   *hcs[0],
			hcValidators: []*types.Validator{
				{
					OperatorAddress: "valoper1",
					Status:          stakingtypes.BondStatusBonded,
					UnbondingEpoch:  0,
					ExchangeRate:    sdk.NewDec(1),
					Delegable:       true,
				},
				{
					OperatorAddress: "valoper2",
					Status:          stakingtypes.BondStatusUnbonding,
					UnbondingEpoch:  types.CurrentUnbondingEpoch(hcs[0].UnbondingFactor, epoch),
					ExchangeRate:    sdk.NewDec(1),
					Delegable:       true,
				},
			},
			validators: []stakingtypes.Validator{
				{
					OperatorAddress:     "valoper3",
					Status:              stakingtypes.Unbonding,
					LiquidShares:        sdk.NewDec(0),
					ValidatorBondShares: sdk.NewDec(0),
					Tokens:              sdk.NewInt(100),
					DelegatorShares:     sdk.NewDec(100),
				},
			},
			expectedDelegableState: map[string]bool{
				"valoper1": true,
				"valoper2": true,
				"valoper3": true,
			},
			err: fmt.Errorf("validator with address %s not registered", "valoper3"),
		},
		{
			name: "UpdateExchangeRate",
			hc:   *hcs[0],
			hcValidators: []*types.Validator{
				{
					OperatorAddress: TestAddress,
					Status:          stakingtypes.BondStatusBonded,
					DelegatedAmount: sdk.NewInt(10),
					ExchangeRate:    sdk.NewDec(2),
					Delegable:       true,
				},
			},
			validators: []stakingtypes.Validator{
				{
					OperatorAddress:     TestAddress,
					Status:              stakingtypes.Bonded,
					LiquidShares:        sdk.NewDec(0),
					ValidatorBondShares: sdk.NewDec(0),
					Tokens:              sdk.NewInt(100),
					DelegatorShares:     sdk.NewDec(100),
				},
			},
			expectedDelegableState: map[string]bool{
				TestAddress: false,
			},
		},
		{
			name: "UpdateExchangeRateZeroShares",
			hc:   *hcs[0],
			hcValidators: []*types.Validator{
				{
					OperatorAddress: TestAddress,
					Status:          stakingtypes.BondStatusBonded,
					DelegatedAmount: sdk.NewInt(10),
					ExchangeRate:    sdk.NewDec(1),
					Delegable:       false,
				},
			},
			validators: []stakingtypes.Validator{
				{
					OperatorAddress:     TestAddress,
					Status:              stakingtypes.Bonded,
					LiquidShares:        sdk.NewDec(0),
					ValidatorBondShares: sdk.NewDec(0),
					Tokens:              sdk.NewInt(100),
					DelegatorShares:     sdk.NewDec(0),
				},
			},
			expectedDelegableState: map[string]bool{
				TestAddress: false,
			},
		}, {
			name: "UpdateState",
			hc:   *hcs[0],
			hcValidators: []*types.Validator{
				{
					OperatorAddress: "valoper2",
					Status:          stakingtypes.BondStatusUnbonding,
					UnbondingEpoch:  types.CurrentUnbondingEpoch(hcs[0].UnbondingFactor, epoch),
					ExchangeRate:    sdk.NewDec(1),
					DelegatedAmount: sdk.NewInt(100),
					Delegable:       true,
				},
			},
			validators: []stakingtypes.Validator{
				{
					OperatorAddress:     "valoper2",
					Status:              stakingtypes.Bonded,
					LiquidShares:        sdk.NewDec(0),
					ValidatorBondShares: sdk.NewDec(0),
					Tokens:              sdk.NewInt(100),
					DelegatorShares:     sdk.NewDec(101),
				},
			},
			expectedDelegableState: map[string]bool{
				"valoper2": true,
			},
			err: fmt.Errorf("error while querying validator valoper2 delegation: decoding bech32 failed: invalid separator index -1"),
		}, {
			name: "ValidatorDelegable",
			hc:   *hcs[0],
			hcValidators: []*types.Validator{
				{
					OperatorAddress: "valoper2",
					ExchangeRate:    sdk.NewDec(1),
					Delegable:       false,
				},
			},
			validators: []stakingtypes.Validator{
				{
					OperatorAddress:     "valoper2",
					LiquidShares:        sdk.NewDec(30),
					ValidatorBondShares: sdk.NewDec(1),
					Tokens:              sdk.NewInt(100),
					DelegatorShares:     sdk.NewDec(100),
				},
			},
			expectedDelegableState: map[string]bool{
				"valoper2": true,
			},
		}, {
			name: "ValidatorNotDelegableCap",
			hc:   *hcs[0],
			hcValidators: []*types.Validator{
				{
					OperatorAddress: "valoper2",
					ExchangeRate:    sdk.NewDec(1),
					Delegable:       true,
				},
			},
			validators: []stakingtypes.Validator{
				{
					OperatorAddress:     "valoper2",
					LiquidShares:        sdk.NewDec(60),
					ValidatorBondShares: sdk.NewDec(1),
					Tokens:              sdk.NewInt(100),
					DelegatorShares:     sdk.NewDec(100),
				},
			},
			expectedDelegableState: map[string]bool{
				"valoper2": false,
			},
		}, {
			name: "ValidatorNotDelegableBondFactor",
			hc:   *hcs[0],
			hcValidators: []*types.Validator{
				{
					OperatorAddress: "valoper2",
					ExchangeRate:    sdk.NewDec(1),
					Delegable:       true,
				},
			},
			validators: []stakingtypes.Validator{
				{
					OperatorAddress:     "valoper2",
					LiquidShares:        sdk.NewDec(30),
					ValidatorBondShares: sdk.NewDec(0),
					Tokens:              sdk.NewInt(100),
					DelegatorShares:     sdk.NewDec(100),
				},
			},
			expectedDelegableState: map[string]bool{
				"valoper2": false,
			},
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			t.hc.Validators = t.hcValidators

			for _, validator := range t.validators {
				err := suite.app.LiquidStakeIBCKeeper.ProcessHostChainValidatorUpdates(suite.ctx, &t.hc, validator)
				suite.Require().Equal(t.err, err)
			}
			if t.err == nil {
				suite.Require().Equal(len(t.validators), len(t.hc.Validators))

				for i, validator := range t.hc.Validators {
					suite.Require().Equal(t.validators[i].OperatorAddress, validator.OperatorAddress)
					suite.Require().Equal(t.validators[i].Status.String(), validator.Status)

					var exchangeRate sdk.Dec
					if t.validators[i].DelegatorShares.IsZero() {
						exchangeRate = sdk.OneDec()
					} else {
						exchangeRate = sdk.NewDecFromInt(t.validators[i].Tokens).Quo(t.validators[i].DelegatorShares)
					}
					suite.Require().Equal(exchangeRate, validator.ExchangeRate)

					if validator.Status == stakingtypes.BondStatusUnbonding {
						suite.Require().Equal(types.CurrentUnbondingEpoch(hcs[0].UnbondingFactor, epoch), validator.UnbondingEpoch)
					} else if validator.Status == stakingtypes.BondStatusBonded {
						suite.Require().Equal(int64(0), validator.UnbondingEpoch)
					}

					suite.Require().Equal(t.expectedDelegableState[validator.OperatorAddress], validator.Delegable)
				}
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestRedistributeValidatorWeight() {
	hcs := suite.app.LiquidStakeIBCKeeper.GetAllHostChains(suite.ctx)

	tc := []struct {
		name         string
		hc           *types.HostChain
		hcValidators []*types.Validator
		validator    *types.Validator
		expected     map[string]sdk.Dec
	}{
		{
			name: "Success",
			hc:   hcs[0],
			hcValidators: []*types.Validator{
				{
					OperatorAddress: "valoper1",
					Weight:          decFromStr("0.6"),
				},
				{
					OperatorAddress: "valoper2",
					Weight:          decFromStr("0.2"),
				},
				{
					OperatorAddress: "valoper3",
					Weight:          decFromStr("0.15"),
				},
				{
					OperatorAddress: "valoper4",
					Weight:          decFromStr("0.15"),
				},
			},
			validator: &types.Validator{
				OperatorAddress: "valoper1",
				Weight:          decFromStr("0.6"),
			},
			expected: map[string]sdk.Dec{
				"valoper1": decFromStr("0"),
				"valoper2": decFromStr("0.4"),
				"valoper3": decFromStr("0.35"),
				"valoper4": decFromStr("0.35"),
			},
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			t.hc.Validators = t.hcValidators

			suite.app.LiquidStakeIBCKeeper.RedistributeValidatorWeight(suite.ctx, t.hc, t.validator)

			suite.Require().Equal(len(t.hcValidators), len(t.expected))

			for _, validator := range t.hc.Validators {
				suite.Require().Equal(t.expected[validator.OperatorAddress], validator.Weight)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGetAllHostChains() {
	hostChains := suite.app.LiquidStakeIBCKeeper.GetAllHostChains(suite.ctx)

	suite.Require().Equal(1, len(hostChains))
}

func (suite *IntegrationTestSuite) TestGetHostChainFromIBCDenom() {
	tc := []struct {
		name     string
		ibcDenom string
		found    bool
	}{
		{
			name:     "Success",
			ibcDenom: "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2",
			found:    true,
		},
		{
			name:     "NotFound",
			ibcDenom: "ibc/1234",
			found:    false,
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			hc, found := suite.app.LiquidStakeIBCKeeper.GetHostChainFromIbcDenom(suite.ctx, t.ibcDenom)

			suite.Require().Equal(t.found, found)
			if found {
				suite.Require().Equal(suite.chainB.ChainID, hc.ChainId)
			} else {
				suite.Require().Equal("", hc.ChainId)
			}

		})
	}
}

func (suite *IntegrationTestSuite) TestGetHostChainFromDelegatorAddress() {
	hcFromChainID, found := suite.app.LiquidStakeIBCKeeper.GetHostChain(suite.ctx, suite.chainB.ChainID)
	suite.Require().True(found)
	tc := []struct {
		name             string
		delegatorAddress string
		found            bool
	}{
		{
			name:             "Success",
			delegatorAddress: hcFromChainID.DelegationAccount.Address,
			found:            true,
		},
		{
			name:             "NotFound",
			delegatorAddress: "valoper1",
			found:            false,
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			hc, found := suite.app.LiquidStakeIBCKeeper.GetHostChainFromDelegatorAddress(suite.ctx, t.delegatorAddress)

			suite.Require().Equal(t.found, found)

			if t.found {
				suite.Require().Equal(suite.chainB.ChainID, hc.ChainId)
			} else {
				suite.Require().Equal("", hc.ChainId)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGetHostChainFromChannelID() {
	hcFromChainID, found := suite.app.LiquidStakeIBCKeeper.GetHostChain(suite.ctx, suite.chainB.ChainID)
	suite.Require().True(found)
	tc := []struct {
		name      string
		channelId string
		found     bool
	}{
		{
			name:      "Success",
			channelId: hcFromChainID.ChannelId,
			found:     true,
		},
		{
			name:      "NotFound",
			channelId: "channel-1",
			found:     false,
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			hc, found := suite.app.LiquidStakeIBCKeeper.GetHostChainFromChannelID(suite.ctx, t.channelId)

			suite.Require().Equal(t.found, found)

			if t.found {
				suite.Require().Equal(suite.chainB.ChainID, hc.ChainId)
			} else {
				suite.Require().Equal("", hc.ChainId)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGetHostChainCValue() {
	hc, found := suite.app.LiquidStakeIBCKeeper.GetHostChain(suite.ctx, suite.chainB.ChainID)
	suite.Require().Equal(true, found)

	suite.Require().Equal(sdk.OneDec(), hc.CValue)

	testAmount := sdk.NewInt64Coin(hc.MintDenom(), 100)
	suite.Require().NoError(testutil.FundModuleAccount(suite.app.BankKeeper, suite.ctx, types.ModuleName, sdk.NewCoins(testAmount)))

	hc.Validators[0].DelegatedAmount = sdk.NewInt(100)

	suite.Require().Equal(sdk.OneDec(), hc.CValue)
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
			name: "Case 1",
			hc: types.HostChain{
				ChainId: suite.chainB.ChainID,
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
			name: "NotFound",
			hc: types.HostChain{
				ChainId: suite.chainB.ChainID,
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
			name: "InvalidRequest",
			hc: types.HostChain{
				ChainId: suite.chainB.ChainID,
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
			err := suite.app.LiquidStakeIBCKeeper.UpdateHostChainValidatorWeight(
				suite.ctx,
				&t.hc,
				t.validatorAddress,
				t.validatorWeight,
			)

			if t.success {
				suite.Require().NoError(err)
				suite.Require().Equal(len(t.hc.Validators), 1)
				suite.Require().Equal(t.hc.Validators[0].Weight, decFromStr(t.validatorWeight))
			} else {
				suite.Require().Error(err)
				suite.Require().Equal(len(t.hc.Validators), 1)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGetHostChainFromHostDenom() {
	tc := []struct {
		name      string
		hostdenom string
		found     bool
	}{
		{
			name:      "Success",
			hostdenom: "uatom",
			found:     true,
		},
		{
			name:      "NotFound",
			hostdenom: "notuatom",
			found:     false,
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			hc, found := suite.app.LiquidStakeIBCKeeper.GetHostChainFromHostDenom(suite.ctx, t.hostdenom)

			suite.Require().Equal(t.found, found)

			if t.found {
				suite.Require().Equal(suite.chainB.ChainID, hc.ChainId)
			} else {
				suite.Require().Equal("", hc.ChainId)
			}
		})
	}
}
