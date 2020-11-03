package server

import (
	"fmt"
	"log"

	"github.com/zfd81/magpie/store"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/zfd81/magpie/memory"
	"github.com/zfd81/magpie/meta"
)

var (
	cache = memory.New()
	env   = NewInstance("magpie", "magpie")
	db    = env.CreateDatabase(meta.DatabaseInfo{
		Name: "taglib",
		Text: "taglib",
	})
)

func UUID() string {
	return uuid.New().String()
}

func write(key string, value []interface{}) int {
	cache.Set(key, value)
	return 1
}

func read(key string) []interface{} {
	return cache.GetSlice(key)
}

func remove(key string) int {
	cache.Remove(key)
	return 1
}

func InitTables() error {
	kvs, err := store.GetWithPrefix([]byte(db.GetPath()))
	cnt := 0
	if err == nil {
		for _, kv := range kvs {
			tbl, err := db.LoadTable(kv.Value)
			if err != nil {
				log.Fatalln(err.Error())
			}
			fmt.Printf("[INFO] Table %s initialized successfully \n", tbl.Name)
			cnt++
		}
		color.Green("[INFO] A total of %d tables were initialized \n", cnt)
	}
	return err
}
