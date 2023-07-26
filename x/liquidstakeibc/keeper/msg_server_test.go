package keeper_test

import (
	"context"
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctfrtypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"

	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/keeper"
	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func (suite *IntegrationTestSuite) Test_msgServer_LiquidStake() {
	pstakeapp, ctx := suite.app, suite.ctx
	hc, found := pstakeapp.LiquidStakeIBCKeeper.GetHostChain(ctx, suite.chainB.ChainID)
	suite.Require().True(found)
	epoch := pstakeapp.EpochsKeeper.GetEpochInfo(suite.chainA.GetContext(), types.DelegationEpoch)
	suite.NotNil(epoch)
	err := pstakeapp.LiquidStakeIBCKeeper.BeforeEpochStart(suite.chainA.GetContext(), epoch.Identifier, epoch.CurrentEpoch)
	suite.Require().NoError(err)

	type args struct {
		goCtx context.Context
		msg   *types.MsgLiquidStake
	}
	tests := []struct {
		name    string
		args    args
		want    *types.MsgLiquidStakeResponse
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				goCtx: ctx,
				msg: &types.MsgLiquidStake{
					DelegatorAddress: suite.chainA.SenderAccount.GetAddress().String(),
					Amount:           sdk.NewInt64Coin(hc.IBCDenom(), 1000),
				},
			},
			want:    &types.MsgLiquidStakeResponse{},
			wantErr: false,
		}, {
			name: "host chain with ibc denom not found",
			args: args{
				goCtx: ctx,
				msg: &types.MsgLiquidStake{
					DelegatorAddress: suite.chainA.SenderAccount.GetAddress().String(),
					Amount:           sdk.NewInt64Coin(hc.HostDenom, 1000),
				},
			},
			want:    nil,
			wantErr: true,
		}, {
			name: "amount less than mint amount",
			args: args{
				goCtx: ctx,
				msg: &types.MsgLiquidStake{
					DelegatorAddress: suite.chainA.SenderAccount.GetAddress().String(),
					Amount:           sdk.NewInt64Coin(hc.IBCDenom(), 0),
				},
			},
			want:    nil,
			wantErr: true,
		}, {
			name: "invalid delegator address",
			args: args{
				goCtx: ctx,
				msg: &types.MsgLiquidStake{
					DelegatorAddress: "invalidaddr",
					Amount:           sdk.NewInt64Coin(hc.IBCDenom(), 1000),
				},
			},
			want:    nil,
			wantErr: true,
		}, {
			name: "failed to send tokens",
			args: args{
				goCtx: ctx,
				msg: &types.MsgLiquidStake{
					DelegatorAddress: suite.chainA.SenderAccounts[1].SenderAccount.GetAddress().String(),
					Amount:           sdk.NewInt64Coin(hc.IBCDenom(), 1000),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			k := keeper.NewMsgServerImpl(suite.app.LiquidStakeIBCKeeper)

			got, err := k.LiquidStake(tt.args.goCtx, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("LiquidStake() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LiquidStake() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func (suite *IntegrationTestSuite) Test_msgServer_LiquidStakeLSM() {
	pstakeapp, ctx := suite.app, suite.ctx
	hc, found := pstakeapp.LiquidStakeIBCKeeper.GetHostChain(ctx, suite.chainB.ChainID)
	suite.Require().True(found)

	lsmIbcDenom := ibctfrtypes.ParseDenomTrace(
		ibctfrtypes.GetPrefixedDenom(hc.PortId, hc.ChannelId, hc.Validators[0].OperatorAddress+"/1"),
	).IBCDenom()
	suite.Require().Equal(nil, transfertypes.ValidateIBCDenom(lsmIbcDenom))

	type args struct {
		goCtx               context.Context
		msg                 *types.MsgLiquidStakeLSM
		chainActive         bool
		lsmActive           bool
		createSecondDeposit bool
	}
	tests := []struct {
		name    string
		args    args
		want    *types.MsgLiquidStakeLSMResponse
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				goCtx: ctx,
				msg: &types.MsgLiquidStakeLSM{
					DelegatorAddress: suite.chainA.SenderAccount.GetAddress().String(),
					Delegations:      []*sdk.Coin{{Denom: lsmIbcDenom, Amount: sdk.NewInt(1000)}},
				},
				chainActive:         true,
				lsmActive:           true,
				createSecondDeposit: false,
			},
			want:    &types.MsgLiquidStakeLSMResponse{},
			wantErr: false,
		}, {
			name: "Invalid IBC denom",
			args: args{
				goCtx: ctx,
				msg: &types.MsgLiquidStakeLSM{
					DelegatorAddress: suite.chainA.SenderAccount.GetAddress().String(),
					Delegations:      []*sdk.Coin{{Denom: "ibc", Amount: sdk.NewInt(1000)}},
				},
				chainActive:         true,
				lsmActive:           true,
				createSecondDeposit: false,
			},
			want:    nil,
			wantErr: true,
		}, {
			name: "Invalid IBC hex hash",
			args: args{
				goCtx: ctx,
				msg: &types.MsgLiquidStakeLSM{
					DelegatorAddress: suite.chainA.SenderAccount.GetAddress().String(),
					Delegations:      []*sdk.Coin{{Denom: "uatom", Amount: sdk.NewInt(1000)}},
				},
				chainActive:         true,
				lsmActive:           true,
				createSecondDeposit: false,
			},
			want:    nil,
			wantErr: true,
		}, {
			name: "No IBC denom trace",
			args: args{
				goCtx: ctx,
				msg: &types.MsgLiquidStakeLSM{
					DelegatorAddress: suite.chainA.SenderAccount.GetAddress().String(),
					Delegations:      []*sdk.Coin{{Denom: "ibc/27394FB092D2EAAD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2", Amount: sdk.NewInt(1000)}},
				},
				chainActive:         true,
				lsmActive:           true,
				createSecondDeposit: false,
			},
			want:    nil,
			wantErr: true,
		}, {
			name: "Host chain not active",
			args: args{
				goCtx: ctx,
				msg: &types.MsgLiquidStakeLSM{
					DelegatorAddress: suite.chainA.SenderAccount.GetAddress().String(),
					Delegations:      []*sdk.Coin{{Denom: lsmIbcDenom, Amount: sdk.NewInt(1000)}},
				},
				chainActive:         false,
				lsmActive:           true,
				createSecondDeposit: false,
			},
			want:    nil,
			wantErr: true,
		}, {
			name: "Host chain LSM flag not active",
			args: args{
				goCtx: ctx,
				msg: &types.MsgLiquidStakeLSM{
					DelegatorAddress: suite.chainA.SenderAccount.GetAddress().String(),
					Delegations:      []*sdk.Coin{{Denom: lsmIbcDenom, Amount: sdk.NewInt(1000)}},
				},
				chainActive:         true,
				lsmActive:           false,
				createSecondDeposit: false,
			},
			want:    nil,
			wantErr: true,
		}, {
			name: "Not enough balance",
			args: args{
				goCtx: ctx,
				msg: &types.MsgLiquidStakeLSM{
					DelegatorAddress: suite.chainA.SenderAccount.GetAddress().String(),
					Delegations:      []*sdk.Coin{{Denom: lsmIbcDenom, Amount: sdk.NewInt(1000000)}},
				},
				chainActive:         true,
				lsmActive:           true,
				createSecondDeposit: false,
			},
			want:    nil,
			wantErr: true,
		}, {
			name: "Deposit already exists",
			args: args{
				goCtx: ctx,
				msg: &types.MsgLiquidStakeLSM{
					DelegatorAddress: suite.chainA.SenderAccount.GetAddress().String(),
					Delegations:      []*sdk.Coin{{Denom: lsmIbcDenom, Amount: sdk.NewInt(1000)}},
				},
				chainActive:         true,
				lsmActive:           true,
				createSecondDeposit: true,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			k := keeper.NewMsgServerImpl(suite.app.LiquidStakeIBCKeeper)

			suite.UpdateChainActive(tt.args.chainActive, hc)
			suite.UpdateChainLSMActive(tt.args.lsmActive, hc)

			if tt.args.createSecondDeposit {
				deposit := &types.LSMDeposit{
					ChainId:          suite.chainB.ChainID,
					Denom:            hc.Validators[0].OperatorAddress + "/1",
					DelegatorAddress: suite.chainA.SenderAccount.GetAddress().String(),
				}
				suite.app.LiquidStakeIBCKeeper.SetLSMDeposit(
					ctx,
					deposit,
				)
			}

			got, err := k.LiquidStakeLSM(tt.args.goCtx, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("LiquidStake() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LiquidStake() got = %v, want %v", got, tt.want)
			}

			suite.UpdateChainActive(true, hc)
			suite.UpdateChainLSMActive(true, hc)
		})
	}
}

func (suite *IntegrationTestSuite) Test_msgServer_LiquidUnstake() {
	pstakeapp, ctx := suite.app, suite.ctx
	hc, found := pstakeapp.LiquidStakeIBCKeeper.GetHostChain(ctx, suite.chainB.ChainID)
	suite.Require().True(found)
	type args struct {
		goCtx context.Context
		msg   *types.MsgLiquidUnstake
	}
	tests := []struct {
		name    string
		args    args
		want    *types.MsgLiquidUnstakeResponse
		wantErr bool
	}{
		{
			name: "No tokens to unstake",
			args: args{
				goCtx: ctx,
				msg: &types.MsgLiquidUnstake{
					DelegatorAddress: suite.chainA.SenderAccount.GetAddress().String(),
					Amount:           sdk.NewInt64Coin(hc.MintDenom(), 100),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			k := keeper.NewMsgServerImpl(suite.app.LiquidStakeIBCKeeper)

			got, err := k.LiquidUnstake(tt.args.goCtx, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("LiquidUnstake() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LiquidUnstake() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func (suite *IntegrationTestSuite) Test_msgServer_Redeem() {
	pstakeapp, ctx := suite.app, suite.ctx
	hc, found := pstakeapp.LiquidStakeIBCKeeper.GetHostChain(ctx, suite.chainB.ChainID)
	suite.Require().True(found)

	type args struct {
		goCtx context.Context
		msg   *types.MsgRedeem
	}
	tests := []struct {
		name    string
		args    args
		want    *types.MsgRedeemResponse
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				goCtx: ctx,
				msg: &types.MsgRedeem{
					DelegatorAddress: suite.chainA.SenderAccount.GetAddress().String(),
					Amount:           sdk.NewInt64Coin(hc.MintDenom(), 100),
				},
			},
			want:    nil,
			wantErr: true,
		}}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			k := keeper.NewMsgServerImpl(pstakeapp.LiquidStakeIBCKeeper)

			got, err := k.Redeem(tt.args.goCtx, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Redeem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Redeem() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func (suite *IntegrationTestSuite) Test_msgServer_RegisterHostChain() {
	pstakeapp, ctx := suite.app, suite.ctx

	type args struct {
		goCtx context.Context
		msg   *types.MsgRegisterHostChain
	}
	tests := []struct {
		name    string
		args    args
		want    *types.MsgRegisterHostChainResponse
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				goCtx: ctx,
				msg: &types.MsgRegisterHostChain{
					Authority:          suite.chainA.SenderAccount.GetAddress().String(),
					ConnectionId:       suite.transferPathAC.EndpointA.ConnectionID,
					DepositFee:         sdk.ZeroDec(),
					RestakeFee:         sdk.ZeroDec(),
					UnstakeFee:         sdk.ZeroDec(),
					RedemptionFee:      sdk.ZeroDec(),
					ChannelId:          suite.transferPathAC.EndpointA.ChannelID,
					PortId:             suite.transferPathAC.EndpointA.ChannelConfig.PortID,
					HostDenom:          "uosmo",
					MinimumDeposit:     sdk.OneInt(),
					UnbondingFactor:    4,
					AutoCompoundFactor: 2,
				},
			},
			want:    &types.MsgRegisterHostChainResponse{},
			wantErr: false,
		}}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			k := keeper.NewMsgServerImpl(pstakeapp.LiquidStakeIBCKeeper)

			got, err := k.RegisterHostChain(tt.args.goCtx, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterHostChain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegisterHostChain() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func (suite *IntegrationTestSuite) Test_msgServer_UpdateHostChain() {
	pstakeapp, ctx := suite.app, suite.ctx
	hc, found := pstakeapp.LiquidStakeIBCKeeper.GetHostChain(ctx, suite.chainB.ChainID)
	suite.Require().True(found)

	type args struct {
		goCtx context.Context
		msg   *types.MsgUpdateHostChain
	}
	tests := []struct {
		name    string
		args    args
		want    *types.MsgUpdateHostChainResponse
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				goCtx: ctx,
				msg: &types.MsgUpdateHostChain{
					Authority: suite.chainA.SenderAccount.GetAddress().String(),
					ChainId:   hc.ChainId,
					Updates: []*types.KVUpdate{{
						Key:   types.KeyActive,
						Value: "true",
					}, {
						Key:   types.KeyAutocompoundFactor,
						Value: "20",
					}},
				},
			},
			want:    &types.MsgUpdateHostChainResponse{},
			wantErr: false,
		}}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			k := keeper.NewMsgServerImpl(pstakeapp.LiquidStakeIBCKeeper)

			got, err := k.UpdateHostChain(tt.args.goCtx, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateHostChain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateHostChain() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func (suite *IntegrationTestSuite) Test_msgServer_UpdateParams() {
	pstakeapp, ctx := suite.app, suite.ctx

	type args struct {
		goCtx context.Context
		msg   *types.MsgUpdateParams
	}
	tests := []struct {
		name    string
		args    args
		want    *types.MsgUpdateParamsResponse
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				goCtx: ctx,
				msg: &types.MsgUpdateParams{
					Authority: suite.chainA.SenderAccount.GetAddress().String(),
					Params: types.Params{
						AdminAddress:     suite.chainA.SenderAccount.GetAddress().String(),
						FeeAddress:       suite.chainA.SenderAccount.GetAddress().String(),
						UpperCValueLimit: sdk.OneDec(),
						LowerCValueLimit: sdk.ZeroDec(),
					},
				},
			},
			want:    &types.MsgUpdateParamsResponse{},
			wantErr: false,
		}}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			k := keeper.NewMsgServerImpl(pstakeapp.LiquidStakeIBCKeeper)
			got, err := k.UpdateParams(tt.args.goCtx, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateParams() got = %v, want %v", got, tt.want)
			}
		})
	}
}
