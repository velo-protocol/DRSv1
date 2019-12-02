package stellar

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/velo-protocol/DRSv1/node/app/constants"
)

func (repo *repo) GetAccountData(stellarAddress string) (map[string]string, error) {
	account, err := repo.HorizonClient.AccountDetail(horizonclient.AccountRequest{
		AccountID: stellarAddress,
	})
	if err != nil {
		return nil, errors.Wrapf(err, constants.ErrGetAccountDetail, stellarAddress)
	}

	return account.Data, nil
}
