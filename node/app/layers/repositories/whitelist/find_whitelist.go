package whitelist

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/velo-labs/cen/libs/errors"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"gitlab.com/velo-labs/cen/node/app/layers/repositories/whitelist/models"
)

func (r repo) FindOneWhitelist(filter entities.WhitelistFilter) (*entities.Whitelist, error) {
	var resultModel models.WhitelistModel
	if err := r.Conn.Where(makeFilterAttr(filter)).First(&resultModel).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, verrors.InternalError{Message: err.Error()}
	}

	result := resultModel.ToEntity()
	return &result, nil
}

func makeFilterAttr(filter entities.WhitelistFilter) (whitelistFilterAttr models.WhitelistModel) {
	whitelistFilterAttr = models.WhitelistModel{}

	if filter.StellarPublicAddress != nil {
		whitelistFilterAttr.StellarPublicAddress = filter.StellarPublicAddress
	}

	if filter.RoleCode != nil {
		whitelistFilterAttr.RoleCode = filter.RoleCode
	}

	return whitelistFilterAttr
}
