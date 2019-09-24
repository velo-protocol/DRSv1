package models

import (
	"github.com/shopspring/decimal"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"time"
)

type CreatePriceEntry struct {
	FeederPublicKey             string
	Asset                       string
	PriceInCurrencyPerAssetUnit decimal.Decimal
	Currency                    string
	CreatedAt                   time.Time
}

func (CreatePriceEntry) TableName() string {
	return constants.PriceEntryTable
}
