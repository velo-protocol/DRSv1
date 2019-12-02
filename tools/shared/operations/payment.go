package operations

import (
	"github.com/stellar/go/txnbuild"
)

func PaymentOp(sourceAccount txnbuild.Account, assetCode string, assetIssuer string, destinationAddress string, amount string) *txnbuild.Payment {
	return &txnbuild.Payment{
		Destination: destinationAddress,
		Amount:      amount,
		Asset: &txnbuild.CreditAsset{
			Code:   assetCode,
			Issuer: assetIssuer,
		},
		SourceAccount: sourceAccount,
	}
}
