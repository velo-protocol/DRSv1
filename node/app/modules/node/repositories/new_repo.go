package repositories

import (
	"github.com/syndtr/goleveldb/leveldb"
	"gitlab.com/velo-labs/cen/app/modules/node"
)

type repository struct {
	LevelConn *leveldb.DB
}

func NewNodeRepository(levelConn *leveldb.DB) node.Repository {
	return &repository{
		LevelConn: levelConn,
	}
}
