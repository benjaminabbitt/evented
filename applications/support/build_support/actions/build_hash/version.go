package build_hash

import (
	"fmt"
	"github.com/benjaminabbitt/evented/applications/support/build_support/actions/root"
	"github.com/spf13/cobra"
)
import "github.com/go-git/go-git/v5" // with go modules enabled (GO111MODULE=on or outside GOPATH)

func init() {
	root.RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use: "hashed_version",
	Run: func(cmd *cobra.Command, args []string) {
		r, _ := git.PlainOpen(args[0])
		workTree, _ := r.Worktree()
		status, _ := workTree.Status()
		//fmt.Print(status)
		if !status.IsClean() {
			fmt.Println(fmt.Sprintf("%s-%s", args[1], "dirty"))
		} else {
			head, _ := r.Head()
			fmt.Println(fmt.Sprintf("%s-%s", args[1], head.Hash().String()[0:7]))
		}
	},
}
