package usecases

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/amount"
	"github.com/stellar/go/txnbuild"
	"gitlab.com/velo-labs/cen/libs/convert"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	vxdr "gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	env "gitlab.com/velo-labs/cen/node/app/environments"
	"gitlab.com/velo-labs/cen/node/app/errors"
	"gitlab.com/velo-labs/cen/node/app/utils"
	"strconv"
)

func (useCase *useCase) RedeemCredit(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*entities.RedeemCreditOutput, nerrors.NodeError) {
	if err := veloTx.VeloOp.Validate(); err != nil {
		return nil, nerrors.ErrInvalidArgument{Message: err.Error()}
	}

	op := veloTx.TxEnvelope().VeloTx.VeloOp.Body.RedeemCreditOp

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

	// get and validate issuer account
	assetIssuerAccount, err := useCase.StellarRepo.GetAccount(op.Issuer.Address())
	if err != nil {
		return nil, nerrors.ErrNotFound{
			Message: errors.Wrap(err, constants.ErrGetIssuerAccount).Error(),
		}
	}
	if len(assetIssuerAccount.Signers) != 3 {
		return nil, nerrors.ErrPrecondition{
			Message: fmt.Sprintf(constants.ErrInvalidIssuerAccount, "signer count must be 3"),
		}
	}
	peggedValueString, err := utils.DecodeBase64(assetIssuerAccount.Data["peggedValue"])
	if err != nil {
		return nil, nerrors.ErrPrecondition{
			Message: fmt.Sprintf(constants.ErrInvalidIssuerAccount, "invalid pegged value format"),
		}
	}
	peggedValueRaw, err := strconv.ParseInt(peggedValueString, 10, 64)
	if err != nil {
		return nil, nerrors.ErrPrecondition{
			Message: fmt.Sprintf(constants.ErrInvalidIssuerAccount, "invalid pegged value format"),
		}
	}
	peggedValue := decimal.New(peggedValueRaw, -7)
	peggedCurrency, err := utils.DecodeBase64(assetIssuerAccount.Data["peggedCurrency"])
	if err != nil {
		return nil, nerrors.ErrPrecondition{
			Message: fmt.Sprintf(constants.ErrInvalidIssuerAccount, "invalid pegged currency format"),
		}
	}

	// derive trusted partner address
	var trustedPartnerAddress string
	var drsKeyFound, issuerKeyFound bool // make sure the other two keys are what we expected
	for _, signer := range assetIssuerAccount.Signers {
		if signer.Key != env.DrsPublicKey && signer.Key != op.Issuer.Address() {
			trustedPartnerAddress = signer.Key
		}

		if signer.Key == env.DrsPublicKey && signer.Weight == 1 {
			drsKeyFound = true
		}
		if signer.Key == op.Issuer.Address() && signer.Weight == 0 {
			issuerKeyFound = true
		}
	}
	if trustedPartnerAddress == "" || !drsKeyFound || !issuerKeyFound {
		return nil, nerrors.ErrPrecondition{
			Message: fmt.Sprintf(constants.ErrInvalidIssuerAccount, "no drs as signer"),
		}
	}

	// get trusted partner account
	trustedPartnerAccount, err := useCase.StellarRepo.GetAccount(trustedPartnerAddress)
	if err != nil {
		return nil, nerrors.ErrPrecondition{
			Message: errors.Wrap(err, constants.ErrGetTrustedPartnerAccountDetail).Error(),
		}
	}

	// verify that trusted partner is in the trusted partner list
	drsAccountData, err := useCase.StellarRepo.GetDrsAccountData()
	if err != nil {
		return nil, nerrors.ErrInternal{
			Message: errors.Wrap(err, constants.ErrGetDrsAccountData).Error(),
		}
	}
	trustedPartnerMetaAddress, err := useCase.StellarRepo.GetAccountDecodedDataByKey(drsAccountData.TrustedPartnerListAddress, trustedPartnerAccount.AccountID)
	if err != nil {
		return nil, nerrors.ErrPrecondition{
			Message: errors.Wrap(err, constants.ErrVerifyTrustedPartnerAccount).Error(),
		}
	}

	// verify that the asset is in the trust partner meta
	_, err = useCase.StellarRepo.GetAccountDecodedDataByKey(trustedPartnerMetaAddress, fmt.Sprintf("%s_%s", op.AssetCode, op.Issuer.Address()))
	if err != nil {
		return nil, nerrors.ErrPrecondition{
			Message: errors.Wrap(err, constants.ErrVerifyAssetCode).Error(),
		}
	}

	//get median price from price account
	medianPrice, err := useCase.StellarRepo.GetMedianPriceFromPriceAccount(drsAccountData.VeloPriceAddress(vxdr.Currency(peggedCurrency)))
	if err != nil {
		return nil, nerrors.ErrPrecondition{Message: errors.Wrap(err, constants.ErrGetPriceOfPeggedCurrency).Error()}
	}
	if !medianPrice.IsPositive() {
		return nil, nerrors.ErrPrecondition{Message: constants.ErrMedianPriceMustBeGreaterThanZero}
	}

	drsKp, err := vconvert.SecretKeyToKeyPair(env.DrsSecretKey)
	if err != nil {
		return nil, nerrors.ErrInternal{
			Message: errors.Wrap(err, constants.ErrDerivedKeyPairFromSeed).Error(),
		}
	}

	// calculate for velo amount
	assetAmount := decimal.New(int64(op.Amount), -7)
	veloAmount := assetAmount.Mul(peggedValue).Div(medianPrice).Truncate(7)

	tx := txnbuild.Transaction{
		SourceAccount: txSenderAccount,
		Network:       env.NetworkPassphrase,
		Timebounds:    txnbuild.NewTimeout(env.StellarTxTimeBoundInMinutes),
		Operations: []txnbuild.Operation{
			&txnbuild.Payment{
				Destination: assetIssuerAccount.AccountID,
				Amount:      amount.String(op.Amount),
				Asset: txnbuild.CreditAsset{
					Code:   op.AssetCode,
					Issuer: op.Issuer.Address(),
				},
				SourceAccount: txSenderAccount,
			},
			&txnbuild.Payment{
				Destination: txSenderAccount.AccountID,
				Amount:      veloAmount.String(),
				Asset: txnbuild.CreditAsset{
					Code:   string(vxdr.AssetVELO),
					Issuer: env.VeloIssuerPublicKey,
				},
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: drsKp.Address(),
				},
			},
		},
	}
	signedTx, err := tx.BuildSignEncode(drsKp)
	if err != nil {
		return nil, nerrors.ErrInternal{
			Message: errors.Wrap(err, constants.ErrBuildAndSignTransaction).Error(),
		}
	}
	return &entities.RedeemCreditOutput{
		SignedStellarTxXdr: signedTx,
	}, nil
}
