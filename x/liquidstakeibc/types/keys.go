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
	CValueEpoch            = "hour"

	// ICA types
	DelegateICAType = "delegate"
	RewardsICAType  = "rewards"

	// ICQ query types
	// /key is required for proof generation
	StakingStoreQuery = "store/staking/key"
	BankStoreQuery    = "store/bank/key"

	LiquidStakeDenomPrefix = "stk"

	IBCTimeoutHeightIncrement uint64 = 1000

	ICATimeoutTimestamp = 15 * time.Minute

	UnbondingStateEpochLimit = 4
)

// Consts for KV updates, update host chain
const (
	KeyAddValidator       string = "add_validator"
	KeyRemoveValidator    string = "remove_validator"
	KeyValidatorSlashing  string = "validator_slashing"
	KeyValidatorWeight    string = "validator_weight"
	KeyDepositFee         string = "deposit_fee"
	KeyRestakeFee         string = "restake_fee"
	KeyUnstakeFee         string = "unstake_fee"
	KeyRedemptionFee      string = "redemption_fee"
	KeyMinimumDeposit     string = "min_deposit"
	KeyActive             string = "active"
	KeySetWithdrawAddress string = "set_withdraw_address"
	KeyAutocompoundFactor string = "autocompound_factor"
)

var (
	HostChainKey          = []byte{0x01}
	DepositKey            = []byte{0x02}
	UnbondingKey          = []byte{0x03}
	UserUnbondingKey      = []byte{0x04}
	ValidatorUnbondingKey = []byte{0x05}
	ParamsKey             = []byte{0x06}
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

func GetDepositStoreKey(chainID string, epochNumber int64) []byte {
	return append([]byte(chainID), []byte(strconv.FormatInt(epochNumber, 10))...)
}
