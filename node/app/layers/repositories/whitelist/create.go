package whitelist

import (
	"gitlab.com/velo-labs/cen/libs/errors"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"gitlab.com/velo-labs/cen/node/app/layers/repositories/whitelist/models"
)

func (r *repo) Create(whitelist entities.Whitelist) (*string, error) {
	model := models.WhitelistModel{
		StellarAddress: &whitelist.StellarAddress,
		Role:           &whitelist.Role,
	}

	if err := r.Conn.Save(&model).Error; err != nil {
		return nil, verrors.InternalError{Message: err.Error()}
	}

	return model.ID, nil
}