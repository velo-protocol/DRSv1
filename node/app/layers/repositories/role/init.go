package role

import (
	"github.com/jinzhu/gorm"
)

type repo struct {
	Conn *gorm.DB
}

// InitRepo CRUD
func InitRepo(Conn *gorm.DB) Repo {
	return &repo{Conn: Conn}
}
