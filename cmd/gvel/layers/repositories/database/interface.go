package database

type Repository interface {
	Init(path string) error
	Save(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	GetAll() ([][]byte, error)
}
