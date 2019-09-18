package models

import (
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
)

type CreateWhiteList struct {
	StellarPublicKey string
	RoleCode         string
}

func (CreateWhiteList) TableName() string {
	return constants.WhiteListTable
}

type GetWhiteListFilter struct {
	StellarPublicKey *string
	RoleCode         *string
}

func (GetWhiteListFilter) TableName() string {
	return constants.WhiteListTable
}

type GetWhiteList struct {
	ID               string `gorm:"primary_key"`
	StellarPublicKey string
	RoleCode         string
}

func (GetWhiteList) TableName() string {
	return constants.WhiteListTable
}

func (m GetWhiteList) ToEntity() (entity entities.WhiteList) {
	return entities.WhiteList{
		ID:               m.ID,
		StellarPublicKey: m.StellarPublicKey,
		RoleCode:         m.RoleCode,
	}
}
