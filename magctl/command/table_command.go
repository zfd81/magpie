package command

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/zfd81/magpie/meta"

	pb "github.com/zfd81/magpie/api/magpiepb"

	"github.com/spf13/cobra"
)

func NewTableCommand() *cobra.Command {
	ac := &cobra.Command{
		Use:   "table <subcommand>",
		Short: "Table related commands",
	}
	ac.AddCommand(newTableAddCommand())
	ac.AddCommand(newTableDeleteCommand())
	ac.AddCommand(newTableDescribeCommand())
	ac.AddCommand(newTableListCommand())
	return ac
}

func newTableAddCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "add <table definition file>",
		Short: "Adds a new table",
		Run:   tableAddCommandFunc,
	}
	return &cmd
}

func newTableDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "del <table name>",
		Short: "Deletes a table",
		Run:   tableDeleteCommandFunc,
	}
}

func newTableDescribeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "desc <table name> [file path]",
		Short: "Describes a table",
		Run:   tableDescribeCommandFunc,
	}
}

func newTableListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists all tables",
		Run:   tableListCommandFunc,
	}
}

func tableAddCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		ExitWithError(ExitBadArgs, fmt.Errorf("table add command requires table definition file as its argument"))
	}
	path := args[0]
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		Errorf("open %s: No such file", path)
		return
	}
	definition, err := ioutil.ReadFile(path)
	if err != nil {
		Errorf(err.Error())
		return
	}
	request := &pb.RpcRequest{}
	request.Data = string(definition)
	resp, err := tableClient.CreateTable(context.Background(), request)
	if err != nil {
		Errorf(err.Error())
		return
	}
	if resp.Code == 200 {
		Print(resp.Message)
	} else {
		Errorf(resp.Message)
	}
}

func tableDeleteCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		ExitWithError(ExitBadArgs, fmt.Errorf("table del command requires table name as its argument"))
	}
	request := &pb.RpcRequest{}
	request.Params = map[string]string{}
	request.Params["name"] = args[0]
	resp, err := tableClient.DeleteTable(context.Background(), request)
	if err != nil {
		Errorf(err.Error())
		return
	}
	if resp.Code == 200 {
		Print(resp.Message)
	} else {
		Errorf(resp.Message)
	}
}

func tableDescribeCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		ExitWithError(ExitBadArgs, fmt.Errorf("table desc command requires table name as its argument"))
	}
	request := &pb.RpcRequest{}
	request.Params = map[string]string{}
	name := args[0]
	request.Params["name"] = name
	resp, err := tableClient.DescribeTable(context.Background(), request)
	if err != nil {
		Errorf(err.Error())
		return
	}
	var definition string
	if resp.Data != "" {
		var tbl meta.TableInfo
		err = json.Unmarshal([]byte(resp.Data), &tbl)
		if err != nil {
			Errorf(err.Error())
			return
		}
		var str bytes.Buffer
		_ = json.Indent(&str, []byte(resp.Data), "", "  ")
		definition = str.String()
	}
	if len(args) > 1 {
		path := args[1]
		err := ioutil.WriteFile(path, []byte(definition), 0666) //写入文件(字节数组)
		if err != nil {
			Errorf("Error exporting table structure:%s", err.Error())
		} else {
			Print("Export table structure succeeded")
		}
	} else {
		fmt.Printf("Table[%s] details:\n", name)
		fmt.Println(definition)
		fmt.Println("")
	}
}

func tableListCommandFunc(cmd *cobra.Command, args []string) {
	request := &pb.RpcRequest{}
	resp, err := tableClient.ListTables(context.Background(), request)
	if err != nil {
		Errorf(err.Error())
		return
	}
	var tbls []meta.TableInfo
	err = json.Unmarshal([]byte(resp.Data), &tbls)
	if err != nil {
		Errorf(err.Error())
		return
	}
	fmt.Println("Table list:")
	fmt.Println("-------------------------------------------------------")
	fmt.Printf("%2s %30s\n", "", "Name")
	fmt.Println("-------------------------------------------------------")
	for i, tbl := range tbls {
		fmt.Printf("%2d %30s\n", i+1, tbl.Name)
	}
	fmt.Println("-------------------------------------------------------")
}
