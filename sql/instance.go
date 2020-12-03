package sql

import (
	"fmt"
	"os"
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

	dataDirectory := filepath.Join(conf.DataDirectory, "data")
	dir, err := os.Stat(dataDirectory)
	//判断数据目录是否存在
	if err != nil || !dir.IsDir() {
		err = os.MkdirAll(dataDirectory, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("mkdir failed![%v]\n", err)
		}
	}
	storage, err := store.NewStoragePool(filepath.Join(conf.DataDirectory, "data", db.FileName()))
	if err != nil {
		return nil, err
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
