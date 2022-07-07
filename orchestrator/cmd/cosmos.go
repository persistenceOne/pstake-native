package cmd

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/orchestrator/config"
	"github.com/persistenceOne/pstake-native/orchestrator/constants"
	"github.com/persistenceOne/pstake-native/orchestrator/helpers"
	"github.com/persistenceOne/pstake-native/orchestrator/orchestrator"
	tendermintService "github.com/tendermint/tendermint/libs/service"
	stdlog "log"
	"time"
)

func InitCosmosChain(homePath string, config config.CosmosConfig) (*orchestrator.CosmosChain, error) {
	chain := &orchestrator.CosmosChain{}
	chain.Key = "unusedKey"
	chain.ChainID = config.ChainID
	chain.GRPCAddr = config.GRPCAddr
	chain.RPCAddr = config.RPCAddr
	chain.AccountPrefix = config.AccountPrefix
	chain.GasAdjustment = config.GasAdjustment
	chain.GasPrices = config.GasPrice
	chain.CustodialAddress = sdk.AccAddress(config.CustodialAddr)
	chain.CoinType = config.CoinType

	err := chain.Init(string(chain.CustodialAddress), homePath, 1*time.Second, nil, true)
	if err != nil {
		return chain, err
	}
	if chain.KeyExists(chain.Key) {
		stdlog.Println("Key Exists")
		err = chain.KeyBase.Delete(chain.Key)
		if err != nil {
			return chain, err
		}
	}

	_, err = helpers.KeyAddOrRestore(*chain, chain.Key, constants.CosmosCoinType)
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
