package stellar_test

import (
	"errors"
	"fmt"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stretchr/testify/assert"
	"github.com/velo-protocol/DRSv1/node/app/constants"
	"github.com/velo-protocol/DRSv1/node/app/testhelpers"
	"testing"
)

func TestRepo_GetAccountData(t *testing.T) {
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

		data, err := helper.repo.GetAccountData(testhelpers.PublicKey1)

		assert.NoError(t, err)
		assert.Equal(t, map[string]string{"key1": "dmFsdWUx"}, data)

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
		assert.Contains(t, err.Error(), fmt.Sprintf(constants.ErrGetAccountDetail, testhelpers.PublicKey1))

		helper.mockedHorizonClient.
			AssertNumberOfCalls(t, "AccountDetail", 1)
	})
}
