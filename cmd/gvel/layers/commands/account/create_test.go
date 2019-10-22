package account_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommandHandler_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockPrompt.EXPECT().
			RequestPassphrase().
			Return("strong_password!")

		helper.mockLogic.EXPECT().
			CreateAccount("strong_password!").
			Return(helper.keyPair, nil)

		helper.accountCommandHandler.Create(nil, nil)

		logEntries := helper.logHook.AllEntries()
		assert.Equal(t, "generating a new stellar account", logEntries[0].Message)
		assert.Equal(t, fmt.Sprintf("%s has been created\n", helper.keyPair.Address()), logEntries[1].Message)
	})
}
