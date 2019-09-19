package usecases

import (
	"context"
	"github.com/AlekSi/pointer"
	"github.com/pkg/errors"
	vconvert "gitlab.com/velo-labs/cen/libs/convert"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
)

func (useCase *useCase) CreateWhiteList(ctx context.Context, veloTxEnvelope *vxdr.VeloTxEnvelope) error {
	txSenderPublicKey := veloTxEnvelope.VeloTx.SourceAccount.Address()
	role := veloTxEnvelope.VeloTx.VeloOp.Body.WhiteListOp.Role
	address := veloTxEnvelope.VeloTx.VeloOp.Body.WhiteListOp.Address.Address()

	txSenderKeyPair, err := vconvert.PublicKeyToKeyPair(txSenderPublicKey)
	if err != nil {
		return errors.Wrap(err, constants.ErrCreateWhiteList.Error())
	}
	if txSenderKeyPair.Hint() != veloTxEnvelope.Signatures[0].Hint {
		return errors.Wrap(constants.ErrBadSignature, constants.ErrCreateWhiteList.Error())
	}

	regulatorEntity, err := useCase.WhitelistRepo.FindOneWhitelist(entities.WhiteListFilter{
		StellarPublicKey: pointer.ToString(txSenderPublicKey),
		RoleCode:         pointer.ToString(string(vxdr.RoleRegulator)),
	})
	if err != nil {
		return errors.Wrap(constants.ErrToGetDataFromDatabase, constants.ErrCreateWhiteList.Error())
	}

	if regulatorEntity == nil {
		return errors.Wrap(constants.ErrUnauthorized, constants.ErrCreateWhiteList.Error())
	}

	roleEntity, err := useCase.WhitelistRepo.FindOneRole(string(role))
	if err != nil {
		return errors.Wrap(constants.ErrToGetDataFromDatabase, constants.ErrCreateWhiteList.Error())
	}

	if roleEntity == nil {
		return errors.Wrap(constants.ErrRoleNotFound, constants.ErrCreateWhiteList.Error())
	}

	_, err = useCase.WhitelistRepo.CreateWhitelist(&entities.WhiteList{
		StellarPublicKey: address,
		RoleCode:         string(role),
	})
	if err != nil {
		return errors.Wrap(constants.ErrToSaveDatabase, constants.ErrCreateWhiteList.Error())
	}

	return nil
}