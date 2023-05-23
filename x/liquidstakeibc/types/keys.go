package types

import (
	"strconv"
	"time"
)

const (
	// ModuleName defines the module name
	ModuleName = "liquidstakeibc"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// DepositModuleAccount DepositModuleAccountName
	DepositModuleAccount = ModuleName + "_deposit_account"

	// UndelegationModuleAccount UndelegationModuleAccountName
	UndelegationModuleAccount = ModuleName + "_undelegation_account"

	// Epoch identifiers
	DelegationEpoch        = "day"
	UndelegationEpoch      = "day"
	RewardsEpochIdentifier = "day"

	// ICA types
	DelegateICAType = "delegate"
	RewardsICAType  = "rewards"

	IBCTimeoutHeightIncrement uint64 = 1000

	ICATimeoutTimestamp = 15 * time.Minute

	UnbondingStateEpochLimit = 4
)

var (
	HostChainKey          = []byte{0x01}
	DepositKey            = []byte{0x02}
	UnbondingKey          = []byte{0x03}
	UserUnbondingKey      = []byte{0x04}
	ValidatorUnbondingKey = []byte{0x05}
)

func GetUnbondingStoreKey(chainID string, epochNumber int64) []byte {
	return append([]byte(chainID), []byte(strconv.FormatInt(epochNumber, 10))...)
}

func GetUserUnbondingStoreKey(chainID, delegatorAddress string, epochNumber int64) []byte {
	return append([]byte(chainID), append([]byte(delegatorAddress), []byte(strconv.FormatInt(epochNumber, 10))...)...)
}

func GetValidatorUnbondingStoreKey(chainID, validatorAddress string, epochNumber int64) []byte {
	return append([]byte(chainID), append([]byte(validatorAddress), []byte(strconv.FormatInt(epochNumber, 10))...)...)
}
