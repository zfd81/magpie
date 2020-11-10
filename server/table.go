package server

import (
	"reflect"
	"strings"

	"github.com/zfd81/magpie/memory"

	"github.com/zfd81/magpie/sql"

	expr "github.com/zfd81/magpie/expression"

	"github.com/spf13/cast"

	"github.com/zfd81/magpie/meta"
)

type Table struct {
	meta.TableInfo
	primaryKeys   []*Column          //主键列
	columnMapping map[string]*Column //列映射
	cache         *memory.Cache
}

func (t *Table) init() {
	t.primaryKeys = make([]*Column, len(t.Keys))
	t.columnMapping = map[string]*Column{}
	for i, col := range t.Columns {
		col.Index = i
		col.Expression = col.Name
		t.columnMapping[col.Name] = NewColumn(*col)
	}
	for i, name := range t.Keys {
		t.primaryKeys[i] = t.columnMapping[name]
	}
	//for _, v := range t.DerivedCols {
	//	t.columnMapping[v.Name] = NewDerivedColumn(*v)
	//}
	t.cache = memory.New()
}

func (t *Table) GetColumn(name string) *Column {
	return t.columnMapping[name]
}

func (t *Table) NewRow() []interface{} {
	return make([]interface{}, len(t.Columns))
}

func (t *Table) RowData(data interface{}) (string, []interface{}) {
	value := reflect.Indirect(reflect.ValueOf(data))
	rowkey := strings.Builder{}
	for _, col := range t.primaryKeys {
		rowkey.WriteString(cast.ToString(value.Index(col.Index).Interface()))
	}
	row := t.NewRow()
	for _, col := range t.columnMapping {
		row[col.Index] = col.Value(value.Index(col.Index).Interface())
	}
	return rowkey.String(), row
}

func (t *Table) RowKey(data map[string]interface{}) string {
	key := strings.Builder{}
	for _, col := range t.primaryKeys {
		val, found := data[col.Name]
		if !found {
			return ""
		}
		key.WriteString(cast.ToString(val))
		delete(data, col.Name)
	}
	return key.String()
}

func (t *Table) BuildExprEnv(row []interface{}) map[string]interface{} {
	env := map[string]interface{}{}
	for _, col := range t.Columns {
		env[col.Name] = row[col.Index]
	}
	return env
}

func (t *Table) Insert(key string, row []interface{}) int {
	t.cache.Set(key, row)
	return 1
}

func (t *Table) DeleteByPrimaryKey(data map[string]interface{}) int {
	cnt := 0
	key := t.RowKey(data)
	if key != "" {
		t.cache.Remove(key)
		cnt++
	}
	return cnt
}

func (t *Table) UpdateByPrimaryKey(data map[string]interface{}) int {
	cnt := 0
	key := t.RowKey(data)
	if key != "" {
		row := t.cache.GetSlice(key)
		if len(row) > 0 {
			for k, v := range data {
				col := t.columnMapping[k]
				if col != nil {
					row[col.Index] = col.Value(v)
				}
			}
			t.cache.Set(key, row)
			cnt++
		}
	}
	return cnt
}

func (t *Table) readRow(data map[string]interface{}) []interface{} {
	key := t.RowKey(data)
	if key == "" {
		return nil
	}
	return t.cache.GetSlice(key)
}

func (t *Table) FindByPrimaryKey(columns []*sql.Field, conditions map[string]interface{}) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	key := t.RowKey(conditions)
	if key != "" {
		row := t.cache.GetSlice(key)
		if len(row) > 0 {
			env := t.BuildExprEnv(row)
			for _, column := range columns {
				val, err := expr.Eval(column.GetExpr(), env)
				if err != nil {
					return result, err
				}
				result[column.GetName()] = val
			}
		}
	}
	return result, nil
}

func (t *Table) Truncate() {
	t.cache.Clear()
}
