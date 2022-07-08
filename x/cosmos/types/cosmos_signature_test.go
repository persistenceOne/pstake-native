package types_test

import (
	"fmt"
	"github.com/persistenceOne/pstake-native/app"
	"github.com/persistenceOne/pstake-native/x/cosmos/types"
	"testing"
)

func TestAccAddressFromBech32(t *testing.T) {
	app.SetAddressPrefixes()
	accAddress, _ := types.AccAddressFromBech32("persistencevaloper183g695ap32wnds5k9xwd3yq997dqxudfz524f3", "persistencevaloper")
	fmt.Println(accAddress.String())
}
