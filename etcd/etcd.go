package etcd

import (
	"context"
	"time"

	"github.com/zfd81/magpie/config"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

type OperType int

//operType操作类型:0 create,1 delete,2 modify
type WatcherFunc func(operType OperType, key []byte, value []byte, createRevision int64, modRevision int64, version int64)

const (
	CREATE OperType = 0
	DELETE OperType = 1
	MODIFY OperType = 2
)

var (
	conf   = config.GetConfig()
	client *clientv3.Client
)

func Connect() error {
	clientConfig := clientv3.Config{
		Endpoints:   conf.Etcd.Endpoints,
		DialTimeout: time.Duration(conf.Etcd.DialTimeout) * time.Second,
	}
	cli, err := clientv3.New(clientConfig)
	if err != nil {
		return err
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	for _, v := range clientConfig.Endpoints {
		_, err = cli.Status(timeoutCtx, v)
		if err == nil {
			break
		}
	}
	if err != nil {
		return err
	}
	client = cli
	return nil
}

func GetClient() *clientv3.Client {
	return client
}

func Put(key, value string) (revision int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.Etcd.RequestTimeout)*time.Second)
	resp, err := client.Put(ctx, key, value)
	cancel()
	if err != nil {
		return -1, err
	}
	return resp.Header.Revision, nil
}

func PutWithLease(key, value string, leaseID clientv3.LeaseID) (revision int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.Etcd.RequestTimeout)*time.Second)
	resp, err := client.Put(ctx, key, value, clientv3.WithLease(leaseID))
	cancel()
	if err != nil {
		return -1, err
	}
	return resp.Header.Revision, nil
}

func PutWithTTL(key, value string, ttl int64) (revision int64, err error) {
	lease, err := client.Grant(context.TODO(), ttl)
	if err != nil {
		return -1, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.Etcd.RequestTimeout)*time.Second)
	resp, err := client.Put(ctx, key, value, clientv3.WithLease(lease.ID))
	cancel()
	if err != nil {
		return -1, err
	}
	return resp.Header.Revision, nil
}

func Del(key string) (revision int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.Etcd.RequestTimeout)*time.Second)
	resp, err := client.Delete(ctx, key)
	cancel()
	if err != nil {
		return -1, err
	}
	return resp.Header.Revision, nil
}

func DelWithPrefix(key string) (revision int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.Etcd.RequestTimeout)*time.Second)
	resp, err := client.Delete(ctx, key, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return -1, err
	}
	return resp.Header.Revision, nil
}

func Get(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.Etcd.RequestTimeout)*time.Second)
	resp, err := client.Get(ctx, key)
	cancel()
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) < 1 {
		return nil, nil
	}
	return resp.Kvs[0].Value, nil
}

func GetWithPrefix(prefix string) ([]*mvccpb.KeyValue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.Etcd.RequestTimeout)*time.Second)
	resp, err := client.Get(ctx, prefix, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return nil, err
	}
	return resp.Kvs, nil
}

func Watch(key string, watcher WatcherFunc) {
	rch := client.Watch(context.Background(), key)
	go func() {
		for wresp := range rch {
			for _, ev := range wresp.Events {
				oper := CREATE
				version := ev.Kv.Version
				if ev.Type == mvccpb.DELETE {
					oper = DELETE
				} else if version > 1 {
					oper = MODIFY
				}
				watcher(oper, ev.Kv.Key, ev.Kv.Value, ev.Kv.CreateRevision, ev.Kv.ModRevision, version)
			}
		}
	}()
}

func WatchWithPrefix(key string, watcher WatcherFunc) {
	rch := client.Watch(context.Background(), key, clientv3.WithPrefix())
	go func() {
		for wresp := range rch {
			for _, ev := range wresp.Events {
				oper := CREATE
				version := ev.Kv.Version
				if ev.Type == mvccpb.DELETE {
					oper = DELETE
				} else if version > 1 {
					oper = MODIFY
				}
				watcher(oper, ev.Kv.Key, ev.Kv.Value, ev.Kv.CreateRevision, ev.Kv.ModRevision, version)
			}
		}
	}()
}
