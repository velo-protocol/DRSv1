package usecases

import (
	"context"
	"github.com/pkg/errors"
	vxdr "gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
)

func (useCase *useCase) CreateWhiteList(ctx context.Context, veloTxEnvelope *vxdr.VeloTxEnvelope) error {
	roleExists, err := useCase.WhitelistRepo.FindOneRole(whitelistEntity.Role)
	if err != nil {
		return nil, err
	}

	if roleExists == nil {
		return nil, errors.Wrap(constants.ErrRoleNotFound, constants.ErrCreateWhiteList.Error())
	}

	dbTx, err := useCase.WhitelistRepo.BeginTx()
	if err != nil {
		return nil, errors.Wrap(constants.ErrorToBeginTransaction, constants.ErrCreateWhiteList.Error())
	}

	entity, err := useCase.WhitelistRepo.CreateWhitelistTx(dbTx, &whitelistEntity)
	if err != nil {
		return nil, err
	}

	_, err = useCase.WhitelistRepo.CommitTx(dbTx)
	if err != nil {
		dbTx.Rollback()
		return nil, err
	}

	return &entity.ID, nil
}
