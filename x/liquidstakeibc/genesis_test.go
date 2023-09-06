package liquidstakeibc_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/pstake-native/v2/app/helpers"
	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc"
	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func TestGenesis(t *testing.T) {
	genesisState := &types.GenesisState{
		Params: types.DefaultParams(),
		HostChains: []*types.HostChain{{
			ChainId:      "chainA-1",
			ConnectionId: "connection-1",
			Params: &types.HostChainLSParams{
				DepositFee:      sdk.ZeroDec(),
				RestakeFee:      sdk.ZeroDec(),
				UnstakeFee:      sdk.ZeroDec(),
				RedemptionFee:   sdk.ZeroDec(),
				LsmValidatorCap: sdk.NewDec(1),
				LsmBondFactor:   sdk.NewDec(-1),
			},
			HostDenom: "uatom",
			ChannelId: "channel-1",
			PortId:    "transfer",
			DelegationAccount: &types.ICAAccount{
				Address:      "",
				Balance:      sdk.Coin{},
				Owner:        "",
				ChannelState: 0,
			},
			RewardsAccount: &types.ICAAccount{
				Address:      "",
				Balance:      sdk.Coin{},
				Owner:        "",
				ChannelState: 0,
			},
			Validators: []*types.Validator{{
				OperatorAddress: authtypes.NewModuleAddressOrBech32Address("testval").String(),
				Status:          stakingtypes.BondStatusBonded,
				Weight:          sdk.OneDec(),
				DelegatedAmount: sdk.NewInt(1221),
				ExchangeRate:    sdk.OneDec(),
				UnbondingEpoch:  0,
			}},
			MinimumDeposit:     sdk.OneInt(),
			CValue:             sdk.OneDec(),
			LastCValue:         sdk.OneDec(),
			UnbondingFactor:    0,
			Active:             false,
			AutoCompoundFactor: sdk.MustNewDecFromStr("2"),
		}},
		Deposits: []*types.Deposit{{
			ChainId:       "chainA-1",
			Amount:        sdk.NewInt64Coin("ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9", 100),
			Epoch:         0,
			State:         0,
			IbcSequenceId: "",
		}},
		Unbondings: []*types.Unbonding{{
			ChainId:       "chainA-1",
			EpochNumber:   0,
			MatureTime:    time.Time{},
			BurnAmount:    sdk.NewInt64Coin("stk/uatom", 10),
			UnbondAmount:  sdk.NewInt64Coin("uatom", 10),
			IbcSequenceId: "",
			State:         0,
		}},
		UserUnbondings: []*types.UserUnbonding{{
			ChainId:      "chainA-1",
			EpochNumber:  0,
			Address:      authtypes.NewModuleAddressOrBech32Address("test").String(),
			StkAmount:    sdk.NewInt64Coin("stk/uatom", 10),
			UnbondAmount: sdk.NewInt64Coin("uatom", 10),
		}},
		ValidatorUnbondings: []*types.ValidatorUnbonding{{
			ChainId:          "chainA-1",
			EpochNumber:      0,
			MatureTime:       time.Time{},
			ValidatorAddress: authtypes.NewModuleAddressOrBech32Address("testval").String(),
			Amount:           sdk.NewInt64Coin("uatom", 1000),
			IbcSequenceId:    "",
		}},
	}

	_, pStakeApp, ctx := helpers.CreateTestApp(t)
	k := pStakeApp.LiquidStakeIBCKeeper
	liquidstakeibc.InitGenesis(ctx, k, genesisState)

	got := liquidstakeibc.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	require.Equal(t, genesisState.Params, got.Params)
	require.Equal(t, genesisState.HostChains, got.HostChains)
	require.Equal(t, genesisState.Deposits, got.Deposits)
}
