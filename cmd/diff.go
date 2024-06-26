package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
	"github.com/underwoo16/gh-ark/gh"
	"github.com/underwoo16/gh-ark/git"
	"github.com/underwoo16/gh-ark/utils"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Create PR from latest commit",
	Long:  `Creates a pull request on GitHub which contains the latest commit and targets origin/master`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDiffCmd()
	},
	Example: `gh ark diff
gh ark diff -l
gh ark diff --create`,
}

var newDiffList bool
var createNew bool

func init() {
	diffCmd.Flags().BoolVarP(&newDiffList, "list", "l", false, "Select commit from list")
	diffCmd.Flags().BoolVarP(&createNew, "create", "c", false, "Create a new PR targeting the repository's default branch")
}

func runDiffCmd() error {
	gitService := git.NewGitService()
	ghService := gh.NewGitHubService()

	if newDiffList {
		return createPullRequestFromList(gitService, ghService)
	}

	return createPullRequestFromLatest(gitService, ghService)
}

func createPullRequestFromList(gitService git.GitService, ghService gh.GitHubService) error {
	commits, _ := gitService.LogFromMainOrMaster()
	choice, _ := ghService.Prompt("Select commit to create PR from", commits[0], commits)

	commit := strings.Fields(commits[choice])[0]
	return createPullRequest(commit, gitService, ghService)
}

func createPullRequestFromLatest(gitService git.GitService, ghService gh.GitHubService) error {
	latestCommit := gitService.LatestCommit()
	return createPullRequest(latestCommit, gitService, ghService)
}

func createPullRequest(commit string, gitService git.GitService, ghService gh.GitHubService) error {
	io := iostreams.System()
	io.StartProgressIndicator()

	if !createNew {
		previousCommit, err := gitService.PreviousCommit()
		if err != nil {
			log.Fatal(err)
		}

		if previousCommit != commit {
			previousBranch := gitService.BuildBranchNameFromCommit(previousCommit)

			existingPR := ghService.GetPullRequestForBranch(previousBranch)
			if existingPR != nil {
				io.StopProgressIndicator()
				fmt.Printf("Updating PR (%s) with commit %s\n", utils.Yellow(previousBranch), utils.Yellow(commit))
				return updatePullRequest(commit, previousCommit, gitService, ghService)
			}
		}
	}

	trunk := gitService.CurrentBranch()

	fmt.Printf("Creating PR from commit %s\n", utils.Yellow(commit))
	branchName := gitService.BuildBranchNameFromCommit(commit)

	if !gitService.LocalBranchExists(branchName) {
		err := gitService.CreateBranch(branchName)
		if err != nil {
			log.Fatal(err)
		}
	}

	// TODO: stash changes or use worktree?
	err := gitService.Switch(branchName)
	if err != nil {
		log.Fatal(err)
	}

	err = gitService.CherryPick(commit)

	if err != nil {
		fmt.Println("Cherry-pick failed, aborting...")
		gitService.AbortCherryPick()
		gitService.Switch(trunk)
		log.Fatal(err)
	}

	err = gitService.PushNewBranch()
	if err != nil {
		fmt.Println("Push failed...")
		log.Fatal(err)
	}

	io.StopProgressIndicator()

	// TODO: handle when user cancels - error code 2?
	err = ghService.CreatePullRequest()
	if err != nil {
		fmt.Println("PR creation failed...")
		log.Fatal(err)
	}

	err = gitService.Switch(trunk)
	if err != nil {
		fmt.Printf("Switching back to %s failed...\n", utils.Green(trunk))
		log.Fatal(err)
	}

	return nil
}
