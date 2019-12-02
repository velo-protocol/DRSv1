package entities

import "github.com/shopspring/decimal"

type UpdatePriceOutput struct {
	SignedStellarTxXdr          string
	CollateralCode              string
	Currency                    string
	PriceInCurrencyPerAssetUnit decimal.Decimal
}
