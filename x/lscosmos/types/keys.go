package types

const (
	// ModuleName defines the module name
	ModuleName = "lscosmos"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_lscosmos"

	// Version defines the current version the IBC module supports
	Version = "lscosmos-1"

	// PortID is the default port id that module binds to
	PortID = "lscosmos"

	// MsgTypeLiquidStake is the type of message to liquid stake
	MsgTypeLiquidStake = "msg_liquid_stake"

	// DepositModuleAccount DepositModuleAccountName
	DepositModuleAccount = "deposit_account"
)

var (
	// PortKey defines the key to store the port ID in store
	PortKey = KeyPrefix("lscosmos-port")

	CosmosIBCParamsKey = []byte{0x01}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
