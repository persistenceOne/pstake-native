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

### Bug Fixes

- [#643](https://github.com/persistenceOne/pstake-native/pull/643) Re-add LSCosmos types to enable parsing of older transactions and gov data.

## [v2.3.2] - 2023-09-07

### Bug Fixes

- [#639](https://github.com/persistenceOne/pstake-native/pull/639) LSM cap room bug.

## [v2.3.1] - 2023-09-06

### Bug Fixes

- [#637](https://github.com/persistenceOne/pstake-native/pull/637) LSM bond factor validation fix.

## [v2.3.0] - 2023-09-06

### Features

- [#594](https://github.com/persistenceOne/pstake-native/pull/594) LSM integration.
- [#631](https://github.com/persistenceOne/pstake-native/pull/631) Add telemetry to measure time taken in begin/end block.

### Bug Fixes

- [#632](https://github.com/persistenceOne/pstake-native/pull/632) LSM cap fix.
- [#621](https://github.com/persistenceOne/pstake-native/pull/621) ICA recreation timeout fix.

### Removed
- [#627](https://github.com/persistenceOne/pstake-native/pull/627) Remove lscosmos, lspersistence in favour of liquidstakeibc.


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
  lscosmos: depracate, remove functional code

## [v2.1.0-rc0] - 2023-04-20

Never released.

### Improvements

- [#411](https://github.com/persistenceOne/pstake-native/pull/411) add admin functionality to disable module incase of
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
- [#550](https://github.com/persistenceOne/pstake-native/pull/550) udpate height
- [#551](https://github.com/persistenceOne/pstake-native/pull/551) add undelegations json as embeded file

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

- [#372](https://github.com/persistenceOne/pstake-native/pull/372) increase ICA timeout to 15mins, so it doesnt timeout
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

