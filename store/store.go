package store

// Driver is the interface that must be implemented by a KV storage.
type Driver interface {
	// Open returns a new Storage.
	// The path is the string for storage specific format.
	Open(path string) error
	CreateTable(name string) error
}

type Storage interface {
	Put(key, value []byte) error
	Get(key []byte) ([]byte, error)
	GetWithPrefix(prefix []byte) ([]*KeyValue, error)
	Delete(key []byte) error
	DeleteWithPrefix(prefix []byte) error
	Count() int
}

type KeyValue struct {
	Key   []byte
	Value []byte
}

func New(path string) (Storage, error) {
	db := &boltdb{}
	db.Open(path)
	db.CreateTable(tableName)
	return db, nil
}
