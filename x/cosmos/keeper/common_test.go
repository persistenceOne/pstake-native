package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/persistenceOne/pstake-native/app"
	"github.com/persistenceOne/pstake-native/app/params"
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

func newTestApp() app.PstakeApp {
	db := tmdb.NewMemDB()
	// encCdc := app.MakeTestEncodingConfig()

	encoding := params.MakeEncodingConfig()
	testApp := app.NewGaiaApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, simapp.DefaultNodeHome, 5, encoding, simapp.EmptyAppOptions{})
	//return testApp, app.NewDefaultGenesisState(encoding.Marshaler)
	return *testApp
}
