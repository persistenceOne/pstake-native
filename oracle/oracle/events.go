package oracle

//
//import (
//	"fmt"
//
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	clienttypes "github.com/cosmos/ibc-go/v2/modules/core/02-client/types"
//	connectiontypes "github.com/cosmos/ibc-go/v2/modules/core/03-connection/types"
//	channeltypes "github.com/cosmos/ibc-go/v2/modules/core/04-channel/types"
//)
//
//.
//func ParseClientIDFromEvents(events sdk.StringEvents) (string, error) {
//	for _, ev := range events {
//		if ev.Type == clienttypes.EventTypeCreateClient {
//			for _, attr := range ev.Attributes {
//				if attr.Key == clienttypes.AttributeKeyClientID {
//					return attr.Value, nil
//				}
//			}
//		}
//	}
//	return "", fmt.Errorf("client identifier event attribute not found")
//}
//
//func ParseConnectionIDFromEvents(events sdk.StringEvents) (string, error) {
//	for _, ev := range events {
//		if ev.Type == connectiontypes.EventTypeConnectionOpenInit ||
//			ev.Type == connectiontypes.EventTypeConnectionOpenTry {
//			for _, attr := range ev.Attributes {
//				if attr.Key == connectiontypes.AttributeKeyConnectionID {
//					return attr.Value, nil
//				}
//			}
//		}
//	}
//	return "", fmt.Errorf("connection identifier event attribute not found")
//}
//
//
//func ParseChannelIDFromEvents(events sdk.StringEvents) (string, error) {
//	for _, ev := range events {
//		if ev.Type == channeltypes.EventTypeChannelOpenInit || ev.Type == channeltypes.EventTypeChannelOpenTry {
//			for _, attr := range ev.Attributes {
//				if attr.Key == channeltypes.AttributeKeyChannelID {
//					return attr.Value, nil
//				}
//			}
//		}
//	}
//	return "", fmt.Errorf("channel identifier event attribute not found")
//}
