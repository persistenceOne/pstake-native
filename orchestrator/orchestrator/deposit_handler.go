package orchestrator

import (
	"context"
	"fmt"
	stdlog "log"
	"strings"
	"time"

	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	tendermintTypes "github.com/tendermint/tendermint/rpc/core/types"
	"google.golang.org/grpc"
)

func (c *CosmosChain) DepositHandler(valAddr string, orcSeeds []string, nativeCliCtx cosmosClient.Context, clientCtx cosmosClient.Context, native *NativeChain, depositHeight int64, protoCodec *codec.ProtoCodec) error {
	var resultTxs []*tendermintTypes.ResultTx
	fmt.Println("Node is connected")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	page := 1
	perPage := 100

	txResult, err := c.Client.TxSearch(ctx, fmt.Sprintf("tx.height=%d", depositHeight), true, &page, &perPage, "asc")
	if err != nil {
		return err
	}
	if txResult.TotalCount <= perPage {
		err := handleTxResult(c, valAddr, orcSeeds, depositHeight, native, nativeCliCtx, clientCtx, txResult.Txs, protoCodec)
		if err != nil {
			return err
		}
		return nil
	}

	resultTxs = append(resultTxs, txResult.Txs...)
	totalPages := (txResult.TotalCount / perPage) + 1
	if txResult.TotalCount%perPage == 0 {
		totalPages = txResult.TotalCount / perPage
	}
	for i := page + 1; i <= totalPages; i++ {
		txResult, err = c.Client.TxSearch(ctx, fmt.Sprintf("tx.height=%d", depositHeight), true, &i, &perPage, "asc")
		if err != nil {
			fmt.Println("Error in getting txs", err)
			return err
		}
		resultTxs = append(resultTxs, txResult.Txs...)
	}
	err = handleTxResult(c, valAddr, orcSeeds, depositHeight, native, nativeCliCtx, clientCtx, resultTxs, protoCodec)
	return nil
}

func handleTxResult(c *CosmosChain, valAddr string, orcSeeds []string, depositHeight int64, native *NativeChain, nativeCliCtx cosmosClient.Context, clientCtx cosmosClient.Context, resultTxs []*tendermintTypes.ResultTx, protoCodec *codec.ProtoCodec) error {
	for i, transaction := range resultTxs {
		fmt.Println("Cosmos Deposit Tx:", transaction.Hash.String(), fmt.Sprintf("(%d)", i+1))
		err := processCustodialDepositTxAndTranslateToNative(c, valAddr, orcSeeds, depositHeight, native, nativeCliCtx, clientCtx, transaction)
		if err != nil {
			fmt.Println("Error in getting custodial deposit tx", err)
			return err
		}
		return nil
	}
	return nil
}

func processCustodialDepositTxAndTranslateToNative(chain *CosmosChain, valAddr string, orcSeeds []string, depositHeight int64, native *NativeChain, nativeCLiCtx cosmosClient.Context, clientCtx cosmosClient.Context, txResult *tendermintTypes.ResultTx) error {
	fmt.Println("Cosmos Deposit Tx")
	if txResult.TxResult.GetCode() == 0 {

		txInterface, err := clientCtx.TxConfig.TxDecoder()(txResult.Tx)
		if err != nil {
			panic(err)
			return err
		}
		tx, ok := txInterface.(signing.Tx)
		if !ok {
			return err
		}
		memo := strings.TrimSpace(tx.GetMemo())

		for _, msg := range tx.GetMsgs() {
			switch txMsg := msg.(type) {
			case *banktypes.MsgSend:
				if txMsg.ToAddress == chain.CustodialAddress.String() {
					for _, coin := range txMsg.Amount {
						//TODO: handle multiple keys for signing
						_, addr := GetPivKeyAddress(native.AccountPrefix, native.CoinType, orcSeeds[0])

						stdlog.Println("orchestrator address, ", addr)
						msg = &cosmosTypes.MsgMintTokensForAccount{
							AddressFromMemo:     memo,
							OrchestratorAddress: addr,
							Amount:              sdk.NewCoin(coin.Denom, coin.Amount),
							TxHash:              txResult.Hash.String(),
							ChainID:             chain.ChainID,
							BlockHeight:         depositHeight,
						}

						fmt.Println(msg)
						fmt.Println("test-test")
						fmt.Println(native)
						txBytes, err := SignNativeTx(orcSeeds[0], native, nativeCLiCtx, msg)
						if err != nil {
							return err
						}
						grpcConn, _ := grpc.Dial(native.GRPCAddr, grpc.WithInsecure())
						if err != nil {
							return err
						}
						defer func(grpcConn *grpc.ClientConn) {
							err := grpcConn.Close()
							if err != nil {

							}
						}(grpcConn)

						txClient := sdkTx.NewServiceClient(grpcConn)

						fmt.Println("client created")

						res, err := txClient.BroadcastTx(context.Background(),
							&sdkTx.BroadcastTxRequest{
								Mode:    sdkTx.BroadcastMode_BROADCAST_MODE_SYNC,
								TxBytes: txBytes,
							},
						)
						if err != nil {
							return err
						}

						fmt.Println(res.TxResponse.Code, res.TxResponse.TxHash, res)

						return nil
					}
				}

			}
		}

	}
	return nil
}
