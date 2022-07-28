package orchestrator

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	prov "github.com/tendermint/tendermint/light/provider/http"
	"google.golang.org/grpc"
)

func TestZ(t *testing.T) {

	custodialAdrr, err := AccAddressFromBech32("cosmos15ddw7dkp56zytf3peshxr8fwn5w76y4g462ql2", "cosmos")
	if err != nil {
		panic(any(err))
	}
	rpcClientC, _ := newRPCClient("http://13.212.166.231:26657", 1*time.Second)
	liteproviderC, _ := prov.New("test", "http://13.212.166.231:26657")
	chainC := &CosmosChain{
		Key:              "unusedKey",
		ChainID:          "test",
		CustodialAddress: custodialAdrr,
		RPCAddr:          "http://13.212.166.231:26657",
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

	clientContextCosmos := client.Context{}.
		WithCodec(cosmosEncodingConfig.Marshaler).
		WithInterfaceRegistry(cosmosEncodingConfig.InterfaceRegistry).
		WithTxConfig(cosmosEncodingConfig.TxConfig).
		WithLegacyAmino(cosmosEncodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authTypes.AccountRetriever{}).
		WithNodeURI(chainC.RPCAddr).
		WithClient(chainC.Client).
		WithHomeDir("./").
		WithViper("").
		WithChainID(chainC.ChainID)

	fmt.Println(clientContextCosmos)

	grpcConnCosmos, _ := grpc.Dial(chainC.GRPCAddr, grpc.WithInsecure())
	defer func(grpcConnCosmos *grpc.ClientConn) {
		err := grpcConnCosmos.Close()
		if err != nil {
			panic(any(err))
		}
	}(grpcConnCosmos)

	txClient := sdkTx.NewServiceClient(grpcConnCosmos)

	fmt.Println("service client created")

	var status string
	TxHash := "2585F0C3CA6AB88DB0CB635F6C5AE4759C4809AD23533D92815CF21123B23320"

loop:
	for timeout := time.After(20 * time.Second); ; {

		select {
		case <-timeout:
			status = "not success"
			break loop
		default:
		}

		res2, err := txClient.GetTx(context.Background(),
			&sdkTx.GetTxRequest{
				Hash: TxHash,
			},
		)
		if err != nil {
			errorS := err.Error()
			ok := strings.Contains(errorS, "not found")
			if ok {
				continue
			} else {
				status = "not success"
				fmt.Println(status)
			}

		}

		txCode := res2.TxResponse.Code

		if txCode == sdkErrors.SuccessABCICode {
			status = "success"
			fmt.Println(status)
			break loop
		} else if txCode == sdkErrors.ErrInvalidSequence.ABCICode() {
			status = "sequence mismatch"
			fmt.Println(status)
			break loop
		} else if txCode == sdkErrors.ErrOutOfGas.ABCICode() {
			status = "gas failure"
			fmt.Println(status)
			break
		} else {
			status = "not success"
			fmt.Println(status)
			break loop
		}
	}

	fmt.Println(status)

}
