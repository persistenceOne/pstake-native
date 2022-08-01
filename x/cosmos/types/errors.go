package types

import sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrInvalid                                  = sdkErrors.Register(ModuleName, 42, "invalid")
	ErrInvalidVote                              = sdkErrors.Register(ModuleName, 43, "invalid vote option")
	ErrInvalidProposal                          = sdkErrors.Register(ModuleName, 44, "invalid proposal sender")
	ErrOrchAddressNotFound                      = sdkErrors.Register(ModuleName, 45, "orchestrator address not found")
	ErrInvalidGenesis                           = sdkErrors.Register(ModuleName, 46, "invalid genesis state")
	ErrUnknownProposal                          = sdkErrors.Register(ModuleName, 47, "unknown proposal")
	ErrInactiveProposal                         = sdkErrors.Register(ModuleName, 48, "inactive proposal")
	ErrTxnNotPresentInOutgoingPool              = sdkErrors.Register(ModuleName, 49, "txn not present in outgoing pool")
	ErrInvalidStatus                            = sdkErrors.Register(ModuleName, 50, "invalid status type")
	ErrModuleNotEnabled                         = sdkErrors.Register(ModuleName, 51, "module not enabled")
	ErrInvalidWithdrawDenom                     = sdkErrors.Register(ModuleName, 52, "invalid withdraw denom")
	ErrInvalidBondDenom                         = sdkErrors.Register(ModuleName, 53, "invalid bond denom")
	ErrInvalidCustodialAddress                  = sdkErrors.Register(ModuleName, 54, "invalid custodial address")
	ErrPubKeyNotFound                           = sdkErrors.Register(ModuleName, 55, "pubKey is empty")
	ErrOrchAddressPresentInSignaturePool        = sdkErrors.Register(ModuleName, 56, "orchestrator address present in signature pool")
	ErrOrcastratorPubkeyIsMultisig              = sdkErrors.Register(ModuleName, 57, "orcastrator pubkey is a multisig key")
	ErrInvalidMultisigPubkey                    = sdkErrors.Register(ModuleName, 58, "multisig pubkey invalid")
	ErrMoreMultisigAccountsBelongToOneValidator = sdkErrors.Register(ModuleName, 59, "More than 1 Multisig subkeys cannot be held by singular validator")
	ErrMultiSigAddressNotFound                  = sdkErrors.Register(ModuleName, 60, "multi sig address not found")
	ErrValidatorOrchestratorMappingNotFound     = sdkErrors.Register(ModuleName, 61, "validator orchestrator mapping not found")
	ErrMoreThanTwoOrchestratorAddressesMapping  = sdkErrors.Register(ModuleName, 62, "not allowed more than two orchestrator for one validator")
	ErrModuleAlreadyEnabled                     = sdkErrors.Register(ModuleName, 63, "module already enabled")
	ErrValidatorNotAllowed                      = sdkErrors.Register(ModuleName, 64, "validator not allowed to add orchestrator address")
)
