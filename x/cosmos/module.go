/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package cosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simTypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	cosmosCli "github.com/persistenceOne/pstake-native/x/cosmos/client/cli"
	cosmosRest "github.com/persistenceOne/pstake-native/x/cosmos/client/rest"
	cosmosKeeper "github.com/persistenceOne/pstake-native/x/cosmos/keeper"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic defines the basic application module used by the cosmos module.
type AppModuleBasic struct {
	cdc codec.Codec
}

var _ module.AppModuleBasic = AppModuleBasic{}

// Name returns the cosmos module's name.
func (AppModuleBasic) Name() string {
	return cosmosTypes.ModuleName
}

// RegisterLegacyAminoCodec registers the cosmos module's types on the given LegacyAmino codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cosmosTypes.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the module's interface types
func (b AppModuleBasic) RegisterInterfaces(registry codecTypes.InterfaceRegistry) {
	cosmosTypes.RegisterInterfaces(registry)
}

// DefaultGenesis returns default genesis state as raw bytes for the cosmos
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(cosmosTypes.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the cosmos module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var data cosmosTypes.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", cosmosTypes.ModuleName, err)
	}

	return cosmosTypes.ValidateGenesis(data)
}

// RegisterRESTRoutes registers the REST routes for the cosmos module.
func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr *mux.Router) {
	cosmosRest.RegisterHandlers(clientCtx, rtr)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the cosmos module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	_ = cosmosTypes.RegisterQueryHandlerClient(context.Background(), mux, cosmosTypes.NewQueryClient(clientCtx))
}

// GetTxCmd returns no root tx command for the cosmos module.
func (AppModuleBasic) GetTxCmd() *cobra.Command { return cosmosCli.NewTxCmd() }

// GetQueryCmd returns the root query command for the cosmos module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cosmosCli.GetQueryCmd()
}

//____________________________________________________________________________

// AppModule implements an application module for the cosmos module.
type AppModule struct {
	AppModuleBasic

	keeper cosmosKeeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper cosmosKeeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		keeper:         keeper,
	}
}

// Name returns the cosmos module's name.
func (AppModule) Name() string {
	return cosmosTypes.ModuleName
}

// RegisterInvariants registers the cosmos module invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route returns the message routing key for the cosmos module.
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(cosmosTypes.RouterKey, NewHandler(am.keeper))
}

// QuerierRoute returns the cosmos module's querier route name.
func (AppModule) QuerierRoute() string {
	return cosmosTypes.QuerierRoute
}

// ConsensusVersion returns the cosmos module's consensus version number.
func (am AppModule) ConsensusVersion() uint64 {
	return 1
}

// LegacyQuerierHandler returns the cosmos module sdk.Querier.
func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return cosmosKeeper.NewQuerier(am.keeper, legacyQuerierCdc)
}

// RegisterServices registers a gRPC query service to respond to the
// module-specific gRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	cosmosTypes.RegisterMsgServer(cfg.MsgServer(), cosmosKeeper.NewMsgServerImpl(am.keeper))
	cosmosTypes.RegisterQueryServer(cfg.QueryServer(), am.keeper)
}

// InitGenesis performs genesis initialization for the cosmos module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState cosmosTypes.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	cosmosKeeper.InitGenesis(ctx, am.keeper, &genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the cosmos
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the cosmos module.
func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {
}

// EndBlock returns the end blocker for the cosmos module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.keeper)
	return []abci.ValidatorUpdate{}
}

//____________________________________________________________________________

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of the cosmos module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	//simulation.RandomizedGenState(simState)
	panic("implement me")
}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simTypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized cosmos param changes for the simulator.
func (AppModule) RandomizedParams(r *rand.Rand) []simTypes.ParamChange {
	//return simulation.ParamChanges(r)
	panic("implement me")
}

// RegisterStoreDecoder registers a decoder for cosmos module's types.
func (AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations doesn't return any cosmos module operation.
func (AppModule) WeightedOperations(_ module.SimulationState) []simTypes.WeightedOperation {
	return nil
}
