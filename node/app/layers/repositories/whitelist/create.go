package whitelist

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/velo-labs/cen/libs/errors"
	"gitlab.com/velo-labs/cen/node/app/entities"
)

func create(dbTx *gorm.DB, whitelist *entities.Whitelist) (*entities.Whitelist, error) {
	if err := dbTx.Save(&whitelist).Error; err != nil {
		return nil, verrors.InternalError{Message: err.Error()}
	}

	return whitelist, nil
}

func (r *repo) CreateTx(dbTx *gorm.DB, whitelist *entities.Whitelist) (*entities.Whitelist, error) {
	return create(dbTx, whitelist)
}

func (r *repo) Create(whitelist *entities.Whitelist) (*entities.Whitelist, error) {
	return create(r.Conn, whitelist)
}


//func (r *repo) Create(ctx context.Context, whitelist entities.Whitelist) (*string, error) {
//	model := models.WhitelistModel{
//		StellarAddress: &whitelist.StellarAddress,
//		Role:           &whitelist.Role,
//	}
//
//	if err := r.Conn.Save(&model).Error; err != nil {
//		return nil, verrors.InternalError{Message: err.Error()}
//	}
//
//	return model.ID, nil
//}
//
//func (r *repo) CreateTx(ctx context.Context, tx, whitelist entities.Whitelist) (*string, error) {
//	return create(ctx, tx, entitiy)
//}