package keeper_test

import (
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"testing"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibctesting "github.com/cosmos/ibc-go/v3/testing"
	"github.com/stretchr/testify/suite"

	"github.com/persistenceOne/pstake-native/app"
	"github.com/persistenceOne/pstake-native/app/helpers"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

var (
	allowListedValidators = types.AllowListedValidators{
		AllowListedValidators: []types.AllowListedValidator{
			{
				ValidatorAddress: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt",
				TargetWeight:     sdk.NewDecWithPrec(33, 2),
			},
			{
				ValidatorAddress: "cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2",
				TargetWeight:     sdk.NewDecWithPrec(33, 2),
			},
			{
				ValidatorAddress: "cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5",
				TargetWeight:     sdk.NewDecWithPrec(34, 2),
			},
		},
	}
	ChainID          = "cosmoshub-4"
	ConnectionID     = "connection-0"
	TransferChannel  = "channel-0"
	TransferPort     = "transfer"
	BaseDenom        = "uatom"
	MintDenom        = "ustkatom"
	MinDeposit       = sdk.NewInt(5)
	PstakeFeeAddress = "persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9"
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

	// set host chain params
	depositFee, err := sdk.NewDecFromStr("0.01")
	suite.NoError(err)

	restakeFee, err := sdk.NewDecFromStr("0.02")
	suite.NoError(err)

	unstakeFee, err := sdk.NewDecFromStr("0.03")
	suite.NoError(err)

	hostChainParams := types.NewHostChainParams(
		ChainID,
		ConnectionID,
		TransferChannel,
		TransferPort,
		BaseDenom,
		MintDenom,
		PstakeFeeAddress,
		MinDeposit,
		depositFee,
		restakeFee,
		unstakeFee,
	)
	suite.app.LSCosmosKeeper.SetHostChainParams(suite.ctx, hostChainParams)
	suite.app.LSCosmosKeeper.SetAllowListedValidators(ctx, allowListedValidators)
}

func (suite *IntegrationTestSuite) TestMintToken() {
	pstakeApp, ctx := suite.app, suite.ctx
	testParams := pstakeApp.LSCosmosKeeper.GetHostChainParams(ctx)

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
