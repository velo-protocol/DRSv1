package stellar

import "github.com/stellar/go/protocols/horizon"

type Repo interface {
	LoadAccount(stellarAddress string) (*horizon.Account, error)
	SubmitTransaction(txB64 string) (*horizon.TransactionSuccess, error)
}
