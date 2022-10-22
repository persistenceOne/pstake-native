package types

import (
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
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

	// MsgTypeRedeem is the type of message redeem
	MsgTypeRedeem = "msg_redeem"

	// MsgTypeClaim is the type of message claim
	MsgTypeClaim = "msg_claim"

	// MsgTypeJumpStart is the type of message Jump start
	MsgTypeJumpStart = "msg_jump_start"

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

	DelegationEpochIdentifier              = "day"
	RewardEpochIdentifier                  = "day"
	UndelegationEpochIdentifier            = "day"
	UndelegationEpochNumberFactor    int64 = 4
	UndelegationCompletionTimeBuffer       = time.Second * 60 //Does tendermint still have time drifts?

	IBCTimeoutHeightIncrement uint64 = 1000
	ICATimeoutTimestamp              = time.Minute

	CosmosValOperPrefix = "cosmosvaloper"
)

var (
	MaxCValue = sdk.MustNewDecFromStr("1.1")
)

var (
	// PortKey defines the key to store the port ID in store

	ModuleEnableKey                 = []byte{0x01}
	HostChainParamsKey              = []byte{0x02}
	AllowListedValidatorsKey        = []byte{0x03}
	DelegationStateKey              = []byte{0x04}
	HostChainRewardAddressKey       = []byte{0x05}
	IBCTransientStoreKey            = []byte{0x06}
	UnbondingEpochCValueKey         = []byte{0x07}
	DelegatorUnbondingEpochEntryKey = []byte{0x08}
	HostAccountsKey                 = []byte{0x09}
)

func GetUnbondingEpochCValueKey(epochNumber int64) []byte {
	return append(UnbondingEpochCValueKey, []byte(strconv.FormatInt(epochNumber, 10))...)
}

func GetDelegatorUnbondingEpochEntryKey(delegatorAddress sdk.AccAddress, epochNumber int64) []byte {
	return append(append(DelegatorUnbondingEpochEntryKey, address.MustLengthPrefix(delegatorAddress)...), []byte(strconv.FormatInt(epochNumber, 10))...)
}

func GetPartialDelegatorUnbondingEpochEntryKey(delegatorAddress sdk.AccAddress) []byte {
	return append(DelegatorUnbondingEpochEntryKey, address.MustLengthPrefix(delegatorAddress)...)
}
