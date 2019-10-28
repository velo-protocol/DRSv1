package usecases

import (
	"context"
	"github.com/AlekSi/pointer"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/txnbuild"
	"gitlab.com/velo-labs/cen/libs/convert"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"gitlab.com/velo-labs/cen/node/app/environments"
	"gitlab.com/velo-labs/cen/node/app/errors"
	"strings"
)

func (useCase *useCase) RebalanceReserve(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*entities.RebalanceOutput, nerrors.NodeError) {
	if err := veloTx.VeloOp.Validate(); err != nil {
		return nil, nerrors.ErrInvalidArgument{Message: err.Error()}
	}

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

	// get drs collateral account
	drsCollateralAccount, err := useCase.StellarRepo.GetDrsAccountData()
	if err != nil {
		return nil, nerrors.ErrInternal{
			Message: errors.Wrap(err, constants.ErrGetDrsAccountData).Error(),
		}
	}

	// get median price thb
	medianPriceThb, err := useCase.StellarRepo.GetMedianPriceFromPriceAccount(drsCollateralAccount.PriceThbVeloAddress)
	if err != nil {
		return nil, nerrors.ErrPrecondition{
			Message: errors.Wrap(err, constants.ErrGetPriceOfPeggedCurrency).Error(),
		}
	}

	// get median price usd
	medianPriceUsd, err := useCase.StellarRepo.GetMedianPriceFromPriceAccount(drsCollateralAccount.PriceUsdVeloAddress)
	if err != nil {
		return nil, nerrors.ErrPrecondition{
			Message: errors.Wrap(err, constants.ErrGetPriceOfPeggedCurrency).Error(),
		}
	}

	// get median price sgd
	medianPriceSgd, err := useCase.StellarRepo.GetMedianPriceFromPriceAccount(drsCollateralAccount.PriceSgdVeloAddress)
	if err != nil {
		return nil, nerrors.ErrPrecondition{
			Message: errors.Wrap(err, constants.ErrGetPriceOfPeggedCurrency).Error(),
		}
	}

	// get tp list data
	tpListData, err := useCase.StellarRepo.GetAccountDecodedData(drsCollateralAccount.TrustedPartnerListAddress)
	if err != nil {
		return nil, nerrors.ErrPrecondition{
			Message: errors.Wrap(err, constants.ErrGetTrustedPartnerListAccountData).Error(),
		}
	}

	var drsCollateralRequiredAmount = decimal.Zero

	// calculate drs collateral required amount
	for _, tpMetaAddress := range tpListData {
		tpMetaData, err := useCase.StellarRepo.GetAccountDecodedData(tpMetaAddress)
		if err != nil {
			return nil, nerrors.ErrPrecondition{
				Message: errors.Wrap(err, constants.ErrGetTrustedPartnerMetaAccountDetail).Error(),
			}
		}

		// calculate drs collateral required amount per tp
		var collateralPerTp = decimal.Zero
		for stableCredit := range tpMetaData {
			assetDetail := strings.Split(stableCredit, "_")
			if len(assetDetail) != 2 {
				return nil, nerrors.ErrPrecondition{Message: constants.ErrVerifyAssetCode}
			}

			issuerAccount, err := useCase.SubUseCase.GetIssuerAccount(ctx, &entities.GetIssuerAccountInput{IssuerAddress: assetDetail[1]})
			if err != nil {
				return nil, nerrors.ErrPrecondition{Message: errors.Wrap(err, constants.ErrGetIssuerAccount).Error()}
			}

			assetPage, err := useCase.StellarRepo.GetAsset(entities.GetAssetInput{
				AssetCode:   assetDetail[0],
				AssetIssuer: assetDetail[1],
			})
			if err != nil || len(assetPage.Embedded.Records) < 1 {
				return nil, nerrors.ErrPrecondition{Message: errors.Wrapf(err, constants.ErrGetAsset, assetDetail[0]).Error()}
			}
			stableAmount, err := decimal.NewFromString(assetPage.Embedded.Records[0].Amount)
			if err != nil {
				return nil, nerrors.ErrPrecondition{Message: errors.Wrapf(err, "invalid stable amount format").Error()}
			}

			var collateralPerCredit decimal.Decimal
			switch vxdr.Currency(issuerAccount.PeggedCurrency) {
			case vxdr.CurrencyTHB:
				collateralPerCredit = stableAmount.Mul(issuerAccount.PeggedValue).Div(medianPriceThb)
			case vxdr.CurrencySGD:
				collateralPerCredit = stableAmount.Mul(issuerAccount.PeggedValue).Div(medianPriceSgd)
			case vxdr.CurrencyUSD:
				collateralPerCredit = stableAmount.Mul(issuerAccount.PeggedValue).Div(medianPriceUsd)
			default:
				return nil, nerrors.ErrInternal{Message: constants.ErrPeggedCurrencyIsNotSupport}

			}
			// sum total drs collateral required amount of tp
			collateralPerTp = collateralPerTp.Add(collateralPerCredit)
		}
		// sum total drs collateral required amount
		drsCollateralRequiredAmount = drsCollateralRequiredAmount.Add(collateralPerTp)
	}

	// get drs collateral amount
	drsCollateralBalances, err := useCase.StellarRepo.GetAccountBalances(env.DrsPublicKey)
	if err != nil {
		return nil, nerrors.ErrInternal{Message: errors.Wrap(err, constants.ErrGetDrsAccountBalance).Error()}
	}

	var drsCollateralAmount = decimal.Zero
	var drsCollateralAssetCode string
	var drsCollateralAssetIssuer string
	for _, balance := range drsCollateralBalances {
		if balance.Code == string(vxdr.AssetVELO) && balance.Issuer == env.VeloIssuerPublicKey {
			balanceDecimal, err := decimal.NewFromString(balance.Balance)
			if err != nil {
				return nil, nerrors.ErrInternal{Message: err.Error()}
			}
			drsCollateralAssetCode = balance.Code
			drsCollateralAssetIssuer = balance.Issuer
			drsCollateralAmount = balanceDecimal
		}
	}
	if drsCollateralAssetIssuer == "" || drsCollateralAssetCode == "" {
		return nil, nerrors.ErrInternal{Message: constants.ErrDrsCollateralTrustlineNotFound}
	}

	rebalanceOutput := &entities.RebalanceOutput{
		Collaterals: []*entities.Collateral{
			{
				AssetCode:      drsCollateralAssetCode,
				AssetIssuer:    drsCollateralAssetIssuer,
				RequiredAmount: drsCollateralRequiredAmount.Truncate(7),
				PoolAmount:     drsCollateralAmount.Truncate(7),
			},
		},
	}

	drsCollateralKP, err := vconvert.SecretKeyToKeyPair(env.DrsSecretKey)
	if err != nil {
		return nil, nerrors.ErrInternal{
			Message: errors.Wrap(err, constants.ErrDerivedKeyPairFromSeed).Error(),
		}
	}

	if drsCollateralAmount.GreaterThan(drsCollateralRequiredAmount) { // collateral amount greater than required amount
		tx := txnbuild.Transaction{
			SourceAccount: txSenderAccount,
			Network:       env.NetworkPassphrase,
			Timebounds:    txnbuild.NewTimeout(env.StellarTxTimeBoundInMinutes),
			Operations: []txnbuild.Operation{
				&txnbuild.Payment{
					Destination: drsCollateralAccount.DrsReserve,
					Amount:      drsCollateralAmount.Sub(drsCollateralRequiredAmount).String(),
					Asset: txnbuild.CreditAsset{
						Code:   drsCollateralAssetCode,
						Issuer: drsCollateralAssetIssuer,
					},
					SourceAccount: &txnbuild.SimpleAccount{
						AccountID: drsCollateralKP.Address(),
					},
				},
			},
		}
		signedTx, err := tx.BuildSignEncode(drsCollateralKP)
		if err != nil {
			return nil, nerrors.ErrInternal{
				Message: errors.Wrap(err, constants.ErrBuildAndSignTransaction).Error(),
			}
		}
		rebalanceOutput.SignedStellarTxXdr = &signedTx
	} else if drsCollateralAmount.LessThan(drsCollateralRequiredAmount) { // required amount greater than collateral amount

		tx := txnbuild.Transaction{
			SourceAccount: txSenderAccount,
			Network:       env.NetworkPassphrase,
			Timebounds:    txnbuild.NewTimeout(env.StellarTxTimeBoundInMinutes),
			Operations: []txnbuild.Operation{
				&txnbuild.Payment{
					Destination: drsCollateralKP.Address(),
					Amount:      drsCollateralRequiredAmount.Sub(drsCollateralAmount).String(),
					Asset: txnbuild.CreditAsset{
						Code:   drsCollateralAssetCode,
						Issuer: drsCollateralAssetIssuer,
					},
					SourceAccount: &txnbuild.SimpleAccount{
						AccountID: drsCollateralAccount.DrsReserve,
					},
				},
			},
		}
		signedTx, err := tx.BuildSignEncode(drsCollateralKP)
		if err != nil {
			return nil, nerrors.ErrInternal{
				Message: errors.Wrap(err, constants.ErrBuildAndSignTransaction).Error(),
			}
		}
		rebalanceOutput.SignedStellarTxXdr = &signedTx
	} else {
		rebalanceOutput.SignedStellarTxXdr = pointer.ToString("")
	}

	return rebalanceOutput, nil
}
