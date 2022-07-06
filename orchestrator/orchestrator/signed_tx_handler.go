package orchestrator

import (
	"context"
	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	authTx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"google.golang.org/grpc"
	stdlog "log"
	"strconv"
	"strings"
	"time"
)

const (
	SUCCESS           = "success"
	NOT_SUCCESS       = "not success"
	GAS_FAILURE       = "gas failure"
	SEQUENCE_MISMATCH = "sequence mismatch"
)

func (n *NativeChain) SignedOutgoingTxHandler(txIdStr, valAddr string, orcSeeds []string, nativeCliCtx cosmosClient.Context, clientCtx cosmosClient.Context, native *NativeChain, chain *CosmosChain) error {
	txId, err := strconv.ParseUint(txIdStr, 10, 64)

	if err != nil {
		return err
	}
	grpcConn, err := grpc.Dial(native.GRPCAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer func(grpcConn *grpc.ClientConn) {
		err := grpcConn.Close()
		if err != nil {
			stdlog.Println("GRPC Connection error")
		}
	}(grpcConn)
	LiquidStakingModuleClient := cosmosTypes.NewQueryClient(grpcConn)

	stdlog.Println("query client connected")

	TxResult, err := LiquidStakingModuleClient.QueryTxByID(context.Background(),
		&cosmosTypes.QueryOutgoingTxByIDRequest{TxID: uint64(txId)},
	)
	if err != nil {
		return err
	}

	SignedTx := TxResult.CosmosTxDetails.Tx
	sigTx := authTx.WrapTx(&SignedTx)

	sigTx1 := sigTx.GetTx()
	signedTxBytes, err := clientCtx.TxConfig.TxEncoder()(sigTx1)
	if err != nil {
		panic(err)
	}

	if err != nil {
		return err
	}
	grpcConnCosmos, err := grpc.Dial(chain.GRPCAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer func(grpcConnCosmos *grpc.ClientConn) {
		err := grpcConnCosmos.Close()
		if err != nil {

		}
	}(grpcConnCosmos)

	txClient := sdkTx.NewServiceClient(grpcConnCosmos)

	stdlog.Println("client created")
	res, err := txClient.BroadcastTx(context.Background(),
		&sdkTx.BroadcastTxRequest{
			Mode:    sdkTx.BroadcastMode_BROADCAST_MODE_SYNC,
			TxBytes: signedTxBytes,
		},
	)
	if err != nil {
		return err
	}
	stdlog.Println(res.TxResponse.Code, res.TxResponse.TxHash, res)

	var status string
	var height int64
	var blockRes *coretypes.ResultBlockResults

	cosmosTxHash := res.TxResponse.TxHash

loop:
	for timeout := time.After(20 * time.Second); ; {

		select {
		case <-timeout:
			status = NOT_SUCCESS
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
				status = NOT_SUCCESS
			}

		}

		txCode := res2.TxResponse.Code

		if txCode == sdkErrors.SuccessABCICode {
			status = SUCCESS
			height = res2.TxResponse.Height
			break loop
		} else if txCode == sdkErrors.ErrInvalidSequence.ABCICode() {
			status = SEQUENCE_MISMATCH
			break loop
		} else if txCode == sdkErrors.ErrOutOfGas.ABCICode() {
			status = GAS_FAILURE
			break
		} else {
			status = NOT_SUCCESS

			break loop
		}
	}

	if status == SUCCESS {
		rpcClient, err := newRPCClient(chain.RPCAddr, 1*time.Second)
		if err != nil {
			return err
		}
		blockRes, err = rpcClient.BlockResults(context.Background(), &height)
		if err != nil {
			return err
		}
	}

	err = SendMsgAck(native, chain, orcSeeds, cosmosTxHash, status, nativeCliCtx, clientCtx, blockRes.TxsResults)
	if err != nil {
		return err
	}

	return nil

}
