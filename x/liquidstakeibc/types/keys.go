package types

const (
	// ModuleName defines the module name
	ModuleName = "liquidstakeibc"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// DepositModuleAccount DepositModuleAccountName
	DepositModuleAccount = ModuleName + "_deposit_account"
)

var (
	HostChainKey = []byte{0x01}
)
