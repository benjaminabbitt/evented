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
		git_root := args[0]
		human_version := args[1]
		r, err := git.PlainOpenWithOptions(git_root, &git.PlainOpenOptions{DetectDotGit: true})
		workTree, err := r.Worktree()
		status, err := workTree.Status()

		if err != nil {
			fmt.Println(err)
		}

		//fmt.Print(status)
		if !status.IsClean() {
			fmt.Println(fmt.Sprintf("%s-%s", human_version, "dirty"))
		} else {
			head, _ := r.Head()
			fmt.Println(fmt.Sprintf("%s-%s", human_version, head.Hash().String()[0:7]))

		}
	},
}
