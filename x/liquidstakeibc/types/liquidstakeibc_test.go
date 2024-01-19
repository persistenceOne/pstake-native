package types_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/persistenceOne/pstake-native/v2/app"
	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func init() {
	app.SetAddressPrefixes()
	types.RegisterInterfaces(codectypes.NewInterfaceRegistry())
}

func TestCurrentUnbondingEpoch(t *testing.T) {
	type args struct {
		factor      int64
		epochNumber int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "1 gets 4",
			args: args{
				factor:      4,
				epochNumber: 1,
			},
			want: 4,
		}, {
			name: "4 gets 4",
			args: args{
				factor:      4,
				epochNumber: 4,
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := types.CurrentUnbondingEpoch(tt.args.factor, tt.args.epochNumber); got != tt.want {
				t.Errorf("CurrentUnbondingEpoch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsUnbondingEpoch(t *testing.T) {
	type args struct {
		factor      int64
		epochNumber int64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid",
			args: args{
				factor:      4,
				epochNumber: 4,
			},
			want: true,
		}, {
			name: "valid",
			args: args{
				factor:      4,
				epochNumber: 3,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := types.IsUnbondingEpoch(tt.args.factor, tt.args.epochNumber); got != tt.want {
				t.Errorf("IsUnbondingEpoch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultDelegateAccountPortOwner(t *testing.T) {
	type args struct {
		chainID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid",
			args: args{chainID: "chain-1"},
			want: "chain-1.delegate",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := types.DefaultDelegateAccountPortOwner(tt.args.chainID); got != tt.want {
				t.Errorf("DefaultDelegateAccountPortOwner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultRewardsAccountPortOwner(t *testing.T) {
	type args struct {
		chainID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid",
			args: args{chainID: "chain-1"},
			want: "chain-1.rewards",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := types.DefaultRewardsAccountPortOwner(tt.args.chainID); got != tt.want {
				t.Errorf("DefaultRewardsAccountPortOwner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeposit_Validate(t *testing.T) {
	type fields struct {
		ChainId       string
		Amount        sdk.Coin
		Epoch         int64
		State         types.Deposit_DepositState
		IbcSequenceId string
	}
	validCoin := sdk.NewInt64Coin("ibc/uatom", 1000)
	invalidCoin := validCoin
	invalidCoin.Amount = sdk.NewInt(-1)
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				ChainId:       "chain-1",
				Amount:        validCoin,
				Epoch:         0,
				State:         0,
				IbcSequenceId: "",
			},
			wantErr: false,
		},
		{
			name: "invalid amount",
			fields: fields{
				ChainId:       "chain-1",
				Amount:        invalidCoin,
				Epoch:         0,
				State:         0,
				IbcSequenceId: "",
			},
			wantErr: true,
		},
		{
			name: "invalid state",
			fields: fields{
				ChainId:       "chain-1",
				Amount:        validCoin,
				Epoch:         0,
				State:         5,
				IbcSequenceId: "",
			},
			wantErr: true,
		},
		{
			name: "invalid state",
			fields: fields{
				ChainId:       "chain-1",
				Amount:        validCoin,
				Epoch:         0,
				State:         1,
				IbcSequenceId: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deposit := &types.Deposit{
				ChainId:       tt.fields.ChainId,
				Amount:        tt.fields.Amount,
				Epoch:         tt.fields.Epoch,
				State:         tt.fields.State,
				IbcSequenceId: tt.fields.IbcSequenceId,
			}
			if err := deposit.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHostChain_Validate(t *testing.T) {
	type fields struct {
		ChainId            string
		ConnectionId       string
		Params             *types.HostChainLSParams
		HostDenom          string
		ChannelId          string
		PortId             string
		DelegationAccount  *types.ICAAccount
		RewardsAccount     *types.ICAAccount
		Validators         []*types.Validator
		MinimumDeposit     math.Int
		CValue             sdk.Dec
		LastCValue         sdk.Dec
		UnbondingFactor    int64
		Active             bool
		AutoCompoundFactor sdk.Dec
	}
	validFields := func() fields {
		return fields{
			ChainId:      "chain-1",
			ConnectionId: "connection-1",
			Params: &types.HostChainLSParams{
				DepositFee:    sdk.ZeroDec(),
				RestakeFee:    sdk.ZeroDec(),
				UnstakeFee:    sdk.ZeroDec(),
				RedemptionFee: sdk.ZeroDec(),
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
				DelegatedAmount: sdk.OneInt(),
				ExchangeRate:    sdk.OneDec(),
				UnbondingEpoch:  0,
			}},
			MinimumDeposit:     sdk.OneInt(),
			CValue:             sdk.OneDec(),
			LastCValue:         sdk.OneDec(),
			UnbondingFactor:    4,
			Active:             false,
			AutoCompoundFactor: sdk.MustNewDecFromStr("2"),
		}
	}

	tests := []struct {
		name    string
		fields  func() fields
		wantErr bool
	}{
		{
			name:    "correct",
			fields:  func() fields { return validFields() },
			wantErr: false,
		},
		{
			name: "invalid params",
			fields: func() fields {
				newfields := validFields()
				newfields.Params.DepositFee = sdk.MustNewDecFromStr("2")
				return newfields
			},
			wantErr: true,
		},
		{
			name: "invalid validator",
			fields: func() fields {
				newfields := validFields()
				newfields.Validators[0].Weight = sdk.MustNewDecFromStr("2")
				return newfields
			},
			wantErr: true,
		},
		{
			name: "invalid cvalue",
			fields: func() fields {
				newfields := validFields()
				newfields.CValue = sdk.MustNewDecFromStr("-2")
				return newfields
			},
			wantErr: true,
		},
		{
			name: "invalid mindeposit",
			fields: func() fields {
				newfields := validFields()
				newfields.MinimumDeposit = sdk.NewInt(-2)
				return newfields
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ttfields := tt.fields()
			hc := &types.HostChain{
				ChainId:            ttfields.ChainId,
				ConnectionId:       ttfields.ConnectionId,
				Params:             ttfields.Params,
				HostDenom:          ttfields.HostDenom,
				ChannelId:          ttfields.ChannelId,
				PortId:             ttfields.PortId,
				DelegationAccount:  ttfields.DelegationAccount,
				RewardsAccount:     ttfields.RewardsAccount,
				Validators:         ttfields.Validators,
				MinimumDeposit:     ttfields.MinimumDeposit,
				CValue:             ttfields.CValue,
				LastCValue:         ttfields.LastCValue,
				UnbondingFactor:    ttfields.UnbondingFactor,
				Active:             ttfields.Active,
				AutoCompoundFactor: ttfields.AutoCompoundFactor,
			}
			if err := hc.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnbonding_Validate(t *testing.T) {
	type fields struct {
		ChainId       string
		EpochNumber   int64
		MatureTime    time.Time
		BurnAmount    sdk.Coin
		UnbondAmount  sdk.Coin
		IbcSequenceId string
		State         types.Unbonding_UnbondingState
	}
	validCoin := sdk.NewInt64Coin("ibc/uatom", 1000)
	invalidCoin := validCoin
	invalidCoin.Amount = sdk.NewInt(-1000)
	validFields := func() fields {
		return fields{
			ChainId:       "chain-1",
			EpochNumber:   0,
			MatureTime:    time.Time{},
			BurnAmount:    validCoin,
			UnbondAmount:  validCoin,
			IbcSequenceId: "",
			State:         0,
		}
	}
	tests := []struct {
		name    string
		fields  func() fields
		wantErr bool
	}{
		{
			name:    "valid",
			fields:  func() fields { return validFields() },
			wantErr: false,
		}, {
			name: "invalid burnamount",
			fields: func() fields {
				field := validFields()
				field.BurnAmount = invalidCoin
				return field
			},
			wantErr: true,
		}, {
			name: "invalid unbound amount",
			fields: func() fields {
				field := validFields()
				field.UnbondAmount = invalidCoin
				return field
			},
			wantErr: true,
		}, {
			name: "invalid unbonding state",
			fields: func() fields {
				field := validFields()
				field.State = 9
				return field
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fields()
			u := &types.Unbonding{
				ChainId:       fields.ChainId,
				EpochNumber:   fields.EpochNumber,
				MatureTime:    fields.MatureTime,
				BurnAmount:    fields.BurnAmount,
				UnbondAmount:  fields.UnbondAmount,
				IbcSequenceId: fields.IbcSequenceId,
				State:         fields.State,
			}
			if err := u.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserUnbonding_Validate(t *testing.T) {
	type fields struct {
		ChainId      string
		EpochNumber  int64
		Address      string
		StkAmount    sdk.Coin
		UnbondAmount sdk.Coin
	}
	validCoin := sdk.NewInt64Coin("stk/uatom", 1000)
	invalidCoin := validCoin
	invalidCoin.Amount = sdk.NewInt(-1000)
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				ChainId:      "chain-1",
				EpochNumber:  0,
				Address:      authtypes.NewModuleAddressOrBech32Address("test").String(),
				StkAmount:    validCoin,
				UnbondAmount: validCoin,
			},
			wantErr: false,
		},
		{
			name: "invalid coin",
			fields: fields{
				ChainId:      "chain-1",
				EpochNumber:  0,
				Address:      authtypes.NewModuleAddressOrBech32Address("test").String(),
				StkAmount:    validCoin,
				UnbondAmount: invalidCoin,
			},
			wantErr: true,
		},
		{
			name: "invalid addr",
			fields: fields{
				ChainId:      "chain-1",
				EpochNumber:  0,
				Address:      "test",
				StkAmount:    validCoin,
				UnbondAmount: validCoin,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ub := &types.UserUnbonding{
				ChainId:      tt.fields.ChainId,
				EpochNumber:  tt.fields.EpochNumber,
				Address:      tt.fields.Address,
				StkAmount:    tt.fields.StkAmount,
				UnbondAmount: tt.fields.UnbondAmount,
			}
			if err := ub.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatorUnbonding_Validate(t *testing.T) {
	type fields struct {
		ChainId          string
		EpochNumber      int64
		MatureTime       time.Time
		ValidatorAddress string
		Amount           sdk.Coin
		IbcSequenceId    string
	}
	validCoin := sdk.NewInt64Coin("uatom", 1000)
	invalidCoin := validCoin
	invalidCoin.Amount = sdk.NewInt(-1000)
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				ChainId:          "chain-1",
				EpochNumber:      0,
				MatureTime:       time.Time{},
				ValidatorAddress: authtypes.NewModuleAddressOrBech32Address("testval").String(),
				Amount:           validCoin,
				IbcSequenceId:    "",
			},
			wantErr: false,
		},
		{
			name: "invalid amount",
			fields: fields{
				ChainId:          "chain-1",
				EpochNumber:      0,
				MatureTime:       time.Time{},
				ValidatorAddress: authtypes.NewModuleAddressOrBech32Address("testval").String(),
				Amount:           invalidCoin,
				IbcSequenceId:    "",
			},
			wantErr: true,
		},
		{
			name: "invalid addr",
			fields: fields{
				ChainId:          "chain-1",
				EpochNumber:      0,
				MatureTime:       time.Time{},
				ValidatorAddress: "testval",
				Amount:           validCoin,
				IbcSequenceId:    "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vb := &types.ValidatorUnbonding{
				ChainId:          tt.fields.ChainId,
				EpochNumber:      tt.fields.EpochNumber,
				MatureTime:       tt.fields.MatureTime,
				ValidatorAddress: tt.fields.ValidatorAddress,
				Amount:           tt.fields.Amount,
				IbcSequenceId:    tt.fields.IbcSequenceId,
			}
			if err := vb.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHostChainLSParams_Validate(t *testing.T) {
	type fields struct {
		DepositFee    sdk.Dec
		RestakeFee    sdk.Dec
		UnstakeFee    sdk.Dec
		RedemptionFee sdk.Dec
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				DepositFee:    sdk.ZeroDec(),
				RestakeFee:    sdk.ZeroDec(),
				UnstakeFee:    sdk.ZeroDec(),
				RedemptionFee: sdk.ZeroDec(),
			},
			wantErr: false,
		},
		{
			name: "invalid deposit fee",
			fields: fields{
				DepositFee:    sdk.MustNewDecFromStr("-1"),
				RestakeFee:    sdk.ZeroDec(),
				UnstakeFee:    sdk.ZeroDec(),
				RedemptionFee: sdk.ZeroDec(),
			},
			wantErr: true,
		},
		{
			name: "invalid restake fee",
			fields: fields{
				DepositFee:    sdk.ZeroDec(),
				RestakeFee:    sdk.MustNewDecFromStr("1.1"),
				UnstakeFee:    sdk.ZeroDec(),
				RedemptionFee: sdk.ZeroDec(),
			},
			wantErr: true,
		},
		{
			name: "invalid unstake fee",
			fields: fields{
				DepositFee:    sdk.ZeroDec(),
				RestakeFee:    sdk.ZeroDec(),
				UnstakeFee:    sdk.MustNewDecFromStr("-1"),
				RedemptionFee: sdk.ZeroDec(),
			},
			wantErr: true,
		},
		{
			name: "invalid redemption fee",
			fields: fields{
				DepositFee:    sdk.ZeroDec(),
				RestakeFee:    sdk.ZeroDec(),
				UnstakeFee:    sdk.ZeroDec(),
				RedemptionFee: sdk.MustNewDecFromStr("1.2"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := &types.HostChainLSParams{
				DepositFee:    tt.fields.DepositFee,
				RestakeFee:    tt.fields.RestakeFee,
				UnstakeFee:    tt.fields.UnstakeFee,
				RedemptionFee: tt.fields.RedemptionFee,
			}
			if err := params.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidator_Validate(t *testing.T) {
	type fields struct {
		OperatorAddress string
		Status          string
		Weight          sdk.Dec
		DelegatedAmount math.Int
		ExchangeRate    sdk.Dec
		UnbondingEpoch  int64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				OperatorAddress: authtypes.NewModuleAddressOrBech32Address("testval").String(),
				Status:          stakingtypes.BondStatusBonded,
				Weight:          sdk.OneDec(),
				DelegatedAmount: sdk.OneInt(),
				ExchangeRate:    sdk.OneDec(),
				UnbondingEpoch:  0,
			},
			wantErr: false,
		},
		{
			name: "invalid operatorAddr",
			fields: fields{
				OperatorAddress: "testval",
				Status:          stakingtypes.BondStatusBonded,
				Weight:          sdk.OneDec(),
				DelegatedAmount: sdk.OneInt(),
				ExchangeRate:    sdk.OneDec(),
				UnbondingEpoch:  0,
			},
			wantErr: true,
		},
		{
			name: "invalid status",
			fields: fields{
				OperatorAddress: authtypes.NewModuleAddressOrBech32Address("testval").String(),
				Status:          "Status random",
				Weight:          sdk.OneDec(),
				DelegatedAmount: sdk.OneInt(),
				ExchangeRate:    sdk.OneDec(),
				UnbondingEpoch:  0,
			},
			wantErr: true,
		},
		{
			name: "invalid weight",
			fields: fields{
				OperatorAddress: authtypes.NewModuleAddressOrBech32Address("testval").String(),
				Status:          stakingtypes.BondStatusBonded,
				Weight:          sdk.MustNewDecFromStr("3"),
				DelegatedAmount: sdk.OneInt(),
				ExchangeRate:    sdk.OneDec(),
				UnbondingEpoch:  0,
			},
			wantErr: true,
		},
		{
			name: "invalid delegated amount",
			fields: fields{
				OperatorAddress: authtypes.NewModuleAddressOrBech32Address("testval").String(),
				Status:          stakingtypes.BondStatusBonded,
				Weight:          sdk.OneDec(),
				DelegatedAmount: sdk.NewInt(-1),
				ExchangeRate:    sdk.OneDec(),
				UnbondingEpoch:  0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &types.Validator{
				OperatorAddress: tt.fields.OperatorAddress,
				Status:          tt.fields.Status,
				Weight:          tt.fields.Weight,
				DelegatedAmount: tt.fields.DelegatedAmount,
				ExchangeRate:    tt.fields.ExchangeRate,
				UnbondingEpoch:  tt.fields.UnbondingEpoch,
			}
			if err := validator.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
