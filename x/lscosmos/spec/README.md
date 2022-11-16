<!--
order: 1
title: "Ls-cosmos overview"
parent:
  title: "ls-cosmos"
-->

# `ls-cosmos`

## Abstarct

ls-cosmos is liquid staking module for the cosmos chain. The purpose of the module is to provide functionality to be
able to liquid-stake ATOM tokens and get representative stkATOM tokens in return.

## Contents

0. **[Introduction](00_introduction.md)
1. **[Concept](01_concepts.md)**
2. **[Proposal](02_proposal.md)
   - [Register Host Chain Proposal](02_proposal.md#register-host-chain-proposal)
   - [Change Min Deposit and Fee Proposal](02_proposal.md#change-min-deposit-and-fee-proposal)
   - [Change Pstake Fee Address Proposal](02_proposal.md#change-pstake-fee-address-proposal)
   - [Change Allow Listed Validators Proposal](02_proposal.md#change-allow-listed-validators-proposal)
3. **[State](03_state.md)**
4. **[Events](04_events.md)**
   - [MsgLiquidStake](04_events.md#msgliquidstake)
   - [MsgJuice](04_events.md#msgjuice)
   - [MsgLiquidUnstake](04_events.md#msgliquidunstake)
   - [MsgRedeem](04_events.md#msgredeem)
   - [MsgClaim](04_events.md#msgclaim)
   - [MsgJumpStart](04_events.md#msgjumpstart)
   - [MsgRecreateICA](04_events.md#msgrecreateica)
5. **[Keeper](05_keeper.md)**
      [KeeperFunctions](05_keeper.md#keeper-functions)
6. **[Messages](06_messages.md)**
    - [MsgLiquidStake](06_messages.md#msgliquidstake)
    - [MsgJuice](06_messages.md#msgjuice)
    - [MsgLiquidUnstake](06_messages.md#msgliquidunstake)
    - [MsgRedeem](06_messages.md#msgredeem)
    - [MsgClaim](06_messages.md#msgclaim)
    - [MsgJumpStart](06_messages.md#msgjumpstart)
    - [MsgRecreateICA](06_messages.md#msgrecreateica)
7. **[Queries](07_queries.md)**
8. **[Future improvements](08_future_improvements.md)**