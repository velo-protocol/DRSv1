package operations

import (
	"github.com/stellar/go/txnbuild"
)

func ManageDataOp(sourceAccount txnbuild.Account, key string, value string) *txnbuild.ManageData {
	return &txnbuild.ManageData{
		Name:          key,
		Value:         []byte(value),
		SourceAccount: sourceAccount,
	}
}
