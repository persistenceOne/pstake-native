package oracle

import (
	"errors"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"strings"
	"testing"
)

func TestA(t *testing.T) {
	privkey, _ := GetPivKeyAddress("persistence", 750, "bomb sand fashion torch return coconut color captain vapor inhale lyrics lady grant ordinary lazy decrease quit devote paddle impulse prize equip hip ball")
	_, _ = Bech32ifyAddressBytes("persistence", sdkTypes.AccAddress(privkey.PubKey().Address()))
	//addrStr2, _ := sdkTypes.Bech32ifyAddressBytes("persistence", addr)
	//fmt.Println(len(addrStr2))
	//fmt.Println(addr.String(), addrStr, addrStr2)
}

func ValAddressFromBech32(address, prefix string) (valAddr sdkTypes.ValAddress, err error) {
	if len(strings.TrimSpace(address)) == 0 {
		return sdkTypes.ValAddress{}, errors.New("empty address string is not allowed")
	}

	bz, err := sdkTypes.GetFromBech32(address, prefix)
	if err != nil {
		return nil, err
	}

	err = sdkTypes.VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return bz, nil
}

//func Bech32ifyAddressBytes(prefix string, address sdkTypes.AccAddress) (string, error) {
//	if address.Empty() {
//		return "", nil
//	}
//	if len(address.Bytes()) == 0 {
//		return "", nil
//	}
//	if len(prefix) == 0 {
//		return "", errors.New("prefix cannot be empty")
//	}
//	return bech32.ConvertAndEncode(prefix, address.Bytes())
//}

func Bech32ifyValAddressBytes(prefix string, address sdkTypes.ValAddress) (string, error) {
	if address.Empty() {
		return "", nil
	}
	if len(address.Bytes()) == 0 {
		return "", nil
	}
	if len(prefix) == 0 {
		return "", errors.New("prefix cannot be empty")
	}
	return bech32.ConvertAndEncode(prefix, address.Bytes())
}
