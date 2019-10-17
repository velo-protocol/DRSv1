package usecases

import (
	"context"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"gitlab.com/velo-labs/cen/node/app/errors"
)

func (useCase *useCase) GetExchangeRate(ctx context.Context, getExchangeRate *entities.GetExchangeRateInput) (*entities.GetExchangeRateOutPut, nerrors.NodeError) {
	panic("implement me")
}
