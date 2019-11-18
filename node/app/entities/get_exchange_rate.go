package entities

import (
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/strkey"
	"github.com/velo-protocol/DRSv1/node/app/constants"
	"regexp"
)

type GetExchangeRateInput struct {
	AssetCode string
	Issuer    string
}

func (input *GetExchangeRateInput) Validate() error {
	if input.AssetCode == "" {
		return errors.Errorf("%s %s", constants.AssetCode, constants.ErrMustNotBeBlank)
	}

	if input.Issuer == "" {
		return errors.Errorf("%s %s", constants.Issuer, constants.ErrMustNotBeBlank)
	}

	if matched, _ := regexp.MatchString(`^[A-Za-z0-9]{1,7}$`, input.AssetCode); !matched {
		return errors.Errorf(constants.ErrInvalidFormat, constants.AssetCode)
	}

	ok := strkey.IsValidEd25519PublicKey(input.Issuer)
	if !ok {
		return errors.Errorf(constants.ErrInvalidFormat, constants.Issuer)
	}
	return nil
}

type GetExchangeRateOutPut struct {
	AssetCode              string
	Issuer                 string
	RedeemablePricePerUnit decimal.Decimal
	RedeemableCollateral   string
}
