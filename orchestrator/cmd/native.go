package cmd

import (
	"log"
	"time"

	"github.com/persistenceOne/pstake-native/orchestrator/config"
	"github.com/persistenceOne/pstake-native/orchestrator/constants"
	"github.com/persistenceOne/pstake-native/orchestrator/orchestrator"
	tendermintService "github.com/tendermint/tendermint/libs/service"
)

func InitNativeChain(homePath string, cfg config.NativeConfig) (*orchestrator.NativeChain, error) {
	chain := &orchestrator.NativeChain{}
	chain.Key = "unusedNativeKey"
	chain.ChainID = cfg.ChainID
	chain.RPCAddr = cfg.RPCAddr
	chain.GRPCAddr = cfg.GRPCAddr
	chain.AccountPrefix = cfg.AccountPrefix
	chain.GasAdjustment = cfg.GasAdjustment
	chain.GasPrices = cfg.GasPrices
	chain.CoinType = cfg.CoinType

	err := chain.Init(homePath, 1*time.Second, nil, true)
	if err != nil {
		return chain, err
	}
	if chain.KeyExists(chain.Key) {
		log.Println("Key Exists for native chain")
		err = chain.KeyBase.Delete(chain.Key)
		if err != nil {
			return chain, err
		}
	}

	_, err = config.KeyAddOrRestoreNative(*chain, chain.Key, constants.NativeCoinType)
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
