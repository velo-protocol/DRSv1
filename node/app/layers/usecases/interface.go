package usecases

import "gitlab.com/velo-labs/cen/node/app/entities"

type UseCase interface {
	SetupAccount(
		issuerCreationTx string,
		peggedValue string,
		peggedCurrency string,
		assetName string,
	) (*entities.Credit, error)
}
