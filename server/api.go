package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/zfd81/magpie/sql"

	"github.com/spf13/cast"

	"github.com/zfd81/magpie/store"

	pb "github.com/zfd81/magpie/proto/magpiepb"

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
	r, err := stream.Recv()
	name := r.Data
	tbl := db.GetTable(name)
	if tbl == nil {
		return stream.SendAndClose(&pb.LoadResponse{
			Code:    400,
			Message: fmt.Sprintf("table %s does not exist", name),
		})
	}
	size := len(tbl.Columns)
	var cnt int64 = 0
	for {
		r, err = stream.Recv()
		if err == io.EOF {
			endTime := time.Now()
			return stream.SendAndClose(&pb.LoadResponse{
				Code:        200,
				Name:        "",
				StartTime:   startTime.Format("2006-01-02 15:04:05.000"),
				EndTime:     endTime.Format("2006-01-02 15:04:05.000"),
				ElapsedTime: int64(time.Since(startTime)),
				RecordCount: cnt,
				Message:     fmt.Sprintf("Data loading complete"),
			})
		}
		if err != nil {
			return err
		}
		fields := strings.SplitN(r.Data, ",", size)
		cnt++
		key, row, err := tbl.RowData(fields)
		if err != nil {
			return stream.SendAndClose(&pb.LoadResponse{
				Code:    400,
				Message: fmt.Sprintf("Data format error in line %d", cnt),
			})
		}
		tbl.Insert(key, row)
	}
	return nil
}

func (s *MagpieServer) Execute(ctx context.Context, request *pb.QueryRequest) (*pb.QueryResponse, error) {

	stmt, err := sql.Parse(request.Sql)
	if err != nil {
		return nil, err
	}

	tableName := stmt.From[0].Name
	tbl := db.GetTable(tableName)
	if tbl == nil {
		return nil, fmt.Errorf("table %s does not exist", tableName)
	}
	conditions := map[string]interface{}{}
	for _, v := range stmt.Where {
		conditions[v.Name] = string(v.Value)
	}

	result, err := tbl.FindByPrimaryKey(stmt.Select, conditions)
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return &pb.QueryResponse{
		Code: 200,
		Data: string(bytes),
	}, nil
}
