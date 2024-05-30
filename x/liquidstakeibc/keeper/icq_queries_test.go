package keeper_test

import (
	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func (suite *IntegrationTestSuite) TestKeeper_QueryHostChainValidator() {
	pstakeApp, ctx := suite.app, suite.ctx
	k := pstakeApp.LiquidStakeIBCKeeper
	hc, found := k.GetHostChain(ctx, suite.chainB.ChainID)
	suite.Require().Equal(found, true)
	type args struct {
		hc               *types.HostChain
		validatorAddress string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				hc:               hc,
				validatorAddress: hc.Validators[0].OperatorAddress,
			},
			wantErr: false,
		},
		{
			name: "Invalid oper addr",
			args: args{
				hc:               hc,
				validatorAddress: "invalid addr",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if err := k.QueryHostChainValidator(ctx, tt.args.hc, tt.args.validatorAddress); (err != nil) != tt.wantErr {
				suite.T().Errorf("QueryHostChainValidator() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestKeeper_QueryValidatorDelegation() {
	pstakeApp, ctx := suite.app, suite.ctx
	k := pstakeApp.LiquidStakeIBCKeeper
	hc, found := k.GetHostChain(ctx, suite.chainB.ChainID)
	suite.Require().Equal(found, true)

	hc2 := types.HostChain{DelegationAccount: &types.ICAAccount{Address: "invalid"}}

	type args struct {
		hc        *types.HostChain
		validator *types.Validator
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				hc:        hc,
				validator: hc.Validators[0],
			},
			wantErr: false,
		}, {
			name: "Invalid delegator addr",
			args: args{
				hc:        &hc2,
				validator: hc.Validators[0],
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if err := k.QueryValidatorDelegation(ctx, tt.args.hc, tt.args.validator); (err != nil) != tt.wantErr {
				suite.T().Errorf("QueryValidatorDelegation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestKeeper_QueryDelegationHostChainAccountBalance() {
	pstakeApp, ctx := suite.app, suite.ctx
	k := pstakeApp.LiquidStakeIBCKeeper
	hc, found := k.GetHostChain(ctx, suite.chainB.ChainID)
	suite.Require().Equal(found, true)

	hc2 := types.HostChain{DelegationAccount: &types.ICAAccount{Address: "invalid"}}

	type args struct {
		hc *types.HostChain
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Success",
			args:    args{hc: hc},
			wantErr: false,
		}, {
			name:    "Invalidadelegator addr",
			args:    args{hc: &hc2},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if err := k.QueryDelegationHostChainAccountBalance(ctx, tt.args.hc); (err != nil) != tt.wantErr {
				suite.T().Errorf("QueryDelegationHostChainAccountBalance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestKeeper_QueryRewardsHostChainAccountBalance() {
	pstakeApp, ctx := suite.app, suite.ctx
	k := pstakeApp.LiquidStakeIBCKeeper
	hc, found := k.GetHostChain(ctx, suite.chainB.ChainID)
	suite.Require().Equal(found, true)

	hc2 := types.HostChain{RewardsAccount: &types.ICAAccount{Address: "invalid"}}

	type args struct {
		hc *types.HostChain
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Success",
			args:    args{hc: hc},
			wantErr: false,
		}, {
			name:    "invalid rewards addr",
			args:    args{hc: &hc2},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if err := k.QueryRewardsHostChainAccountBalance(ctx, tt.args.hc); (err != nil) != tt.wantErr {
				suite.T().Errorf("QueryRewardsHostChainAccountBalance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestKeeper_QueryNonCompoundableRewardsHostChainAccountBalance() {
	pstakeApp, ctx := suite.app, suite.ctx
	k := pstakeApp.LiquidStakeIBCKeeper
	hc, found := k.GetHostChain(ctx, suite.chainB.ChainID)
	suite.Require().Equal(found, true)

	hc.RewardParams = &types.RewardParams{Destination: "cosmos1g4sr6pcr68v8ng8hfg4pj852cg6kg4cwe40wuw8nxdxl67xp0vusjfs3n0", Denom: "uatom"}
	k.SetHostChain(ctx, hc)

	hc2 := types.HostChain{RewardsAccount: &types.ICAAccount{Address: "invalid"}}

	type args struct {
		hc *types.HostChain
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Success",
			args:    args{hc: hc},
			wantErr: false,
		}, {
			name:    "invalid rewards addr",
			args:    args{hc: &hc2},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if err := k.QueryRewardsHostChainAccountBalance(ctx, tt.args.hc); (err != nil) != tt.wantErr {
				suite.T().Errorf("QueryRewardsHostChainAccountBalance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestKeeper_QueryValidatorDelegationUpdate() {
	pstakeApp, ctx := suite.app, suite.ctx
	k := pstakeApp.LiquidStakeIBCKeeper
	hc, found := k.GetHostChain(ctx, suite.chainB.ChainID)
	suite.Require().Equal(found, true)

	hc2 := types.HostChain{DelegationAccount: &types.ICAAccount{Address: "invalid"}}

	type args struct {
		hc        *types.HostChain
		validator *types.Validator
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				hc:        hc,
				validator: hc.Validators[0],
			},
			wantErr: false,
		}, {
			name: "Invalid delegator addr",
			args: args{
				hc:        &hc2,
				validator: hc.Validators[0],
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if err := k.QueryValidatorDelegationUpdate(ctx, tt.args.hc, tt.args.validator); (err != nil) != tt.wantErr {
				suite.T().Errorf("QueryValidatorDelegationUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
