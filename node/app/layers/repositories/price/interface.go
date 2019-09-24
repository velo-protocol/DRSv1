package price

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/velo-labs/cen/node/app/entities"
)

type Repo interface {
	CreatePriceEntryTx(dbTx *gorm.DB, whitelist *entities.CreatePriceEntry) (*entities.CreatePriceEntry, error)
	CreatePriceEntry(whitelist *entities.CreatePriceEntry) (*entities.CreatePriceEntry, error)
}
