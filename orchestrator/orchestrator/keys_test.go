package oracle

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"testing"
)

func TestS(t *testing.T) {
	k, address, _ := createMemoryKeyFromMnemonic("april patch recipe debate remove hurdle concert gesture design near predict enough color tail business imitate twelve february punch cheap vanish december cool wheel")
	fmt.Println(k)
	fmt.Println(address)

	hd.CreateHDPath(750, 0, 0)
	fmt.Println("sss")
}
