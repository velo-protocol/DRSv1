package account_test

import (
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stellar/go/keypair"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/commands/account"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/mocks"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/mocks"

	"testing"
)

type helper struct {
	accountCommandHandler *account.CommandHandler
	mockLogic             *mocks.MockLogic
	mockPrompt            *mockutils.MockPrompt
	mockController        *gomock.Controller
	keyPair               *keypair.Full
	logHook               *test.Hook
	done                  func()
}

func initTest(t *testing.T) *helper {
	mockCtrl := gomock.NewController(t)
	mockLogic := mocks.NewMockLogic(mockCtrl)
	mockPrompt := mockutils.NewMockPrompt(mockCtrl)
	keyPair, _ := keypair.Random()

	logger, hook := test.NewNullLogger()
	console.Logger = logger

	return &helper{
		accountCommandHandler: account.NewCommandHandler(mockLogic, mockPrompt),
		mockLogic:             mockLogic,
		mockPrompt:            mockPrompt,
		mockController:        mockCtrl,
		keyPair:               keyPair,
		logHook:               hook,
		done: func() {
			hook.Reset()
		},
	}
}
