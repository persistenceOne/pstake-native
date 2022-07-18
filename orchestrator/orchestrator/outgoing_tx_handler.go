package orchestrator

import (
	"context"
	stdlog "log"
	"strconv"

	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	"google.golang.org/grpc"
)

func (n *NativeChain) OutgoingTxHandler(txIdstr string, valAddr string, orcSeeds []string, nativeCliCtx cosmosClient.Context, clientCtx cosmosClient.Context, native *NativeChain, chain *CosmosChain, cHeight uint64) error {
	txId, err := strconv.ParseUint(txIdstr, 10, 64)

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

	OutgoingTx := TxResult.CosmosTxDetails.GetTx()

	signerAddress := TxResult.CosmosTxDetails.SignerAddress

	signature, err := GetSignBytesForCosmos(orcSeeds[0], chain, clientCtx, OutgoingTx, signerAddress)
	if err != nil {
		return err
	}
	_, addr := GetPivKeyAddress(native.AccountPrefix, native.CoinType, orcSeeds[0])

	grpcConnCos, err := grpc.Dial(native.GRPCAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer func(grpcConnCos *grpc.ClientConn) {
		err := grpcConnCos.Close()
		if err != nil {

		}
	}(grpcConnCos)

	txClient := sdkTx.NewServiceClient(grpcConnCos)

	stdlog.Println("client created")

	msg := &cosmosTypes.MsgSetSignature{
		OrchestratorAddress: addr,
		OutgoingTxID:        txId,
		Signature:           signature,
		BlockHeight:         int64(cHeight),
	}

	txBytes, err := SignNativeTx(orcSeeds[0], native, nativeCliCtx, msg)
	if err != nil {
		return err
	}

	res, err := txClient.BroadcastTx(context.Background(),
		&sdkTx.BroadcastTxRequest{
			Mode:    sdkTx.BroadcastMode_BROADCAST_MODE_SYNC,
			TxBytes: txBytes,
		},
	)
	if err != nil {
		return err
	}
	stdlog.Println(res.TxResponse.Code, res.TxResponse.TxHash, res)
	return nil
}
