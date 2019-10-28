package usecases

import (
	"context"
	vtxnbuild "gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"gitlab.com/velo-labs/cen/node/app/errors"
)

func (useCase *useCase) RebalanceReserve(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*entities.RebalanceOutput, nerrors.NodeError) {
	return nil, nil
}
