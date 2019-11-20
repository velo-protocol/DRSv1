package account_test

import (
	"bou.ke/monkey"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/cobra"
	"github.com/stellar/go/keypair"
	"github.com/velo-protocol/DRSv1/cmd/gvel/layers/commands/account"
	"github.com/velo-protocol/DRSv1/cmd/gvel/layers/mocks"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/console"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/mocks"
	"os"

	"testing"
)

type helper struct {
	accountCommandHandler *account.CommandHandler
	mockLogic             *mocks.MockLogic
	mockPrompt            *mockutils.MockPrompt
	mockConfig            *mockutils.MockConfiguration
	mockController        *gomock.Controller
	keyPair               *keypair.Full
	logHook               *test.Hook
	tableLogHook          *test.Hook
	done                  func()

	cmd        *cobra.Command
	createCmd  *cobra.Command
	listCmd    *cobra.Command
	defaultCmd *cobra.Command
	importCmd  *cobra.Command
	exportCmd  *cobra.Command
}

func initTest(t *testing.T) *helper {
	mockCtrl := gomock.NewController(t)
	mockLogic := mocks.NewMockLogic(mockCtrl)
	mockPrompt := mockutils.NewMockPrompt(mockCtrl)
	mockConfig := mockutils.NewMockConfiguration(mockCtrl)
	keyPair, _ := keypair.Random()

	handler := account.NewCommandHandler(mockLogic, mockPrompt, mockConfig)
	cmd := handler.Command()

	// logger
	logger, hook := test.NewNullLogger()
	console.Logger = logger

	// table logger
	tableLogger, tableLogHook := test.NewNullLogger()
	console.TableLogger = tableLogger

	// overwrite os.Exit
	monkey.Patch(os.Exit, func(code int) { panic(code) })

	// to omit what loader print
	console.DefaultLoadWriter = console.Logger.Out

	return &helper{
		accountCommandHandler: handler,
		mockLogic:             mockLogic,
		mockPrompt:            mockPrompt,
		mockController:        mockCtrl,
		keyPair:               keyPair,
		logHook:               hook,
		tableLogHook:          tableLogHook,
		done: func() {
			hook.Reset()
			monkey.UnpatchAll()
		},

		cmd:        cmd,
		createCmd:  cmd.Commands()[0],
		defaultCmd: cmd.Commands()[1],
		exportCmd:  cmd.Commands()[2],
		importCmd:  cmd.Commands()[3],
		listCmd:    cmd.Commands()[4],
	}
}
