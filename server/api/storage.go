package api

import (
	"context"
	"encoding/json"

	"github.com/spf13/cast"
	pb "github.com/zfd81/magpie/proto/magpiepb"
	"github.com/zfd81/magpie/store"
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
