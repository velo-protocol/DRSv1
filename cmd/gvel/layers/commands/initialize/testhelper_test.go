package initialize_test

import (
	"bou.ke/monkey"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/velo-protocol/DRSv1/cmd/gvel/layers/commands/initialize"
	"github.com/velo-protocol/DRSv1/cmd/gvel/layers/mocks"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/console"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/mocks"
	"os"
	"testing"
)

type helper struct {
	commandHandler *initialize.CommandHandler
	loggerHook     *test.Hook
	mockLogic      *mocks.MockLogic
	mockPrompt     *mockutils.MockPrompt
	mockConfig     *mockutils.MockConfiguration
	mockController *gomock.Controller
	done           func()
}

func initTest(t *testing.T) *helper {
	mockCtrl := gomock.NewController(t)
	mockLogic := mocks.NewMockLogic(mockCtrl)
	mockPrompt := mockutils.NewMockPrompt(mockCtrl)
	mockConfig := mockutils.NewMockConfiguration(mockCtrl)

	logger, hook := test.NewNullLogger()
	console.Logger = logger

	monkey.Patch(os.Exit, func(code int) { panic(code) })

	// to omit what loader print
	console.DefaultLoadWriter = console.Logger.Out

	return &helper{
		commandHandler: initialize.NewCommandHandler(mockLogic, mockPrompt, mockConfig),
		mockLogic:      mockLogic,
		mockPrompt:     mockPrompt,
		mockConfig:     mockConfig,
		mockController: mockCtrl,
		loggerHook:     hook,
		done: func() {
			hook.Reset()
			monkey.UnpatchAll()
		},
	}
}
