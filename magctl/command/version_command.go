package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints the version of magpie",
		Run:   versionCommandFunc,
	}
}

func versionCommandFunc(cmd *cobra.Command, args []string) {
	fmt.Println("magpie version:", Version)
}
