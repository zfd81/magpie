package store

import (
	"fmt"
)

type StoragePool struct {
	pool []Storage
	path string
}

func (sp *StoragePool) Open(path string) error {
	for i := 0; i < poolSize; i++ {
		db := &boltdb{}
		err := db.Open(fmt.Sprintf("%s_%d.db", path, i))
		if err != nil {
			sp.pool = make([]Storage, poolSize)
			return err
		}
		sp.pool[i] = db
	}
	sp.path = path
	return nil
}

func (sp *StoragePool) CreateTable(name string) error {
	for _, db := range sp.pool {
		err := db.CreateTable(name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sp *StoragePool) BatchPut(table string, kvs []KeyValue) error {
	return nil
}

func (sp *StoragePool) Put(table string, key, value []byte) error {
	index := PageIndex(key)
	db := sp.GetStorage(index)
	if db == nil {
		return fmt.Errorf("Storage not found: index out of range [%d] with length %d", poolSize, poolSize)
	}
	return db.Put(table, key, value)
}

func (sp *StoragePool) Get(table string, key []byte) ([]byte, error) {
	index := PageIndex(key)
	db := sp.GetStorage(index)
	if db == nil {
		return nil, fmt.Errorf("Storage not found: index out of range [%d] with length %d", poolSize, poolSize)
	}
	return db.Get(table, key)
}

func (sp *StoragePool) GetWithPrefix(table string, prefix []byte) ([]*KeyValue, error) {
	var vals []*KeyValue
	for _, db := range sp.pool {
		kvs, err := db.GetWithPrefix(table, prefix)
		if err != nil {
			return nil, err
		}
		vals = append(vals, kvs...)
	}
	return vals, nil
}

func (sp *StoragePool) Delete(table string, key []byte) error {
	index := PageIndex(key)
	db := sp.GetStorage(index)
	if db == nil {
		return fmt.Errorf("Storage not found: index out of range [%d] with length %d", poolSize, poolSize)
	}
	return db.Delete(table, key)
}

func (sp *StoragePool) DeleteWithPrefix(table string, prefix []byte) error {
	for _, db := range sp.pool {
		err := db.DeleteWithPrefix(table, prefix)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sp *StoragePool) Truncate(table string) error {
	for _, db := range sp.pool {
		err := db.Truncate(table)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sp *StoragePool) Iterator(table string, f func(k, v string) error) error {
	for _, db := range sp.pool {
		if err := db.Iterator(table, f); err != nil {
			return err
		}
	}
	return nil
}

func (sp *StoragePool) IteratorWithPrefix(table string, prefix []byte, f func(k, v string) error) error {
	for _, db := range sp.pool {
		if err := db.IteratorWithPrefix(table, prefix, f); err != nil {
			return err
		}
	}
	return nil
}

func (sp *StoragePool) Count(table string) int {
	cnt := 0
	for _, db := range sp.pool {
		cnt = cnt + db.Count(table)
	}
	return cnt
}

func (sp *StoragePool) GetStorage(index int) Storage {
	if index < poolSize {
		return sp.pool[index]
	}
	return nil
}
