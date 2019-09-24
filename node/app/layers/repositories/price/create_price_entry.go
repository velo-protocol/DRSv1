package price

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"gitlab.com/velo-labs/cen/node/app/layers/repositories/price/models"
)

func createPriceEntry(dbTx *gorm.DB, priceEntry *entities.CreatePriceEntry) (*entities.CreatePriceEntry, error) {
	createPriceEntryModel := &models.CreatePriceEntry{
		FeederPublicKey:             priceEntry.FeederPublicKey,
		Asset:                       priceEntry.Asset,
		PriceInCurrencyPerAssetUnit: priceEntry.PriceInCurrencyPerAssetUnit,
		Currency:                    priceEntry.Currency,
	}

	if err := dbTx.Save(createPriceEntryModel).Error; err != nil {
		return nil, err
	}

	return priceEntry, nil
}

func (r *repo) CreatePriceEntryTx(dbTx *gorm.DB, priceEntry *entities.CreatePriceEntry) (*entities.CreatePriceEntry, error) {
	return createPriceEntry(dbTx, priceEntry)
}

func (r *repo) CreatePriceEntry(priceEntry *entities.CreatePriceEntry) (*entities.CreatePriceEntry, error) {
	return createPriceEntry(r.Conn, priceEntry)
}
