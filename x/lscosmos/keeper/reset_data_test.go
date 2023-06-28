package keeper_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/simapp"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/keeper"
	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
)

// ////////// fetch data from json:
// curl -X GET -H "Content-Type: application/json" -H "x-cosmos-block-height: 15884400" "https://rest.cosmos.audit.one/cosmos/staking/v1beta1/delegators/cosmos13t4996czrgft9gw43epuwauccrldu5whx6uprjdmvsmuf7ylg8yqcxgzk3/unbonding_delegations"
// creationheight 15826532 => epoch 232, 15881865 => 236

func TestParseHostAccountUnbondings(t *testing.T) {
	mintDenom := "stk/uatom"
	baseDenom := "uatom"
	// create a map to quickly access each undelegation epoch entry and initialise it
	hostAccountUndelegationsMap := keeper.ParseHostAccountUnbondings(mintDenom, baseDenom)
	fmt.Println(hostAccountUndelegationsMap)
}

var (
	Addr1 = sdk.AccAddress("test1_______________")
	Addr2 = sdk.AccAddress("test2_______________")
	Addr3 = sdk.AccAddress("test3_______________")
	Addr4 = sdk.AccAddress("test4_______________")
)

func (suite *IntegrationTestSuite) TestFork() {
	pstakeApp, ctx := suite.app, suite.ctx
	k := pstakeApp.LSCosmosKeeper
	k.GetUndelegationModuleAccount(ctx)
	allowListedVals := k.GetAllowListedValidators(ctx) //3 validators
	hcp := k.GetHostChainParams(ctx)

	var hostAccountDelegations []types.HostAccountDelegation
	for _, val := range allowListedVals.AllowListedValidators {
		hostAccountDelegations = append(hostAccountDelegations,
			types.NewHostAccountDelegation(val.ValidatorAddress,
				sdk.NewInt64Coin(hcp.BaseDenom, 234098780),
			),
		)
	}
	var hostAccountUndelegations []types.HostAccountUndelegation
	for i := 200; i < 244; i = i + 4 {
		var undelegationEntries []types.UndelegationEntry
		for _, val := range allowListedVals.AllowListedValidators {
			undelegationEntries = append(undelegationEntries, types.UndelegationEntry{
				ValidatorAddress: val.ValidatorAddress,
				Amount:           sdk.NewInt64Coin(hcp.BaseDenom, int64(i)*920780), //any amount is fine.
			})
		}
		stkAtoms := sdk.NewInt64Coin(hcp.MintDenom, 908780*3*int64(i))
		completionTime := 24 * (time.Duration(i) - 200)
		hostAccountUndelegations = append(hostAccountUndelegations, types.HostAccountUndelegation{
			EpochNumber:             int64(i),
			TotalUndelegationAmount: stkAtoms,
			CompletionTime:          time.Now().Add(time.Hour * completionTime), // 0 days, 4 days, 8days etc..
			UndelegationEntries:     undelegationEntries,
		})
		atoms := sdk.NewInt64Coin(hcp.MintDenom, 920780*3*int64(i))

		k.SetUnbondingEpochCValue(ctx, types.UnbondingEpochCValue{
			EpochNumber:    int64(i),
			STKBurn:        stkAtoms,
			AmountUnbonded: atoms,
			IsMatured:      false,
			IsFailed:       false,
		})
		k.SetDelegatorUnbondingEpochEntry(ctx, types.DelegatorUnbondingEpochEntry{
			DelegatorAddress: Addr1.String(),
			EpochNumber:      int64(i),
			Amount:           sdk.NewInt64Coin(hcp.MintDenom, 908780*int64(i)),
		})
		k.SetDelegatorUnbondingEpochEntry(ctx, types.DelegatorUnbondingEpochEntry{
			DelegatorAddress: Addr2.String(),
			EpochNumber:      int64(i),
			Amount:           sdk.NewInt64Coin(hcp.MintDenom, 908780*int64(i)),
		})
		k.SetDelegatorUnbondingEpochEntry(ctx, types.DelegatorUnbondingEpochEntry{
			DelegatorAddress: Addr3.String(),
			EpochNumber:      int64(i),
			Amount:           sdk.NewInt64Coin(hcp.MintDenom, 508780*int64(i)),
		})
		k.SetDelegatorUnbondingEpochEntry(ctx, types.DelegatorUnbondingEpochEntry{
			DelegatorAddress: Addr4.String(),
			EpochNumber:      int64(i),
			Amount:           sdk.NewInt64Coin(hcp.MintDenom, 400000*int64(i)),
		})
	}
	// Remove completion time for epoch with index 1, 4, 7
	failedUndelegation1 := hostAccountUndelegations[1]
	hostAccountUndelegations[1].CompletionTime = time.Time{}
	failedUndelegation2 := hostAccountUndelegations[4]
	hostAccountUndelegations[4].CompletionTime = time.Time{}
	failedUndelegation3 := hostAccountUndelegations[7]
	hostAccountUndelegations[7].CompletionTime = time.Time{}
	// delete epoch 8
	deletedUndelegation := hostAccountUndelegations[8]
	hostAccountUndelegations = append(hostAccountUndelegations[:8], hostAccountUndelegations[8+1:]...)
	k.FailUnbondingEpochCValue(ctx, deletedUndelegation.EpochNumber, sdk.NewInt64Coin(hcp.MintDenom, 0))

	k.SetDelegationState(ctx, types.DelegationState{
		HostDelegationAccountBalance: sdk.NewCoins(sdk.NewCoin(hcp.BaseDenom, sdk.NewInt(1000000000))),
		HostChainDelegationAddress:   "cosmos13t4996czrgft9gw43epuwauccrldu5whx6uprjdmvsmuf7ylg8yqcxgzk3",
		HostAccountDelegations:       hostAccountDelegations,
		HostAccountUndelegations:     hostAccountUndelegations,
	})

	// add amount to protocol account
	fundProtocol := sdk.NewCoins(sdk.NewInt64Coin(hcp.MintDenom, 1000000000))
	protocolAccAddr := sdk.MustAccAddressFromBech32(keeper.PROTOCOL_ACC)
	acc := pstakeApp.AccountKeeper.NewAccountWithAddress(ctx, protocolAccAddr)
	pstakeApp.AccountKeeper.SetAccount(ctx, acc)
	suite.Require().NoError(simapp.FundAccount(pstakeApp.BankKeeper, ctx, protocolAccAddr, fundProtocol))

	// add amount to undelegation account
	stkatomswithModule := failedUndelegation1.TotalUndelegationAmount.Add(failedUndelegation2.TotalUndelegationAmount).Add(failedUndelegation3.TotalUndelegationAmount)

	suite.Require().NoError(simapp.FundAccount(pstakeApp.BankKeeper, ctx, authtypes.NewModuleAddress(types.UndelegationModuleAccount), sdk.NewCoins(stkatomswithModule)))

	err := k.ClaimFailed(ctx, types.DelegatorUnbondingEpochEntry{
		DelegatorAddress: Addr1.String(),
		EpochNumber:      deletedUndelegation.EpochNumber,
		Amount: sdk.NewCoin(
			deletedUndelegation.TotalUndelegationAmount.Denom,
			deletedUndelegation.TotalUndelegationAmount.Amount.QuoRaw(3),
		), // 1/3rd undelegation is always Addr1
	})
	suite.Require().NoError(err)

	err = k.Fork(ctx)
	suite.Require().NoError(err)
}
