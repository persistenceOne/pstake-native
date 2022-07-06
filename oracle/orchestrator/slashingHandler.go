package orchestrator

import (
	"context"
	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	"google.golang.org/grpc"
	stdlog "log"
)

func (c *CosmosChain) SlashingHandler(slash string, orcSeeds []string, valAddr string, nativeCliCtx cosmosClient.Context, native *NativeChain, chain *CosmosChain, cHeight int64) error {
	custodialAddr, err := Bech32ifyAddressBytes(chain.AccountPrefix, chain.CustodialAddress)
	if err != nil {
		stdlog.Println(err)
		return err
	}
	SetSDKConfigPrefix(chain.AccountPrefix)
	slashedValAddress, err := AccAddressFromBech32(slash, chain.AccountPrefix)

	if err != nil {
		stdlog.Println(err)
		return err
	}

	grpcConn, err := grpc.Dial(chain.GRPCAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer func(grpcConn *grpc.ClientConn) {
		err := grpcConn.Close()
		if err != nil {
			stdlog.Println("GRPC Connection error")
		}
	}(grpcConn)

	if err != nil {
		stdlog.Println("GRPC Connection failed")
		return err
	}

	stakingQueryClient := stakingTypes.NewQueryClient(grpcConn)

	stdlog.Println("staking query client connected")

	BondedTokensQueryResult, err := stakingQueryClient.Delegation(context.Background(),
		&stakingTypes.QueryDelegationRequest{
			DelegatorAddr: custodialAddr,
			ValidatorAddr: string(slashedValAddress),
		},
	)

	if err != nil {
		return err
	}

	BondedDelegations := BondedTokensQueryResult.DelegationResponse.Balance

	_, addr := GetPivKeyAddress(native.AccountPrefix, native.CoinType, orcSeeds[0])
	msg := &cosmosTypes.MsgSlashingEventOnCosmosChain{
		ValidatorAddress:    valAddr,
		CurrentDelegation:   BondedDelegations,
		OrchestratorAddress: addr,
		SlashType:           "",
		ChainID:             chain.ChainID,
		BlockHeight:         cHeight,
	}

	txBytes, err := SignNativeTx(orcSeeds[0], native, nativeCliCtx, msg)

	if err != nil {
		return err
	}

	grpcConnN, err := grpc.Dial(native.GRPCAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer func(grpcConnN *grpc.ClientConn) {
		err := grpcConnN.Close()
		if err != nil {

		}
	}(grpcConnN)

	txClient := sdkTx.NewServiceClient(grpcConnN)

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
