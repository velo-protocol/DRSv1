package collateral_test

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/velo-protocol/DRSv1/cmd/gvel/entity"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/console"
	"strings"
	"testing"
)

func TestCommandHandler_GetHealthCheck(t *testing.T) {

	t.Run("success", func(t *testing.T) {

		helper := initTest(t)
		defer helper.done()

		helper.mockLogic.EXPECT().
			GetCollateralHealthCheck().
			Return(&entity.GetCollateralHealthCheckOutput{
				Asset:          "kBEAM (GD4K...)",
				PoolAmount:     "1.0000000",
				RequiredAmount: "1.5000000",
			}, nil)

		helper.collateralCommandHandler.GetHealthCheck(helper.healthCheckCmd, nil)
		lines := strings.Split(helper.tableLogHook.LastEntry().Message, "\n")

		assert.Contains(t, lines[1], "ASSET")
		assert.Contains(t, lines[1], "COLLATERAL POOL")
		assert.Contains(t, lines[1], "REQUIRED COLLATERAL")

		assert.Contains(t, lines[3], "kBEAM (GD4K...)")
		assert.Contains(t, lines[3], "1.0000000")
		assert.Contains(t, lines[3], "1.5000000")
	})

	t.Run("fail, get health check returns error", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockLogic.EXPECT().
			GetCollateralHealthCheck().
			Return(nil, errors.New("some error has occurred"))

		assert.PanicsWithValue(t, console.ExitError, func() {
			helper.collateralCommandHandler.GetHealthCheck(helper.healthCheckCmd, nil)
		})

		logEntries := helper.logHook.Entries

		assert.Equal(t, errors.New("some error has occurred").Error(), logEntries[0].Message)
	})
}
