package logic_test

import (
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/velo-protocol/DRSv1/cmd/gvel/entity"
	"testing"
)

func TestLogic_CreateAccount(t *testing.T) {
	t.Run("success, default flag is set to true", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockStellar.EXPECT().
			GetFreeLumens(gomock.Any()).Return(nil)

		helper.mockDB.EXPECT().
			Save(gomock.Any(), gomock.Any()).Return(nil)

		helper.mockConfiguration.EXPECT().
			GetDefaultAccount().Return("GBVI3QZYXCWQBSGZ4TNJOHDZ5KZYGZOVSE46TVAYJYTMNCGW2PWLWO73")

		helper.mockConfiguration.EXPECT().
			SetDefaultAccount(gomock.Any()).Return(nil)

		output, err := helper.logic.CreateAccount(&entity.CreateAccountInput{
			Passphrase:          "strong_password!",
			SetAsDefaultAccount: true,
		})

		assert.NoError(t, err)
		assert.NotNil(t, output.GeneratedKeyPair)
		assert.True(t, output.IsDefault)
	})

	t.Run("success, default flag is set to false, but no default account is defined before", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockStellar.EXPECT().
			GetFreeLumens(gomock.Any()).Return(nil)

		helper.mockDB.EXPECT().
			Save(gomock.Any(), gomock.Any()).Return(nil)

		helper.mockConfiguration.EXPECT().
			GetDefaultAccount().Return("")

		helper.mockConfiguration.EXPECT().
			SetDefaultAccount(gomock.Any()).Return(nil)

		output, err := helper.logic.CreateAccount(&entity.CreateAccountInput{
			Passphrase:          "strong_password!",
			SetAsDefaultAccount: false,
		})

		assert.NoError(t, err)
		assert.NotNil(t, output.GeneratedKeyPair)
		assert.True(t, output.IsDefault)
	})

	t.Run("error - failed to get free lumens from friendbot", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockStellar.EXPECT().
			GetFreeLumens(gomock.Any()).Return(errors.New("error happens here"))

		_, err := helper.logic.CreateAccount(&entity.CreateAccountInput{
			Passphrase: "strong_password!",
		})
		assert.Error(t, err)
	})

	t.Run("error - failed to save an account to level db", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockStellar.EXPECT().
			GetFreeLumens(gomock.Any()).Return(nil)

		helper.mockDB.EXPECT().
			Save(gomock.Any(), gomock.Any()).Return(errors.New("error happens here"))

		_, err := helper.logic.CreateAccount(&entity.CreateAccountInput{
			Passphrase: "strong_password!",
		})
		assert.Error(t, err)
	})

	t.Run("error - failed to write config", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockStellar.EXPECT().
			GetFreeLumens(gomock.Any()).Return(nil)

		helper.mockDB.EXPECT().
			Save(gomock.Any(), gomock.Any()).Return(errors.New("error happens here"))

		helper.mockConfiguration.EXPECT().
			GetDefaultAccount().Return("GBVI3QZYXCWQBSGZ4TNJOHDZ5KZYGZOVSE46TVAYJYTMNCGW2PWLWO73")

		helper.mockConfiguration.EXPECT().
			SetDefaultAccount(gomock.Any()).Return(errors.New("error happens here"))

		_, err := helper.logic.CreateAccount(&entity.CreateAccountInput{
			Passphrase:          "strong_password!",
			SetAsDefaultAccount: true,
		})
		assert.Error(t, err)
	})
}
