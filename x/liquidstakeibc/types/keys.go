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

	// Epoch identifiers
	DelegationEpoch = "day"

	// ICA types
	DelegateICAType = "delegate"
	RewardsICAType  = "rewards"

	IBCTimeoutHeightIncrement uint64 = 1000

	ICATimeoutTimestamp = 15 * time.Minute
)

var (
	HostChainKey = []byte{0x01}
	DepositKey   = []byte{0x02}
)
