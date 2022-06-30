package oracle

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txD "github.com/cosmos/cosmos-sdk/types/tx"
	tx2 "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	prov "github.com/tendermint/tendermint/light/provider/http"
	"google.golang.org/grpc"
	logg "log"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestC(t *testing.T) {

	seed := "april patch recipe debate remove hurdle concert gesture design near predict enough color tail business imitate twelve february punch cheap vanish december cool wheel"
	_, addr := GetSDKPivKeyAndAddressR("persistence", 118, seed)

	rpcClient, _ := newRPCClient("http://13.229.64.99:26657", 1*time.Second)
	liteprovider, _ := prov.New("native", "http://13.229.64.99:26657")
	chain := &NativeChain{
		Key:           "unusedNativeKey",
		ChainID:       "native",
		RPCAddr:       "http://13.229.64.99:26657",
		AccountPrefix: "persistence",
		GasAdjustment: 1.0,
		GasPrices:     "0.025stake",
		GRPCAddr:      "13.229.64.99:9090",
		CoinType:      118,
		HomePath:      "",
		KeyBase:       nil,
		Client:        rpcClient,
		Encoding:      params.EncodingConfig{},
		Provider:      liteprovider,
		address:       nil,
		logger:        nil,
		timeout:       0,
		debug:         false,
	}

	nativeEncodingConfig := chain.MakeEncodingConfig()
	chain.Encoding = nativeEncodingConfig
	chain.logger = defaultChainLogger()

	clientContextNative := client.Context{}.
		WithFromAddress(sdk.AccAddress(addr)).
		WithCodec(nativeEncodingConfig.Marshaler).
		WithInterfaceRegistry(nativeEncodingConfig.InterfaceRegistry).
		WithTxConfig(nativeEncodingConfig.TxConfig).
		WithLegacyAmino(nativeEncodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authTypes.AccountRetriever{}).
		WithNodeURI(chain.RPCAddr).
		WithClient(chain.Client).
		WithHomeDir("./").
		WithViper("").
		WithChainID(chain.ChainID)

	rpcClientC, _ := newRPCClient("http://13.212.166.231:26657", 1*time.Second)
	liteproviderC, _ := prov.New("test", "http://13.212.166.231:26657")
	chainC := &CosmosChain{
		Key:           "unusedKey",
		ChainID:       "test",
		RPCAddr:       "http://13.212.166.231:26657",
		AccountPrefix: "cosmos",
		GasAdjustment: 1.0,
		GasPrices:     "0.025stake",
		GRPCAddr:      "13.212.166.231:9090",
		CoinType:      118,
		HomePath:      "",
		KeyBase:       nil,
		Client:        rpcClientC,
		Encoding:      params.EncodingConfig{},
		Provider:      liteproviderC,
		address:       nil,
		logger:        nil,
		timeout:       0,
		debug:         false,
	}

	cosmosEncodingConfig := chainC.MakeEncodingConfig()
	chain.Encoding = cosmosEncodingConfig
	chain.logger = defaultChainLogger()

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

	txId, err := strconv.ParseUint("1", 10, 64)

	if err != nil {
		panic(err)
	}

	grpcConn, err := grpc.Dial(chain.GRPCAddr, grpc.WithInsecure())
	defer func(grpcConn *grpc.ClientConn) {
		err := grpcConn.Close()
		if err != nil {
			logg.Println("GRPC Connection error")
		}
	}(grpcConn)
	LiquidStakingModuleClient := cosmosTypes.NewQueryClient(grpcConn)

	fmt.Println("query client connected")

	TxResult, err := LiquidStakingModuleClient.QueryTxByID(context.Background(),
		&cosmosTypes.QueryOutgoingTxByIDRequest{TxID: uint64(txId)},
	)

	SignedTx := TxResult.CosmosTxDetails.Tx
	sigTx := tx2.WrapTx(&SignedTx)

	sigTx1 := sigTx.GetTx()

	signedTxBytes, err := clientContextCosmos.TxConfig.TxEncoder()(sigTx1)
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}
	grpcConnCosmos, _ := grpc.Dial(chainC.GRPCAddr, grpc.WithInsecure())
	defer func(grpcConnCosmos *grpc.ClientConn) {
		err := grpcConnCosmos.Close()
		if err != nil {

		}
	}(grpcConnCosmos)

	txClient := txD.NewServiceClient(grpcConnCosmos)

	fmt.Println("client created")
	res, err := txClient.BroadcastTx(context.Background(),
		&txD.BroadcastTxRequest{
			Mode:    txD.BroadcastMode_BROADCAST_MODE_SYNC,
			TxBytes: signedTxBytes,
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.TxResponse.Code, res.TxResponse.TxHash, res)
	status := "success"
	cosmosTxHash := res.TxResponse.TxHash
	if res.TxResponse.Code == 0 {
		status = "keeper failure"
	}
	err = SendMsgAcknowledgement(chain, chainC, []string{seed}, cosmosTxHash, status, clientContextNative, clientContextNative)
	if err != nil {
		panic(err)
	}

}
