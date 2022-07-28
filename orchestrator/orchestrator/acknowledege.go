package orchestrator

import (
	"context"
	stdlog "log"

	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	abciTypes "github.com/tendermint/tendermint/abci/types"
	"google.golang.org/grpc"
)

func SendMsgAck(native *NativeChain, cosmosChain *CosmosChain, orcSeeds []string, txHash string, status string,
	nativeCliCtx cosmosClient.Context, clientCtx cosmosClient.Context, blockResults []*abciTypes.ResponseDeliverTx) error {
	_, addr := GetPivKeyAddress(native.AccountPrefix, native.CoinType, orcSeeds[0])

	valDetails := GetValidatorDetails(cosmosChain)
	if status == SUCCESS {
		val, err := PopulateRewards(cosmosChain, valDetails, blockResults)
		if err != nil {
			return err
		}
		valDetails = val
	}

	SetSDKConfigPrefix(cosmosChain.ChainID)
	address, err, ok := GetMultiSigAddress(native, cosmosChain)
	if err != nil {
		return err
	}

	if ok {
		acc, seq, err := clientCtx.AccountRetriever.GetAccountNumberSequence(clientCtx, address)
		if err != nil {
			return err
		}

		msg := &cosmosTypes.MsgTxStatus{
			OrchestratorAddress: addr,
			TxHash:              txHash,
			Status:              status,
			SequenceNumber:      seq,
			AccountNumber:       acc,
			ValidatorDetails:    valDetails,
		}

		txBytes, err := SignNativeTx(orcSeeds[0], native, nativeCliCtx, msg)

		if err != nil {
			return err
		}

		grpcConn, err := grpc.Dial(native.GRPCAddr, grpc.WithInsecure())

		if err != nil {
			panic(any(err))
		}
		defer func(grpcConn *grpc.ClientConn) {
			err := grpcConn.Close()
			if err != nil {

			}
		}(grpcConn)

		txClient := sdkTx.NewServiceClient(grpcConn)

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

		if err != nil {
			return err
		}

		return nil
	}

	return nil

}

func GetMultiSigAddress(chain *NativeChain, chainC *CosmosChain) (types.AccAddress, error, bool) {
	var txId uint64

	grpcConn, err := grpc.Dial(chain.GRPCAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err, false
	}
	defer func(grpcConn *grpc.ClientConn) {
		err := grpcConn.Close()
		if err != nil {
			stdlog.Println("GRPC Connection error")
		}
	}(grpcConn)

	if err != nil {
		stdlog.Println("GRPC Connection failed")
		panic(any(err))
	}
	cosmosQueryClient := cosmosTypes.NewQueryClient(grpcConn)
	stdlog.Println("staking query client connected")

	ActiveTxID, err := cosmosQueryClient.ActiveTxn(context.Background(),
		&cosmosTypes.QueryActiveTxnRequest{},
	)
	if err != nil {
		return nil, err, false
	}

	txId = ActiveTxID.GetTxID()
	if txId != 0 {
		TxResult, err := cosmosQueryClient.QueryTxByID(context.Background(),
			&cosmosTypes.QueryOutgoingTxByIDRequest{TxID: uint64(txId)},
		)

		signerAddress := TxResult.CosmosTxDetails.SignerAddress
		SetSDKConfigPrefix(chainC.AccountPrefix)

		signerAddr, err := AccAddressFromBech32(signerAddress, chainC.AccountPrefix)
		if err != nil {
			return nil, err, false
		}

		return signerAddr, nil, true

	}

	return nil, nil, false
}
