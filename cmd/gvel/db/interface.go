package db

type DB interface {
	Save(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	GetAll() ([][]byte, error)
}
