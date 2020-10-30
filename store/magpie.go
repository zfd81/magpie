package store

import (
	"github.com/zfd81/magpie/util/etcd"
)

const storagePath = "@magpie"

var magpieDB Storage

func Put(key, value []byte) error {
	_, err := etcd.Put(storagePath+string(key), string(value))
	return err
}

func Get(key []byte) ([]byte, error) {
	return etcd.Get(storagePath + string(key))
}

func GetWithPrefix(prefix []byte) ([]*KeyValue, error) {
	kvs, err := etcd.GetWithPrefix(storagePath + string(prefix))
	if err != nil {
		return nil, err
	}
	var vals []*KeyValue
	for _, kv := range kvs {
		vals = append(vals, &KeyValue{
			Key:   kv.Key[len(storagePath):],
			Value: kv.Value,
		})
	}
	return vals, nil
}

func Delete(key []byte) error {
	_, err := etcd.Del(storagePath + string(key))
	return err
}

func DeleteWithPrefix(prefix []byte) error {
	_, err := etcd.DelWithPrefix(storagePath + string(prefix))
	return err
}

func Count() int {
	kvs, err := etcd.GetWithPrefix(storagePath)
	if err != nil {
		return -1
	}
	return len(kvs)
}
