package database

import (
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
)

type levelDbDatabase struct {
	conn *leveldb.DB
}

func NewLevelDbDatabase(path string) (*levelDbDatabase, error) {
	conn, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open leveldb")
	}

	return &levelDbDatabase{
		conn: conn,
	}, nil
}

func (database *levelDbDatabase) Save(key []byte, value []byte) error {
	return database.conn.Put(key, value, nil)
}

func (database *levelDbDatabase) Get(key []byte) ([]byte, error) {
	return database.conn.Get(key, nil)
}

func (database *levelDbDatabase) GetAll() ([][]byte, error) {
	var all [][]byte

	iter := database.conn.NewIterator(nil, nil)
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
