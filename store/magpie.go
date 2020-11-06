package store

import (
	"github.com/zfd81/magpie/config"
	"github.com/zfd81/magpie/util/etcd"
)

var (
	magpieDB Storage
	conf     = config.GetConfig()
)

func Put(key, value []byte) error {
	_, err := etcd.Put(conf.Directory+string(key), string(value))
	return err
}

func Get(key []byte) ([]byte, error) {
	return etcd.Get(conf.Directory + string(key))
}

func GetWithPrefix(prefix []byte) ([]*KeyValue, error) {
	kvs, err := etcd.GetWithPrefix(conf.Directory + string(prefix))
	if err != nil {
		return nil, err
	}
	var vals []*KeyValue
	for _, kv := range kvs {
		vals = append(vals, &KeyValue{
			Key:   kv.Key[len(conf.Directory):],
			Value: kv.Value,
		})
	}
	return vals, nil
}

func Delete(key []byte) error {
	_, err := etcd.Del(conf.Directory + string(key))
	return err
}

func DeleteWithPrefix(prefix []byte) error {
	_, err := etcd.DelWithPrefix(conf.Directory + string(prefix))
	return err
}

func Count() int {
	kvs, err := etcd.GetWithPrefix(conf.Directory)
	if err != nil {
		return -1
	}
	return len(kvs)
}
