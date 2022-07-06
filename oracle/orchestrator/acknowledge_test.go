package orchestrator

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	prov "github.com/tendermint/tendermint/light/provider/http"
	"google.golang.org/grpc"
	stdlog "log"
	"os"
	"testing"
	"time"
)

func TestM(t *testing.T) {
	rpcaddr := "http://13.229.64.99:26657"
	grpcaddr := "13.229.64.99:9090"
	seed := "april patch recipe debate remove hurdle concert gesture design near predict enough color tail business imitate twelve february punch cheap vanish december cool wheel"
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

	cosmosrpc := "http://13.212.166.231:26657"
	cosmosgrpc := "13.212.166.231:9090"
	rpcClientC, _ := newRPCClient(cosmosrpc, 1*time.Second)
	liteproviderC, _ := prov.New("test", cosmosrpc)

	cusTodialAddress, _ := AccAddressFromBech32("cosmos15ddw7dkp56zytf3peshxr8fwn5w76y4g462ql2", "cosmos")
	chainC := &CosmosChain{
		Key:              "unusedNativeKey",
		CustodialAddress: cusTodialAddress,
		ChainID:          "test",
		RPCAddr:          cosmosrpc,
		AccountPrefix:    "cosmos",
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

	ValDetails := GetValidatorDetails(chainC)

	SetSDKConfigPrefix(chainC.AccountPrefix)
	address, err, flag := GetMultiSigAddress(chain, chainC)
	if err != nil {
		panic(err)
	}

	if flag == "pass" {
		acc, seq, err := clientContextCosmos.AccountRetriever.GetAccountNumberSequence(clientContextCosmos, address)

		if err != nil {
			panic(err)
		}

		msg := &cosmosTypes.MsgTxStatus{
			OrchestratorAddress: addr,
			TxHash:              "2309DBA86D984925C45DE6C3A697E114303C948556E38EDBE1ED338CEC878A06",
			Status:              "success",
			SequenceNumber:      seq,
			AccountNumber:       acc,
			ValidatorDetails:    ValDetails,
		}

		txBytes, err := SignNativeTx(seed, chain, clientContextNative, msg)

		if err != nil {
			panic(err)
		}

		grpcConn, _ := grpc.Dial(chain.GRPCAddr, grpc.WithInsecure())
		defer func(grpcConn *grpc.ClientConn) {
			err := grpcConn.Close()
			if err != nil {

			}
		}(grpcConn)

		txClient := sdkTx.NewServiceClient(grpcConn)

		res, err := txClient.BroadcastTx(context.Background(),
			&sdkTx.BroadcastTxRequest{
				Mode:    sdkTx.BroadcastMode_BROADCAST_MODE_SYNC,
				TxBytes: txBytes,
			},
		)
		if err != nil {
			panic(err)

		}

		stdlog.Println(res.TxResponse.Code, res.TxResponse.TxHash, res)

		if err != nil {
			panic(err)
		}
	}

}
