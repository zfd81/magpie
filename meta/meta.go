package meta

import (
	"encoding/json"
	"log"

	"github.com/zfd81/magpie/store"
)

const (
	PathSeparator  = "/"
	NameSeparator  = "."
	InstanceSuffix = ".ins"
	DatabaseSuffix = ".db"
	TableSuffix    = ".tbl"
	MetaPath       = "/meta"
	StoragePath    = "meta.db"
)

type MetaInfo interface {
	GetMName() string
	GetPath() string
	Store() error
	Load() error
}

var (
	metaDB store.Storage
)

func init() {
	storage, err := store.New(StoragePath)
	if err != nil {
		log.Panicln(err)
	}
	metaDB = storage
}

func StoreMetadata(info MetaInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	return metaDB.Put([]byte(info.GetPath()), data)
}

func LoadMetadata(info MetaInfo) error {
	data, err := ReadMetadata(info.GetPath())
	if err != nil {
		return err
	}
	return json.Unmarshal(data, info)
}

func ReadMetadata(path string) ([]byte, error) {
	return metaDB.Get([]byte(path))
}
