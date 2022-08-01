package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName is the name of the lspersistence module
	ModuleName = "lspersistence"

	// RouterKey is the message router key for the lspersistence module
	RouterKey = ModuleName

	// StoreKey is the default store key for the lspersistence module
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the lspersistence module
	QuerierRoute = ModuleName
)

var (
	// Keys for store prefixes
	LiquidValidatorsKey = []byte{0xc0} // prefix for each key to a liquid validator
)

// GetLiquidValidatorKey creates the key for the liquid validator with address
// VALUE: lspersistence/LiquidValidator
func GetLiquidValidatorKey(operatorAddr sdk.ValAddress) []byte {
	return append(LiquidValidatorsKey, address.MustLengthPrefix(operatorAddr)...)
}
