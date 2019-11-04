package logic_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	spec "gitlab.com/velo-labs/cen/grpc"
	"testing"
)

func TestLogic_GetCollateralHealthCheck(t *testing.T) {

	var (
		assetCode      = "kBEAM"
		assetIssuer    = "GC3COBQESTRET2AXK5ADR63L7LOMEZWDPODW4F2Z7Y44TTEOTRBSKXQ3"
		poolAmount     = "1.0000000"
		requiredAmount = "2.0000000"
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
			GetCollateralHealthCheck(context.Background(), gomock.AssignableToTypeOf(&spec.GetCollateralHealthCheckRequest{})).
			Return(&spec.GetCollateralHealthCheckReply{
				AssetCode:      assetCode,
				AssetIssuer:    assetIssuer,
				PoolAmount:     poolAmount,
				RequiredAmount: requiredAmount,
			}, nil)

		output, err := helper.logic.GetCollateralHealthCheck()
		assert.NoError(t, err)
		assert.Equal(t, assetCode, output.AssetCode)
		assert.Equal(t, assetIssuer, output.AssetIssuer)
		assert.Equal(t, poolAmount, output.PoolAmount)
		assert.Equal(t, requiredAmount, output.RequiredAmount)

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

		output, err := helper.logic.GetCollateralHealthCheck()

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
			GetCollateralHealthCheck(context.Background(), gomock.AssignableToTypeOf(&spec.GetCollateralHealthCheckRequest{})).
			Return(nil, errors.New("failed to get collateral health check"))

		output, err := helper.logic.GetCollateralHealthCheck()

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "failed to get collateral health check")
	})
}
