package oracle

import (
	"context"
	"fmt"
	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/persistenceOne/pstake-native/oracle/constants"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	logg "log"
)

func GetValidatorDetails(chain *CosmosChain) []cosmosTypes.ValidatorDetails {
	var ValidatorDetailsArr []cosmosTypes.ValidatorDetails

	custodialAddr, err := Bech32ifyAddressBytes(chain.AccountPrefix, chain.CustodialAddress)
	if err != nil {
		logg.Println(err)
		panic(err)

	}
	grpcConn, err := grpc.Dial(chain.GRPCAddr, grpc.WithInsecure())
	defer func(grpcConn *grpc.ClientConn) {
		err := grpcConn.Close()
		if err != nil {
			logg.Println("GRPC Connection error")
		}
	}(grpcConn)

	if err != nil {
		logg.Println("GRPC Connection failed")
		panic(err)
	}

	stakingQueryClient := stakingTypes.NewQueryClient(grpcConn)

	logg.Println("staking query client connected")

	BondedTokensQueryResult, err := stakingQueryClient.DelegatorDelegations(context.Background(),
		&stakingTypes.QueryDelegatorDelegationsRequest{
			DelegatorAddr: custodialAddr,
			Pagination:    nil,
		},
	)

	if err != nil {
		logg.Println("cannot get total delegations")
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
				logg.Println("cannot get unbonding delegations")
				panic(err)
			}

		}

		Unbondingtokens := types.NewInt(0)

		if flag == true {
			UnBondingEntries := UnbondingTokensQueryResult.Unbond.Entries
			for _, Entry := range UnBondingEntries {
				Unbondingtokens.Add(Entry.Balance)
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

func GetAccountDetails(cosmosClient cosmosClient.Context, chain *CosmosChain, addr string) (seqNum, accountNum uint64) {
	grpcConn, err := grpc.Dial(chain.GRPCAddr, grpc.WithInsecure())
	defer func(grpcConn *grpc.ClientConn) {
		err := grpcConn.Close()
		if err != nil {
			logg.Println("GRPC Connection error")
		}
	}(grpcConn)

	if err != nil {
		logg.Println("GRPC Connection failed")
		panic(err)

	}
	authQueryClient := authTypes.NewQueryClient(grpcConn)

	fmt.Println("account query client connected")

	AccountDetails, err := authQueryClient.Account(context.Background(),
		&authTypes.QueryAccountRequest{Address: addr},
	)

	if err != nil {
		logg.Println("cannot get accounts")

	}
	var account authTypes.AccountI
	Account := AccountDetails.Account

	err = cosmosClient.InterfaceRegistry.UnpackAny(Account, &account)

	if err != nil {

		fmt.Println("err unmarshaling ANY account")

	}

	return account.GetSequence(), account.GetAccountNumber()

}
