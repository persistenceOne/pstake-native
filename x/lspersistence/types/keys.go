package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName is the name of the liquidstaking module
	ModuleName = "lspersistence"

	// RouterKey is the message router key for the liquidstaking module
	RouterKey = ModuleName

	// StoreKey is the default store key for the liquidstaking module
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the liquidstaking module
	QuerierRoute = ModuleName
)

var (
	// Keys for store prefixes
	ParamsKey           = []byte{0x00}
	LiquidValidatorsKey = []byte{0xc0} // prefix for each key to a liquid validator
)

// GetLiquidValidatorKey creates the key for the liquid validator with address
// VALUE: lspersistence/LiquidValidator
func GetLiquidValidatorKey(operatorAddr sdk.ValAddress) []byte {
	return append(LiquidValidatorsKey, address.MustLengthPrefix(operatorAddr)...)
}
