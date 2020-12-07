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

	storage, err := store.New(filepath.Join(conf.DataDirectory, fmt.Sprintf("%s.db", db.FileName())))
	if err != nil {
		return nil, err
	}
	db.storage = storage
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
