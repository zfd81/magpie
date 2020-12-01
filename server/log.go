package server

import (
	"encoding/json"
	"time"

	"github.com/golang/protobuf/proto"

	pb "github.com/zfd81/magpie/proto/magpiepb"

	log "github.com/sirupsen/logrus"

	"github.com/zfd81/magpie/store"
)

const (
	LogTableName = "_log"
)

type Log struct {
	Index uint64
	Team  []byte
	Data  []byte
}

var (
	storage store.Storage
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

func Append(entry *pb.Entry) uint64 {
	bytes, err := proto.Marshal(entry)
	if err != nil {
		log.Error("Marshal to struct error: %v", err)
		return 0
	}
	storage.Put(LogTableName, []byte(entry.Timestamp), bytes)
	return 1
}

func Remove(date string) error {
	return storage.DeleteWithPrefix(LogTableName, []byte(date))
}

func OpenLogStorage(path string) error {
	db, err := store.New(path)
	if err != nil {
		return err
	}
	err = db.CreateTable(LogTableName)
	if err != nil {
		return err
	}
	storage = db
	return nil
}
