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
	"os"
	"testing"
	"time"
)

func Test1(t *testing.T) {
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
		CoinType:      750,
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

	seed := "bomb sand fashion torch return coconut color captain vapor inhale lyrics lady grant ordinary lazy decrease quit devote paddle impulse prize equip hip ball"
	msg := &cosmosTypes.MsgMintTokensForAccount{
		AddressFromMemo:     "persistence1uv0stzuxn5ar3enrkmnaqh8jnz6uflqy88ydwe",
		OrchestratorAddress: "persistence1kma3k7lzgg2rjtymda0kn7z3nvc677ca36wapd",
		Amount:              sdk.NewCoin("stake", sdk.NewInt(10)),
		TxHash:              "testt",
		BlockHeight:         int64(67),
		ChainID:             "test",
	}

	addr, _ := AccAddressFromBech32("persistence1kma3k7lzgg2rjtymda0kn7z3nvc677ca36wapd", "persistence")
	fmt.Println(msg)

	clientContextNative := client.Context{}.
		WithFromAddress(addr).
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
	//cfg := sdk.GetConfig()
	//print(cfg)
	txBytes, err := SignNativeTx(seed, chain, clientContextNative, msg)
	if err != nil {
		panic(err)
	}

	grpcConn, _ := grpc.Dial(chain.GRPCAddr, grpc.WithInsecure())
	defer func(grpcConn *grpc.ClientConn) {
		err := grpcConn.Close()
		if err != nil {
			panic(err)
		}
	}(grpcConn)

	txClient := txD.NewServiceClient(grpcConn)
	fmt.Println("client created")

	res, err := txClient.BroadcastTx(context.Background(),
		&txD.BroadcastTxRequest{
			TxBytes: txBytes,
			Mode:    txD.BroadcastMode_BROADCAST_MODE_SYNC,
		})

	if err != nil {
		panic(err)
	}

	fmt.Println(res.TxResponse.Code, res.TxResponse)

}
