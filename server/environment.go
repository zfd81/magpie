package server

import (
	"fmt"
	"log"

	"github.com/zfd81/magpie/store"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/zfd81/magpie/meta"
)

var (
	env = NewInstance("magpie", "magpie")
	db  = env.CreateDatabase(meta.DatabaseInfo{
		Name: "taglib",
		Text: "taglib",
	})
)

func UUID() string {
	return uuid.New().String()
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
