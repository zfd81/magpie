package cluster

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/fatih/color"

	log "github.com/sirupsen/logrus"

	"github.com/zfd81/magpie/server"

	"google.golang.org/grpc"

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
)

var (
	leaseID  clientv3.LeaseID
	node     *Node
	leaderId string
	members  = make(map[string]*Node)
	teams    = make(map[string]*Team)
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

func GetTeam(name string) *Team {
	value, found := teams[name]
	if found {
		return value
	}
	t := &Team{}
	teams[name] = t
	return t
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

	err = DataSync() //和同一团队的leader进行数据同步
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Panic("Data synchronization error")
	}

	//节点注册并参与选举
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
		addNode(key, value)
	} else if operType == etcd.DELETE {
		removeNode(key)
	}
}

func addNode(key []byte, value []byte) *Node {
	node := &Node{}
	err := json.Unmarshal(value, node)
	if err != nil {
		log.WithFields(log.Fields{
			"key": string(key),
			"err": err.Error(),
		}).Error("Failed to add node")
		return nil
	}
	node.Connect()
	members[node.Id] = node    //添加到成员列表中
	team := GetTeam(node.Team) //获得团队
	team.AddMember(node)       //添加到团队中
	return node
}

func removeNode(key []byte) int {
	id := NodeId(key)
	value, found := members[id]
	if found {
		team := GetTeam(value.Team) //获得团队
		team.RemoveMember(id)       //从团队中移除
		delete(members, id)         //从成员列表中移除
		return 1
	}
	return 0
}

func DataSync() error {
	log.Info(color.New(color.FgGreen).SprintFunc()("Start data synchronization >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"))
	team := GetTeam(node.Team)
	if !team.IsLeader(node) {
		leader := team.GetLeader()
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", leader.Address, leader.Port), grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(3*time.Second))
		if err != nil {
			log.WithFields(log.Fields{
				"err": err.Error(),
			}).Panic("Data synchronization error, did not connect")
			return err
		}
		defer conn.Close()
		c := pb.NewClusterClient(conn)
		request := &pb.RpcRequest{
			Params: map[string]string{},
		}
		db := server.GetDatabase("")
		for name, _ := range db.Tables {
			request.Params["name"] = name
			stream, err := c.DataSync(context.Background(), request)
			if err != nil {
				log.WithFields(log.Fields{
					"table": name,
					"err":   err.Error(),
				}).Panic("Data synchronization error")
				return err
			}
			tbl := db.GetTable(name)
			size := len(tbl.Columns)
			for {
				res, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.WithFields(log.Fields{
						"table": name,
						"err":   err.Error(),
					}).Error("Data synchronization error")
					return err
				}
				fields := strings.SplitN(res.Data, ",", size)
				key, row := tbl.RowData(fields)
				tbl.Insert(key, row)
			}
			log.WithFields(log.Fields{
				"table": name,
			}).Info("Table data synchronization succeeded")
		}
	}
	return nil
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
