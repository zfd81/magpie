package api

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/zfd81/magpie/server"

	"github.com/zfd81/magpie/cluster"
	pb "github.com/zfd81/magpie/proto/magpiepb"
)

type MagpieServer struct{}

func (s *MagpieServer) Load(stream pb.Magpie_LoadServer) error {
	startTime := time.Now()
	r, err := stream.Recv()
	name := r.Data
	tbl := server.GetDatabase("").GetTable(name)
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
			log.WithFields(log.Fields{
				"table":   name,
				"elapsed": time.Since(startTime),
			}).Info("Data loaded successfully")
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
		key, row := tbl.RowData(fields)
		tbl.Insert(key, row)
	}
	return nil
}

func (s *MagpieServer) Execute(ctx context.Context, request *pb.QueryRequest) (*pb.QueryResponse, error) {
	res, err := server.Execute(request.Sql)
	if err != nil {
		return nil, err
	}
	if request.QueryType != pb.QueryType_SELECT {
		cluster.Broadcast(request.Sql)
	}
	resp := &pb.QueryResponse{
		Code: 200,
		Data: res,
	}
	if request.QueryType == pb.QueryType_SELECT {
		resp.DataType = pb.DataType_MAP
	} else {
		resp.DataType = pb.DataType_INT
	}
	return resp, nil
}
