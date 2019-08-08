package client

import "gitlab.com/velo-labs/cen/app/entities"

type UseCase interface {
	Setup(
		setupXdr string,
		peggedValue string,
		peggedCurrency string,
		assetName string,
	) (*entities.Mint, error)
}
