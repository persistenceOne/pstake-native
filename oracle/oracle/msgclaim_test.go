package oracle

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestF(t *testing.T) {
	//grpcConn, _ := grpc.Dial("grpc-cosmoshub-ia.notional.ventures:443", grpc.WithInsecure())
	//defer func(grpcConn *grpc.ClientConn) {
	//	err := grpcConn.Close()
	//	if err != nil {
	//	}
	//}(grpcConn)

	rpcClientC, _ := newRPCClient("https://rpc.cosmoshub-4.audit.one:443", 1*time.Second)
	height := int64(11093208)

	blockResults, err := rpcClientC.BlockResults(context.Background(), &height)

	tmp := blockResults.TxsResults
	fmt.Println(tmp)

	for _, eve := range tmp {
		eve1 := eve.Events
		for _, events := range eve1 {
			fmt.Println("{")
			data := events
			tmp2 := data.Attributes

			fmt.Print(data.Type)
			for _, dt := range tmp2 {

				dt3 := dt.String()

				//dt4 := dt.GetValue()
				fmt.Printf("--> %v <--", dt3)
				fmt.Println()
				fmt.Println("}")
				fmt.Println()
			}
		}
	}

	//txSlice := events["tx"]

	if err != nil {
		panic(err)
	}

}
