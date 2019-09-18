package whitelist

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/velo-labs/cen/libs/errors"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"gitlab.com/velo-labs/cen/node/app/layers/repositories/whitelist/models"
)

func (r *repo) FindOneRole(role string) (*entities.Role, error) {
	var resultModel models.RoleModel
	if err := r.Conn.Where("code = ?", role).First(&resultModel).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, verrors.InternalError{Message: err.Error()}
	}

	result := resultModel.ToEntity()
	return &result, nil
}