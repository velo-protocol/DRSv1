package logic_test

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogic_ListAccount(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		helper := initTest(t)

		mockedStellarAccountsBytes := stellarAccountsBytes()

		helper.mockDB.EXPECT().
			GetAll().
			Return(mockedStellarAccountsBytes, nil)
		helper.mockConfiguration.EXPECT().
			GetDefaultAccount().
			Return("GA...")

		accounts, err := helper.logic.ListAccount()

		assert.NoError(t, err)
		assert.NotEmpty(t, accounts)
		assert.Equal(t, ((*accounts)[0]).EncryptedSeed, []byte("fake-seed"))
		assert.True(t, ((*accounts)[0]).IsDefault)
	})
	t.Run("happy, default account not found", func(t *testing.T) {
		helper := initTest(t)

		mockedStellarAccountsBytes := stellarAccountsBytes()

		helper.mockDB.EXPECT().
			GetAll().
			Return(mockedStellarAccountsBytes, nil)
		helper.mockConfiguration.EXPECT().
			GetDefaultAccount().
			Return("")

		accounts, err := helper.logic.ListAccount()

		assert.NoError(t, err)
		assert.NotEmpty(t, accounts)
		assert.Equal(t, ((*accounts)[0]).EncryptedSeed, []byte("fake-seed"))
		assert.False(t, ((*accounts)[0]).IsDefault)
	})

	t.Run("error - failed to load accounts from db", func(t *testing.T) {
		helper := initTest(t)

		helper.mockDB.EXPECT().GetAll().Return(nil, errors.New("error here"))

		accounts, err := helper.logic.ListAccount()

		assert.Empty(t, accounts)
		assert.Error(t, err)
	})

	t.Run("error - failed to unmarshal stored data to entity", func(t *testing.T) {
		helper := initTest(t)

		helper.mockDB.EXPECT().
			GetAll().
			Return([][]byte{[]byte("fake")}, nil)
		helper.mockConfiguration.EXPECT().
			GetDefaultAccount().
			Return("GA...")
		accounts, err := helper.logic.ListAccount()

		assert.Error(t, err)
		assert.Empty(t, accounts)
	})
}
