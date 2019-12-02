package entities

import "github.com/shopspring/decimal"

type RebalanceOutput struct {
	Collaterals        []*Collateral
	SignedStellarTxXdr *string
}

type Collateral struct {
	AssetCode      string
	AssetIssuer    string
	RequiredAmount decimal.Decimal
	PoolAmount     decimal.Decimal
}
