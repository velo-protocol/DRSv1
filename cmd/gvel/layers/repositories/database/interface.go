package database

type Repository interface {
	Save(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	GetAll() ([][]byte, error)
}
