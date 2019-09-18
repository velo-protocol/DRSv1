package models

import (
	"gitlab.com/velo-labs/cen/node/app/entities"
	"time"
)

type CreateWhiteList struct {
	ID                   *string `gorm:"primary_key"`
	StellarPublicAddress *string
	RoleCode             *string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            *time.Time
}

func (CreateWhiteList) TableName() string {
	return "whitelists"
}

type GetWhiteList struct{
	ID                   *string `gorm:"primary_key"`
	StellarPublicAddress *string
	RoleCode             *string
}

type WhitelistListModel []GetWhiteList

func (GetWhiteList) TableName() string {
	return "whitelists"
}

func (m GetWhiteList) ToEntity() (entity entities.Whitelist) {
	return entities.Whitelist{
		ID:                   *m.ID,
		StellarPublicAddress: *m.StellarPublicAddress,
		RoleCode:             *m.RoleCode,
	}
}

func (m WhitelistListModel) ToEntities() []entities.Whitelist {
	var results []entities.Whitelist
	for _, v := range m {
		results = append(results, v.ToEntity())
	}
	return results
}