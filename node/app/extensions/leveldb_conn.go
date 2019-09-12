package extensions

import (
	"github.com/syndtr/goleveldb/leveldb"
	env "gitlab.com/velo-labs/cen/node/app/environments"
)

func ConnLevelDB() *leveldb.DB {
	db, err := leveldb.OpenFile(env.LevelDBPath, nil)
	if err != nil {
		panic("failed to connect to leveldb")
	}

	return db
}
