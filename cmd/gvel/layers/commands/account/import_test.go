package account_test

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/velo-protocol/DRSv1/cmd/gvel/entity"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/console"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/validation"
	"testing"
)

func TestCommandHandler_ImportAccount(t *testing.T) {
	t.Run("success, as default true", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		kp := helper.keyPair
		_ = helper.importCmd.Flags().Set("default", "true")

		helper.mockPrompt.EXPECT().
			RequestHiddenString("Please enter seed key of the address", gomock.AssignableToTypeOf(validation.ValidateSeedKey)).
			Return("GB...")

		helper.mockPrompt.EXPECT().
			RequestPassphrase().
			Return("strong_password!")

		helper.mockLogic.EXPECT().
			ImportAccount(&entity.ImportAccountInput{
				Passphrase:   "strong_password!",
				SeedKey:      "GB...",
				SetAsDefault: true,
			}).
			Return(&entity.ImportAccountOutput{
				ImportedKeyPair: kp,
				IsDefault:       true,
			}, nil)

		helper.accountCommandHandler.ImportAccount(helper.importCmd, nil)

		logEntries := helper.logHook.AllEntries()
		msg := fmt.Sprintf("Account with address %s is now the default account. Please keep your passphrase safe.", kp.Address())
		assert.Equal(t, msg, logEntries[0].Message)
	})

	t.Run("success, as default false", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		kp := helper.keyPair
		_ = helper.importCmd.Flags().Set("default", "false")

		helper.mockPrompt.EXPECT().
			RequestHiddenString("Please enter seed key of the address", gomock.AssignableToTypeOf(validation.ValidateSeedKey)).
			Return("GB...")

		helper.mockPrompt.EXPECT().
			RequestPassphrase().
			Return("strong_password!")

		helper.mockLogic.EXPECT().
			ImportAccount(&entity.ImportAccountInput{
				Passphrase:   "strong_password!",
				SeedKey:      "GB...",
				SetAsDefault: false,
			}).
			Return(&entity.ImportAccountOutput{
				ImportedKeyPair: kp,
				IsDefault:       false,
			}, nil)

		helper.accountCommandHandler.ImportAccount(helper.importCmd, nil)

		logEntries := helper.logHook.AllEntries()
		msg := fmt.Sprintf("Add account with address %s to gvel. Please keep your passphrase safe.", kp.Address())
		assert.Equal(t, msg, logEntries[0].Message)
	})

	t.Run("error, logic.ImportAccount returns error", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		_ = helper.importCmd.Flags().Set("default", "false")

		helper.mockPrompt.EXPECT().
			RequestHiddenString("Please enter seed key of the address", gomock.AssignableToTypeOf(validation.ValidateSeedKey)).
			Return("GB...")

		helper.mockPrompt.EXPECT().
			RequestPassphrase().
			Return("strong_password!")

		helper.mockLogic.EXPECT().
			ImportAccount(&entity.ImportAccountInput{
				Passphrase:   "strong_password!",
				SeedKey:      "GB...",
				SetAsDefault: false,
			}).
			Return(nil, errors.New("some error has occurred"))

		assert.PanicsWithValue(t, console.ExitError, func() {
			helper.accountCommandHandler.ImportAccount(helper.importCmd, nil)
		})

		assert.Equal(t, "some error has occurred", helper.logHook.LastEntry().Message)
	})
}
