package credit_test

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
	"testing"
)

func TestCommandHandler_GetExchangeRate(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockPrompt.EXPECT().
			RequestString("Please input asset code of the stable credit", nil).
			Return("kBEAM")
		helper.mockPrompt.EXPECT().
			RequestString("Please input issuing account of the stable credit", nil).
			Return("GA...")

		helper.mockLogic.EXPECT().
			GetExchangeRate(gomock.AssignableToTypeOf(&entity.GetExchangeRateInput{})).
			Return(&entity.GetExchangeRateOutput{
				AssetCode:              "kBEAM",
				Issuer:                 "GA...",
				RedeemablePricePerUnit: "2",
				RedeemableCollateral:   "VELO",
			}, nil)

		helper.creditCommandHandler.GetExchangeRate(helper.mintCmd, nil)

		logEntries := helper.logHook.Entries
		assert.Equal(t, fmt.Sprintf("You will get %s %s for 1 %s.", "2", "VELO", "kBEAM"), logEntries[0].Message)
	})

	t.Run("error, logic.GetExchangeRate returns error", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockPrompt.EXPECT().
			RequestString("Please input asset code of the stable credit", nil).
			Return("kBEAM")
		helper.mockPrompt.EXPECT().
			RequestString("Please input issuing account of the stable credit", nil).
			Return("GA...")

		helper.mockLogic.EXPECT().
			GetExchangeRate(gomock.AssignableToTypeOf(&entity.GetExchangeRateInput{})).
			Return(nil, errors.New("some error has occurred"))

		assert.PanicsWithValue(t, console.ExitError, func() {
			helper.creditCommandHandler.GetExchangeRate(helper.mintCmd, nil)
		})

	})

}
