package leveldb

import (
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"gitlab.com/velo-labs/cen/cmd/gvel/db"
)

type Conn struct {
	db *leveldb.DB
}

func NewLevelDB(path string) (db.DB, error) {
	conn := &Conn{}
	var err error

	conn.db, err = leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open leveldb")
	}

	return conn, nil
}

func (c *Conn) Save(key []byte, value []byte) error {
	return c.db.Put(key, value, nil)
}

func (c *Conn) Get(key []byte) ([]byte, error) {
	return c.db.Get(key, nil)
}

func (c *Conn) GetAll() ([][]byte, error) {
	var all [][]byte

	iter := c.db.NewIterator(nil, nil)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	iter.Release()

	err := iter.Error()
	if err != nil {
		return nil, err
	}

	return all, nil
}
