package collateral_test

import (
	"bou.ke/monkey"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/cobra"
	"github.com/stellar/go/keypair"
	"github.com/velo-protocol/DRSv1/cmd/gvel/layers/commands/collateral"
	"github.com/velo-protocol/DRSv1/cmd/gvel/layers/mocks"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/console"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/mocks"
	"os"
	"testing"
)

type helper struct {
	collateralCommandHandler *collateral.CommandHandler
	mockLogic                *mocks.MockLogic
	mockPrompt               *mockutils.MockPrompt
	mockConfig               *mockutils.MockConfiguration
	mockController           *gomock.Controller
	keyPair                  *keypair.Full
	logHook                  *test.Hook
	tableLogHook             *test.Hook
	done                     func()

	cmd            *cobra.Command
	healthCheckCmd *cobra.Command
	rebalanceCmd   *cobra.Command
}

func initTest(t *testing.T) *helper {
	mockCtrl := gomock.NewController(t)
	mockLogic := mocks.NewMockLogic(mockCtrl)
	mockPrompt := mockutils.NewMockPrompt(mockCtrl)
	mockConfig := mockutils.NewMockConfiguration(mockCtrl)
	keyPair, _ := keypair.Random()

	handler := collateral.NewCommandHandler(mockLogic, mockPrompt, mockConfig)
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
		collateralCommandHandler: handler,
		mockLogic:                mockLogic,
		mockPrompt:               mockPrompt,
		mockController:           mockCtrl,
		keyPair:                  keyPair,
		logHook:                  hook,
		tableLogHook:             tableLogHook,
		done: func() {
			hook.Reset()
			monkey.UnpatchAll()
		},

		cmd:            cmd,
		healthCheckCmd: cmd.Commands()[0],
		rebalanceCmd:   cmd.Commands()[1],
	}
}
