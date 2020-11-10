package meta

import (
	"fmt"
)

type TableInfo struct {
	Id      string        `json:"id,omitempty"`
	Name    string        `json:"name"`
	Text    string        `json:"text,omitempty"`
	Comment string        `json:"comment,omitempty"`
	Charset string        `json:"charset,omitempty"`
	Columns []*ColumnInfo `json:"cols"`
	Keys    []string      `json:"keys"`
	Indexes []string      `json:"indexes,omitempty"`
	//DerivedCols []*ColumnInfo `json:"dcols"`
	Database DatabaseInfo `json:"-"`
}

func (t *TableInfo) GetMName() string {
	return fmt.Sprintf("%s%s", t.Name, TableSuffix)
}

func (t *TableInfo) GetPath() string {
	return fmt.Sprintf("%s%s%s", t.Database.GetPath(), PathSeparator, t.GetMName())
}

func (i *TableInfo) Store() error {
	return StoreMetadata(i)
}

func (i *TableInfo) Load() error {
	return LoadMetadata(i)
}

func (i *TableInfo) Remove() error {
	return RemoveMetadata(i)
}

func (t *TableInfo) CreateColumn(name string, dataType string) *ColumnInfo {
	col := &ColumnInfo{
		Name:       name,
		Text:       name,
		DataType:   dataType,
		Expression: name,
	}
	col.Index = len(t.Columns)
	t.Columns = append(t.Columns, col)
	return col
}

func (t *TableInfo) RemoveColumn(name string) *TableInfo {
	for i, v := range t.Columns { //删除基础列
		if v.Name == name {
			t.Columns = append(t.Columns[:i], t.Columns[i+1:]...)
			break
		}
	}
	for i, v := range t.Columns {
		v.Index = i
	}
	for i, v := range t.Indexes { //删除索引
		if v == name {
			t.Indexes = append(t.Indexes[:i], t.Indexes[i+1:]...)
			break
		}
	}
	return t
}

func (t *TableInfo) ModifyColumn(col *ColumnInfo) *TableInfo {
	for i, v := range t.Columns {
		if v.Name == col.Name {
			col.Index = v.Index
			col.Expression = col.Name
			t.Columns[i] = col
			break
		}
	}
	return t
}

func (t *TableInfo) GetColumnIndex(name string) int {
	for _, v := range t.Columns {
		if v.Name == name {
			return v.Index
		}
	}
	return -1
}
