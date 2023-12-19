package types

import (
	"encoding/binary"
	"time"
)

const (
	// ModuleName defines the module name
	ModuleName = "ratesync"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_ratesync"

	LiquidStakeAllowAllDenoms = "*"
	LiquidStakeEpoch          = "day"
	DefaultPortOwnerPrefix    = "pstake_ratesync_"

	ICATimeoutTimestamp = 60 * time.Minute
)

var (
	HostChainIDKeyPrefix = []byte{0x01}
	HostChainKeyPrefix   = []byte{0x02}
	ParamsKeyPrefix      = []byte{0x00}
)

// HostChainKey returns the store key to retrieve a Chain from the index fields
func HostChainKey(
	id uint64,
) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}
