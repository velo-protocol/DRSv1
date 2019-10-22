package logic_test

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/logic"
	"os"
	"testing"
)

func TestLogic_Init(t *testing.T) {
	t.Run("happy - setConfigFile", func(t *testing.T) {
		lo := logic.NewLogic(nil, nil)

		err := lo.Init("./.velo")
		assert.NoError(t, err)

		_, err = os.Stat("./.velo")
		assert.NoError(t, err)

		err = os.RemoveAll("./.velo")
		assert.NoError(t, err)
	})
}
