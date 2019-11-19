package entities

import "github.com/shopspring/decimal"

type UpdatePriceOutput struct {
	SignedStellarTxXdr          string
	Asset                       string
	Currency                    string
	PriceInCurrencyPerAssetUnit decimal.Decimal
}
