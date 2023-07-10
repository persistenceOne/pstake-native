package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func TestParams_Validate(t *testing.T) {
	type fields struct {
		AdminAddress     string
		FeeAddress       string
		UpperCValueLimit sdk.Dec
		LowerCValueLimit sdk.Dec
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Valid Params",
			fields: fields{
				AdminAddress:     types.DefaultAdminAddress,
				FeeAddress:       types.DefaultFeeAddress,
				UpperCValueLimit: sdk.OneDec(),
				LowerCValueLimit: sdk.ZeroDec(),
			},
			wantErr: false,
		},
		{
			name: "Invalid admin address",
			fields: fields{
				AdminAddress:     "",
				FeeAddress:       types.DefaultFeeAddress,
				UpperCValueLimit: sdk.OneDec(),
				LowerCValueLimit: sdk.ZeroDec(),
			},
			wantErr: true,
		},
		{
			name: "invalid fee address",
			fields: fields{
				AdminAddress:     types.DefaultAdminAddress,
				FeeAddress:       "",
				UpperCValueLimit: sdk.OneDec(),
				LowerCValueLimit: sdk.ZeroDec(),
			},
			wantErr: true,
		},
		{
			name: "Invalid Lower Limit",
			fields: fields{
				AdminAddress:     types.DefaultAdminAddress,
				FeeAddress:       types.DefaultFeeAddress,
				UpperCValueLimit: sdk.OneDec(),
				LowerCValueLimit: sdk.OneDec(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &types.Params{
				AdminAddress:     tt.fields.AdminAddress,
				FeeAddress:       tt.fields.FeeAddress,
				UpperCValueLimit: tt.fields.UpperCValueLimit,
				LowerCValueLimit: tt.fields.LowerCValueLimit,
			}
			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
