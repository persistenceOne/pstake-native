package oracle

import (
	"context"
	"errors"
	"fmt"
	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cosmosTypes "github.com/persistenceOne/pStake-native/x/cosmos/types"

	"github.com/tendermint/tendermint/types"
	"time"
)

func (c *CosmosChain) ProposalHandler(valAddr string, orcSeeds []string, nativeCliCtx cosmosClient.Context, clientCtx cosmosClient.Context, native *NativeChain, depositHeight int64, protoCodec *codec.ProtoCodec) error {

	query := "active_proposal"
	subscriber := fmt.Sprintf("%s-subscriber", "oracle")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	eventChan, err := c.Client.Subscribe(ctx, subscriber, query, 1000)
	if err != nil {
		return err
	}

	select {
	case event := <-eventChan:
		proposal := (event.Data.(types.TMEventData))
		cosmosTypes.MsgMakeProposal{}

	case <-ctx.Done():
		return errors.New("timed out waiting for event")
	}

}
