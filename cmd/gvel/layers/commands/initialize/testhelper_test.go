package initialize_test

import (
	"bou.ke/monkey"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus/hooks/test"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/commands/initialize"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/mocks"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/mocks"
	"os"
	"testing"
)

type helper struct {
	commandHandler *initialize.CommandHandler
	loggerHook     *test.Hook
	mockLogic      *mocks.MockLogic
	mockPrompt     *mockutils.MockPrompt
	mockController *gomock.Controller
	done           func()
}

func initTest(t *testing.T) *helper {
	mockCtrl := gomock.NewController(t)
	mockLogic := mocks.NewMockLogic(mockCtrl)
	mockPrompt := mockutils.NewMockPrompt(mockCtrl)

	logger, hook := test.NewNullLogger()
	console.Logger = logger

	monkey.Patch(os.Exit, func(code int) { panic(code) })

	return &helper{
		commandHandler: initialize.NewCommandHandler(mockLogic, mockPrompt),
		mockLogic:      mockLogic,
		mockPrompt:     mockPrompt,
		mockController: mockCtrl,
		loggerHook:     hook,
		done: func() {
			hook.Reset()
			monkey.UnpatchAll()
		},
	}
}
