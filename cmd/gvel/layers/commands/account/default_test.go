package account_test

import (
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/validation"
	"testing"
)

func TestCommandHandler_Default(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockPrompt.EXPECT().
			RequestString("Please input public address to be made default", gomock.AssignableToTypeOf(validation.ValidateStellarAddress)).
			Return("GA...")

		helper.mockLogic.EXPECT().
			SetDefaultAccount(&entity.SetDefaultAccountInput{
				Account: "GA...",
			}).
			Return(&entity.SetDefaultAccountOutput{
				Account: "GA...",
			}, nil)

		helper.accountCommandHandler.Default(helper.defaultCmd, nil)

		logEntries := helper.logHook.AllEntries()
		assert.Equal(t, "GA... is now set as the default account for signing transaction.", logEntries[0].Message)
	})
	t.Run("fail, logic.SetDefaultAccount returns error", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockPrompt.EXPECT().
			RequestString("Please input public address to be made default", gomock.AssignableToTypeOf(validation.ValidateStellarAddress)).
			Return("GA...")

		helper.mockLogic.EXPECT().
			SetDefaultAccount(&entity.SetDefaultAccountInput{
				Account: "GA...",
			}).
			Return(nil, errors.New("some error has occurred"))

		assert.PanicsWithValue(t, console.ExitError, func() {
			helper.accountCommandHandler.Default(helper.defaultCmd, nil)
		})
	})
}
