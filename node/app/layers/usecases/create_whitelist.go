package usecases

import (
	"context"
	"github.com/AlekSi/pointer"
	"github.com/pkg/errors"
	vxdr "gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
)

func (useCase *useCase) CreateWhiteList(ctx context.Context, veloTxEnvelope *vxdr.VeloTxEnvelope) error {
	pkSender := veloTxEnvelope.VeloTx.SourceAccount.Address()
	role := veloTxEnvelope.VeloTx.VeloOp.Body.WhiteListOp.Role
	address := veloTxEnvelope.VeloTx.VeloOp.Body.WhiteListOp.Address.Address()

	regulatorExists, err := useCase.WhitelistRepo.FindOneWhitelist(entities.WhitelistFilter{
		StellarPublicAddress: pointer.ToString(pkSender),
		RoleCode: pointer.ToString(string(vxdr.RoleRegulator)),
	})
	if err != nil {
		return err
	}

	if regulatorExists == nil {
		return errors.Wrap(constants.ErrRoleIsNotValid, constants.ErrCreateWhiteList.Error())
	}

	roleExists, err := useCase.WhitelistRepo.FindOneRole(string(role))
	if err != nil {
		return err
	}

	if roleExists == nil {
		return errors.Wrap(constants.ErrRoleNotFound, constants.ErrCreateWhiteList.Error())
	}

	dbTx := useCase.WhitelistRepo.BeginTx()
	if err != nil {
		return errors.Wrap(constants.ErrorToBeginTransaction, constants.ErrCreateWhiteList.Error())
	}

	_, err = useCase.WhitelistRepo.CreateWhitelistTx(dbTx, &entities.Whitelist{
		StellarPublicAddress: address,
		RoleCode: string(role),
	})
	if err != nil {
		return err
	}

	err = useCase.WhitelistRepo.CommitTx(dbTx)
	if err != nil {
		dbTx.Rollback()
		return err
	}

	return nil
}
