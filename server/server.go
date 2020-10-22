package server

import "github.com/zfd81/magpie/meta"

func CreateTable(info meta.TableInfo) *Table {
	return db.CreateTable(info)
}
