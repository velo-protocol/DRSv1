package usecases

import (
	"context"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	env "gitlab.com/velo-labs/cen/node/app/environments"
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
			Message: errors.Wrap(err, constants.ErrMedianPriceMustBeGreaterThanZero).Error(),
		}
	}

	// get median price usd
	medianPriceUsd, err := useCase.StellarRepo.GetMedianPriceFromPriceAccount(drsAccount.PriceUsdVeloAddress)
	if err != nil {
		return nil, nerrors.ErrPrecondition{
			Message: errors.Wrap(err, constants.ErrMedianPriceMustBeGreaterThanZero).Error(),
		}
	}

	// get median price sgd
	medianPriceSgd, err := useCase.StellarRepo.GetMedianPriceFromPriceAccount(drsAccount.PriceSgdVeloAddress)
	if err != nil {
		return nil, nerrors.ErrPrecondition{
			Message: errors.Wrap(err, constants.ErrMedianPriceMustBeGreaterThanZero).Error(),
		}
	}

	// get tp list data
	tpListData, err := useCase.StellarRepo.GetAccountDecodedData(drsAccount.TrustedPartnerListAddress)
	if err != nil {
		return nil, nerrors.ErrPrecondition{
			Message: errors.Wrap(err, constants.ErrGetTrustedPartnerListAccountData).Error(),
		}
	}

	var collateral = decimal.Zero

	// calculate collateral amount
	for _, tpMetaAddress := range tpListData {
		tpMetaData, err := useCase.StellarRepo.GetAccountDecodedData(tpMetaAddress)
		if err != nil {
			return nil, nerrors.ErrPrecondition{
				Message: errors.Wrap(err, constants.ErrGetTrustedPartnerMetaAccountDetail).Error(),
			}
		}

		// calculate collateral amount per tp
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

			if !stableAmount.GreaterThan(decimal.Zero) {
				return nil, nerrors.ErrPrecondition{Message: constants.ErrStableCreditAmountMustBeGreaterThanZero}
			}

			var collateralPerStable = decimal.Decimal{}
			switch vxdr.Currency(issuerAccount.PeggedCurrency) {
			case vxdr.CurrencyTHB:
				collateralPerStable = stableAmount.Mul(issuerAccount.PeggedValue).Div(medianPriceThb)
			case vxdr.CurrencySGD:
				collateralPerStable = stableAmount.Mul(issuerAccount.PeggedValue).Div(medianPriceSgd)
			case vxdr.CurrencyUSD:
				collateralPerStable = stableAmount.Mul(issuerAccount.PeggedValue).Div(medianPriceUsd)

			}
			// sum total collateral of tp
			collateralPerTp = collateralPerTp.Add(collateralPerStable)
		}
		// sum total collateral amount
		collateral = collateral.Add(collateralPerTp)
	}

	//get drs reserve account
	account, err := useCase.StellarRepo.GetAccount(drsAccount.DrsReserve)
	if err != nil {
		return nil, nerrors.ErrInternal{Message: constants.ErrGetDrsReserveAccountDetail}
	}

	// get drs reserve collateral amount
	var poolAmount = decimal.Zero
	for _, balance := range account.Balances {
		if balance.Code == string(vxdr.AssetVELO) && balance.Issuer == env.VeloIssuerPublicKey {
			balanceDecimal, err := decimal.NewFromString(balance.Balance)
			if err != nil {
				return nil, nerrors.ErrInternal{Message: err.Error()}
			}

			poolAmount = balanceDecimal
		}
	}

	return &entities.GetCollateralHealthCheckOutput{
		AssetCode:      string(vxdr.AssetVELO),
		AssetIssuer:    env.VeloIssuerPublicKey,
		RequiredAmount: collateral,
		PoolAmount:     poolAmount,
	}, nil
}
