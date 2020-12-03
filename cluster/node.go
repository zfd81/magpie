package cluster

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/zfd81/magpie/server"

	log "github.com/sirupsen/logrus"

	pb "github.com/zfd81/magpie/proto/magpiepb"
	"google.golang.org/grpc"
)

type Node struct {
	Id          string `json:"id"`
	Address     string `json:"addr"`
	Port        int32  `json:"port"`
	Team        string `json:"team"`
	StartUpTime int64  `json:"start-up-time"`
	LeaderFlag  bool   `json:"leader-flag,omitempty"`
	logClient   pb.LogClient
}

func (n *Node) Connect() error {
	address := fmt.Sprintf("%s:%d", n.Address, n.Port)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	n.logClient = pb.NewLogClient(conn)
	return nil
}

func (n *Node) Log(entry *pb.Entry) error {
	log.WithFields(log.Fields{
		"to":        fmt.Sprintf("%s:%d", n.Address, n.Port),
		"timestamp": entry.Timestamp,
	}).Info(entry.Data)
	_, err := n.logClient.Apply(context.Background(), entry)
	return err
}

func (n *Node) DataSync(table string) error {
	team := GetTeam(n.Team)
	if !team.IsLeader(n) {
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
		startTime := time.Now()
		request.Params["name"] = table
		stream, err := c.DataSync(context.Background(), request)
		if err != nil {
			log.WithFields(log.Fields{
				"Table": table,
				"Err":   err.Error(),
			}).Panic("Data synchronization error")
			return err
		}
		tbl := db.GetTable(table)
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.WithFields(log.Fields{
					"Table": table,
					"Err":   err.Error(),
				}).Error("Data synchronization error")
				return err
			}
			row := tbl.NewRow()
			row.Load(res.Data, ",")
			tbl.Insert(row)
		}
		log.WithFields(log.Fields{
			"Table":   table,
			"Elapsed": time.Since(startTime),
		}).Info("Table data synchronization succeeded")
	}
	return nil
}

func NewNode(id string) *Node {
	return &Node{
		Id:         id,
		LeaderFlag: false,
	}
}

func NodeId(key []byte) string {
	return string(key)[len(MemberPath())+1:]
}
