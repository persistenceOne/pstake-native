package oracle

import (
	"context"
	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	txD "github.com/cosmos/cosmos-sdk/types/tx"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	"google.golang.org/grpc"
	logg "log"
	"strconv"
)

func (n *NativeChain) OutgoingTxHandler(txIdstr string, valAddr string, orcSeeds []string, nativeCliCtx cosmosClient.Context, clientCtx cosmosClient.Context, native *NativeChain, chain *CosmosChain, cHeight uint64) error {
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

	logg.Println("query client connected")

	TxResult, err := LiquidStakingModuleClient.QueryTxByID(context.Background(),
		&cosmosTypes.QueryOutgoingTxByIDRequest{TxID: uint64(txId)},
	)

	OutgoingTx := TxResult.CosmosTxDetails.GetTx()

	signerAddress := TxResult.CosmosTxDetails.SignerAddress

	signature, err := GetSignBytesForCosmos(orcSeeds[0], chain, clientCtx, OutgoingTx, signerAddress)
	_, addr := GetSDKPivKeyAndAddressR(native.AccountPrefix, native.CoinType, orcSeeds[0])

	if err != nil {
		return err
	}

	grpcConnCos, _ := grpc.Dial(native.GRPCAddr, grpc.WithInsecure())
	defer func(grpcConnCos *grpc.ClientConn) {
		err := grpcConnCos.Close()
		if err != nil {

		}
	}(grpcConnCos)

	txClient := txD.NewServiceClient(grpcConnCos)

	logg.Println("client created")

	msg := &cosmosTypes.MsgSetSignature{
		OrchestratorAddress: addr,
		OutgoingTxID:        txId,
		Signature:           signature,
		BlockHeight:         int64(cHeight),
	}

	txBytes, err := SignNativeTx(orcSeeds[0], native, nativeCliCtx, msg)

	res, err := txClient.BroadcastTx(context.Background(),
		&txD.BroadcastTxRequest{
			Mode:    txD.BroadcastMode_BROADCAST_MODE_SYNC,
			TxBytes: txBytes,
		},
	)
	if err != nil {
		return err
	}
	logg.Println(res.TxResponse.Code, res.TxResponse.TxHash, res)

	if err != nil {
		return err
	}

	return nil
}
