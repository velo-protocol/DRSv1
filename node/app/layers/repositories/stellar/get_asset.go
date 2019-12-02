package stellar

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
	"github.com/velo-protocol/DRSv1/node/app/constants"
	"github.com/velo-protocol/DRSv1/node/app/entities"
)

func (repo *repo) GetAsset(getAssetInput entities.GetAssetInput) (*horizon.AssetsPage, error) {

	assetRequest := horizonclient.AssetRequest{
		ForAssetCode:   getAssetInput.AssetCode,
		ForAssetIssuer: getAssetInput.AssetIssuer,
	}

	if getAssetInput.Order != nil {
		assetRequest.Order = horizonclient.Order(*getAssetInput.Order)
	}

	if getAssetInput.Cursor != nil {
		assetRequest.Cursor = *getAssetInput.Cursor
	}

	if getAssetInput.Limit != nil {
		assetRequest.Limit = *getAssetInput.Limit
	}

	asset, err := repo.HorizonClient.Assets(assetRequest)
	if err != nil {
		return nil, errors.Wrapf(err, constants.ErrGetAsset, getAssetInput.AssetCode)
	}

	return &asset, nil
}
