package cluster

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	pb "github.com/zfd81/magpie/proto/magpiepb"
	"google.golang.org/grpc"
)

type Node struct {
	Id          string `json:"id"`
	Address     string `json:"addr"`
	Port        int64  `json:"port"`
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

func NewNode(id string) *Node {
	return &Node{
		Id:         id,
		LeaderFlag: false,
	}
}

func NodeId(key []byte) string {
	return string(key)[len(MemberPath())+1:]
}
