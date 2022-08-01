package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/simapp"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ibcTransferTypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	"github.com/stretchr/testify/suite"

	"github.com/persistenceOne/pstake-native/app"
	"github.com/persistenceOne/pstake-native/app/helpers"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	app        *app.PstakeApp
	ctx        sdk.Context
	govHandler govtypes.Handler
}

func (suite *IntegrationTestSuite) SetupTest() {
	_, app, ctx := helpers.CreateTestApp()

	keeper := app.LSCosmosKeeper

	params := types.DefaultParams()
	keeper.SetParams(ctx, params)

	suite.app = &app
	suite.ctx = ctx
}

func testProposal(
	title,
	description,
	connection,
	channel,
	transfer,
	uatom,
	ustkatom,
	minDeposit,
	pStakeFee string) *types.RegisterCosmosChainProposal {
	return types.NewRegisterCosmosChainProposal(
		title,
		description,
		connection,
		channel,
		transfer,
		uatom,
		ustkatom,
		minDeposit,
		pStakeFee,
	)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) TestMintToken() {
	pstakeApp, ctx := suite.app, suite.ctx

	testParams := types.RegisterCosmosChainProposal{
		Title:                "register cosmos chain proposal",
		Description:          "this proposal register cosmos chain params in the chain",
		IBCConnection:        "test connection",
		TokenTransferChannel: "test-channel-1",
		TokenTransferPort:    "test-transfer",
		BaseDenom:            "uatom",
		MintDenom:            "ustkatom",
		MinDeposit:           "5",
		PStakeFee:            "0.1",
	}

	ibcDenom := ibcTransferTypes.GetPrefixedDenom(testParams.TokenTransferPort, testParams.TokenTransferChannel, testParams.BaseDenom)
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
