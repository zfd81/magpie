package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zfd81/magpie/cluster"
	"github.com/zfd81/magpie/etcd"

	"github.com/zfd81/magpie/server"

	pb "github.com/zfd81/magpie/proto/magpiepb"
)

type ClusterServer struct{}

func (c *ClusterServer) DataSync(request *pb.RpcRequest, stream pb.Cluster_DataSyncServer) error {
	name := request.Params["name"]
	tbl := server.GetDatabase("").GetTable(name)
	if tbl == nil {
		return fmt.Errorf("table %s does not exist", name)
	}
	return tbl.FindAll(func(k, v string) error {
		stream.Send(&pb.StreamResponse{Data: v})
		return nil
	})
}

func (c *ClusterServer) ListMembers(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {
	kvs, err := etcd.GetWithPrefix(cluster.MemberPath())
	if err != nil {
		return nil, err
	}
	var nodes []*cluster.Node
	var teams = make(map[string]*cluster.Team)
	for _, kv := range kvs {
		n := &cluster.Node{}
		err := json.Unmarshal(kv.Value, n)
		if err != nil {
			return nil, err
		}
		var team *cluster.Team
		value, found := teams[n.Team]
		if found {
			team = value
		} else {
			team = &cluster.Team{}
			teams[n.Team] = team
		}
		team.AddMember(n)
		nodes = append(nodes, n)
	}
	for _, n := range nodes {
		team := teams[n.Team]
		n.LeaderFlag = team.IsLeader(n)
	}
	bytes, err := json.Marshal(nodes)
	if err != nil {
		return nil, err
	}
	return &pb.RpcResponse{
		Code: 200,
		Data: string(bytes),
	}, nil
}

func (c *ClusterServer) MemberStatus(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {
	tbls := server.GetDatabase("").Tables
	tblInfos := []map[string]interface{}{}
	for _, tbl := range tbls {
		info := map[string]interface{}{}
		info["name"] = tbl.Name
		cc, rc, size := tbl.Status()
		info["colCount"] = cc
		info["rowCount"] = rc
		info["tblSize"] = size
		tblInfos = append(tblInfos, info)
	}
	bytes, err := json.Marshal(tblInfos)
	if err != nil {
		return nil, err
	}
	return &pb.RpcResponse{
		Code: 200,
		Data: string(bytes),
	}, nil
}
