package configs

import (
	"os"

	lensclient "github.com/strangelove-ventures/lens/client"
)

func GetPersistenceClient(home, keyhome string) (*lensclient.ChainClient, error) {
	ccc := &lensclient.ChainClientConfig{
		Key:     "default",
		ChainID: "core-1",
		//RPCAddr:        "https://rpc.persistence.audit.one:443",
		RPCAddr:        "https://rpc.core.persistence.one:443",
		GRPCAddr:       "https://grpc.persistence.audit.one:443",
		AccountPrefix:  "persistence",
		KeyringBackend: "test",
		GasAdjustment:  1.2,
		GasPrices:      "0.025uxprt",
		KeyDirectory:   keyhome,
		Debug:          true,
		Timeout:        "20s",
		BlockTimeout:   "10s",
		OutputFormat:   "json",
		SignModeStr:    "direct",
		MinGasAmount:   0,
		ExtraCodecs:    nil,
		Modules:        ModuleBasics,
		Slip44:         0,
	}
	return lensclient.NewChainClient(nil, ccc, home, os.Stdin, os.Stdout)

}
func GetCosmosClient(home, keyhome string) (*lensclient.ChainClient, error) {
	ccc := &lensclient.ChainClientConfig{
		Key:            "default",
		ChainID:        "cosmoshub-4",
		RPCAddr:        "https://rpc.cosmos.audit.one:443",
		GRPCAddr:       "grpc.cosmos.dragonstake.io:443",
		AccountPrefix:  "cosmos",
		KeyringBackend: "test",
		GasAdjustment:  1.2,
		GasPrices:      "0.0001uatom",
		KeyDirectory:   keyhome,
		Debug:          true,
		Timeout:        "20s",
		BlockTimeout:   "10s",
		OutputFormat:   "json",
		SignModeStr:    "direct",
		MinGasAmount:   0,
		ExtraCodecs:    nil,
		Modules:        ModuleBasics,
		Slip44:         0,
	}
	return lensclient.NewChainClient(nil, ccc, home, os.Stdin, os.Stdout)

}
