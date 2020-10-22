package server

import "github.com/zfd81/magpie/meta"

type Instance struct {
	meta.InstanceInfo
	Databases map[string]*Database
}

func (i *Instance) CreateDatabase(info meta.DatabaseInfo) *Database {
	db := &Database{
		DatabaseInfo: meta.DatabaseInfo{
			Name:     info.Name,
			Text:     info.Text,
			Instance: i.InstanceInfo,
		},
		Tables: map[string]*Table{},
	}
	i.Databases[db.Name] = db
	return db
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
