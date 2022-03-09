package keeper

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types"
)

func ConvertTxResponseFromString(s string) types.TxResponse {
	a := json.Unmarshal([]byte(s), &types.TxResponse{})
	fmt.Println(a)
	return types.TxResponse{}
}
