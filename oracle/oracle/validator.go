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
	logg "log"
)

func GetValidatorDetails(chain *CosmosChain) []cosmosTypes.ValidatorDetails {
	var ValidatorDetailsArr []cosmosTypes.ValidatorDetails

	custodialAddr := chain.CustodialAddress.String()
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

	fmt.Println("staking query client connected")

	BondedTokensQueryResult, err := stakingQueryClient.DelegatorDelegations(context.Background(),
		&stakingTypes.QueryDelegatorDelegationsRequest{
			DelegatorAddr: custodialAddr,
			Pagination:    nil,
		},
	)

	if err != nil {
		logg.Println("cannot get total delegations")
	}

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
			logg.Println("cannot get unbonding delegations")

		}

		UnBondingEntries := UnbondingTokensQueryResult.Unbond.Entries
		Unbondingtokens := types.NewInt(0)
		for _, Entry := range UnBondingEntries {
			Unbondingtokens.Add(Entry.Balance)
		}

		newEntry := cosmosTypes.ValidatorDetails{
			ValidatorAddress: valAddr,
			BondedTokens:     BondedTokens,
			UnbondingTokens:  types.NewCoin(constants.CosmosDenom, Unbondingtokens),
		}
		ValidatorDetailsArr = append(ValidatorDetailsArr, newEntry)
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
