package keeper

import (
	multisig2 "github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	"github.com/cosmos/cosmos-sdk/crypto/types/multisig"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

type OutgoingSignaturePoolKeyAndValue struct {
	txID                       uint64
	OutgoingSignaturePoolValue cosmosTypes.OutgoingSignaturePoolValue
}

/*
Adds the signature entry to the signature pool store with the given validator address.
Performs the following actions :
  1. Checks if the store has the key, if it has the key then it appends the signature.
  2. If not present in the store then creates a new entry.
*/
func (k Keeper) addToOutgoingSignaturePool(ctx sdk.Context, singleSignature cosmosTypes.SingleSignatureDataForOutgoingPool, txID uint64, orchestratorAddress sdk.AccAddress, validatorAddress sdk.ValAddress) error {
	outgoingSignaturePoolStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOutgoingSignaturePoolKey)
	key := cosmosTypes.UInt64Bytes(txID)
	if outgoingSignaturePoolStore.Has(key) {
		var outgoingSignaturePoolValue cosmosTypes.OutgoingSignaturePoolValue
		k.cdc.MustUnmarshal(outgoingSignaturePoolStore.Get(key), &outgoingSignaturePoolValue)

		if !outgoingSignaturePoolValue.Find(validatorAddress.String()) {
			return sdkErrors.Wrap(cosmosTypes.ErrOrchAddressPresentInSignaturePool, validatorAddress.String())
		}
		outgoingSignaturePoolValue.SingleSignatures = append(outgoingSignaturePoolValue.SingleSignatures, singleSignature)
		outgoingSignaturePoolValue.UpdateValues(validatorAddress.String(), k.GetTotalValidatorOrchestratorCount(ctx))
		outgoingSignaturePoolValue.OrchestratorAddresses = append(outgoingSignaturePoolValue.OrchestratorAddresses, orchestratorAddress.String())

		outgoingSignaturePoolStore.Set(key, k.cdc.MustMarshal(&outgoingSignaturePoolValue))
		return nil
	}
	outgoingSignaturePoolValue := cosmosTypes.NewOutgoingSignaturePoolValue(singleSignature, validatorAddress, orchestratorAddress)
	outgoingSignaturePoolStore.Set(key, k.cdc.MustMarshal(&outgoingSignaturePoolValue))
	return nil
}

// Gets all the entries from the outgoing signature pool
func (k Keeper) getAllFromOutgoingSignaturePool(ctx sdk.Context) (list []OutgoingSignaturePoolKeyAndValue, err error) {
	outgoingSignaturePoolStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOutgoingSignaturePoolKey)
	iterator := outgoingSignaturePoolStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var outgoingSignaturePoolValue cosmosTypes.OutgoingSignaturePoolValue
		err = k.cdc.Unmarshal(iterator.Value(), &outgoingSignaturePoolValue)
		if err != nil {
			return list, err
		}
		txID := cosmosTypes.UInt64FromBytes(iterator.Key())
		list = append(list, OutgoingSignaturePoolKeyAndValue{txID, outgoingSignaturePoolValue})
	}
	return list, err
}

// Removes the entry corresponding to the given txID from the outgoing signature pool
func (k Keeper) removeFromOutgoingSignaturePool(ctx sdk.Context, txID uint64) {
	outgoingSignaturePoolStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOutgoingSignaturePoolKey)
	key := cosmosTypes.UInt64Bytes(txID)
	outgoingSignaturePoolStore.Delete(key)
}

/*
ProcessAllSignature processes all the outgoing signature entries
This function is called every EndBlocker to perform the defined set of actions as mentioned below :
   1. Get the list of all outgoing signature entries
   2. Checks if the signatures sent have crossed the threshold.
   3. If majority is reached and other conditions match then the signature is added to the transaction.
   4. Once the signature is added, a signed outgoing txn event is generated.
*/
func (k Keeper) ProcessAllSignature(ctx sdk.Context) {
	outgoingSignaturePool, err := k.getAllFromOutgoingSignaturePool(ctx)
	if err != nil {
		panic(err)
	}

	for _, os := range outgoingSignaturePool {
		ka, ok := k.GetAccountState(ctx, k.GetCurrentAddress(ctx)).GetPubKey().(multisig.PubKey)
		if !ok {
			panic("not able to convert to pubkey")
		}
		if os.OutgoingSignaturePoolValue.Counter >= uint64(ka.GetThreshold()) {
			multisigAcc := k.GetAccountState(ctx, k.GetCurrentAddress(ctx))
			multisigPub := multisigAcc.GetPubKey().(*multisig2.LegacyAminoPubKey)
			multisigSig := multisig.NewMultisig(len(multisigPub.PubKeys))

			for i, sig := range os.OutgoingSignaturePoolValue.SingleSignatures {
				externalSig := cosmosTypes.ConvertSingleSignatureDataForOutgoingPoolToSingleSignatureData(sig)
				orchAddress, err := sdk.AccAddressFromBech32(os.OutgoingSignaturePoolValue.OrchestratorAddresses[i])
				if err != nil {
					panic(err)
				}
				account := k.AuthKeeper.GetAccount(ctx, orchAddress)
				if err := multisig.AddSignatureFromPubKey(multisigSig, &externalSig, account.GetPubKey(), multisigPub.GetPubKeys()); err != nil {
					panic(err)
				}
			}

			sigV2 := signingtypes.SignatureV2{
				PubKey:   multisigPub,
				Data:     multisigSig,
				Sequence: multisigAcc.GetSequence(),
			}

			cosmosTx, err := k.GetTxnFromOutgoingPoolByID(ctx, os.txID)
			if err != nil {
				panic(err)
			}

			err = cosmosTx.CosmosTxDetails.SetSignatures(sigV2)
			if err != nil {
				panic(err)
			}

			err = k.SetOutgoingTxnSignaturesAndEmitEvent(ctx, cosmosTx.CosmosTxDetails, os.txID)
			if err != nil {
				panic(err)
			}

			k.removeFromOutgoingSignaturePool(ctx, os.txID)
		}
	}
}
