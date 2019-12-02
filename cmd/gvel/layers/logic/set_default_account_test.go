package logic_test

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/velo-protocol/DRSv1/cmd/gvel/entity"
	"testing"
)

func TestLogic_SetDefaultAccount(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockDB.EXPECT().
			Get([]byte("GA...")).
			Return(stellarAccountsBytes(), nil)
		helper.mockConfiguration.EXPECT().
			SetDefaultAccount("GA...").
			Return(nil)

		output, err := helper.logic.SetDefaultAccount(&entity.SetDefaultAccountInput{
			Account: "GA...",
		})

		assert.NoError(t, err)
		assert.Equal(t, "GA...", output.Account)
	})
	t.Run("error, database returns error", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockDB.EXPECT().
			Get([]byte("GA...")).
			Return(nil, errors.New("some error has occurred"))

		_, err := helper.logic.SetDefaultAccount(&entity.SetDefaultAccountInput{
			Account: "GA...",
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "address GA... is not found in gvel")
	})
	t.Run("error, config returns error", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockDB.EXPECT().
			Get([]byte("GA...")).
			Return(stellarAccountsBytes(), nil)
		helper.mockConfiguration.EXPECT().
			SetDefaultAccount("GA...").
			Return(errors.New("some error has occurred"))

		_, err := helper.logic.SetDefaultAccount(&entity.SetDefaultAccountInput{
			Account: "GA...",
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to write config file")
	})
}
