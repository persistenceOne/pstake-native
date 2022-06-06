package utils

import (
	"fmt"
	"github.com/persistenceOne/pstake-native/oracle/constants"
	"testing"
)

func TestA(t *testing.T) {
	a, b := GetSDKPivKeyAndAddress(constants.Seed[0])
	fmt.Println(a, b)
}
