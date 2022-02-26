package cmd

import (
	"fmt"
	"github.com/benjaminabbitt/evented/support"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	cobra.OnInitialize(initConfig)
}

var rootCmd = &cobra.Command{
	Use:   "sender",
	Short: "root command",
	Long:  "long root command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(os.Stdout, "In Command Run")
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	log := support.Log()
	defer log.Sync()

	//config := configuration.Configuration{}
	//config.Initialize(log)
}

func initConfig() {
	fmt.Fprintln(os.Stdout, "In initConfig")
}
