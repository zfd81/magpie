package store

import (
	"bytes"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

const (
	// Permissions to use on the db file. This is only used if the
	// database file does not exist and needs to be created.
	dbFileMode = 0600
)

type boltdb struct {
	db   *bolt.DB
	path string
}

func (db *boltdb) Open(path string) error {
	database, err := bolt.Open(path, dbFileMode, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	db.db = database
	db.path = path
	return nil
}

func (db *boltdb) CreateTable(name string) error {
	return db.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return fmt.Errorf("create table: %s", err)
		}
		return nil
	})
}

func (db *boltdb) Put(table string, key, value []byte) error {
	return db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		return b.Put(key, value)
	})
}

func (db *boltdb) BatchPut(table string, kvs []KeyValue) error {
	return db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		for _, kv := range kvs {
			if err := b.Put(kv.Key, kv.Value); err != nil {
				return err
			}
		}
		return nil
	})
}

func (db *boltdb) Get(table string, key []byte) ([]byte, error) {
	var val []byte
	err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		val = b.Get(key)
		return nil
	})
	return val, err
}

func (db *boltdb) GetWithPrefix(table string, prefix []byte) ([]*KeyValue, error) {
	var vals []*KeyValue
	err := db.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(table)).Cursor()
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			vals = append(vals, &KeyValue{
				Key:   k,
				Value: v,
			})
		}
		return nil
	})
	return vals, err
}

func (db *boltdb) Delete(table string, key []byte) error {
	return db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		return b.Delete(key)
	})
}

func (db *boltdb) DeleteWithPrefix(table string, prefix []byte) error {
	return db.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		c := b.Cursor()
		for k, _ := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = c.Next() {
			err := b.Delete(k)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (db *boltdb) Truncate(table string) error {
	return db.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			if err := b.Delete(k); err != nil {
				return err
			}
		}
		return nil
	})
}

func (db *boltdb) Iterator(table string, f func(k, v string) error) error {
	return db.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(table)).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			err := f(string(k), string(v))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (db *boltdb) IteratorWithPrefix(table string, prefix []byte, f func(k, v string) error) error {
	return db.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(table)).Cursor()
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			err := f(string(k), string(v))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (db *boltdb) Count(table string) int {
	cnt := 0
	db.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(table)).Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			cnt++
		}
		return nil
	})
	return cnt
}
