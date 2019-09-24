package usecases

import (
	"context"
	"fmt"
	"github.com/AlekSi/pointer"
	"gitlab.com/velo-labs/cen/libs/convert"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"gitlab.com/velo-labs/cen/node/app/errors"
	"strings"
)

func (useCase *useCase) CreateWhiteList(ctx context.Context, veloTx *vtxnbuild.VeloTx) nerrors.NodeError {
	if err := veloTx.VeloOp.Validate(); err != nil {
		return nerrors.ErrInvalidArgument{Message: err.Error()}
	}

	txSenderPublicKey := veloTx.TxEnvelope().VeloTx.SourceAccount.Address()
	role := veloTx.TxEnvelope().VeloTx.VeloOp.Body.WhiteListOp.Role
	address := veloTx.TxEnvelope().VeloTx.VeloOp.Body.WhiteListOp.Address.Address()

	txSenderKeyPair, err := vconvert.PublicKeyToKeyPair(txSenderPublicKey)
	if err != nil {
		return nerrors.ErrUnAuthenticated{Message: err.Error()}
	}
	if veloTx.TxEnvelope().Signatures == nil {
		return nerrors.ErrUnAuthenticated{Message: constants.ErrSignatureNotFound}
	}
	if txSenderKeyPair.Hint() != veloTx.TxEnvelope().Signatures[0].Hint {
		return nerrors.ErrUnAuthenticated{Message: constants.ErrSignatureNotMatchSourceAccount}
	}

	regulatorEntity, err := useCase.WhitelistRepo.FindOneWhitelist(entities.WhiteListFilter{
		StellarPublicKey: pointer.ToString(txSenderPublicKey),
		RoleCode:         pointer.ToString(string(vxdr.RoleRegulator)),
	})
	if err != nil {
		return nerrors.ErrInternal{Message: err.Error()}
	}
	if regulatorEntity == nil {
		return nerrors.ErrPermissionDenied{
			Message: fmt.Sprintf(constants.ErrFormatSignerNotHavePermission, constants.VeloOpWhiteList),
		}
	}

	roleEntity, err := useCase.WhitelistRepo.FindOneRole(string(role))
	if err != nil {
		return nerrors.ErrInternal{Message: err.Error()}
	}
	if roleEntity == nil {
		return nerrors.ErrNotFound{Message: constants.ErrRoleNotFound}
	}

	_, err = useCase.WhitelistRepo.CreateWhitelist(&entities.WhiteList{
		StellarPublicKey: address,
		RoleCode:         string(role),
	})

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates") {
			return nerrors.ErrAlreadyExists{
				Message: fmt.Sprintf(constants.ErrWhiteListAlreadyWhiteListed, txSenderPublicKey, vxdr.RoleMap[role]),
			}
		}
		return nerrors.ErrInternal{Message: constants.ErrToSaveDatabase}
	}

	return nil
}
