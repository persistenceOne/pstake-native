/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package types

import (
	"encoding/binary"
	"encoding/hex"
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

	StorageWindow = 40000

	MinGasFee = 800000
	MaxGasFee = 4000000

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

	KeyEpochStoreForUndelegation = "EpochStoreForUndelegation"

	OutgoingTXPoolKey = []byte{0x01}

	ValidatorOrchestratorStoreKey = []byte{0x02}

	ProposalStoreKey = []byte{0x03}

	ProposalIDKey = []byte{0x04}

	ProposalsKeyPrefix = []byte{0x05}

	ActiveProposalQueuePrefix = []byte{0x06}

	VotesKeyPrefix = []byte{0x07}

	HashAndIDStore = []byte{0x08}

	KeyOutgoingUnbondStore = []byte{0x9}

	KeyStakingEpochStore = []byte{0xA}

	KeyCurrentEpochRewardsStore = []byte{0xB}

	KeyEpochStoreForWithdrawSuccess = []byte{0xC}

	KeyUndelegateSuccessStore = []byte{0xD}

	KeyWithdrawStore = []byte{0x2F}

	KeyOutgoingSignaturePoolKey = []byte{0xF}

	KeyMultisigAccountStore = []byte{0x11}

	KeyCurrentMultisigAddress = []byte{0x12}

	KeyTransactionQueue = []byte{0x13}

	KeyCosmosValidatorWeights = []byte{0x14}

	KeyNativeValidatorWeights = []byte{0x15}

	KeySlashingStore = []byte{0x16}

	KeyMintTokenStore = []byte{0x17}

	KeyOrchestratorLastUpdateHeightNative = []byte{0x18}

	KeyOrchestratorLastUpdateHeightCosmos = []byte{0x19}

	KeyMintedAmount = []byte{0x1A}

	KeyVirtuallyStakedAmount = []byte{0x1B}

	KeyStakedAmount = []byte{0x1C}

	KeyVirtuallyUnbonded = []byte{0x1D}

	KeyCosmosBalances = []byte{0x1E}
)

func GetEpochStoreForUndelegationKey(epochNumber int64) []byte {
	return append([]byte(KeyEpochStoreForUndelegation), Int64Bytes(epochNumber)...)
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
