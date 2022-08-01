package cmd

import (
	stdlog "log"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/orchestrator/config"
	"github.com/persistenceOne/pstake-native/orchestrator/constants"
	"github.com/persistenceOne/pstake-native/orchestrator/orchestrator"
	tendermintService "github.com/tendermint/tendermint/libs/service"
)

func InitCosmosChain(homePath string, cfg config.CosmosConfig) (*orchestrator.CosmosChain, error) {
	chain := &orchestrator.CosmosChain{}
	chain.Key = "unusedKey"
	chain.ChainID = cfg.ChainID
	chain.GRPCAddr = cfg.GRPCAddr
	chain.RPCAddr = cfg.RPCAddr
	chain.AccountPrefix = cfg.AccountPrefix
	chain.GasAdjustment = cfg.GasAdjustment
	chain.GasPrices = cfg.GasPrice
	chain.CustodialAddress = sdk.AccAddress(cfg.CustodialAddr)
	chain.CoinType = cfg.CoinType

	err := chain.Init(string(chain.CustodialAddress), homePath, 1*time.Second, nil, true)
	if err != nil {
		return chain, err
	}
	if chain.KeyExists(chain.Key) {
		stdlog.Println("Key Exists for Cosmos Chain")
		err = chain.KeyBase.Delete(chain.Key)
		if err != nil {
			return chain, err
		}
	}

	_, err = config.KeyAddOrRestore(*chain, chain.Key, constants.CosmosCoinType)
	if err != nil {
		return chain, err
	}
	if err = chain.Start(); err != nil {
		if err != tendermintService.ErrAlreadyStarted {
			return chain, err
		}

	}
	return chain, nil
}
