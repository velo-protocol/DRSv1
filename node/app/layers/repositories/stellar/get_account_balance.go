package stellar

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
	"gitlab.com/velo-labs/cen/node/app/constants"
)

func (repo *repo) GetAccountBalances(stellarAddress string) ([]horizon.Balance, error) {
	account, err := repo.HorizonClient.AccountDetail(horizonclient.AccountRequest{
		AccountID: stellarAddress,
	})
	if err != nil {
		return nil, errors.Wrapf(err, constants.ErrGetAccountDetail, stellarAddress)
	}

	return account.Balances, nil
}
