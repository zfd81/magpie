package command

import (
	"time"

	pb "github.com/zfd81/magpie/proto/magpiepb"

	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"
)

type GlobalFlags struct {
	Endpoints []string
	User      string
	Password  string
}

var (
	storageClient pb.StorageClient
	magpieClient  pb.MagpieClient
)

func GetConnection() *grpc.ClientConn {
	address := globalFlags.Endpoints[0]
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(3*time.Second))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn
}

func GetTableClient() pb.MetaClient {
	return pb.NewMetaClient(GetConnection())
}

func GetClusterClient() pb.ClusterClient {
	return pb.NewClusterClient(GetConnection())
}
