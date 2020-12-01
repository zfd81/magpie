package sql

import (
	"fmt"
	"path/filepath"

	"github.com/zfd81/magpie/config"
	"github.com/zfd81/magpie/meta"
	"github.com/zfd81/magpie/store"
)

var (
	conf = config.GetConfig()
)

type Instance struct {
	meta.InstanceInfo
	Databases map[string]*Database
}

func (i *Instance) CreateDatabase(info meta.DatabaseInfo) (*Database, error) {
	db := &Database{
		DatabaseInfo: meta.DatabaseInfo{
			Name:     info.Name,
			Text:     info.Text,
			Instance: i.InstanceInfo,
		},
		Tables: map[string]*Table{},
	}
	i.Databases[db.Name] = db
	var storage, err = store.NewStoragePool(filepath.Join(conf.DataDirectory, db.FileName()))
	if err != nil {
		storage, err = store.NewStoragePool(filepath.Join(conf.DataDirectory, fmt.Sprintf("%s%d", db.FileName(), conf.Port)))
		if err != nil {
			return nil, err
		}
	}
	db.storagePool = storage
	return db, nil
}

func (i *Instance) GetDatabase(name string) *Database {
	return i.Databases[name]
}

func NewInstance(name, text string) *Instance {
	return &Instance{
		InstanceInfo: meta.InstanceInfo{
			Name: name,
			Text: text,
		},
		Databases: map[string]*Database{},
	}
}
