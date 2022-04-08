package types_test

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"testing"
)

func TestNewMsgWithdrawStkAsset(t *testing.T) {
	addr := types.NewModuleAddress("cosmos")
	fmt.Println(addr.String())
	//cosmos1fjlpjuttr2nn5e7ufv5vxsu3s7d4qvjeula3h5
}
