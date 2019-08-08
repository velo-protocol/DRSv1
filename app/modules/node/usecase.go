package node

import "gitlab.com/velo-labs/cen/app/entities"

type UseCase interface {
	Setup(
		issuerCreationTx string,
		peggedValue string,
		peggedCurrency string,
		assetName string,
	) (*entities.Credit, error)
}
