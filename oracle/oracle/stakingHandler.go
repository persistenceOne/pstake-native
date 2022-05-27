package oracle

import (
	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
)

func (n *NativeChain) StakingHandler(valAddr string, orcSeeds []string, nativeCliCtx cosmosClient.Context, clientCtx cosmosClient.Context, native *NativeChain, depositHeight int64, protoCodec *codec.ProtoCodec) error {

}
