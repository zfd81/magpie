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
	bufferSize = conf.BufferSize
	batchSize  = conf.WriteBatchSize
	poolSize   = conf.StoragePoolSize
	fleet      *Fleet
	counter    int32 = -1 //计数器
)

type boat struct {
	channel  chan []string
	capacity int
	data     []string
	mu       sync.RWMutex
}

func (b *boat) Channel() chan []string {
	return b.channel
}

func (b *boat) Size() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.data)
}

func (b *boat) Add(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.data = append(b.data, key)
	if len(b.data) == b.capacity {
		b.channel <- b.data
		b.data = []string{}
	}
}

func (b *boat) Clear() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	cnt := len(b.data)
	if cnt > 0 {
		b.channel <- b.data
		b.data = []string{}
	}
	return cnt
}

type Fleet struct {
	boats []*boat
	//capacity int
}

func (f *Fleet) GetBoat(index int) *boat {
	return f.boats[index]
}

func (f *Fleet) SetBoat(index int, boat *boat) {
	f.boats[index] = boat
}

func (f *Fleet) HasTask() bool {
	for _, b := range f.boats {
		if len(b.channel) > 0 {
			return true
		}
	}
	return false
}

func (f *Fleet) Clear() int {
	cnt := 0
	for _, b := range f.boats {
		cnt = cnt + b.Clear()
	}
	return cnt
}

func trigger(tbl *Table) {
	log.Info("Magpie server start write trigger")
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		<-ticker.C
		val := atomic.AddInt32(&counter, -1)
		if val < 0 {
			if fleet.HasTask() {
				atomic.StoreInt32(&counter, 20)
			} else {
				cnt := fleet.Clear()
				if cnt > 0 {
					log.Info("Trigger inserted data successfully: ", cnt)
				}
				break
			}
		}
	}
	log.Info("Magpie server exit write trigger")
}

func writer(tbl *Table, index int) {
	storage := tbl.db.GetStorage(index)
	stream := fleet.GetBoat(index).channel
	for keys := range stream {
		kvs := make([]store.KeyValue, 0, batchSize)
		for _, k := range keys {
			kvs = append(kvs, store.KeyValue{[]byte(k), []byte(tbl.cache.GetString(k))})
		}
		err := storage.BatchPut(tbl.Name, kvs)
		if err != nil {
			log.Error("Data storage error: ", err)
		} else {
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
	db            *Database          //存储
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

	fleet = &Fleet{
		boats: make([]*boat, poolSize),
	}
	for i := 0; i < poolSize; i++ {
		b := &boat{
			channel:  make(chan []string, bufferSize),
			capacity: batchSize,
		}
		fleet.SetBoat(i, b)
		index := i
		go writer(t, index)
	}
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
	fleet.GetBoat(store.PageIndex([]byte(key))).Add(key)
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
		err := t.db.storage.Delete(t.Name, []byte(key))
		if err == nil {
			return 1
		}
	}
	return 0
}

func (t *Table) UpdateByPrimaryKey(data map[string]string) int {
	key := t.rowKey(data)
	if key != "" {
		bytes, err := t.db.storage.Get(t.Name, []byte(key))
		if err == nil {
			row := t.NewRow()
			row.Load(string(bytes), FieldSeparator)
			for k, v := range data {
				col := t.columnMapping[k]
				if col != nil {
					row.Set(col.Index, v)
				}
			}
			err := t.db.storage.Put(t.Name, []byte(key), []byte(row.String()))
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
		bytes, err := t.db.storage.Get(t.Name, []byte(key))
		if err != nil || bytes == nil {
			if value, found := t.cache.Get(key); found {
				bytes = []byte(cast.ToString(value))
			} else {
				return result, nil
			}
		}
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
	return result, nil
}

func (t *Table) FindAll(f func(k, v string) error) error {
	return t.db.storage.Iterator(t.Name, f)
}

func (t *Table) Truncate() {
	t.db.storage.Truncate(t.Name)
}

func (t *Table) Status() (int, int, int) {
	colCount := len(t.Columns)
	size := 0
	rowCount := t.db.storage.Count(t.Name)
	return colCount, rowCount, size / 1024
}
