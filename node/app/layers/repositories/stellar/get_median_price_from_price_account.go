package stellar

import (
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/velo-protocol/DRSv1/node/app/constants"
	"github.com/velo-protocol/DRSv1/node/app/environments"
	"github.com/velo-protocol/DRSv1/node/app/layers/repositories/stellar/models"
	"sort"
	"time"
)

func (repo *repo) GetMedianPriceFromPriceAccount(priceAccountAddress string) (decimal.Decimal, error) {
	account, err := repo.HorizonClient.AccountDetail(horizonclient.AccountRequest{
		AccountID: priceAccountAddress,
	})
	if err != nil {
		return decimal.Zero, errors.Wrapf(err, constants.ErrGetAccountDetail, priceAccountAddress)
	}

	now := time.Now()

	// Parse data and filter out invalid data
	var prices []models.Price
	for key, encodedValue := range account.Data {
		price, err := models.NewPrice(key, encodedValue)
		if err != nil {
			continue
		}

		// if the difference between now and the price exceed ValidPriceBoundary (15 min default),
		// the price will be excluded from calculation
		if now.Sub(time.Unix(price.UnixTimestamp, 0)) > env.ValidPriceBoundary {
			continue
		}

		prices = append(prices, *price)
	}

	if len(prices) == 0 {
		return decimal.Zero, errors.Errorf(constants.ErrNoValidPrice, priceAccountAddress)
	}

	// Sort by value
	sort.SliceStable(prices, func(i, j int) bool {
		return prices[i].Value < prices[j].Value
	})

	var medianValue decimal.Decimal
	if len(prices)%2 == 0 { // even amount of samples
		// (leftMiddle + rightMiddle) / 2
		leftMiddle := decimal.New(prices[(len(prices)/2)-1].Value, -7)
		rightMiddle := decimal.New(prices[len(prices)/2].Value, -7)
		medianValue = leftMiddle.Add(rightMiddle).Div(decimal.New(2, 0))
	} else { // odd amount of samples
		medianValue = decimal.New(prices[len(prices)/2].Value, -7)
	}

	return medianValue.Truncate(7), nil
}
