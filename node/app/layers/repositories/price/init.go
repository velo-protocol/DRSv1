package price

import (
	"github.com/jinzhu/gorm"
)

type repo struct {
	Conn *gorm.DB
}

// InitRepo CRUD
func InitRepo(conn *gorm.DB) Repo {
	return &repo{Conn: conn}
}
