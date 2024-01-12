package types

import (
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"

	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func TestTypes(t *testing.T) {
	// hostchain
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

	// features
	features := ValidHostChainInMsg(0).Features
	require.NoError(t, features.ValdidateBasic())
	features.LiquidStakeIBC.FeatureType = FeatureType_LIQUID_STAKE
	require.Error(t, features.ValdidateBasic())
	features = ValidHostChainInMsg(0).Features
	features.LiquidStakeIBC.Enabled = true
	require.Error(t, features.ValdidateBasic())

	features = ValidHostChainInMsg(0).Features
	features.LiquidStake.FeatureType = FeatureType_LIQUID_STAKE_IBC
	require.Error(t, features.ValdidateBasic())
	features = ValidHostChainInMsg(0).Features
	features.LiquidStake.Enabled = true
	require.Error(t, features.ValdidateBasic())

	// liquidstakeFeature
	lsfeature := ValidHostChainInMsg(0).Features.LiquidStake
	require.NoError(t, lsfeature.ValdidateBasic())
	lsfeature.Instantiation = InstantiationState_INSTANTIATION_INITIATED
	require.Error(t, lsfeature.ValdidateBasic())

	lsfeature = ValidHostChainInMsg(0).Features.LiquidStake
	lsfeature.Enabled = true
	require.Error(t, lsfeature.ValdidateBasic())
	lsfeature.ContractAddress = authtypes.NewModuleAddress("contract").String()
	require.Error(t, lsfeature.ValdidateBasic())

	lsfeature = ValidHostChainInMsg(0).Features.LiquidStake
	lsfeature.CodeID = 1 // non zero
	lsfeature.Instantiation = InstantiationState_INSTANTIATION_INITIATED
	require.NoError(t, lsfeature.ValdidateBasic())
	lsfeature.Enabled = true
	require.Error(t, lsfeature.ValdidateBasic())
	lsfeature.ContractAddress = authtypes.NewModuleAddress("contract").String()
	require.Error(t, lsfeature.ValdidateBasic())

	lsfeature = ValidHostChainInMsg(0).Features.LiquidStake
	lsfeature.CodeID = 1 // non zero
	lsfeature.Instantiation = InstantiationState_INSTANTIATION_COMPLETED
	require.Error(t, lsfeature.ValdidateBasic())
	lsfeature.ContractAddress = authtypes.NewModuleAddress("contract").String()
	require.NoError(t, lsfeature.ValdidateBasic())
	lsfeature.ContractAddress = "cosmos1xxxxxx"
	require.Error(t, lsfeature.ValdidateBasic())

	lsfeature = ValidHostChainInMsg(0).Features.LiquidStake
	lsfeature.Denoms = []string{"*", "stk/uxprt"}
	require.Error(t, lsfeature.ValdidateBasic())

	lsfeature = ValidHostChainInMsg(0).Features.LiquidStake
	require.Equal(t, false, lsfeature.AllowsAllDenoms())
	require.Equal(t, false, lsfeature.AllowsDenom("stk/uxprt"))
	lsfeature.Denoms = []string{"*"}
	require.NoError(t, lsfeature.ValdidateBasic())
	require.Equal(t, true, lsfeature.AllowsAllDenoms())
	require.Equal(t, true, lsfeature.AllowsDenom("stk/uxprt"))
	lsfeature.Denoms = []string{"*", "stk/uxprt"}
	require.Equal(t, false, lsfeature.AllowsAllDenoms())
	require.Equal(t, true, lsfeature.AllowsDenom("stk/uxprt"))

	lsfeature = ValidHostChainInMsg(0).Features.LiquidStake
	lsfeature2 := ValidHostChainInMsg(0).Features.LiquidStake
	require.Equal(t, true, lsfeature.Equals(lsfeature2))

	lsfeature2 = ValidHostChainInMsg(0).Features.LiquidStake
	lsfeature2.Enabled = true
	require.Equal(t, false, lsfeature.Equals(lsfeature2))
	lsfeature2.FeatureType = FeatureType_LIQUID_STAKE_IBC
	require.Equal(t, false, lsfeature.Equals(lsfeature2))
	lsfeature2.Denoms = []string{"*"}
	require.Equal(t, false, lsfeature.Equals(lsfeature2))
	lsfeature2.ContractAddress = "cosmos1xxx"
	require.Equal(t, false, lsfeature.Equals(lsfeature2))
	lsfeature2.Instantiation = InstantiationState_INSTANTIATION_COMPLETED
	require.Equal(t, false, lsfeature.Equals(lsfeature2))
	lsfeature2.CodeID = 1
	require.Equal(t, false, lsfeature.Equals(lsfeature2))

	require.Equal(t, "pstake_ratesync_1", DefaultPortOwner(1))
	require.Equal(t, "icacontroller-pstake_ratesync_1", MustICAPortIDFromOwner(DefaultPortOwner(1)))
	require.Panics(t, func() {
		MustICAPortIDFromOwner("")
	})

	owner, err := OwnerFromPortID("icacontroller-pstake_ratesync_1")
	require.NoError(t, err)
	require.Equal(t, "pstake_ratesync_1", owner)

	owner, err = OwnerFromPortID("ica-pstake_ratesync_1")
	require.Error(t, err)
	require.Equal(t, "", owner)

	id, err := IDFromPortID("icacontroller-pstake_ratesync_1")
	require.NoError(t, err)
	require.Equal(t, uint64(1), id)

	id, err = IDFromPortID("icacontroller-pstake_ratesync1")
	require.Error(t, err)
	require.Equal(t, uint64(0), id)
}
