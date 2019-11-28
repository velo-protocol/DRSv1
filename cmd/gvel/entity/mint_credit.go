package entity

import (
	"github.com/stellar/go/protocols/horizon"
)

type MintCreditInput struct {
	AssetCodeToBeMinted string
	CollateralAssetCode string
	CollateralAmount    string
	Passphrase          string
}

type MintCreditOutput struct {
	AssetCodeToBeMinted        string
	CollateralAssetCode        string
	CollateralAmount           string
	AssetIssuerToBeIssued      string
	AssetDistributorToBeIssued string
	AssetAmountToBeIssued      string
	SourceAddress              string
	TxResult                   *horizon.TransactionSuccess
}
