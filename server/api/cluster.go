package api

import (
	"bytes"
	"fmt"

	"github.com/zfd81/magpie/server"

	"github.com/spf13/cast"
	pb "github.com/zfd81/magpie/proto/magpiepb"
)

type ClusterServer struct{}

func (c *ClusterServer) DataSync(request *pb.RpcRequest, stream pb.Cluster_DataSyncServer) error {
	name := request.Params["name"]
	tbl := server.GetDatabase("").GetTable(name)
	if tbl == nil {
		return fmt.Errorf("table %s does not exist", name)
	}
	tbl.FindAll(func(k string, v interface{}) {
		var str bytes.Buffer
		row := cast.ToSlice(v)
		for i, v := range row {
			if i > 0 {
				str.WriteString(",")
			}
			str.WriteString(cast.ToString(v))
		}
		stream.Send(&pb.StreamResponse{Data: str.String()})
	})
	return nil
}
