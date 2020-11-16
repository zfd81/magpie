package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zfd81/magpie/server"

	"github.com/zfd81/magpie/meta"
	pb "github.com/zfd81/magpie/proto/magpiepb"
)

type MetaServer struct{}

func (s *MetaServer) CreateTable(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {
	var info meta.TableInfo
	err := json.Unmarshal([]byte(request.Data), &info)
	if err != nil {
		return nil, err
	}
	server.GetDatabase("").CreateTable(info)
	msg := fmt.Sprintf("Table %s created successfully", info.Name)
	fmt.Println(msg)
	return &pb.RpcResponse{
		Code:    200,
		Message: msg,
	}, nil
}

func (s *MetaServer) DeleteTable(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {
	name := request.Params["name"]
	err := server.GetDatabase("").DeleteTable(name)
	if err != nil {
		return nil, err
	}
	msg := fmt.Sprintf("Table %s deleted successfully", name)
	fmt.Println(msg)
	return &pb.RpcResponse{
		Code:    200,
		Message: msg,
	}, nil
}

func (s *MetaServer) DescribeTable(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {
	tbl := server.GetDatabase("").DescribeTable(request.Params["name"])
	if tbl.Name == "" {
		return &pb.RpcResponse{
			Code: 200,
		}, nil
	}
	bytes, err := json.Marshal(tbl)
	if err != nil {
		return nil, err
	}
	return &pb.RpcResponse{
		Code: 200,
		Data: string(bytes),
	}, nil
}

func (s *MetaServer) ListTables(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {
	tbls := server.GetDatabase("").ListTables()
	bytes, err := json.Marshal(tbls)
	if err != nil {
		return nil, err
	}
	return &pb.RpcResponse{
		Code: 200,
		Data: string(bytes),
	}, nil
}
