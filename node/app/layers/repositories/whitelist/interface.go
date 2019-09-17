package whitelist

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/velo-labs/cen/node/app/entities"
)

type Repo interface {
	BeginTx() (*gorm.DB, error)
	CommitTx(dbtx *gorm.DB) (bool, error)

	CreateWhitelistTx(dbTx *gorm.DB, whitelist *entities.Whitelist) (*entities.Whitelist, error)
	CreateWhitelist(whitelist *entities.Whitelist) (*entities.Whitelist, error)

	FindOneWhitelist(filter entities.WhitelistFilter) (*entities.Whitelist, error)
	FindOneRole(role string) (*entities.Role, error)
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