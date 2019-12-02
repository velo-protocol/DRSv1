package stellar

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/velo-protocol/DRSv1/node/app/constants"
	"github.com/velo-protocol/DRSv1/node/app/entities"
	"github.com/velo-protocol/DRSv1/node/app/environments"
	"github.com/velo-protocol/DRSv1/node/app/layers/repositories/stellar/models"
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
