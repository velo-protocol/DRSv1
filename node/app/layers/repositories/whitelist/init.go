package whitelist

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/velo-labs/cen/node/app/entities"
)

type repo struct {
	Conn *gorm.DB
}

// InitRepo CRUD
func InitRepo(Conn *gorm.DB) Repo {
	return &repo{Conn: Conn}
}

// Repo CRUD interface
type Repo interface {
	Create(whitelist entities.Whitelist) (*string, error)
}