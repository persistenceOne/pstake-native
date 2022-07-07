package orchestrator

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/persistenceOne/pstake-native/orchestrator/constants"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	prov "github.com/tendermint/tendermint/light/provider/http"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	stdlog "log"
	"testing"
	"time"
)

func TestE2EValDetails(t *testing.T) {

	custodialAdrr, err := AccAddressFromBech32("cosmos15ddw7dkp56zytf3peshxr8fwn5w76y4g462ql2", "cosmos")
	if err != nil {
		return
	}
	rpcClientC, _ := newRPCClient("http://13.212.166.231:26657", 1*time.Second)
	liteproviderC, _ := prov.New("native", "http://13.212.166.231:26657")
	chainC := &CosmosChain{
		Key:              "unusedNativeKey",
		ChainID:          "test",
		RPCAddr:          "http://13.212.166.231:26657",
		CustodialAddress: custodialAdrr,
		AccountPrefix:    "cosmos",
		GasAdjustment:    1.0,
		GasPrices:        "0.025stake",
		GRPCAddr:         "13.212.166.231:9090",
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

	//clientContextCosmos := client.Context{}.
	//	WithCodec(cosmosEncodingConfig.Marshaler).
	//	WithInterfaceRegistry(cosmosEncodingConfig.InterfaceRegistry).
	//	WithTxConfig(cosmosEncodingConfig.TxConfig).
	//	WithLegacyAmino(cosmosEncodingConfig.Amino).
	//	WithInput(os.Stdin).
	//	WithAccountRetriever(authTypes.AccountRetriever{}).
	//	WithNodeURI(chainC.RPCAddr).
	//	WithClient(chainC.Client).
	//	WithHomeDir("./").
	//	WithViper("").
	//	WithChainID(chainC.ChainID)

	var ValidaDetailsArr []cosmosTypes.ValidatorDetails

	custodialAddr, err := Bech32ifyAddressBytes(chainC.AccountPrefix, chainC.CustodialAddress)

	grpcConn, err := grpc.Dial(chainC.GRPCAddr, grpc.WithInsecure())
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

	BondedTokensQueryResult, err := stakingQueryClient.DelegatorDelegations(context.Background(),
		&stakingTypes.QueryDelegatorDelegationsRequest{
			DelegatorAddr: custodialAddr,
			Pagination:    nil,
		},
	)

	if err != nil {
		stdlog.Println("cannot get total delegations")
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

		Unbondingtokens := sdk.NewInt(0)

		if flag == true {
			UnBondingEntries := UnbondingTokensQueryResult.Unbond.Entries
			for _, Entry := range UnBondingEntries {
				Unbondingtokens.Add(Entry.Balance)
			}
		}

		newEntry := cosmosTypes.ValidatorDetails{
			ValidatorAddress: valAddr,
			BondedTokens:     BondedTokens,
			UnbondingTokens:  sdk.NewCoin(constants.CosmosDenom, Unbondingtokens),
		}
		ValidaDetailsArr = append(ValidaDetailsArr, newEntry)
		flag = true

	}
	fmt.Println(ValidaDetailsArr)
}
