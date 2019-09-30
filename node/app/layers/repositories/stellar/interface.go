package stellar

import "github.com/stellar/go/protocols/horizon"

type Repo interface {
	GetAccount(stellarAddress string) (*horizon.Account, error)
	GetAccounts(stellarAddresses ...string) ([]horizon.Account, error)
	GetAccountData(stellarAddress string) (map[string]string, error)
	GetAccountDecodedData(stellarAddress string) (map[string]string, error)
	SubmitTransaction(txB64 string) (*horizon.TransactionSuccess, error)
}
