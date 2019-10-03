package stellar_test

import (
	"errors"
	"fmt"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/testhelpers"
	"testing"
)

func TestRepo_GetAccounts(t *testing.T) {
	testhelpers.InitEnv()

	t.Run("success", func(t *testing.T) {
		helper := initTest()

		helper.mockedHorizonClient.
			On("AccountDetail", horizonclient.AccountRequest{
				AccountID: testhelpers.PublicKey1,
			}).
			Return(horizon.Account{
				AccountID: testhelpers.PublicKey1,
			}, nil)
		helper.mockedHorizonClient.
			On("AccountDetail", horizonclient.AccountRequest{
				AccountID: testhelpers.PublicKey2,
			}).
			Return(horizon.Account{
				AccountID: testhelpers.PublicKey2,
			}, nil)
		helper.mockedHorizonClient.
			On("AccountDetail", horizonclient.AccountRequest{
				AccountID: testhelpers.PublicKey3,
			}).
			Return(horizon.Account{
				AccountID: testhelpers.PublicKey3,
			}, nil)

		accounts, err := helper.repo.GetAccounts(
			testhelpers.PublicKey1,
			testhelpers.PublicKey2,
			testhelpers.PublicKey3,
		)

		assert.NoError(t, err)
		assert.Equal(t, testhelpers.PublicKey1, accounts[0].AccountID)
		assert.Equal(t, testhelpers.PublicKey2, accounts[1].AccountID)
		assert.Equal(t, testhelpers.PublicKey3, accounts[2].AccountID)

		helper.mockedHorizonClient.
			AssertNumberOfCalls(t, "AccountDetail", 3)
	})
	t.Run("error, fail to get account detail of public key 3", func(t *testing.T) {
		helper := initTest()

		helper.mockedHorizonClient.
			On("AccountDetail", horizonclient.AccountRequest{
				AccountID: testhelpers.PublicKey1,
			}).
			Return(horizon.Account{
				AccountID: testhelpers.PublicKey1,
			}, nil)
		helper.mockedHorizonClient.
			On("AccountDetail", horizonclient.AccountRequest{
				AccountID: testhelpers.PublicKey2,
			}).
			Return(horizon.Account{
				AccountID: testhelpers.PublicKey2,
			}, nil)
		helper.mockedHorizonClient.
			On("AccountDetail", horizonclient.AccountRequest{
				AccountID: testhelpers.PublicKey3,
			}).
			Return(horizon.Account{}, errors.New("some error has occurred"))

		_, err := helper.repo.GetAccounts(
			testhelpers.PublicKey1,
			testhelpers.PublicKey2,
			testhelpers.PublicKey3,
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("fail to get account detail of %s", testhelpers.PublicKey3))

		helper.mockedHorizonClient.
			AssertNumberOfCalls(t, "AccountDetail", 3)
	})
	t.Run("error, fail to get account detail of public key 2 and 3, account 2 has higher priority in error message", func(t *testing.T) {
		helper := initTest()

		helper.mockedHorizonClient.
			On("AccountDetail", horizonclient.AccountRequest{
				AccountID: testhelpers.PublicKey1,
			}).
			Return(horizon.Account{
				AccountID: testhelpers.PublicKey1,
			}, nil)
		helper.mockedHorizonClient.
			On("AccountDetail", horizonclient.AccountRequest{
				AccountID: testhelpers.PublicKey2,
			}).
			Return(horizon.Account{}, errors.New("some error has occurred"))
		helper.mockedHorizonClient.
			On("AccountDetail", horizonclient.AccountRequest{
				AccountID: testhelpers.PublicKey3,
			}).
			Return(horizon.Account{}, errors.New("some error has occurred"))

		_, err := helper.repo.GetAccounts(
			testhelpers.PublicKey1,
			testhelpers.PublicKey2,
			testhelpers.PublicKey3,
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf(constants.ErrGetAccountDetail, testhelpers.PublicKey2))

		helper.mockedHorizonClient.
			AssertNumberOfCalls(t, "AccountDetail", 3)
	})
}
