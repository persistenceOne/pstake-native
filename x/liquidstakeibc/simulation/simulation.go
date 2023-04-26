package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	simappparams "github.com/persistenceOne/pstake-native/v2/app/params"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

// Simulation operation weights constants
const (
	OpWeightMsgDummy = "op_weight_msg_dummy" //nolint:gosec
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONCodec, ak types.AccountKeeper,
) simulation.WeightedOperations {
	var weightMsgDummy int
	appParams.GetOrGenerate(cdc, OpWeightMsgDummy, &weightMsgDummy, nil,
		func(_ *rand.Rand) {
			weightMsgDummy = simappparams.DefaultWeightMsgDummy
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgDummy,
			SimulateMsgDummy(ak),
		),
	}
}

// SimulateMsgDummy tests and runs a single msg dummy where both
// accounts already exist.
func SimulateMsgDummy(ak types.AccountKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		from, skip := randomDummyFields(r, ctx, accs, ak)

		if skip {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(&types.MsgDummy{}), "skip all txns"), nil, nil
		}

		msg := types.NewMsgDummy(from.Address)

		err := sendMsgDummy(r, app, ak, msg, ctx, chainID, []cryptotypes.PrivKey{from.PrivKey})
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "invalid txn"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// sendMsgDummy sends a transaction with a MsgDummy from a provided random account.
func sendMsgDummy(
	r *rand.Rand, app *baseapp.BaseApp, ak types.AccountKeeper,
	msg *types.MsgDummy, ctx sdk.Context, chainID string, privkeys []cryptotypes.PrivKey,
) error {
	var (
		fees sdk.Coins
		err  error
	)

	from, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return err
	}

	account := ak.GetAccount(ctx, from)

	txGen := simappparams.MakeEncodingConfig().TxConfig
	tx, err := helpers.GenSignedMockTx(
		r,
		txGen,
		[]sdk.Msg{msg},
		fees,
		helpers.DefaultGenTxGas,
		chainID,
		[]uint64{account.GetAccountNumber()},
		[]uint64{account.GetSequence()},
		privkeys...,
	)
	if err != nil {
		return err
	}

	_, _, err = app.SimDeliver(txGen.TxEncoder(), tx)
	if err != nil {
		return err
	}

	return nil
}

// randomDummyFields returns the sender
func randomDummyFields(
	r *rand.Rand, ctx sdk.Context, accs []simtypes.Account, ak types.AccountKeeper,
) (simtypes.Account, bool) {
	from, _ := simtypes.RandomAcc(r, accs)

	// disallow sending money to yourself

	acc := ak.GetAccount(ctx, from.Address)
	if acc == nil {
		return from, true
	}

	return from, false
}
