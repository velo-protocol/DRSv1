package whitelist

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/velo-labs/cen/node/app/entities"
)

type Repo interface {
	BeginTx() (*gorm.DB, error)
	CommitTx(dbtx *gorm.DB) (bool, error)

	CreateTx(dbTx *gorm.DB, whitelist *entities.Whitelist) (*entities.Whitelist, error)
	Create(whitelist *entities.Whitelist) (*entities.Whitelist, error)

	FindOne(filter entities.WhitelistFilter) (*entities.Whitelist, error)
}


//func UseCase()  {
//	uc := new UseCase
//
//	dbTx := uc.WhiteListRepo.BeginTx()
//
//	createEntity := uc.WhiteListRepo.Create(ctx, dbTx, entity)
//
//	dbTx
//
//}