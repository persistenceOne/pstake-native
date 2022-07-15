package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/x/cosmos/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSignature(t *testing.T) {
	type testState struct {
		accAddressPrefix string
		valAddressPrefix string
		accAddress       string
		valAddress       string
	}
	type expectedState struct {
		accAddress string
		valAddress string
		err        error
		err1       error
		err2       error
		err3       error
	}
	testMatrix := []struct {
		given    testState
		expected expectedState
	}{
		{
			given: testState{
				accAddressPrefix: "cosmos",
				accAddress:       "cosmos1hcqg5wj9t42zawqkqucs7la85ffyv08lum327c",
				valAddressPrefix: "cosmosvaloper",
				valAddress:       "cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2",
			},
			expected: expectedState{
				accAddress: "cosmos1hcqg5wj9t42zawqkqucs7la85ffyv08lum327c",
				valAddress: "cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2",
				err:        nil,
				err1:       nil,
				err2:       nil,
				err3:       nil,
			},
		},
		{
			given: testState{
				accAddressPrefix: "cosmos",
				accAddress:       "cosmos1hcqg5wj9t42zawqkqucs7la85ffyvum327c",
				valAddressPrefix: "cosmosvaloper",
				valAddress:       "cosmosvaloper1lcck2cxh7dzg53kysg9ktdrsjj6jfwlnm2",
			},
			expected: expectedState{},
		},
		{
			given: testState{
				accAddressPrefix: "persistence",
				accAddress:       "persistence1kszhtg25k2p55mr8l7zxcs2v4vv44xq49qh9p4",
				valAddressPrefix: "persistencevaloper",
				valAddress:       "persistencevaloper1pss7nxeh3f9md2vuxku8q99femnwdjtcga9tdu",
			},
			expected: expectedState{
				accAddress: "persistence1kszhtg25k2p55mr8l7zxcs2v4vv44xq49qh9p4",
				valAddress: "persistencevaloper1pss7nxeh3f9md2vuxku8q99femnwdjtcga9tdu",
				err:        nil,
				err1:       nil,
				err2:       nil,
				err3:       nil,
			},
		},
		{
			given: testState{
				accAddressPrefix: "persistence",
				accAddress:       "persistence1kszhtg25k2p55mr8l2v4vv44xq49qh9p4",
				valAddressPrefix: "persistencevaloper",
				valAddress:       "persistencevaloper1pss7nxeh3f9xku8q99femnwdjtcga9tdu",
			},
			expected: expectedState{},
		},
		{
			given: testState{
				accAddressPrefix: "cosmos",
				accAddress:       "persistence1kszhtg25k2p55mr8l2v4vv44xq49qh9p4",
				valAddressPrefix: "cosmosvaloper",
				valAddress:       "persistencevaloper1pss7nxeh3f9md2vuxku8q99femnwdjtcga9tdu",
			},
			expected: expectedState{},
		},
	}

	for _, test := range testMatrix {
		accAddress, err := types.AccAddressFromBech32(test.given.accAddress, test.given.accAddressPrefix)
		if accAddress != nil || err == nil {
			require.Equal(t, test.expected.err, err)

			err1 := sdk.VerifyAddressFormat(accAddress)
			require.Equal(t, test.expected.err1, err1)

			accAddress1, err4 := types.Bech32ifyAddressBytes(test.given.accAddressPrefix, accAddress)
			if accAddress1 != "" || err4 == nil {
				require.Equal(t, test.expected.accAddress, accAddress1)
			} else {
				require.Error(t, err4)
			}
		} else {
			require.Error(t, err)
		}

		valAddress, err2 := types.ValAddressFromBech32(test.given.valAddress, test.given.valAddressPrefix)
		if valAddress != nil || err2 == nil {
			require.Equal(t, test.expected.err2, err2)

			err3 := sdk.VerifyAddressFormat(accAddress)
			require.Equal(t, test.expected.err3, err3)

			valAddress1, err4 := types.Bech32ifyValAddressBytes(test.given.valAddressPrefix, valAddress)
			if valAddress1 != "" || err4 == nil {
				require.Equal(t, test.expected.valAddress, valAddress1)
			} else {
				require.Error(t, err4)
			}
		} else {
			require.Error(t, err2)
		}
	}
}
