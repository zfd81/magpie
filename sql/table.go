package sql

import (
	"strings"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/zfd81/magpie/memory"

	"github.com/zfd81/magpie/store"

	expr "github.com/zfd81/magpie/sql/expression"

	"github.com/spf13/cast"

	"github.com/zfd81/magpie/meta"
)

var (
	bufferSize       = conf.BufferSize
	batchSize        = conf.WriteBatchSize
	dataStream       = make(chan []string, bufferSize)
	batch            = &boat{}
	counter    int32 = -1 //计数器
)

type boat struct {
	//capacity int
	data []string
	mu   sync.RWMutex
}

func (b boat) Size() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.data)
}

func (b *boat) Add(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.data = append(b.data, key)
}

func (b *boat) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.data = []string{}
}

func (b *boat) Data() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.data
}

func trigger(tbl *Table) {
	log.Info("Magpie server start trigger")
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		<-ticker.C
		val := atomic.AddInt32(&counter, -1)
		if val < 0 {
			if len(dataStream) > 0 {
				atomic.StoreInt32(&counter, 20)
			} else {
				var kvs []store.KeyValue
				for _, k := range batch.data {
					kvs = append(kvs, store.KeyValue{[]byte(k), []byte(tbl.cache.GetString(k))})
				}
				if len(kvs) > 0 {
					err := tbl.db.BatchPut(tbl.Name, kvs)
					if err == nil {
						for _, kv := range kvs {
							tbl.cache.Remove(string(kv.Key))
						}
					}
					log.Info("Trigger inserted data successfully: ", len(kvs))
				}
				break
			}
		}
	}
	log.Info("Magpie server exit trigger")
}

func writer(tbl *Table) {
	for keys := range dataStream {
		kvs := make([]store.KeyValue, batchSize)
		for i, k := range keys {
			kvs[i] = store.KeyValue{[]byte(k), []byte(tbl.cache.GetString(k))}
		}
		err := tbl.db.BatchPut(tbl.Name, kvs)
		if err == nil {
			for _, k := range keys {
				tbl.cache.Remove(k)
			}
		}
	}
}

type Table struct {
	meta.TableInfo
	primaryKeys   []*Column          //主键列
	columnMapping map[string]*Column //列映射
	db            store.Storage      //存储
	cache         *memory.Cache      //缓存
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

	t.cache = memory.New()
	go writer(t)
}

func (t *Table) rowKey(data map[string]string) string {
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

func (t *Table) buildExprEnv(row []string) map[string]interface{} {
	env := map[string]interface{}{}
	for _, col := range t.Columns {
		env[col.Name] = ConversionFuncs[strings.ToUpper(col.DataType)](row[col.Index])
	}
	return env
}

func (t *Table) NewRow() *Row {
	row := &Row{
		data:     make([]string, len(t.Columns)),
		capacity: len(t.Columns),
	}
	return row
}

func (t *Table) GetColumn(name string) *Column {
	return t.columnMapping[name]
}

func (t *Table) Insert(row *Row) int {
	key := t.rowkeyFunc(row.data)
	t.cache.Set(key, row.String())
	batch.Add(key)
	if batch.Size() == batchSize {
		dataStream <- batch.Data()
		batch.Clear()
	}
	if atomic.CompareAndSwapInt32(&counter, -1, 10) {
		go trigger(t)
	} else {
		atomic.StoreInt32(&counter, 10)
	}
	return 1
}

func (t *Table) DeleteByPrimaryKey(data map[string]string) int {
	key := t.rowKey(data)
	if key != "" {
		_, found := t.cache.Get(key)
		if found {
			t.cache.Remove(key)
			return 1
		}
		err := t.db.Delete(t.Name, []byte(key))
		if err == nil {
			return 1
		}
	}
	return 0
}

func (t *Table) UpdateByPrimaryKey(data map[string]string) int {
	key := t.rowKey(data)
	if key != "" {
		bytes, err := t.db.Get(t.Name, []byte(key))
		if err == nil {
			row := t.NewRow()
			row.Load(string(bytes), FieldSeparator)
			for k, v := range data {
				col := t.columnMapping[k]
				if col != nil {
					row.Set(col.Index, v)
				}
			}
			err := t.db.Put(t.Name, []byte(key), []byte(row.String()))
			if err == nil {
				return 1
			}
		}
	}
	return 0
}

func (t *Table) FindByPrimaryKey(columns []*Field, conditions map[string]string) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	key := t.rowKey(conditions)
	if key != "" {
		bytes, err := t.db.Get(t.Name, []byte(key))
		if err != nil || bytes == nil {
			value, found := t.cache.Get(key)
			if found {
				bytes = []byte(cast.ToString(value))
			}
		}
		if bytes != nil {
			row := t.NewRow()
			row.Load(string(bytes), FieldSeparator)
			env := t.buildExprEnv(row.Data())
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
