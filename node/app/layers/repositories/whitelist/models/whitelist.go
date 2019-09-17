package models

import (
	"gitlab.com/velo-labs/cen/node/app/entities"
	"time"
)

type CreateWhiteList struct {

}

type GetWhiteList struct{

}

type WhitelistModel struct {
	ID                *string `gorm:"primary_key"`
	StellarAddress    *string
	Role              *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
}

type WhitelistListModel []WhitelistModel

func (WhitelistModel) TableName() string {
	return "whitelists"
}

func (m WhitelistModel) ToEntity() (entity entities.Whitelist) {
	return entities.Whitelist{
		ID:                *m.ID,
		StellarAddress:    *m.StellarAddress,
		Role:              *m.Role,
	}
}

func (m WhitelistListModel) ToEntities() []entities.Whitelist {
	var results []entities.Whitelist
	for _, v := range m {
		results = append(results, v.ToEntity())
	}
	return results
}