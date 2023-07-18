package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icqtypes "github.com/persistenceOne/persistence-sdk/v2/x/interchainquery/types"
	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/keeper"
)

func (suite *IntegrationTestSuite) TestValidatorCallback() {
	pstakeApp, ctx := suite.app, suite.ctx
	k := pstakeApp.LiquidStakeIBCKeeper
	hc, found := k.GetHostChain(ctx, suite.chainB.ChainID)
	suite.Require().Equal(found, true)

	makeData := func(validator stakingtypes.Validator) []byte {
		return stakingtypes.MustMarshalValidator(pstakeApp.AppCodec(), &validator)
	}
	type args struct {
		ctx   sdk.Context
		data  []byte
		query icqtypes.Query
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				data: makeData(stakingtypes.Validator{
					OperatorAddress: hc.Validators[0].OperatorAddress,
					Status:          0,
					Tokens:          sdk.NewInt(100),
					DelegatorShares: sdk.NewDec(100),
				}),
				query: icqtypes.Query{ChainId: hc.ChainId},
			},
			wantErr: false,
		},
		{
			name: "Invalid Chain ID",
			args: args{
				data: makeData(stakingtypes.Validator{
					OperatorAddress: hc.Validators[0].OperatorAddress,
					Status:          0,
					Tokens:          sdk.NewInt(100),
					DelegatorShares: sdk.NewDec(100),
				}),
				query: icqtypes.Query{ChainId: "invalid-1"},
			},
			wantErr: true,
		}, {
			name: "Invalid Data",
			args: args{
				data:  []byte("invalid data"),
				query: icqtypes.Query{ChainId: hc.ChainId},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if err := keeper.ValidatorCallback(k, ctx, tt.args.data, tt.args.query); (err != nil) != tt.wantErr {
				suite.T().Errorf("ValidatorCallback() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestDelegationCallback() {
	pstakeApp, ctx := suite.app, suite.ctx
	k := pstakeApp.LiquidStakeIBCKeeper
	hc, found := k.GetHostChain(ctx, suite.chainB.ChainID)
	suite.Require().Equal(found, true)

	for i := range hc.Validators {
		hc.Validators[i].DelegatedAmount = sdk.NewInt(100)
	}
	k.SetHostChain(ctx, hc)

	makeData := func(delegation stakingtypes.Delegation) []byte {
		return stakingtypes.MustMarshalDelegation(pstakeApp.AppCodec(), delegation)

	}
	type args struct {
		ctx   sdk.Context
		data  []byte
		query icqtypes.Query
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				data: makeData(stakingtypes.Delegation{
					DelegatorAddress: hc.DelegationAccount.Address,
					ValidatorAddress: hc.Validators[0].OperatorAddress,
					Shares:           sdk.NewDec(100),
				}),
				query: icqtypes.Query{ChainId: hc.ChainId},
			},
			wantErr: false,
		}, {
			name: "slashed validator",
			args: args{
				data: makeData(stakingtypes.Delegation{
					DelegatorAddress: hc.DelegationAccount.Address,
					ValidatorAddress: hc.Validators[0].OperatorAddress,
					Shares:           sdk.NewDec(10),
				}),
				query: icqtypes.Query{ChainId: hc.ChainId},
			},
			wantErr: false,
		}, {
			name: "Invalid Chain ID",
			args: args{
				data: makeData(stakingtypes.Delegation{
					DelegatorAddress: hc.DelegationAccount.Address,
					ValidatorAddress: hc.Validators[0].OperatorAddress,
					Shares:           sdk.NewDec(100),
				}),
				query: icqtypes.Query{ChainId: "invalid-1"},
			},
			wantErr: true,
		}, {
			name: "Invalid Data",
			args: args{
				data:  []byte("invalid data"),
				query: icqtypes.Query{ChainId: hc.ChainId},
			},
			wantErr: true,
		}, {
			name: "Invalid validator",
			args: args{
				data: makeData(stakingtypes.Delegation{
					DelegatorAddress: hc.DelegationAccount.Address,
					ValidatorAddress: "does not exist",
					Shares:           sdk.NewDec(100),
				}),
				query: icqtypes.Query{ChainId: hc.ChainId},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if err := keeper.DelegationCallback(k, ctx, tt.args.data, tt.args.query); (err != nil) != tt.wantErr {
				suite.T().Errorf("DelegationCallback() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestDelegationAccountBalanceCallback() {
	pstakeApp, ctx := suite.app, suite.ctx
	k := pstakeApp.LiquidStakeIBCKeeper
	hc, found := k.GetHostChain(ctx, suite.chainB.ChainID)
	suite.Require().Equal(found, true)

	makeData := func(denom string, amount int64) []byte {
		coin := sdk.NewInt64Coin(denom, amount)
		return pstakeApp.AppCodec().MustMarshal(&coin)
	}
	type args struct {
		data  []byte
		query icqtypes.Query
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				data:  makeData(hc.HostDenom, 100),
				query: icqtypes.Query{ChainId: hc.ChainId},
			},
			wantErr: false,
		}, {
			name: "invalid chain id",
			args: args{
				data:  makeData(hc.HostDenom, 100),
				query: icqtypes.Query{ChainId: "Invalid Chain ID"},
			},
			wantErr: true,
		}, {
			name: "invalid data",
			args: args{
				data:  []byte("invalid"),
				query: icqtypes.Query{ChainId: hc.ChainId},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if err := keeper.DelegationAccountBalanceCallback(k, ctx, tt.args.data, tt.args.query); (err != nil) != tt.wantErr {
				suite.T().Errorf("DelegationAccountBalanceCallback() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestRewardsAccountBalanceCallback() {
	pstakeApp, ctx := suite.app, suite.ctx
	k := pstakeApp.LiquidStakeIBCKeeper
	hc, found := k.GetHostChain(ctx, suite.chainB.ChainID)
	suite.Require().Equal(found, true)

	makeData := func(denom string, amount int64) []byte {
		coin := sdk.NewInt64Coin(denom, amount)
		return pstakeApp.AppCodec().MustMarshal(&coin)
	}
	for i := range hc.Validators {
		hc.Validators[i].DelegatedAmount = sdk.NewInt(1000000)
	}
	k.SetHostChain(ctx, hc)
	type args struct {
		ctx   sdk.Context
		data  []byte
		query icqtypes.Query
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success, hits the cap",
			args: args{
				data:  makeData(hc.HostDenom, 100),
				query: icqtypes.Query{ChainId: hc.ChainId},
			},
			wantErr: false,
		}, {
			name: "Success, does not hit the cap",
			args: args{
				data:  makeData(hc.HostDenom, 1),
				query: icqtypes.Query{ChainId: hc.ChainId},
			},
			wantErr: false,
		}, {
			name: "invalid chain id",
			args: args{
				data:  makeData(hc.HostDenom, 100),
				query: icqtypes.Query{ChainId: "Invalid Chain ID"},
			},
			wantErr: true,
		}, {
			name: "invalid data",
			args: args{
				data:  []byte("invalid"),
				query: icqtypes.Query{ChainId: hc.ChainId},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if err := keeper.RewardsAccountBalanceCallback(k, ctx, tt.args.data, tt.args.query); (err != nil) != tt.wantErr {
				suite.T().Errorf("RewardsAccountBalanceCallback() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
