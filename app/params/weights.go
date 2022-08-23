package params

// Default simulation operation weights for messages and gov proposals
const (
	DefaultWeightMsgSend                        int = 100
	DefaultWeightMsgMultiSend                   int = 10
	DefaultWeightMsgSetWithdrawAddress          int = 50
	DefaultWeightMsgWithdrawDelegationReward    int = 50
	DefaultWeightMsgWithdrawValidatorCommission int = 50
	DefaultWeightMsgFundCommunityPool           int = 50
	DefaultWeightMsgDeposit                     int = 100
	DefaultWeightMsgVote                        int = 67
	DefaultWeightMsgUnjail                      int = 100
	DefaultWeightMsgCreateValidator             int = 100
	DefaultWeightMsgEditValidator               int = 5
	DefaultWeightMsgDelegate                    int = 100
	DefaultWeightMsgUndelegate                  int = 100
	DefaultWeightMsgBeginRedelegate             int = 100

	DefaultWeightCommunitySpendProposal int = 5
	DefaultWeightTextProposal           int = 5
	DefaultWeightParamChangeProposal    int = 5

	//  Params from lspersistence module
	DefaultLSPWeightMsgLiquidStake                    int = 80
	DefaultLSPWeightMsgLiquidUnstake                  int = 30
	DefaultLSPWeightAddWhitelistValidatorsProposal    int = 50
	DefaultLSPWeightUpdateWhitelistValidatorsProposal int = 5
	DefaultLSPWeightDeleteWhitelistValidatorsProposal int = 5
	DefaultLSPWeightCompleteRedelegationUnbonding     int = 30
	DefaultLSPWeightTallyWithLiquidStaking            int = 30
)
