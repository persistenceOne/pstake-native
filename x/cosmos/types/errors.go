package types

import sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrInvalid                                  = sdkErrors.Register(ModuleName, 42, "invalid")
	ErrInvalidVote                              = sdkErrors.Register(ModuleName, 43, "invalid vote option")
	ErrInvalidProposal                          = sdkErrors.Register(ModuleName, 44, "invalid proposal sender")
	ErrOrchAddressNotFound                      = sdkErrors.Register(ModuleName, 48, "orchestrator address not found")
	ErrInvalidGenesis                           = sdkErrors.Register(ModuleName, 49, "invalid genesis state")
	ErrUnknownProposal                          = sdkErrors.Register(ModuleName, 50, "unknown proposal")
	ErrInactiveProposal                         = sdkErrors.Register(ModuleName, 51, "inactive proposal")
	ErrTxnNotPresentInOutgoingPool              = sdkErrors.Register(ModuleName, 52, "txn not present in outgoing pool")
	ErrInvalidStatus                            = sdkErrors.Register(ModuleName, 54, "invalid status type")
	ErrModuleNotEnabled                         = sdkErrors.Register(ModuleName, 56, "module not enabled")
	ErrInvalidWithdrawDenom                     = sdkErrors.Register(ModuleName, 57, "invalid withdraw denom")
	ErrInvalidBondDenom                         = sdkErrors.Register(ModuleName, 58, "invalid bond denom")
	ErrInvalidCustodialAddress                  = sdkErrors.Register(ModuleName, 59, "invalid custodial address")
	ErrPubKeyNotFound                           = sdkErrors.Register(ModuleName, 61, "pubKey is empty")
	ErrOrchAddressPresentInSignaturePool        = sdkErrors.Register(ModuleName, 62, "orchestrator address present in signature pool")
	ErrOrcastratorPubkeyIsMultisig              = sdkErrors.Register(ModuleName, 63, "orcastrator pubkey is a multisig key")
	ErrInvalidMultisigPubkey                    = sdkErrors.Register(ModuleName, 64, "multisig pubkey invalid")
	ErrMoreMultisigAccountsBelongToOneValidator = sdkErrors.Register(ModuleName, 65, "More than 1 Multisig subkeys cannot be held by singular validator")
	ErrMultiSigAddressNotFound                  = sdkErrors.Register(ModuleName, 66, "multi sig address not found")
	ErrValidatorOrchestratorMappingNotFound     = sdkErrors.Register(ModuleName, 67, "validator orchestrator mapping not found")
	ErrMoreThanTwoOrchestratorAddressesMapping  = sdkErrors.Register(ModuleName, 68, "not allowed more than two orchestrator for one validator")
	ErrModuleAlreadyEnabled                     = sdkErrors.Register(ModuleName, 69, "module already enabled")
	ErrValidatorNotAllowed                      = sdkErrors.Register(ModuleName, 70, "validator not allowed to add orchestrator address")
)
