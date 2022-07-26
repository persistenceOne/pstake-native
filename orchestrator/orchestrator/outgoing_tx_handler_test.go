package orchestrator

import (
	"context"
	"fmt"
	stdlog "log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	prov "github.com/tendermint/tendermint/light/provider/http"
	"google.golang.org/grpc"
)

func TestB(t *testing.T) {
	rpcaddr := "http://18.139.224.127:26657"
	grpcaddr := "18.139.224.127:9090"
	seed := "bomb sand fashion torch return coconut color captain vapor inhale lyrics lady grant ordinary lazy decrease quit devote paddle impulse prize equip hip ball"
	_, addr := GetPivKeyAddress("persistence", 118, seed)

	rpcClient, _ := newRPCClient(rpcaddr, 1*time.Second)
	liteprovider, _ := prov.New("native", rpcaddr)
	chain := &NativeChain{
		Key:           "unusedNativeKey",
		ChainID:       "native",
		RPCAddr:       grpcaddr,
		AccountPrefix: "persistence",
		GasAdjustment: 1.0,
		GasPrices:     "0.025stake",
		GRPCAddr:      grpcaddr,
		CoinType:      118,
		HomePath:      "",
		KeyBase:       nil,
		Client:        rpcClient,
		Encoding:      params.EncodingConfig{},
		Provider:      liteprovider,
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

	custodialAdrr, err := AccAddressFromBech32("cosmos15ddw7dkp56zytf3peshxr8fwn5w76y4g462ql2", "cosmos")
	cosmosrpc := "http://13.212.166.231:26657"
	cosmosgrpc := "13.212.166.231:9090"
	rpcClientC, _ := newRPCClient(cosmosrpc, 1*time.Second)
	liteproviderC, _ := prov.New("test", cosmosrpc)
	chainC := &CosmosChain{
		Key:              "unusedNativeKey",
		ChainID:          "test",
		RPCAddr:          cosmosrpc,
		AccountPrefix:    "cosmos",
		CustodialAddress: custodialAdrr,
		GasAdjustment:    1.0,
		GasPrices:        "0.025stake",
		GRPCAddr:         cosmosgrpc,
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
		panic(any(err))
	}

	grpcConn, err := grpc.Dial(grpcaddr, grpc.WithInsecure())
	defer func(grpcConn *grpc.ClientConn) {
		err := grpcConn.Close()
		if err != nil {
			stdlog.Println("GRPC Connection error")
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
	_, addr = GetPivKeyAddress("persistence", 118, seed)

	if err != nil {
		panic(any(err))
	}

	grpcConnCos, _ := grpc.Dial(chain.GRPCAddr, grpc.WithInsecure())
	defer func(grpcConnCos *grpc.ClientConn) {
		err := grpcConnCos.Close()
		if err != nil {

		}
	}(grpcConnCos)

	txClient := sdkTx.NewServiceClient(grpcConnCos)

	fmt.Println("client created")

	msg := &cosmosTypes.MsgSetSignature{
		OrchestratorAddress: addr,
		OutgoingTxID:        txId,
		Signature:           signature,
		BlockHeight:         int64(120),
	}

	txBytes, err := SignNativeTx(seed, chain, clientContextNative, msg)

	res, err := txClient.BroadcastTx(context.Background(),
		&sdkTx.BroadcastTxRequest{
			Mode:    sdkTx.BroadcastMode_BROADCAST_MODE_BLOCK,
			TxBytes: txBytes,
		},
	)
	if err != nil {
		panic(any(err))
	}
	fmt.Println(res.TxResponse.Code, res.TxResponse.TxHash, res)

	//err = SendMsgAck(native, chain, orcSeeds, res.TxResponse.TxHash, valAddr, nativeCliCtx, clientCtx)

	if err != nil {
		panic(any(err))
	}

}
