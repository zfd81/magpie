package command

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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

	request := &pb.QueryRequest{}
	sql := strings.TrimSpace(strings.Join(args, " "))
	request.Sql = sql
	sql = strings.ToUpper(sql)
	if strings.HasPrefix(sql, "SELECT") {
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
	resp, err := magpieClient.Execute(context.Background(), request)
	if err != nil {
		Errorf(err.Error())
		return
	}
	fmt.Println(resp.Data)
}
