package send

import (
	"github.com/benjaminabbitt/evented/applications/todo/commands/send/commands/create"
	"github.com/spf13/cobra"
)

func init() {
	Cmd.Flags().String("host", "localhost", "The host with which to connect")
	Cmd.Flags().Int("port", 1737, "The port on which to connect")

	Cmd.AddCommand(create.Cmd)
}

var Cmd = &cobra.Command{
	Use:   "send",
	Short: "Sends an evented command",
	Long:  `Sends an evented command to the location and with the data specified`,
}
