package entity

import "github.com/stellar/go/protocols/horizon"

type RedeemCreditInput struct {
	AssetCodeToBeRedeemed   string
	AssetIssuerToBeRedeemed string
	AmountToBeRedeemed      string
	Passphrase              string
}

type RedeemCreditOutput struct {
	AssetCodeToBeRedeemed   string
	AssetIssuerToBeRedeemed string
	AmountToBeRedeemed      string
	CollateralCode          string
	CollateralIssuer        string
	CollateralAmount        string
	TxResult                *horizon.TransactionSuccess
}
