/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package types

import (
	"strconv"
	"strings"

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

	MintDenom = "pstake"
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

	AddressAndAmountKey = "AddressAndAmountKey"

	MintingPoolStoreKey = "MintingPoolStoreKey"

	OrchestratorValidatorStoreKey = "OrchestratorValidatorStoreKey"
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

func GetDestinationAddressAndAmountKey(destinationAddress sdk.AccAddress, coins sdk.Coins) string {
	amount := make([]byte, 32)
	amount = []byte(coins[0].Amount.String())

	a := append(destinationAddress.Bytes(), amount...)
	return ConvertByteArrToString(a)
}

func ValidVoteOption(option VoteOption) bool {
	if option == OptionYes ||
		option == OptionAbstain ||
		option == OptionNo ||
		option == OptionNoWithVeto {
		return true
	}
	return false
}
