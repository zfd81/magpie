package cluster

import (
	log "github.com/sirupsen/logrus"
	"github.com/zfd81/magpie/etcd"
)

//监听集群Leader节点变化
func WatchLeader() {
	etcd.Watch(LeaderPath(), func(operType etcd.OperType, key []byte, value []byte, createRevision int64, modRevision int64, version int64) {
		log.Info("Team leader election completed")
	})
}

//监听集群Members节点变化
func WatchMembers() {
	etcd.WatchWithPrefix(MemberPath(), func(operType etcd.OperType, key []byte, value []byte, createRevision int64, modRevision int64, version int64) {
		if operType == etcd.CREATE {
			log.WithFields(log.Fields{
				"id": NodeId(key),
			}).Info("Add node")
			addNode(key, value)
		} else if operType == etcd.MODIFY {
			log.WithFields(log.Fields{
				"id": NodeId(key),
			}).Info("Modify node")
			addNode(key, value)
		} else if operType == etcd.DELETE {
			log.WithFields(log.Fields{
				"id": NodeId(key),
			}).Info("Delete node")
			removeNode(key)
		}
	})
}
