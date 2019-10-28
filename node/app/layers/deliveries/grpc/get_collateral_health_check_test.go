package grpc_test

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/grpc"
	"gitlab.com/velo-labs/cen/node/app/entities"
	nerrors "gitlab.com/velo-labs/cen/node/app/errors"
	"testing"
)

func TestHandler_GetCollateralHealthCheck(t *testing.T) {

	var (
		assetCode      = "vTHB"
		assetIssuer    = "GCQDOOHRLBZW2A6COOMMWI5RAKGEZPBXSGZ6L6WA7M7GK3ZMHODDRAS3"
		requiredAmount = decimal.NewFromFloat(150.2230124)
		poolAmount     = decimal.NewFromFloat(250.0092210)
	)

	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		helper.mockController.Finish()

		helper.mockUseCase.EXPECT().GetCollateralHealthCheck(context.Background()).Return(&entities.GetCollateralHealthCheckOutput{
			AssetCode:      assetCode,
			AssetIssuer:    assetIssuer,
			RequiredAmount: requiredAmount,
			PoolAmount:     poolAmount,
		}, nil)

		getCollateralHealthCheckOutput, err := helper.handler.GetCollateralHealthCheck(context.Background(), &grpc.GetCollateralHealthCheckEmpty{})
		assert.NoError(t, err)
		assert.NotNil(t, getCollateralHealthCheckOutput)
		assert.Equal(t, assetCode, getCollateralHealthCheckOutput.AssetCode)
		assert.Equal(t, assetIssuer, getCollateralHealthCheckOutput.AssetIssuer)
		assert.Equal(t, requiredAmount.StringFixed(7), getCollateralHealthCheckOutput.RequiredAmount)
		assert.Equal(t, poolAmount.StringFixed(7), getCollateralHealthCheckOutput.PoolAmount)

	})

	t.Run("error, use case return error", func(t *testing.T) {
		helper := initTest(t)
		helper.mockController.Finish()

		helper.mockUseCase.EXPECT().GetCollateralHealthCheck(context.Background()).Return(nil, nerrors.ErrInternal{Message: "some error has occurred"})

		getCollateralHealthCheckOutput, err := helper.handler.GetCollateralHealthCheck(context.Background(), &grpc.GetCollateralHealthCheckEmpty{})
		assert.Error(t, err)
		assert.Nil(t, getCollateralHealthCheckOutput)
	})
}
