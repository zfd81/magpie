package server

import (
	"strings"

	"github.com/zfd81/magpie/memory"

	"github.com/zfd81/magpie/sql"

	expr "github.com/zfd81/magpie/expression"

	"github.com/spf13/cast"

	"github.com/zfd81/magpie/meta"
)

type DataConversionFunc func(fields []string) ([]interface{}, error)

type Table struct {
	meta.TableInfo
	primaryKeys        []*Column          //主键列
	columnMapping      map[string]*Column //列映射
	dataConversionFunc DataConversionFunc
	cache              *memory.Cache
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

func (t *Table) RowData(fields []string) (string, []interface{}, error) {
	rowkey := strings.Builder{}
	for _, col := range t.primaryKeys {
		rowkey.WriteString(cast.ToString(fields[col.Index]))
	}
	row, err := t.dataConversionFunc(fields)
	if err != nil {
		return "", nil, err
	}
	return rowkey.String(), row, nil
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
	key := t.RowKey(data)
	if key == "" {
		return 0
	}
	t.cache.Remove(key)
	return 1
}

func (t *Table) UpdateByPrimaryKey(data map[string]interface{}) int {
	key := t.RowKey(data)
	if key == "" {
		return 0
	}
	row := t.cache.GetSlice(key)
	for k, v := range data {
		col := t.columnMapping[k]
		if col != nil {
			row[col.Index] = col.Value(v)
		}
	}
	t.cache.Set(key, row)
	return 1
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
		env := t.BuildExprEnv(row)
		for _, column := range columns {
			n, e := ParseColumn(column)
			val, err := expr.Eval(e, env)
			if err != nil {
				return result, err
			}
			result[n] = val
		}
	}
	return result, nil
}

func (t *Table) Truncate() {
	t.cache.Clear()
}

func BuildingDataConversionFunc(table *Table) DataConversionFunc {
	var funs []func(field interface{}) (interface{}, error)
	for _, col := range table.Columns {
		if strings.ToUpper(col.DataType) == meta.DataTypeString {
			funs = append(funs, func(field interface{}) (interface{}, error) {
				return cast.ToStringE(field)
			})
		} else if strings.ToUpper(col.DataType) == meta.DataTypeInteger {
			funs = append(funs, func(field interface{}) (interface{}, error) {
				return cast.ToIntE(field)
			})
		} else if strings.ToUpper(col.DataType) == meta.DataTypeBool {
			funs = append(funs, func(field interface{}) (interface{}, error) {
				return cast.ToBoolE(field)
			})
		} else {
			funs = append(funs, func(field interface{}) (interface{}, error) {
				return cast.ToStringE(field)
			})
		}
	}
	return func(fields []string) ([]interface{}, error) {
		var data []interface{}
		for i, v := range funs {
			val, err := v(fields[i])
			if err != nil {
				return data, nil
			}
			data = append(data, val)
		}
		return data, nil
	}
}

func ParseColumn(field *sql.Field) (name string, expr string) {
	if field.As == "" {
		name = field.Name
		expr = field.Name
	} else {
		name = field.As
		expr = field.Expr
	}
	return
}
