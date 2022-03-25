package oracle

import (
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

func (c *CosmosChain) DepositTxEventForBlock(BlockHeight int64) error {
	client, err := rpchttp.New(c.RPCAddr, "/websocket")
	if err != nil {
		return err
	}
	err = client.Start()
	if err != nil {
		return err
	}
	defer client.Stop()
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	//defer cancel()

	//query := "tm.event = 'Tx' AND transfer.recipient = '" + string(c.CustodialAddress) + "' AND tx.height = '" + string(BlockHeight) + "'"
	////txs, err := client.Subscribe(ctx, "orchestrator", query)
	//if err != nil {
	//	return err
	//}
	//
	//go func() {
	//	for e := range txs {
	//		//relay to native chain
	//		fmt.Println("got ", e.Data.(types.EventDataTx))
	//	}
	//}()
	return nil
}
