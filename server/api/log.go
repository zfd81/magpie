package api

import (
	"context"
	"fmt"

	"github.com/zfd81/magpie/server"

	log "github.com/sirupsen/logrus"

	pb "github.com/zfd81/magpie/proto/magpiepb"
)

type LogServer struct{}

func (l *LogServer) Apply(ctx context.Context, entry *pb.Entry) (*pb.RpcResponse, error) {
	log.WithFields(log.Fields{
		"from":      fmt.Sprintf("%s:%d", entry.Address, entry.Port),
		"timestamp": entry.Timestamp,
	}).Info(entry.Data)
	res, err := server.Execute(entry.Data)
	if err != nil {
		return nil, err
	}
	server.Append(entry)
	resp := &pb.RpcResponse{
		Code: 200,
		Data: res,
	}
	return resp, nil
}
