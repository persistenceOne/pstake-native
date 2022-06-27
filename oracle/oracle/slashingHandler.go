package oracle

import (
	"context"
	"fmt"
	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txD "github.com/cosmos/cosmos-sdk/types/tx"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	"google.golang.org/grpc"
	logg "log"
)

func (c *CosmosChain) SlashingHandler(slash string, orcSeeds []string, valAddr string, nativeCliCtx cosmosClient.Context, native *NativeChain, chain *CosmosChain, cHeight int64) error {
	custodialAddr := chain.CustodialAddress.String()
	slashedValAddress, _ := sdk.AccAddressFromHex(slash)

	grpcConn, err := grpc.Dial(chain.GRPCAddr, grpc.WithInsecure())
	defer func(grpcConn *grpc.ClientConn) {
		err := grpcConn.Close()
		if err != nil {
			logg.Println("GRPC Connection error")
		}
	}(grpcConn)

	if err != nil {
		logg.Println("GRPC Connection failed")
		panic(err)
	}

	stakingQueryClient := stakingTypes.NewQueryClient(grpcConn)

	fmt.Println("staking query client connected")

	BondedTokensQueryResult, err := stakingQueryClient.Delegation(context.Background(),
		&stakingTypes.QueryDelegationRequest{
			DelegatorAddr: custodialAddr,
			ValidatorAddr: string(slashedValAddress),
		},
	)

	BondedDelegations := BondedTokensQueryResult.DelegationResponse.Balance
	//valAddr := BondedTokensQueryResult.DelegationResponse.Delegation.ValidatorAddress

	_, addr := GetSDKPivKeyAndAddress(orcSeeds[0])
	msg := &cosmosTypes.MsgSlashingEventOnCosmosChain{
		ValidatorAddress:    valAddr,
		CurrentDelegation:   BondedDelegations,
		OrchestratorAddress: string(addr),
		SlashType:           "",
		ChainID:             chain.ChainID,
		BlockHeight:         cHeight,
	}

	txBytes, err := SignNativeTx(orcSeeds[0], native, nativeCliCtx, msg)

	if err != nil {
		return err
	}

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
