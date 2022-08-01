package orchestrator

import (
	"context"
	stdlog "log"
	"strconv"
	"time"

	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	"github.com/tendermint/tendermint/types"
)

func (c *CosmosChain) DepositTxEventForBlock(blockHeight int64) error {
	client, err := rpchttp.New(c.RPCAddr, "/websocket")
	if err != nil {
		return err
	}
	err = client.Start()
	if err != nil {
		return err
	}
	defer client.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	query := "tm.event = 'Tx' AND transfer.recipient = '" + string(c.CustodialAddress) + "' AND tx.height = '" + strconv.FormatInt(blockHeight, 10) + "'"
	txs, err := client.Subscribe(ctx, "orchestrator", query)
	if err != nil {
		return err
	}

	go func() {
		for e := range txs {
			//relay to native chain
			stdlog.Println("got ", e.Data.(types.EventDataTx))
		}
	}()
	return nil
}

func (c *CosmosChain) ActiveProposalEventHandler(blockHeight int64) error {
	client, err := rpchttp.New(c.RPCAddr, "/websocket")
	if err != nil {
		return err
	}
	err = client.Start()
	if err != nil {
		return err
	}
	defer func(client *rpchttp.HTTP) {
		err := client.Stop()
		if err != nil {

		}
	}(client)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	query := "active_proposal"

	proposals, err := client.Subscribe(ctx, "orchestrator", query)

	if err != nil {
		return err
	}

	go func() {
		for e := range proposals {
			//relay to native chain
			stdlog.Println("got ", e.Data.(types.EventDataTx))
		}
	}()
	return nil

}
