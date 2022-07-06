package orchestrator

import (
	"context"
	"fmt"
	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/persistenceOne/pstake-native/oracle/constants"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	abciTypes "github.com/tendermint/tendermint/abci/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	stdlog "log"
	"strings"
)

func GetValidatorDetails(chain *CosmosChain) []cosmosTypes.ValidatorDetails {
	var ValidatorDetailsArr []cosmosTypes.ValidatorDetails

	custodialAddr, err := Bech32ifyAddressBytes(chain.AccountPrefix, chain.CustodialAddress)
	if err != nil {
		stdlog.Println(err)
		panic(err)

	}
	grpcConn, err := grpc.Dial(chain.GRPCAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer func(grpcConn *grpc.ClientConn) {
		err := grpcConn.Close()
		if err != nil {
			stdlog.Println("GRPC Connection error")
		}
	}(grpcConn)

	if err != nil {
		stdlog.Println("GRPC Connection failed")
		panic(err)
	}

	stakingQueryClient := stakingTypes.NewQueryClient(grpcConn)

	stdlog.Println("staking query client connected")

	BondedTokensQueryResult, err := stakingQueryClient.DelegatorDelegations(context.Background(),
		&stakingTypes.QueryDelegatorDelegationsRequest{
			DelegatorAddr: custodialAddr,
			Pagination:    nil,
		},
	)

	if err != nil {
		panic(err)
	}
	flag := true
	for _, Delegations := range BondedTokensQueryResult.DelegationResponses {
		valAddr := Delegations.Delegation.ValidatorAddress
		BondedTokens := Delegations.Balance
		UnbondingTokensQueryResult, err := stakingQueryClient.UnbondingDelegation(context.Background(),
			&stakingTypes.QueryUnbondingDelegationRequest{
				DelegatorAddr: custodialAddr,
				ValidatorAddr: valAddr,
			},
		)

		if err != nil {
			statusErr := status.Convert(err)

			if statusErr.Code() == 5 {
				flag = false
			} else {
				stdlog.Println("cannot get unbonding delegations")
				panic(err)
			}

		}

		Unbondingtokens := types.ZeroInt()

		if flag == true {
			UnBondingEntries := UnbondingTokensQueryResult.Unbond.Entries
			for _, Entry := range UnBondingEntries {
				Unbondingtokens = Unbondingtokens.Add(Entry.Balance)
			}
		}

		newEntry := cosmosTypes.ValidatorDetails{
			ValidatorAddress: valAddr,
			BondedTokens:     BondedTokens,
			UnbondingTokens:  types.NewCoin(constants.CosmosDenom, Unbondingtokens),
		}
		ValidatorDetailsArr = append(ValidatorDetailsArr, newEntry)
		flag = true
	}

	return ValidatorDetailsArr
}

func PopulateRewards(chain *CosmosChain, valDetails []cosmosTypes.ValidatorDetails, blockRes []*abciTypes.ResponseDeliverTx) ([]cosmosTypes.ValidatorDetails, error) {
	rewardMap := make(map[string]sdkTypes.Coin)
	var valAddress string
	var amountCoin sdkTypes.Coin

	custodialAddr, err := Bech32ifyAddressBytes(chain.AccountPrefix, chain.CustodialAddress)
	if err != nil {
		return nil, err
	}

	for _, txLog := range blockRes {
		eventList := txLog.Events
		logString := txLog.Log
		if strings.Contains(logString, "/cosmos.staking.v1beta1.MsgDelegate") {
			for _, events := range eventList {
				currEvent := events
				if currEvent.Type == "transfer" && string(currEvent.Attributes[0].Value) == custodialAddr {
					amountCoin, err = sdkTypes.ParseCoinNormalized(string(currEvent.Attributes[2].Value))
					if err != nil {
						return nil, err
					}
				}
				if currEvent.Type == "delegate" {
					valAddress = string(currEvent.Attributes[0].Value)
					val, ok := rewardMap[valAddress]
					if !ok {
						rewardMap[valAddress] = amountCoin
					} else {
						rewardMap[valAddress] = val.Add(amountCoin)
					}
				}
			}
		}
	}
	for _, valDetails := range valDetails {
		valAddr := valDetails.ValidatorAddress
		valDetails.RewardsCollected = rewardMap[valAddr]
	}

	return valDetails, nil
}

func GetAccountDetails(cosmosClient cosmosClient.Context, chain *CosmosChain, addr string) (seqNum, accountNum uint64) {
	grpcConn, err := grpc.Dial(chain.GRPCAddr, grpc.WithInsecure())
	defer func(grpcConn *grpc.ClientConn) {
		err := grpcConn.Close()
		if err != nil {
			stdlog.Println("GRPC Connection error")
		}
	}(grpcConn)

	if err != nil {
		stdlog.Println("GRPC Connection failed")
		panic(err)

	}
	authQueryClient := authTypes.NewQueryClient(grpcConn)

	fmt.Println("account query client connected")

	AccountDetails, err := authQueryClient.Account(context.Background(),
		&authTypes.QueryAccountRequest{Address: addr},
	)

	if err != nil {
		stdlog.Println("cannot get accounts")

	}
	var account authTypes.AccountI
	Account := AccountDetails.Account

	err = cosmosClient.InterfaceRegistry.UnpackAny(Account, &account)

	if err != nil {

		stdlog.Println("err unmarshalling ANY account")

	}

	return account.GetSequence(), account.GetAccountNumber()

}
