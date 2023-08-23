use cosmwasm_schema::{cw_serde, QueryResponses};
use cosmwasm_std::{Coin, CustomMsg};
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

#[cw_serde]
pub struct InstantiateMsg {
    pub count: i32,
}

#[cw_serde]
pub enum ExecuteMsg {
    Increment {},
    Reset { count: i32 },
    LiquidStake { receiver: String },
}

// #[cw_serde]
#[non_exhaustive]
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum PstakeMsg {
    MsgLiquidStake {
        delegator_address: String,
        amount: Coin,
    }
}
impl CustomMsg for PstakeMsg{}


// /// This is just a demo place so we can test custom message handling
// #[derive(Debug, Clone, Serialize, Deserialize, JsonSchema, PartialEq)]
// #[serde(rename = "snake_case")]
// pub enum CustomMsg {
//     MsgLiquidStake {
//         delegator_address: String,
//         amount: Coin,
//     }
// }


#[cw_serde]
#[derive(QueryResponses)]
pub enum QueryMsg {
    // GetCount returns the current count as a json-encoded number
    #[returns(GetCountResponse)]
    GetCount {},
}

// We define a custom struct for each query response
#[cw_serde]
pub struct GetCountResponse {
    pub count: i32,
}
