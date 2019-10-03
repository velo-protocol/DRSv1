package stellar

import (
	"github.com/stellar/go/protocols/horizon"
	"gitlab.com/velo-labs/cen/node/app/entities"
)

type Repo interface {
	GetAccount(stellarAddress string) (*horizon.Account, error)
	GetAccounts(stellarAddresses ...string) ([]horizon.Account, error)
	GetAccountData(stellarAddress string) (map[string]string, error)
	GetAccountDecodedData(stellarAddress string) (map[string]string, error)
	GetAccountDecodedDataByKey(stellarAddress string, key string) (string, error)
	GetDrsAccountData() (*entities.DrsAccountData, error)
	SubmitTransaction(txB64 string) (*horizon.TransactionSuccess, error)
}
