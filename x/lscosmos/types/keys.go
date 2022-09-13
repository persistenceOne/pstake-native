package types

import (
	"time"

	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
)

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

	// MsgTypeLiquidStake is the type of message to liquid stake
	MsgTypeLiquidStake = "msg_liquid_stake"

	// MsgTypeJuice is the type of message Juice
	MsgTypeJuice = "msg_juice"

	// MsgTypeLiquidUnstake is the type of message liquid unstake
	MsgTypeLiquidUnstake = "msg_liquid_unstake"

	// DepositModuleAccount DepositModuleAccountName
	DepositModuleAccount = ModuleName + "_pstake_deposit_account"

	// DelegationModuleAccount DelegationModuleAccountName
	DelegationModuleAccount = ModuleName + "_pstake_delegation_account"

	// RewardModuleAccount RewardModuleAccountName
	RewardModuleAccount = ModuleName + "_pstake_reward_account"

	// UndelegationModuleAccount UndelegationModuleAccountName,
	// This account will not be a part of maccPerms - Deny list, since it receives undelegated tokens.
	UndelegationModuleAccount = ModuleName + "_pstake_undelegation_account"

	// RewardBoosterModuleAccount RewardBoosterModuleAccountName
	RewardBoosterModuleAccount = ModuleName + "_reward_booster_account"

	DelegationEpochIdentifier   = "day"
	RewardEpochIdentifier       = "day"
	UndelegationEpochIdentifier = "week"

	IBCTimeoutHeightIncrement uint64 = 100
	ICATimeoutTimestamp              = time.Minute * 5

	CosmosValOperPrefix = "cosmosvaloper"
)

var (
	DelegationAccountPortID, _ = icatypes.NewControllerPortID(DelegationModuleAccount)
	RewardAccountPortID, _     = icatypes.NewControllerPortID(RewardModuleAccount)
)
var (
	// PortKey defines the key to store the port ID in store

	ModuleEnableKey           = []byte{0x01}
	HostChainParamsKey        = []byte{0x02}
	AllowListedValidatorsKey  = []byte{0x03}
	DelegationStateKey        = []byte{0x04}
	HostChainRewardAddressKey = []byte{0x05}
	IBCTransientStoreKey      = []byte{0x06}
)
