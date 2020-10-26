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
	tbl.dataConversionFunc = BuildingDataConversionFunc(tbl) //构建行数据转换函数
	d.Tables[tbl.Name] = tbl
	return tbl
}
