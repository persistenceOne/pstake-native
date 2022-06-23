package orc

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/oracle/configuration"
	"github.com/persistenceOne/pstake-native/oracle/constants"
	"github.com/persistenceOne/pstake-native/oracle/helpers"
	"github.com/persistenceOne/pstake-native/oracle/oracle"
	tendermintService "github.com/tendermint/tendermint/libs/service"
	"time"
)

func InitCosmosChain(homePath string, config configuration.CosmosConfig) (*oracle.CosmosChain, error) {
	chain := &oracle.CosmosChain{}
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
		fmt.Println("Key Exists")
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
