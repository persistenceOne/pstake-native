package oracle

import (
	"context"
	"fmt"
	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	txD "github.com/cosmos/cosmos-sdk/types/tx"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	"google.golang.org/grpc"
)

func SendMsgAcknowledgement(native *NativeChain, cosmosChain *CosmosChain, orcSeeds []string, TxHash string, valAddr string, nativeCliCtx cosmosClient.Context, clientCtx cosmosClient.Context) error {
	_, addr := GetSDKPivKeyAndAddress(orcSeeds[0])

	ValDetails := GetValidatorDetails(cosmosChain)

	msg := &cosmosTypes.MsgTxStatus{
		OrchestratorAddress: addr.String(),
		TxHash:              TxHash,
		Status:              "success",
		SequenceNumber:      1,
		AccountNumber:       1,
		ValidatorDetails:    ValDetails,
	}

	txBytes, err := SignNativeTx(orcSeeds[0], native, nativeCliCtx, msg)

	if err != nil {
		return err
	}

	grpcConn, _ := grpc.Dial(native.GRPCAddr, grpc.WithInsecure())
	defer func(grpcConn *grpc.ClientConn) {
		err := grpcConn.Close()
		if err != nil {

		}
	}(grpcConn)

	txClient := txD.NewServiceClient(grpcConn)

	res, err := txClient.BroadcastTx(context.Background(),
		&txD.BroadcastTxRequest{
			Mode:    txD.BroadcastMode_BROADCAST_MODE_SYNC,
			TxBytes: txBytes,
		},
	)
	if err != nil {
		return err
	}

	fmt.Println(res.TxResponse.Code, res.TxResponse.TxHash, res)

	return nil

}
