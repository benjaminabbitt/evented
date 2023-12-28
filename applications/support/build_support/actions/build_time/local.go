package build_time

import (
	"fmt"
	"github.com/benjaminabbitt/evented/applications/support/build_support/actions/root"
	"github.com/spf13/cobra"
	"time"
)

func init() {
	root.RootCmd.AddCommand(timeCmd)
}

var timeCmd = &cobra.Command{
	Use: "now",
	Run: func(cmd *cobra.Command, args []string) {
		now := time.Now()
		t := now.Format(time.RFC3339)
		fmt.Println(t)
	},
}
