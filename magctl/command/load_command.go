package command

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	pb3 "github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
	pb "github.com/zfd81/magpie/proto/magpiepb"
)

const (
	OneK = 1024
	OneM = 1024 * OneK
	OneG = 1024 * OneM
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
	conn := GetConnection()
	defer conn.Close()
	magpieClient = pb.NewMagpieClient(conn)
	stream, err := magpieClient.Load(context.Background())
	if err != nil {
		Errorf(err.Error())
		return
	}
	name := args[0]
	path := args[1]
	info, err := os.Stat(path)
	fileSize := info.Size()
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
		resp, err := stream.CloseAndRecv()
		if err != nil {
			Errorf(err.Error())
			return
		}
		Errorf(resp.Message)
		return
	}

	dataStream := make(chan string, 200)
	count := 100
	cnt := 0
	bar := pb3.Simple.Start(count)
	go func() {
		for i := 0; i < count-1; i++ {
			bar.Increment()
			cnt++
			time.Sleep(time.Duration(fileSize*10/(7*OneM)) * time.Millisecond)
		}
	}()
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for scanner.Scan() {
			dataStream <- scanner.Text()
		}
		close(dataStream)
	}()

	for data := range dataStream {
		err = stream.Send(&pb.StreamRequest{Data: data})
		if err != nil {
			resp, err := stream.CloseAndRecv()
			if err != nil {
				Errorf(err.Error())
				return
			}
			Errorf(resp.Message)
			return
		}
	}

	wg.Wait()
	resp, err := stream.CloseAndRecv()
	if err != nil {
		Errorf(err.Error())
		return
	}
	if resp.Code != 200 {
		Errorf(resp.Message)
	} else {
		bar.Add(count - cnt)
		bar.Finish()
		Print("Start time: %s", resp.StartTime)
		Print("End time: %s", resp.EndTime)
		Print("Elapsed time: %v", time.Duration(resp.ElapsedTime))
		Print("Record Count: %d", resp.RecordCount)
		Print(resp.Message)
	}
}
