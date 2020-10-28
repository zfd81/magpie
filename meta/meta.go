package meta

import (
	"encoding/json"

	"github.com/zfd81/magpie/store"
)

const (
	PathSeparator  = "/"
	NameSeparator  = "."
	InstanceSuffix = ".ins"
	DatabaseSuffix = ".db"
	TableSuffix    = ".tbl"
	MetaPath       = "/meta"
)

type MetaInfo interface {
	GetMName() string
	GetPath() string
	Store() error
	Load() error
	Remove() error
}

func StoreMetadata(info MetaInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	return store.Put([]byte(info.GetPath()), data)
}

func LoadMetadata(info MetaInfo, args ...string) error {
	var path string
	if args != nil && len(args) > 0 {
		path = args[0]
	} else {
		path = info.GetPath()
	}
	data, err := store.Get([]byte(path))
	if err != nil {
		return err
	}
	return json.Unmarshal(data, info)
}

func RemoveMetadata(info MetaInfo) error {
	return store.Delete([]byte(info.GetPath()))
}
