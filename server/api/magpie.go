package api

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/zfd81/magpie/store"

	"github.com/zfd81/magpie/config"

	"github.com/zfd81/magpie/sql"

	log "github.com/sirupsen/logrus"

	"github.com/zfd81/magpie/server"

	"github.com/zfd81/magpie/cluster"
	pb "github.com/zfd81/magpie/proto/magpiepb"
)

var conf = config.GetConfig()

type MagpieServer struct{}

func (s *MagpieServer) Load(stream pb.Magpie_LoadServer) error {
	startTime := time.Now()
	r, err := stream.Recv()
	name := r.Data //获得表名
	db := server.GetDatabase("")
	tbl := db.GetTable(name)
	if tbl == nil {
		return stream.SendAndClose(&pb.LoadResponse{
			Code:    400,
			Message: fmt.Sprintf("table %s does not exist", name),
		})
	}

	var mu sync.RWMutex //读写锁
	chs := make([]chan []*sql.Row, conf.StoragePoolSize)
	rowsPool := make([][]*sql.Row, conf.StoragePoolSize)
	for i := 0; i < conf.StoragePoolSize; i++ {
		chs[i] = make(chan []*sql.Row, 50)
		rowsPool[i] = []*sql.Row{}
	}

	wg := sync.WaitGroup{}
	var cnt int64 = 0

	add := func(count int64) {
		mu.Lock()
		cnt = cnt + count
		mu.Unlock()
	}

	for p := 0; p < conf.StoragePoolSize; p++ {
		wg.Add(1)
		index := p
		go func() {
			defer wg.Done()
			for rows := range chs[index] {
				kvs := make([]store.KeyValue, len(rows))
				for i, row := range rows {
					kvs[i] = row.KeyValue()
				}
				err := db.GetStorage(index).BatchPut(name, kvs)
				if err == nil {
					add(int64(len(rows)))
				}
			}
		}()
	}

	num := 1000

	for {
		r, err = stream.Recv()
		if err == io.EOF {
			for i, rows := range rowsPool {
				if len(rows) > 0 {
					chs[i] <- rows
				}
				close(chs[i])
			}
			break
		}
		if err != nil {
			for _, c := range chs {
				close(c)
			}
			return err
		}
		row := tbl.NewRow()
		row.Load(r.Data, ",")
		index := db.GetStorageIndex(row.Key())
		rows := &rowsPool[index]
		*rows = append(*rows, row)
		if len(*rows) == num {
			chs[index] <- *rows
			rowsPool[index] = []*sql.Row{}
		}
	}

	wg.Wait()
	endTime := time.Now()
	log.WithFields(log.Fields{
		"table":   name,
		"elapsed": time.Since(startTime),
	}).Info("Data loaded successfully")
	cluster.Broadcast("load table " + strings.TrimSpace(name))
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
