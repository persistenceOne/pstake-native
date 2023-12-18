package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/persistenceOne/pstake-native/v2/x/ratesync/keeper"
	"github.com/persistenceOne/pstake-native/v2/x/ratesync/types"
)

func (suite *IntegrationTestSuite) setupMsgServer() (types.MsgServer, context.Context) {
	k, ctx := suite.app.RatesyncKeeper, suite.ctx
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}

func (suite *IntegrationTestSuite) TestMsgServer() {
	ms, ctx := suite.setupMsgServer()
	suite.Require().NotNil(ms)
	suite.Require().NotNil(ctx)
}

func (suite *IntegrationTestSuite) TestChainMsgServerCreate() {
	k, ctx := suite.app.RatesyncKeeper, suite.ctx
	srv := keeper.NewMsgServerImpl(*k)
	wctx := sdk.WrapSDKContext(ctx)

	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateHostChain{Authority: GovAddress.String(),
			HostChain: types.HostChain{ID: uint64(i + 1)},
		}
		_, err := srv.CreateHostChain(wctx, expected)
		suite.Require().NoError(err)
		_, found := k.GetHostChain(ctx,
			expected.HostChain.ID,
		)
		suite.Require().True(found)
	}
}

func (suite *IntegrationTestSuite) TestChainMsgServerUpdate() {
	k, ctx := suite.app.RatesyncKeeper, suite.ctx
	hc := createNChain(k, ctx, 1)[0]
	hc2 := types.HostChain{ID: 1}
	tests := []struct {
		desc    string
		request *types.MsgUpdateHostChain
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdateHostChain{Authority: GovAddress.String(),
				HostChain: hc,
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgUpdateHostChain{Authority: "B",
				HostChain: hc,
			},
			err: sdkerrors.ErrorInvalidSigner,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgUpdateHostChain{Authority: GovAddress.String(),
				HostChain: hc2,
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		suite.T().Run(tc.desc, func(t *testing.T) {
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)
			expected := &types.MsgCreateHostChain{Authority: GovAddress.String(),
				HostChain: hc,
			}

			_, err := srv.UpdateHostChain(wctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				_, found := k.GetHostChain(ctx,
					expected.HostChain.ID,
				)
				suite.Require().True(found)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestChainMsgServerDelete() {
	k, ctx := suite.app.RatesyncKeeper, suite.ctx
	_ = createNChain(k, ctx, 1)[0]

	tests := []struct {
		desc    string
		request *types.MsgDeleteHostChain
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgDeleteHostChain{Authority: GovAddress.String(),
				ID: 0,
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgDeleteHostChain{Authority: "B",
				ID: 0,
			},
			err: sdkerrors.ErrorInvalidSigner,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgDeleteHostChain{Authority: GovAddress.String(),
				ID: 10,
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		suite.T().Run(tc.desc, func(t *testing.T) {
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)

			_, err := srv.DeleteHostChain(wctx, tc.request)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
			} else {
				suite.Require().NoError(err)
				_, found := k.GetHostChain(ctx,
					tc.request.ID,
				)
				suite.Require().False(found)
			}
		})
	}
}
