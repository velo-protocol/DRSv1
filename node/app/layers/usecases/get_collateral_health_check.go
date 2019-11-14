package usecases

import (
	"context"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"gitlab.com/velo-labs/cen/node/app/environments"
	"gitlab.com/velo-labs/cen/node/app/errors"
	"strings"
)

func (useCase *useCase) GetCollateralHealthCheck(ctx context.Context) (*entities.GetCollateralHealthCheckOutput, nerrors.NodeError) {

	// get drs account
	drsAccount, err := useCase.StellarRepo.GetDrsAccountData()
	if err != nil {
		return nil, nerrors.ErrInternal{
			Message: errors.Wrap(err, constants.ErrGetDrsAccountData).Error(),
		}
	}

	// get median price thb
	medianPriceThb, err := useCase.StellarRepo.GetMedianPriceFromPriceAccount(drsAccount.PriceThbVeloAddress)
	if err != nil {
		return nil, nerrors.ErrPrecondition{
			Message: errors.Wrap(err, constants.ErrGetPriceOfPeggedCurrency).Error(),
		}
	}

	// get median price usd
	medianPriceUsd, err := useCase.StellarRepo.GetMedianPriceFromPriceAccount(drsAccount.PriceUsdVeloAddress)
	if err != nil {
		return nil, nerrors.ErrPrecondition{
			Message: errors.Wrap(err, constants.ErrGetPriceOfPeggedCurrency).Error(),
		}
	}

	// get median price sgd
	medianPriceSgd, err := useCase.StellarRepo.GetMedianPriceFromPriceAccount(drsAccount.PriceSgdVeloAddress)
	if err != nil {
		return nil, nerrors.ErrPrecondition{
			Message: errors.Wrap(err, constants.ErrGetPriceOfPeggedCurrency).Error(),
		}
	}

	// get tp list data
	tpListData, err := useCase.StellarRepo.GetAccountDecodedData(drsAccount.TrustedPartnerListAddress)
	if err != nil {
		return nil, nerrors.ErrPrecondition{
			Message: errors.Wrap(err, constants.ErrGetTrustedPartnerListAccountData).Error(),
		}
	}

	var drsCollateralRequiredAmount = decimal.Zero

	// calculate drs collateral required amount amount
	for _, tpMetaAddress := range tpListData {
		tpMetaData, err := useCase.StellarRepo.GetAccountDecodedData(tpMetaAddress)
		if err != nil {
			return nil, nerrors.ErrPrecondition{
				Message: errors.Wrap(err, constants.ErrGetTrustedPartnerMetaAccountDetail).Error(),
			}
		}

		// calculate drs collateral required amount amount per tp
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
			if err != nil {
				return nil, nerrors.ErrPrecondition{Message: errors.Wrapf(err, constants.ErrGetAsset, assetDetail[0]).Error()}
			}

			if len(assetPage.Embedded.Records) < 1 {
				continue
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
		// sum total drs collateral required amount amount
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

	return &entities.GetCollateralHealthCheckOutput{
		AssetCode:      string(vxdr.AssetVELO),
		AssetIssuer:    env.VeloIssuerPublicKey,
		RequiredAmount: drsCollateralRequiredAmount.Truncate(7),
		PoolAmount:     drsCollateralAmount.Truncate(7),
	}, nil
}
