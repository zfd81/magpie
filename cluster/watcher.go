package cluster

import (
	log "github.com/sirupsen/logrus"
	"github.com/zfd81/magpie/etcd"
)

//监听集群Leader节点变化
func WatchLeader() {
	etcd.Watch(LeaderPath(), func(operType etcd.OperType, key []byte, value []byte, createRevision int64, modRevision int64, version int64) {
		if operType == etcd.CREATE || operType == etcd.MODIFY {
			log.WithFields(log.Fields{
				"id": string(value),
			}).Info("Team leader election completed")
		}
	})
}

//监听集群Members节点变化
func WatchMembers() {
	etcd.WatchWithPrefix(MemberPath(), func(operType etcd.OperType, key []byte, value []byte, createRevision int64, modRevision int64, version int64) {
		if operType == etcd.CREATE {
			log.WithFields(log.Fields{
				"id": NodeId(key),
			}).Info("New member join")
			addNode(key, value)
		} else if operType == etcd.MODIFY {
			log.WithFields(log.Fields{
				"id": NodeId(key),
			}).Info("Member status change")
			addNode(key, value)
		} else if operType == etcd.DELETE {
			log.WithFields(log.Fields{
				"id": NodeId(key),
			}).Info("Member exit")
			removeNode(key)
		}
	})
}
