package meta

import "fmt"

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

func (d *DatabaseInfo) CreateTable(name string) *TableInfo {
	tbl := &TableInfo{
		Name:     name,
		Text:     name,
		Columns:  make([]*ColumnInfo, 0, 10),
		Database: d,
	}
	return tbl
}
