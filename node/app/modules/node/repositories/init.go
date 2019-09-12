package repositories

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/velo-labs/cen/node/app/modules/node"
)

type repository struct {
	DB *gorm.DB
}

func NewNodeRepository(dbConn *gorm.DB ) node.Repository {
	return &repository{
		DB: dbConn,
	}
}
