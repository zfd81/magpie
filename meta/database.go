package meta

import (
	"encoding/json"
	"fmt"

	"github.com/zfd81/magpie/store"
)

type DatabaseInfo struct {
	Name     string       `json:"name"`
	Text     string       `json:"text"`
	Comment  string       `json:"comment,omitempty"`
	Charset  string       `json:"charset"`
	Instance InstanceInfo `json:"-"`
}

func (d *DatabaseInfo) GetMName() string {
	return fmt.Sprintf("%s%s", d.Name, DatabaseSuffix)
}

func (d *DatabaseInfo) GetPath() string {
	return fmt.Sprintf("%s%s%s", d.Instance.GetPath(), PathSeparator, d.GetMName())
}

func (i *DatabaseInfo) Store() error {
	return StoreMetadata(i)
}

func (i *DatabaseInfo) Load() error {
	return LoadMetadata(i)
}

func (i *DatabaseInfo) Remove() error {
	return RemoveMetadata(i)
}

func (d *DatabaseInfo) CreateTable(name string) *TableInfo {
	tbl := &TableInfo{
		Name:     name,
		Text:     name,
		Columns:  make([]*ColumnInfo, 0, 10),
		Database: *d,
	}
	return tbl
}

func (d *DatabaseInfo) ListTables() []*TableInfo {
	var tbls []*TableInfo
	kvs, err := store.GetWithPrefix([]byte(d.GetPath()))
	if err != nil {
		return tbls
	}
	for _, kv := range kvs {
		tbl := &TableInfo{}
		err = json.Unmarshal(kv.Value, tbl)
		if err != nil {
			return tbls
		}
		tbls = append(tbls, tbl)
	}
	return tbls
}
