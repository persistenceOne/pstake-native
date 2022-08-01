package orchestrator

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	prov "github.com/tendermint/tendermint/light/provider/http"
)

func TestE2EAck(t *testing.T) {
	rpcaddr := "http://localhost:36657"
	grpcaddr := "localhost:9010"
	seed := "april patch recipe debate remove hurdle concert gesture design near predict enough color tail business imitate twelve february punch cheap vanish december cool wheel"
	_, addr := GetPivKeyAddress("persistence", 118, seed)

	rpcClient, _ := newRPCClient(rpcaddr, 1*time.Second)
	liteprovider, _ := prov.New("pstaked", rpcaddr)
	chain := &NativeChain{
		Key:           "unusedNativeKey",
		ChainID:       "pstaked",
		RPCAddr:       rpcaddr,
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

	cosmosrpc := "http://localhost:12003"
	cosmosgrpc := "localhost:12344"
	rpcClientC, _ := newRPCClient(cosmosrpc, 1*time.Second)
	liteproviderC, _ := prov.New("test", cosmosrpc)

	cusTodialAddress, _ := AccAddressFromBech32("cosmos15ddw7dkp56zytf3peshxr8fwn5w76y4g462ql2", "cosmos")
	chainC := &CosmosChain{
		Key:              "unusedKey",
		CustodialAddress: cusTodialAddress,
		ChainID:          "gaiad",
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

	rpcClientC, _ = newRPCClient("http://localhost:12003", 1*time.Second)
	height := int64(3657)

	blockResults, err := rpcClientC.BlockResults(context.Background(), &height)

	blockRes := blockResults.TxsResults

	TxhAsh := "3871A626178FDAE0863581E491A47437FBC1592F8783D83624D298E07CDCCB54"

	err = SendMsgAck(chain, chainC, []string{seed}, TxhAsh, "success", clientContextNative, clientContextCosmos, blockRes)
	if err != nil {
		fmt.Println(err, "ddddd")
	}

}
