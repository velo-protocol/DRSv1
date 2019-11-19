package stellar_test

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/protocols/horizon/base"
	"github.com/stretchr/testify/assert"
	"github.com/velo-protocol/DRSv1/libs/xdr"
	"github.com/velo-protocol/DRSv1/node/app/constants"
	"github.com/velo-protocol/DRSv1/node/app/testhelpers"
	"testing"
)

func TestRepo_GetAccountBalances(t *testing.T) {
	var (
		assetType1   = "credit_alphanum4"
		assetType2   = "credit_alphanum12"
		balance      = "100.0000000"
		assetCode1   = string(vxdr.AssetVELO)
		assetCode2   = "kTONG"
		assetIssuer1 = testhelpers.PublicKey1
		assetIssuer2 = testhelpers.PublicKey2
	)

	t.Run("success", func(t *testing.T) {
		helper := initTest()

		helper.mockedHorizonClient.
			On("AccountDetail", horizonclient.AccountRequest{
				AccountID: testhelpers.PublicKey1,
			}).
			Return(horizon.Account{
				AccountID: testhelpers.PublicKey1,
				Balances: []horizon.Balance{
					{
						Balance: balance,
						Asset: base.Asset{
							Type:   assetType1,
							Code:   assetCode1,
							Issuer: assetIssuer1,
						},
					},
				},
			}, nil)

		balances, err := helper.repo.GetAccountBalances(testhelpers.PublicKey1)
		assert.NoError(t, err)
		assert.NotEmpty(t, balances)
		assert.Equal(t, balance, balances[0].Balance)
		assert.Equal(t, assetType1, balances[0].Asset.Type)
		assert.Equal(t, assetCode1, balances[0].Asset.Code)
		assert.Equal(t, assetIssuer1, balances[0].Asset.Issuer)
	})

	t.Run("success", func(t *testing.T) {
		helper := initTest()

		helper.mockedHorizonClient.
			On("AccountDetail", horizonclient.AccountRequest{
				AccountID: testhelpers.PublicKey1,
			}).
			Return(horizon.Account{
				AccountID: testhelpers.PublicKey1,
				Balances: []horizon.Balance{
					{
						Balance: balance,
						Asset: base.Asset{
							Type:   assetType1,
							Code:   assetCode1,
							Issuer: assetIssuer1,
						},
					},
					{
						Balance: balance,
						Asset: base.Asset{
							Type:   assetType2,
							Code:   assetCode2,
							Issuer: assetIssuer2,
						},
					},
				},
			}, nil)

		balances, err := helper.repo.GetAccountBalances(testhelpers.PublicKey1)
		assert.NoError(t, err)
		assert.NotEmpty(t, balances)
		assert.Len(t, balances, 2)
	})

	t.Run("error, fail to get account detail", func(t *testing.T) {
		helper := initTest()

		helper.mockedHorizonClient.
			On("AccountDetail", horizonclient.AccountRequest{
				AccountID: testhelpers.PublicKey1,
			}).
			Return(horizon.Account{}, errors.New("some error has occurred"))

		balances, err := helper.repo.GetAccountBalances(testhelpers.PublicKey1)

		assert.Error(t, err)
		assert.Nil(t, balances)
		assert.Contains(t, err.Error(), fmt.Sprintf(constants.ErrGetAccountDetail, testhelpers.PublicKey1))

		helper.mockedHorizonClient.
			AssertNumberOfCalls(t, "AccountDetail", 1)
	})
}
