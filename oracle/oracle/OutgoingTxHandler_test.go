package oracle

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txD "github.com/cosmos/cosmos-sdk/types/tx"
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

func TestB(t *testing.T) {

	seed := "marble allow december print trial know resource cry next segment twice nose because steel omit confirm hair extend shrimp seminar one minor phone deputy"
	_, addr := GetSDKPivKeyAndAddressR("persistence", 118, seed)

	rpcClient, _ := newRPCClient("http://10.128.36.249:26657", 1*time.Second)
	liteprovider, _ := prov.New("native", "http://10.128.36.249:26657")
	chain := &NativeChain{
		Key:           "unusedNativeKey",
		ChainID:       "test",
		RPCAddr:       "http://10.128.36.249:26657",
		AccountPrefix: "persistence",
		GasAdjustment: 1.0,
		GasPrices:     "0.025stake",
		GRPCAddr:      "10.128.36.249:9090",
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
	liteproviderC, _ := prov.New("native", "http://13.212.166.231:26657")
	chainC := &CosmosChain{
		Key:           "unusedNativeKey",
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

	grpcConn, err := grpc.Dial("10.128.36.249:9090", grpc.WithInsecure())
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

	//ac,seq,err := clientCtx.AccountRetriever.GetAccount()
	OutgoingTx := TxResult.CosmosTxDetails.GetTx()

	fmt.Println(OutgoingTx.Body.Messages)

	signerAddress := TxResult.CosmosTxDetails.SignerAddress

	signature, err := GetSignBytesForCosmos(seed, chainC, clientContextCosmos, OutgoingTx, signerAddress)
	_, addr = GetSDKPivKeyAndAddressR("persistence", 118, seed)

	if err != nil {
		panic(err)
	}

	grpcConnCos, _ := grpc.Dial(chain.GRPCAddr, grpc.WithInsecure())
	defer func(grpcConnCos *grpc.ClientConn) {
		err := grpcConnCos.Close()
		if err != nil {

		}
	}(grpcConnCos)

	txClient := txD.NewServiceClient(grpcConnCos)

	fmt.Println("client created")

	msg := &cosmosTypes.MsgSetSignature{
		OrchestratorAddress: addr,
		OutgoingTxID:        txId,
		Signature:           signature,
		BlockHeight:         int64(120),
	}

	txBytes, err := SignNativeTx(seed, chain, clientContextNative, msg)

	res, err := txClient.BroadcastTx(context.Background(),
		&txD.BroadcastTxRequest{
			Mode:    txD.BroadcastMode_BROADCAST_MODE_BLOCK,
			TxBytes: txBytes,
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.TxResponse.Code, res.TxResponse.TxHash, res)

	//err = SendMsgAcknowledgement(native, chain, orcSeeds, res.TxResponse.TxHash, valAddr, nativeCliCtx, clientCtx)

	if err != nil {
		panic(err)
	}

}
