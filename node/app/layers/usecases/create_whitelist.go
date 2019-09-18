package usecases

import (
	"context"
	"github.com/AlekSi/pointer"
	"github.com/pkg/errors"
	"gitlab.com/velo-labs/cen/libs/convert"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
)

func (useCase *useCase) CreateWhiteList(ctx context.Context, veloTxEnvelope *vxdr.VeloTxEnvelope) error {
	txSenderPublicKey := veloTxEnvelope.VeloTx.SourceAccount.Address()
	role := veloTxEnvelope.VeloTx.VeloOp.Body.WhiteListOp.Role
	address := veloTxEnvelope.VeloTx.VeloOp.Body.WhiteListOp.Address.Address()

	regulatorEntity, err := useCase.WhitelistRepo.FindOneWhitelist(entities.WhitelistFilter{
		StellarPublicAddress: pointer.ToString(txSenderPublicKey),
		RoleCode: pointer.ToString(string(vxdr.RoleRegulator)),
	})
	if err != nil {
		return err
	}

	if regulatorEntity == nil {
		return errors.Wrap(constants.ErrUnauthorized, constants.ErrCreateWhiteList.Error())
	}

	txSenderKeyPair, err := vconvert.PublicKeyToKeyPair(txSenderPublicKey)
	if err != nil {
		return errors.Wrap(err, constants.ErrCreateWhiteList.Error())
	}
	if txSenderKeyPair.Hint() != veloTxEnvelope.Signatures[0].Hint {
		return errors.Wrap(constants.ErrBadSignature, constants.ErrCreateWhiteList.Error())
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
