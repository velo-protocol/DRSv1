package stellar

import "github.com/stellar/go/protocols/horizon"

type Repository interface {
	SubmitTransaction(txB64 string) (*horizon.TransactionSuccess, error)
	LoadAccount(stellarAddress string) (*horizon.Account, error)
}
