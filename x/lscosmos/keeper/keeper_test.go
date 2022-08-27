package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibctesting "github.com/cosmos/ibc-go/v3/testing"
	"github.com/stretchr/testify/suite"

	"github.com/persistenceOne/pstake-native/app"
	"github.com/persistenceOne/pstake-native/app/helpers"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

func init() {
	ibctesting.DefaultTestingAppInit = helpers.SetupTestingApp
}

type IntegrationTestSuite struct {
	suite.Suite

	app        *app.PstakeApp
	ctx        sdk.Context
	govHandler govtypes.Handler

	coordinator *ibctesting.Coordinator
	chainA      *ibctesting.TestChain
	chainB      *ibctesting.TestChain
	path        *ibctesting.Path
}

func newPstakeAppPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort

	return path
}

func GetPstakeApp(chain *ibctesting.TestChain) *app.PstakeApp {
	app1, ok := chain.App.(*app.PstakeApp)
	if !ok {
		panic("not pstake app")
	}

	return app1
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) SetupTest() {
	_, pstakeApp, ctx := helpers.CreateTestApp()

	keeper := pstakeApp.LSCosmosKeeper

	params := types.DefaultParams()
	keeper.SetParams(ctx, params)

	suite.app = &pstakeApp
	suite.ctx = ctx

	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(2))

	suite.path = newPstakeAppPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupConnections(suite.path)
}

func (suite *IntegrationTestSuite) TestMintToken() {
	pstakeApp, ctx := suite.app, suite.ctx

	testParams := types.RegisterCosmosChainProposal{
		Title:                 "register cosmos chain proposal",
		Description:           "this proposal register cosmos chain params in the chain",
		ModuleEnabled:         true,
		ConnectionID:          "test connection",
		TransferChannel:       "test-channel-1",
		TransferPort:          "transfer",
		BaseDenom:             "uatom",
		MintDenom:             "ustkatom",
		MinDeposit:            sdk.OneInt().MulRaw(5),
		AllowListedValidators: types.AllowListedValidators{AllowListedValidators: []types.AllowListedValidator{{ValidatorAddress: "addr", TargetWeight: sdk.OneDec()}}},
		PstakeDepositFee:      sdk.ZeroDec(),
		PstakeRestakeFee:      sdk.ZeroDec(),
		PstakeUnstakeFee:      sdk.ZeroDec(),
	}

	ibcDenom := ibctransfertypes.GetPrefixedDenom(testParams.TransferPort, testParams.TransferChannel, testParams.BaseDenom)
	balanceOfIbcToken := sdk.NewInt64Coin(ibcDenom, 100)
	mintAmountDec := balanceOfIbcToken.Amount.ToDec().Mul(pstakeApp.LSCosmosKeeper.GetCValue(ctx))
	toBeMintedTokens, _ := sdk.NewDecCoinFromDec(testParams.MintDenom, mintAmountDec).TruncateDecimal()

	addr := sdk.AccAddress("addr_______________")
	acc := pstakeApp.AccountKeeper.NewAccountWithAddress(ctx, addr)
	pstakeApp.AccountKeeper.SetAccount(ctx, acc)
	suite.Require().NoError(simapp.FundAccount(pstakeApp.BankKeeper, ctx, addr, sdk.NewCoins(balanceOfIbcToken)))

	suite.Require().NoError(pstakeApp.LSCosmosKeeper.MintTokens(ctx, toBeMintedTokens, addr))

	currBalance := pstakeApp.BankKeeper.GetBalance(ctx, addr, testParams.MintDenom)

	suite.Require().Equal(toBeMintedTokens, currBalance)
}
