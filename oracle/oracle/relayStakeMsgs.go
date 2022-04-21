package oracle

//
//import sdk "github.com/cosmos/cosmos-sdk/types"
//
//type DeliverMsgsAction struct {
//	DeliverMsgs []DeliverMsg
//}
//
//type RelayStakeMsgs struct {
//	CosmosMsg    []sdk.Msg `json:"cosmos-msg"`
//	MaxTxSize    uint64    `json:"max-tx-size"`
//	MaxMsgLength uint64    `json:"max-msg-length"`
//
//	Success bool `json:"success"`
//}
//
//func NewRelayStakeMsgs() *RelayStakeMsgs {
//	return &RelayStakeMsgs{
//		CosmosMsg: []sdk.Msg{},
//		Success:   false,
//	}
//}
//
//func (r *RelayStakeMsgs) Ready() bool {
//	if r == nil {
//		return false
//	}
//
//	if len(r.CosmosMsg) == 0 {
//		return false
//	}
//	return true
//}
