package stellar

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"gitlab.com/velo-labs/cen/node/app/entities"
	env "gitlab.com/velo-labs/cen/node/app/environments"
)

func (repo *repo) GetDrsAccountData() (*entities.DrsAccountData, error) {
	account, err := repo.HorizonClient.AccountDetail(horizonclient.AccountRequest{
		AccountID: env.DrsPublicKey,
	})
	if err != nil {
		return nil, errors.Wrap(err, "fail to get account detail of drs account")
	}

	drsAccountData := new(entities.DrsAccountData)
	err = mapstructure.Decode(account.Data, drsAccountData)
	if err != nil {
		return nil, errors.Wrap(err, "fail to map drs account data to entity")
	}

	return drsAccountData, nil
}
