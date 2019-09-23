package usecases

import (
	"context"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/errors"
)

type UseCase interface {
	SetupCredit(ctx context.Context, veloTxEnvelope *vxdr.VeloTxEnvelope) (*string, nerrors.NodeError)
	CreateWhiteList(ctx context.Context, veloTxEnvelope *vxdr.VeloTxEnvelope) nerrors.NodeError
}
