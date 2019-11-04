package entity

import "github.com/stellar/go/protocols/horizon"

type RedeemCreditInput struct {
	AssetCode   string
	AssetIssuer string
	Amount      string
	Passphrase  string
}

type RedeemCreditOutput struct {
	AssetCode   string
	AssetIssuer string
	Amount      string
	TxResult    *horizon.TransactionSuccess
}
