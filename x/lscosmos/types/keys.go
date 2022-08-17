package types

import icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"

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
	DepositModuleAccount = ModuleName + "_pstake_deposit_account"

	// DelegationModuleAccount DelegationModuleAccountName
	DelegationModuleAccount = ModuleName + "_pstake_delegation_account"

	// RewardModuleAccount RewardModuleAccountName
	RewardModuleAccount = ModuleName + "_pstake_reward_account"

	// UndelegationModuleAccount UndelegationModuleAccountName,
	// This account will not be a part of maccPerms - Deny list, since it receives undelegated tokens.
	UndelegationModuleAccount = ModuleName + "_pstake_undelegation_account"

	DelegationEpochIdentifier   = "day"
	RewardEpochIdentifier       = "day"
	UndelegationEpochIdentifier = "week"
)

var (
	DelegationAccountPortID, _   = icatypes.NewControllerPortID(DelegationModuleAccount)
	RewardAccountPortID, _       = icatypes.NewControllerPortID(RewardModuleAccount)
	UndelegationAccountPortID, _ = icatypes.NewControllerPortID(UndelegationModuleAccount)
)
var (
	// PortKey defines the key to store the port ID in store
	PortKey = KeyPrefix("lscosmos-port")

	ModuleEnableKey          = []byte{0x01}
	CosmosIBCParamsKey       = []byte{0x02}
	AllowListedValidatorsKey = []byte{0x03}
	DelegationAccountKey     = []byte{0x04}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
