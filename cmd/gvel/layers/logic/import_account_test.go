package logic_test

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
	"testing"
)

func TestLogic_ImportAccount(t *testing.T) {

	t.Run("success, default flag is set to true", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockDB.EXPECT().
			Get([]byte("GBMD3RER2POVG774HW34A6FYKPTRPXSPHIKUEOSVQZO5RMLCF7FMVI45")).
			Return(nil, errors.New("leveldb: not found"))

		helper.mockConfiguration.EXPECT().
			SetDefaultAccount(gomock.Any()).Return(nil)

		helper.mockDB.EXPECT().
			Save(gomock.Any(), gomock.Any()).Return(nil)

		helper.mockConfiguration.EXPECT().
			GetDefaultAccount().Return("GBVI3QZYXCWQBSGZ4TNJOHDZ5KZYGZOVSE46TVAYJYTMNCGW2PWLWO73")

		helper.mockConfiguration.EXPECT().
			SetDefaultAccount(gomock.Any()).Return(nil)

		output, err := helper.logic.ImportAccount(&entity.ImportAccountInput{
			Passphrase:   "strong_password!",
			SeedKey:      "SDMP3BACNGUQKTNW4NJL2Z4CYG4FFW46MGBU54JNMAHD7LRHKLZTDJBJ",
			SetAsDefault: true,
		})

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.True(t, output.IsDefault)
		assert.NotEmpty(t, output.ImportedKeyPair)
	})

	t.Run("success, default flag is set to false", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockDB.EXPECT().
			Get([]byte("GBMD3RER2POVG774HW34A6FYKPTRPXSPHIKUEOSVQZO5RMLCF7FMVI45")).
			Return(nil, errors.New("leveldb: not found"))

		helper.mockConfiguration.EXPECT().
			SetDefaultAccount(gomock.Any()).Return(nil)

		helper.mockDB.EXPECT().
			Save(gomock.Any(), gomock.Any()).Return(nil)

		helper.mockConfiguration.EXPECT().
			GetDefaultAccount().Return("GBVI3QZYXCWQBSGZ4TNJOHDZ5KZYGZOVSE46TVAYJYTMNCGW2PWLWO73")

		helper.mockConfiguration.EXPECT().
			SetDefaultAccount(gomock.Any()).Return(nil)

		output, err := helper.logic.ImportAccount(&entity.ImportAccountInput{
			Passphrase:   "strong_password!",
			SeedKey:      "SDMP3BACNGUQKTNW4NJL2Z4CYG4FFW46MGBU54JNMAHD7LRHKLZTDJBJ",
			SetAsDefault: false,
		})

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.False(t, output.IsDefault)
		assert.NotEmpty(t, output.ImportedKeyPair)
	})

	t.Run("error, failed to convert seed key to key pair", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockDB.EXPECT().
			Get([]byte("GBMD3RER2POVG774HW34A6FYKPTRPXSPHIKUEOSVQZO5RMLCF7FMVI45")).
			Return(nil, errors.New("leveldb: not found"))

		helper.mockConfiguration.EXPECT().
			SetDefaultAccount(gomock.Any()).Return(nil)

		helper.mockDB.EXPECT().
			Save(gomock.Any(), gomock.Any()).Return(nil)

		output, err := helper.logic.ImportAccount(&entity.ImportAccountInput{
			Passphrase:   "strong_password!",
			SeedKey:      "SEEDKEY",
			SetAsDefault: true,
		})

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "failed to convert seed key to key pair")
	})

	t.Run("error, failed to get account from db", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockDB.EXPECT().
			Get([]byte("GBMD3RER2POVG774HW34A6FYKPTRPXSPHIKUEOSVQZO5RMLCF7FMVI45")).
			Return(nil, errors.New("some error has occurred"))

		output, err := helper.logic.ImportAccount(&entity.ImportAccountInput{
			Passphrase:   "strong_password!",
			SeedKey:      "SDMP3BACNGUQKTNW4NJL2Z4CYG4FFW46MGBU54JNMAHD7LRHKLZTDJBJ",
			SetAsDefault: true,
		})

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "failed to get account from db")
	})

	t.Run("error, account already exist", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockDB.EXPECT().
			Get([]byte("GBMD3RER2POVG774HW34A6FYKPTRPXSPHIKUEOSVQZO5RMLCF7FMVI45")).
			Return(stellarAccountsBytes(), nil)

		output, err := helper.logic.ImportAccount(&entity.ImportAccountInput{
			Passphrase:   "strong_password!",
			SeedKey:      "SDMP3BACNGUQKTNW4NJL2Z4CYG4FFW46MGBU54JNMAHD7LRHKLZTDJBJ",
			SetAsDefault: true,
		})

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Equal(t, fmt.Sprintf("account %s is already exist", "GBMD3RER2POVG774HW34A6FYKPTRPXSPHIKUEOSVQZO5RMLCF7FMVI45"), err.Error())
	})

	t.Run("error, failed to write config file", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockDB.EXPECT().
			Get([]byte("GBMD3RER2POVG774HW34A6FYKPTRPXSPHIKUEOSVQZO5RMLCF7FMVI45")).
			Return(nil, errors.New("leveldb: not found"))

		helper.mockConfiguration.EXPECT().
			SetDefaultAccount(gomock.Any()).
			Return(errors.New("some error has occurred"))

		output, err := helper.logic.ImportAccount(&entity.ImportAccountInput{
			Passphrase:   "strong_password!",
			SeedKey:      "SDMP3BACNGUQKTNW4NJL2Z4CYG4FFW46MGBU54JNMAHD7LRHKLZTDJBJ",
			SetAsDefault: true,
		})

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "failed to write config file")
	})

	t.Run("error, failed to save stellar account", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockDB.EXPECT().
			Get([]byte("GBMD3RER2POVG774HW34A6FYKPTRPXSPHIKUEOSVQZO5RMLCF7FMVI45")).
			Return(nil, errors.New("leveldb: not found"))

		helper.mockConfiguration.EXPECT().
			SetDefaultAccount(gomock.Any()).Return(nil)

		helper.mockDB.EXPECT().
			Save(gomock.Any(), gomock.Any()).Return(errors.New("some error has occurred"))

		output, err := helper.logic.ImportAccount(&entity.ImportAccountInput{
			Passphrase:   "strong_password!",
			SeedKey:      "SDMP3BACNGUQKTNW4NJL2Z4CYG4FFW46MGBU54JNMAHD7LRHKLZTDJBJ",
			SetAsDefault: true,
		})

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "failed to save stellar account")
	})

	t.Run("error, failed to set default account", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockDB.EXPECT().
			Get([]byte("GBMD3RER2POVG774HW34A6FYKPTRPXSPHIKUEOSVQZO5RMLCF7FMVI45")).
			Return(nil, errors.New("leveldb: not found"))

		helper.mockConfiguration.EXPECT().
			SetDefaultAccount(gomock.Any()).Return(nil)

		helper.mockDB.EXPECT().
			Save(gomock.Any(), gomock.Any()).Return(nil)

		helper.mockConfiguration.EXPECT().
			GetDefaultAccount().Return("GBVI3QZYXCWQBSGZ4TNJOHDZ5KZYGZOVSE46TVAYJYTMNCGW2PWLWO73")

		helper.mockConfiguration.EXPECT().
			SetDefaultAccount(gomock.Any()).Return(errors.New("some error has occurred"))

		output, err := helper.logic.ImportAccount(&entity.ImportAccountInput{
			Passphrase:   "strong_password!",
			SeedKey:      "SDMP3BACNGUQKTNW4NJL2Z4CYG4FFW46MGBU54JNMAHD7LRHKLZTDJBJ",
			SetAsDefault: true,
		})

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "failed to write config file")
	})
}
