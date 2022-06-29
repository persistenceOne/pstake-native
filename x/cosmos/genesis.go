/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package cosmos

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) *GenesisState {
	//TODO
	params := keeper.GetParams(ctx)
	return NewGenesisState(params)
}
