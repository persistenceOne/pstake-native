package keeper_test

import (
	multisig2 "github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/persistenceOne/pstake-native/app"
	"github.com/persistenceOne/pstake-native/app/helpers"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	valAddr1         = sdkTypes.ValAddress("Val1")
	valAddr2         = sdkTypes.ValAddress("Val2")
	valAddrNotExists = sdkTypes.ValAddress("ValNotExists")
	valAddrInvalid   = sdkTypes.ValAddress(make([]byte, address.MaxAddrLen+1))

	orchAddr1       = sdkTypes.AccAddress("orch1")
	orchAddr2       = sdkTypes.AccAddress("orch2,0")
	orchAddr21      = sdkTypes.AccAddress("orch2,1")
	orchAddrInvalid = sdkTypes.AccAddress(make([]byte, address.MaxAddrLen+1))
)

func SetupAccounts(t *testing.T, app app.PstakeApp, ctx sdkTypes.Context) {
	//Adds incorrect pubkeys so that conversion doesnt happen
	accountOrch1 := app.AccountKeeper.NewAccountWithAddress(ctx, orchAddr1)
	err := accountOrch1.SetPubKey(&secp256k1.PubKey{Key: orchAddr1})
	require.Nil(t, err)
	accountOrch2 := app.AccountKeeper.NewAccountWithAddress(ctx, orchAddr2)
	err = accountOrch2.SetPubKey(&secp256k1.PubKey{Key: orchAddr2})
	require.Nil(t, err)

	app.AccountKeeper.SetAccount(ctx, accountOrch1)
	app.AccountKeeper.SetAccount(ctx, accountOrch2)
}
func SetupValidators(t *testing.T, app app.PstakeApp, ctx sdkTypes.Context) {

	app.StakingKeeper.SetValidator(ctx, stakingTypes.Validator{
		OperatorAddress: valAddr1.String(),
	})
	app.StakingKeeper.SetValidator(ctx, stakingTypes.Validator{
		OperatorAddress: valAddr2.String(),
	})
	err := app.CosmosKeeper.SetValidatorOrchestrator(ctx, valAddr1, orchAddr1)
	require.Nil(t, err, "Could not set valAddr1")

	err = app.CosmosKeeper.SetValidatorOrchestrator(ctx, valAddr2, orchAddr2)
	require.Nil(t, err, "Could not set valAddr2")

}

func TestTest(t *testing.T) {
	//_, pstakeApp, ctx := helpers.CreateTestInput()
	//accountOrch1 := pstakeApp.AccountKeeper.NewAccountWithAddress(ctx, orchAddr1)
	//_ = accountOrch1.SetPubKey(&secp256k1.PubKey{Key: orchAddr1})
	//accountOrch2 := pstakeApp.AccountKeeper.NewAccountWithAddress(ctx, orchAddr2)
	//_ = accountOrch2.SetPubKey(&secp256k1.PubKey{Key: orchAddr2})
	//
	//require.Equal(t, accountOrch2.GetPubKey().Address(), orchAddr1.Bytes())
}

func TestCheckValidator(t *testing.T) {
	_, pstakeApp, ctx := helpers.CreateTestApp()
	keeper := pstakeApp.CosmosKeeper

	SetupValidators(t, pstakeApp, ctx)

	val1, ok := keeper.CheckValidator(ctx, valAddr1)
	require.Equal(t, valAddr1, val1)
	require.Equal(t, true, ok)

	valNotExists, ok := keeper.CheckValidator(ctx, valAddrNotExists)
	require.Nil(t, valNotExists)
	require.Equal(t, false, ok)

	valInvalid, ok := keeper.CheckValidator(ctx, valAddrInvalid)
	require.Nil(t, valInvalid)
	require.Equal(t, false, ok)

}

func TestSetValidatorOrchestrator(t *testing.T) {
	_, pstakeApp, ctx := helpers.CreateTestApp()
	keeper := pstakeApp.CosmosKeeper
	require.Panics(t, func() { _ = keeper.SetValidatorOrchestrator(ctx, valAddr1, orchAddrInvalid) })
	require.Panics(t, func() { _ = keeper.SetValidatorOrchestrator(ctx, valAddrInvalid, orchAddr1) })
	require.Error(t, keeper.SetValidatorOrchestrator(ctx, valAddrNotExists, orchAddr1))
	SetupValidators(t, pstakeApp, ctx)
	require.Error(t, keeper.SetValidatorOrchestrator(ctx, valAddr1, orchAddr1))
	require.Nil(t, keeper.SetValidatorOrchestrator(ctx, valAddr2, orchAddr21))
}

func TestGetTotalValidatorOrchestratorCount(t *testing.T) {
	_, pstakeApp, ctx := helpers.CreateTestApp()
	keeper := pstakeApp.CosmosKeeper
	count := keeper.GetTotalValidatorOrchestratorCount(ctx)
	require.Equal(t, int64(0), count)

	SetupValidators(t, pstakeApp, ctx)
	count = keeper.GetTotalValidatorOrchestratorCount(ctx)
	require.Equal(t, int64(2), count)
}

func TestRemoveValidatorOrchestrator(t *testing.T) {
	_, pstakeApp, ctx := helpers.CreateTestApp()
	keeper := pstakeApp.CosmosKeeper

	// inconsistent with SetValidatorOrchestrator // FIX this //TODO
	require.Error(t, keeper.RemoveValidatorOrchestrator(ctx, valAddr1, orchAddrInvalid), "Orch address invalid")
	require.Error(t, keeper.RemoveValidatorOrchestrator(ctx, valAddrInvalid, orchAddr1), "ValidatorAddress invalid")
	require.Panics(t, func() { _ = keeper.RemoveValidatorOrchestrator(ctx, valAddr1, orchAddr1) }, "pub key for orch address not found")
	SetupAccounts(t, pstakeApp, ctx)
	require.Panics(t, func() { _ = keeper.RemoveValidatorOrchestrator(ctx, valAddr1, orchAddr1) }, "multisig is not set")
	SetupValidators(t, pstakeApp, ctx)
	require.Panics(t, func() { _ = keeper.RemoveValidatorOrchestrator(ctx, valAddr1, orchAddr1) }, "multisig is not set")

	pubkeyOrch1 := pstakeApp.AccountKeeper.GetAccount(ctx, orchAddr1).GetPubKey()
	multisigPubkey := multisig2.NewLegacyAminoPubKey(int(1), []types.PubKey{pubkeyOrch1})
	multisigAddr := sdkTypes.AccAddress("multisigAddr")
	multiSigAcc := pstakeApp.AccountKeeper.NewAccountWithAddress(ctx, multisigAddr)
	err := multiSigAcc.SetPubKey(multisigPubkey)
	require.Nil(t, err)
	pstakeApp.CosmosKeeper.SetAccountState(ctx, multiSigAcc)
	pstakeApp.CosmosKeeper.SetCurrentAddress(ctx, multisigAddr)

	require.Error(t, keeper.RemoveValidatorOrchestrator(ctx, valAddr1, orchAddr1), "OrchAddress present in multisig")

	require.Error(t, keeper.RemoveValidatorOrchestrator(ctx, valAddrNotExists, orchAddr2), "validator does not exist")
	require.Error(t, keeper.RemoveValidatorOrchestrator(ctx, valAddr2, orchAddr2), "cannot remove the only orch-validator mapping")

	accountOrch21 := pstakeApp.AccountKeeper.NewAccountWithAddress(ctx, orchAddr21)
	err = accountOrch21.SetPubKey(&secp256k1.PubKey{Key: orchAddr21})
	require.Nil(t, err)

	pstakeApp.AccountKeeper.SetAccount(ctx, accountOrch21)
	require.Nil(t, keeper.SetValidatorOrchestrator(ctx, valAddr2, orchAddr21), "--setup")

	require.Nil(t, keeper.RemoveValidatorOrchestrator(ctx, valAddr2, orchAddr21), "")
}
