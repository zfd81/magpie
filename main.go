package main

import (
	"fmt"
	"net"
	"time"

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
	port int
)

func init() {
	rootCmd.Flags().IntVar(&port, "port", config.GetConfig().Port, "Port to run the server")
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05", //时间格式
	})
}

func startCommandFunc(cmd *cobra.Command, args []string) {
	server.InitTables()                 //初始化表
	cluster.Register(time.Now().Unix()) // 集群注册
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterStorageServer(s, &server.StorageServer{})
	pb.RegisterTableServer(s, &server.TableServer{})
	pb.RegisterMagpieServer(s, &server.MagpieServer{})
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
