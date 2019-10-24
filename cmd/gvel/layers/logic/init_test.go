package logic_test

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLogic_Init(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		err := helper.logic.Init("./.velo")
		assert.NoError(t, err)

		_, err = os.Stat("./.velo/config.json")
		assert.NoError(t, err)

		_, err = os.Stat("./.velo/db/account")
		assert.NoError(t, err)
	})

	t.Run("fail, setupConfigFile returns error", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		// force setupConfigFile to return error
		viper.Set("initialized", true)

		err := helper.logic.Init("./.velo")
		assert.Error(t, err)
	})
}
