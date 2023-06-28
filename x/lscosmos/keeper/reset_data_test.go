package keeper_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/persistenceOne/pstake-native/v2/app"
	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/keeper"
	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
)

const (
	STARTING_EPOCH int64 = 200
	EPOCH_INTERVAL int64 = 4
	EPOCH_NUM      int64 = 10
	FUNDING_AMOUNT int64 = 1000000000

	ADDR1_UNDELEGATION int64 = 4500
	ADDR2_UNDELEGATION int64 = 7500
	ADDR3_UNDELEGATION int64 = 9000
	ADDR4_UNDELEGATION int64 = 9000

	DELEGATION_MODULE_ACCOUNT_ADDRESS string = "cosmos13t4996czrgft9gw43epuwauccrldu5whx6uprjdmvsmuf7ylg8yqcxgzk3"
)

var (
	Addr1 = sdk.AccAddress("test1_______________")
	Addr2 = sdk.AccAddress("test2_______________")
	Addr3 = sdk.AccAddress("test3_______________")
	Addr4 = sdk.AccAddress("test4_______________")

	FAILED_EPOCHS        = []int64{204, 216, 224}
	FORCE_TO_FAIL_EPOCHS = []int64{216}
	DELETED_EPOCHS       = []int64{232}
)

// test case where there is enough balance to cover claims in both protocol and module unbonding accounts.
func (suite *IntegrationTestSuite) TestForkSuccessFirstCase() {
	pstakeApp, ctx := suite.app, suite.ctx

	k, totalStkAtomsNeeded := PrepareTest(ctx, pstakeApp)

	// fund the protocol account
	protocolAccAddr := sdk.MustAccAddressFromBech32(keeper.PROTOCOL_ACC)
	pstakeApp.AccountKeeper.SetAccount(ctx, pstakeApp.AccountKeeper.NewAccountWithAddress(ctx, protocolAccAddr))
	suite.Require().NoError(
		simapp.FundAccount(
			pstakeApp.BankKeeper,
			ctx,
			protocolAccAddr,
			sdk.NewCoins(sdk.NewInt64Coin(k.GetHostChainParams(ctx).MintDenom, FUNDING_AMOUNT)),
		),
	)

	// fund the undelegation account
	suite.Require().NoError(
		simapp.FundAccount(
			pstakeApp.BankKeeper,
			ctx,
			authtypes.NewModuleAddress(types.UndelegationModuleAccount),
			sdk.NewCoins(sdk.NewInt64Coin(k.GetHostChainParams(ctx).MintDenom, FUNDING_AMOUNT)),
		),
	)

	err := k.Fork(ctx)
	suite.Require().NoError(err)

	// check the protocol balance hasn't been touched
	suite.Require().Equal(
		sdk.NewInt(FUNDING_AMOUNT),
		pstakeApp.BankKeeper.GetBalance(ctx, protocolAccAddr, k.GetHostChainParams(ctx).MintDenom).Amount,
	)

	// check force to fail epochs are marked as failed
	for _, epoch := range FORCE_TO_FAIL_EPOCHS {
		suite.Require().Equal(true, k.GetUnbondingEpochCValue(ctx, epoch).IsFailed)
	}

	for _, epoch := range DELETED_EPOCHS {
		// check unbonding epoch c value is not failed anymore
		suite.Require().Equal(false, k.GetUnbondingEpochCValue(ctx, epoch).IsFailed)

		// check host account undelegation for epoch has been created
		hostAccountUndelegationForEpoch, err := k.GetHostAccountUndelegationForEpoch(ctx, epoch)
		suite.Require().NoError(err)

		// check host account undelegation for epoch total amount is the stkATOM amount
		suite.Require().Equal(
			k.GetUnbondingEpochCValue(ctx, epoch).STKBurn,
			hostAccountUndelegationForEpoch.TotalUndelegationAmount,
		)
	}

	// SECOND PART OF THE TEST. THIS PART PERFORMS CHECKS ON THE STATE AFTER (SIMULATING)
	// CLEARING PACKETS ON THE RELAYER AND AFTER USERS CLAIM THEIR TOKENS.

	// simulate the packets being cleared in the relayer
	for _, epoch := range FAILED_EPOCHS {
		unbondingEpochCValue := k.GetUnbondingEpochCValue(ctx, epoch)
		unbondingEpochCValue.IsFailed = true
		k.SetUnbondingEpochCValue(ctx, unbondingEpochCValue)
	}

	// claim all the failed epochs
	for _, epoch := range FAILED_EPOCHS {
		err = k.ClaimFailed(ctx, types.DelegatorUnbondingEpochEntry{
			DelegatorAddress: Addr1.String(),
			EpochNumber:      epoch,
			Amount:           sdk.NewInt64Coin(k.GetHostChainParams(ctx).MintDenom, epoch*ADDR1_UNDELEGATION),
		})
		suite.Require().NoError(err)

		err = k.ClaimFailed(ctx, types.DelegatorUnbondingEpochEntry{
			DelegatorAddress: Addr2.String(),
			EpochNumber:      epoch,
			Amount:           sdk.NewInt64Coin(k.GetHostChainParams(ctx).MintDenom, epoch*ADDR2_UNDELEGATION),
		})
		suite.Require().NoError(err)

		err = k.ClaimFailed(ctx, types.DelegatorUnbondingEpochEntry{
			DelegatorAddress: Addr3.String(),
			EpochNumber:      epoch,
			Amount:           sdk.NewInt64Coin(k.GetHostChainParams(ctx).MintDenom, epoch*ADDR3_UNDELEGATION),
		})
		suite.Require().NoError(err)

		err = k.ClaimFailed(ctx, types.DelegatorUnbondingEpochEntry{
			DelegatorAddress: Addr4.String(),
			EpochNumber:      epoch,
			Amount:           sdk.NewInt64Coin(k.GetHostChainParams(ctx).MintDenom, epoch*ADDR3_UNDELEGATION),
		})
		suite.Require().NoError(err)
	}

	// check the undelegation module account has had the total stkATOM amount removed
	suite.Require().Equal(
		sdk.NewInt(FUNDING_AMOUNT).Sub(totalStkAtomsNeeded.Amount),
		pstakeApp.BankKeeper.GetBalance(
			ctx,
			authtypes.NewModuleAddress(types.UndelegationModuleAccount),
			k.GetHostChainParams(ctx).MintDenom,
		).Amount,
	)
}

// test case where there is enough balance to cover claims in the protocol account but not in the module unbonding account.
func (suite *IntegrationTestSuite) TestForkSuccessSecondCase() {
	pstakeApp, ctx := suite.app, suite.ctx

	k, totalStkAtomsNeeded := PrepareTest(ctx, pstakeApp)

	// fund the protocol account
	protocolAccAddr := sdk.MustAccAddressFromBech32(keeper.PROTOCOL_ACC)
	pstakeApp.AccountKeeper.SetAccount(ctx, pstakeApp.AccountKeeper.NewAccountWithAddress(ctx, protocolAccAddr))
	suite.Require().NoError(
		simapp.FundAccount(
			pstakeApp.BankKeeper,
			ctx,
			protocolAccAddr,
			sdk.NewCoins(sdk.NewInt64Coin(k.GetHostChainParams(ctx).MintDenom, FUNDING_AMOUNT)),
		),
	)

	err := k.Fork(ctx)
	suite.Require().NoError(err)

	// check the protocol balance has been deducted by the total stkATOM amount needed for claims
	suite.Require().Equal(
		sdk.NewInt(FUNDING_AMOUNT).Sub(totalStkAtomsNeeded.Amount),
		pstakeApp.BankKeeper.GetBalance(ctx, protocolAccAddr, k.GetHostChainParams(ctx).MintDenom).Amount,
	)

	// check force to fail epochs are marked as failed
	for _, epoch := range FORCE_TO_FAIL_EPOCHS {
		suite.Require().Equal(true, k.GetUnbondingEpochCValue(ctx, epoch).IsFailed)
	}

	for _, epoch := range DELETED_EPOCHS {
		// check unbonding epoch c value is not failed anymore
		suite.Require().Equal(false, k.GetUnbondingEpochCValue(ctx, epoch).IsFailed)

		// check host account undelegation for epoch has been created
		hostAccountUndelegationForEpoch, err := k.GetHostAccountUndelegationForEpoch(ctx, epoch)
		suite.Require().NoError(err)

		// check host account undelegation for epoch total amount is the stkATOM amount
		suite.Require().Equal(
			k.GetUnbondingEpochCValue(ctx, epoch).STKBurn,
			hostAccountUndelegationForEpoch.TotalUndelegationAmount,
		)
	}

	// SECOND PART OF THE TEST. THIS PART PERFORMS CHECKS ON THE STATE AFTER (SIMULATING)
	// CLEARING PACKETS ON THE RELAYER AND AFTER USERS CLAIM THEIR TOKENS.

	// simulate the packets being cleared in the relayer
	for _, epoch := range FAILED_EPOCHS {
		unbondingEpochCValue := k.GetUnbondingEpochCValue(ctx, epoch)
		unbondingEpochCValue.IsFailed = true
		k.SetUnbondingEpochCValue(ctx, unbondingEpochCValue)
	}

	// claim all the failed epochs
	for _, epoch := range FAILED_EPOCHS {
		err = k.ClaimFailed(ctx, types.DelegatorUnbondingEpochEntry{
			DelegatorAddress: Addr1.String(),
			EpochNumber:      epoch,
			Amount:           sdk.NewInt64Coin(k.GetHostChainParams(ctx).MintDenom, epoch*ADDR1_UNDELEGATION),
		})
		suite.Require().NoError(err)

		err = k.ClaimFailed(ctx, types.DelegatorUnbondingEpochEntry{
			DelegatorAddress: Addr2.String(),
			EpochNumber:      epoch,
			Amount:           sdk.NewInt64Coin(k.GetHostChainParams(ctx).MintDenom, epoch*ADDR2_UNDELEGATION),
		})
		suite.Require().NoError(err)

		err = k.ClaimFailed(ctx, types.DelegatorUnbondingEpochEntry{
			DelegatorAddress: Addr3.String(),
			EpochNumber:      epoch,
			Amount:           sdk.NewInt64Coin(k.GetHostChainParams(ctx).MintDenom, epoch*ADDR3_UNDELEGATION),
		})
		suite.Require().NoError(err)

		err = k.ClaimFailed(ctx, types.DelegatorUnbondingEpochEntry{
			DelegatorAddress: Addr4.String(),
			EpochNumber:      epoch,
			Amount:           sdk.NewInt64Coin(k.GetHostChainParams(ctx).MintDenom, epoch*ADDR3_UNDELEGATION),
		})
		suite.Require().NoError(err)
	}

	// check the undelegation module account has just used the required amount for claims
	suite.Require().Equal(
		sdk.ZeroInt(),
		pstakeApp.BankKeeper.GetBalance(
			ctx,
			authtypes.NewModuleAddress(types.UndelegationModuleAccount),
			k.GetHostChainParams(ctx).MintDenom,
		).Amount,
	)
}

// test case where there is not enough balance to cover claims in the protocol account but there is in the module unbonding account.
func (suite *IntegrationTestSuite) TestForkSuccessThirdCaseCase() {
	pstakeApp, ctx := suite.app, suite.ctx

	k, totalStkAtomsNeeded := PrepareTest(ctx, pstakeApp)

	// fund the undelegation account
	suite.Require().NoError(
		simapp.FundAccount(
			pstakeApp.BankKeeper,
			ctx,
			authtypes.NewModuleAddress(types.UndelegationModuleAccount),
			sdk.NewCoins(sdk.NewInt64Coin(k.GetHostChainParams(ctx).MintDenom, FUNDING_AMOUNT)),
		),
	)

	err := k.Fork(ctx)
	suite.Require().NoError(err)

	// check the protocol balance hasn't been touched
	protocolAccAddr := sdk.MustAccAddressFromBech32(keeper.PROTOCOL_ACC)
	suite.Require().Equal(
		sdk.ZeroInt(),
		pstakeApp.BankKeeper.GetBalance(ctx, protocolAccAddr, k.GetHostChainParams(ctx).MintDenom).Amount,
	)

	// check force to fail epochs are marked as failed
	for _, epoch := range FORCE_TO_FAIL_EPOCHS {
		suite.Require().Equal(true, k.GetUnbondingEpochCValue(ctx, epoch).IsFailed)
	}

	for _, epoch := range DELETED_EPOCHS {
		// check unbonding epoch c value is not failed anymore
		suite.Require().Equal(false, k.GetUnbondingEpochCValue(ctx, epoch).IsFailed)

		// check host account undelegation for epoch has been created
		hostAccountUndelegationForEpoch, err := k.GetHostAccountUndelegationForEpoch(ctx, epoch)
		suite.Require().NoError(err)

		// check host account undelegation for epoch total amount is the stkATOM amount
		suite.Require().Equal(
			k.GetUnbondingEpochCValue(ctx, epoch).STKBurn,
			hostAccountUndelegationForEpoch.TotalUndelegationAmount,
		)
	}

	// SECOND PART OF THE TEST. THIS PART PERFORMS CHECKS ON THE STATE AFTER (SIMULATING)
	// CLEARING PACKETS ON THE RELAYER AND AFTER USERS CLAIM THEIR TOKENS.

	// simulate the packets being cleared in the relayer
	for _, epoch := range FAILED_EPOCHS {
		unbondingEpochCValue := k.GetUnbondingEpochCValue(ctx, epoch)
		unbondingEpochCValue.IsFailed = true
		k.SetUnbondingEpochCValue(ctx, unbondingEpochCValue)
	}

	// claim all the failed epochs
	for _, epoch := range FAILED_EPOCHS {
		err = k.ClaimFailed(ctx, types.DelegatorUnbondingEpochEntry{
			DelegatorAddress: Addr1.String(),
			EpochNumber:      epoch,
			Amount:           sdk.NewInt64Coin(k.GetHostChainParams(ctx).MintDenom, epoch*ADDR1_UNDELEGATION),
		})
		suite.Require().NoError(err)

		err = k.ClaimFailed(ctx, types.DelegatorUnbondingEpochEntry{
			DelegatorAddress: Addr2.String(),
			EpochNumber:      epoch,
			Amount:           sdk.NewInt64Coin(k.GetHostChainParams(ctx).MintDenom, epoch*ADDR2_UNDELEGATION),
		})
		suite.Require().NoError(err)

		err = k.ClaimFailed(ctx, types.DelegatorUnbondingEpochEntry{
			DelegatorAddress: Addr3.String(),
			EpochNumber:      epoch,
			Amount:           sdk.NewInt64Coin(k.GetHostChainParams(ctx).MintDenom, epoch*ADDR3_UNDELEGATION),
		})
		suite.Require().NoError(err)

		err = k.ClaimFailed(ctx, types.DelegatorUnbondingEpochEntry{
			DelegatorAddress: Addr4.String(),
			EpochNumber:      epoch,
			Amount:           sdk.NewInt64Coin(k.GetHostChainParams(ctx).MintDenom, epoch*ADDR3_UNDELEGATION),
		})
		suite.Require().NoError(err)
	}

	// check the undelegation module account has had the total stkATOM amount removed
	suite.Require().Equal(
		sdk.NewInt(FUNDING_AMOUNT).Sub(totalStkAtomsNeeded.Amount),
		pstakeApp.BankKeeper.GetBalance(
			ctx,
			authtypes.NewModuleAddress(types.UndelegationModuleAccount),
			k.GetHostChainParams(ctx).MintDenom,
		).Amount,
	)
}

// test case where there is not enough balance to cover claims in both protocol and module unbonding accounts.
func (suite *IntegrationTestSuite) TestForkUnsuccessful() {
	pstakeApp, ctx := suite.app, suite.ctx

	k, _ := PrepareTest(ctx, pstakeApp)

	suite.Require().Panics(func() { k.Fork(ctx) })
}

func PrepareTest(ctx sdk.Context, app *app.PstakeApp) (k keeper.Keeper, totalStkAtomsNeeded sdk.Coin) {
	k = app.LSCosmosKeeper
	k.GetUndelegationModuleAccount(ctx)
	allowListedVals := k.GetAllowListedValidators(ctx).AllowListedValidators //3 validators
	hcp := k.GetHostChainParams(ctx)

	// create the delegations, total delegations 600ATOM (600.000.000 uatom), 200.000.000 uatom per validator
	var hostAccountDelegations []types.HostAccountDelegation
	for _, val := range allowListedVals {
		hostAccountDelegations = append(hostAccountDelegations,
			types.NewHostAccountDelegation(val.ValidatorAddress,
				sdk.NewInt64Coin(hcp.BaseDenom, 200000000),
			),
		)
	}

	// create the undelegations
	// there will be 10 epochs in the test, which of 5 will be corrupted and 5 will not
	var hostAccountUndelegations []types.HostAccountUndelegation
	for epoch := STARTING_EPOCH; epoch < STARTING_EPOCH+(EPOCH_INTERVAL*EPOCH_NUM); epoch += EPOCH_INTERVAL {
		// create the undelegation entries
		var undelegationEntries []types.UndelegationEntry
		for _, val := range allowListedVals {
			undelegationEntries = append(undelegationEntries, types.UndelegationEntry{
				ValidatorAddress: val.ValidatorAddress,
				Amount:           sdk.NewInt64Coin(hcp.BaseDenom, epoch*10000), // undelegate starting at 2.000.000 uatom
			})
		}

		// create the host account undelegations
		atoms := sdk.NewInt64Coin(hcp.BaseDenom, epoch*10000*3)
		stkAtoms := sdk.NewInt64Coin(hcp.MintDenom, epoch*10000*3)
		completionTime := 24 * (time.Duration(epoch) - 200)
		hostAccountUndelegations = append(hostAccountUndelegations, types.HostAccountUndelegation{
			EpochNumber:             epoch,
			TotalUndelegationAmount: stkAtoms,
			CompletionTime:          time.Now().Add(time.Hour * completionTime), // 0 days, 4 days, 8days etc..
			UndelegationEntries:     undelegationEntries,
		})

		// create the unbonding epoch c value
		k.SetUnbondingEpochCValue(ctx, types.UnbondingEpochCValue{
			EpochNumber:    epoch,
			STKBurn:        stkAtoms,
			AmountUnbonded: atoms,
			IsMatured:      false,
			IsFailed:       false,
		})

		// create the delegator unbonding entries, they total at epoch*10.000 atom, which is the total being undelegated
		k.SetDelegatorUnbondingEpochEntry(ctx, types.DelegatorUnbondingEpochEntry{
			DelegatorAddress: Addr1.String(),
			EpochNumber:      epoch,
			Amount:           sdk.NewInt64Coin(hcp.MintDenom, epoch*ADDR1_UNDELEGATION),
		})
		k.SetDelegatorUnbondingEpochEntry(ctx, types.DelegatorUnbondingEpochEntry{
			DelegatorAddress: Addr2.String(),
			EpochNumber:      epoch,
			Amount:           sdk.NewInt64Coin(hcp.MintDenom, epoch*ADDR2_UNDELEGATION),
		})
		k.SetDelegatorUnbondingEpochEntry(ctx, types.DelegatorUnbondingEpochEntry{
			DelegatorAddress: Addr3.String(),
			EpochNumber:      epoch,
			Amount:           sdk.NewInt64Coin(hcp.MintDenom, epoch*ADDR3_UNDELEGATION),
		})
		k.SetDelegatorUnbondingEpochEntry(ctx, types.DelegatorUnbondingEpochEntry{
			DelegatorAddress: Addr4.String(),
			EpochNumber:      epoch,
			Amount:           sdk.NewInt64Coin(hcp.MintDenom, epoch*ADDR4_UNDELEGATION),
		})
	}

	// prepare the failed epochs mature time
	totalStkAtomsNeeded = sdk.NewCoin(hcp.MintDenom, sdk.ZeroInt())
	for _, epoch := range FAILED_EPOCHS {
		hostAccountUndelegations[(epoch%200)/4].CompletionTime = time.Time{}
		totalStkAtomsNeeded.Amount = totalStkAtomsNeeded.Amount.Add(
			hostAccountUndelegations[(epoch%200)/4].TotalUndelegationAmount.Amount,
		)
	}

	// remove the host account undelegation for the deleted epochs
	for _, epoch := range DELETED_EPOCHS {
		hostAccountUndelegations = append(hostAccountUndelegations[:(epoch%200)/4], hostAccountUndelegations[(epoch%200)/4+1:]...)
		k.FailUnbondingEpochCValue(ctx, epoch, sdk.NewInt64Coin(hcp.MintDenom, 0))
	}

	// create the delegation state of the module
	k.SetDelegationState(ctx, types.DelegationState{
		HostDelegationAccountBalance: sdk.NewCoins(sdk.NewCoin(hcp.BaseDenom, sdk.NewInt(1000000000))),
		HostChainDelegationAddress:   DELEGATION_MODULE_ACCOUNT_ADDRESS,
		HostAccountDelegations:       hostAccountDelegations,
		HostAccountUndelegations:     hostAccountUndelegations,
	})

	return
}

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
