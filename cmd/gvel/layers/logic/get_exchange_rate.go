package logic

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
	spec "gitlab.com/velo-labs/cen/grpc"
)

func (lo *logic) GetExchangeRate(input *entity.GetExchangeRateInput) (*entity.GetExchangeRateOutput, error) {
	defaultAccount := lo.AppConfig.GetDefaultAccount()
	accountBytes, err := lo.DB.Get([]byte(defaultAccount))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account from db")
	}

	var account entity.StellarAccount
	err = json.Unmarshal(accountBytes, &account)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal account")
	}

	result, err := lo.Velo.Client(nil).GetExchangeRate(context.Background(), &spec.GetExchangeRateRequest{
		AssetCode: input.AssetCode,
		Issuer:    input.Issuer,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get exchange rate")
	}

	return &entity.GetExchangeRateOutput{
		AssetCode:              result.AssetCode,
		Issuer:                 result.Issuer,
		RedeemablePricePerUnit: result.RedeemablePricePerUnit,
		RedeemableCollateral:   result.RedeemableCollateral,
	}, nil
}
