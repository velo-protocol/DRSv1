package credit_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
	"testing"
)

func TestCommandHandler_Setup(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockPrompt.EXPECT().
			RequestString("Please input asset code", nil).
			Return("vTHB")
		helper.mockPrompt.EXPECT().
			RequestString("Please input pegged value", nil).
			Return("1")
		helper.mockPrompt.EXPECT().
			RequestString("Please input pegged currency", nil).
			Return("THB")
		helper.mockPrompt.EXPECT().
			RequestHiddenString("Please enter passphrase", nil).
			Return("password")
		helper.mockLogic.EXPECT().
			SetupCredit(gomock.AssignableToTypeOf(&entity.SetupCreditInput{})).
			Return(&entity.SetupCreditOutput{
				AssetCode:      "vTHB",
				PeggedValue:    "1",
				PeggedCurrency: "THB",
				SourceAddress:  "GA...",
				TxResult:       &horizon.TransactionSuccess{Hash: "AAA..."},
			}, nil)

		helper.creditCommandHandler.Setup(helper.setupCmd, nil)

		logEntries := helper.logHook.Entries
		assert.Equal(t, "Stable credit vTHB set up for account GA... successfully.", logEntries[0].Message)
		assert.Equal(t, "Stellar Transaction Hash ✉️ : AAA...", logEntries[1].Message)
	})

	t.Run("error, logic.SetupCredit returns error", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockPrompt.EXPECT().
			RequestString("Please input asset code", nil).
			Return("vTHB")
		helper.mockPrompt.EXPECT().
			RequestString("Please input pegged value", nil).
			Return("1")
		helper.mockPrompt.EXPECT().
			RequestString("Please input pegged currency", nil).
			Return("THB")
		helper.mockPrompt.EXPECT().
			RequestHiddenString("Please enter passphrase", nil).
			Return("password")
		helper.mockLogic.EXPECT().
			SetupCredit(gomock.AssignableToTypeOf(&entity.SetupCreditInput{})).
			Return(nil, errors.New("some error has occurred"))

		assert.PanicsWithValue(t, console.ExitError, func() {
			helper.creditCommandHandler.Setup(helper.setupCmd, nil)
		})
	})
}
