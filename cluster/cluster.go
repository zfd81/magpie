package cluster

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	pb "github.com/zfd81/magpie/proto/magpiepb"

	"github.com/zfd81/magpie/config"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/zfd81/magpie/util/etcd"
)

const (
	ClusterDirectory = "/cluster"
	LeaderDirectory  = "/leader"
	MemberDirectory  = "/members"
	LogDirectory     = "/mlog"
)

var (
	leaseID  clientv3.LeaseID
	node     *Node
	leaderId string
	members  = make(map[string]*Node)
	conf     = config.GetConfig()
)

func GetClusterPath() string {
	return conf.Directory + ClusterDirectory
}

func GetLeaderPath() string {
	return GetClusterPath() + LeaderDirectory
}

func GetMemberPath() string {
	return GetClusterPath() + MemberDirectory
}

func GetLogPath() string {
	return GetClusterPath() + LogDirectory
}

func GetNode() *Node {
	return node
}

func GetMembers() map[string]*Node {
	return members
}

func Register(startUpTime int64) error {
	ip, err := externalIP()
	if err != nil {
		return err
	}
	cli := etcd.GetClient()
	session, err := concurrency.NewSession(cli, concurrency.WithTTL(conf.Cluster.HeartbeatInterval))
	if err != nil {
		return err
	}
	node = NewNode(fmt.Sprintf("%x", session.Lease()))
	node.Address = ip.String()
	node.Port = conf.Port
	node.Team = conf.Team
	node.StartUpTime = startUpTime

	//获得集群领导者目录
	lpath := GetLeaderPath()

	//获得集群成员结点目录
	mpath := GetMemberPath()

	//监听集群leader结点变化
	etcd.Watch(lpath, func(operType etcd.OperType, key []byte, value []byte, createRevision int64, modRevision int64, version int64) {
		leaderId = string(value)
	})

	//监听集群结点变化
	etcd.WatchWithPrefix(mpath, clusterWatcher)

	StartScheduler() //启动计划程序

	//加载现有结点
	kvs, err := etcd.GetWithPrefix(mpath)
	if err == nil {
		for _, kv := range kvs {
			n := addNode(kv.Key, kv.Value)
			if n.LeaderFlag {
				leaderId = n.Id
			}
		}
	}

	//结点注册并参与选举
	data, err := json.Marshal(node)
	if err != nil {
		return err
	}
	elect := concurrency.NewElection(session, mpath)
	go func() {
		//竞选 Leader，直到成为 Leader 函数才返回
		if err := elect.Campaign(context.Background(), string(data)); err != nil {
			fmt.Println(err)
		} else {
			node.LeaderFlag = true
			if _, err = etcd.PutWithLease(lpath, node.Id, session.Lease()); err != nil {
				fmt.Println(err)
			}
		}
	}()
	return err
}

func clusterWatcher(operType etcd.OperType, key []byte, value []byte, createRevision int64, modRevision int64, version int64) {
	if operType == etcd.CREATE {
		addNode(key, value)
	} else if operType == etcd.MODIFY {
	} else if operType == etcd.DELETE {
		delete(members, NodeId(key))
	}
}

func addNode(key []byte, value []byte) *Node {
	node := &Node{}
	err := json.Unmarshal(value, node)
	if err == nil {
		node.Connect()
		members[node.Id] = node
	}
	return node
}

func NodeId(key []byte) string {
	return string(key)[len(GetMemberPath())+1:]
}

func externalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}
	return ip
}

func Broadcast(command string) {
	entry := &pb.Entry{
		Index:     0,
		Data:      command,
		Address:   node.Address,
		Port:      node.Port,
		Team:      conf.Team,
		Timestamp: time.Now().Format("20060102150405.000"),
	}
	for _, v := range members {
		if v.Id != node.Id && v.Team == node.Team {
			partner := v
			go func() {
				for i := 0; i < 3; i++ {
					err := partner.Log(entry)
					if err == nil {
						return
					}
				}
			}()
		}
	}
}
