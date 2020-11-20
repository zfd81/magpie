package cluster

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/zfd81/magpie/server"

	"google.golang.org/grpc"

	pb "github.com/zfd81/magpie/proto/magpiepb"

	"github.com/zfd81/magpie/config"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/zfd81/magpie/etcd"
)

const (
	ClusterDirectory = "/cluster"
	LeaderDirectory  = "/leader"
	MemberDirectory  = "/members"
)

var (
	mu      sync.RWMutex
	leaseID clientv3.LeaseID
	node    *Node
	members = make(map[string]*Node)
	teams   = make(map[string]*Team)
	conf    = config.GetConfig()
)

func ClusterPath() string {
	return conf.Directory + ClusterDirectory
}

func LeaderPath() string {
	return ClusterPath() + LeaderDirectory
}

func MemberPath() string {
	return ClusterPath() + MemberDirectory
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

func InitMembers() error {
	//加载现有结点
	kvs, err := etcd.GetWithPrefix(MemberPath())
	if err != nil {
		return err
	}
	for _, kv := range kvs {
		addNode(kv.Key, kv.Value)
	}
	log.Info("Cluster member initialization successful.")
	return nil
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

	//节点注册并参与选举
	data, err := json.Marshal(node)
	if err != nil {
		return err
	}

	elect := concurrency.NewElection(session, MemberPath())
	go func() {
		//竞选 Leader，直到成为 Leader 函数才返回
		if err := elect.Campaign(context.Background(), string(data)); err != nil {
			log.WithFields(log.Fields{
				"id":  node.Id,
				"err": err,
			}).Error("Campaign leader error")
		} else {
			node.LeaderFlag = true
			//竞选成功后更改节点状态
			if _, err = etcd.PutWithLease(LeaderPath(), node.Id, session.Lease()); err != nil {
				fmt.Println(err)
			}
		}
	}()
	log.Info("Cluster registration successful.")
	return nil
}

func addNode(key []byte, value []byte) *Node {
	n := &Node{}
	err := json.Unmarshal(value, n)
	if err != nil {
		log.WithFields(log.Fields{
			"key": string(key),
			"err": err.Error(),
		}).Error("Failed to add node")
		return nil
	}
	if n.Id != node.Id {
		n.Connect()
	}
	mu.Lock()
	members[n.Id] = n       //添加到成员列表中
	team := GetTeam(n.Team) //获得团队
	team.AddMember(n)       //添加到团队中
	mu.Unlock()
	return n
}

func removeNode(key []byte) int {
	id := NodeId(key)
	value, found := members[id]
	if found {
		mu.Lock()
		team := GetTeam(value.Team) //获得团队
		team.RemoveMember(id)       //从团队中移除
		delete(members, id)         //从成员列表中移除
		mu.Unlock()
		return 1
	}
	return 0
}

func DataSync() error {
	log.Info("Start data synchronization:")
	team := GetTeam(node.Team)
	if !team.IsLeader(node) {
		leader := team.GetLeader()
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", leader.Address, leader.Port), grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(3*time.Second))
		if err != nil {
			log.WithFields(log.Fields{
				"Err": err.Error(),
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
			startTime := time.Now()
			request.Params["name"] = name
			stream, err := c.DataSync(context.Background(), request)
			if err != nil {
				log.WithFields(log.Fields{
					"Table": name,
					"Err":   err.Error(),
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
						"Table": name,
						"Err":   err.Error(),
					}).Error("Data synchronization error")
					return err
				}
				fields := strings.SplitN(res.Data, ",", size)
				key, row := tbl.RowData(fields)
				tbl.Insert(key, row)
			}
			log.WithFields(log.Fields{
				"Table":   name,
				"Elapsed": time.Since(startTime),
			}).Info("- Table data synchronization succeeded")
		}
	}
	log.Info("Data synchronization complete.")
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
