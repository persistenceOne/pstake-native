package orc

import (
	"fmt"
	"github.com/persistenceOne/pStake-native/oracle/configuration"
	"github.com/persistenceOne/pStake-native/oracle/constants"
	"github.com/persistenceOne/pStake-native/oracle/helpers"
	"github.com/persistenceOne/pStake-native/oracle/oracle"
	tendermintService "github.com/tendermint/tendermint/libs/service"
	"time"
)

func InitNativeChain(homePath string) (*oracle.NativeChain, error) {
	chain := &oracle.NativeChain{}
	chain.Key = "unusedNativeKey"
	chain.ChainID = configuration.GetConfig().NativeConfig.ChainID
	chain.RPCAddr = configuration.GetConfig().NativeConfig.RPCAddr
	chain.AccountPrefix = configuration.GetConfig().NativeConfig.AccountPrefix
	chain.GasAdjustment = configuration.GetConfig().NativeConfig.GasAdjustment
	chain.GasPrices = configuration.GetConfig().NativeConfig.GasPrices
	//= sdk.AccAddress(configuration.GetConfig().CosmosConfig.CustodialAddr)

	err := chain.Init(homePath, 1*time.Second, nil, true)
	if err != nil {
		return chain, err
	}
	if chain.KeyExists(chain.Key) {
		fmt.Println("Key Exists")
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
