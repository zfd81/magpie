package server

import (
	"github.com/google/uuid"
	"github.com/zfd81/magpie/memory"
	"github.com/zfd81/magpie/meta"
)

const (
	DataTypeString  = "STRING"
	DataTypeInteger = "INT"
	DataTypeFloat   = "FLOAT"
	DataTypeBool    = "BOOL"
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
