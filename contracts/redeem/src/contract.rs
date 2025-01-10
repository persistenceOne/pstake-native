// use cw2::set_contract_version;
use crate::error::ContractError;
use crate::msg::{ExecuteMsg, InstantiateMsg, QueryMsg};
use crate::state::{CfgData, ADMIN, CONFIG};
#[cfg(not(feature = "library"))]
use cosmwasm_std::entry_point;
use cosmwasm_std::{coins, to_json_binary, Addr, BankMsg, Binary, CosmosMsg, Deps, DepsMut, Env, MessageInfo, Order, Response, StdError, StdResult};
/*
// version info for migration info
const CONTRACT_NAME: &str = "crates.io:redeem";
const CONTRACT_VERSION: &str = env!("CARGO_PKG_VERSION");
*/

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn instantiate(
    deps: DepsMut,
    _env: Env,
    _info: MessageInfo,
    msg: InstantiateMsg,
) -> Result<Response, ContractError> {
    ADMIN.save(deps.storage, &msg.admin)?;
    Ok(Response::default())
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn execute(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response, ContractError> {
    match msg {
        ExecuteMsg::AddConfig { base_asset, stk_asset, exchange_rate } => {
            let admin = ADMIN.load(deps.storage)?;
            if info.sender != admin {
                return Err(ContractError::Unauthorized {});
            }
            let cfg_data: CfgData = CfgData { base_asset, exchange_rate };
            CONFIG.save(deps.storage, &stk_asset, &cfg_data)?;
            Ok(Response::default())
        }

        ExecuteMsg::UpdateConfig { base_asset, stk_asset, exchange_rate } => {
            let admin = ADMIN.load(deps.storage)?;
            if info.sender != admin {
                return Err(ContractError::Unauthorized {});
            }
            CONFIG.update(deps.storage, &stk_asset, |current_value| -> StdResult<CfgData> {
                let data = current_value.unwrap_or(CfgData {
                    base_asset,
                    exchange_rate,
                });

                Ok(data)
            })?;
            Ok(Response::default())
        }
        ExecuteMsg::Redeem {} => {
            if info.funds.len() != 1 {
                return Err(ContractError::Std(StdError::generic_err("only one coin allowed")));
            }
            let coin = info.clone().funds.pop().unwrap();
            let cfg = CONFIG.load(deps.storage, &coin.denom)?;
            let amt = coin.amount.checked_div_floor(cfg.exchange_rate).unwrap();
            let amt_to_send = amt.clone().u128();
            let send_msg = CosmosMsg::Bank(BankMsg::Send {
                to_address: info.sender.into(),
                amount: coins(amt_to_send, cfg.base_asset),
            });
            let burn_msg = CosmosMsg::Bank(BankMsg::Burn {
                amount: info.funds
            });
            let res = Response::new()
                .add_messages(vec![burn_msg, send_msg])
                .add_attribute("action", "redeem_stk");
            Ok(res)
        }
    }
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn query(deps: Deps, _env: Env, msg: QueryMsg) -> StdResult<Binary> {
    match msg {
        QueryMsg::Config { stk_asset } => {
            to_json_binary(&CONFIG.load(deps.storage, &stk_asset)?)
        }
        QueryMsg::Configs {} => {
            let configs: Vec<(Addr, CfgData)> = CONFIG
                .range(deps.storage, None, None, Order::Ascending)
                .map(|item| {
                    item.map(|(key, value)| {
                        let addr = Addr::unchecked(key); // Convert String to Addr
                        (addr, value)
                    })
                })
                .collect::<StdResult<Vec<(Addr, CfgData)>>>()?;

            to_json_binary(&configs)
        }
    }
}

#[cfg(test)]
mod tests {
    use crate::contract;
    use crate::msg::QueryMsg;
    use crate::state::CfgData;
    use cosmwasm_std::{coins, Addr, Decimal};
    use cw_multi_test::{App, ContractWrapper, Executor};
    use std::fmt::{Debug, Pointer};
    use std::str::FromStr;

    #[test]
    fn test_allow_lp_token() {
        let admin = String::from("admin");
        let coins1 = vec![
            coins(1000_000_000, "uxprt")[0].clone(),
            coins(1000_000_000, "stk/uatom")[0].clone(),
            coins(1000_000_000, "ibc/uatom")[0].clone(),
        ];
        let admin_addr = Addr::unchecked(admin.clone());
        let mut app = App::new(|router, _, storage| {
            // initialization moved to App construction
            router.bank.init_balance(storage, &admin_addr, coins1).unwrap()
        });

        let redeem_contract = Box::new(
            ContractWrapper::new_with_empty(
                contract::execute,
                contract::instantiate,
                contract::query,
            )
        );
        let code_id = app.store_code(redeem_contract);

        let msg = contract::InstantiateMsg {
            admin: admin_addr.clone(),

        };
        let contract_i = app
            .instantiate_contract(
                code_id,
                admin_addr.clone(),
                &msg,
                &[],
                String::from("C1"),
                None,
            )
            .unwrap();

        let ex = Decimal::from_str("0.5").unwrap();
        let res = app.execute_contract(
            admin_addr.clone(), contract_i.clone(), &contract::ExecuteMsg::AddConfig {
                base_asset: "ibc/uatom".to_string(),
                stk_asset: "stk/uatom".to_string(),
                exchange_rate: ex,
            }, &vec![
                coins(1000_000_00, "ibc/uatom")[0].clone(),
            ],
        ).unwrap();
        println!("{:?}", res);

        let q1: Vec<(Addr, CfgData)> = app.wrap().query_wasm_smart(&contract_i.clone(), { &QueryMsg::Configs {} }).unwrap();
        println!("{:?}", q1);

        let res1 = app.execute_contract(
            admin_addr.clone(), contract_i.clone(), &contract::ExecuteMsg::Redeem {},
            &vec![
                coins(1000_000, "stk/uatom")[0].clone(),
            ],
        ).unwrap();
        println!("{:?}", res1);
    }
}
