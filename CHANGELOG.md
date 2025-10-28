<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry should ideally include a tag and
the Github issue reference in the following format:

* (<tag>) \#<issue-number> message

The issue numbers will later be link-ified during the release process so you do
not have to worry about including a link manually, but you can if you wish.


"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Removed" for now removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking CLI commands and REST routes.
"State Machine Breaking" for breaking the AppState

Ref: https://keepachangelog.com/en/1.0.0/
-->

# Changelog

## [Unreleased]
### Features


## [v6.0.0]

### Features

- update-cosmos-sdk to 0.53.x, use cosmos-sdk/x/epochs ([#1034](https://github.com/persistenceOne/pstake-native/pull/1034))

## [v5.1.0]

### Chores

- remove cosmos-sdk fork usage ([#1023](https://github.com/persistenceOne/pstake-native/pull/1023))


## [v5.0.0]

### Refactoring

- replace sdk-lsm to use gaia/liquid ([#1017](https://github.com/persistenceOne/pstake-native/pull/1017))
- stop rebalancing as it introduces redelegations ([#1020](https://github.com/persistenceOne/pstake-native/pull/1020))

## [v4.0.0]

### Features

- **ci**: use ci to autogenerate changelog ([#974](https://github.com/persistenceOne/pstake-native/pull/974))
- update deps to sdk v0.50.x-lsm ([#988](https://github.com/persistenceOne/pstake-native/pull/988))

### Bug Fixes

- **ci**: update test-coverage ([#975](https://github.com/persistenceOne/pstake-native/pull/975))

### Refactoring

- move modules interfaces/codecs encodings ([#981](https://github.com/persistenceOne/pstake-native/pull/981))
- remove module usage for deprecated modules ([#982](https://github.com/persistenceOne/pstake-native/pull/982))
- remove ibc-go direct dep ([#998](https://github.com/persistenceOne/pstake-native/pull/998))

## [v3.0.0]

### Deleted 

- [958](https://github.com/persistenceOne/pstake-native/pull/958) Delete liquidstakeibc, ratesync module code.

## [v2.16.0] - 2024-12-24

### Improvements

- [855](https://github.com/persistenceOne/pstake-native/pull/855) Add condition for not allowing zero delegation unbondings icq
- [900](https://github.com/persistenceOne/pstake-native/pull/900) Add feature to deprecate liquidstakeibc

## [v2.15.0] - 2024-05-30

### Improvements

- [831](https://github.com/persistenceOne/pstake-native/pull/831) Add amino tags for protobuf msgs for compiling in js
  using telescope
- [841](https://github.com/persistenceOne/pstake-native/pull/841) add type ForceUpdateValidatorDelegations to
  MsgUpdateHostChain
- [842](https://github.com/persistenceOne/pstake-native/pull/842) liquidstake: move rebalancing from begin block to day
  epoch

## [v2.13.0] - 2024-05-01

### Bug Fixes

- [815](https://github.com/persistenceOne/pstake-native/pull/815) Not escape merkle paths for proof verification.
- [822](https://github.com/persistenceOne/pstake-native/pull/822) liquidstake: redelegation to follow msg router instead
  of keeper.

## [v2.12.0] - 2024-04-05

### Improvements

- [800](https://github.com/persistenceOne/pstake-native/pull/800) Remove proxy account usage.
- [790](https://github.com/persistenceOne/pstake-native/pull/790) Make liquidstake module LSM Cap compliant.

### Bug Fixes

- [792](https://github.com/persistenceOne/pstake-native/pull/792) Use GetHostChainFromHostDenom in ICA Transfer
  unsuccessfulAck instead of GetHostChainFromDelegatorAddress as Rewards account too uses ICA Transfer to autocompound
- [795](https://github.com/persistenceOne/pstake-native/pull/795) Reject zero weight validator LSM shares for
  liquidstakeibc

## [v2.11.0] - 2024-03-12

### Features

- [783](https://github.com/persistenceOne/pstake-native/pull/783) Move rewards autocompounding to hourly epoch.
- [784](https://github.com/persistenceOne/pstake-native/pull/784) Allow admin address to update params.

### Improvements

- [773](https://github.com/persistenceOne/pstake-native/pull/773) Improve logging.

### Bug Fixes

- [774](https://github.com/persistenceOne/pstake-native/pull/774) Fix liquidstake params test.

## [v2.10.0] - 2024-02-21

### Features

- [760](https://github.com/persistenceOne/pstake-native/pull/760) Calculate C Value after autocompounding / slashing.
- [758](https://github.com/persistenceOne/pstake-native/pull/758) Dynamic C Value limit updates.
- [757](https://github.com/persistenceOne/pstake-native/pull/757) Change ibc transfer to use timeoutTimestamp instead of
  timeoutHeight
- [756](https://github.com/persistenceOne/pstake-native/pull/756) Add query for singular host-chain in liquidstakeibc
- [755](https://github.com/persistenceOne/pstake-native/pull/755) Add channel-id, port to ratesync host-chains and
  liquidstake instantiate.

### Bug Fixes

- [766](https://github.com/persistenceOne/pstake-native/pull/766) stkxprt audit fixes
- [761](https://github.com/persistenceOne/pstake-native/pull/761) Use counterparty channels instead of self chain for
  ratesync-instantiate

## [v2.9.1] - 2024-01-26

### Bug Fixes

- [753](https://github.com/persistenceOne/pstake-native/pull/753) Set default bounds for c value

## [v2.9.0] - 2024-01-26

### Features

- [737](https://github.com/persistenceOne/pstake-native/pull/737) Unhandled errors.
- [736](https://github.com/persistenceOne/pstake-native/pull/736) Check for host denom duplicates.
- [733](https://github.com/persistenceOne/pstake-native/pull/733) Add more validation for host-chain.
- [732](https://github.com/persistenceOne/pstake-native/pull/732) Move c value bounds to per-chain params.
- [729](https://github.com/persistenceOne/pstake-native/pull/729) Add rewards account query (hence autocompound)
  OnChanOpenAck.
- [727](https://github.com/persistenceOne/pstake-native/pull/727) Send LSM redeem messages in chunks.
- [721](https://github.com/persistenceOne/pstake-native/pull/721) Add Query host chain user unbondings.

### Bug Fixes

- [752](https://github.com/persistenceOne/pstake-native/pull/752) Use correct existing delegation amount.
- [751](https://github.com/persistenceOne/pstake-native/pull/751) Set LSM bond factor as -1 by default.
- [750](https://github.com/persistenceOne/pstake-native/pull/750) Shares to tokens.
- [734](https://github.com/persistenceOne/pstake-native/pull/734) Host chain duplication check.
- [731](https://github.com/persistenceOne/pstake-native/pull/731) Set limit to LSM deposit filtering.
- [730](https://github.com/persistenceOne/pstake-native/pull/730) Fix deposit validate.
- [728](https://github.com/persistenceOne/pstake-native/pull/728) Fix prevent users from liquid-staking funds by
  removing the Deposit entry.
- [726](https://github.com/persistenceOne/pstake-native/pull/726) Fix minimal unbondings.
- [725](https://github.com/persistenceOne/pstake-native/pull/725) Fix Incorrect bookkeeping of validator’s delegated
  amount upon redelegation
- [720](https://github.com/persistenceOne/pstake-native/pull/720) Fix unbondings loop.
- [719](https://github.com/persistenceOne/pstake-native/pull/719) Fix afterEpoch hooks to take LiquidStake feature
  instead of LiquidStakeIBC

## [v2.8.2] - 2024-01-09

### Bug Fixes

- [715](https://github.com/persistenceOne/pstake-native/pull/715) Fix stuck unbondings.

## [v2.8.1] - 2023-12-21

### Bug Fixes

- [707](https://github.com/persistenceOne/pstake-native/pull/707) Fix liquidstakeibc redeem edge case for protecting
  cValue

## [v2.8.0] - 2023-12-20

### Features

- [703](https://github.com/persistenceOne/pstake-native/pull/703) Add ratesync module.

## [v2.7.x] - 2023-12-15

### Features

- [687](https://github.com/persistenceOne/pstake-native/pull/687) Add queries for redelegations and redelegation txs.
- [696](https://github.com/persistenceOne/pstake-native/pull/696) Add capability to swap rewards.
- [697](https://github.com/persistenceOne/pstake-native/pull/697) Add hooks for liquidstakeibc c_value updates.

### Bug Fixes

- [685](https://github.com/persistenceOne/pstake-native/pull/685) Fix rebalancing to happen outside of unbonding epochs

## [v2.6.0] - 2023-11-26

### Features

- [#680](https://github.com/persistenceOne/pstake-native/pull/680) Add rebalancing

## [v2.5.0] - 2023-10-20

### Features

- [#667](https://github.com/persistenceOne/pstake-native/pull/667) Monitoring Events.

### Improvements

- [#668](https://github.com/persistenceOne/pstake-native/pull/668) Update ICA timeout.

### Bug Fixes

- [#665](https://github.com/persistenceOne/pstake-native/pull/665) LSM deposit timeout fix.

## [v2.4.0] - 2023-09-13

### Bug Fixes

- [#652](https://github.com/persistenceOne/pstake-native/pull/652) Register MsgLiquidStakeLSM into amino codec.

## [v2.3.3] - 2023-09-07

### Bug Fixes

- [#643](https://github.com/persistenceOne/pstake-native/pull/643) Re-add LSCosmos types to enable parsing of older
  transactions and gov data.

## [v2.3.2] - 2023-09-07

### Bug Fixes

- [#639](https://github.com/persistenceOne/pstake-native/pull/639) LSM cap room bug.

## [v2.3.1] - 2023-09-06

### Bug Fixes

- [#637](https://github.com/persistenceOne/pstake-native/pull/637) LSM bond factor validation fix.

## [v2.3.0] - 2023-09-06

### Features

- [#594](https://github.com/persistenceOne/pstake-native/pull/594) LSM integration.
- [#631](https://github.com/persistenceOne/pstake-native/pull/631) Add telemetry to measure time taken in begin/end
  block.

### Bug Fixes

- [#632](https://github.com/persistenceOne/pstake-native/pull/632) LSM cap fix.
- [#621](https://github.com/persistenceOne/pstake-native/pull/621) ICA recreation timeout fix.

### Removed

- [#627](https://github.com/persistenceOne/pstake-native/pull/627) Remove lscosmos, lspersistence in favour of
  liquidstakeibc.

## [v2.2.2] - 2023-08-07

### Improvements

- [#604](https://github.com/persistenceOne/pstake-native/pull/604) Create an empty deposit when a host chain is
  registered.

## [v2.2.1] - 2023-07-27

### Improvements

- [#603](https://github.com/persistenceOne/pstake-native/pull/603) Add extra code for lscosmos migration to account for
  unrelayed packets.

## [v2.2.0] - 2023-07-24

### Features

- [#560](https://github.com/persistenceOne/pstake-native/pull/560) liquidstakeibc: allow localhost client type.
- [#524](https://github.com/persistenceOne/pstake-native/pull/524) [LiquidStakeIbc] Query Updates
- [#474](https://github.com/persistenceOne/pstake-native/pull/474) Liquidstakeibc Unstake
- [#475](https://github.com/persistenceOne/pstake-native/pull/475) LiquidstakeIbc Claim
- [#479](https://github.com/persistenceOne/pstake-native/pull/479) Liquidstakeibc - Redeem
- [#481](https://github.com/persistenceOne/pstake-native/pull/481) Liquidstakeibc - Automatically claim failed
  unbondings
- [#482](https://github.com/persistenceOne/pstake-native/pull/482) LiquidstakeIbc - Autocompounding
- [#485](https://github.com/persistenceOne/pstake-native/pull/485) Liquidstakeibc - host chain activation
- [#486](https://github.com/persistenceOne/pstake-native/pull/486) Liquidstakeibc - Slashing
- [#488](https://github.com/persistenceOne/pstake-native/pull/488) Liquidstakeibc - update unbonding queries
- [#487](https://github.com/persistenceOne/pstake-native/pull/487) Liquidstakeibc - Withdraw delegator rewards
- [#513](https://github.com/persistenceOne/pstake-native/pull/513) [LiquidStakeIbc] ICQ Proofs
- [#511](https://github.com/persistenceOne/pstake-native/pull/511) use params store in lspersistence instead of params
  module
- [#478](https://github.com/persistenceOne/pstake-native/pull/478) stkxprt: add fees for staking
- [#483](https://github.com/persistenceOne/pstake-native/pull/483) stkxprt: add fees for unstake.
- [#480](https://github.com/persistenceOne/pstake-native/pull/480) stkxprt: add fees for restake.

- [#514](https://github.com/persistenceOne/pstake-native/pull/514) add proto checks to github actions

### Improvements

- [#588](https://github.com/persistenceOne/pstake-native/pull/588) Cleanup repository
- [#492](https://github.com/persistenceOne/pstake-native/pull/492),[#558](https://github.com/persistenceOne/pstake-native/pull/558),[#561](https://github.com/persistenceOne/pstake-native/pull/561),[#563](https://github.com/persistenceOne/pstake-native/pull/563),[#564](https://github.com/persistenceOne/pstake-native/pull/564),[#565](https://github.com/persistenceOne/pstake-native/pull/565),[#566](https://github.com/persistenceOne/pstake-native/pull/566),[#567](https://github.com/persistenceOne/pstake-native/pull/567),[#568](https://github.com/persistenceOne/pstake-native/pull/568),[#569](https://github.com/persistenceOne/pstake-native/pull/569),[#570](https://github.com/persistenceOne/pstake-native/pull/570),[#571](https://github.com/persistenceOne/pstake-native/pull/571),[#573](https://github.com/persistenceOne/pstake-native/pull/573),[#574](https://github.com/persistenceOne/pstake-native/pull/574),[#576](https://github.com/persistenceOne/pstake-native/pull/576),[#577](https://github.com/persistenceOne/pstake-native/pull/577),[#580](https://github.com/persistenceOne/pstake-native/pull/580),[#581](https://github.com/persistenceOne/pstake-native/pull/581),[#582](https://github.com/persistenceOne/pstake-native/pull/582),[#585](https://github.com/persistenceOne/pstake-native/pull/585)
  unit tests for liquidstake ibc
- [#587](https://github.com/persistenceOne/pstake-native/pull/587) [LiquidStakeIbc] e2e test setup
- [#522](https://github.com/persistenceOne/pstake-native/pull/522) [LiquidStakeIbc] Cleanups
- [#501](https://github.com/persistenceOne/pstake-native/pull/501) liquidstakeibc: use icaaccount.owner instead of
  keeper func,
- [#518](https://github.com/persistenceOne/pstake-native/pull/518),[#528](https://github.com/persistenceOne/pstake-native/pull/528),[#491](https://github.com/persistenceOne/pstake-native/pull/491),[#536](https://github.com/persistenceOne/pstake-native/pull/536),[#508](https://github.com/persistenceOne/pstake-native/pull/508),
  add migrations.
- [#459](https://github.com/persistenceOne/pstake-native/pull/459) LS module refactor - Part 1
- [#493](https://github.com/persistenceOne/pstake-native/pull/493) refactor queries for lscosmos for proofs
- [#539](https://github.com/persistenceOne/pstake-native/pull/539) ci-separate module tests with e2e.

### Bug Fixes

- [#556](https://github.com/persistenceOne/pstake-native/pull/556) LiquidStakeIbc notional Audit
- [#540](https://github.com/persistenceOne/pstake-native/pull/540) [LiquidStakeIbc] Fix Auto Slashing Mechanism
- [#542](https://github.com/persistenceOne/pstake-native/pull/542) [LiquidstakeIbc] Limit Auto Compounding
- [#533](https://github.com/persistenceOne/pstake-native/pull/533) [Liquidstakeibc] remove deposits when chain is
  disabled
- [#535](https://github.com/persistenceOne/pstake-native/pull/535) [Liquidstakeibc] add init validator delegation to
  deal with updates to validator set.
- [#499](https://github.com/persistenceOne/pstake-native/pull/499) [LiquidStakeIbc] Fixes
- [#512](https://github.com/persistenceOne/pstake-native/pull/512) fix icq unmarshal for icq unmarshal.

### Removed

- [#516](https://github.com/persistenceOne/pstake-native/pull/516),[#517](https://github.com/persistenceOne/pstake-native/pull/517),[#515](https://github.com/persistenceOne/pstake-native/pull/515)
  lscosmos: deprecate, remove functional code

## [v2.1.0-rc0] - 2023-04-20

Never released.

### Improvements

- [#411](https://github.com/persistenceOne/pstake-native/pull/411) add admin functionality to disable module in case of
  failure.
- [#410](https://github.com/persistenceOne/pstake-native/pull/410) reset IBC state instead of retrying IBC.
- [#422](https://github.com/persistenceOne/pstake-native/pull/422) sdkv46
- [#440](https://github.com/persistenceOne/pstake-native/pull/440) implement slashing handling
- [#446](https://github.com/persistenceOne/pstake-native/pull/446) add ibcfee
- [#447](https://github.com/persistenceOne/pstake-native/pull/447) remove register_host_chain_proposal.json.
- [#451](https://github.com/persistenceOne/pstake-native/pull/451) Update protogen using buf

### Bug Fixes

- [#405](https://github.com/persistenceOne/pstake-native/pull/405) audit: inconsistent state fixes
- [#413](https://github.com/persistenceOne/pstake-native/pull/413) fix resetICA like #405

### Removed

- [#409](https://github.com/persistenceOne/pstake-native/pull/409) remove governance proposal

## [v2.0.1, v2.0.2] - 2023-06-29

### Bug Fixes

- [#545](https://github.com/persistenceOne/pstake-native/pull/545) fix failed unbondings on mainnet
- [#549](https://github.com/persistenceOne/pstake-native/pull/549) add test for GetUnstakingEpochForPacket
- [#550](https://github.com/persistenceOne/pstake-native/pull/550) update height
- [#551](https://github.com/persistenceOne/pstake-native/pull/551) add undelegations json as embedded file

## [v2.0.0] - 2023-02-18

### Improvements

- [#397](https://github.com/persistenceOne/pstake-native/pull/397) use default auth ante.go
- [#399](https://github.com/persistenceOne/pstake-native/pull/399) add equal condition for undelegation.CompletionTime
  while checking mature undelegations
- [#401](https://github.com/persistenceOne/pstake-native/pull/401) bind baseDenom and mintDenom.
- [#403](https://github.com/persistenceOne/pstake-native/pull/403) disallow jumpstarting module from getting a second
  chance once enabled

### Removed

- [#392](https://github.com/persistenceOne/pstake-native/pull/392) remove fork logic

## [v1.3.0] -2022-12-22

### Improvements

- [#373](https://github.com/persistenceOne/pstake-native/pull/373) Added introduction spec.
- [#374](https://github.com/persistenceOne/pstake-native/pull/374) Updating concepts spec.
- [#375](https://github.com/persistenceOne/pstake-native/pull/375) Updating events spec.
- [#376](https://github.com/persistenceOne/pstake-native/pull/376) Updating keeper spec.
- [#377](https://github.com/persistenceOne/pstake-native/pull/377) Updated message spec.
- [#378](https://github.com/persistenceOne/pstake-native/pull/378) Restructuring spec folder.
- [#383](https://github.com/persistenceOne/pstake-native/pull/383) Adding individual weights validity check for
  validators.
- [#381](https://github.com/persistenceOne/pstake-native/pull/381) Oak report fixes.
- [#385](https://github.com/persistenceOne/pstake-native/pull/385) cap restake to 25%.
- [#390](https://github.com/persistenceOne/pstake-native/pull/390) Negative coin error fix.
- [#386](https://github.com/persistenceOne/pstake-native/pull/386)  reverts ibc from delegation account to deposit
  account.

### Removed

- [#384](https://github.com/persistenceOne/pstake-native/pull/384) remove msgjuice.

## [v1.2.1] -2022-12-5

### Improvements

- [#383](https://github.com/persistenceOne/pstake-native/pull/383) allow rest requests via grpc.

## [v1.2.0] -2022-11-9

### Features

- [#369](https://github.com/persistenceOne/pstake-native/pull/369) add recreate-ica tx.

### Improvements

- [#372](https://github.com/persistenceOne/pstake-native/pull/372) increase ICA timeout to 15mins, so it doesn't timeout
  during upgrades

## [v1.1.0] -2022-10-24

### Improvements

- [#358](https://github.com/persistenceOne/pstake-native/pull/358) add bounds for max fee limits, add allowlisted
  validators deduplication.
- [#359](https://github.com/persistenceOne/pstake-native/pull/359) added comments for all functions for godoc.

### Bug Fixes

- [#364](https://github.com/persistenceOne/pstake-native/pull/364) fix deterministic ordering of array.

## [v1.0.0] -2022-10-22

### Features

- liquid staking module for cosmoshub-4