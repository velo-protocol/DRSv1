package usecases

import (
	"context"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"gitlab.com/velo-labs/cen/node/app/errors"
)

type UseCase interface {
	SetupAccount(
		issuerCreationTx string,
		peggedValue string,
		peggedCurrency string,
		assetName string,
	) (*entities.Credit, error)
	CreateWhiteList(ctx context.Context, veloTxEnvelope *vxdr.VeloTxEnvelope) nerrors.NodeError
}
