package server

import (
	"github.com/google/uuid"
	"github.com/zfd81/magpie/memory"
	"github.com/zfd81/magpie/meta"
)

const (
	DataTypeString  = "STRING"
	DataTypeInteger = "INT"
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

func write(key string, value []interface{}) {
	cache.Set(key, value)
}

func read(key string) []interface{} {
	return cache.GetSlice(key)
}
