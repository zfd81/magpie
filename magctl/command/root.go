package command

import (
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "github.com/zfd81/magpie/proto/magpiepb"

	"github.com/spf13/cobra"
	"github.com/zfd81/magpie/errors"
)

const (
	Version        = "1.0"
	cliName        = "magctl"
	cliDescription = "A simple command line client for magpie."

	defaultDialTimeout      = 2 * time.Second
	defaultCommandTimeOut   = 5 * time.Second
	defaultKeepAliveTime    = 2 * time.Second
	defaultKeepAliveTimeOut = 6 * time.Second
)

var (
	globalFlags = GlobalFlags{}
	rootCmd     = &cobra.Command{
		Use:        cliName,
		Short:      cliDescription,
		SuggestFor: []string{"rockctl"},
	}
	conn          *grpc.ClientConn
	storageClient pb.StorageClient
	tableClient   pb.TableClient
	magpieClient  pb.MagpieClient
)

func init() {
	rootCmd.PersistentFlags().StringSliceVar(&globalFlags.Endpoints, "endpoints", []string{"127.0.0.1:8843"}, "gRPC endpoints")
	rootCmd.PersistentFlags().StringVar(&globalFlags.User, "user", "", "username[:password] for authentication (prompt if password is not supplied)")
	rootCmd.PersistentFlags().StringVar(&globalFlags.Password, "password", "", "password for authentication (if this option is used, --user option shouldn't include password)")

	rootCmd.AddCommand(
		NewVersionCommand(),
		NewTableCommand(),
		NewStoreCommand(),
		NewLoadCommand(),
		NewExecuteCommand(),
	)
}

func Execute() {
	address := globalFlags.Endpoints[0]
	// Set up a connection to the server.
	c, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	conn = c
	defer conn.Close()
	storageClient = pb.NewStorageClient(conn)
	tableClient = pb.NewTableClient(conn)
	magpieClient = pb.NewMagpieClient(conn)
	if err := rootCmd.Execute(); err != nil {
		ExitWithError(ExitError, err)
	}
}

func Print(format string, msgs ...interface{}) {
	fmt.Printf("[INFO] %s \n", fmt.Sprintf(format, msgs...))
}

func Errorf(format string, msgs ...interface{}) {
	fmt.Printf("[ERROR] %s \n", errors.ErrorStyleFunc(fmt.Sprintf(format, msgs...)))
}
