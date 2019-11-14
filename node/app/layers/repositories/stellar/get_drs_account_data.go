package stellar

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"gitlab.com/velo-labs/cen/node/app/environments"
	"gitlab.com/velo-labs/cen/node/app/layers/repositories/stellar/models"
)

func (repo *repo) GetDrsAccountData() (*entities.DrsAccountData, error) {
	account, err := repo.HorizonClient.AccountDetail(horizonclient.AccountRequest{
		AccountID: env.DrsPublicKey,
	})
	if err != nil {
		return nil, errors.Wrap(err, constants.ErrGetDrsAccountDetail)
	}

	drsAccountDataModel := models.DrsAccountData(account.Data)
	return drsAccountDataModel.Entity()
}
