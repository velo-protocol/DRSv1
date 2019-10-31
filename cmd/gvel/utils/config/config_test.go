package config

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/cmd/gvel/constants"
	"os"
	"testing"
)

func TestInitConfigFile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		helper := initTest()
		defer helper.done()

		err := helper.config.InitConfigFile("./.gvel")
		assert.NoError(t, err)

		_, err = os.Stat("./.gvel/config.json")
		assert.NoError(t, err)

		assert.True(t, viper.GetBool("initialized"))
		assert.Equal(t, ".gvel/db/account", viper.GetString("accountDbPath"))
		assert.Equal(t, constants.DefaultHorizonUrl, viper.GetString("horizonUrl"))
	})

	t.Run("success, config file is already exist, no error", func(t *testing.T) {
		helper := initTest()
		defer helper.done()

		helper.config.viper.Set("initialized", true)

		err := helper.config.InitConfigFile("./.gvel")
		assert.NoError(t, err)
		assert.Equal(t, helper.loggerHook.LastEntry().Message, "config file found")
	})
}
