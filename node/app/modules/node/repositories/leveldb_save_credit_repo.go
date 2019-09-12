package repositories

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/velo-labs/cen/node/app/entities"
	env "gitlab.com/velo-labs/cen/node/app/environments"
	"gitlab.com/velo-labs/cen/node/app/modules/node/repositories/models"
)

func (repo *repository) SaveCredit(credit entities.Credit) error {
	creditModel, err := new(models.Credit).Parse(&credit)
	if err != nil {
		return errors.Wrap(err, "failed to parse credit entity to credit model")
	}

	creditJson, err := json.Marshal(creditModel)
	if err != nil {
		return err
	}

	err = repo.LevelConn.Put([]byte(env.CreditPrefix+credit.CreditOwnerAddress+credit.AssetName), creditJson, nil)

	return nil
}
