/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package types

import (
	"encoding/binary"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/tendermint/tendermint/crypto/tmhash"
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
	QueryTxByID     = "txByID"
	QueryProposal   = "proposal"
	QueryVote       = "vote"
	QueryVotes      = "votes"

	MintDenom  = "pstake" //TODO shift to params
	StakeDenom = "uatom"  //TODO shift to params

	MinimumRatioForMajority = 0.66

	StorageWindow = 100 //TODO : Revert Back to 100

	Bech32Prefix = "cosmos"
)

var (
	KeyValidatorAddress = "KeyValidatorAddress"

	SequenceKeyPrefix = "SequenceKeyPrefix"

	KeyLastTXPoolID = SequenceKeyPrefix + "lastTxPoolId"

	KeyCosmosValidatorSet = []byte{0x01}

	KeyTotalDelegationTillDate = []byte{0x02}

	OutgoingTXPoolKey = []byte{0x03}

	AddressAndAmountStoreKey = []byte{0x04}

	MintingPoolStoreKey = []byte{0x05}

	ValidatorOrchestratorStoreKey = []byte{0x07}

	ProposalStoreKey = []byte{0x08}

	ProposalIDKey = []byte{0x09}

	VotingParams = []byte{0xA}

	ProposalsKeyPrefix = []byte{0xB}

	ActiveProposalQueuePrefix = []byte{0xC}

	VotesKeyPrefix = []byte{0xD}

	HashAndIDStore = []byte{0xE}

	KeyOutgoingUnbondStore = []byte{0xF}

	KeyStakingEpochStore = []byte{0x10}

	KeyMintingEpochStore = []byte{0x12}

	KeyRewardsStore = []byte{0x13}

	KeyCurrentEpochRewardsStore = []byte{0x14}

	KeyEpochStoreForUndelegation = "EpochStoreForUndelegation"

	KeyEpochStoreForWithdrawSuccess = []byte{0x10}

	KeyUndelegateSuccessStore = []byte{0x11}

	KeyWithdrawStore = []byte{0x12}

	KeyOutgoingTxSignature = []byte{0x11}

	KeyOutgoingSignaturePoolKey = []byte{0x12}
)

func GetEpochStoreForUndelegationKey(epochNumber int64) []byte {
	return append([]byte(KeyEpochStoreForUndelegation), Int64Bytes(epochNumber)...)
}

func ConvertByteArrToString(value []byte) string {
	var ret strings.Builder
	for i := 0; i < len(value); i++ {
		ret.WriteString(string(value[i]))
	}
	return ret.String()
}

func GetChainIDTxHashBlockHeightKey(chainID string, blockHeight int64, txHash string) string {
	return chainID + strconv.FormatInt(blockHeight, 10) + txHash
}

func GetChainIDAndBlockHeightKey(chainID string, blockHeight int64) string {
	return chainID + strconv.FormatInt(blockHeight, 10)
}

func GetDestinationAddressAmountAndTxHashKey(destinationAddress sdkTypes.AccAddress, coins sdkTypes.Coins, txHash string) string {
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
	return append(ActiveProposalQueuePrefix, sdkTypes.FormatTimeBytes(endTime)...)
}

// VoteKey key of a specific vote from the store
func VoteKey(proposalID uint64, voterAddr sdkTypes.AccAddress) []byte {
	return append(VotesKey(proposalID), address.MustLengthPrefix(voterAddr.Bytes())...)
}

// VotesKey gets the first part of the votes key based on the proposalID
func VotesKey(proposalID uint64) []byte {
	return append(VotesKeyPrefix, GetProposalIDBytes(proposalID)...)
}

// OutgoingTxSignatureKey forms a key from txID to the substore
func OutgoingTxSignatureKey(txID uint64) []byte {
	return append(KeyOutgoingTxSignature, GetProposalIDBytes(txID)...)
}

func BytesToHexUpper(bz []byte) string {
	return hex.EncodeToString(tmhash.Sum(bz))
}
