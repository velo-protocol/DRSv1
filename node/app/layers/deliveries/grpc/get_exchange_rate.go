package grpc

import (
	"context"
	spec "gitlab.com/velo-labs/cen/grpc"
	"gitlab.com/velo-labs/cen/node/app/entities"
)

func (handler handler) GetExchangeRate(ctx context.Context, getExchangeRate *spec.GetExchangeRateRequest) (*spec.GetExchangeRateReply, error) {

	entity := &entities.GetExchangeRateInput{
		AssetCode: getExchangeRate.AssetCode,
		Issuer:    getExchangeRate.Issuer,
	}

	getExchangeRateOutput, err := handler.UseCase.GetExchangeRate(ctx, entity)
	if err != nil {
		return nil, err.GRPCError()
	}

	return &spec.GetExchangeRateReply{
		AssetCode:              getExchangeRateOutput.AssetCode,
		Issuer:                 getExchangeRateOutput.Issuer,
		RedeemableCollateral:   getExchangeRateOutput.RedeemableCollateral,
		RedeemablePricePerUnit: getExchangeRateOutput.RedeemablePricePerUnit.StringFixed(7),
	}, nil
}
