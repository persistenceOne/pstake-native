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

func (k Keeper) addToOutgoingSignaturePool(ctx sdk.Context, singleSignature cosmosTypes.SingleSignatureDataForOutgoingPool, txID uint64, orchestratorAddress sdk.AccAddress) error {
	outgoingSignaturePoolStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOutgoingSignaturePoolKey)
	key := cosmosTypes.UInt64Bytes(txID)
	if outgoingSignaturePoolStore.Has(key) {
		var outgoingSignaturePoolValue cosmosTypes.OutgoingSignaturePoolValue
		err := outgoingSignaturePoolValue.Unmarshal(outgoingSignaturePoolStore.Get(key))
		if err != nil {
			return err
		}
		if !outgoingSignaturePoolValue.Find(orchestratorAddress.String()) {
			return sdkErrors.Wrap(cosmosTypes.ErrOrchAddressPresentInSignaturePool, orchestratorAddress.String())
		}
		outgoingSignaturePoolValue.SingleSignatures = append(outgoingSignaturePoolValue.SingleSignatures, singleSignature)
		outgoingSignaturePoolValue.AddAndIncrement(orchestratorAddress.String())

		bz, err := outgoingSignaturePoolValue.Marshal()
		if err != nil {
			return err
		}
		outgoingSignaturePoolStore.Set(key, bz)
	}
	outgoingSignaturePoolValue := cosmosTypes.NewOutgoingSignaturePoolValue(singleSignature, orchestratorAddress)
	bz, err := outgoingSignaturePoolValue.Marshal()
	if err != nil {
		return err
	}
	outgoingSignaturePoolStore.Set(key, bz)
	return nil
}

func (k Keeper) getAllFromOutgoingSignaturePool(ctx sdk.Context) (list []OutgoingSignaturePoolKeyAndValue, err error) {
	outgoingSignaturePoolStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOutgoingSignaturePoolKey)
	iterator := outgoingSignaturePoolStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var outgoingSignaturePoolValue cosmosTypes.OutgoingSignaturePoolValue
		err = outgoingSignaturePoolValue.Unmarshal(iterator.Value())
		if err != nil {
			return list, err
		}
		txID := cosmosTypes.UInt64FromBytes(iterator.Key())
		list = append(list, OutgoingSignaturePoolKeyAndValue{txID, outgoingSignaturePoolValue})
	}
	return list, err
}

func (k Keeper) removeFromOutgoingSignaturePool(ctx sdk.Context, txID uint64) {
	outgoingSignaturePoolStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOutgoingSignaturePoolKey)
	key := cosmosTypes.UInt64Bytes(txID)
	outgoingSignaturePoolStore.Delete(key)
}

func (k Keeper) ProcessAllSignature(ctx sdk.Context) error {
	outgoingSignaturePool, err := k.getAllFromOutgoingSignaturePool(ctx)
	if err != nil {
		return err
	}
	params := k.GetParams(ctx)
	for _, os := range outgoingSignaturePool {
		if os.OutgoingSignaturePoolValue.Counter >= params.MultisigThreshold {
			//TODO ***IMPORTANT*** : discuss and implement how to maintain multisig sequence number.
			custodialAddress, err := cosmosTypes.AccAddressFromBech32(k.GetParams(ctx).CustodialAddress, cosmosTypes.Bech32Prefix)
			if err != nil {
				return err
			}
			multisigAccount := k.authKeeper.GetAccount(ctx, custodialAddress)
			multisigPub := multisigAccount.GetPubKey().(*multisig2.LegacyAminoPubKey)
			multisigSig := multisig.NewMultisig(len(multisigPub.PubKeys))

			for i, sig := range os.OutgoingSignaturePoolValue.SingleSignatures {
				externalSig := cosmosTypes.ConvertSingleSignatureDataForOutgoingPoolToSingleSignatureData(sig)
				orchAddress, err := sdk.AccAddressFromBech32(os.OutgoingSignaturePoolValue.OrchestratorAddresses[i])
				if err != nil {
					return err
				}
				account := k.authKeeper.GetAccount(ctx, orchAddress)
				if err := multisig.AddSignatureFromPubKey(multisigSig, &externalSig, account.GetPubKey(), multisigPub.GetPubKeys()); err != nil {
					return err
				}
			}

			sigV2 := signingtypes.SignatureV2{
				PubKey:   multisigPub,
				Data:     multisigSig,
				Sequence: multisigAccount.GetSequence(),
			}

			cosmosTx, err := k.getTxnFromOutgoingPoolByID(ctx, os.txID)
			if err != nil {
				return err
			}

			err = cosmosTx.CosmosTxDetails.SetSignatures(sigV2)
			if err != nil {
				return err
			}

			err = k.setOutgoingTxnSignaturesAndEmitEvent(ctx, cosmosTx.CosmosTxDetails, os.txID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
