package sql

import (
	"encoding/json"

	"github.com/zfd81/magpie/store"

	"github.com/zfd81/magpie/meta"
)

type Database struct {
	meta.DatabaseInfo
	Tables  map[string]*Table
	storage store.Storage
}

func (d *Database) FileName() string {
	return d.Instance.Name + "-" + d.Name
}

func (d *Database) CreateTable(info meta.TableInfo) *Table {
	tbl := &Table{
		TableInfo: meta.TableInfo{
			Name:     info.Name,
			Text:     info.Text,
			Columns:  info.Columns,
			Keys:     info.Keys,
			Indexes:  info.Indexes,
			Database: d.DatabaseInfo,
		},
		db: d.storage,
	}
	tbl.init()
	tbl.Store() //保存元数据
	d.Tables[tbl.Name] = tbl
	d.storage.CreateTable(tbl.Name)
	return tbl
}

func (d *Database) LoadTable(bytes []byte) (*Table, error) {
	info := &meta.TableInfo{}
	err := json.Unmarshal(bytes, info)
	if err != nil {
		return nil, err
	}
	info.Database = d.DatabaseInfo
	tbl := &Table{
		TableInfo: *info,
		db:        d.storage,
	}
	tbl.init()
	d.Tables[tbl.Name] = tbl
	d.storage.CreateTable(tbl.Name)
	return tbl, nil
}

func (d *Database) DeleteTable(name string) error {
	tbl := d.Tables[name]
	if tbl != nil {
		tbl.Truncate()         //清空数据
		tbl.Remove()           //删除元数据
		delete(d.Tables, name) //删除表映射
	}
	return nil
}

func (d *Database) GetTable(name string) *Table {
	return d.Tables[name]
}

func (d *Database) DescribeTable(name string) meta.TableInfo {
	path := d.GetPath() + meta.PathSeparator + name + meta.TableSuffix
	tbl := meta.TableInfo{}
	meta.LoadMetadata(&tbl, path)
	return tbl
}

//func (d *Database) GetStorage(index int) store.Storage {
//	return d.storagePool.GetStorage(index)
//}
