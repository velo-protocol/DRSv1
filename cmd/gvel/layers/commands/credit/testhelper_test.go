package credit_test

import (
	"bou.ke/monkey"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/cobra"
	"github.com/stellar/go/keypair"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/commands/credit"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/mocks"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/mocks"
	"os"

	"testing"
)

type helper struct {
	creditCommandHandler *credit.CommandHandler
	mockLogic            *mocks.MockLogic
	mockPrompt           *mockutils.MockPrompt
	mockConfig           *mockutils.MockConfiguration
	mockController       *gomock.Controller
	keyPair              *keypair.Full
	logHook              *test.Hook
	tableLogHook         *test.Hook
	done                 func()

	cmd      *cobra.Command
	setupCmd *cobra.Command
}

func initTest(t *testing.T) *helper {
	mockCtrl := gomock.NewController(t)
	mockLogic := mocks.NewMockLogic(mockCtrl)
	mockPrompt := mockutils.NewMockPrompt(mockCtrl)
	mockConfig := mockutils.NewMockConfiguration(mockCtrl)
	keyPair, _ := keypair.Random()

	handler := credit.NewCommandHandler(mockLogic, mockPrompt, mockConfig)
	cmd := handler.Command()

	// logger
	logger, hook := test.NewNullLogger()
	console.Logger = logger

	// table logger
	tableLogger, tableLogHook := test.NewNullLogger()
	console.TableLogger = tableLogger

	// overwrite os.Exit
	monkey.Patch(os.Exit, func(code int) { panic(code) })

	return &helper{
		creditCommandHandler: handler,
		mockLogic:            mockLogic,
		mockPrompt:           mockPrompt,
		mockController:       mockCtrl,
		keyPair:              keyPair,
		logHook:              hook,
		tableLogHook:         tableLogHook,
		done: func() {
			hook.Reset()
			monkey.UnpatchAll()
		},

		cmd:      cmd,
		setupCmd: cmd.Commands()[0],
	}
}
