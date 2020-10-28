package main

import (
	"fmt"
	"log"
	"net"

	"github.com/zfd81/magpie/server"

	pb "github.com/zfd81/magpie/api/magpiepb"

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
}

func startCommandFunc(cmd *cobra.Command, args []string) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMagpieServer(s, &server.Server{})
	pb.RegisterStorageServer(s, &server.StorageServer{})
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
