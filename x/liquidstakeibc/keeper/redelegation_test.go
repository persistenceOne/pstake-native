package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var redelegation = &stakingtypes.Redelegation{
	DelegatorAddress:    "delAddr",
	ValidatorSrcAddress: "valSrcAddr",
	ValidatorDstAddress: "valDstAddr",
	Entries: []stakingtypes.RedelegationEntry{{
		CreationHeight:          1,
		CompletionTime:          time.Now(),
		InitialBalance:          sdk.Int{},
		SharesDst:               sdk.Dec{},
		UnbondingId:             0,
		UnbondingOnHoldRefCount: 0,
	}},
}

func (suite *IntegrationTestSuite) TestSetGetRedelegations() {

	suite.app.LiquidStakeIBCKeeper.SetRedelegations(
		suite.ctx, suite.chainB.ChainID,
		[]*stakingtypes.Redelegation{redelegation},
	)

	redelegations, found := suite.app.LiquidStakeIBCKeeper.GetRedelegations(
		suite.ctx,
		suite.chainB.ChainID,
	)

	suite.Require().Equal(true, found)
	suite.Require().Equal(redelegations.ChainID, suite.chainB.ChainID)
	suite.Require().Equal(redelegations.Redelegations[0].ValidatorSrcAddress, redelegation.ValidatorSrcAddress)
	suite.Require().Equal(redelegations.Redelegations[0].ValidatorDstAddress, redelegation.ValidatorDstAddress)
}

func (suite *IntegrationTestSuite) TestAddRedelegationEntry() {

	suite.app.LiquidStakeIBCKeeper.SetRedelegations(
		suite.ctx, suite.chainB.ChainID,
		[]*stakingtypes.Redelegation{redelegation},
	)
	tc := []struct {
		name       string
		chainID    string
		msg        stakingtypes.MsgBeginRedelegate
		resp       stakingtypes.MsgBeginRedelegateResponse
		lenEntries uint32
		index      int
	}{
		{
			name:    "Success",
			chainID: suite.chainB.ChainID,
			msg: stakingtypes.MsgBeginRedelegate{
				DelegatorAddress:    "delAddr",
				ValidatorSrcAddress: "valSrcAddr",
				ValidatorDstAddress: "valDstAddr",
				Amount:              sdk.Coin{},
			},
			resp:       stakingtypes.MsgBeginRedelegateResponse{CompletionTime: time.Now().Add(time.Second)},
			lenEntries: 2,
			index:      0,
		},
		{
			name:    "success",
			chainID: suite.chainB.ChainID,
			msg: stakingtypes.MsgBeginRedelegate{
				DelegatorAddress:    "delAddr",
				ValidatorSrcAddress: "valSrcAddr",
				ValidatorDstAddress: "valDstAddr",
				Amount:              sdk.Coin{},
			},
			resp:       stakingtypes.MsgBeginRedelegateResponse{CompletionTime: time.Now().Add(time.Second)},
			lenEntries: 3,
			index:      0,
		}, {
			name:    "new chainid",
			chainID: suite.chainC.ChainID,
			msg: stakingtypes.MsgBeginRedelegate{
				DelegatorAddress:    "delAddr2",
				ValidatorSrcAddress: "valSrcAddr2",
				ValidatorDstAddress: "valDstAddr2",
				Amount:              sdk.Coin{},
			},
			resp:       stakingtypes.MsgBeginRedelegateResponse{CompletionTime: time.Now().Add(time.Second)},
			lenEntries: 1,
			index:      0,
		}, {
			name:    "new chainid",
			chainID: suite.chainC.ChainID,
			msg: stakingtypes.MsgBeginRedelegate{
				DelegatorAddress:    "delAddr2",
				ValidatorSrcAddress: "valSrcAddr3",
				ValidatorDstAddress: "valDstAddr2",
				Amount:              sdk.Coin{},
			},
			resp:       stakingtypes.MsgBeginRedelegateResponse{CompletionTime: time.Now().Add(time.Second)},
			lenEntries: 1,
			index:      0,
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {

			suite.app.LiquidStakeIBCKeeper.AddRedelegationEntry(
				suite.ctx,
				t.chainID,
				t.msg,
				t.resp,
			)

			redelegations, _ := suite.app.LiquidStakeIBCKeeper.GetRedelegations(
				suite.ctx,
				t.chainID,
			)

			suite.Require().Equal(t.lenEntries, uint32(len(redelegations.Redelegations[t.index].Entries)))
		})
	}
}
