package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/underwoo16/gh-diffstack/gh"
	"github.com/underwoo16/gh-diffstack/git"
	"github.com/underwoo16/gh-diffstack/utils"
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

var updateList bool

func init() {
	// TODO: add flags
	// TODO: add ability to target PR by number
	// TODO: add ability to target PR by branch name
	// TODO: add flag to select commit sha from list
	updateCmd.Flags().BoolVarP(&updateList, "list", "l", false, "Select commit from list")
	// TODO: add flag to squash commit in PR
	// TODO: if no arg or flag, prompt user to select PR to update
}

func runUpdateCmd(args []string) error {
	gitService := git.NewGitService()
	ghService := gh.NewGitHubService()
	return updatePullRequest(args, gitService, ghService)
}

func updatePullRequest(args []string, gitService git.GitService, ghService gh.GitHubService) error {
	trunk := gitService.CurrentBranch()
	latestCommit := gitService.LatestCommit()
	pullRequestCommit := args[0]

	branchName := gitService.BuildBranchNameFromCommit(pullRequestCommit)

	pullRequest := ghService.GetPullRequestForBranch(branchName)

	if pullRequest == nil {
		return fmt.Errorf("no pull request found for stack%s", branchName)
	}

	fmt.Printf("Updating pull request:\n%s <- %s\n%s\n", utils.Green(pullRequest.BaseRefName), utils.Yellow(pullRequest.HeadRefName), utils.Blue(pullRequest.Url))

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

	fmt.Printf("Pull request succesfully updated.\n")

	return nil
}
