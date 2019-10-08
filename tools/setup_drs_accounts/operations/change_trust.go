package operations

import (
	"github.com/stellar/go/txnbuild"
)

func ChangeTrustOp(sourceAccount txnbuild.Account, assetCode, assetIssuer string) *txnbuild.ChangeTrust {

	return &txnbuild.ChangeTrust{
		Line: txnbuild.CreditAsset{
			Code:   assetCode,
			Issuer: assetIssuer,
		},
		SourceAccount: sourceAccount,
	}
}
