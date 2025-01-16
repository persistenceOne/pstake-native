use cosmwasm_schema::cw_serde;
use cosmwasm_std::{Addr, Decimal};
use cw_storage_plus::{Item, Map};

#[cw_serde]
pub struct CfgData {
    pub base_asset: String,
    pub exchange_rate: Decimal,
}

pub const CONFIG: Map<&str, CfgData> = Map::new("config");
pub const ADMIN: Item<Addr> = Item::new("admin");