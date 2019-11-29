package entities

import "github.com/shopspring/decimal"

type MintCreditOutput struct {
	SignedStellarTxXdr         string
	AssetAmountToBeIssued      decimal.Decimal
	AssetCodeToBeIssued        string
	AssetIssuerToBeMinted      string
	AssetDistributorToBeMinted string
	CollateralAmount           decimal.Decimal
	CollateralAssetCode        string
}
