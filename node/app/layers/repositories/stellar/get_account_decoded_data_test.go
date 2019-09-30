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

func TestRepo_GetAccountDecodedData(t *testing.T) {
	testhelpers.InitEnv()

	t.Run("success", func(t *testing.T) {
		helper := initTest()

		helper.mockedHorizonClient.
			On("AccountDetail", horizonclient.AccountRequest{
				AccountID: testhelpers.PublicKey1,
			}).
			Return(horizon.Account{
				AccountID: testhelpers.PublicKey1,
				Data:      map[string]string{"key1": "dmFsdWUx"},
			}, nil)

		decodedData, err := helper.repo.GetAccountDecodedData(testhelpers.PublicKey1)

		assert.NoError(t, err)
		assert.Equal(t, map[string]string{"key1": "value1"}, decodedData)

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

		_, err := helper.repo.GetAccountData(testhelpers.PublicKey1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("fail to get account detail of %s", testhelpers.PublicKey1))

		helper.mockedHorizonClient.
			AssertNumberOfCalls(t, "AccountDetail", 1)
	})

	t.Run("error, fail to decode data", func(t *testing.T) {
		helper := initTest()

		helper.mockedHorizonClient.
			On("AccountDetail", horizonclient.AccountRequest{
				AccountID: testhelpers.PublicKey1,
			}).
			Return(horizon.Account{
				AccountID: testhelpers.PublicKey1,
				Data:      map[string]string{"key1": "BAD_B64_VALUE"},
			}, nil)

		_, err := helper.repo.GetAccountDecodedData(testhelpers.PublicKey1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), `fail to decode data "key1"`)

		helper.mockedHorizonClient.
			AssertNumberOfCalls(t, "AccountDetail", 1)
	})
}
