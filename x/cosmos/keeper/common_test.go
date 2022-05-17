package keeper_test

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/app"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
	"time"
)

var defaultConsensusParams = &abci.ConsensusParams{
	Block: &abci.BlockParams{
		MaxBytes: 200000,
		MaxGas:   2000000,
	},
	Evidence: &tmproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &tmproto.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}

func newTestApp(isCheckTx bool, withGenesis bool) app.PstakeApp {
	db := tmdb.NewMemDB()
	// encCdc := app.MakeTestEncodingConfig()

	encoding := app.MakeEncodingConfig()
	testApp := app.NewGaiaApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, app.DefaultNodeHome, 5, encoding, simapp.EmptyAppOptions{})
	genesis := app.GenesisState{}
	if withGenesis {
		genesis = app.NewDefaultGenesisState()
	}

	if !isCheckTx {
		// InitChain must be called to stop deliverState from being nil
		stateBytes, err := json.MarshalIndent(genesis, "", " ")
		if err != nil {
			panic(err)
		}
		// Initialize the chain
		testApp.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: defaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}
	return *testApp
}

func createTestInput() (*codec.LegacyAmino, app.PstakeApp, sdk.Context) {
	app := newTestApp(false, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	return app.LegacyAmino(), app, ctx
}
