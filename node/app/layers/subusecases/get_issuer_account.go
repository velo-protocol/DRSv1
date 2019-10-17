package subusecases

import (
	"context"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"gitlab.com/velo-labs/cen/node/app/utils"
	"strconv"
)

func (subUseCase *subUseCase) GetIssuerAccount(ctx context.Context, input *entities.GetIssuerAccountInput) (*entities.GetIssuerAccountOutput, error) {
	// get and validate issuer account
	assetIssuerAccount, err := subUseCase.StellarRepo.GetAccount(input.IssuerAddress)
	if err != nil {
		return nil, errors.Wrap(err, constants.ErrGetIssuerAccount)
	}
	if len(assetIssuerAccount.Signers) != 3 {
		return nil, errors.Errorf(constants.ErrInvalidIssuerAccount, "signer count must be 3")
	}
	peggedValueString, err := utils.DecodeBase64(assetIssuerAccount.Data["peggedValue"])
	if err != nil {
		return nil, errors.Errorf(constants.ErrInvalidIssuerAccount, "invalid pegged value format")
	}
	peggedValueRaw, err := strconv.ParseInt(peggedValueString, 10, 64)
	if err != nil {
		return nil, errors.Errorf(constants.ErrInvalidIssuerAccount, "invalid pegged value format")
	}
	peggedValue := decimal.New(peggedValueRaw, -7)
	peggedCurrency, err := utils.DecodeBase64(assetIssuerAccount.Data["peggedCurrency"])
	if err != nil {
		return nil, errors.Errorf(constants.ErrInvalidIssuerAccount, "invalid pegged currency format")
	}
	assetCode, err := utils.DecodeBase64(assetIssuerAccount.Data["assetCode"])
	if err != nil {
		return nil, errors.Errorf(constants.ErrInvalidIssuerAccount, "invalid asset code format")
	}

	return &entities.GetIssuerAccountOutput{
		Account:        assetIssuerAccount,
		PeggedValue:    peggedValue,
		PeggedCurrency: peggedCurrency,
		AssetCode:      assetCode,
	}, nil
}
