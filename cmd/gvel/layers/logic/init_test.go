package logic_test

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogic_Init(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockConfiguration.EXPECT().
			InitConfigFile("./").
			Return(nil)

		helper.mockConfiguration.EXPECT().
			GetAccountDbPath().
			Return("./db/accounts")

		helper.mockDB.EXPECT().
			Init("./db/accounts").
			Return(nil)

		err := helper.logic.Init("./")
		assert.NoError(t, err)
	})

	t.Run("fail, setupConfigFile returns error", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockConfiguration.EXPECT().
			InitConfigFile("./").
			Return(errors.New("some error has occurred"))

		err := helper.logic.Init("./")
		assert.Error(t, err)
	})

	t.Run("fail, database init returns error", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockConfiguration.EXPECT().
			InitConfigFile("./").
			Return(nil)

		helper.mockConfiguration.EXPECT().
			GetAccountDbPath().
			Return("./db/accounts")

		helper.mockDB.EXPECT().
			Init("./db/accounts").
			Return(errors.New("some error has occurred"))

		err := helper.logic.Init("./")
		assert.Error(t, err)
	})
}
