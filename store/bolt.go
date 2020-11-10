package store

import (
	"bytes"
	"fmt"

	"github.com/boltdb/bolt"
)

const (
	// Permissions to use on the db file. This is only used if the
	// database file does not exist and needs to be created.
	dbFileMode = 0600
	tableName  = "_tag"
)

type boltdb struct {
	db   *bolt.DB
	path string
}

func (db *boltdb) Open(path string) error {
	database, err := bolt.Open(path, dbFileMode, nil)
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

func (db *boltdb) Put(key, value []byte) error {
	return db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tableName))
		return b.Put(key, value)
	})
}

func (db *boltdb) Get(key []byte) ([]byte, error) {
	var val []byte
	err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tableName))
		val = b.Get(key)
		return nil
	})
	return val, err
}

func (db *boltdb) GetWithPrefix(prefix []byte) ([]*KeyValue, error) {
	var vals []*KeyValue
	err := db.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(tableName)).Cursor()
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

func (db *boltdb) Delete(key []byte) error {
	return db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tableName))
		return b.Delete(key)
	})
}

func (db *boltdb) DeleteWithPrefix(prefix []byte) error {
	return db.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tableName))
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

func (db *boltdb) Count() int {
	cnt := 0
	db.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(tableName)).Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			cnt++
		}
		return nil
	})
	return cnt
}

func NewBolt() *boltdb {
	return &boltdb{}
}
