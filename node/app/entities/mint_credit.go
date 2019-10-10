package entities

import "github.com/shopspring/decimal"

type MintCreditOutput struct {
	SignedStellarTxXdr string
	MintAmount         decimal.Decimal
	MintCurrency       string
	CollateralAmount   decimal.Decimal
	CollateralAsset    string
}
