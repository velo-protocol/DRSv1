package entities

import (
	"github.com/shopspring/decimal"
	"github.com/stellar/go/protocols/horizon"
)

type GetIssuerAccountInput struct {
	IssuerAddress string
}

type GetIssuerAccountOutput struct {
	Account        *horizon.Account
	PeggedValue    decimal.Decimal
	PeggedCurrency string
	AssetCode      string
}
