package entities

import (
	"github.com/shopspring/decimal"
)

type GetCollateralHealthCheckOutput struct {
	AssetCode      string
	AssetIssuer    string
	RequiredAmount decimal.Decimal
	PoolAmount     decimal.Decimal
}
