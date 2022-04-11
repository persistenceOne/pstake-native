<!--
order: 10
-->

# Parameters

The cosmos liquid staking module contains the following parameters:

| Key                               | Type                    | Example                |
|-----------------------------------|-------------------------|------------------------|
| MinMintingAmount                  |        | "259200000000000"      |
| MaxMintingAmount                  | uint16                  | 100                    |
| MinBurningAmount                  | uint16                  | 7                      |
| MaxBurningAmount                  | uint16                  | 3                      |
| MaxValidatorToDelegate            | string                  | "stake"                |
| ValidatorSetCosmosChain           | []WeightedAddressCosmos | "0.000000000000000000" |
| ValidatorSetNativeChain           | []WeightedAddress | "0.000000000000000000" |
| WeightedDeveloperRewardsReceivers | []WeightedAddress | "0.000000000000000000" |
| DistributionProportion            | []DistributionProportions | "0.000000000000000000" |
| Epochs                            | int64 | "0.000000000000000000" |
| MaxIncomingAndOutgoingTxns        | int64 | "0.000000000000000000" |
| CosmosProposalParams              | CosmosChainProposalParams |  |
| DelegationThreshold               | string | "0.000000000000000000" |
| ModuleEnabled                     | bool                    | false |
