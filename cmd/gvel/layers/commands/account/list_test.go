package account_test

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
	"strings"
	"testing"
)

func TestCommandHandler_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockLogic.EXPECT().
			ListAccount().
			Return(&[]entity.StellarAccount{
				{
					Address:   "GA...",
					IsDefault: false,
				},
				{
					Address:   "GB...",
					IsDefault: true,
				},
				{
					Address:   "GC...",
					IsDefault: false,
				},
			}, nil)

		helper.accountCommandHandler.List(helper.listCmd, nil)
		lines := strings.Split(helper.tableLogHook.LastEntry().Message, "\n")

		assert.Contains(t, lines[1], "INDEX")
		assert.Contains(t, lines[1], "ADDRESS")
		assert.Contains(t, lines[1], "DEFAULT")

		assert.Contains(t, lines[3], "GA...")
		assert.Contains(t, lines[4], "GB...")
		assert.Contains(t, lines[5], "GC...")

		assert.Contains(t, lines[3], "false")
		assert.Contains(t, lines[4], "true")
		assert.Contains(t, lines[5], "false")
	})
	t.Run("fail, list account returns error", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockLogic.EXPECT().
			ListAccount().
			Return(nil, errors.New("some error has occurred"))

		assert.PanicsWithValue(t, console.ExitError, func() {
			helper.accountCommandHandler.List(helper.listCmd, nil)
		})
	})
}
