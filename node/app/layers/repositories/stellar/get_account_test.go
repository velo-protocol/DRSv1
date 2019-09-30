package stellar_test

import (
	"errors"
	"fmt"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/node/app/testhelpers"
	"testing"
)

func TestRepo_GetAccount(t *testing.T) {
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

		account, err := helper.repo.GetAccount(testhelpers.PublicKey1)

		assert.NoError(t, err)
		assert.Equal(t, testhelpers.PublicKey1, account.AccountID)

		helper.mockedHorizonClient.
			AssertNumberOfCalls(t, "AccountDetail", 1)
	})
	t.Run("error, fail to get account detail", func(t *testing.T) {
		helper := initTest()

		helper.mockedHorizonClient.
			On("AccountDetail", horizonclient.AccountRequest{
				AccountID: testhelpers.PublicKey1,
			}).
			Return(horizon.Account{}, errors.New("some error has occurred"))

		_, err := helper.repo.GetAccount(testhelpers.PublicKey1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("fail to get account detail of %s", testhelpers.PublicKey1))

		helper.mockedHorizonClient.
			AssertNumberOfCalls(t, "AccountDetail", 1)
	})
}
