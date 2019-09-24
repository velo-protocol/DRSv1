package entities

import "github.com/shopspring/decimal"

type CreatePriceEntry struct {
	FeederPublicKey             string
	Asset                       string
	PriceInCurrencyPerAssetUnit decimal.Decimal
	Currency                    string
}
