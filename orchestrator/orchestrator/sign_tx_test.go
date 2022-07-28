package orchestrator

import (
	"context"
	"fmt"
	stdlog "log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	authTx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	prov "github.com/tendermint/tendermint/light/provider/http"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"google.golang.org/grpc"
)

func TestE2ESignedTx(t *testing.T) {

	seed := []string{"bomb sand fashion torch return coconut color captain vapor inhale lyrics lady grant ordinary lazy decrease quit devote paddle impulse prize equip hip ball",
		"april patch recipe debate remove hurdle concert gesture design near predict enough color tail business imitate twelve february punch cheap vanish december cool wheel",
		"soft brown armed regret hip few ancient control steel bright basic swamp sentence present immune napkin orbit giggle year another crowd essence noble dice",
		"road gallery tooth script volcano deputy summer acid bulk anger fatigue notable secret blood bean apology burger rookie rug bench away dutch secret upper"}
	_, addr := GetPivKeyAddress("persistence", 118, seed[2])

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

	txId, err := strconv.ParseUint("2", 10, 64)

	if err != nil {
		panic(err)
	}

	grpcConn, err := grpc.Dial(chain.GRPCAddr, grpc.WithInsecure())
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

	SignedTx := TxResult.CosmosTxDetails.Tx
	sigTx := authTx.WrapTx(&SignedTx)

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
			panic(err)
		}
	}(grpcConnCosmos)

	txClient := sdkTx.NewServiceClient(grpcConnCosmos)

	fmt.Println("service client created")
	res, err := txClient.BroadcastTx(context.Background(),
		&sdkTx.BroadcastTxRequest{
			Mode:    sdkTx.BroadcastMode_BROADCAST_MODE_SYNC,
			TxBytes: signedTxBytes,
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.TxResponse.Code, res.TxResponse.TxHash, res)
	var status string
	var height int64
	var BlockResults *coretypes.ResultBlockResults
	cosmosTxHash := res.TxResponse.TxHash

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
				Hash: cosmosTxHash,
			},
		)
		if err != nil {
			errorS := err.Error()
			ok := strings.Contains(errorS, "not found")
			if ok {
				continue loop
			} else {
				status = "not success"
			}

		}

		txCode := res2.TxResponse.Code

		if txCode == sdkErrors.SuccessABCICode {
			status = "success"
			height = res2.TxResponse.Height
			break loop
		} else if txCode == sdkErrors.ErrInvalidSequence.ABCICode() {
			status = "sequence mismatch"
			break loop
		} else if txCode == sdkErrors.ErrOutOfGas.ABCICode() {
			status = "gas failure"
			break
		} else {
			status = "not success"

			break loop
		}

	}

	if status == "success" && height != 0 {
		BlockResults, _ = rpcClient.BlockResults(context.Background(), &height)
	}

	err = SendMsgAck(chain, chainC, []string{seed[2]}, cosmosTxHash, status, clientContextNative, clientContextCosmos, BlockResults.TxsResults)
	if err != nil {
		panic(err)
	}

}
