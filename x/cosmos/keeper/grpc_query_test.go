package keeper_test

import (
	gocontext "context"
	"github.com/cosmos/cosmos-sdk/baseapp"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/authz"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/persistenceOne/pstake-native/app/helpers"
	"github.com/persistenceOne/pstake-native/x/cosmos/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestKeeper_QueryTxByID(t *testing.T) {
	msg := stakingTypes.MsgUndelegate{
		DelegatorAddress: "address1",
		ValidatorAddress: "address2",
		Amount:           sdk.Coin{},
	}
	anyMsg, _ := codecTypes.NewAnyWithValue(&msg)

	execMsg := authz.MsgExec{
		Grantee: "address1",
		Msgs:    []*codecTypes.Any{anyMsg},
	}

	execMsgAny, _ := codecTypes.NewAnyWithValue(&execMsg)

	tx := types.CosmosTx{
		Tx: sdkTx.Tx{
			Body: &sdkTx.TxBody{
				Messages:      []*codecTypes.Any{execMsgAny},
				Memo:          "",
				TimeoutHeight: 0,
			},
			AuthInfo: &sdkTx.AuthInfo{
				SignerInfos: nil,
				Fee: &sdkTx.Fee{
					Amount:   nil,
					GasLimit: 200000,
					Payer:    "",
				},
			},
			Signatures: nil,
		},
	}

	_, app, ctx := helpers.CreateTestApp()
	keeper := app.CosmosKeeper
	keeper.SetNewTxnInOutgoingPool(ctx, 1, tx)
	tx1, _ := keeper.GetTxnFromOutgoingPoolByID(ctx, 1)
	err := tx1.CosmosTxDetails.Tx.UnpackInterfaces(app.AppCodec())
	require.NoError(t, err)

	_ = tx1.CosmosTxDetails.GetTx()

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.CosmosKeeper)
	queryClient := types.NewQueryClient(queryHelper)
	res, err := queryClient.QueryTxByID(gocontext.Background(), &types.QueryOutgoingTxByIDRequest{TxID: 1})
	require.NoError(t, err)

	getTx := res.CosmosTxDetails.GetTx()
	err = getTx.UnpackInterfaces(app.AppCodec())
	require.NoError(t, err)

	exec := getTx.GetMsgs()[0].(*authz.MsgExec)

	for _, im := range exec.Msgs {
		require.IsType(t, "string", im.TypeUrl)
	}
}
