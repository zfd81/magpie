package command

import (
	"bufio"
	"context"
	"fmt"
	"os"

	pb "github.com/zfd81/magpie/api/magpiepb"

	"github.com/spf13/cobra"
)

func NewLoadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "load <table name> <data file>",
		Short: "Load data file",
		Run:   loadCommandFunc,
	}
	return cmd
}

func loadCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		ExitWithError(ExitBadArgs, fmt.Errorf("load command requires table name and data file as its argument"))
	}
	stream, err := magpieClient.Load(context.Background())
	if err != nil {
		Errorf(err.Error())
		return
	}
	name := args[0]
	path := args[1]
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		Errorf("open %s: No such file", path)
		return
	}
	file, err := os.Open(path)
	if err != nil {
		Errorf("Read file %s failed:", path, err.Error())
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	err = stream.Send(&pb.StreamRequest{Data: name})
	if err != nil {
		Errorf(err.Error())
		return
	}
	cnt := 0
	for scanner.Scan() {
		cnt++
		err = stream.Send(&pb.StreamRequest{Data: scanner.Text()})
		if err != nil {
			Errorf(err.Error())
			return
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		Errorf(err.Error())
		return
	}
	fmt.Println(resp.Data)
}
