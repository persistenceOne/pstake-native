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

	// DelegationEpochIdentifier is the identifier for delegation epoch
	DelegationEpochIdentifier = "day"

	// RewardEpochIdentifier is the identifier for rewards epoch
	RewardEpochIdentifier = "day"

	// UndelegationEpochIdentifier is the identifier for undelegation epoch
	UndelegationEpochIdentifier = "day"

	// UndelegationEpochNumberFactor is the undelegation epoch number factor
	UndelegationEpochNumberFactor int64 = 4

	// UndelegationCompletionTimeBuffer is the undeleagation completion time buffer
	UndelegationCompletionTimeBuffer = time.Second * 60 //Does tendermint still have time drifts?

	// IBCTimeoutHeightIncrement is the IBC timeout height incerement
	IBCTimeoutHeightIncrement uint64 = 1000

	// ICATimeoutTimestamp is the ICA timeout time stamp
	ICATimeoutTimestamp = time.Minute

	// CosmosValOperPrefix is the prefix for cosmos validator address
	CosmosValOperPrefix = "cosmosvaloper"
)

// fee limits
var (
	MaxPstakeDepositFee    = sdk.MustNewDecFromStr("0.5")
	MaxPstakeRestakeFee    = sdk.MustNewDecFromStr("0.2")
	MaxPstakeUnstakeFee    = sdk.MustNewDecFromStr("0.5")
	MaxPstakeRedemptionFee = sdk.MustNewDecFromStr("0.2")
	MaxCValue              = sdk.MustNewDecFromStr("1.1")
)

var (
	// PortKey defines the key to store the port ID in store

	ModuleEnableKey                 = []byte{0x01} // key for module state
	HostChainParamsKey              = []byte{0x02} // key for host chain params
	AllowListedValidatorsKey        = []byte{0x03} // key for allow listed validators
	DelegationStateKey              = []byte{0x04} // key for delegation state
	HostChainRewardAddressKey       = []byte{0x05} // key for host chain address
	IBCTransientStoreKey            = []byte{0x06} // key for IBC transient store
	UnbondingEpochCValueKey         = []byte{0x07} // prefix for unbodning epoch c value store
	DelegatorUnbondingEpochEntryKey = []byte{0x08} // prefix for delegator unbonding epoch entry
	HostAccountsKey                 = []byte{0x09} // key for host accounts
)

// GetUnbondingEpochCValueKey returns a slice of byte made of UnbondingEpochCValueKey and epoch number
// coverted to bytes
func GetUnbondingEpochCValueKey(epochNumber int64) []byte {
	return append(UnbondingEpochCValueKey, []byte(strconv.FormatInt(epochNumber, 10))...)
}

// GetDelegatorUnbondingEpochEntryKey returns a slice of byte made of DelegatorUnbondingEpochEntryKey,
// delegator address as bytes and epoch number converted to bytes
func GetDelegatorUnbondingEpochEntryKey(delegatorAddress sdk.AccAddress, epochNumber int64) []byte {
	return append(append(DelegatorUnbondingEpochEntryKey, address.MustLengthPrefix(delegatorAddress)...), []byte(strconv.FormatInt(epochNumber, 10))...)
}

// GetPartialDelegatorUnbondingEpochEntryKey returns a slice of byte made of DelegatorUnbondingEpochEntryKey
// and delegator address as bytes
func GetPartialDelegatorUnbondingEpochEntryKey(delegatorAddress sdk.AccAddress) []byte {
	return append(DelegatorUnbondingEpochEntryKey, address.MustLengthPrefix(delegatorAddress)...)
}
