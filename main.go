package main

import (
	"fmt"
	"net"
	"time"

	"github.com/zfd81/magpie/etcd"

	"github.com/zfd81/magpie/server/api"

	log "github.com/sirupsen/logrus"

	"github.com/zfd81/magpie/cluster"
	"github.com/zfd81/magpie/server"

	pb "github.com/zfd81/magpie/proto/magpiepb"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zfd81/magpie/config"
	cmd "github.com/zfd81/magpie/magctl/command"
	"google.golang.org/grpc"
)

var (
	rootCmd = &cobra.Command{
		Use:        "parrot",
		Short:      "parrot server",
		SuggestFor: []string{"parrot"},
		Run:        startCommandFunc,
	}
	port int64
)

func init() {
	rootCmd.Flags().Int64Var(&port, "port", config.GetConfig().Port, "Port to run the server")
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05", //时间格式
	})
	//log.SetLevel(log.InfoLevel)
}

func startCommandFunc(cmd *cobra.Command, args []string) {
	config.GetConfig().Port = port
	err := server.OpenLogStorage("magpie-log.db")
	if err != nil {
		err = server.OpenLogStorage(fmt.Sprintf("magpie-log-%d.db", port))
		if err != nil {
			log.Panic(err)
		}
	}
	err = etcd.Connect() //连接etcd服务器
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Panic("Etcd connection error")
	}
	err = server.InitMetadata() //初始化表
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Panic("Metadata initialization failed")
	}
	cluster.WatchLeader()                     //监听集群Leader节点变化
	cluster.WatchMembers()                    //监听集群Members节点变化
	err = cluster.Register(time.Now().Unix()) // 集群注册
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Panic("Cluster registration failed")
	}
	err = cluster.InitMembers() //初始化集群成员信息
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Initialization member information error")
	}
	err = cluster.DataSync() //和同一团队的leader进行数据同步
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Panic("Data synchronization error")
	}
	//schedule.StartScheduler() //启动计划程序

	time.Sleep(time.Duration(2) * time.Second)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterStorageServer(s, &api.StorageServer{})
	pb.RegisterMetaServer(s, &api.MetaServer{})
	pb.RegisterLogServer(s, &api.LogServer{})
	pb.RegisterMagpieServer(s, &api.MagpieServer{})
	pb.RegisterClusterServer(s, &api.ClusterServer{})
	log.Infof("Magpie server listening on: %d", config.GetConfig().Port)
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		cmd.ExitWithError(cmd.ExitError, err)
	}
}

func main() {
	color.Green(config.GetConfig().Banner)
	Execute()
}
