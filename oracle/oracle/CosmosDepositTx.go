package oracle

import github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"

type CosmoDepositTx struct {
	AddressFromMemo     string                                   `protobuf:"bytes,1,opt,name=address_from_memo,json=addressFromMemo,proto3" json:"address_from_memo,omitempty" yaml:"address_from_memo"`
	OrchestratorAddress string                                   `protobuf:"bytes,2,opt,name=orchestrator_address,json=orchestratorAddress,proto3" json:"orchestrator_address,omitempty" yaml:"orchestrator_address"`
	Amount              github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,3,rep,name=amount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Coin,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"amount"`
	TxHash              string                                   `protobuf:"bytes,4,opt,name=tx_hash,json=txHash,proto3" json:"tx_hash,omitempty" yaml:"tx_hash"`
	ChainID             string                                   `protobuf:"bytes,5,opt,name=chain_i_d,json=chainID,proto3" json:"chain_i_d,omitempty" yaml:"chain_id"`
	BlockHeight         int64                                    `protobuf:"varint,6,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty" yaml:"block_height"`
}

//
//func (c CosmoSDepositTx) Reset() {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (c CosmoSDepositTx) String() string {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (c CosmoSDepositTx) ProtoMessage() {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (c CosmoSDepositTx) ValidateBasic() error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (c CosmoSDepositTx) GetSigners() []github_com_cosmos_cosmos_sdk_types.AccAddress {
//	//TODO implement me
//	panic("implement me")
//}
//
//var _ github_com_cosmos_cosmos_sdk_types.Msg = &CosmoSDepositTx{}
