package oracle

import (
	"context"
	"fmt"
	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	txD "github.com/cosmos/cosmos-sdk/types/tx"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/persistenceOne/pstake-native/oracle/utils"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	"google.golang.org/grpc"
	logg "log"
	"strconv"
)

func (c *CosmosChain) ProposalHandler(propId string, orcSeeds []string, nativeCliCtx cosmosClient.Context, native *NativeChain, chain *CosmosChain, cHeight int64) error {
	propID, err := strconv.ParseUint(propId, 10, 64)

	if err != nil {
		return err
	}

	grpcConn, err := grpc.Dial(chain.GRPCAddr, grpc.WithInsecure())
	defer func(grpcConn *grpc.ClientConn) {
		err := grpcConn.Close()
		if err != nil {
			logg.Println("GRPC Connection error")
		}
	}(grpcConn)

	GovClient := govtypes.NewQueryClient(grpcConn)

	fmt.Println("gov query client connected")

	PropResult, err := GovClient.Proposal(context.Background(),
		&govtypes.QueryProposalRequest{ProposalId: propID},
	)

	Proposal := PropResult.Proposal

	_, addr := utils.GetSDKPivKeyAndAddress(orcSeeds[0])
	msg := &cosmosTypes.MsgMakeProposal{
		Title:               Proposal.GetTitle(),
		Description:         Proposal.ProposalType(),
		OrchestratorAddress: addr.String(),
		ProposalID:          Proposal.ProposalId,
		ChainID:             chain.ChainID,
		BlockHeight:         cHeight,
		VotingStartTime:     Proposal.VotingStartTime,
		VotingEndTime:       Proposal.VotingEndTime,
	}

	txBytes, err := utils.SignNativeTx(orcSeeds[0], native, nativeCliCtx, msg)

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
