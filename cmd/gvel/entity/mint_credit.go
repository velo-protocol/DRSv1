package entity

import (
	"github.com/stellar/go/protocols/horizon"
)

type MintCreditInput struct {
	AssetToBeMinted     string
	CollateralAssetCode string
	CollateralAmount    string
	Passphrase          string
}

type MintCreditOutput struct {
	AssetToBeMinted     string
	CollateralAssetCode string
	CollateralAmount    string
	SourceAddress       string
	TxResult            *horizon.TransactionSuccess
}
