package usecases

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/txnbuild"
	"gitlab.com/velo-labs/cen/libs/convert"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	env "gitlab.com/velo-labs/cen/node/app/environments"
	"gitlab.com/velo-labs/cen/node/app/errors"
	"gitlab.com/velo-labs/cen/node/app/utils"
	"strconv"
	"strings"
)

func (useCase *useCase) MintCredit(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*entities.MintCreditOutput, nerrors.NodeError) {
	if err := veloTx.VeloOp.Validate(); err != nil {
		return nil, nerrors.ErrInvalidArgument{Message: err.Error()}
	}

	op := veloTx.TxEnvelope().VeloTx.VeloOp.Body.MintCreditOp

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

	drsAccountData, err := useCase.StellarRepo.GetDrsAccountData()
	if err != nil {
		return nil, nerrors.ErrInternal{
			Message: errors.Wrap(err, constants.ErrGetDrsAccountData).Error(),
		}
	}

	// validate trusted partner role
	trustedPartnerListData, err := useCase.StellarRepo.GetAccountData(drsAccountData.TrustedPartnerListAddress)
	if err != nil {
		return nil, nerrors.ErrInternal{
			Message: errors.Wrap(err, constants.ErrGetTrustedPartnerListAccountData).Error(),
		}
	}

	trustedPartnerMetaEncodedAddress, ok := trustedPartnerListData[txSenderKeyPair.Address()]
	if !ok {
		return nil, nerrors.ErrPermissionDenied{
			Message: fmt.Sprintf(constants.ErrFormatSignerNotHavePermission, constants.VeloOpMintCredit),
		}
	}

	trustedPartnerMetaAddress, err := utils.DecodeBase64(trustedPartnerMetaEncodedAddress)
	if err != nil {
		return nil, nerrors.ErrInternal{
			Message: errors.Wrapf(err, constants.ErrToDecodeData, trustedPartnerMetaEncodedAddress).Error(),
		}
	}

	// get trusted partner meta
	trustedPartnerMeta, err := useCase.StellarRepo.GetAccountData(trustedPartnerMetaAddress)
	if err != nil {
		return nil, nerrors.ErrInternal{Message: errors.Wrapf(err, constants.ErrGetTrustedPartnerMetaAccountDetail).Error()}
	}

	var issuerAccount string
	var distributionAccount string
	for key, value := range trustedPartnerMeta {
		assetDetail := strings.Split(key, "_")
		if assetDetail[0] == op.AssetCodeToBeIssued {
			issuerAccount = assetDetail[1]
			distributionAccount, err = utils.DecodeBase64(value)
			if err != nil {
				return nil, nerrors.ErrInternal{Message: errors.Wrapf(err, constants.ErrToDecodeData, value).Error()}
			}
			break
		}

	}
	if issuerAccount == "" || distributionAccount == "" {
		return nil, nerrors.ErrPrecondition{
			Message: errors.Errorf(constants.ErrAssetCodeToBeIssuedNotSetup, op.AssetCodeToBeIssued).Error(),
		}
	}

	// get issuer account data
	issuerAccountData, err := useCase.StellarRepo.GetAccountDecodedData(issuerAccount)
	if err != nil {
		return nil, nerrors.ErrInternal{Message: errors.Wrapf(err, constants.ErrGetIssuerAccountDetail).Error()}
	}

	peggedCurrency, ok := issuerAccountData["peggedCurrency"]
	if !ok {
		return nil, nerrors.ErrPrecondition{
			Message: errors.Errorf(constants.ErrAssetCodeToBeIssuedNotSetup, op.AssetCodeToBeIssued).Error(),
		}
	}

	peggedValueString, ok := issuerAccountData["peggedValue"]
	if !ok {
		return nil, nerrors.ErrPrecondition{
			Message: errors.Errorf(constants.ErrAssetCodeToBeIssuedNotSetup, op.AssetCodeToBeIssued).Error(),
		}
	}

	//get median price from price account
	medianPrice, err := useCase.StellarRepo.GetMedianPriceFromPriceAccount(drsAccountData.VeloPriceAddress(vxdr.Currency(peggedCurrency)))
	if err != nil {
		return nil, nerrors.ErrPrecondition{Message: errors.Wrap(err, constants.ErrGetPriceOfPeggedCurrency).Error()}
	}

	peggedValueRaw, err := strconv.ParseInt(peggedValueString, 10, 64)
	if err != nil {
		return nil, nerrors.ErrInternal{Message: err.Error()}
	}
	peggedValue := decimal.New(peggedValueRaw, -7)
	if !peggedValue.GreaterThan(decimal.Zero) {
		return nil, nerrors.ErrPrecondition{Message: constants.ErrPeggedValueMustBeGreaterThanZero}
	}

	collateralAmount := decimal.New(int64(op.CollateralAmount), -7)
	collateralAsset := string(op.CollateralAssetCode)
	collateralAssetIssuer := env.VeloIssuerPublicKey

	mintAmount := collateralAmount.Mul(medianPrice).Div(peggedValue)

	drsKp, err := vconvert.SecretKeyToKeyPair(env.DrsSecretKey)
	if err != nil {
		return nil, nerrors.ErrInternal{
			Message: errors.Wrap(err, constants.ErrDerivedKeyPairFromSeed).Error(),
		}
	}

	tx := txnbuild.Transaction{
		SourceAccount: txSenderAccount,
		Network:       env.NetworkPassphrase,
		Timebounds:    txnbuild.NewTimeout(env.StellarTxTimeBoundInMinutes),
		Operations: []txnbuild.Operation{
			&txnbuild.Payment{
				Destination: drsKp.Address(),
				Amount:      collateralAmount.String(),
				Asset: txnbuild.CreditAsset{
					Code:   collateralAsset,
					Issuer: collateralAssetIssuer,
				},
				SourceAccount: txSenderAccount,
			},
			&txnbuild.Payment{
				Destination: distributionAccount,
				Amount:      mintAmount.String(),
				Asset: txnbuild.CreditAsset{
					Code:   op.AssetCodeToBeIssued,
					Issuer: issuerAccount,
				},
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: issuerAccount,
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
	return &entities.MintCreditOutput{
		SignedStellarTxXdr: signedTx,
		MintAmount:         mintAmount,
		MintCurrency:       op.AssetCodeToBeIssued,
		CollateralAmount:   collateralAmount,
		CollateralAsset:    collateralAsset,
	}, nil
}
