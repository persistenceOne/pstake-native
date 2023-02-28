package cmd

import (
	"errors"
	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/server"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	ibcclienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	tmcfg "github.com/tendermint/tendermint/config"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	pstakeApp "github.com/persistenceOne/pstake-native/v2/app"
	"github.com/persistenceOne/pstake-native/v2/app/params"
)

// NewRootCmd creates a new root command for simd. It is called once in the
// main function.
func NewRootCmd() (*cobra.Command, params.EncodingConfig) {
	encodingConfig := pstakeApp.MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithHomeDir(pstakeApp.DefaultNodeHome).
		WithViper("")

	rootCmd := &cobra.Command{
		Use:   "pstaked",
		Short: "Stargate Cosmos Hub App",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err = client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customTemplate, custompStakeConfig := initAppConfig()
			customTMConfig := initTendermintConfig()

			return server.InterceptConfigsPreRunHandler(cmd, customTemplate, custompStakeConfig, customTMConfig)
		},
	}

	initRootCmd(rootCmd, encodingConfig)

	return rootCmd, encodingConfig
}

// initTendermintConfig helps to override default Tendermint Config values.
// return tmcfg.DefaultConfig if no custom configuration is required for the application.
func initTendermintConfig() *tmcfg.Config {
	cfg := tmcfg.DefaultConfig()

	// these values put a higher strain on node memory
	// cfg.P2P.MaxNumInboundPeers = 100
	// cfg.P2P.MaxNumOutboundPeers = 40

	return cfg
}
func initAppConfig() (string, interface{}) {
	srvCfg := serverconfig.DefaultConfig()
	srvCfg.StateSync.SnapshotInterval = 1000
	srvCfg.StateSync.SnapshotKeepRecent = 10
	srvCfg.MinGasPrices = "0stake,0uxprt"
	return params.CustomConfigTemplate, params.CustomAppConfig{
		Config: *srvCfg,
		BypassMinFeeMsgTypes: []string{
			sdk.MsgTypeURL(&ibcchanneltypes.MsgRecvPacket{}),
			sdk.MsgTypeURL(&ibcchanneltypes.MsgAcknowledgement{}),
			sdk.MsgTypeURL(&ibcclienttypes.MsgUpdateClient{}),
		},
	}
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig params.EncodingConfig) {
	cfg := sdk.GetConfig()

	cfg.Seal()

	rootCmd.AddCommand(
		genutilcli.InitCmd(pstakeApp.ModuleBasics, pstakeApp.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, pstakeApp.DefaultNodeHome),
		genutilcli.GenTxCmd(pstakeApp.ModuleBasics, encodingConfig.TxConfig, banktypes.GenesisBalancesIterator{}, pstakeApp.DefaultNodeHome),
		genutilcli.ValidateGenesisCmd(pstakeApp.ModuleBasics),
		AddGenesisAccountCmd(pstakeApp.DefaultNodeHome),
		tmcli.NewCompletionCmd(rootCmd, true),
		testnetCmd(pstakeApp.ModuleBasics, banktypes.GenesisBalancesIterator{}),
		debug.Cmd(),
		config.Cmd(),
	)

	ac := appCreator{
		encCfg: encodingConfig,
	}
	server.AddCommands(rootCmd, pstakeApp.DefaultNodeHome, ac.newApp, ac.appExport, addModuleInitFlags)

	// add keybase, auxiliary RPC, query, and tx child commands
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		queryCommand(),
		txCommand(),
		keys.Commands(pstakeApp.DefaultNodeHome),
	)
}

func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetAccountCmd(),
		rpc.ValidatorCommand(),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(),
		authcmd.QueryTxCmd(),
	)

	pstakeApp.ModuleBasics.AddQueryCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetMultiSignBatchCmd(),
		authcmd.GetValidateSignaturesCommand(),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
	)

	pstakeApp.ModuleBasics.AddTxCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

type appCreator struct {
	encCfg params.EncodingConfig
}

func (ac appCreator) newApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
) servertypes.Application {
	baseappOptions := server.DefaultBaseappOptions(appOpts)

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	return pstakeApp.NewpStakeApp(
		logger, db, traceStore, true, skipUpgradeHeights,
		cast.ToString(appOpts.Get(flags.FlagHome)),
		cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
		ac.encCfg,
		appOpts,
		baseappOptions...,
	)
}

func (ac appCreator) appExport(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
) (servertypes.ExportedApp, error) {
	var pStakeApp *pstakeApp.PstakeApp
	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, errors.New("application home not set")
	}

	if height != -1 {
		pStakeApp = pstakeApp.NewpStakeApp(logger, db, traceStore, false, map[int64]bool{}, homePath, uint(1), ac.encCfg, appOpts)

		if err := pStakeApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	} else {
		pStakeApp = pstakeApp.NewpStakeApp(logger, db, traceStore, true, map[int64]bool{}, homePath, uint(1), ac.encCfg, appOpts)
	}

	return pStakeApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs)

}
