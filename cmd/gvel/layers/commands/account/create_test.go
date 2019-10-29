package account_test

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
	"testing"
)

func TestCommandHandler_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		_ = helper.createCmd.Flags().Set("default", "true")

		helper.mockPrompt.EXPECT().
			RequestPassphrase().
			Return("strong_password!")

		helper.mockLogic.EXPECT().
			CreateAccount(&entity.CreateAccountInput{
				Passphrase:          "strong_password!",
				SetAsDefaultAccount: true,
			}).
			Return(&entity.CreateAccountOutput{
				GeneratedKeyPair: helper.keyPair,
				IsDefault:        true,
			}, nil)

		helper.accountCommandHandler.Create(helper.createCmd, nil)

		logEntries := helper.logHook.AllEntries()
		assert.Contains(t, logEntries[0].Message, fmt.Sprintf("A new account is created with address"))
		assert.Contains(t, logEntries[0].Message, fmt.Sprintf("Please remember to keep your passphrase safe. You will not be able to recover this passphrase."))
	})
	t.Run("fail, logic.CreateAccount returns error", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockPrompt.EXPECT().
			RequestPassphrase().
			Return("strong_password!")

		helper.mockLogic.EXPECT().
			CreateAccount(&entity.CreateAccountInput{
				Passphrase:          "strong_password!",
				SetAsDefaultAccount: false,
			}).
			Return(nil, errors.New("some error has occurred"))

		assert.PanicsWithValue(t, console.ExitError, func() {
			helper.accountCommandHandler.Create(helper.createCmd, nil)
		})

		assert.Equal(t, "some error has occurred", helper.logHook.LastEntry().Message)
	})
}
