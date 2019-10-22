package logic_test

import (
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogic_CreateAccount(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		helper := initTest(t)

		helper.mockFriendBot.EXPECT().
			GetFreeLumens(gomock.Any()).Return(nil)

		helper.mockDB.EXPECT().
			Save(gomock.Any(), gomock.Any()).Return(nil)

		kp, err := helper.logic.CreateAccount("1234")

		assert.NoError(t, err)
		assert.NotEmpty(t, kp)
	})

	t.Run("error - failed to get free lumens from friendbot", func(t *testing.T) {
		helper := initTest(t)

		helper.mockFriendBot.EXPECT().
			GetFreeLumens(gomock.Any()).Return(errors.New("error happens here"))

		kp, err := helper.logic.CreateAccount("1234")

		assert.Error(t, err)
		assert.Empty(t, kp)
	})

	t.Run("error - failed to save an account to level db", func(t *testing.T) {
		helper := initTest(t)

		helper.mockFriendBot.EXPECT().
			GetFreeLumens(gomock.Any()).Return(nil)

		helper.mockDB.EXPECT().
			Save(gomock.Any(), gomock.Any()).Return(errors.New("error happens here"))

		kp, err := helper.logic.CreateAccount("1234")

		assert.Error(t, err)
		assert.Empty(t, kp)
	})
}
