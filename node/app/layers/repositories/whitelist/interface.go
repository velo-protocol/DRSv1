package whitelist

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/velo-labs/cen/node/app/entities"
)

type Repo interface {
	BeginTx() *gorm.DB
	CommitTx(dbTx *gorm.DB) error

	CreateWhitelistTx(dbTx *gorm.DB, whitelist *entities.WhiteList) (*entities.WhiteList, error)
	CreateWhitelist(whitelist *entities.WhiteList) (*entities.WhiteList, error)
	FindOneWhitelist(filter entities.WhiteListFilter) (*entities.WhiteList, error)

	FindOneRole(role string) (*entities.Role, error)
}
