/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package types

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"strconv"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// ModuleName Module Name
	ModuleName = "cosmos"

	// DefaultParamspace params keeper
	DefaultParamspace = ModuleName

	// StoreKey is the default store key for cosmos
	StoreKey = ModuleName

	// RouterKey is the module name router key
	RouterKey = ModuleName

	// QuerierRoute is the querier route for the cosmos store.
	QuerierRoute = StoreKey

	// QueryParameters Query endpoints supported by the cosmos querier
	QueryParameters = "parameters"
	QueryTxByID     = "byTxID"

	MintDenom  = "pstake" //TODO shift to params
	StakeDenom = "uatom"  //TODO shift to params

	MinimumRatioForMajority = 0.66

	StorageWindow = 50
)

var (
	KeyOrchestratorAddress = "KeyOrchestratorAddress"
	KeyAccAddress          = "KeyAccAddress"
	OutgoingTxPrefix       = []byte{0x01}
	IncomingTxPrefix       = []byte{0x02}

	// SequenceKeyPrefix indexes different txids
	SequenceKeyPrefix = "SequenceKeyPrefix"

	// KeyLastTXPoolID indexes the lastTxPoolID
	KeyLastTXPoolID = SequenceKeyPrefix + "lastTxPoolId"

	// OutgoingTXPoolKey indexes the last nonce for the outgoing tx pool
	OutgoingTXPoolKey = "OutgoingTXPoolKey"

	AddressAndAmountStoreKey = "AddressAndAmountKey"

	MintingPoolStoreKey = "MintingPoolStoreKey"

	OrchestratorValidatorStoreKey = "OrchestratorValidatorStoreKey"

	ProposalStoreKey = "ProposalStoreKey"

	ProposalIDKey = []byte("ProposalIDKey")

	VotingParams = []byte("VotingParams")

	ProposalsKeyPrefix = []byte("ProposalsKeyPrefix")

	ActiveProposalQueuePrefix = []byte("ActiveProposalQueuePrefix")

	VotesKeyPrefix = []byte("VotesKeyPrefix")

	HashAndIDStore = []byte("HashAndIDStore")
)

func ConvertByteArrToString(value []byte) string {
	var ret strings.Builder
	for i := 0; i < len(value); i++ {
		ret.WriteString(string(value[i]))
	}
	return ret.String()
}

func GetOrchestratorAddressKey(orc sdk.AccAddress) string {
	if err := sdk.VerifyAddressFormat(orc); err != nil {
		panic(sdkErrors.Wrap(err, "invalid orchestrator address"))
	}
	return KeyOrchestratorAddress + string(orc.Bytes())
}

func GetChainIDTxHashBlockHeightKey(chainID string, blockHeight int64, txHash string) string {
	return chainID + strconv.FormatInt(blockHeight, 10) + txHash
}

func GetOutgoingTxPoolKey(fee sdk.Coin, id uint64) string {
	// sdkInts have a size limit of 255 bits or 32 bytes
	// therefore this will never panic and is always safe
	amount := make([]byte, 32)
	amount = []byte(fee.Amount.String())

	a := append(amount, UInt64Bytes(id)...)
	b := append([]byte(OutgoingTXPoolKey), a...)
	return ConvertByteArrToString(b)
}

func GetDestinationAddressAmountAndTxHashKey(destinationAddress sdk.AccAddress, coins sdk.Coins, txHash string) string {
	amount := make([]byte, 32)
	amount = []byte(coins[0].Amount.String())

	a := append(destinationAddress.Bytes(), amount...)
	b := append([]byte(txHash), a...)
	return ConvertByteArrToString(b)
}

// GetProposalIDFromBytes returns proposalID in uint64 format from a byte array
func GetProposalIDFromBytes(bz []byte) (proposalID uint64) {
	return binary.BigEndian.Uint64(bz)
}

// ProposalKey1 gets a specific proposal from the store
func ProposalKey1(proposalID uint64) []byte {
	return append(ProposalsKeyPrefix, GetProposalIDBytes(proposalID)...)
}

// GetProposalIDBytes returns the byte representation of the proposalID
func GetProposalIDBytes(proposalID uint64) (proposalIDBz []byte) {
	proposalIDBz = make([]byte, 8)
	binary.BigEndian.PutUint64(proposalIDBz, proposalID)
	return
}

// ActiveProposalQueueKey returns the key for a proposalID in the activeProposalQueue
func ActiveProposalQueueKey(proposalID uint64, endTime time.Time) []byte {
	return append(ActiveProposalByTimeKey(endTime), GetProposalIDBytes(proposalID)...)
}

// ActiveProposalByTimeKey gets the active proposal queue key by endTime
func ActiveProposalByTimeKey(endTime time.Time) []byte {
	return append(ActiveProposalQueuePrefix, sdk.FormatTimeBytes(endTime)...)
}

// VoteKey key of a specific vote from the store
func VoteKey(proposalID uint64, voterAddr sdk.AccAddress) []byte {
	return append(VotesKey(proposalID), address.MustLengthPrefix(voterAddr.Bytes())...)
}

// VotesKey gets the first part of the votes key based on the proposalID
func VotesKey(proposalID uint64) []byte {
	return append(VotesKeyPrefix, GetProposalIDBytes(proposalID)...)
}

func BytesToHexUpper(bz []byte) string {
	return hex.EncodeToString(tmhash.Sum(bz))
}
