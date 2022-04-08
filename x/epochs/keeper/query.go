package keeper

import (
	sdkCodec "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/persistenceOne/pstake-native/x/epochs/types"
	abciTypes "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper, legacyQuerierCdc *sdkCodec.LegacyAmino) sdk.Querier {
	return func(ctx sdkTypes.Context, path []string, req abciTypes.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryEpochInfos:
			return queryEpochInfos(ctx, req, k, legacyQuerierCdc)
		default:
			return nil, sdkErrors.Wrapf(sdkErrors.ErrUnknownRequest, "unknown %s query endpoint", types.ModuleName)
		}
	}
}

func queryEpochInfos(ctx sdkTypes.Context, req abciTypes.RequestQuery, k Keeper, legacyQuerierCdc *sdkCodec.LegacyAmino) ([]byte, error) {
	var epochInfosRequest types.QueryEpochsInfoRequest

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &epochInfosRequest)
	if err != nil {
		return nil, sdkErrors.Wrap(sdkErrors.ErrJSONUnmarshal, err.Error())
	}

	epochInfo := k.GetEpochInfo(ctx, "3minutes")

	bz, err := sdkCodec.MarshalJSONIndent(legacyQuerierCdc, epochInfo)
	if err != nil {
		return nil, sdkErrors.Wrap(sdkErrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
