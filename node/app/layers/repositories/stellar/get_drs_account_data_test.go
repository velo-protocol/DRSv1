package stellar_test

import (
	"encoding/base64"
	"errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"gitlab.com/velo-labs/cen/node/app/testhelpers"
	"testing"
)

func TestRepo_GetDrsAccountData(t *testing.T) {
	testhelpers.InitEnv()

	t.Run("success", func(t *testing.T) {
		helper := initTest()

		helper.mockedHorizonClient.
			On("AccountDetail", horizonclient.AccountRequest{
				AccountID: testhelpers.DrsPublicKey,
			}).
			Return(horizon.Account{
				AccountID: testhelpers.DrsPublicKey,
				Data: map[string]string{
					"TrustedPartnerList": base64.StdEncoding.EncodeToString([]byte("pk1")),
					"RegulatorList":      base64.StdEncoding.EncodeToString([]byte("pk2")),
					"PriceFeederList":    base64.StdEncoding.EncodeToString([]byte("pk3")),
					"Price[USD-VELO]":    base64.StdEncoding.EncodeToString([]byte("pk4")),
					"Price[THB-VELO]":    base64.StdEncoding.EncodeToString([]byte("pk5")),
					"Price[SGD-VELO]":    base64.StdEncoding.EncodeToString([]byte("pk6")),
				},
			}, nil)

		drsAccountData, err := helper.repo.GetDrsAccountData()

		assert.NoError(t, err)
		assert.Equal(t, entities.DrsAccountData{
			TrustedPartnerListAddress: "pk1",
			RegulatorListAddress:      "pk2",
			PriceFeederListAddress:    "pk3",
			PriceUsdVeloAddress:       "pk4",
			PriceThbVeloAddress:       "pk5",
			PriceSgdVeloAddress:       "pk6",
			Base64Decoded:             true,
		}, *drsAccountData)

		helper.mockedHorizonClient.
			AssertNumberOfCalls(t, "AccountDetail", 1)
	})
	t.Run("error, fail to get account detail", func(t *testing.T) {
		helper := initTest()

		helper.mockedHorizonClient.
			On("AccountDetail", horizonclient.AccountRequest{
				AccountID: testhelpers.DrsPublicKey,
			}).
			Return(horizon.Account{}, errors.New("some error has occurred"))

		_, err := helper.repo.GetDrsAccountData()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "fail to get account detail of drs account")

		helper.mockedHorizonClient.
			AssertNumberOfCalls(t, "AccountDetail", 1)
	})
}
