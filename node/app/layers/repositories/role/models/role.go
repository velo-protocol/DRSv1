package models

import (
	"gitlab.com/velo-labs/cen/node/app/entities"
	"time"
)

// RoleModel for talking with DB
type RoleModel struct {
	ID        *string `gorm:"primary_key"`
	Name      *string
	Code      *string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// TableName to tell GORM know `RoleModel` must use `ambassadors` table
func (RoleModel) TableName() string {
	return "roles"
}

// ToEntity convert model to entity
func (m RoleModel) ToEntity() (entity entities.Role) {
	return entities.Role{
		ID:   *m.ID,
		Name: *m.Name,
		Code: *m.Code,
	}
}
