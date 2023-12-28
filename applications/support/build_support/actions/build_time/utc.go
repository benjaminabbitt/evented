package build_time

import (
	"fmt"
	"github.com/benjaminabbitt/evented/applications/support/build_support/actions/root"
	"github.com/spf13/cobra"
	"time"
)

func init() {
	root.RootCmd.AddCommand(utcCmd)
}

var utcCmd = &cobra.Command{
	Use: "utcNow",
	Run: func(cmd *cobra.Command, args []string) {
		now := time.Now().UTC()
		t := now.Format(time.RFC3339)
		fmt.Println(t)
	},
}
