package mlog

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/zfd81/magpie/util/etcd"

	"github.com/zfd81/magpie/store"
)

const (
	LogTableName = "_log"
)

var (
	storage store.Storage
	Key     string
	Node    string
)

func timestamp() string {
	return time.Now().Format("20060102150405.000")
}

type Entry struct {
	Data      string
	Node      string
	Timestamp string
}

func (l *Entry) Marshal() (bytes []byte) {
	bytes, _ = json.Marshal(l)
	return
}

type Logger struct {
}

func Append(entry *Entry) uint64 {
	storage.Put([]byte(entry.Timestamp), entry.Marshal())
	return 1
}

func Remove(date string) error {
	return storage.DeleteWithPrefix([]byte(date))
}

func init() {
	db := store.NewBolt()
	err := db.Open("magpie-mlog.db")
	if err != nil {
		log.Panic(err)
	}
	err = db.CreateTable(LogTableName)
	if err != nil {
		log.Panic(err)
	}
	storage = db
}

func SendLog(data string) error {
	entry := &Entry{
		Data:      data,
		Node:      Node,
		Timestamp: timestamp(),
	}
	_, err := etcd.Put(Key, string(entry.Marshal()))
	return err
}
