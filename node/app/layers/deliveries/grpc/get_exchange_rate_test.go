package grpc_test

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	spec "github.com/velo-protocol/DRSv1/grpc"
	"github.com/velo-protocol/DRSv1/node/app/entities"
	"github.com/velo-protocol/DRSv1/node/app/errors"
	"testing"
)

func TestHandler_GetExchangeRate(t *testing.T) {

	var (
		assetCode              = "vTHB"
		issuer                 = "GCQDOOHRLBZW2A6COOMMWI5RAKGEZPBXSGZ6L6WA7M7GK3ZMHODDRAS3"
		redeemAblePrucePerUnit = decimal.NewFromFloat(1000.25)
		redeemableCollateral   = "VELO"
	)

	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		helper.mockController.Finish()

		getExchangeRateInput := &entities.GetExchangeRateInput{
			AssetCode: assetCode,
			Issuer:    issuer,
		}

		helper.mockUseCase.EXPECT().GetExchangeRate(context.Background(), getExchangeRateInput).Return(&entities.GetExchangeRateOutPut{
			AssetCode:              assetCode,
			Issuer:                 issuer,
			RedeemablePricePerUnit: redeemAblePrucePerUnit,
			RedeemableCollateral:   redeemableCollateral,
		}, nil)

		reply, err := helper.handler.GetExchangeRate(context.Background(), &spec.GetExchangeRateRequest{
			AssetCode: assetCode,
			Issuer:    issuer,
		})
		assert.NoError(t, err)
		assert.NotNil(t, reply)
		assert.Equal(t, assetCode, reply.AssetCode)
		assert.Equal(t, issuer, reply.Issuer)
		assert.Equal(t, redeemableCollateral, reply.RedeemableCollateral)
		assert.Equal(t, redeemAblePrucePerUnit.StringFixed(7), reply.RedeemablePricePerUnit)
	})

	t.Run("error, use case return error", func(t *testing.T) {
		helper := initTest(t)
		helper.mockController.Finish()

		getExchangeRateInput := &entities.GetExchangeRateInput{
			AssetCode: assetCode,
			Issuer:    issuer,
		}

		helper.mockUseCase.EXPECT().GetExchangeRate(context.Background(), getExchangeRateInput).Return(nil, nerrors.ErrInternal{Message: "some error has occurred"})

		_, err := helper.handler.GetExchangeRate(context.Background(), &spec.GetExchangeRateRequest{
			AssetCode: assetCode,
			Issuer:    issuer,
		})

		assert.Error(t, err)
	})
}
