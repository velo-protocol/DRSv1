package stellar_test

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/support/render/problem"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/testhelpers"
	"testing"
)

func TestRepo_GetAccount(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		helper := initTest()

		helper.mockedHorizonClient.
			On("AccountDetail", horizonclient.AccountRequest{
				AccountID: testhelpers.PublicKey1,
			}).
			Return(horizon.Account{
				AccountID: testhelpers.PublicKey1,
			}, nil)

		account, err := helper.repo.GetStellarAccount(testhelpers.PublicKey1)

		assert.NoError(t, err)
		assert.Equal(t, testhelpers.PublicKey1, account.AccountID)

		helper.mockedHorizonClient.
			AssertNumberOfCalls(t, "AccountDetail", 1)
	})
	t.Run("error, fail to get account detail, horizon error", func(t *testing.T) {
		helper := initTest()

		helper.mockedHorizonClient.
			On("AccountDetail", horizonclient.AccountRequest{
				AccountID: testhelpers.PublicKey1,
			}).
			Return(horizon.Account{}, &horizonclient.Error{
				Problem: problem.NotFound,
			})

		_, err := helper.repo.GetStellarAccount(testhelpers.PublicKey1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf(constants.ErrGetAccountDetail, testhelpers.PublicKey1))
		assert.Contains(t, err.Error(), problem.NotFound.Detail)

		helper.mockedHorizonClient.
			AssertNumberOfCalls(t, "AccountDetail", 1)
	})
	t.Run("error, fail to get account detail, non-horizon error", func(t *testing.T) {
		helper := initTest()

		helper.mockedHorizonClient.
			On("AccountDetail", horizonclient.AccountRequest{
				AccountID: testhelpers.PublicKey1,
			}).
			Return(horizon.Account{}, errors.New("some error has occurred"))

		_, err := helper.repo.GetStellarAccount(testhelpers.PublicKey1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf(constants.ErrGetAccountDetail, testhelpers.PublicKey1))
		assert.Contains(t, err.Error(), "some error has occurred")

		helper.mockedHorizonClient.
			AssertNumberOfCalls(t, "AccountDetail", 1)
	})
}
