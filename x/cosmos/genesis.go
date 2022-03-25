/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package cosmos

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/x/cosmos/types"
)

//// InitGenesis new cosmos genesis
//func InitGenesis(ctx sdk.Context, keeper Keeper, data *GenesisState) {
//	keeper.SetParams(ctx, data.Params)
//	keeper.SetProposalID(ctx, 1)
//	//keeper.SetVotingParams(ctx, data.Params.CosmosProposalParams)
//	//TODO add remaining
//}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) *GenesisState {
	params := keeper.GetParams(ctx)
	return NewGenesisState(params, nil, nil, types.OutgoingTx{})
}
