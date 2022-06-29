package keeper_test

import (
	"fmt"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	"strings"
	"testing"
)

//status: ""
//tx:
//auth_info:
//fee:
//amount: []
//gas_limit: "400000"
//granter: ""
//payer: ""
//signer_infos: []
//body:
//extension_options: []
//memo: ""
//messages:
//- '@type': /cosmos.authz.v1beta1.MsgExec
//grantee: cosmos1g42ycjrd7dzu2r9af3xnjw09q3pfscffe73nnn
//msgs:
//- '@type': /cosmos.staking.v1beta1.MsgDelegate
//amount:
//amount: "5000000"
//denom: stake
//delegator_address: cosmos15ddw7dkp56zytf3peshxr8fwn5w76y4g462ql2
//validator_address: cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt
//- '@type': /cosmos.staking.v1beta1.MsgDelegate
//amount:
//amount: "2000000"
//denom: stake
//delegator_address: cosmos15ddw7dkp56zytf3peshxr8fwn5w76y4g462ql2
//validator_address: cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2
//- '@type': /cosmos.staking.v1beta1.MsgDelegate
//amount:
//amount: "2000000"
//denom: stake
//delegator_address: cosmos15ddw7dkp56zytf3peshxr8fwn5w76y4g462ql2
//validator_address: cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2
//- '@type': /cosmos.staking.v1beta1.MsgDelegate
//amount:
//amount: "1000000"
//denom: stake
//delegator_address: cosmos15ddw7dkp56zytf3peshxr8fwn5w76y4g462ql2
//validator_address: cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5
//non_critical_extension_options: []
//timeout_height: "0"
//signatures: []
//txHash: ""

func TestKeeper_SetOutgoingTxnSignaturesAndEmitEvent(t *testing.T) {

	bytes := []byte("cosmos")
	txHash := cosmosTypes.BytesToHexUpper(bytes)
	txHash1 := strings.ToUpper(txHash)
	fmt.Println(txHash1)

}
