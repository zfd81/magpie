package store

import (
	"log"
)

const storagePath = "magpie.db"

var magpieDB Storage

func init() {
	storage, err := New(storagePath)
	if err != nil {
		log.Panicln(err)
	}
	magpieDB = storage
}

func Put(key, value []byte) error {
	return magpieDB.Put(key, value)
}

func Get(key []byte) ([]byte, error) {
	return magpieDB.Get(key)
}

func GetWithPrefix(prefix []byte) ([]*KeyValue, error) {
	return magpieDB.GetWithPrefix(prefix)
}

func Delete(key []byte) error {
	return magpieDB.Delete(key)
}

func DeleteWithPrefix(prefix []byte) error {
	return magpieDB.DeleteWithPrefix(prefix)
}

func Count() int {
	return magpieDB.Count()
}
