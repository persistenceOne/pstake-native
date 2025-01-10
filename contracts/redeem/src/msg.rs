use crate::state::CfgData;
use cosmwasm_schema::{cw_serde, QueryResponses};
use cosmwasm_std::{Addr, Decimal};

#[cw_serde]
pub struct InstantiateMsg {
    pub admin: Addr,
}

#[cw_serde]
pub enum ExecuteMsg {
    AddConfig { base_asset: String, stk_asset: String, exchange_rate: Decimal },
    UpdateConfig { base_asset: String, stk_asset: String, exchange_rate: Decimal },
    Redeem {},
}

#[cw_serde]
#[derive(QueryResponses)]
pub enum QueryMsg {
    #[returns(Vec<ConfigResponse>)]
    Configs {},

    #[returns(CfgData)]
    Config { stk_asset: String },
}

#[cw_serde]
pub struct ConfigResponse {
    address: Addr,
    cfg: CfgData,
}