package server

import (
	"github.com/zfd81/magpie/meta"
)

type Database struct {
	meta.DatabaseInfo
	Tables map[string]*Table
}

func (d *Database) CreateTable(info meta.TableInfo) *Table {
	tbl := &Table{
		TableInfo: meta.TableInfo{
			Name:     info.Name,
			Text:     info.Text,
			Columns:  info.Columns,
			Keys:     info.Keys,
			Database: d.DatabaseInfo,
		},
	}
	tbl.init()
	tbl.dataConversionFunc = BuildingDataConversionFunc(tbl) //构建行数据转换函数
	tbl.Store()                                              //添加元数据
	d.Tables[tbl.Name] = tbl
	return tbl
}

func (d *Database) DeleteTable(name string) error {
	tbl := d.Tables[name]
	tbl.Truncate()         //清空数据
	tbl.Remove()           //删除元数据
	delete(d.Tables, name) //删除表映射
	return nil
}

func (d *Database) DescribeTable(name string) meta.TableInfo {
	tbl := d.Tables[name]
	tbl.Load()
	return tbl.TableInfo
}
