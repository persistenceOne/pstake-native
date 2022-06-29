package oracle

import (
	"context"
	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	txD "github.com/cosmos/cosmos-sdk/types/tx"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	"google.golang.org/grpc"
	logg "log"
)

func SendMsgAcknowledgement(native *NativeChain, cosmosChain *CosmosChain, orcSeeds []string, TxHash string, status string, nativeCliCtx cosmosClient.Context, clientCtx cosmosClient.Context) error {

	_, addr := GetSDKPivKeyAndAddressR(native.AccountPrefix, native.CoinType, orcSeeds[0])

	address, err := AccAddressFromBech32(addr, native.AccountPrefix)
	if err != nil {
		return err
	}
	ValDetails := GetValidatorDetails(cosmosChain)
	acc, seq, err := clientCtx.AccountRetriever.GetAccountNumberSequence(clientCtx, address)

	if err != nil {
		return err
	}

	msg := &cosmosTypes.MsgTxStatus{
		OrchestratorAddress: addr,
		TxHash:              TxHash,
		Status:              status,
		SequenceNumber:      seq,
		AccountNumber:       acc,
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

	logg.Println(res.TxResponse.Code, res.TxResponse.TxHash, res)

	if err != nil {
		return err
	}

	return nil

}
