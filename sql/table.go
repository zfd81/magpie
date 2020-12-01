package sql

import (
	"strings"

	"github.com/zfd81/magpie/store"

	expr "github.com/zfd81/magpie/sql/expression"

	"github.com/spf13/cast"

	"github.com/zfd81/magpie/meta"
)

const (
	FieldSeparator = "|"
)

type Row struct {
	data     []string
	keyFunc  func(data []string) string
	capacity int
}

func (r *Row) Append(data string) {
	r.data = append(r.data, data)
}

func (r *Row) Set(index int, val string) {
	r.data[index] = val
}

func (r *Row) Get(index int) string {
	return r.data[index]
}

func (r *Row) Data() []string {
	return r.data
}

func (r *Row) Key() string {
	return r.keyFunc(r.data)
}

func (r *Row) String() string {
	return strings.Join(r.data, FieldSeparator)
}

func (r *Row) KeyValue() store.KeyValue {
	return store.KeyValue{
		Key:   []byte(r.Key()),
		Value: []byte(r.String()),
	}
}

func (r *Row) Load(line, sep string) {
	r.data = strings.SplitN(line, sep, r.capacity)
}

type Table struct {
	meta.TableInfo
	primaryKeys   []*Column          //主键列
	columnMapping map[string]*Column //列映射
	db            store.Storage
	rowkeyFunc    func(data []string) string
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

	keyIndexs := make([]int, len(t.primaryKeys))
	for i, col := range t.primaryKeys {
		keyIndexs[i] = col.Index
	}
	t.rowkeyFunc = func(data []string) string {
		rowkey := strings.Builder{}
		for _, v := range keyIndexs {
			rowkey.WriteString(data[v])
		}
		return rowkey.String()
	}
}

func (t *Table) GetColumn(name string) *Column {
	return t.columnMapping[name]
}

func (t *Table) NewRow() *Row {
	row := &Row{
		data:     make([]string, len(t.Columns)),
		keyFunc:  t.rowkeyFunc,
		capacity: len(t.Columns),
	}
	return row
}

func (t *Table) RowKey(data map[string]string) string {
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

func (t *Table) BuildExprEnv(row []string) map[string]interface{} {
	env := map[string]interface{}{}
	for _, col := range t.Columns {
		env[col.Name] = ConversionFuncs[strings.ToUpper(col.DataType)](row[col.Index])
	}
	return env
}

func (t *Table) Insert(row *Row) int {
	rowkey := strings.Builder{}
	for _, col := range t.primaryKeys {
		rowkey.WriteString(row.Get(col.Index))
	}
	err := t.db.Put(t.Name, []byte(rowkey.String()), []byte(row.String()))
	if err == nil {
		return 1
	}
	return 0
}

func (t *Table) BatchInsert(rows []Row) int {
	kvs := make([]store.KeyValue, len(rows))
	for i, row := range rows {
		kvs[i] = row.KeyValue()
	}
	err := t.db.BatchPut(t.Name, kvs)
	if err == nil {
		return len(rows)
	}
	return 0
}

func (t *Table) DeleteByPrimaryKey(data map[string]string) int {
	rowkey := t.RowKey(data)
	if rowkey != "" {
		err := t.db.Delete(t.Name, []byte(rowkey))
		if err == nil {
			return 1
		}
	}
	return 0
}

func (t *Table) UpdateByPrimaryKey(data map[string]string) int {
	rowkey := t.RowKey(data)
	if rowkey != "" {
		bytes, err := t.db.Get(t.Name, []byte(rowkey))
		if err == nil {
			row := t.NewRow()
			row.Load(string(bytes), FieldSeparator)
			for k, v := range data {
				col := t.columnMapping[k]
				if col != nil {
					row.Set(col.Index, v)
				}
			}
			err := t.db.Put(t.Name, []byte(rowkey), []byte(row.String()))
			if err == nil {
				return 1
			}
		}
	}
	return 0
}

//func (t *Table) readRow(data map[string]interface{}) []interface{} {
//	key := t.RowKey(data)
//	if key == "" {
//		return nil
//	}
//	return t.cache.GetSlice(key)
//}

func (t *Table) FindByPrimaryKey(columns []*Field, conditions map[string]string) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	rowkey := t.RowKey(conditions)
	if rowkey != "" {
		bytes, err := t.db.Get(t.Name, []byte(rowkey))
		if err == nil && bytes != nil {
			row := t.NewRow()
			row.Load(string(bytes), FieldSeparator)
			env := t.BuildExprEnv(row.Data())
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

func (t *Table) FindAll(f func(k, v string) error) error {
	return t.db.Iterator(t.Name, f)
}

func (t *Table) Truncate() {
	t.db.Truncate(t.Name)
}

func (t *Table) Status() (int, int, int) {
	colCount := len(t.Columns)
	size := 0
	rowCount := t.db.Count(t.Name)
	return colCount, rowCount, size / 1024
}
