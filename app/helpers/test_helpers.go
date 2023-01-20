package helpers

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctesting "github.com/cosmos/ibc-go/v4/testing"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	tmdb "github.com/tendermint/tm-db"

	"github.com/persistenceOne/pstake-native/app"
)

// SimAppChainID hardcoded chainID for simulation
const (
	SimAppChainID = "pstake-app"
)

var DefaultConsensusParams = &abci.ConsensusParams{
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
	testApp := app.NewpStakeApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, app.DefaultNodeHome, 5, encoding, simapp.EmptyAppOptions{})
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
				ConsensusParams: DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}
	return *testApp
}

func CreateTestApp() (*codec.LegacyAmino, app.PstakeApp, sdk.Context) {
	testApp := newTestApp(false, false)
	ctx := testApp.BaseApp.NewContext(false, tmproto.Header{})

	return testApp.LegacyAmino(), testApp, ctx
}

type EmptyAppOptions struct{}

func (EmptyAppOptions) Get(o string) interface{} { return nil }

func Setup(t *testing.T, isCheckTx bool, invCheckPeriod uint) *app.PstakeApp {
	t.Helper()

	testApp, genesisState := setup(!isCheckTx, invCheckPeriod)
	if !isCheckTx {
		// InitChain must be called to stop deliverState from being nil
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		require.NoError(t, err)

		// Initialize the chain
		testApp.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return testApp
}

func setup(withGenesis bool, invCheckPeriod uint) (*app.PstakeApp, app.GenesisState) {
	db := tmdb.NewMemDB()
	encCdc := app.MakeEncodingConfig()
	testApp := app.NewpStakeApp(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		app.DefaultNodeHome,
		invCheckPeriod,
		encCdc,
		EmptyAppOptions{},
	)
	if withGenesis {
		return testApp, app.NewDefaultGenesisState()
	}

	return testApp, app.GenesisState{}
}

// SetupTestingApp initializes the IBC-go testing application
func SetupTestingApp() (ibctesting.TestingApp, map[string]json.RawMessage) {
	db := tmdb.NewMemDB()
	newpStakeApp := app.NewpStakeApp(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		app.DefaultNodeHome,
		5,
		app.MakeEncodingConfig(),
		EmptyAppOptions{},
	)
	return newpStakeApp, app.NewDefaultGenesisState()
}
