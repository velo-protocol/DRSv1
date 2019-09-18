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

	regulatorEntity, err := useCase.WhitelistRepo.FindOneWhitelist(entities.WhitelistFilter{
		StellarPublicAddress: pointer.ToString(pkSender),
		RoleCode: pointer.ToString(string(vxdr.RoleRegulator)),
	})
	if err != nil {
		return err
	}

	if regulatorEntity == nil {
		return errors.Wrap(constants.ErrUnauthorized, constants.ErrCreateWhiteList.Error())
	}

	roleEntity, err := useCase.WhitelistRepo.FindOneRole(string(role))
	if err != nil {
		return err
	}

	if roleEntity == nil {
		return errors.Wrap(constants.ErrRoleNotFound, constants.ErrCreateWhiteList.Error())
	}

	_, err = useCase.WhitelistRepo.CreateWhitelist(&entities.Whitelist{
		StellarPublicAddress: address,
		RoleCode: string(role),
	})
	if err != nil {
		return err
	}

	return nil
}
