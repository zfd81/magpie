package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/spf13/cast"

	"github.com/zfd81/magpie/store"

	pb "github.com/zfd81/magpie/api/magpiepb"

	"github.com/zfd81/magpie/meta"
)

type StorageServer struct{}

func (s *StorageServer) Get(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {
	key := request.Params["key"]
	prefix := cast.ToBool(request.Params["prefix"])
	data := map[string]string{}
	if prefix {
		kvs, err := store.GetWithPrefix([]byte(key))
		if err != nil {
			return nil, err
		}
		for _, kv := range kvs {
			data[string(kv.Key)] = string(kv.Value)
		}
	} else {
		bytes, err := store.Get([]byte(key))
		if err != nil {
			return nil, err
		}
		if bytes != nil {
			data[key] = string(bytes)
		}
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return &pb.RpcResponse{
		Code: 200,
		Data: string(bytes),
	}, nil
}

func (s *StorageServer) Count(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {
	cnt := store.Count()
	return &pb.RpcResponse{
		Code: 200,
		Data: cast.ToString(cnt),
	}, nil
}

type TableServer struct{}

func (s *TableServer) CreateTable(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {
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

func (s *TableServer) DeleteTable(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {
	name := request.Params["name"]
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

func (s *TableServer) DescribeTable(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {
	tbl := db.DescribeTable(request.Params["name"])
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

func (s *TableServer) ListTables(ctx context.Context, request *pb.RpcRequest) (*pb.RpcResponse, error) {
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

type MagpieServer struct{}

func (s *MagpieServer) Load(stream pb.Magpie_LoadServer) error {
	startTime := time.Now()
	fmt.Println(startTime)
	cnt := 0
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			//endTime := time.Now()
			return stream.SendAndClose(&pb.RpcResponse{
				Code: 200,
				Data: cast.ToString(cnt),
			})
		}
		if err != nil {
			return err
		}
		cnt++
		fmt.Println(r.Data)
	}
	fmt.Println(cnt)
	return nil
}
