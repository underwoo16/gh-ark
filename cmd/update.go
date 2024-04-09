package cmd

import (
	"fmt"
	"strings"

	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
	"github.com/underwoo16/gh-ark/gh"
	"github.com/underwoo16/gh-ark/git"
	"github.com/underwoo16/gh-ark/utils"
)

var updateCmd = &cobra.Command{
	Use:   "upr",
	Short: "Update an existing PR with the latest commit",
	Long:  `Updates a pull request on GitHub with the latest commit and squashes the latest commit into the commit associated with the existing PR`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runUpdateCmd(args)
	},
	Example: `gh ark update $commit_sha`,
}

var updatePrList bool

func init() {
	// TODO: add flags
	// TODO: add ability to target PR by number
	// TODO: add ability to target PR by branch name
	// TODO: add flag to select commit sha from list
	updateCmd.Flags().BoolVarP(&updatePrList, "list", "l", false, "Select commit from list")
	// TODO: add flag to squash commit in PR
	// TODO: if no arg or flag, prompt user to select PR to update
}

func runUpdateCmd(args []string) error {
	gitService := git.NewGitService()
	ghService := gh.NewGitHubService()

	if updatePrList {
		return updatePullRequestList(gitService, ghService)
	}

	return updatePullRequest(args[0], gitService, ghService)
}

func updatePullRequestList(gitService git.GitService, ghService gh.GitHubService) error {
	commits, _ := gitService.LogFromMainOrMaster()
	choice, _ := ghService.Prompt("Select commit to update PR with", commits[0], commits)
	commit := strings.Fields(commits[choice])[0]
	return updatePullRequest(commit, gitService, ghService)
}

func updatePullRequest(pullRequestCommit string, gitService git.GitService, ghService gh.GitHubService) error {
	trunk := gitService.CurrentBranch()
	latestCommit := gitService.LatestCommit()

	branchName := gitService.BuildBranchNameFromCommit(pullRequestCommit)
	fmt.Println("Attempting to update pull request for branch:", branchName)

	io := iostreams.System()
	io.StartProgressIndicator()
	pullRequest := ghService.GetPullRequestForBranch(branchName)

	if pullRequest == nil {
		return fmt.Errorf("no pull request found for stack: %s", branchName)
	}

	io.StopProgressIndicator()

	fmt.Printf("Updating pull request:\n%s <- %s\n%s\n", utils.Green(pullRequest.BaseRefName), utils.Yellow(pullRequest.HeadRefName), utils.Blue(pullRequest.Url))

	io.StartProgressIndicator()
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

	io.StopProgressIndicator()
	fmt.Printf("Pull request succesfully updated.\n")

	return nil
}
