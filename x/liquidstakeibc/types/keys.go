package types

import (
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
	DelegationEpoch            = "day"
	UndelegationEpoch          = "day"
	RewardsEpochIdentifier     = "day"
	RedelegationEpochIdentifer = "day"
	CValueEpoch                = "hour"

	// ICA types
	DelegateICAType = "delegate"
	RewardsICAType  = "rewards"

	// ICQ query types
	// /key is required for proof generation
	StakingStoreQuery = "store/staking/key"
	BankStoreQuery    = "store/bank/key"

	// Host chain flags
	LSMFlag = "lsm"

	LiquidStakeDenomPrefix = "stk"

	IBCTimeoutTimestamp = 120 * time.Minute

	ICAMessagesChunkSize = 10

	IBCPrefix = "ibc" + "/"

	UnbondingStateEpochLimit = 4

	LSMDepositFilterLimit = 10000

	Percentage int64 = 100

	DaysInYear int64 = 365

	CValueDynamicLowerDiff int64 = 2

	CValueDynamicUpperDiff int64 = 10
)

// Consts for KV updates, update host chain
const (
	KeyAddValidator                string = "add_validator"
	KeyRemoveValidator             string = "remove_validator"
	KeyValidatorUpdate             string = "validator_update"
	KeyValidatorWeight             string = "validator_weight"
	KeyDepositFee                  string = "deposit_fee"
	KeyRestakeFee                  string = "restake_fee"
	KeyUnstakeFee                  string = "unstake_fee"
	KeyRedemptionFee               string = "redemption_fee"
	KeyLSMValidatorCap             string = "lsm_validator_cap"
	KeyLSMBondFactor               string = "lsm_bond_factor"
	KeyMaxEntries                  string = "max_entries"
	KeyUpperCValueLimit            string = "upper_c_value_limit"
	KeyLowerCValueLimit            string = "lower_c_value_limit"
	KeyRedelegationAcceptableDelta string = "redelegation_acceptable_delta"
	KeyMinimumDeposit              string = "min_deposit"
	KeyActive                      string = "active"
	KeySetWithdrawAddress          string = "set_withdraw_address"
	KeyAutocompoundFactor          string = "autocompound_factor"
	KeyFlags                       string = "flags"
	KeyRewardParams                string = "reward_params"
	KeyForceUpdateValidator        string = "force_update_validator"
	KeyForceUnbond                 string = "force_unbond"
	KeyForceICATransfer            string = "force_ica_transfer"
	KeyForceICATransferRewards     string = "force_ica_transfer_rewards"
	KeyForceTransferDeposits       string = "force_transfer_deposits"
	KeyForceTransferUnbonded       string = "force_transfer_unbonded"
	KeyForceFailUnbond             string = "force_fail_unbond"
)

var (
	HostChainKey          = []byte{0x01}
	DepositKey            = []byte{0x02}
	UnbondingKey          = []byte{0x03}
	UserUnbondingKey      = []byte{0x04}
	ValidatorUnbondingKey = []byte{0x05}
	ParamsKey             = []byte{0x06}
	LSMDepositKey         = []byte{0x07}
	RedelegationsKey      = []byte{0x08}
	RedelegationTxKey     = []byte{0x09}
)
