package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/zfd81/magpie/cluster"

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
	command := entry.Data
	var result string
	if strings.HasPrefix(command, "load table") {
		splitted := strings.SplitN(command, " ", 3)
		cluster.CurrentNode().DataSync(splitted[2])
	} else {
		res, err := server.Execute(entry.Data)
		if err != nil {
			return nil, err
		}
		result = res
	}
	server.Append(entry)
	resp := &pb.RpcResponse{
		Code: 200,
		Data: result,
	}
	return resp, nil
}
