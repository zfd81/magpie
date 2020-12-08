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
	port int32
)

func init() {
	rootCmd.Flags().Int32Var(&port, "port", config.GetConfig().Port, "Port to run the server")
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05", //时间格式
	})
	//log.SetLevel(log.InfoLevel)
}

func startCommandFunc(cmd *cobra.Command, args []string) {
	config.GetConfig().Port = port

	//打印配置信息
	log.Info("Magpie version: ", config.GetConfig().Version)
	log.Info("Magpie data directory: ", config.GetConfig().DataDirectory)
	log.Info("Magpie data buffer size: ", config.GetConfig().BufferSize)
	log.Info("Magpie data write batch size: ", config.GetConfig().WriteBatchSize)
	log.Info("Magpie data storage pool size: ", config.GetConfig().StoragePoolSize)
	log.Info("Magpie etcd endpoints: ", config.GetConfig().Etcd.Endpoints)

	//打开数据存储库
	err := server.InitStorage()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Fatalln("Opening database file error")
	}

	//打开日志库
	err = server.OpenLogStorage("magpie-log.db")
	if err != nil {
		err = server.OpenLogStorage(fmt.Sprintf("magpie-log-%d.db", port))
		if err != nil {
			log.Fatalln(err)
		}
	}

	//连接etcd服务器
	err = etcd.Connect()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Fatalln("Etcd connection error")
	}

	//开启RPC服务
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

	//初始化表
	err = server.InitMetadata()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Fatalln("Metadata initialization failed")
	}

	//监听集群Leader节点变化
	cluster.WatchLeader()

	//监听集群Members节点变化
	cluster.WatchMembers()

	time.Sleep(time.Duration(2) * time.Second)

	// 集群注册
	err = cluster.Register(time.Now().Unix())
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Fatalln("Cluster registration failed")
	}

	time.Sleep(time.Duration(1) * time.Second)

	//初始化集群成员信息
	err = cluster.InitMembers()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatalln("Initialization member information error")
	}

	//和同一团队的leader进行数据同步
	err = cluster.DataSync()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Fatalln("Data synchronization error")
	}

	//schedule.StartScheduler() //启动计划程序

	log.Infof("Magpie server started successfully, listening on: %d", config.GetConfig().Port)

	//开启RPC服务
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
