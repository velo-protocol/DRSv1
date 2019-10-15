package init

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/cmd/gvel/constant"
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

		_, err := os.Stat(constant.DefaultConfigFilePath)
		assert.NoError(t, err)

		err = os.RemoveAll(constant.DefaultConfigFilePath)
		assert.NoError(t, err)
	})
}
