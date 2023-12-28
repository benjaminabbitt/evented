package build_hash

import (
	"fmt"
	"github.com/benjaminabbitt/evented/applications/support/build_support/actions/root"
	"github.com/spf13/cobra"
)
import "github.com/go-git/go-git/v5" // with go modules enabled (GO111MODULE=on or outside GOPATH)

var HumanVersion string
var GitRoot string

func init() {
	const humanVersionName = "human_version"
	const humanVersionShorthand = "v"
	versionCmd.Flags().StringVarP(&HumanVersion, humanVersionName, humanVersionShorthand, "0.0.0", "The human preferred version string.  Typically semantic version.")
	const gitRootName = "git_root"
	const gitRootShorthand = "r"
	versionCmd.Flags().StringVarP(&GitRoot, gitRootName, gitRootShorthand, "./", "The path of the git repository root")
	root.RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "hashed_version",
	Short: "Generates a hashed version with the human readable provided version string and git git short hash or a dirty marker",
	Long:  `Sends an evented event to the location and with the data specified`,
	Run: func(cmd *cobra.Command, args []string) {
		r, err := git.PlainOpenWithOptions(GitRoot, &git.PlainOpenOptions{DetectDotGit: true})
		workTree, err := r.Worktree()
		status, err := workTree.Status()

		if err != nil {
			fmt.Println(err)
		}

		if !status.IsClean() {
			fmt.Println("Dirty Code" + status.String())
			fmt.Println(fmt.Sprintf("%s-%s", HumanVersion, "dirty"))
		} else {
			head, _ := r.Head()
			fmt.Println(fmt.Sprintf("%s-%s", HumanVersion, head.Hash().String()[0:7]))

		}
	},
}
