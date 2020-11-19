package command

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cast"

	"github.com/zfd81/magpie/cluster"

	pb "github.com/zfd81/magpie/proto/magpiepb"

	"github.com/spf13/cobra"
)

func NewMemberCommand() *cobra.Command {
	ac := &cobra.Command{
		Use:   "member <subcommand>",
		Short: "Member related commands",
	}
	ac.AddCommand(newMemberStatusCommand())
	ac.AddCommand(newMemberListCommand())
	return ac
}

func newMemberStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "stat <member id>",
		Short: "View member status",
		Run:   memberStatusCommandFunc,
	}
}

func newMemberListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists all members",
		Run:   memberListCommandFunc,
	}
}

func memberStatusCommandFunc(cmd *cobra.Command, args []string) {
	//if len(args) < 1 {
	//	ExitWithError(ExitBadArgs, fmt.Errorf("member desc command requires member name as its argument"))
	//}
	request := &pb.RpcRequest{}
	resp, err := GetClusterClient().MemberStatus(context.Background(), request)
	if err != nil {
		Errorf(err.Error())
		return
	}
	var tblInfos []map[string]interface{}
	err = json.Unmarshal([]byte(resp.Data), &tblInfos)
	if err != nil {
		Errorf(err.Error())
		return
	}
	fmt.Println("+--------------------+--------------+-----------+-----------+")
	fmt.Printf("%1s %18s %1s %12s %1s %9s %1s %9s %1s\n", "|", "TABLE NAME    ", "|", "COLUMN COUNT", "|", "ROW COUNT", "|", "SIZE(KB)", "|")
	fmt.Println("+--------------------+--------------+-----------+-----------+")
	for _, n := range tblInfos {
		fmt.Printf("%1s %18s %1s %12s %1s %9s %1s %9s %1s\n", "|", n["name"], "|", cast.ToString(n["colCount"]), "|", cast.ToString(n["rowCount"]), "|", cast.ToString(n["tblSize"]), "|")
	}
	fmt.Println("+--------------------+--------------+-----------+-----------+")
}

func memberListCommandFunc(cmd *cobra.Command, args []string) {
	request := &pb.RpcRequest{}
	resp, err := GetClusterClient().ListMembers(context.Background(), request)
	if err != nil {
		Errorf(err.Error())
		return
	}
	var nodes []*cluster.Node
	err = json.Unmarshal([]byte(resp.Data), &nodes)
	if err != nil {
		Errorf(err.Error())
		return
	}
	fmt.Println("+----------------------+------------------+------------+-----------+---------------------+")
	fmt.Printf("%1s %20s %1s %16s %1s %10s %1s %9s %1s %19s %1s\n", "|", "ENDPOINT      ", "|", "ID       ", "|", "TEAM   ", "|", "IS LEADER", "|", "START-UP TIME   ", "|")
	fmt.Println("+----------------------+------------------+------------+-----------+---------------------+")
	for _, n := range nodes {
		fmt.Printf("%1s %20s %1s %16s %1s %10s %1s %9t %1s %19s %1s\n", "|", fmt.Sprintf("%s:%d", n.Address, n.Port), "|", n.Id, "|", n.Team, "|", n.LeaderFlag, "|", time.Unix(int64(n.StartUpTime), 0).Format("2006-01-02 15:04:05"), "|")
	}
	fmt.Println("+----------------------+------------------+------------+-----------+---------------------+")
}
