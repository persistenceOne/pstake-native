package oracle

import (
	"context"
	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	txD "github.com/cosmos/cosmos-sdk/types/tx"
	tx2 "github.com/cosmos/cosmos-sdk/x/auth/tx"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	"google.golang.org/grpc"
	logg "log"
	"strconv"
	"strings"
	"time"
)

func (n *NativeChain) SignedOutgoingTxHandler(txIdStr, valAddr string, orcSeeds []string, nativeCliCtx cosmosClient.Context, clientCtx cosmosClient.Context, native *NativeChain, chain *CosmosChain) error {
	txId, err := strconv.ParseUint(txIdStr, 10, 64)

	if err != nil {
		return err
	}
	grpcConn, err := grpc.Dial(native.GRPCAddr, grpc.WithInsecure())
	defer func(grpcConn *grpc.ClientConn) {
		err := grpcConn.Close()
		if err != nil {
			logg.Println("GRPC Connection error")
		}
	}(grpcConn)
	LiquidStakingModuleClient := cosmosTypes.NewQueryClient(grpcConn)

	logg.Println("query client connected")

	TxResult, err := LiquidStakingModuleClient.QueryTxByID(context.Background(),
		&cosmosTypes.QueryOutgoingTxByIDRequest{TxID: uint64(txId)},
	)

	SignedTx := TxResult.CosmosTxDetails.Tx
	sigTx := tx2.WrapTx(&SignedTx)

	sigTx1 := sigTx.GetTx()
	signedTxBytes, err := clientCtx.TxConfig.TxEncoder()(sigTx1)
	if err != nil {
		panic(err)
	}

	if err != nil {
		return err
	}
	grpcConnCosmos, _ := grpc.Dial(chain.GRPCAddr, grpc.WithInsecure())
	defer func(grpcConnCosmos *grpc.ClientConn) {
		err := grpcConnCosmos.Close()
		if err != nil {

		}
	}(grpcConnCosmos)

	txClient := txD.NewServiceClient(grpcConnCosmos)

	logg.Println("client created")
	res, err := txClient.BroadcastTx(context.Background(),
		&txD.BroadcastTxRequest{
			Mode:    txD.BroadcastMode_BROADCAST_MODE_SYNC,
			TxBytes: signedTxBytes,
		},
	)
	if err != nil {
		return err
	}
	logg.Println(res.TxResponse.Code, res.TxResponse.TxHash, res)
	var status string
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
			&txD.GetTxRequest{
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

	err = SendMsgAcknowledgement(native, chain, orcSeeds, cosmosTxHash, status, nativeCliCtx, clientCtx)
	if err != nil {
		return err
	}

	return nil

}
