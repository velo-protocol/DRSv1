package stellar_test

import (
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stretchr/testify/assert"
	"github.com/velo-protocol/DRSv1/node/app/constants"
	"github.com/velo-protocol/DRSv1/node/app/testhelpers"
	"strconv"
	"testing"
	"time"
)

func TestRepo_GetMedianPriceFromPriceAccount(t *testing.T) {
	encodePrice := func(timestamp string, price string) string {
		return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s_%s", timestamp, price)))
	}

	t.Run("success", func(t *testing.T) {
		t.Run("price data count is even", func(t *testing.T) {
			helper := initTest()
			now := strconv.FormatInt(time.Now().Unix(), 10)

			helper.mockedHorizonClient.
				On("AccountDetail", horizonclient.AccountRequest{
					AccountID: testhelpers.PriceUsdVeloPublicKey,
				}).
				Return(horizon.Account{
					AccountID: testhelpers.PriceUsdVeloPublicKey,
					Data: map[string]string{
						testhelpers.PublicKey1: encodePrice(now, "14000000"), // price 1.4
						testhelpers.PublicKey2: encodePrice(now, "15000000"), // price 1.5
						testhelpers.PublicKey3: encodePrice(now, "11000000"), // price 1.1
						testhelpers.PublicKey4: encodePrice(now, "8000000"),  // price 0.8
					},
				}, nil)

			value, err := helper.repo.GetMedianPriceFromPriceAccount(testhelpers.PriceUsdVeloPublicKey)

			assert.NoError(t, err)
			assert.Equal(t, decimal.New(12500000, -7), value) // price 1.25

			helper.mockedHorizonClient.
				AssertNumberOfCalls(t, "AccountDetail", 1)
		})
		t.Run("price data count is odd", func(t *testing.T) {
			helper := initTest()
			now := strconv.FormatInt(time.Now().Unix(), 10)

			helper.mockedHorizonClient.
				On("AccountDetail", horizonclient.AccountRequest{
					AccountID: testhelpers.PriceUsdVeloPublicKey,
				}).
				Return(horizon.Account{
					AccountID: testhelpers.PriceUsdVeloPublicKey,
					Data: map[string]string{
						testhelpers.PublicKey1: encodePrice(now, "14000000"), // price 1.4
						testhelpers.PublicKey2: encodePrice(now, "15000000"), // price 1.5
						testhelpers.PublicKey3: encodePrice(now, "11000000"), // price 1.1
					},
				}, nil)

			value, err := helper.repo.GetMedianPriceFromPriceAccount(testhelpers.PriceUsdVeloPublicKey)

			assert.NoError(t, err)
			assert.Equal(t, decimal.New(14000000, -7), value) // price 1.4

			helper.mockedHorizonClient.
				AssertNumberOfCalls(t, "AccountDetail", 1)
		})
		t.Run("result must be truncated to -7 precision", func(t *testing.T) {
			helper := initTest()
			now := strconv.FormatInt(time.Now().Unix(), 10)

			helper.mockedHorizonClient.
				On("AccountDetail", horizonclient.AccountRequest{
					AccountID: testhelpers.PriceUsdVeloPublicKey,
				}).
				Return(horizon.Account{
					AccountID: testhelpers.PriceUsdVeloPublicKey,
					Data: map[string]string{
						testhelpers.PublicKey1: encodePrice(now, "14000009"), // price 1.4000009
						testhelpers.PublicKey2: encodePrice(now, "15000008"), // price 1.5000008
					},
				}, nil)

			value, err := helper.repo.GetMedianPriceFromPriceAccount(testhelpers.PriceUsdVeloPublicKey)

			assert.NoError(t, err)
			assert.Equal(t, decimal.New(14500008, -7), value) // price 1.45000085 ~ truncated 1.4500008

			helper.mockedHorizonClient.
				AssertNumberOfCalls(t, "AccountDetail", 1)
		})
	})

	t.Run("error, fail to get account detail", func(t *testing.T) {
		helper := initTest()

		helper.mockedHorizonClient.
			On("AccountDetail", horizonclient.AccountRequest{
				AccountID: testhelpers.PriceUsdVeloPublicKey,
			}).
			Return(horizon.Account{}, errors.New("some error has occurred"))

		_, err := helper.repo.GetMedianPriceFromPriceAccount(testhelpers.PriceUsdVeloPublicKey)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf(constants.ErrGetAccountDetail, testhelpers.PriceUsdVeloPublicKey))

		helper.mockedHorizonClient.
			AssertNumberOfCalls(t, "AccountDetail", 1)
	})

	t.Run("error, no valid price in price account", func(t *testing.T) {
		t.Run("invalid format of price data", func(t *testing.T) {
			helper := initTest()

			helper.mockedHorizonClient.
				On("AccountDetail", horizonclient.AccountRequest{
					AccountID: testhelpers.PriceUsdVeloPublicKey,
				}).
				Return(horizon.Account{
					AccountID: testhelpers.PriceUsdVeloPublicKey,
					Data: map[string]string{
						testhelpers.PublicKey1: "BAD_B64_VALUE",
					},
				}, nil)

			_, err := helper.repo.GetMedianPriceFromPriceAccount(testhelpers.PriceUsdVeloPublicKey)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), fmt.Sprintf("no valid price in price account %s", testhelpers.PriceUsdVeloPublicKey))

			helper.mockedHorizonClient.
				AssertNumberOfCalls(t, "AccountDetail", 1)
		})
		t.Run("timestamp of price data is outdated", func(t *testing.T) {
			helper := initTest()

			helper.mockedHorizonClient.
				On("AccountDetail", horizonclient.AccountRequest{
					AccountID: testhelpers.PriceUsdVeloPublicKey,
				}).
				Return(horizon.Account{
					AccountID: testhelpers.PriceUsdVeloPublicKey,
					Data: map[string]string{
						testhelpers.PublicKey1: encodePrice("0", "14000000"), // time 0, price 1.5
					},
				}, nil)

			_, err := helper.repo.GetMedianPriceFromPriceAccount(testhelpers.PriceUsdVeloPublicKey)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), fmt.Sprintf("no valid price in price account %s", testhelpers.PriceUsdVeloPublicKey))

			helper.mockedHorizonClient.
				AssertNumberOfCalls(t, "AccountDetail", 1)
		})
	})
}
