package build_hash

import (
	"fmt"
	"github.com/benjaminabbitt/evented/applications/support/build_support/actions/root"
	"github.com/spf13/cobra"
)
import "github.com/go-git/go-git/v5" // with go modules enabled (GO111MODULE=on or outside GOPATH)

var Human_version string
var Git_root string

func init() {
	const human_version_name = "human_version"
	const human_version_shorthand = "v"
	versionCmd.Flags().StringVarP(&Human_version, human_version_name, human_version_shorthand, "0.0.0", "The human preferred version string.  Typically semantic version.")
	const git_root_name = "git_root"
	const git_root_shorthand = "r"
	versionCmd.Flags().StringVarP(&Git_root, git_root_name, git_root_shorthand, "./", "The path of the git repository root")
	root.RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "hashed_version",
	Short: "Generates a hashed version with the human readable provided version string and git git short hash or a dirty marker",
	Long:  `Sends an evented event to the location and with the data specified`,
	Run: func(cmd *cobra.Command, args []string) {
		r, err := git.PlainOpenWithOptions(Git_root, &git.PlainOpenOptions{DetectDotGit: true})
		workTree, err := r.Worktree()
		status, err := workTree.Status()

		if err != nil {
			fmt.Println(err)
		}

		if !status.IsClean() {
			fmt.Println(fmt.Sprintf("%s-%s", Human_version, "dirty"))
		} else {
			head, _ := r.Head()
			fmt.Println(fmt.Sprintf("%s-%s", Human_version, head.Hash().String()[0:7]))

		}
	},
}
