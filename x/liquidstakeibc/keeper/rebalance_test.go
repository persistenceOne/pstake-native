package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/proto"

	"github.com/persistenceOne/pstake-native/v3/x/liquidstakeibc/types"
)

func (suite *IntegrationTestSuite) TestKeeper_Rebalance() {
	suite.SetupHostChainAB()
	hc, _ := suite.app.LiquidStakeIBCKeeper.GetHostChain(suite.ctx, suite.chainB.ChainID)

	type fields struct {
		validators []*types.Validator
	}
	type args struct {
		epoch           int64
		acceptableDelta math.Int
		maxEntries      uint32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []proto.Message
	}{
		{
			name: "Success",
			fields: fields{
				validators: []*types.Validator{{
					OperatorAddress: "valA",
					Status:          stakingtypes.Bonded.String(),
					Weight:          sdk.MustNewDecFromStr("0.5"),
					DelegatedAmount: sdk.NewInt(1000000),
					ExchangeRate:    sdk.OneDec(),
					UnbondingEpoch:  0,
					Delegable:       true,
				}, {
					OperatorAddress: "valB",
					Status:          stakingtypes.Bonded.String(),
					Weight:          sdk.MustNewDecFromStr("0.5"),
					DelegatedAmount: sdk.NewInt(0),
					ExchangeRate:    sdk.OneDec(),
					UnbondingEpoch:  0,
					Delegable:       true,
				}},
			},
			args: args{
				epoch:           int64(1),
				acceptableDelta: sdk.NewInt(100000),
				maxEntries:      uint32(7),
			},
			want: []proto.Message{&stakingtypes.MsgBeginRedelegate{
				DelegatorAddress:    hc.DelegationAccount.Address,
				ValidatorSrcAddress: "valA",
				ValidatorDstAddress: "valB",
				Amount:              sdk.NewCoin(HostDenom, sdk.NewInt(500000)),
			}},
		}, {
			name: "Success",
			fields: fields{
				validators: []*types.Validator{{
					OperatorAddress: "valA",
					Status:          stakingtypes.Bonded.String(),
					Weight:          sdk.MustNewDecFromStr("0.5"),
					DelegatedAmount: sdk.NewInt(1000000),
					ExchangeRate:    sdk.OneDec(),
					UnbondingEpoch:  0,
					Delegable:       true,
				}, {
					OperatorAddress: "valB",
					Status:          stakingtypes.Bonded.String(),
					Weight:          sdk.MustNewDecFromStr("0.5"),
					DelegatedAmount: sdk.NewInt(0),
					ExchangeRate:    sdk.OneDec(),
					UnbondingEpoch:  0,
					Delegable:       true,
				}},
			},
			args: args{
				epoch:           int64(1),
				acceptableDelta: sdk.NewInt(10000000),
				maxEntries:      uint32(7),
			},
			want: []proto.Message(nil),
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			k := suite.app.LiquidStakeIBCKeeper
			hc, _ = k.GetHostChain(suite.ctx, suite.chainB.ChainID)
			hc.Params.MaxEntries = tt.args.maxEntries
			hc.Params.RedelegationAcceptableDelta = tt.args.acceptableDelta
			hc.Validators = tt.fields.validators
			k.SetHostChain(suite.ctx, hc)
			msgs := k.GenerateRedelegateMsgs(suite.ctx, *hc)
			suite.Require().Equal(tt.want, msgs)
			suite.NotPanics(func() { k.RebalanceWorkflow(suite.ctx, hc.UnbondingFactor) })
			suite.NotPanics(func() { k.RebalanceWorkflow(suite.ctx, hc.UnbondingFactor+1) })
		})
	}
}
