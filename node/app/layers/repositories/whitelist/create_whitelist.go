package whitelist

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/velo-labs/cen/libs/errors"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"gitlab.com/velo-labs/cen/node/app/layers/repositories/whitelist/models"
)

func createWhitelist(dbTx *gorm.DB, whitelist *entities.Whitelist) (*entities.Whitelist, error) {
	createWhiteListModel := models.CreateWhiteList{
		StellarPublicAddress: &whitelist.StellarPublicAddress,
		RoleCode: &whitelist.RoleCode,
	}

	if err := dbTx.Save(&createWhiteListModel).Error; err != nil {
		return nil, verrors.InternalError{Message: err.Error()}
	}

	return whitelist, nil
}

func (r *repo) CreateWhitelistTx(dbTx *gorm.DB, whitelist *entities.Whitelist) (*entities.Whitelist, error) {
	return createWhitelist(dbTx, whitelist)
}

func (r *repo) CreateWhitelist(whitelist *entities.Whitelist) (*entities.Whitelist, error) {
	return createWhitelist(r.Conn, whitelist)
}
