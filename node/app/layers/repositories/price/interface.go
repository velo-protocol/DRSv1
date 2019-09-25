package price

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/velo-labs/cen/node/app/entities"
)

type Repo interface {
	CreatePriceEntryTx(dbTx *gorm.DB, priceEntry *entities.CreatePriceEntry) (*entities.CreatePriceEntry, error)
	CreatePriceEntry(priceEntry *entities.CreatePriceEntry) (*entities.CreatePriceEntry, error)
}
