package sql

import (
	"strings"

	"github.com/antonmedv/expr"

	"github.com/spf13/cast"
	"github.com/zfd81/magpie/meta"
)

var ConversionFuncs = map[string]func(val interface{}) interface{}{
	meta.DataTypeString: func(val interface{}) interface{} {
		return cast.ToString(val)
	},
	meta.DataTypeInteger: func(val interface{}) interface{} {
		return cast.ToInt(val)
	},
	meta.DataTypeFloat: func(val interface{}) interface{} {
		return cast.ToFloat64(val)
	},
	meta.DataTypeBool: func(val interface{}) interface{} {
		return cast.ToBool(val)
	},
}

type Column struct {
	meta.ColumnInfo
	handler   func(interface{}) interface{}
	IsDerived bool
}

func (c *Column) Value(val interface{}) interface{} {
	return ConversionFuncs[strings.ToUpper(c.DataType)](val)
}

func NewColumn(info meta.ColumnInfo) *Column {
	col := &Column{
		ColumnInfo: info,
	}
	ConversionFunc := ConversionFuncs[strings.ToUpper(col.DataType)]
	col.handler = func(val interface{}) interface{} {
		return ConversionFunc(val)
	}
	col.IsDerived = false
	return col
}

func NewDerivedColumn(info meta.ColumnInfo) *Column {
	col := &Column{
		ColumnInfo: info,
	}
	program, _ := expr.Compile(col.Expression)
	col.handler = func(env interface{}) interface{} {
		val, _ := expr.Run(program, env)
		return val
	}
	col.IsDerived = true
	return col
}
