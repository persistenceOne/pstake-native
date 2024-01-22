package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func TestParams_Validate(t *testing.T) {
	type fields struct {
		AdminAddress sdk.AccAddress
		FeeAddress   sdk.AccAddress
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Valid Params",
			fields: fields{
				AdminAddress: types.DefaultAdminAddress,
				FeeAddress:   types.DefaultFeeAddress,
			},
			wantErr: false,
		},
		{
			name: "Invalid admin address",
			fields: fields{
				AdminAddress: sdk.AccAddress{},
				FeeAddress:   types.DefaultFeeAddress,
			},
			wantErr: true,
		},
		{
			name: "invalid fee address",
			fields: fields{
				AdminAddress: types.DefaultAdminAddress,
				FeeAddress:   sdk.AccAddress{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &types.Params{
				AdminAddress: tt.fields.AdminAddress.String(),
				FeeAddress:   tt.fields.FeeAddress.String(),
			}
			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
