package stellar

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
	"github.com/velo-protocol/DRSv1/node/app/constants"
)

func (repo *repo) GetAccount(stellarAddress string) (*horizon.Account, error) {
	account, err := repo.HorizonClient.AccountDetail(horizonclient.AccountRequest{
		AccountID: stellarAddress,
	})
	if err != nil {
		return nil, errors.Wrapf(err, constants.ErrGetAccountDetail, stellarAddress)
	}

	return &account, nil
}
