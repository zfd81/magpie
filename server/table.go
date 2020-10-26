package server

import (
	"strings"

	"github.com/spf13/cast"

	"github.com/zfd81/magpie/meta"
)

type DataConversionFunc func(fields []string) ([]interface{}, error)

type Table struct {
	meta.TableInfo
	keyPrefix          string             //key前缀
	primaryKeys        []*Column          //主键
	columnMapping      map[string]*Column //列映射
	dataConversionFunc DataConversionFunc
}

func (t *Table) init() {
	t.keyPrefix = t.Database.Name + "." + t.Name //初始化key前缀
	t.primaryKeys = make([]*Column, len(t.Keys))
	t.columnMapping = map[string]*Column{}
	for _, v := range t.Columns {
		t.columnMapping[v.Name] = NewColumn(*v)
	}
	for i, name := range t.Keys {
		t.primaryKeys[i] = t.columnMapping[name]
	}
	for _, v := range t.DerivedCols {
		t.columnMapping[v.Name] = NewDerivedColumn(*v)
	}
}

func (t *Table) RowData(fields []string) (string, []interface{}, error) {
	rowkey := strings.Builder{}
	rowkey.WriteString(t.keyPrefix)
	for _, col := range t.primaryKeys {
		rowkey.WriteString("_")
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
	key.WriteString(t.keyPrefix)
	for _, col := range t.primaryKeys {
		val, found := data[col.Name]
		if !found {
			return ""
		}
		key.WriteString("_")
		key.WriteString(cast.ToString(val))
		delete(data, col.Name)
	}
	return key.String()
}

func (t *Table) Insert(key string, row []interface{}) int {
	return write(key, row)
}

func (t *Table) DeleteByPrimaryKey(data map[string]interface{}) int {
	key := t.RowKey(data)
	if key == "" {
		return 0
	}
	return remove(key)
}

func (t *Table) UpdateByPrimaryKey(data map[string]interface{}) int {
	key := t.RowKey(data)
	if key == "" {
		return 0
	}
	row := read(key)
	for k, v := range data {
		col := t.columnMapping[k]
		if col != nil {
			row[col.Index] = col.Value(v)
		}
	}
	return write(key, row)
}

func (t *Table) readRow(data map[string]interface{}) []interface{} {
	key := t.RowKey(data)
	if key == "" {
		return nil
	}
	return read(key)
}

func (t *Table) FindByPrimaryKey(names []string, data map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	row := t.RowKey(data)
	if len(row) > 0 {
		rowMap := map[string]interface{}{}
		for _, v := range t.Columns {
			col := t.columnMapping[v.Name]
			rowMap[col.Name] = row[col.Index]
		}
		for _, v := range t.DerivedCols {
			col := t.columnMapping[v.Name]
			val := col.Value(rowMap)
			rowMap[col.Name] = val
		}
		for _, name := range names {
			result[name] = rowMap[name]
		}
	}
	return result
}

func (t *Table) Truncate() int {
	return cache.RemoveWithPrefix(t.keyPrefix)
}

func BuildingDataConversionFunc(table *Table) DataConversionFunc {
	var funs []func(field interface{}) (interface{}, error)
	for _, col := range table.Columns {
		if strings.ToUpper(col.DataType) == DataTypeString {
			funs = append(funs, func(field interface{}) (interface{}, error) {
				return cast.ToStringE(field)
			})
		} else if strings.ToUpper(col.DataType) == DataTypeInteger {
			funs = append(funs, func(field interface{}) (interface{}, error) {
				return cast.ToIntE(field)
			})
		} else if strings.ToUpper(col.DataType) == DataTypeBool {
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
