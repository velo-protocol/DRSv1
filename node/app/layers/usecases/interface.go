package usecases

import (
	"context"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/node/app/errors"
)

type UseCase interface {
	SetupCredit(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*string, nerrors.NodeError)
	CreateWhiteList(ctx context.Context, veloTx *vtxnbuild.VeloTx) nerrors.NodeError
}
