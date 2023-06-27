package keeper_test

import (
	"fmt"
	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/keeper"
	"testing"
	"time"
)

// ////////// fetch data from json:
// curl -X GET -H "Content-Type: application/json" -H "x-cosmos-block-height: 15884400" "https://rest.cosmos.audit.one/cosmos/staking/v1beta1/delegators/cosmos13t4996czrgft9gw43epuwauccrldu5whx6uprjdmvsmuf7ylg8yqcxgzk3/unbonding_delegations"
// creationheight 15826532 => epoch 232, 15881865 => 236
type UnbondingEntry struct {
	CreationHeight          string    `json:"creation_height"`
	CompletionTime          time.Time `json:"completion_time"`
	InitialBalance          string    `json:"initial_balance"`
	Balance                 string    `json:"balance"`
	UnbondingId             string    `json:"unbonding_id"`
	UnbondingOnHoldRefCount string    `json:"unbonding_on_hold_ref_count"`
}

type UnbondingResponse struct {
	DelegatorAddress string           `json:"delegator_address"`
	ValidatorAddress string           `json:"validator_address"`
	Entries          []UnbondingEntry `json:"entries"`
}

type Unbondings struct {
	UnbondingResponses []UnbondingResponse `json:"unbonding_responses"`
}

func TestParseHostAccountUnbondings(t *testing.T) {
	mintDenom := "stk/uatom"
	baseDenom := "uatom"
	// create a map to quickly access each undelegation epoch entry and initialise it
	hostAccountUndelegationsMap := keeper.ParseHostAccountUnbondings(mintDenom, baseDenom)
	fmt.Println(hostAccountUndelegationsMap)
}
