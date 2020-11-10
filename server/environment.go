package server

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/zfd81/magpie/sql"

	"github.com/spf13/cast"

	"github.com/zfd81/magpie/store"

	"github.com/google/uuid"
	"github.com/zfd81/magpie/meta"
)

var (
	env = NewInstance("magpie", "magpie")
	db  = env.CreateDatabase(meta.DatabaseInfo{
		Name: "taglib",
		Text: "taglib",
	})
)

func UUID() string {
	return uuid.New().String()
}

func InitTables() error {
	kvs, err := store.GetWithPrefix([]byte(db.GetPath()))
	cnt := 0
	if err == nil {
		for _, kv := range kvs {
			tbl, err := db.LoadTable(kv.Value)
			if err != nil {
				log.Panic(err)
			}
			log.WithFields(log.Fields{
				"name": tbl.Name,
			}).Info("Table initialized successfully")
			cnt++
		}
		log.Infof("A total of %d tables were initialized \n", cnt)
	}
	return err
}

func Execute(query string) (string, error) {
	stmt, err := sql.Parse(query)
	if err != nil {
		return "", err
	}
	switch stmt := stmt.(type) {
	case *sql.SelectStatement:
		tableName := stmt.From[0].Name
		tbl := db.GetTable(tableName)
		if tbl == nil {
			return "", fmt.Errorf("table %s does not exist", tableName)
		}
		conditions := map[string]interface{}{}
		for _, v := range stmt.Where {
			conditions[v.Name] = string(v.Value)
		}
		result, err := tbl.FindByPrimaryKey(stmt.Select, conditions)
		if err != nil {
			return "", err
		}
		bytes, err := json.Marshal(result)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	case *sql.InsertStatement:
		tableName := stmt.Table
		tbl := db.GetTable(tableName)
		if tbl == nil {
			return "", fmt.Errorf("table %s does not exist", tableName)
		}
		cnt := 0
		cols := stmt.Columns
		if cols == nil {
			for _, row := range stmt.Rows {
				key, data := tbl.RowData(*row)
				if err != nil {
					return "", err
				}
				tbl.Insert(key, data)
				cnt++
			}
		} else {
			for _, row := range stmt.Rows {
				datas := tbl.NewRow()
				for i, name := range cols {
					col := tbl.GetColumn(name)
					datas[col.Index] = row.Get(i)
				}
				key, data := tbl.RowData(datas)
				if err != nil {
					return "", err
				}
				tbl.Insert(key, data)
				cnt++
			}

		}
		return cast.ToString(cnt), nil
	case *sql.DeleteStatement:
		tableName := stmt.Table
		tbl := db.GetTable(tableName)
		if tbl == nil {
			return "", fmt.Errorf("table %s does not exist", tableName)
		}
		conditions := map[string]interface{}{}
		for _, v := range stmt.Where {
			conditions[v.Name] = string(v.Value)
		}
		return cast.ToString(tbl.DeleteByPrimaryKey(conditions)), nil
	case *sql.UpdateStatement:
		tableName := stmt.Table
		tbl := db.GetTable(tableName)
		if tbl == nil {
			return "", fmt.Errorf("table %s does not exist", tableName)
		}
		conditions := map[string]interface{}{}
		for _, v := range stmt.Fields {
			conditions[v.Name] = string(v.Expr)
		}
		for _, v := range stmt.Where {
			conditions[v.Name] = string(v.Value)
		}
		return cast.ToString(tbl.UpdateByPrimaryKey(conditions)), nil
	default:
		return "", fmt.Errorf("unsupported syntax: ", query)
	}
}
