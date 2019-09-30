package operations

import (
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
)

func CreateAccountOp(sourceAccount txnbuild.Account, startAmount string) (*txnbuild.CreateAccount, *keypair.Full) {
	newAccount, err := keypair.Random()
	if err != nil {
		panic(err)
	}

	return &txnbuild.CreateAccount{
		SourceAccount: sourceAccount,
		Destination:   newAccount.Address(),
		Amount:        startAmount,
	}, newAccount
}
