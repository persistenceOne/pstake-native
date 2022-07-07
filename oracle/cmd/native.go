package orc

import (
	"github.com/persistenceOne/pstake-native/oracle/config"
	"github.com/persistenceOne/pstake-native/oracle/constants"
	"github.com/persistenceOne/pstake-native/oracle/helpers"
	"github.com/persistenceOne/pstake-native/oracle/orchestrator"
	tendermintService "github.com/tendermint/tendermint/libs/service"
	"log"
	"time"
)

func InitNativeChain(homePath string, config config.NativeConfig) (*orchestrator.NativeChain, error) {
	chain := &orchestrator.NativeChain{}
	chain.Key = "unusedNativeKey"
	chain.ChainID = config.ChainID
	chain.RPCAddr = config.RPCAddr
	chain.GRPCAddr = config.GRPCAddr
	chain.AccountPrefix = config.AccountPrefix
	chain.GasAdjustment = config.GasAdjustment
	chain.GasPrices = config.GasPrices
	chain.CoinType = config.CoinType

	err := chain.Init(homePath, 1*time.Second, nil, true)
	if err != nil {
		return chain, err
	}
	if chain.KeyExists(chain.Key) {
		log.Println("Key Exists")
		err = chain.KeyBase.Delete(chain.Key)
		if err != nil {
			return chain, err
		}
	}

	_, err = helpers.KeyAddOrRestoreNative(*chain, chain.Key, constants.NativeCoinType)
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
