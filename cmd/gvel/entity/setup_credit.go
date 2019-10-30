package entity

import "github.com/stellar/go/protocols/horizon"

type SetupCreditInput struct {
	Passphrase     string
	AssetCode      string
	PeggedValue    string
	PeggedCurrency string
}

type SetupCreditOutput struct {
	AssetCode      string
	PeggedValue    string
	PeggedCurrency string
	SourceAddress  string
	TxResult       *horizon.TransactionSuccess
}
