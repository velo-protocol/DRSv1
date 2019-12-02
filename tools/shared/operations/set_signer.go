package operations

import (
	"github.com/stellar/go/txnbuild"
)

func SetSignerOp(sourceAccount txnbuild.Account, signerPublicAddress string, weight txnbuild.Threshold,
	masterWeight, lowThreshold, mediumThreshold, highThreshold *txnbuild.Threshold) *txnbuild.SetOptions {
	return &txnbuild.SetOptions{
		MasterWeight:    masterWeight,
		LowThreshold:    lowThreshold,
		MediumThreshold: mediumThreshold,
		HighThreshold:   highThreshold,
		Signer: &txnbuild.Signer{
			Address: signerPublicAddress,
			Weight:  weight,
		},
		SourceAccount: sourceAccount,
	}
}
