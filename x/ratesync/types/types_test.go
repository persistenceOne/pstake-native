package types

import (
	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTypes(t *testing.T) {
	hc := ValidHostChainInMsg(0)
	require.NoError(t, hc.ValidateBasic())

	hc2 := ValidHostChainInMsg(0)
	hc2.ConnectionID = "invalid chars $%&/"
	require.Error(t, hc2.ValidateBasic())

	hc2 = ValidHostChainInMsg(0)
	hc2.ICAAccount.Owner = DefaultPortOwner(1)
	require.NoError(t, hc2.ValidateBasic())
	hc2.ICAAccount.Owner = "anything else"
	require.Error(t, hc2.ValidateBasic())
	hc2.ICAAccount.Owner = "pstake_ratesync_notint"
	require.Error(t, hc2.ValidateBasic())

	hc2 = ValidHostChainInMsg(0)
	hc2.Features.LiquidStake.Enabled = true
	require.Error(t, hc2.ValidateBasic())
	hc2.ICAAccount.Address = "someAddr"
	require.Error(t, hc2.ValidateBasic())

	hc2 = ValidHostChainInMsg(0)
	hc2.ICAAccount.ChannelState = liquidstakeibctypes.ICAAccount_ICA_CHANNEL_CREATED
	require.Error(t, hc2.ValidateBasic())
	hc2.ICAAccount.Address = "someAddr"
	require.Error(t, hc2.ValidateBasic())

	hc2 = ValidHostChainInMsg(0)
	hc2.Features.LiquidStake.FeatureType = FeatureType_LIQUID_STAKE_IBC
	require.Error(t, hc2.ValidateBasic())

	require.False(t, hc2.IsActive())
	hc2.Features.LiquidStake.Enabled = true
	require.True(t, hc2.IsActive())

}
