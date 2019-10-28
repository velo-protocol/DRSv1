package database

import (
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDbDatabase struct {
	conn *leveldb.DB
}

func NewLevelDb(path string) (*LevelDbDatabase, error) {
	conn, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open leveldb")
	}

	return &LevelDbDatabase{
		conn: conn,
	}, nil
}

func (database *LevelDbDatabase) Init(path string) error {
	_, err := NewLevelDb(path)
	return err
}

func (database *LevelDbDatabase) Save(key []byte, value []byte) error {
	return database.conn.Put(key, value, nil)
}

func (database *LevelDbDatabase) Get(key []byte) ([]byte, error) {
	return database.conn.Get(key, nil)
}

func (database *LevelDbDatabase) GetAll() ([][]byte, error) {
	var all [][]byte

	iter := database.conn.NewIterator(nil, nil)
	for iter.Next() {
		val := make([]byte, len(iter.Value()))
		copy(val[:], iter.Value()[:])
		all = append(all, val)
	}
	iter.Release()

	err := iter.Error()
	if err != nil {
		return nil, err
	}

	return all, nil
}
