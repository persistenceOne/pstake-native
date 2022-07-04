package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/x/cosmos/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewMsgWithdrawStkAsset(t *testing.T) {
	//_, app, ctx := helpers.CreateTestApp()
	//keeper := app.CosmosKeeper

	name := "persistenceValidator"
	valAddress, _ := sdk.ValAddressFromBech32(name)
	accAddress, _ := sdk.AccAddressFromBech32(name)

	msg := types.NewMsgSetOrchestrator(valAddress, accAddress)
	require.NoError(t, nil, msg.ValidateBasic())

	msg1 := types.NewMsgRemoveOrchestrator(valAddress, accAddress)
	require.NoError(t, nil, msg1.ValidateBasic())

	msg2 := types.NewMsgWithdrawStkAsset(accAddress, accAddress, sdk.NewCoin("test", sdk.NewInt(0)))
	require.NoError(t, nil, msg2.ValidateBasic())

	msg3 := types.NewMsgSetSignature(accAddress, 1, []byte{}, 1)
	require.NoError(t, nil, msg3.ValidateBasic())
}
