package orchestrator

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	prov "github.com/tendermint/tendermint/light/provider/http"
)

func TestLogicForRewardMap(t *testing.T) {

	rpcClientC, _ := newRPCClient("http://13.229.229.141:12001", 1*time.Second)
	height := int64(3657)

	blockResults, err := rpcClientC.BlockResults(context.Background(), &height)

	blockRes := blockResults.TxsResults

	rewardMap := make(map[string]sdkTypes.Coin)
	var valAddress string
	var amountCoin sdkTypes.Coin
	for _, txLog := range blockRes {
		eventList := txLog.Events
		logString := txLog.Log

		if strings.Contains(logString, "/cosmos.authz.v1beta1.MsgExec") {
			for _, events := range eventList {
				currEvent := events

				if currEvent.Type == "transfer" && string(currEvent.Attributes[0].Value) == "cosmos15ddw7dkp56zytf3peshxr8fwn5w76y4g462ql2" {

					amountCoin, err = sdkTypes.ParseCoinNormalized(string(currEvent.Attributes[2].Value))
					if err != nil {
						panic(any(err))
					}

				}
				if currEvent.Type == "delegate" || currEvent.Type == "undelegate" {
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

	custodialAdrr, err := AccAddressFromBech32("cosmos15ddw7dkp56zytf3peshxr8fwn5w76y4g462ql2", "cosmos")
	liteproviderC, _ := prov.New("native", "http://13.229.229.141:12001")
	chainC := &CosmosChain{
		Key:              "unusedNativeKey",
		ChainID:          "test",
		RPCAddr:          "http://13.229.229.141:12001",
		CustodialAddress: custodialAdrr,
		AccountPrefix:    "cosmos",
		GasAdjustment:    1.0,
		GasPrices:        "0.025stake",
		GRPCAddr:         "13.229.229.141:12344",
		CoinType:         118,
		HomePath:         "",
		KeyBase:          nil,
		Client:           rpcClientC,
		Encoding:         params.EncodingConfig{},
		Provider:         liteproviderC,
		address:          nil,
		logger:           nil,
		timeout:          0,
		debug:            false,
	}

	cosmosEncodingConfig := chainC.MakeEncodingConfig()
	chainC.Encoding = cosmosEncodingConfig
	chainC.logger = defaultChainLogger()

	valDetails := GetValidatorDetails(chainC)

	vaDetails, err := PopulateRewards(chainC, valDetails, blockRes)
	if err != nil {
		panic(any(err))
	}

	fmt.Println(vaDetails)

	if err != nil {
		panic(any(err))
	}

}
