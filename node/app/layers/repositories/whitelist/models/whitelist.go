package models

import (
	"gitlab.com/velo-labs/cen/node/app/entities"
	"time"
)

type CreateWhiteList struct {

}

type GetWhiteList struct{

}

// WhitelistModel for talking with DB
type WhitelistModel struct {
	ID                *string `gorm:"primary_key"`
	StellarAddress    *string
	Role              *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
}

// WhitelistListModel array of whitelist model
type WhitelistListModel []WhitelistModel

// TableName to tell GORM know `WhitelistModel` must use `whitelists` table
func (WhitelistModel) TableName() string {
	return "whitelists"
}

// ToEntity convert model to entity
func (m WhitelistModel) ToEntity() (entity entities.Whitelist) {
	return entities.Whitelist{
		ID:                *m.ID,
		StellarAddress:    *m.StellarAddress,
		Role:              *m.Role,
	}
}

// ToEntities convert model to entities
func (m WhitelistListModel) ToEntities() []entities.Whitelist {
	var results []entities.Whitelist
	for _, v := range m {
		results = append(results, v.ToEntity())
	}
	return results
}