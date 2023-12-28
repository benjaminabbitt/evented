package root

import (
	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize()
}

var RootCmd = &cobra.Command{
	Use: "build-data",
}

func Execute() error {
	return RootCmd.Execute()
}
