package cmd

import (
	"fmt"
	"github.com/benjaminabbitt/evented/applications/support/sender/actions/root"
	"github.com/spf13/cobra"
)

func init() {
	sendCmd.Flags().String("host", "localhost", "The host with which to connect")
	sendCmd.Flags().Int("port", 1737, "The port on which to connect")
	root.RootCmd.AddCommand(sendCmd)
}

var sendCmd = &cobra.Command{
	Use:   "command",
	Short: "Sends an evented command",
	Long:  `Sends an evented command to the location and with the data specified`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO")
	},
}
