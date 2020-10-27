package server

import (
	"context"
	"encoding/json"
	"fmt"

	pb "github.com/zfd81/magpie/api/magpiepb"

	"github.com/zfd81/magpie/meta"
)

type Server struct{}

func (s *Server) CreateTable(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {
	var info meta.TableInfo
	err := json.Unmarshal([]byte(request.Data), &info)
	if err != nil {
		return nil, err
	}
	db.CreateTable(info)
	msg := fmt.Sprintf("Table %s created successfully", info.Name)
	fmt.Println(msg)
	return &pb.RpcResponse{
		Code:    200,
		Message: msg,
	}, nil
}

func (s *Server) DeleteTable(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {
	name := request.Attributes["name"]
	err := db.DeleteTable(name)
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

func (s *Server) DescribeTable(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {
	tbl := db.DescribeTable(request.Attributes["name"])
	bytes, err := json.Marshal(tbl)
	if err != nil {
		return nil, err
	}
	return &pb.RpcResponse{
		Code: 200,
		Data: string(bytes),
	}, nil
}

func (s *Server) ListTables(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {
	tbls := db.ListTables()
	bytes, err := json.Marshal(tbls)
	if err != nil {
		return nil, err
	}
	return &pb.RpcResponse{
		Code: 200,
		Data: string(bytes),
	}, nil
}
