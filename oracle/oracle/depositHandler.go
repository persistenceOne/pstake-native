package oracle

import (
	"context"
	"fmt"
	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txD "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/persistenceOne/pStake-native/oracle/constants"
	"github.com/persistenceOne/pStake-native/oracle/utils"
	cosmosTypes "github.com/persistenceOne/pStake-native/x/cosmos/types"
	tendermintTypes "github.com/tendermint/tendermint/rpc/core/types"
	"google.golang.org/grpc"
	"strings"
	"time"
)

func (c *CosmosChain) DepositHandler(nativeCliCtx cosmosClient.Context, clientCtx cosmosClient.Context, native *NativeChain, depositHeight int64, protoCodec *codec.ProtoCodec) error {
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
		err := handleTxResult(c, depositHeight, native, nativeCliCtx, clientCtx, txResult.Txs, protoCodec)
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
	err = handleTxResult(c, depositHeight, native, nativeCliCtx, clientCtx, resultTxs, protoCodec)
	return nil
}

func handleTxResult(c *CosmosChain, depositHeight int64, native *NativeChain, nativeCliCtx cosmosClient.Context, clientCtx cosmosClient.Context, resultTxs []*tendermintTypes.ResultTx, protoCodec *codec.ProtoCodec) error {
	for i, transaction := range resultTxs {
		fmt.Println("Cosmos Deposit Tx:", transaction.Hash.String(), fmt.Sprintf("(%d)", i+1))
		err := processCustodialDepositTxAndTranslateToNative(c, depositHeight, native, nativeCliCtx, clientCtx, transaction)
		if err != nil {
			fmt.Println("Error in getting custodial deposit tx", err)
			return err
		}
		return nil
	}
	return nil
}

func processCustodialDepositTxAndTranslateToNative(chain *CosmosChain, depositHeight int64, native *NativeChain, nativeCLiCtx cosmosClient.Context, clientCtx cosmosClient.Context, txResult *tendermintTypes.ResultTx) error {
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
		fmt.Println(memo + "<--memo")

		for _, msg := range tx.GetMsgs() {
			switch txMsg := msg.(type) {
			case *banktypes.MsgSend:
				if txMsg.ToAddress == chain.CustodialAddress.String() {
					for _, coin := range txMsg.Amount {
						_, addr := utils.GetSDKPivKeyAndAddress()
						fmt.Println(addr)
						fmt.Println(coin)
						msg = &cosmosTypes.MsgMintTokensForAccount{
							AddressFromMemo:     memo,
							OrchestratorAddress: addr.String(),
							Amount:              sdk.NewCoins(sdk.NewCoin(coin.Denom, coin.Amount)),
							TxHash:              txResult.Hash.String(),
							ChainID:             chain.ChainID,
							BlockHeight:         depositHeight,
						}

						fmt.Println(msg.String(), "<--nativeMsg")

						txBytes, err := utils.SignMintTx(nativeCLiCtx, msg)
						fmt.Println(txBytes, "<--signedMsg")
						if err != nil {
							return err
						}
						grpcConn, _ := grpc.Dial(constants.NativeGRPCAddr, grpc.WithInsecure())
						defer func(grpcConn *grpc.ClientConn) {
							err := grpcConn.Close()
							if err != nil {

							}
						}(grpcConn)

						txClient := txD.NewServiceClient(grpcConn)

						fmt.Println("client created")

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
						//return nil
						//ctx := chain.CLIContext(0).WithFromAddress(from).WithBroadcastMode(configuration.GetAppConfig().Tendermint.BroadcastMode)

						return nil
					}
				}

			}
		}

	}
	return nil
}