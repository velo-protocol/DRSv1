package initialize_test

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/cmd/gvel/constants"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
	"testing"
)

func TestCommandHandler_Init(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockLogic.EXPECT().
			Init(constants.DefaultConfigFilePath).
			Return(nil)

		helper.commandHandler.Init(nil, nil)

		logEntries := helper.loggerHook.AllEntries()
		assert.Equal(t, "gvel has been initialized\n", logEntries[0].Message)
		assert.Equal(t, fmt.Sprintf("using config file at: %s\n", constants.DefaultConfigFilePath), logEntries[1].Message)
	})

	t.Run("fail, config already exist", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		viper.Set("initialized", true)

		assert.PanicsWithValue(t, console.ExitError, func() {
			helper.commandHandler.Init(nil, nil)
		})
		assert.Equal(t, fmt.Sprintf("gvel has already been initialized, configuration can be found at %s", constants.DefaultConfigFilePath), helper.loggerHook.LastEntry().Message)
	})

	t.Run("fail, logic returns error", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockLogic.EXPECT().
			Init(constants.DefaultConfigFilePath).
			Return(errors.New("some error has occurred"))

		assert.PanicsWithValue(t, console.ExitError, func() {
			helper.commandHandler.Init(nil, nil)
		})
	})
}
