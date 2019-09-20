package whitelist

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"gitlab.com/velo-labs/cen/node/app/layers/repositories/whitelist/models"
)

func (r *repo) FindOneRole(role string) (*entities.Role, error) {
	var resultModel models.GetRole
	if err := r.Conn.Where("code = ?", role).First(&resultModel).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, errors.New(constants.ErrToGetDataFromDatabase)
	}

	result := resultModel.ToEntity()
	return &result, nil
}
