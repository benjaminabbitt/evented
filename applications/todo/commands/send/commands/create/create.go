package create

import (
	"fmt"
	actx2 "github.com/benjaminabbitt/evented/applications/todo/actx"
	"github.com/dsnet/try"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use: "create",
	Run: func(cmd *cobra.Command, args []string) {
		actx := cmd.Context().(actx2.TodoSendContext)
		log := actx.Log
		log.Info("test")
		host, _ := cmd.Parent().Flags().GetString("host")
		port, _ := cmd.Parent().Flags().GetUint16("port")
		fmt.Println("in Create")
		create := CreateCommmand(try.E1(uuid.NewRandom()), "test")
		response := Send(actx, host, port, create)
		log.Info(response)
	},
}
