package utils

import (
	"fmt"
	"testing"
)

func TestA(t *testing.T) {
	a, b := GetSDKPivKeyAndAddress()
	fmt.Println(a, b)
}
