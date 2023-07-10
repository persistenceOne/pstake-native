package types_test

import (
	"reflect"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func TestHostChain_GetHostChainTotalDelegations(t *testing.T) {

	tests := []struct {
		name      string
		hostChain func() types.HostChain
		want      math.Int
	}{
		{
			name: "count",
			hostChain: func() types.HostChain {
				hc := validHostChain()
				hc.Validators = append(hc.Validators, makeVal("one"))
				return *hc
			},
			want: sdk.NewInt(3),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hc := tt.hostChain()
			if got := hc.GetHostChainTotalDelegations(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetHostChainTotalDelegations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHostChain_GetValidator(t *testing.T) {
	foundvaloper := authtypes.NewModuleAddressOrBech32Address("testval1").String()

	type args struct {
		operatorAddress string
	}
	tests := []struct {
		name      string
		hostChain func() types.HostChain
		args      args
		want      *types.Validator
		want1     bool
	}{
		{
			name: "findone",
			hostChain: func() types.HostChain {
				hc := validHostChain()
				hc.Validators = append(hc.Validators, makeVal(foundvaloper))
				return *hc
			},
			args:  args{operatorAddress: foundvaloper},
			want:  makeVal(foundvaloper),
			want1: true,
		},
		{
			name: "not found",
			hostChain: func() types.HostChain {
				hc := validHostChain()
				hc.Validators = append(hc.Validators, makeVal(authtypes.NewModuleAddressOrBech32Address("testval4").String()))
				return *hc
			},
			args:  args{operatorAddress: foundvaloper},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hc := tt.hostChain()
			got, got1 := hc.GetValidator(tt.args.operatorAddress)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValidator() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetValidator() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func validHostChain() *types.HostChain {
	return &types.HostChain{
		ChainId:           "chain-1",
		ConnectionId:      "connection-1",
		Params:            nil,
		HostDenom:         "uatom",
		ChannelId:         "channel-1",
		PortId:            "transfer",
		DelegationAccount: nil,
		RewardsAccount:    nil,
		Validators: []*types.Validator{{
			OperatorAddress: authtypes.NewModuleAddressOrBech32Address("testval2").String(),
			Status:          stakingtypes.BondStatusBonded,
			Weight:          sdk.MustNewDecFromStr("0.5"),
			DelegatedAmount: sdk.OneInt(),
			ExchangeRate:    sdk.OneDec(),
			UnbondingEpoch:  0,
		}, {
			OperatorAddress: authtypes.NewModuleAddressOrBech32Address("testval3").String(),
			Status:          stakingtypes.BondStatusBonded,
			Weight:          sdk.MustNewDecFromStr("0.3"),
			DelegatedAmount: sdk.OneInt(),
			ExchangeRate:    sdk.OneDec(),
			UnbondingEpoch:  0,
		}},
		MinimumDeposit:     sdk.OneInt(),
		CValue:             sdk.OneDec(),
		LastCValue:         sdk.OneDec(),
		UnbondingFactor:    0,
		Active:             false,
		AutoCompoundFactor: sdk.Dec{},
	}

}
func makeVal(valoperAddr string) *types.Validator {
	return &types.Validator{
		OperatorAddress: valoperAddr,
		Status:          stakingtypes.BondStatusBonded,
		Weight:          sdk.MustNewDecFromStr("0.2"),
		DelegatedAmount: sdk.OneInt(),
		ExchangeRate:    sdk.OneDec(),
		UnbondingEpoch:  0,
	}
}
