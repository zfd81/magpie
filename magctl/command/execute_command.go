package command

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	pb "github.com/zfd81/magpie/proto/magpiepb"

	"github.com/spf13/cobra"
)

type Condition map[string]interface{}

func (c *Condition) String() string {
	bytes, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return string(bytes)
}

type Stmt struct {
	Name       string
	Items      []string
	Conditions Condition
}

func NewExecuteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exec <sql>",
		Short: "Execute SQL",
		Run:   executeCommandFunc,
	}
	return cmd
}

func executeCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		ExitWithError(ExitBadArgs, fmt.Errorf("execute command requires sql as its argument"))
	}
	startTime := time.Now() //计算当前时间
	sql := strings.TrimSpace(strings.Join(args, " "))
	Print(sql)
	request := &pb.Request{
		Sql: sql,
	}
	isUpdate := true
	sql = strings.ToUpper(sql)
	if strings.HasPrefix(sql, "SELECT") {
		isUpdate = false
		request.QueryType = pb.QueryType_SELECT
	} else if strings.HasPrefix(sql, "INSERT") {
		request.QueryType = pb.QueryType_INSERT
	} else if strings.HasPrefix(sql, "DELETE") {
		request.QueryType = pb.QueryType_DELETE
	} else if strings.HasPrefix(sql, "UPDATE") {
		request.QueryType = pb.QueryType_UPDATE
	} else {
		Errorf("syntax error: %s", request.Sql)
		return
	}
	conn := GetConnection()
	defer conn.Close()
	magpieClient = pb.NewMagpieClient(conn)
	var resp *pb.Response
	var err error
	if isUpdate {
		resp, err = magpieClient.Update(context.Background(), request)
	} else {
		resp, err = magpieClient.Query(context.Background(), request)
	}
	if err != nil {
		Errorf(err.Error())
		return
	}
	Print(string(resp.Data))
	Print("time cost: %v\n", time.Since(startTime))
}
