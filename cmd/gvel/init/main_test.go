package init

import (
	"github.com/stretchr/testify/assert"
	_default "gitlab.com/velo-labs/cen/cmd/gvel/default"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	t.Run("happy - NewInitCmd", func(t *testing.T) {
		cmd := NewInitCmd()
		assert.NotEmpty(t, cmd)
	})

	t.Run("happy - initRunner", func(t *testing.T) {
		initRunner(nil, nil)

		_, err := os.Stat(_default.DefaultConfigFilePath)
		assert.NoError(t, err)

		err = os.RemoveAll(_default.DefaultConfigFilePath)
		assert.NoError(t, err)
	})

	t.Run("happy - setConfigFile", func(t *testing.T) {
		err := setConfigFile("./.velo")
		assert.NoError(t, err)

		_, err = os.Stat("./.velo")
		assert.NoError(t, err)

		err = os.RemoveAll("./.velo")
		assert.NoError(t, err)
	})
}