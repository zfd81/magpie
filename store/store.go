package store

import "github.com/zfd81/magpie/store/hashcode"

var (
	poolSize = conf.StoragePoolSize
)

// Driver is the interface that must be implemented by a KV storage.
type Driver interface {
	// Open returns a new Storage.
	// The path is the string for storage specific format.
	Open(path string) error
}

type Storage interface {
	CreateTable(name string) error
	Put(table string, key, value []byte) error
	BatchPut(table string, kvs []KeyValue) error
	Get(table string, key []byte) ([]byte, error)
	GetWithPrefix(table string, prefix []byte) ([]*KeyValue, error)
	Delete(table string, key []byte) error
	DeleteWithPrefix(table string, prefix []byte) error
	Truncate(table string) error
	Iterator(table string, f func(k, v string) error) error
	IteratorWithPrefix(table string, prefix []byte, f func(k, v string) error) error
	Count(table string) int
}

type KeyValue struct {
	Key   []byte
	Value []byte
}

func New(path string) (Storage, error) {
	db := &boltdb{}
	err := db.Open(path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewStoragePool(path string) (*StoragePool, error) {
	pool := &StoragePool{
		pool: make([]Storage, poolSize),
	}
	err := pool.Open(path)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func PageIndex(key []byte) int {
	return hashcode.Hash(key) % poolSize
}
