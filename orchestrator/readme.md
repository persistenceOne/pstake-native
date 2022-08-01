for building: 
`````go build orc.go`````

initialise command:
`````./orc init --orcHome "homepath"`````

start command:
`````./orc start --orcHome "homepath"`````

dummy config.toml:

````    
ValAddress = "csomssssssssssssssssssssssssss"
OrcSeeds = ["axis decline final suggest denial erupt satisfy weekend utility fo$

[CosmosConfig]
    ChainID = "test"
    CustodialAddr = "cosmos15vm0p2x990762txvsrpr26ya54p5qlz9xqlw5z"
    Denom = "stake"
    GRPCAddr = "127.0.0.1:9090"
    RPCAddr = "127.0.0.1:26657"
    AccountPrefix = "cosmos"
    GasAdjustment = 1.0
    GasPrice = "0.025stake"
    CoinType = 118                                  
    
[NativeConfig]
    ChainID = "test1"
    Denom = "stake"
    RPCAddr = "127.0.0.1:9090"
    GRPCAddr = "127.0.0.1:26657"
    AccountPrefix = "cosmos"
    GasAdjustment = 1.0
    GasPrices = "0.025stake"
    CoinType = 118

````