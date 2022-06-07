package oracle

import (
	"context"
	"fmt"
	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	txD "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	"google.golang.org/grpc"
	logg "log"
	"strconv"
)

func (n *NativeChain) OutgoingTxHandler(txIdstr string, valAddr string, orcSeeds []string, nativeCliCtx cosmosClient.Context, clientCtx cosmosClient.Context, native *NativeChain, chain *CosmosChain, depositHeight int64, protoCodec *codec.ProtoCodec) error {
	txId, err := strconv.ParseUint(txIdstr, 10, 64)

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

	fmt.Println("query client connected")

	TxResult, err := LiquidStakingModuleClient.QueryTxByID(context.Background(),
		&cosmosTypes.QueryOutgoingTxByIDRequest{TxID: uint64(txId)},
	)

	OutgoingTx := TxResult.CosmosTxDetails.Tx
	TxBytes, err := clientCtx.TxConfig.TxEncoder()(&OutgoingTx)
	if err != nil {
		panic(err)
		return err
	}

	txInterface, err := clientCtx.TxConfig.TxDecoder()(TxBytes)
	if err != nil {
		panic(err)
		return err
	}

	tx, ok := txInterface.(signing.Tx)
	if !ok {
		return err
	}
	for _, msg := range tx.GetMsgs() {
		fmt.Println(msg.String())

		signedTxBytes, err := SignCosmosTx(orcSeeds[0], chain, nativeCliCtx, msg)
		if err != nil {
			return err
		}
		grpcConn, _ := grpc.Dial(chain.GRPCAddr, grpc.WithInsecure())
		defer func(grpcConn *grpc.ClientConn) {
			err := grpcConn.Close()
			if err != nil {

			}
		}(grpcConn)

		txClient := txD.NewServiceClient(grpcConn)

		fmt.Println("client created")

		res, err := txClient.BroadcastTx(context.Background(),
			&txD.BroadcastTxRequest{
				Mode:    txD.BroadcastMode_BROADCAST_MODE_SYNC,
				TxBytes: signedTxBytes,
			},
		)
		if err != nil {
			return err
		}
		fmt.Println(res.TxResponse.Code, res.TxResponse.TxHash, res)

		cosmosTxHash := res.TxResponse.TxHash

		err = SendMsgAcknowledgement(native, chain, orcSeeds, cosmosTxHash, valAddr, nativeCliCtx, clientCtx)
		if err != nil {
			return err
		}
	}
	return nil
}
