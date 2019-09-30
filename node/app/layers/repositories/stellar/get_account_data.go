package stellar

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
)

func (repo *repo) GetAccountData(stellarAddress string) (map[string]string, error) {
	account, err := repo.HorizonClient.AccountDetail(horizonclient.AccountRequest{
		AccountID: stellarAddress,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "fail to get account detail of %s", stellarAddress)
	}

	return account.Data, nil
}
