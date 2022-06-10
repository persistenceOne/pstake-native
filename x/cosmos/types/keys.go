/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package types

import (
	"encoding/binary"
	"encoding/hex"
	"strconv"
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
	QueryProposals  = "proposals"
	QueryVote       = "vote"
	QueryVotes      = "votes"

	StorageWindow = 20000

	Bech32Prefix = "cosmos"

	Bech32PrefixAccAddr  = Bech32Prefix
	Bech32PrefixAccPub   = Bech32Prefix + sdkTypes.PrefixPublic
	Bech32PrefixValAddr  = Bech32Prefix + sdkTypes.PrefixValidator + sdkTypes.PrefixOperator
	Bech32PrefixValPub   = Bech32Prefix + sdkTypes.PrefixValidator + sdkTypes.PrefixOperator + sdkTypes.PrefixPublic
	Bech32PrefixConsAddr = Bech32Prefix + sdkTypes.PrefixValidator + sdkTypes.PrefixConsensus
	Bech32PrefixConsPub  = Bech32Prefix + sdkTypes.PrefixValidator + sdkTypes.PrefixConsensus + sdkTypes.PrefixPublic
)

var (
	MinimumRatioForMajority = sdkTypes.NewDec(66).Quo(sdkTypes.NewDec(100))

	SequenceKeyPrefix = "SequenceKeyPrefix"

	KeyLastTXPoolID = SequenceKeyPrefix + "lastTxPoolId"

	KeyTotalDelegationTillDate = []byte{0x02}

	OutgoingTXPoolKey = []byte{0x03}

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

	KeyEpochStoreForWithdrawSuccess = []byte{0x15}

	KeyUndelegateSuccessStore = []byte{0x16}

	KeyWithdrawStore = []byte{0x17}

	KeyOutgoingTxSignature = []byte{0x18}

	KeyOutgoingSignaturePoolKey = []byte{0x19}

	KeyMultisigAccountStore = []byte{0x20}

	KeyCurrentMultisigAddress = []byte{0x21}

	KeyTransactionQueue = []byte{0x22}

	KeyCosmosValidatorWeights = []byte{0x23}

	KeyNativeValidatorWeights = []byte{0x24}

	KeySlashingStore = []byte{0x25}

	KeyMintTokenStore = []byte{0x26}

	KeyOracleLastUpdateHeightNative = []byte{0x27}

	KeyOracleLastUpdateHeightCosmos = []byte{0x28}
)

func GetEpochStoreForUndelegationKey(epochNumber int64) []byte {
	return append([]byte(KeyEpochStoreForUndelegation), Int64Bytes(epochNumber)...)
}

func GetChainIDAndBlockHeightKey(chainID string, blockHeight int64) string {
	return chainID + strconv.FormatInt(blockHeight, 10)
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

// MultisigAccountStoreKey turn an address to key used to get it from the account store
func MultisigAccountStoreKey(addr sdkTypes.AccAddress) []byte {
	return append(KeyMultisigAccountStore, addr.Bytes()...)
}

// CurrentMultisigAddressKey turn an address to that is expected to send current txns
func CurrentMultisigAddressKey() []byte {
	return KeyCurrentMultisigAddress
}

func BytesToHexUpper(bz []byte) string {
	return hex.EncodeToString(tmhash.Sum(bz))
}
