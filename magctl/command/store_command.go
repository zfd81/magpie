package command

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cast"

	pb "github.com/zfd81/magpie/proto/magpiepb"

	"github.com/spf13/cobra"
)

func NewStoreCommand() *cobra.Command {
	ac := &cobra.Command{
		Use:   "store <subcommand>",
		Short: "Store related commands",
	}
	ac.AddCommand(newStoreGetCommand())
	return ac
}

var (
	getPrefix    bool
	getKeysOnly  bool
	getCountOnly bool
)

func newStoreGetCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "get <key>",
		Short: "Gets the key from store",
		Run:   storeGetCommandFunc,
	}
	cmd.Flags().BoolVar(&getPrefix, "prefix", false, "Get keys with matching prefix")
	cmd.Flags().BoolVar(&getKeysOnly, "keys-only", false, "Get only the keys")
	cmd.Flags().BoolVar(&getCountOnly, "count-only", false, "Get only the count")
	return &cmd
}

func storeGetCommandFunc(cmd *cobra.Command, args []string) {
	request := &pb.RpcRequest{}
	if getCountOnly && len(args) < 1 {
		resp, err := storageClient.Count(context.Background(), request)
		if err != nil {
			Errorf(err.Error())
			return
		}
		fmt.Println(resp.Data)
		return
	}
	if len(args) < 1 {
		ExitWithError(ExitBadArgs, fmt.Errorf("store get command needs one argument as key"))
	}
	request.Params = map[string]string{}
	request.Params["key"] = args[0]
	request.Params["prefix"] = cast.ToString(getPrefix)
	resp, err := storageClient.Get(context.Background(), request)
	if err != nil {
		Errorf(err.Error())
		return
	}
	var kvs map[string]string
	err = json.Unmarshal([]byte(resp.Data), &kvs)
	if err != nil {
		Errorf(err.Error())
		return
	}
	if getCountOnly {
		fmt.Println(len(kvs))
	} else {
		for k, v := range kvs {
			fmt.Println(k)
			fmt.Println(v)
			fmt.Println("")
		}
	}
}
