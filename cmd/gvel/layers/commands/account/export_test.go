package account_test

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/validation"
	"testing"
)

func TestCommandHandler_ExportAccount(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()
		kp := helper.keyPair

		helper.mockPrompt.EXPECT().
			RequestString("Please input the public key you want to export", gomock.AssignableToTypeOf(validation.ValidateStellarAddress)).
			Return("GB...")
		helper.mockPrompt.EXPECT().
			RequestHiddenString("🔑 Please input the passphrase of the account", nil).
			Return("strong_password!")

		helper.mockLogic.EXPECT().
			ExportAccount(&entity.ExportAccountInput{
				PublicKey:  "GB...",
				Passphrase: "strong_password!",
			}).
			Return(&entity.ExportAccountOutput{
				ExportedKeyPair: kp,
			}, nil)

		helper.accountCommandHandler.ExportAccount(helper.exportCmd, nil)

		logEntries := helper.logHook.AllEntries()
		publicKey := fmt.Sprintf("Your public key is: %s", kp.Address())
		seedKey := fmt.Sprintf("Your seed key is: %s", kp.Seed())

		assert.Equal(t, publicKey, logEntries[0].Message)
		assert.Equal(t, seedKey, logEntries[1].Message)
	})

	t.Run("error, logic.ExportAccount returns error", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockPrompt.EXPECT().
			RequestString("Please input the public key you want to export", gomock.AssignableToTypeOf(validation.ValidateStellarAddress)).
			Return("GB...")

		helper.mockPrompt.EXPECT().
			RequestHiddenString("🔑 Please input the passphrase of the account", nil).
			Return("strong_password!")

		helper.mockLogic.EXPECT().
			ExportAccount(&entity.ExportAccountInput{
				PublicKey:  "GB...",
				Passphrase: "strong_password!",
			}).
			Return(nil, errors.New("some error has occurred"))

		assert.PanicsWithValue(t, console.ExitError, func() {
			helper.accountCommandHandler.ExportAccount(helper.exportCmd, nil)
		})

		assert.Equal(t, "some error has occurred", helper.logHook.LastEntry().Message)
	})
}
