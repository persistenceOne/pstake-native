package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewParamChangeJSON(t *testing.T) {
	rcj := NewRegisterChainJSON(
		"title",
		"description",
		true,
		"cosmoshub-4",
		"connection",
		"channel-1",
		"transfer",
		"uatom",
		"ustkatom",
		"5",
		types.AllowListedValidators{
			AllowListedValidators: []types.AllowListedValidator{{
				ValidatorAddress: "Valaddr",
				TargetWeight:     sdk.OneDec(),
			}}},
		"0.0",
		"0.0",
		"0.0",
		"1000stake",
	)
	require.Equal(t, "title", rcj.Title)
	require.Equal(t, "description", rcj.Description)
	require.Equal(t, "cosmoshub-4", rcj.ChainID)
	require.Equal(t, "connection", rcj.ConnectionID)
	require.Equal(t, "channel-1", rcj.TransferChannel)
	require.Equal(t, "transfer", rcj.TransferPort)
	require.Equal(t, "uatom", rcj.BaseDenom)
	require.Equal(t, "ustkatom", rcj.MintDenom)
	require.Equal(t, "5", rcj.MinDeposit)
	require.Equal(t, "0.0", rcj.PstakeDepositFee)
	require.Equal(t, "0.0", rcj.PstakeRestakeFee)
	require.Equal(t, "0.0", rcj.PstakeUnstakeFee)
	require.Equal(t, "1000stake", rcj.Deposit)
}
