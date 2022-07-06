package orchestrator

import (
	"context"
	"fmt"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"strings"
	"testing"
	"time"
)

func TestLogicForRewardMap(t *testing.T) {

	rpcClientC, _ := newRPCClient("https://rpc.cosmoshub-4.audit.one:443", 1*time.Second)
	height := int64(11157228)

	blockResults, err := rpcClientC.BlockResults(context.Background(), &height)

	blockRes := blockResults.TxsResults

	rewardMap := make(map[string]sdkTypes.Coin)
	var valAddress string
	var amountCoin sdkTypes.Coin
	for _, txLog := range blockRes {
		eventList := txLog.Events
		logString := txLog.Log

		if strings.Contains(logString, "/cosmos.staking.v1beta1.MsgDelegate") {
			for _, events := range eventList {
				currEvent := events

				if currEvent.Type == "transfer" && string(currEvent.Attributes[0].Value) == "cosmos1q4xrgp7uluenws5dflmy7y7wlwcnerwzx4numj" {

					amountCoin, err = sdkTypes.ParseCoinNormalized(string(currEvent.Attributes[2].Value))
					if err != nil {
						panic(err)
					}

				}
				if currEvent.Type == "delegate" {
					valAddress = string(currEvent.Attributes[0].Value)
					val, ok := rewardMap[valAddress]
					if !ok {
						rewardMap[valAddress] = amountCoin
						fmt.Println(rewardMap)
					} else {
						rewardMap[valAddress] = val.Add(amountCoin)
					}
				}
			}
		}
	}

	fmt.Println(rewardMap)

	//txSlice := events["tx"]

	if err != nil {
		panic(err)
	}

}
