package server

import (
	"strings"

	"github.com/spf13/cast"

	"github.com/zfd81/magpie/meta"
)

type PrimaryKeyFunc func(builder strings.Builder, fields []string)
type DataConversionFunc func(fields []string) ([]interface{}, error)

type Table struct {
	meta.TableInfo
	primaryKeyFunc     PrimaryKeyFunc
	dataConversionFunc DataConversionFunc
}

func (t *Table) RowID(fields []string) string {
	var builder strings.Builder
	builder.WriteString(t.Database.Name)
	builder.WriteString("_")
	builder.WriteString(t.Name)
	builder.WriteString("_")
	t.primaryKeyFunc(builder, fields)
	return builder.String()
}

func (t *Table) RowData(fields []string) ([]interface{}, error) {
	return t.dataConversionFunc(fields)
}

func BuildingPrimaryKeyFunc(table *Table) PrimaryKeyFunc {
	var positions []int
	for i, col := range table.Columns {
		for _, key := range table.Keys {
			if col.Name == key {
				positions = append(positions, i)
			}
		}
	}
	return func(builder strings.Builder, fields []string) {
		for _, v := range positions {
			builder.WriteString("_")
			builder.WriteString(cast.ToString(fields[v]))
		}
	}
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
