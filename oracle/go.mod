module github.com/persistenceOne/pstake-native/oracle

go 1.16

require (
	github.com/BurntSushi/toml v0.4.1
	github.com/cosmos/cosmos-sdk v0.44.5
	github.com/cosmos/relayer v1.0.0
	github.com/gogo/protobuf v1.3.3
	github.com/persistenceOne/pstake-native v0.0.0-20220601081531-903d8e733ccd
	github.com/spf13/cobra v1.3.0
	github.com/tendermint/tendermint v0.34.15
	google.golang.org/grpc v1.42.0
)

require (
	github.com/cosmos/ibc-go/v2 v2.0.3 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
replace (
	github.com/persistenceOne/pstake-native v0.0.0-20220601081531-903d8e733ccd => ../../pStake-native
)
