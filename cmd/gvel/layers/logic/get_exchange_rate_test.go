package logic_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/velo-protocol/DRSv1/cmd/gvel/entity"
	"github.com/velo-protocol/DRSv1/grpc"
	"testing"
)

func TestLogic_GetExchangeRate(t *testing.T) {
	var (
		assetCode              = "kBEAM"
		issuer                 = "GC3COBQESTRET2AXK5ADR63L7LOMEZWDPODW4F2Z7Y44TTEOTRBSKXQ3"
		redeemablePricePerUnit = "1.0000000"
		redeemableCollateral   = "VELO"
	)

	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockConfiguration.EXPECT().
			GetDefaultAccount().
			Return(stellarAccountEntity().Address)
		helper.mockDB.EXPECT().
			Get([]byte(stellarAccountEntity().Address)).
			Return(stellarAccountsBytes(), nil)
		helper.mockVelo.EXPECT().
			Client(nil).
			Return(helper.mockVeloClient)
		helper.mockVeloClient.EXPECT().
			GetExchangeRate(context.Background(), gomock.AssignableToTypeOf(&grpc.GetExchangeRateRequest{
				AssetCode: assetCode,
				Issuer:    issuer,
			})).
			Return(&grpc.GetExchangeRateReply{
				AssetCode:              assetCode,
				Issuer:                 issuer,
				RedeemablePricePerUnit: redeemablePricePerUnit,
				RedeemableCollateral:   redeemableCollateral,
			}, nil)

		output, err := helper.logic.GetExchangeRate(&entity.GetExchangeRateInput{
			AssetCode: assetCode,
			Issuer:    issuer,
		})

		assert.NoError(t, err)
		assert.Equal(t, assetCode, output.AssetCode)
		assert.Equal(t, issuer, output.Issuer)
		assert.Equal(t, redeemablePricePerUnit, output.RedeemablePricePerUnit)
		assert.Equal(t, redeemableCollateral, output.RedeemableCollateral)
	})

	t.Run("error, database returns error", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockConfiguration.EXPECT().
			GetDefaultAccount().
			Return(stellarAccountEntity().Address)
		helper.mockDB.EXPECT().
			Get([]byte(stellarAccountEntity().Address)).
			Return(nil, errors.New("some error has occurred"))

		output, err := helper.logic.GetExchangeRate(&entity.GetExchangeRateInput{
			AssetCode: assetCode,
			Issuer:    issuer,
		})

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "failed to get account from db")
	})

	t.Run("error, velo node client returns error", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockConfiguration.EXPECT().
			GetDefaultAccount().
			Return(stellarAccountEntity().Address)
		helper.mockDB.EXPECT().
			Get([]byte(stellarAccountEntity().Address)).
			Return(stellarAccountsBytes(), nil)
		helper.mockVelo.EXPECT().
			Client(nil).
			Return(helper.mockVeloClient)
		helper.mockVeloClient.EXPECT().
			GetExchangeRate(context.Background(), gomock.AssignableToTypeOf(&grpc.GetExchangeRateRequest{
				AssetCode: assetCode,
				Issuer:    issuer,
			})).
			Return(nil, errors.New("failed to get exchange rate"))

		output, err := helper.logic.GetExchangeRate(&entity.GetExchangeRateInput{
			AssetCode: assetCode,
			Issuer:    issuer,
		})

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "failed to get exchange rate")
	})
}
