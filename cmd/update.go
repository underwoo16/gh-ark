package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/underwoo16/gh-diffstack/gh"
	"github.com/underwoo16/gh-diffstack/git"
)

var updateCmd = &cobra.Command{
	Use:   "upr",
	Short: "Update an existing PR with the latest commit",
	Long:  `Updates a pull request on GitHub with the latest commit and squashes the latest commit into the commit associated with the existing PR`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runUpdateCmd(args)
	},
	Example: `gh-diffstack update $commit_sha`,
}

func init() {
	// TODO: add flags
	// TODO: add ability to target PR by number
	// TODO: add ability to target PR by branch name
	// TODO: add flag to select commit sha from list
	// TODO: add flag to squash commit in PR
	// TODO: if no arg or flag, prompt user to select PR to update
}

func runUpdateCmd(args []string) error {
	fmt.Println("Updating PR with latest commit...")

	gitService := git.NewGitService()
	ghService := gh.NewGitHubService()
	return updatePullRequest(args, gitService, ghService)
}

func updatePullRequest(args []string, gitService git.GitService, ghService gh.GitHubService) error {
	trunk := gitService.CurrentBranch()
	latestCommit := gitService.LatestCommit()
	pullRequestCommit := args[0]

	branchName := gitService.BuildBranchNameFromCommit(pullRequestCommit)

	err := gitService.Switch(branchName)
	if err != nil {
		return err
	}

	err = gitService.CherryPick(latestCommit)
	if err != nil {
		fmt.Println("Cherry-pick failed, aborting...")
		gitService.AbortCherryPick()
		gitService.Switch(trunk)
		return err
	}

	err = gitService.Push()
	if err != nil {
		return err
	}

	err = gitService.Switch(trunk)
	if err != nil {
		return err
	}

	err = gitService.AmendCommitWithFixup(pullRequestCommit)
	if err != nil {
		return err
	}

	err = gitService.RebaseInteractiveAutosquash(pullRequestCommit)
	if err != nil {
		return err
	}

	return nil
}
