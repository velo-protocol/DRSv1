package usecases

import (
	"context"
	vxdr "gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/entities"
)

type UseCase interface {
	SetupAccount(
		issuerCreationTx string,
		peggedValue string,
		peggedCurrency string,
		assetName string,
	) (*entities.Credit, error)

	CreateWhiteList(ctx context.Context, veloTxEnvelope *vxdr.VeloTxEnvelope) error
}
