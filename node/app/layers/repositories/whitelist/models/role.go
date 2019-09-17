package models

import (
	"gitlab.com/velo-labs/cen/node/app/entities"
	"time"
)

type RoleModel struct {
	ID        *string `gorm:"primary_key"`
	Name      *string
	Code      *string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (RoleModel) TableName() string {
	return "roles"
}

func (m RoleModel) ToEntity() (entity entities.Role) {
	return entities.Role{
		ID:   *m.ID,
		Name: *m.Name,
		Code: *m.Code,
	}
}
