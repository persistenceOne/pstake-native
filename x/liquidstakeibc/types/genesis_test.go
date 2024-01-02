package types_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState func() *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: func() *types.GenesisState { return types.DefaultGenesisState() },
			valid:    true,
		},
		{
			desc:     "invalid genesis state, params not set",
			genState: func() *types.GenesisState { return &types.GenesisState{} },
			valid:    false,
		},
		{
			desc:     "Valid State with all fields present",
			genState: func() *types.GenesisState { return ValidGenesis() },
			valid:    true,
		},
		{
			desc: "Multiple host chains with same chain-id",
			genState: func() *types.GenesisState {
				genesis := ValidGenesis()
				genesis.HostChains = append(genesis.HostChains, &types.HostChain{ChainId: "chainA-1"})
				return genesis
			},
			valid: false,
		},
		{
			desc: "host chains invalid",
			genState: func() *types.GenesisState {
				genesis := ValidGenesis()
				genesis.HostChains[0].CValue = sdk.MustNewDecFromStr("-1")
				return genesis
			},
			valid: false,
		},
		{
			desc: "deposits of non existent chain-id",
			genState: func() *types.GenesisState {
				genesis := ValidGenesis()
				genesis.Deposits = append(genesis.Deposits, &types.Deposit{
					ChainId:       "nonExistent-1",
					Amount:        sdk.NewInt64Coin("ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9", 100),
					Epoch:         0,
					State:         0,
					IbcSequenceId: ""})
				return genesis
			},
			valid: false,
		},
		{
			desc: "invalid deposit",
			genState: func() *types.GenesisState {
				genesis := ValidGenesis()
				genesis.Deposits = append(genesis.Deposits, &types.Deposit{
					ChainId:       "chainA-1",
					Amount:        sdk.NewInt64Coin("ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9", 100),
					Epoch:         0,
					State:         4,
					IbcSequenceId: ""})
				return genesis
			},
			valid: false,
		}, {
			desc: "invalid amount denom",
			genState: func() *types.GenesisState {
				genesis := ValidGenesis()
				genesis.Deposits = append(genesis.Deposits, &types.Deposit{
					ChainId:       "chainA-1",
					Amount:        sdk.NewInt64Coin("uatom", 100),
					Epoch:         0,
					State:         0,
					IbcSequenceId: ""})
				return genesis
			},
			valid: false,
		},
		{
			desc: "unbondings of non existent chain-id",
			genState: func() *types.GenesisState {
				genesis := ValidGenesis()
				genesis.Unbondings = append(genesis.Unbondings, &types.Unbonding{ChainId: "nonExistent-1"})
				return genesis
			},
			valid: false,
		},
		{
			desc: "invalid BurnAmount",
			genState: func() *types.GenesisState {
				genesis := ValidGenesis()
				genesis.Unbondings = append(genesis.Unbondings, &types.Unbonding{
					ChainId:       "chainA-1",
					EpochNumber:   0,
					MatureTime:    time.Time{},
					BurnAmount:    sdk.NewInt64Coin("ibc/uatom", 10),
					UnbondAmount:  sdk.NewInt64Coin("uatom", 10),
					IbcSequenceId: "",
					State:         0,
				})
				return genesis
			},
			valid: false,
		}, {
			desc: "invalid unbound amount",
			genState: func() *types.GenesisState {
				genesis := ValidGenesis()
				genesis.Unbondings = append(genesis.Unbondings, &types.Unbonding{
					ChainId:       "chainA-1",
					EpochNumber:   0,
					MatureTime:    time.Time{},
					BurnAmount:    sdk.NewInt64Coin("stk/uatom", 10),
					UnbondAmount:  sdk.NewInt64Coin("ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9", 10),
					IbcSequenceId: "",
					State:         0,
				})
				return genesis
			},
			valid: false,
		}, {
			desc: "invalid unbonding state",
			genState: func() *types.GenesisState {
				genesis := ValidGenesis()
				genesis.Unbondings = append(genesis.Unbondings, &types.Unbonding{
					ChainId:       "chainA-1",
					EpochNumber:   0,
					MatureTime:    time.Time{},
					BurnAmount:    sdk.NewInt64Coin("stk/uatom", 0),
					UnbondAmount:  sdk.NewInt64Coin("uatom", 0),
					IbcSequenceId: "",
					State:         6,
				})
				return genesis
			},
			valid: false,
		},
		{
			desc: "user unbondings of non existent chain-id",
			genState: func() *types.GenesisState {
				genesis := ValidGenesis()
				genesis.UserUnbondings = append(genesis.UserUnbondings, &types.UserUnbonding{ChainId: "nonExistent-1"})
				return genesis
			},
			valid: false,
		},
		{
			desc: "user unbondings incorrect stkamount denom",
			genState: func() *types.GenesisState {
				genesis := ValidGenesis()
				genesis.UserUnbondings = append(genesis.UserUnbondings,
					&types.UserUnbonding{
						ChainId:      "chainA-1",
						EpochNumber:  0,
						Address:      authtypes.NewModuleAddressOrBech32Address("test2").String(),
						StkAmount:    sdk.NewInt64Coin("uatom", 10),
						UnbondAmount: sdk.NewInt64Coin("uatom", 10),
					})
				return genesis
			},
			valid: false,
		}, {
			desc: "user unbondings incorrect unboundAmount denom",
			genState: func() *types.GenesisState {
				genesis := ValidGenesis()
				genesis.UserUnbondings = append(genesis.UserUnbondings,
					&types.UserUnbonding{
						ChainId:      "chainA-1",
						EpochNumber:  0,
						Address:      authtypes.NewModuleAddressOrBech32Address("test2").String(),
						StkAmount:    sdk.NewInt64Coin("stk/uatom", 10),
						UnbondAmount: sdk.NewInt64Coin("ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9", 10),
					})
				return genesis
			},
			valid: false,
		}, {
			desc: "user unbondings invalid",
			genState: func() *types.GenesisState {
				genesis := ValidGenesis()
				genesis.UserUnbondings = append(genesis.UserUnbondings,
					&types.UserUnbonding{
						ChainId:      "chainA-1",
						EpochNumber:  0,
						Address:      "",
						StkAmount:    sdk.NewInt64Coin("stk/uatom", 10),
						UnbondAmount: sdk.NewInt64Coin("uatom", 10),
					})
				return genesis
			},
			valid: false,
		},
		{
			desc: "validator unbonding of non existent chain-id",
			genState: func() *types.GenesisState {
				genesis := ValidGenesis()
				genesis.ValidatorUnbondings = append(genesis.ValidatorUnbondings, &types.ValidatorUnbonding{ChainId: "nonExistent-1"})
				return genesis
			},
			valid: false,
		},
		{
			desc: "validator unbonding incorrect amount denom",
			genState: func() *types.GenesisState {
				genesis := ValidGenesis()
				genesis.ValidatorUnbondings = append(genesis.ValidatorUnbondings,
					&types.ValidatorUnbonding{
						ChainId:          "chainA-1",
						EpochNumber:      0,
						MatureTime:       time.Time{},
						ValidatorAddress: authtypes.NewModuleAddressOrBech32Address("testval2").String(),
						Amount:           sdk.NewInt64Coin("ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9", 1000),
						IbcSequenceId:    "",
					})
				return genesis
			},
			valid: false,
		},
		{
			desc: "validator unbonding invalid",
			genState: func() *types.GenesisState {
				genesis := ValidGenesis()
				genesis.ValidatorUnbondings = append(genesis.ValidatorUnbondings,
					&types.ValidatorUnbonding{
						ChainId:          "chainA-1",
						EpochNumber:      0,
						MatureTime:       time.Time{},
						ValidatorAddress: "",
						Amount:           sdk.NewInt64Coin("uatom", 1000),
						IbcSequenceId:    "",
					})
				return genesis
			},
			valid: false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState().Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func ValidGenesis() *types.GenesisState {
	return &types.GenesisState{
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
				ExchangeRate:    sdk.Dec{},
				UnbondingEpoch:  0,
			}},
			MinimumDeposit:     sdk.OneInt(),
			CValue:             sdk.OneDec(),
			LastCValue:         sdk.Dec{},
			UnbondingFactor:    0,
			Active:             false,
			AutoCompoundFactor: sdk.Dec{},
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
}
