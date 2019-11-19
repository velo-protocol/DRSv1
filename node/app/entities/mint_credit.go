package entities

import "github.com/shopspring/decimal"

type MintCreditOutput struct {
	SignedStellarTxXdr    string
	AssetAmountToBeIssued decimal.Decimal
	AssetCodeToBeIssued   string
	CollateralAmount      decimal.Decimal
	CollateralAssetCode   string
}
