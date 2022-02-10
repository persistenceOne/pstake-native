/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package cosmos

import (
	"github.com/persistenceOne/pstake-native/x/cosmos/keeper"
	"github.com/persistenceOne/pstake-native/x/cosmos/types"
)

const (
	ModuleName        = types.ModuleName
	DefaultParamspace = types.DefaultParamspace
	StoreKey          = types.StoreKey
)

var (
	// functions aliases

	NewKeeper       = keeper.NewKeeper
	NewGenesisState = types.NewGenesisState
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
	Params       = types.Params
)
