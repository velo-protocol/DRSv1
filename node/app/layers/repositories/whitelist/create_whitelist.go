package whitelist

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"gitlab.com/velo-labs/cen/node/app/layers/repositories/whitelist/models"
)

func createWhitelist(dbTx *gorm.DB, whitelist *entities.WhiteList) (*entities.WhiteList, error) {
	createWhiteListModel := &models.CreateWhiteList{
		StellarPublicKey: whitelist.StellarPublicKey,
		RoleCode:         whitelist.RoleCode,
	}

	if err := dbTx.Save(createWhiteListModel).Error; err != nil {
		return nil, constants.ErrToSaveDatabase
	}

	return whitelist, nil
}

func (r *repo) CreateWhitelistTx(dbTx *gorm.DB, whitelist *entities.WhiteList) (*entities.WhiteList, error) {
	return createWhitelist(dbTx, whitelist)
}

func (r *repo) CreateWhitelist(whitelist *entities.WhiteList) (*entities.WhiteList, error) {
	return createWhitelist(r.Conn, whitelist)
}
