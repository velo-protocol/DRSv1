package usecases

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"gitlab.com/velo-labs/cen/node/app/errors"
)

func (useCase *useCase) GetExchangeRate(ctx context.Context, input *entities.GetExchangeRateInput) (*entities.GetExchangeRateOutPut, nerrors.NodeError) {
	// validate get exchange rate body
	if err := input.Validate(); err != nil {
		return nil, nerrors.ErrInvalidArgument{
			Message: err.Error(),
		}
	}

	// get and validate issuer account
	getIssuerAccountOutput, err := useCase.SubUseCase.GetIssuerAccount(ctx, &entities.GetIssuerAccountInput{
		IssuerAddress: input.Issuer,
	})
	if err != nil {
		return nil, nerrors.ErrPrecondition{
			Message: err.Error(),
		}
	}

	// get trusted partner from issuer account
	getTrustedPartnerFromIssuerAccountOutput, err := useCase.SubUseCase.GetTrustedPartnerFromIssuerAccount(ctx, &entities.GetTrustedPartnerFromIssuerAccountInput{
		IssuerAccount: getIssuerAccountOutput.Account,
	})
	if err != nil {
		return nil, nerrors.ErrPrecondition{Message: err.Error()}
	}

	// verify that trusted partner is in the trusted partner list
	drsAccountData, err := useCase.StellarRepo.GetDrsAccountData()
	if err != nil {
		return nil, nerrors.ErrInternal{
			Message: errors.Wrap(err, constants.ErrGetDrsAccountData).Error(),
		}
	}
	trustedPartnerMetaAddress, err := useCase.StellarRepo.GetAccountDecodedDataByKey(drsAccountData.TrustedPartnerListAddress, getTrustedPartnerFromIssuerAccountOutput.TrustedPartnerAccount.AccountID)
	if err != nil {
		return nil, nerrors.ErrPrecondition{
			Message: errors.Wrap(err, constants.ErrVerifyTrustedPartnerAccount).Error(),
		}
	}

	// verify that the asset is in the trust partner meta
	_, err = useCase.StellarRepo.GetAccountDecodedDataByKey(trustedPartnerMetaAddress, fmt.Sprintf("%s_%s", input.AssetCode, input.Issuer))
	if err != nil {
		return nil, nerrors.ErrPrecondition{
			Message: errors.Wrap(err, constants.ErrVerifyAssetCode).Error(),
		}
	}

	//get median price from price account
	medianPrice, err := useCase.StellarRepo.GetMedianPriceFromPriceAccount(drsAccountData.VeloPriceAddress(vxdr.Currency(getIssuerAccountOutput.PeggedCurrency)))
	if err != nil {
		return nil, nerrors.ErrPrecondition{Message: errors.Wrap(err, constants.ErrGetPriceOfPeggedCurrency).Error()}
	}
	if !medianPrice.IsPositive() {
		return nil, nerrors.ErrPrecondition{Message: constants.ErrMedianPriceMustBeGreaterThanZero}
	}

	// calculate collateral price
	collateralPrice := getIssuerAccountOutput.PeggedValue.Div(medianPrice).Truncate(7)

	return &entities.GetExchangeRateOutPut{
		AssetCode:              input.AssetCode,
		Issuer:                 input.Issuer,
		RedeemablePricePerUnit: collateralPrice,
		RedeemableCollateral:   string(vxdr.AssetVELO),
	}, nil
}
