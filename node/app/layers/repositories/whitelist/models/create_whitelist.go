package models

import "time"


type WhitelistModel struct {
	ID                *string `gorm:"primary_key"`
	StellarAddress    *string
	Role              *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
}