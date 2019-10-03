package usecases

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stellar/go/txnbuild"
	vconvert "gitlab.com/velo-labs/cen/libs/convert"
	vtxnbuild "gitlab.com/velo-labs/cen/libs/txnbuild"
	vxdr "gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/constants"
	env "gitlab.com/velo-labs/cen/node/app/environments"
	nerrors "gitlab.com/velo-labs/cen/node/app/errors"
	"gitlab.com/velo-labs/cen/node/app/utils"
	"time"
)

func (useCase *useCase) UpdatePrice(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*string, nerrors.NodeError) {
	if err := veloTx.VeloOp.Validate(); err != nil {
		return nil, nerrors.ErrInvalidArgument{Message: err.Error()}
	}

	// validate tx signature
	priceUpdateOp := veloTx.TxEnvelope().VeloTx.VeloOp.Body.PriceUpdateOp
	txSenderPublicKey := veloTx.TxEnvelope().VeloTx.SourceAccount.Address()
	txSenderKeyPair, err := vconvert.PublicKeyToKeyPair(txSenderPublicKey)
	if err != nil {
		return nil, nerrors.ErrInvalidArgument{Message: err.Error()}
	}
	if veloTx.TxEnvelope().Signatures == nil {
		return nil, nerrors.ErrUnAuthenticated{Message: constants.ErrSignatureNotFound}
	}
	if txSenderKeyPair.Hint() != veloTx.TxEnvelope().Signatures[0].Hint {
		return nil, nerrors.ErrUnAuthenticated{Message: constants.ErrSignatureNotMatchSourceAccount}
	}

	// get tx sender account
	txSenderAccount, err := useCase.StellarRepo.GetAccount(veloTx.SourceAccount.GetAccountID())
	if err != nil {
		return nil, nerrors.ErrNotFound{
			Message: errors.Wrap(err, constants.ErrGetSenderAccount).Error(),
		}
	}

	// get drs account
	drsAccountData, err := useCase.StellarRepo.GetDrsAccountData()
	if err != nil {
		return nil, nerrors.ErrInternal{
			Message: errors.Wrap(err, constants.ErrGetDrsAccountData).Error(),
		}
	}

	// get price feeder list account
	priceFeederListAccountData, err := useCase.StellarRepo.GetAccountData(drsAccountData.PriceFeederListAddress)
	if err != nil {
		return nil, nerrors.ErrInternal{
			Message: errors.Wrap(err, constants.ErrGetPriceFeederListAccountData).Error(),
		}
	}

	// validate tx sender role, in which must be price feeder
	priceFeederCurrencyBase64, ok := priceFeederListAccountData[txSenderKeyPair.Address()]
	if !ok {
		return nil, nerrors.ErrPermissionDenied{
			Message: fmt.Sprintf(constants.ErrFormatSignerNotHavePermission, constants.VeloOpPriceUpdate),
		}
	}

	registeredCurrency, err := utils.DecodeBase64(priceFeederCurrencyBase64)
	if vxdr.Currency(registeredCurrency) != priceUpdateOp.Currency {
		return nil, nerrors.ErrInvalidArgument{Message: constants.ErrCurrencyMustMatchWithRegisteredCurrency}
	}

	// prepare tx
	drsKp, err := vconvert.SecretKeyToKeyPair(env.DrsSecretKey)
	if err != nil {
		return nil, nerrors.ErrInternal{
			Message: errors.Wrap(err, constants.ErrDerivedKeyPairFromSeed).Error(),
		}
	}
	tx := txnbuild.Transaction{
		SourceAccount: txSenderAccount,
		Operations: []txnbuild.Operation{
			&txnbuild.ManageData{
				Name:  txSenderKeyPair.Address(),
				Value: []byte(fmt.Sprintf("%d_%d", time.Now().Unix(), priceUpdateOp.PriceInCurrencyPerAssetUnit)),
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: drsAccountData.VeloPriceAddress(vxdr.Currency(registeredCurrency)),
				},
			},
		},
		Network:    env.NetworkPassphrase,
		Timebounds: txnbuild.NewTimeout(env.StellarTxTimeBoundInMinutes),
	}

	signedTxXdr, err := tx.BuildSignEncode(drsKp)
	if err != nil {
		return nil, nerrors.ErrInternal{
			Message: errors.Wrap(err, constants.ErrBuildAndSignTransaction).Error(),
		}
	}

	return &signedTxXdr, nil
}
