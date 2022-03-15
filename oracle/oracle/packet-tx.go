package oracle

//
//import (
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	"github.com/cosmos/oracle/oracle"
//	"time"
//)
//
//func SendStakeMsg(native *oracle.Chain, cosmos *oracle.Chain, amount sdk.Coin, delegatorAddress, validatorAddress string, toHeightOffset uint64, toTimeOffset time.Duration) error {
//	var (
//		timeoutHeight    uint64
//		timeoutTimestamp uint64
//	)
//
//	cosmosH, err := cosmos.QueryLatestHeight()
//	if err != nil {
//		return err
//	}
//	h, err := cosmos.GetIBCUpdateHeader(native, cosmosH)
//	if err != nil {
//		return err
//	}
//
//	switch {
//	case toHeightOffset > 0 && toTimeOffset > 0:
//		timeoutHeight = uint64(h.Header.Height) + toHeightOffset
//		timeoutTimestamp = uint64(time.Now().Add(toTimeOffset).UnixNano())
//	case toHeightOffset > 0:
//		timeoutHeight = uint64(h.Header.Height) + toHeightOffset
//		timeoutTimestamp = 0
//	case toTimeOffset > 0:
//		timeoutHeight = 0
//		timeoutTimestamp = uint64(time.Now().Add(toTimeOffset).UnixNano())
//	case toHeightOffset == 0 && toTimeOffset == 0:
//		timeoutHeight = uint64(h.Header.Height + 1000)
//		timeoutTimestamp = 0
//	}
//	txs := RelayMsgs{
//		Src: []sdk.Msg{MsgDelegate(cosmos, amount, delegatorAddress, validatorAddress)},
//	}
//
//}
//
//func SendUnStakeMsg(native *oracle.Chain, cosmos *oracle.Chain, amount sdk.Coin, delegatorAddress, validatorAddress string, toHeightOffset uint64, toTimeOffset time.Duration) error {
//	var (
//		timeoutHeight    uint64
//		timeoutTimestamp uint64
//	)
//
//	cosmosH, err := cosmos.QueryLatestHeight()
//	if err != nil {
//		return err
//	}
//	h, err := cosmos.GetIBCUpdateHeader(native, cosmosH)
//	if err != nil {
//		return err
//	}
//
//	switch {
//	case toHeightOffset > 0 && toTimeOffset > 0:
//		timeoutHeight = uint64(h.Header.Height) + toHeightOffset
//		timeoutTimestamp = uint64(time.Now().Add(toTimeOffset).UnixNano())
//	case toHeightOffset > 0:
//		timeoutHeight = uint64(h.Header.Height) + toHeightOffset
//		timeoutTimestamp = 0
//	case toTimeOffset > 0:
//		timeoutHeight = 0
//		timeoutTimestamp = uint64(time.Now().Add(toTimeOffset).UnixNano())
//	case toHeightOffset == 0 && toTimeOffset == 0:
//		timeoutHeight = uint64(h.Header.Height + 1000)
//		timeoutTimestamp = 0
//	}
//	txs := RelayMsgs{
//		Src: []sdk.Msg{MsgUndelegate(cosmos, amount, delegatorAddress, validatorAddress)},
//	}
//
//}
//
////func (c CosmosChain) GetChainID() string {
////	return "cosmos"
////}
//
////func (c CosmosChain) SendStakeMsg(dst *Chain, amount sdk.Coin, dstAddr string, toHeightOffset uint64, toTimeOffset time.Duration) error {
////	return nil
//////}
////
////func (c *Chain) SendUnstakeMsg(dst *Chain, amount sdk.Coin, dstAddr string, toHeightOffset uint64, toTimeOffset time.Duration) error {
////	return nil
////}
