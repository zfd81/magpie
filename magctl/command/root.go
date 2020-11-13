package command

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/zfd81/magpie/errs"
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
)

func init() {
	rootCmd.PersistentFlags().StringSliceVar(&globalFlags.Endpoints, "endpoints", []string{"127.0.0.1:8143"}, "gRPC endpoints")
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
	if err := rootCmd.Execute(); err != nil {
		ExitWithError(ExitError, err)
	}
}

func Print(format string, msgs ...interface{}) {
	fmt.Printf("[INFO] %s \n", fmt.Sprintf(format, msgs...))
}

func Errorf(format string, msgs ...interface{}) {
	fmt.Printf("[ERROR] %s \n", errs.ErrorStyleFunc(fmt.Sprintf(format, msgs...)))
}
