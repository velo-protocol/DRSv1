package entities

import "github.com/shopspring/decimal"

type RedeemCreditOutput struct {
	SignedStellarTxXdr      string
	AssetCodeToBeRedeemed   string
	AssetIssuerToBeRedeemed string
	AssetAmountToBeRedeemed decimal.Decimal
	CollateralCode          string
	CollateralIssuer        string
	CollateralAmount        decimal.Decimal
}
