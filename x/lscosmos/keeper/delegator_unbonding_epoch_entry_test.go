package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

func (suite *IntegrationTestSuite) TestKVStore() {
	app, ctx := suite.app, suite.ctx

	addr1, err := sdk.AccAddressFromBech32("persistence1826wkxx8wv7mfnank8l6xu9rxm7kg8rvvk4e0a")
	suite.NoError(err)

	addr2, err := sdk.AccAddressFromBech32("persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9")
	suite.NoError(err)

	addr3, err := sdk.AccAddressFromBech32("persistence1lngwr8ymx3q6gtsff2h8407mawz9azp6kmut02")
	suite.NoError(err)

	testCases := []struct {
		address                                                   sdk.AccAddress
		delegatorUnbondingEpochEntry, expectedUnbondingEpochEntry types.DelegatorUnbondingEpochEntry
		repeatEntry                                               int
	}{
		{
			address: addr1,
			delegatorUnbondingEpochEntry: types.DelegatorUnbondingEpochEntry{
				DelegatorAddress: addr1.String(),
				EpochNumber:      1,
				Amount:           sdk.NewInt64Coin("stkAtom", 10000),
			},
			expectedUnbondingEpochEntry: types.DelegatorUnbondingEpochEntry{
				DelegatorAddress: addr1.String(),
				EpochNumber:      1,
				Amount:           sdk.NewInt64Coin("stkAtom", 20000),
			},
			repeatEntry: 2,
		},
		{
			address: addr2,
			delegatorUnbondingEpochEntry: types.DelegatorUnbondingEpochEntry{
				DelegatorAddress: addr2.String(),
				EpochNumber:      2,
				Amount:           sdk.NewInt64Coin("stkAtom", 10000),
			},
			expectedUnbondingEpochEntry: types.DelegatorUnbondingEpochEntry{
				DelegatorAddress: addr2.String(),
				EpochNumber:      2,
				Amount:           sdk.NewInt64Coin("stkAtom", 10000),
			},
			repeatEntry: 1,
		},
		{
			address: addr3,
			delegatorUnbondingEpochEntry: types.DelegatorUnbondingEpochEntry{
				DelegatorAddress: addr3.String(),
				EpochNumber:      1,
				Amount:           sdk.NewInt64Coin("stkAtom", 10000),
			},
			expectedUnbondingEpochEntry: types.DelegatorUnbondingEpochEntry{
				DelegatorAddress: addr3.String(),
				EpochNumber:      1,
				Amount:           sdk.NewInt64Coin("stkAtom", 40000),
			},
			repeatEntry: 4,
		},
	}

	keeper := app.LSCosmosKeeper

	for _, tc := range testCases {

		for i := 1; i <= tc.repeatEntry; i++ {
			keeper.AddDelegatorUnbondingEpochEntry(
				ctx,
				tc.address,
				tc.delegatorUnbondingEpochEntry.EpochNumber,
				tc.delegatorUnbondingEpochEntry.Amount)
		}

		list := keeper.IterateDelegatorUnbondingEpochEntry(ctx, tc.address)
		suite.Equal(1, len(list))
		suite.Equal(tc.expectedUnbondingEpochEntry.Amount, list[0].Amount)
		suite.Equal(tc.expectedUnbondingEpochEntry.DelegatorAddress, list[0].DelegatorAddress)
		suite.Equal(tc.expectedUnbondingEpochEntry.EpochNumber, list[0].EpochNumber)

		entry := keeper.GetDelegatorUnbondingEpochEntry(ctx, tc.address, list[0].EpochNumber)
		suite.Equal(tc.expectedUnbondingEpochEntry.Amount, entry.Amount)
		suite.Equal(tc.expectedUnbondingEpochEntry.DelegatorAddress, entry.DelegatorAddress)
		suite.Equal(tc.expectedUnbondingEpochEntry.EpochNumber, entry.EpochNumber)

		keeper.RemoveDelegatorUnbondingEpochEntry(ctx, tc.address, list[0].EpochNumber)
		entry = keeper.GetDelegatorUnbondingEpochEntry(ctx, tc.address, list[0].EpochNumber)
		suite.Equal(int64(0), entry.EpochNumber)
	}
}
