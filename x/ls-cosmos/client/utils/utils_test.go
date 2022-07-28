package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewParamChangeJSON(t *testing.T) {
	rcj := NewRegisterChainJSON(
		"title",
		"description",
		"connection",
		"channel-1",
		"transfer",
		"uatom",
		"ustkatom",
		"1000stake",
	)
	require.Equal(t, "title", rcj.Title)
	require.Equal(t, "description", rcj.Description)
	require.Equal(t, "connection", rcj.IBCConnection)
	require.Equal(t, "channel-1", rcj.TokenTransferChannel)
	require.Equal(t, "transfer", rcj.TokenTransferPort)
	require.Equal(t, "uatom", rcj.BaseDenom)
	require.Equal(t, "ustkatom", rcj.MintDenom)
	require.Equal(t, "1000stake", rcj.Deposit)
}
